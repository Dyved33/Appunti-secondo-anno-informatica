#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include<stdio.h>
#include<string.h>
#include<signal.h>
#include<unistd.h>
#include<errno.h>
#include<stdlib.h>

void SigIntHandler(int iSignalCode) {
  printf("\nHandler: Ricevuto signal %d (%s)\n\n", iSignalCode, strsignal(iSignalCode));
  return;
}

int main() {
	struct sigaction sa;

	memset(&sa, '\0', sizeof(struct sigaction));
	sa.sa_handler = SigIntHandler;
	if(sigaction(SIGINT, &sa, NULL) == -1)
		perror("Main: errore in sigaction");
	printf("Premere Ctrl+C per eseguire l'handler di SIGINT\n");
	
	printf("Creazione ed apertura fifo...\n");	
	mkfifo("mypipe", S_IRUSR|S_IWUSR);
	/* if(open("mypipe", O_RDONLY) == -1) {
		char s[128]; 
		sprintf(s, "Main: errore in open, errno = %d", errno);
		perror(s);
		exit(EXIT_FAILURE);
	}	 */
	while ((open("mypipe", O_RDONLY) == -1) && (errno == EINTR));	
	
	return EXIT_SUCCESS;
}
