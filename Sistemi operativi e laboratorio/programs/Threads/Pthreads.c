#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <unistd.h>
#include <time.h>

/*****************************************************************************************
* Funzione che ritorna un intero random tra 1 e 9.
* I più curiosi possono andare a guardare:
* - man clock_gettime
* - man timespec
*****************************************************************************************/
long GetRandValue_between_1_9() {
	struct timespec ts;
	clock_gettime(CLOCK_REALTIME, &ts);
	return ((ts.tv_nsec % 9L) + 1L);
}

/* Funzione eseguita dai 2 threads */
void *vai(void *p) { // In realtà riceviamo il numero identificante il thread, non un indirizzo!
	long lThread = (long)p, lInterval = GetRandValue_between_1_9();
	
	for(int i = 1; i <= 10; i++) {
		printf("Thread %ld: i = %d\n", lThread, i);
		usleep(lInterval * 100000L); 	// Attendiamo lInterval decimi di secondo.
										// usleep(3) riceve come argomento il numero
										// di microsecondi da attendere
	}
	printf("Thread %ld: END FOR\n", lThread);
	pthread_exit((void *)lInterval); // In realtà torniamo il numero di decimi di secondo attesi, non un indirizzo!
}

int main() {
	pthread_t th1, th2;
	long lVal, lInterval = GetRandValue_between_1_9();

	// Creazione di due threads
	pthread_create(&th1, NULL, vai, (void **)1L); // In realtà passiamo il numero identificante il thread, non un indirizzo!
	pthread_create(&th2, NULL, vai, (void **)2L);

	for(int i = 1; i <= 10; i++) {
		printf("Main: i = %d\n", i);
		usleep(lInterval * 100000L);
	}
	printf("Main: END FOR\n");
	
	pthread_join(th1, (void **)&lVal); // Attendiamo il primo thread
	printf("Main: valore tornato dal thread 1 = %ld\n", lVal);
	pthread_join(th2, (void **)&lVal);
	printf("Main: valore tornato dal thread 2 = %ld. My interval = %ld\n", lVal, lInterval);

	return EXIT_SUCCESS;
}
