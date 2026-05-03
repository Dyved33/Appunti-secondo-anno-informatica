#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <unistd.h>
#include <semaphore.h>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/stat.h>

#define MYPIPE "my_pipe"

sem_t sem;

void *LeggiPipe(void *pVoid)
{
	int fd, iCount;
	char p[64];
	printf("LeggiPipe started\n");

	if ((fd = open(MYPIPE, O_RDONLY)) != -1)
	{
		sem_post(&sem);
		iCount = read(fd, p, sizeof(p) - 1);
		p[iCount] = '\0';
		printf("\nLeggiPipe: ho letto dalla pipe \"%s\"\n", p);
		close(fd);
	}
	else
	{
		perror("LeggiPipe: errore in open");
	}

	printf("\nLeggiPipe, exiting...\n");
	pthread_exit((void *)NULL);
}

void *ScriviPipe(void *pVoid)
{
	int fd;
	printf("ScriviPipe started\n");

	if ((fd = open(MYPIPE, O_WRONLY)) != -1)
	{
		sem_wait(&sem);
		if (write(fd, "Prova", 5) == -1)
		{
			perror("ScriviPipe: errore in write");
		}
		close(fd);
	}
	else
	{
		perror("LeggiPipe: errore in open");
	}

	printf("\nScriviPipe, exiting...\n");
	pthread_exit((void *)NULL);
}

int main()
{
	pthread_t thLeggi, thScrivi;

	sem_init(&sem, 0, 0);

	mkfifo(MYPIPE, S_IRUSR | S_IWUSR);
	pthread_create(&thLeggi, NULL, LeggiPipe, NULL);
	pthread_create(&thScrivi, NULL, ScriviPipe, NULL);

	pthread_join(thLeggi, (void **)NULL);
	pthread_join(thScrivi, (void **)NULL);

	sem_destroy(&sem);

	printf("\nMain, exiting...\n");

	return 0;
}