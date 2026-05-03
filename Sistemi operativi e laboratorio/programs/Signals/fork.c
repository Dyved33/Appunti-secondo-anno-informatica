#include <sys/types.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <fcntl.h>
#include <stdio.h>
#include <string.h>
#include <signal.h>
#include <unistd.h>
#include <errno.h>
#include <stdlib.h>

#define ATTESA 8

volatile pid_t pidParent; // type qualifier that instructs the compiler that pid-t value may change unexpectedly
// that's because this variable is used inside a signal handler and it could be called (and modified) at any time

/*********************************************************************************/
// Notate che, quando Ctrl+C viene premuto, entrambi i processi ricevono SIGINT.
// I segnali generati da tastiera sono inviati a tutti i processi
// facenti parte del PROCESS GROUP
/*********************************************************************************/
void SigIntHandler(int iSignalCode)
{
	char *p;
	if (getpid() == pidParent)
		p = "Parent";
	else
		p = "Child";
	printf("\n%s: Ricevuto signal %d (%s)\n\n", p, iSignalCode, strsignal(iSignalCode));
	return;
}

int main()
{
	struct sigaction sa;

	pidParent = getpid();
	memset(&sa, '\0', sizeof(struct sigaction)); // memset is a standard C library function used to fill a block of memory with a specified byte value
	sa.sa_handler = SigIntHandler;
	// sa.sa_flags = SA_RESTART;
	if (sigaction(SIGINT, &sa, NULL) == -1)
		perror("Main: errore in sigaction");
	printf("Premere Ctrl+C per eseguire l'handler di SIGINT\n");

	pid_t pid = fork();
	if (pid == 0)
	{
		// child
		printf("Child: inizio attesa...\n");
		int iSleep = ATTESA;
		/*********************************************************************************/
		// La funzione sleep(3) torna se:
		// - l'intervallo di attesa è spirato. Torna 0
		// - viene interrotta da un signal handler. Torna il numero
		//   di secondi che rimanevano ancora da attendere.
		// "man 3 sleep" per maggiori informazioni.
		// Il richiamo in un ciclo garantisce l'attesa per il numero
		// di secondi (più o meno...) richiesto anche in caso di interruzione per un
		// signal handler
		/*********************************************************************************/
		while ((iSleep = sleep(iSleep)) > 0)
			;
		printf("Child: fine attesa, exiting...\n");
	}
	else
	{
		// parent
		/*************************************************************************/
		// La seguente waitpid() si blocca finchè il child specificato non cambia
		// stato (termina). Tale blocco mette in attesa il processo.
		// La ricezione di un signal gestito con un handler (nel nostro caso SIGINT)
		// provoca l'interruzione di waitpid() che ritorna -1 con errno settata a EINTR.
		// Per maggiori informazioni vedi "man 7 signal".
		// Senza correzioni il parent termina prima del child in caso
		// di richiamo del signal handler.
		waitpid(pid, NULL, 0);
		// Per avere un comportamento corretto sostituite la riga soprastante
		// con questa sottostante:
		//         while(waitpid(pid, NULL, 0) == -1 && errno == EINTR);
		//
		// Oppure, in alternativa, decommentate la riga:
		//              sa.sa_flags = SA_RESTART;
		/*************************************************************************/
		printf("Parent: child terminato, exiting...!\n");
	}

	return EXIT_SUCCESS;
}
