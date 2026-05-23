#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <unistd.h>
#include <string.h>
#include <semaphore.h>
#include <stdbool.h>
#include <sys/wait.h>
#include <sys/mman.h>
#include <sys/stat.h>        /* For mode constants */
#include <fcntl.h>

#define MY_SHM_NAME			"/my_shm"
#define MAPPED_FILENAME		"./mapped_file"

int bFileMapping = false;

struct MyShm_t {
	sem_t sem;
	char sBuf[128];
} stMyShm, *pShm;

int main() {
	int iRet = EXIT_SUCCESS, fd;
	pid_t pid;
	
	/**************************************************************************************************
	* La variabile bFileMapping discerne tra mappare una shared memory o mappare un file in memoria.
	* Una volta chiamata shm_open() o open() il codice è lo stesso in entrambi i casi. I dati vengono
	* scambiati comunque in memoria ma in un caso transitano in una shared memory, nell'altro vengono
	* (anche) scritti nel file.
	* Nell'area di memoria mappata c'è anche un semaforo con cui si sincronizza l'accesso alla
	* mappatura stessa (vedi man 7 shm_overview).
	**************************************************************************************************/
	if(!bFileMapping)
		fd = shm_open(MY_SHM_NAME, O_RDWR|O_CREAT, S_IRUSR|S_IWUSR);
	else
		fd = open(MAPPED_FILENAME, O_RDWR|O_CREAT, S_IRUSR|S_IWUSR);
	if(fd != -1) {
		if(ftruncate(fd, sizeof(struct MyShm_t)) != -1) {
			if((pShm = mmap(NULL, sizeof(struct MyShm_t), PROT_READ|PROT_WRITE, MAP_SHARED, fd, 0)) != MAP_FAILED) {
				close(fd);
				sem_init(&(pShm->sem), 1, 0);
				pid = fork();
				if(pid == 0) {
					sem_wait(&(pShm->sem));
					printf("Child 1: ho riletto \"%s\"\n", pShm->sBuf);
					munmap(pShm, sizeof(struct MyShm_t));
					return EXIT_SUCCESS;
				}
				else {
					pid = fork();
					if(pid == 0) {
						sem_wait(&(pShm->sem));
						printf("Child 2: ho riletto \"%s\"\n", pShm->sBuf);
						munmap(pShm, sizeof(struct MyShm_t));
						return EXIT_SUCCESS;
					}
					else {
						strcpy(pShm->sBuf, "Primo contenuto della memoria");
						sem_post(&(pShm->sem));
						/*******************************************************************************/
						// Provate a commentare la seguente usleep(), le cose potrebbero non tornare....
						// Nel caso a Voi la spiegazione...
						usleep(50000);
						/*******************************************************************************/
						strcpy(pShm->sBuf, "Secondo contenuto della memoria");
						sem_post(&(pShm->sem));
						wait(NULL);
					}
				}
				
				sem_destroy(&(pShm->sem));
				munmap(pShm, sizeof(struct MyShm_t));
			}
			else {
				perror("Errore in mmap");
				iRet = EXIT_FAILURE;
			}
		}
		else {
			perror("Errore in ftruncate");
			iRet = EXIT_FAILURE;
		}
		
		if(!bFileMapping) shm_unlink(MY_SHM_NAME);
	}
	else {
		perror("Errore in apertura");
		iRet = EXIT_FAILURE;
	}
	
	return iRet;
}