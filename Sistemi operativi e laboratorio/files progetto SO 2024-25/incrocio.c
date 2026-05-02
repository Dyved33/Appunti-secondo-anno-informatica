/*
 * incrocio.c — Processo "Incrocio"
 *
 * Responsabilità:
 *   1. Crea e inizializza tutte le risorse IPC (shared memory + semafori POSIX).
 *   2. Attende che il garage segnali la presenza di 4 automobili (sem_ready).
 *   3. Sceglie, una alla volta, l'auto autorizzata a passare (GetNextCar).
 *   4. Scrive il numero di strada dell'auto in incrocio.txt.
 *   5. Segnala all'auto scelta che può attraversare (sem_go[i]).
 *   6. Attende la conferma di avvenuto attraversamento (sem_crossed).
 *   7. Ripete i passi 3‑6 per tutte e 4 le auto, poi torna al passo 2.
 *   8. Alla ricezione di SIGTERM: avvisa tutti gli altri processi e termina per ultimo.
 *
 * IPC:
 *   Shared Memory  SHM_NAME        : condivisione di PID e direzioni
 *   Semaforo       SEM_READY       : garage → incrocio  (4 auto pronte)
 *   Semafori       SEM_GO_FMT      : incrocio → auto i  (puoi attraversare)
 *   Semaforo       SEM_CROSSED     : auto → incrocio    (ho attraversato)
 *   Semaforo       SEM_GARAGE_EXIT : garage → incrocio  (garage terminato)
 *
 * Compilazione:
 *   gcc -Wall -Wextra -o incrocio incrocio.c -lrt -lpthread
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <errno.h>
#include <signal.h>
#include <fcntl.h>
#include <time.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <semaphore.h>

#include "common.h"
#include "incrocio.h"

/* ── Risorse IPC globali (accessibili nel cleanup e nel signal handler) ── */
static SharedData *shm        = NULL;
static int         shm_fd     = -1;
static sem_t      *sem_ready  = SEM_FAILED;
static sem_t      *sem_crossed= SEM_FAILED;
static sem_t      *sem_go[NUM_CARS];
static sem_t      *sem_garage_exit = SEM_FAILED;

/* Flag volatile: modificato dal signal handler, letto dal main loop */
static volatile sig_atomic_t terminate_flag = 0;

/* ── Prototipi funzioni locali ─────────────────────────────────────────── */
static void cleanup(void);
static void handler_sigterm(int sig);
static void append_to_file(const char *filename, int road);
static int  sem_wait_interruptible(sem_t *sem);

/* ══════════════════════════════════════════════════════════════════════════
 * cleanup — rilascia e rimuove tutte le risorse IPC
 * ══════════════════════════════════════════════════════════════════════════ */
static void cleanup(void)
{
    char sem_name[32];
    int  i;

    /* Semaforo "ready" */
    if (sem_ready != SEM_FAILED) {
        sem_close(sem_ready);
        sem_unlink(SEM_READY);
    }
    /* Semaforo "crossed" */
    if (sem_crossed != SEM_FAILED) {
        sem_close(sem_crossed);
        sem_unlink(SEM_CROSSED);
    }
    /* Semaforo "garage_exit" */
    if (sem_garage_exit != SEM_FAILED) {
        sem_close(sem_garage_exit);
        sem_unlink(SEM_GARAGE_EXIT);
    }
    /* Semafori "go" per ogni auto */
    for (i = 0; i < NUM_CARS; i++) {
        if (sem_go[i] != SEM_FAILED) {
            sem_close(sem_go[i]);
            snprintf(sem_name, sizeof(sem_name), SEM_GO_FMT, i);
            sem_unlink(sem_name);
        }
    }
    /* Shared memory */
    if (shm != NULL && shm != MAP_FAILED)
        munmap(shm, sizeof(SharedData));
    if (shm_fd != -1)
        close(shm_fd);
    shm_unlink(SHM_NAME);
}

/* ══════════════════════════════════════════════════════════════════════════
 * handler_sigterm — gestisce SIGTERM: attiva il flag di terminazione
 *
 * Nota: i signal handler devono essere async-signal-safe; limitiamo le
 * operazioni a un semplice assegnamento a volatile sig_atomic_t.
 * ══════════════════════════════════════════════════════════════════════════ */
static void handler_sigterm(int sig)
{
    (void)sig;
    terminate_flag = 1;
}

/* ══════════════════════════════════════════════════════════════════════════
 * append_to_file — aggiunge "road\n" in coda al file specificato
 * ══════════════════════════════════════════════════════════════════════════ */
