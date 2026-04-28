#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>

/***************************************************************************
* Soluzione con MUTEX. Vedi "man pthread_mutex_init"
***************************************************************************/
long counter = 0L, lMax = 100000L;   // variabile condivisa
pthread_mutex_t my_mutex = PTHREAD_MUTEX_INITIALIZER;

void *worker(void *arg)
{
    for (long i = 0L; i < lMax; i++) {
		pthread_mutex_lock(&my_mutex);
        counter++;   // <-- NON atomico. Sezione critica
		pthread_mutex_unlock(&my_mutex);
    }
    return NULL;
}

int main(int argc, char *argv[])
{
    pthread_t t1, t2;
	
	if(argc > 1)
		lMax = atol(argv[1]);

    pthread_create(&t1, NULL, worker, NULL);
    pthread_create(&t2, NULL, worker, NULL);

    pthread_join(t1, NULL);
    pthread_join(t2, NULL);

    printf("counter = %d\n", counter);
	pthread_mutex_destroy(&my_mutex);
    return 0;
}
