#include <stdio.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <semaphore.h>
#include <errno.h>

int main()
{
    sem_t *sem = sem_open("/my_sem", O_CREAT | O_RDWR | O_EXCL, S_IRUSR | S_IWUSR, 0);
    if (sem == SEM_FAILED)
    {
        if (errno == EEXIST)
        {
            printf("Semaphore already exists.\n");
        }
        else
        {
            perror("sem_open");
        }
    }
    else
    {
        sem_close(sem);
    }
    return 0;
}