#include <sys/wait.h>
#include <stdio.h>
#include <unistd.h>
int main()
{
	pid_t pid;
	/* fork a child process */
	pid = fork();
	if (pid < 0) { /* error occurred */
		fprintf(stderr, "Fork Failed");
		return 1;
	}	
	else if (pid == 0) { /* child process */
		printf("Child: my pid = %d, parent pid = %d\n", getpid(), getppid());//execlp("/bin/ls","ls",NULL);
		sleep(5); //mette il processo in attesa passiva (no consumo risorse) per 5 secondi
        //_exit(0); //termina il processo figlio senza eseguire le operazioni di pulizia del processo (come la chiusura dei file aperti)
	}
	else { /* parent process */
		/* parent will wait for the child to complete */
		wait(NULL);
		printf("Parent: my PID = %d, parent pid = %d. Pid del child %d\n",getpid(),getppid(),pid );
        //sleep(10);
	}
    sleep(10);
	return 0;
}