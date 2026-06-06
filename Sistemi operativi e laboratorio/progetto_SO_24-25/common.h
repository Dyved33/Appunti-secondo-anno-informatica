/*
 * common.h — Definizioni condivise tra tutti i processi dell'applicazione "Incrocio"
 *
 * Contiene:
 *   - Costanti per i nomi delle risorse IPC (shared memory e semafori POSIX)
 *   - La struttura SharedData, mappata in memoria condivisa
 *
 * IPC utilizzati:
 *   - POSIX Shared Memory  : scambio di PID e direzioni tra i processi
 *   - POSIX Named Semaphore: sincronizzazione tra garage, incrocio e automobili
 */

#ifndef COMMON_H
#define COMMON_H

#include <sys/types.h>   /* pid_t */

/* ── Numero di automobili/strade ─────────────────────────────────── */
#define NUM_CARS 4

/* ── Nomi risorse IPC POSIX (visibili a tutti i processi per nome) ── */
#define SHM_NAME          "/incrocio_shm"       /* shared memory             */
#define SEM_READY         "/sem_incrocio_ready" /* garage → incrocio: auto pronte */
#define SEM_CROSSED       "/sem_incrocio_cross" /* auto → incrocio: ho attraversato */
#define SEM_GARAGE_EXIT   "/sem_garage_exit"    /* garage → incrocio: sto uscendo */
#define SEM_GO_FMT        "/sem_go_%d"          /* incrocio → auto i: puoi passare */

/* ── File di log ─────────────────────────────────────────────────── */
#define INCROCIO_TXT "incrocio.txt"
#define AUTO_TXT     "auto.txt"

/*
 * SharedData — struttura in memoria condivisa
 *
 * Garage scrive:  garage_pid, car_pids[], directions[]
 * Incrocio scrive: incrocio_pid, terminate
 * Lettura da tutti i processi.
 *
 * Accesso serializzato tramite semafori: non serve mutex aggiuntivo
 * perché garage scrive directions[] PRIMA di postare SEM_READY,
 * e incrocio legge DOPO aver consumato SEM_READY.
 */
typedef struct {
    pid_t incrocio_pid;              /* PID del processo incrocio           */
    pid_t garage_pid;                /* PID del processo garage             */
    pid_t car_pids[NUM_CARS];        /* PID delle automobili correnti       */
    int   directions[NUM_CARS];      /* strada di svolta di ogni auto       */
    volatile int terminate;          /* flag di terminazione globale        */
} SharedData;

#endif /* COMMON_H */
