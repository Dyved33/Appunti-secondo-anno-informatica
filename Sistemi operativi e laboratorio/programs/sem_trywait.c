#include <stdio.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <semaphore.h>
#include <errno.h>
#include <linux/time.h>

int main()
{
    sem_t *sem = sem_open("/my_sem", O_RDWR);
    if (sem == SEM_FAILED)
    {
        perror("sem_open");
    }
    else
    {
        struct timespec ts;
        clock_gettime(CLOCK_REALTIME, &ts);
        ts.tv_sec += 5; // Attende al massimo 5 secondi

        if (sem_timedwait(sem, &ts) == 0)
        {
            printf("Semaforo verde\n");
        }
        else
        {
            if (errno == ETIMEDOUT)
            {
                printf("Semaforo %s rosso\n", "/my_sem");
            }
            else
            {
                perror("Errore in sem_timedwait");
            }
            sem_close(sem);
        }
    }
    return 0;
}