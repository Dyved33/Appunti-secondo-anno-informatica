#include<stdio.h>
#include<string.h>
#include<stdlib.h>
#include<signal.h>
#include<unistd.h>

void SigIntHandler(int iSignalCode) {
  printf("Handler: ricevuto signal %d (%s)\n", iSignalCode, strsignal(iSignalCode));
  return;
}
int main() {
	struct sigaction sa;
	int i;
	
	memset(&sa, '\0', sizeof(struct sigaction));
	sa.sa_handler = SigIntHandler;	// Provate a variare questa assegnazione
									// come specificato nelle due righe seguenti
	//sa.sa_handler = SIG_DFL;
	//sa.sa_handler = SIG_IGN;
	sigaction(SIGINT, &sa, NULL);
	printf("Premere Ctrl+C per eseguire la \"disposition\" di SIGINT\n");
	printf("Se avete definito un handler per Ctrl+C o l'avete ignorato,\n");
	printf("Potete uscire con Ctrl+\\ (provoca l'invio di SIGQUIT)\n");
	for(i = 1; ; i++) {
		sleep(1);
		printf("%d\n", i);
	}
	return EXIT_SUCCESS;
}
