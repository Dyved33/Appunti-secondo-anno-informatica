#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <stdatomic.h>

/***************************************************************************
* Soluzione con _Atomic (C11: Concurrent support library).
* Vedi https://en.cppreference.com/w/c/thread.html
***************************************************************************/
_Atomic long counter = 0L, lMax = 100000L;   // variabile condivisa

void *worker(void *arg)
{
    for (long i = 0L; i < lMax; i++) {
        counter++;   // <-- NON atomico. Sezione critica
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
    return 0;
}
