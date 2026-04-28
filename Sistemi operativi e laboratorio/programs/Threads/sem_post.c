#include <stdio.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <semaphore.h>
#include <errno.h>

int main()
{
    sem_t *sem = sem_open("/my_sem", O_RDWR);
    if (sem == SEM_FAILED)
    {
        perror("sem_open");
    }
    else
    {
        sem_post(sem);
        sem_close(sem);
    }
    return 0;
}