static void append_to_file(const char *filename, int road)
{
    FILE *f = fopen(filename, "a");
    if (f == NULL) {
        perror("fopen");
        return;
    }
    fprintf(f, "%d\n", road);
    fclose(f);
}

/* ══════════════════════════════════════════════════════════════════════════
 * sem_wait_interruptible — sem_wait che gestisce EINTR correttamente:
 *   - Riprova se l'interruzione NON è dovuta a terminazione.
 *   - Restituisce -1 (con terminate_flag == 1) se dobbiamo uscire.
 *   - Restituisce  0 al successo.
 * ══════════════════════════════════════════════════════════════════════════ */
static int sem_wait_interruptible(sem_t *sem)
{
    while (sem_wait(sem) == -1) {
        if (errno == EINTR) {
            if (terminate_flag) return -1;   /* terminazione richiesta  */
            continue;                        /* altro segnale: riprova  */
        }
        perror("sem_wait");
        return -1;
    }
    return 0;
}

/* ══════════════════════════════════════════════════════════════════════════
 * main
 * ══════════════════════════════════════════════════════════════════════════ */
int main(void)
{
    int  i;
    char sem_name[32];

    srand((unsigned)time(NULL));

    /* Inizializza array semafori go */
    for (i = 0; i < NUM_CARS; i++)
        sem_go[i] = SEM_FAILED;

    /* ── Installa handler SIGTERM ──────────────────────────────────────── */
    struct sigaction sa;
    memset(&sa, 0, sizeof(sa));
    sa.sa_handler = handler_sigterm;
    sigemptyset(&sa.sa_mask);
    /* SA_RESTART NON impostato: vogliamo che le syscall vengano interrotte */
    if (sigaction(SIGTERM, &sa, NULL) == -1) {
        perror("sigaction");
        exit(EXIT_FAILURE);
    }

    /* ── Pulizia risorse IPC da esecuzioni precedenti ──────────────────── */
    /* Ignoriamo gli errori: le risorse potrebbero non esistere */
    shm_unlink(SHM_NAME);
    sem_unlink(SEM_READY);
    sem_unlink(SEM_CROSSED);
    sem_unlink(SEM_GARAGE_EXIT);
    for (i = 0; i < NUM_CARS; i++) {
        snprintf(sem_name, sizeof(sem_name), SEM_GO_FMT, i);
        sem_unlink(sem_name);
    }

    /* ── Crea e inizializza la Shared Memory ───────────────────────────── */
    shm_fd = shm_open(SHM_NAME, O_CREAT | O_RDWR, 0666);
    if (shm_fd == -1) {
        perror("shm_open (incrocio)");
        exit(EXIT_FAILURE);
    }
    if (ftruncate(shm_fd, (off_t)sizeof(SharedData)) == -1) {
        perror("ftruncate");
        cleanup();
        exit(EXIT_FAILURE);
    }
    shm = mmap(NULL, sizeof(SharedData),
               PROT_READ | PROT_WRITE, MAP_SHARED, shm_fd, 0);
    if (shm == MAP_FAILED) {
        perror("mmap (incrocio)");
        cleanup();
        exit(EXIT_FAILURE);
    }

    /* Inizializza la struttura in shared memory */
    memset((void *)shm, 0, sizeof(SharedData));
    shm->incrocio_pid = getpid();
    shm->terminate    = 0;
    for (i = 0; i < NUM_CARS; i++) {
        shm->directions[i] = -1;   /* -1 = nessuna auto ancora creata */
        shm->car_pids[i]   = 0;
    }

    /* ── Crea i semafori POSIX named ───────────────────────────────────── */
    /*
     * Valore iniziale 0 per tutti i semafori di sincronizzazione:
     * il primo processo che fa sem_wait si bloccherà finché
     * un altro non farà sem_post.
     */
    sem_ready = sem_open(SEM_READY, O_CREAT | O_RDWR, 0666, 0);
    if (sem_ready == SEM_FAILED) {
        perror("sem_open (ready)");
        cleanup();
        exit(EXIT_FAILURE);
    }

    sem_crossed = sem_open(SEM_CROSSED, O_CREAT | O_RDWR, 0666, 0);
    if (sem_crossed == SEM_FAILED) {
        perror("sem_open (crossed)");
        cleanup();
        exit(EXIT_FAILURE);
    }

    sem_garage_exit = sem_open(SEM_GARAGE_EXIT, O_CREAT | O_RDWR, 0666, 0);
    if (sem_garage_exit == SEM_FAILED) {
        perror("sem_open (garage_exit)");
        cleanup();
        exit(EXIT_FAILURE);
    }

    for (i = 0; i < NUM_CARS; i++) {
        snprintf(sem_name, sizeof(sem_name), SEM_GO_FMT, i);
        sem_go[i] = sem_open(sem_name, O_CREAT | O_RDWR, 0666, 0);
        if (sem_go[i] == SEM_FAILED) {
            perror("sem_open (go_i)");
            cleanup();
            exit(EXIT_FAILURE);
        }
    }

    printf("[INCROCIO] Avviato (PID=%d). Risorse IPC create. "
           "In attesa del garage...\n\n", getpid());
    fflush(stdout);

    /* ════════════════════════════════════════════════════════════════════
     * Loop principale
     * ════════════════════════════════════════════════════════════════════ */
    while (!terminate_flag) {

        /* ── Passo 1: attesa passiva della notifica dal garage ─────────── */
        printf("[INCROCIO] In attesa di nuove automobili...\n");
        fflush(stdout);

        if (sem_wait_interruptible(sem_ready) == -1)
            break;   /* terminazione richiesta */

        if (terminate_flag) break;

        /* ── Passo 2‑7: gestione attraversamento delle 4 auto ──────────── */
        /*
         * Copia locale delle direzioni: incrocio lavora sulla copia
         * e marca -1 le auto già passate, senza toccare la shared memory
         * (che il garage potrebbe sovrascrivere nel ciclo successivo).
         */
        int dirs[NUM_CARS];
        for (i = 0; i < NUM_CARS; i++)
            dirs[i] = shm->directions[i];

        printf("[INCROCIO] 4 automobili all'incrocio. "
               "Gestisco il passaggio...\n");
        fflush(stdout);

        int count;
        for (count = 0; count < NUM_CARS && !terminate_flag; count++) {

            /* Passo 2: determina l'auto che può passare */
            int car = GetNextCar(dirs);
            if (car == -1) {
                fprintf(stderr, "[INCROCIO] ERRORE: GetNextCar ha restituito -1\n");
                break;
            }

            /* Passo 3: scrive in incrocio.txt il numero della strada */
            append_to_file(INCROCIO_TXT, car);

            /* Passo 4: messaggio a video */
            printf("[INCROCIO] Autorizza automobile %d "
                   "(strada %d → strada %d)\n",
                   car, car, dirs[car]);
            fflush(stdout);

            /* Passo 5: segnala all'auto che può attraversare */
            sem_post(sem_go[car]);

            /* Passo 6: attende che l'auto abbia attraversato */
            if (sem_wait_interruptible(sem_crossed) == -1)
                goto terminazione;

            /* Segna l'auto come già passata nella copia locale */
            dirs[car] = -1;

            printf("[INCROCIO] Automobile %d ha completato "
                   "l'attraversamento.\n", car);
            fflush(stdout);
        }

        printf("[INCROCIO] Tutte le automobili hanno attraversato.\n\n");
        fflush(stdout);
    }

terminazione:
    /* ════════════════════════════════════════════════════════════════════
     * Terminazione: incrocio deve essere l'ultimo a terminare
     * ════════════════════════════════════════════════════════════════════ */
    printf("[INCROCIO] SIGTERM ricevuto. Avvio sequenza di terminazione...\n");
    fflush(stdout);

    /* Imposta flag di terminazione nella shared memory */
    shm->terminate = 1;

    /* Sblocca eventuali auto in attesa sul semaforo "go" */
    for (i = 0; i < NUM_CARS; i++) {
        /* Invia SIGTERM all'auto se ancora in esecuzione */
        if (shm->car_pids[i] > 0) {
            kill(shm->car_pids[i], SIGTERM);
        }
        /* Post sul semaforo per sbloccare sem_wait nell'auto */
        sem_post(sem_go[i]);
    }

    /* Invia SIGTERM al garage e attende la sua terminazione */
    if (shm->garage_pid > 0) {
        printf("[INCROCIO] Invio SIGTERM al garage (PID=%d)...\n",
               (int)shm->garage_pid);
        fflush(stdout);
        kill(shm->garage_pid, SIGTERM);

        /* Attesa passiva con timeout (5 secondi) */
        struct timespec ts;
        clock_gettime(CLOCK_REALTIME, &ts);
        ts.tv_sec += 5;
        if (sem_timedwait(sem_garage_exit, &ts) == -1) {
            if (errno == ETIMEDOUT)
                fprintf(stderr, "[INCROCIO] Timeout: garage non ha risposto.\n");
            else if (errno != EINTR)
                perror("sem_timedwait (garage_exit)");
        } else {
            printf("[INCROCIO] Garage terminato correttamente.\n");
        }
    }

    printf("[INCROCIO] Rilascio risorse IPC e terminazione.\n");
    fflush(stdout);

    cleanup();
    return EXIT_SUCCESS;
}
