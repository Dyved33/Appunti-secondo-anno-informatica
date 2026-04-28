#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/wait.h>

#define		READ_END		0
#define		WRITE_END		1

/***********************************************************************************/
// Programma che simula il comportamento di bash quando l'utente richiede
// l'esecuzione della linea di comando:
//		cat /etc/passwd | sort       			oppure
//		cat /etc/passwd | sort > pippo.out		(vedi oltre)
//
// Se volete verificare che bash effettivamente usi queste funzioni potete:
// - aprire un secondo terminale
// - eseguire quest'altra linea di comando:
//		strace -f -e clone,dup2,pipe2,openat,execve -p <pid bash>
//			dove <pid bash> è il PID della shell da cui 'cat /etc/passwd | sort'
//			viene eseguito
/***********************************************************************************/
int main(void) {
	pid_t pid_cat, pid_sort;
	int pipefd[2], iRet = EXIT_SUCCESS;

	if(pipe(pipefd) == 0) {
		pid_cat = fork();
		if(pid_cat == 0) {
			// Processo "cat"
			char *sArgs[] = { "cat", "/etc/passwd", NULL };
			close(pipefd[READ_END]);
			dup2(pipefd[WRITE_END], STDOUT_FILENO);
			/************************************************************************/
			// Decommentate le righe seguenti per far riscontrare a cat un errore di
			// "File non trovato" e reindirizzare anche standard error nella pipe
			/*****/
			//sArgs[1] = "/etc/passwdd";
			//dup2(pipefd[WRITE_END], STDERR_FILENO);
			/************************************************************************/
			execve("/usr/bin/cat", sArgs, NULL);
			perror("Errore in execve() cat");
			return EXIT_FAILURE;
		}
		else {
			pid_sort = fork();
			if(pid_sort == 0) {
				// Processo "sort"
				char *sArgs[] = { "sort", NULL };
				close(pipefd[WRITE_END]);
				dup2(pipefd[READ_END], STDIN_FILENO);
				/************************************************************************/
				// Decommentate le righe seguenti per redirigere l'outout di "sort"
				// nel file "pippo.out"
				/*****/
				//int fd = open("pippo.out", O_WRONLY|O_TRUNC|O_CREAT, S_IRUSR|S_IWUSR);
				//dup2(fd, STDOUT_FILENO);
				/************************************************************************/
				execve("/usr/bin/sort", sArgs, NULL);
				perror("Errore in execve() sort");
				return EXIT_FAILURE;
			}
			else {
				// Parent
				int wstatus;
				close(pipefd[READ_END]);
				close(pipefd[WRITE_END]);
				waitpid(pid_cat, &wstatus, 0);
				printf("********************************** Cat  terminato, codice di ritorno %d ****************************\n", WEXITSTATUS(wstatus));
				waitpid(pid_sort, &wstatus, 0);
				printf("********************************** Sort terminato, codice di ritorno %d ****************************\n", WEXITSTATUS(wstatus));
			}
		}
	}
	else {
		perror("Errore in creazione pipe");
		iRet = EXIT_FAILURE;
	}

	return iRet;
}
