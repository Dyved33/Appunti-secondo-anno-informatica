#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <unistd.h>

void *vai(void *p)
{
    printf("Threads: saluti!\n");

    for (int i = 0; i <= 10; i++)
    {
        printf("Threads %d: i = %d\n", *((int *)p), i);
        usleep(1000000); // Attende 1 secondo perché questa funzione lavora come sleep ma prende come parametro i microsecondi
    }
    pthread_exit((void *)5); // Restituisce un valore
}

int main(void)
{
    pthread_t th1, th2;
    int index1, index2;
    long lVal;

    index1 = 1;
    pthread_create(&th1, NULL, vai, &index1);

    index2 = 2;
    pthread_create(&th2, NULL, vai, &index2);

    for (int i = 0; i <= 10; i++)
    {
        printf("Main: i = %d\n", i);
        usleep(1000000); // Attende 1 secondo
    }

    pthread_join(th1, (void **)&lVal); // Attende che il thread finisca
    // sleep(1); // Attende che il thread finisca
    printf("Main: valore tornato dal thread: %ld\n", lVal);
    return EXIT_SUCCESS;
}