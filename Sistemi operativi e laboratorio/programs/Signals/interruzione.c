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
	/////////////////////sa.sa_flags = SA_RESTART; // VEDI COMMENTI SUCCESSIVI
	if(sigaction(SIGINT, &sa, NULL) == -1)
		perror("Main: errore in sigaction");
	printf("Premere Ctrl+C per eseguire l'handler di SIGINT\n");
	
	printf("Creazione ed apertura fifo...\n");	
	mkfifo("mypipe", S_IRUSR|S_IWUSR); 	// Se la pipe esiste già questa funzione fallisce ma non controlliamo...
										// Qui siamo concentrati su altro....
	/*************************************************************************/									
	// La seguente open() si blocca finchè qualcun altro apre nella modalità
	// opposta (write) la pipe. Tale blocco mette in attesa il processo.
	// La ricezione di un signal gestito con un handler (nel nostro caso SIGINT)
	// provoca l'interruzione di open() che ritorna -1 con errno settata a EINTR.
	// Per maggiori informazioni vedi "man 7 signal".
	if(open("mypipe", O_RDONLY) == -1) {
		char s[128]; 
		sprintf(s, "Main: errore in open, errno = %d", errno);
		perror(s);
		exit(EXIT_FAILURE);
	}	
	/*************************************************************************/
	
	/*************************************************************************/
	// Di seguito riportiamo due possibili soluzioni per evitare che la open()
	// ritorni per errore quando viene ricevuto un signal oggetto di handler.
	//
	// Soluzione 1:
	// decommentate la riga:
	//         sa.sa_flags = SA_RESTART;   
	// La open() viene automaticamente restartata. Potete verificarlo eseguendo strace in
	// un altro terminale
	//
	// Soluzione 2:
	// sostituire il soprastante richiamo di open() nel seguente modo:
	/* int iOpenRet;
	while(((iOpenRet = open("mypipe", O_RDONLY)) == -1) && (errno == EINTR));
	if(iOpenRet == -1) {
		char s[128]; 
		sprintf(s, "Main: errore in open, errno = %d", errno);
		perror(s);
		exit(EXIT_FAILURE);
	} */
	/*************************************************************************/
	
	return EXIT_SUCCESS;
}
