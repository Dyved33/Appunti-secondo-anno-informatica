#include <sys/wait.h>
#include <stdio.h>
#include <unistd.h>
#include <sys/stat.h>
#include <fcntl.h>

int main()
{
	pid_t pid;
	/* char *args[] = {"ls", "-l", NULL}; */ // Array di argomenti per cat
	/* char *args[] = {"cat", NULL} */
	char *args[] = {"go.sh", NULL}; 
	char *envs[] = {"MY_VAR=addio",NULL}; // Array di variabili d'ambiente (vuoto)
	/* fork a child process */
	pid = fork();
	if (pid < 0) { /* error occurred */
		fprintf(stderr, "Fork Failed");
		return 1;
	}	
	else if (pid == 0) { /* child process */
		/* execlp("/bin/ls","ls",NULL); */
		int fd = open("etc/passwd", O_RDONLY);
		printf("Child. my pid = %d\n", getpid());
		int x = execve("/usr/bin/ls", args, envs); 
		if (x) fprintf(stderr, "Errore:""Program failed to execute\n""|[%d]\n");
		else fprintf(stdout, "Program executed successfully\n");
		perror("Errore in execve");
	}
	else { /* parent process */
		/* parent will wait for the child to complete */
		wait(NULL);
		printf("Child Complete");
	}
	return 0;
}
