/*
 * garage.c — Processo "Garage"
 *
 * Responsabilità:
 *   1. Si connette alle risorse IPC create dal processo incrocio.
 *   2. Loop infinito:
 *        a. Estrae le direzioni e crea 4 processi automobile (fork + execl).
 *        b. Stampa un messaggio per ogni auto creata.
 *        c. Notifica il processo incrocio (sem_ready).
 *        d. Attende (passivamente con waitpid) che tutte e 4 le auto terminino.
 *        e. Attende 1 secondo (sleep).
 *   3. Alla ricezione di SIGTERM: termina i figli ancora in vita,
 *      segnala il completamento all'incrocio e termina.
 *
 * I processi automobile vengono creati con fork(2) + execl(3):
 *   - fork() crea il nuovo processo
 *   - execl() sostituisce l'immagine del figlio con "./automobile"
 *   - Il numero di strada e la direzione di svolta sono passati come argv
 *
 * Compilazione:
 *   gcc -Wall -Wextra -o garage garage.c -lrt -lpthread
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
#include <sys/wait.h>
#include <semaphore.h>

#include "common.h"
#include "incrocio.h"

/* ── Risorse IPC globali ────────────────────────────────────────────────── */
static SharedData *shm           = NULL;
static int         shm_fd        = -1;
static sem_t      *sem_ready     = SEM_FAILED;
static sem_t      *sem_garage_exit = SEM_FAILED;

/* PID delle auto del ciclo corrente (usato anche nella terminazione) */
static pid_t current_car_pids[NUM_CARS];

/* Flag di terminazione */
static volatile sig_atomic_t terminate_flag = 0;

/* ── Prototipi locali ────────────────────────────────────────────────────── */
static void cleanup_garage(void);
static void handler_sigterm(int sig);
static void kill_remaining_children(void);

/* ══════════════════════════════════════════════════════════════════════════
 * cleanup_garage — chiude le risorse IPC (NON le rimuove: lo fa incrocio)
 * ══════════════════════════════════════════════════════════════════════════ */
static void cleanup_garage(void)
{
    if (sem_ready       != SEM_FAILED) sem_close(sem_ready);
    if (sem_garage_exit != SEM_FAILED) sem_close(sem_garage_exit);
    if (shm != NULL && shm != MAP_FAILED)
        munmap(shm, sizeof(SharedData));
    if (shm_fd != -1)
        close(shm_fd);
}

/* ══════════════════════════════════════════════════════════════════════════
 * handler_sigterm — imposta il flag di terminazione
 * ══════════════════════════════════════════════════════════════════════════ */
static void handler_sigterm(int sig)
{
    (void)sig;
    terminate_flag = 1;
}

/* ══════════════════════════════════════════════════════════════════════════
 * kill_remaining_children — invia SIGTERM e raccoglie i figli ancora vivi
 * ══════════════════════════════════════════════════════════════════════════ */
static void kill_remaining_children(void)
{
    int i;
    for (i = 0; i < NUM_CARS; i++) {
        if (current_car_pids[i] != 0) {
            kill(current_car_pids[i], SIGTERM);
        }
    }
    /* Attesa terminazione (waitpid di tutti i figli rimasti) */
    for (i = 0; i < NUM_CARS; i++) {
        if (current_car_pids[i] != 0) {
            waitpid(current_car_pids[i], NULL, 0);
            current_car_pids[i] = 0;
            shm->car_pids[i]    = 0;
        }
    }
}

/* ══════════════════════════════════════════════════════════════════════════
 * main
 * ══════════════════════════════════════════════════════════════════════════ */
