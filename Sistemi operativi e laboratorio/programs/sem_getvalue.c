#include <stdio.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <semaphore.h>
#include <errno.h>

int main()
{
    sem_t *sem = sem_open("/my_sem", O_RDONLY);
    if (sem == SEM_FAILED)
    {
        perror("sem_open");
    }
    else
    {
        int iVal;
        sem_getvalue(sem, &iVal);
        printf("Semaphore value: %s = %d\n", "/my_sem", iVal);
        sem_close(sem);
    }
    return 0;
}