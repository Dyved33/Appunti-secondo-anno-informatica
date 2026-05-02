/*
 * automobile.c — Processo "Automobile"
 *
 * Creato da garage con fork(2) + execl(3).
 *
 * Riceve come argomenti:
 *   argv[1] — numero della strada di provenienza (0‥3)
 *   argv[2] — numero della strada in cui svoltare (0‥3, ≠ argv[1])
 *
 * Responsabilità:
 *   1. Si connette alle risorse IPC (shared memory + semafori POSIX).
 *   2. Attende passivamente il semaforo "go" specifico per la propria strada.
 *   3. Attraversa l'incrocio (simulato con breve pausa).
 *   4. Aggiunge il numero della propria strada al file auto.txt.
 *   5. Segnala all'incrocio l'avvenuto attraversamento (sem_crossed).
 *   6. Rilascia le risorse e termina.
 *
 * Compilazione:
 *   gcc -Wall -Wextra -o automobile automobile.c -lrt -lpthread
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <errno.h>
#include <signal.h>
#include <fcntl.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <semaphore.h>

#include "common.h"

/* ── Risorse IPC locali ─────────────────────────────────────────────────── */
static SharedData *shm        = NULL;
static int         shm_fd     = -1;
static sem_t      *sem_go     = SEM_FAILED;   /* semaforo per questa auto  */
static sem_t      *sem_crossed= SEM_FAILED;   /* segnala avvenuto passaggio */

static volatile sig_atomic_t terminate_flag = 0;

/* ── Prototipi locali ────────────────────────────────────────────────────── */
static void cleanup_auto(void);
static void handler_sigterm(int sig);
static void append_to_file(const char *filename, int road);

/* ══════════════════════════════════════════════════════════════════════════
 * cleanup_auto — chiude le risorse IPC (NON le rimuove: lo fa incrocio)
 * ══════════════════════════════════════════════════════════════════════════ */