int main(void)
{
    int i;

    srand((unsigned)(time(NULL) ^ (unsigned long)getpid()));

    /* Inizializza array PID auto */
    for (i = 0; i < NUM_CARS; i++)
        current_car_pids[i] = 0;

    /* ── Installa handler SIGTERM ──────────────────────────────────────── */
    struct sigaction sa;
    memset(&sa, 0, sizeof(sa));
    sa.sa_handler = handler_sigterm;
    sigemptyset(&sa.sa_mask);
    if (sigaction(SIGTERM, &sa, NULL) == -1) {
        perror("sigaction (garage)");
        exit(EXIT_FAILURE);
    }

    /* ── Connessione alla Shared Memory (creata da incrocio) ───────────── */
    /*
     * incrocio viene avviato per primo dallo script; attendiamo che crei
     * la shared memory prima di procedere (retry con pausa).
     */
    int retry = 0;
    while ((shm_fd = shm_open(SHM_NAME, O_RDWR, 0666)) == -1) {
        if (errno != ENOENT || ++retry > 20) {   /* max ~4 secondi di attesa */
            perror("shm_open (garage)");
            exit(EXIT_FAILURE);
        }
        usleep(200000);   /* attesa passiva 200 ms */
    }

    shm = mmap(NULL, sizeof(SharedData),
               PROT_READ | PROT_WRITE, MAP_SHARED, shm_fd, 0);
    if (shm == MAP_FAILED) {
        perror("mmap (garage)");
        cleanup_garage();
        exit(EXIT_FAILURE);
    }

    /* Registra il proprio PID nella shared memory */
    shm->garage_pid = getpid();

    /* ── Apertura semafori (già creati da incrocio) ─────────────────────── */
    sem_ready = sem_open(SEM_READY, 0);
    if (sem_ready == SEM_FAILED) {
        perror("sem_open ready (garage)");
        cleanup_garage();
        exit(EXIT_FAILURE);
    }

    sem_garage_exit = sem_open(SEM_GARAGE_EXIT, 0);
    if (sem_garage_exit == SEM_FAILED) {
        perror("sem_open garage_exit (garage)");
        cleanup_garage();
        exit(EXIT_FAILURE);
    }

    printf("[GARAGE]   Avviato (PID=%d). Connesso alle risorse IPC.\n\n",
           getpid());
    fflush(stdout);

    /* ════════════════════════════════════════════════════════════════════
     * Loop principale
     * ════════════════════════════════════════════════════════════════════ */
    while (!terminate_flag) {

        printf("[GARAGE]   Creo 4 nuove automobili...\n");
        fflush(stdout);

        /* ── Passo 1 & 2: estrai direzioni e crea le 4 auto ───────────── */
        for (i = 0; i < NUM_CARS && !terminate_flag; i++) {

            /* Estrae a sorte la direzione di svolta (≠ strada i) */
            int direction = EstraiDirezione(i);

            /* Scrive direzione e resetta PID nella shared memory */
            shm->directions[i] = direction;
            shm->car_pids[i]   = 0;

            /* Messaggio a video prima della fork */
            printf("[GARAGE]   Automobile %d creata: strada %d → strada %d\n",
                   i, i, direction);
            fflush(stdout);

            /* fork(): crea il nuovo processo */
            pid_t pid = fork();
            if (pid < 0) {
                perror("fork (garage)");
                cleanup_garage();
                exit(EXIT_FAILURE);
            }

            if (pid == 0) {
                /*
                 * ── PROCESSO FIGLIO (futuro automobile) ────────────────
                 * execl() sostituisce questa immagine con "./automobile".
                 * I parametri road e direction sono passati come stringhe argv.
                 * La shared memory e i semafori saranno aperti per nome
                 * dal processo automobile stesso.
                 */
                char road_s[4], dir_s[4];
                snprintf(road_s, sizeof(road_s), "%d", i);
                snprintf(dir_s,  sizeof(dir_s),  "%d", direction);

                execl("./automobile", "automobile", road_s, dir_s, NULL);

                /* execl() non dovrebbe mai tornare */
                perror("execl automobile");
                exit(EXIT_FAILURE);
            }

            /* ── PROCESSO PADRE (garage) ─────────────────────────────── */
            current_car_pids[i] = pid;
            shm->car_pids[i]    = pid;
        }

        if (terminate_flag) break;

        /* ── Passo 3: notifica l'incrocio che 4 auto sono pronte ───────── */
        printf("[GARAGE]   Notifico l'incrocio: 4 auto al quadrivio.\n");
        fflush(stdout);
        sem_post(sem_ready);

        /* ── Passo 4: attesa passiva della terminazione di tutti i figli── */
        /*
         * Ogni processo automobile termina dopo aver attraversato l'incrocio.
         * waitpid() è un'attesa passiva (il processo si sospende).
         * Gestisce EINTR: se arriva un segnale durante waitpid, riprova
         * (a meno che terminate_flag sia stato impostato).
         */
        for (i = 0; i < NUM_CARS; i++) {
            if (current_car_pids[i] != 0) {
                int ret;
                do {
                    ret = waitpid(current_car_pids[i], NULL, 0);
                } while (ret == -1 && errno == EINTR && !terminate_flag);

                if (ret > 0) {
                    /* Figlio terminato correttamente */
                    current_car_pids[i] = 0;
                    shm->car_pids[i]    = 0;
                }
                /* Se terminate_flag è stato impostato, usciamo dal loop
                 * con i PID rimanenti ancora in current_car_pids[]:
                 * la sezione di terminazione li gestirà. */
            }
        }

        if (terminate_flag) break;

        printf("[GARAGE]   Tutte le auto hanno attraversato. "
               "Attendo 1 secondo...\n\n");
        fflush(stdout);

        /* ── Passo 5: attesa passiva di 1 secondo ──────────────────────── */
        /*
         * sleep() è interrotto da SIGTERM (il nostro handler non usa
         * SA_RESTART), quindi la funzione ritorna anticipatamente e
         * terminate_flag viene controllato all'inizio del while.
         */
        sleep(1);
    }

    /* ════════════════════════════════════════════════════════════════════
     * Terminazione
     * ════════════════════════════════════════════════════════════════════ */
    printf("[GARAGE]   SIGTERM ricevuto. Termino i figli rimasti...\n");
    fflush(stdout);

    /* Termina eventuali auto ancora in esecuzione */
    kill_remaining_children();

    printf("[GARAGE]   Terminato. Notifico l'incrocio.\n");
    fflush(stdout);

    /* Segnala all'incrocio che il garage ha completato la terminazione */
    sem_post(sem_garage_exit);

    cleanup_garage();
    return EXIT_SUCCESS;
}
