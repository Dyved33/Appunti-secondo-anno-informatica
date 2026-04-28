#include <sys/wait.h>
#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>
#include <stdlib.h>

// ########################################################
// ls -l con execve()
int ls_l() {
	char *args[] = { "ls", "-l", NULL };
	char *envs[] = { NULL };
	printf("Child ls_l: pid = %d\nRunning \"ls -l\". Premere Enter...", getpid());
	fgetc(stdin);
	execve("/usr/bin/ls", args, envs);
	perror("ls_l, errore in execve()");
	exit(EXIT_FAILURE);
}

// ########################################################
// ls -l con execlp(), consulta anche la variabile di environment PATH.
// Gli argomenti sono passati come lista e non come array
int ls_l_execlp() {
	printf("Child ls_l_execlp: pid = %d\nRunning \"ls -l\" by execlp().Premere Enter...", getpid());
	fgetc(stdin);
	execlp("ls", "ls", "-l", NULL);
	perror("ls_l_execlp, errore in execlp()");
	exit(EXIT_FAILURE);
}

// ########################################################
// go.sh con execve(). Passiamo la variabile di environment "MY_VAR"
int go_sh() {
	char *args[] = { "go.sh", NULL };
	char *envs[] = { "MY_VAR=Addio" };
	printf("Child go_sh: pid = %d\nRunning \"go.sh\". Premere Enter...", getpid());
	fgetc(stdin);
	execve("./go.sh", args, envs);
	perror("go_sh, errore in execve()");
	exit(EXIT_FAILURE);
}

// ########################################################
// cat /etc/passwd
// Il child "eredita" un duplicato del file aperto nel parent;
// in go2.sh c'è la duplicazione del descrittore 3 in 0;
// cat senza argomenti legge dal descrittore 0 (standard input)
int go2_sh(int fd) {
	char *args[] = { "go2.sh", NULL };
	char *envs[] = { NULL };
	printf("go2_sh: pid = %d\nRunning \"go2.sh\" (cat /etc/passwd).Premere Enter...", getpid());
	fgetc(stdin);
	execve("./go2.sh", args, envs);
	perror("go2_sh, errore in execve()");
	exit(EXIT_FAILURE);
}

int main()
{
	pid_t pid;
	int fd;
	
	int (*fptr[])() = { &ls_l, &go_sh, &ls_l_execlp, &go2_sh };
	for(int i = 0; i < sizeof(fptr)/sizeof(fptr[0]); i++) {
		printf("########################################################\n");
		pid = fork();
		if (pid < 0) { /* error occurred */
			fprintf(stderr, "Fork Failed");
			return EXIT_FAILURE;
		}	
		if (pid == 0) { /* child process */
			if(fptr[i] != &go2_sh)
				fptr[i](); 
			else {
				fd = open("/etc/passwd", O_RDONLY);
				fptr[i](fd);
			} 
			return EXIT_FAILURE; // Non dovrebbe mai essere eseguita
		}
		else { /* parent process */
			/* parent will wait for the child to complete */
			wait(NULL);
			printf("\nParent: child %d completato\n", pid);
		}		
	}	

	close(fd);
	return EXIT_SUCCESS;
}