static void cleanup_auto(void)
{
    if (sem_go      != SEM_FAILED) sem_close(sem_go);
    if (sem_crossed != SEM_FAILED) sem_close(sem_crossed);
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
 * append_to_file — aggiunge "road\n" in coda al file specificato
 * ══════════════════════════════════════════════════════════════════════════ */
static void append_to_file(const char *filename, int road)
{
    FILE *f = fopen(filename, "a");
    if (f == NULL) {
        perror("fopen (automobile)");
        return;
    }
    fprintf(f, "%d\n", road);
    fclose(f);
}

/* ══════════════════════════════════════════════════════════════════════════
 * main
 * ══════════════════════════════════════════════════════════════════════════ */
int main(int argc, char *argv[])
{
    /* ── Verifica argomenti ────────────────────────────────────────────── */
    if (argc != 3) {
        fprintf(stderr, "Uso: automobile <road> <direction>\n");
        exit(EXIT_FAILURE);
    }

    int road      = atoi(argv[1]);
    int direction = atoi(argv[2]);

    if (road < 0 || road >= NUM_CARS || direction < 0 || direction >= NUM_CARS) {
        fprintf(stderr, "[AUTO ?] Parametri non validi: road=%d direction=%d\n",
                road, direction);
        exit(EXIT_FAILURE);
    }

    /* ── Installa handler SIGTERM ──────────────────────────────────────── */
    struct sigaction sa;
    memset(&sa, 0, sizeof(sa));
    sa.sa_handler = handler_sigterm;
    sigemptyset(&sa.sa_mask);
    if (sigaction(SIGTERM, &sa, NULL) == -1) {
        perror("sigaction (automobile)");
        exit(EXIT_FAILURE);
    }

    /* ── Connessione alla Shared Memory ────────────────────────────────── */
    shm_fd = shm_open(SHM_NAME, O_RDWR, 0666);
    if (shm_fd == -1) {
        perror("shm_open (automobile)");
        exit(EXIT_FAILURE);
    }

    shm = mmap(NULL, sizeof(SharedData),
               PROT_READ | PROT_WRITE, MAP_SHARED, shm_fd, 0);
    if (shm == MAP_FAILED) {
        perror("mmap (automobile)");
        close(shm_fd);
        exit(EXIT_FAILURE);
    }

    /* ── Apertura semafori ──────────────────────────────────────────────── */
    /* Semaforo "go" specifico per questa strada */
    char sem_name[32];
    snprintf(sem_name, sizeof(sem_name), SEM_GO_FMT, road);
    sem_go = sem_open(sem_name, 0);
    if (sem_go == SEM_FAILED) {
        perror("sem_open go (automobile)");
        cleanup_auto();
        exit(EXIT_FAILURE);
    }

    /* Semaforo comune "crossed" (uno per volta: incrocio serializza) */
    sem_crossed = sem_open(SEM_CROSSED, 0);
    if (sem_crossed == SEM_FAILED) {
        perror("sem_open crossed (automobile)");
        cleanup_auto();
        exit(EXIT_FAILURE);
    }

    printf("[AUTO   %d] Avviata (PID=%d, strada %d → strada %d). "
           "In attesa del via...\n",
           road, getpid(), road, direction);
    fflush(stdout);

    /* ════════════════════════════════════════════════════════════════════
     * Passo 1: attesa passiva del semaforo "go"
     *
     * sem_wait() si blocca finché il processo incrocio non esegue
     * sem_post(sem_go[road]).  Se arriva SIGTERM, sem_wait restituisce
     * -1 con errno == EINTR: in quel caso verifichiamo terminate_flag.
     * ════════════════════════════════════════════════════════════════════ */
    while (sem_wait(sem_go) == -1) {
        if (errno == EINTR) {
            if (terminate_flag) {
                printf("[AUTO   %d] Terminazione richiesta prima del via.\n",
                       road);
                fflush(stdout);
                cleanup_auto();
                return EXIT_SUCCESS;
            }
            continue;   /* altro segnale: riprova */
        }
        perror("sem_wait go (automobile)");
        cleanup_auto();
        exit(EXIT_FAILURE);
    }

    /* Controlla terminate_flag anche dopo un sem_wait riuscito
     * (incrocio potrebbe aver fatto sem_post solo per sbloccarci
     *  durante la terminazione) */
    if (terminate_flag) {
        printf("[AUTO   %d] Terminazione richiesta (sblocco da incrocio).\n",
               road);
        fflush(stdout);
        cleanup_auto();
        return EXIT_SUCCESS;
    }

    /* ════════════════════════════════════════════════════════════════════
     * Passo 2: attraversa l'incrocio
     *
     * Simulato con una breve pausa: in un sistema reale qui ci sarebbe
     * la logica di movimento del veicolo.
     * ════════════════════════════════════════════════════════════════════ */
    printf("[AUTO   %d] Sto attraversando l'incrocio "
           "(strada %d → strada %d)...\n",
           road, road, direction);
    fflush(stdout);

    usleep(150000);   /* pausa passiva da 150 ms per simulare l'attraversamento */

    /* ════════════════════════════════════════════════════════════════════
     * Passo 3: scrive il numero della propria strada in auto.txt
     *
     * Il file viene aperto in append (O_APPEND garantisce atomicità
     * se le scritture sono entro PIPE_BUF; qui usiamo fopen("a") che
     * cerca il descrittore in append mode).
     * ════════════════════════════════════════════════════════════════════ */
    append_to_file(AUTO_TXT, road);

    /* ════════════════════════════════════════════════════════════════════
     * Passo 4: segnala all'incrocio che ha attraversato e termina
     * ════════════════════════════════════════════════════════════════════ */
    sem_post(sem_crossed);

    printf("[AUTO   %d] Attraversamento completato. Esco.\n", road);
    fflush(stdout);

    cleanup_auto();
    return EXIT_SUCCESS;
}
