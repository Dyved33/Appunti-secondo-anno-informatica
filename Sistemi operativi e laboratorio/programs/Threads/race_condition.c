#include <stdio.h>
#include <pthread.h>
//include <stdatomic.h> // Per usare variabili atomiche, se necessario

pthread_mutex_t my_mutex = PTHREAD_MUTEX_INITIALIZER; // my_mutex per proteggere l'accesso a counter
// se uso la libreria stdatomic.h, posso dichiarare counter come _Atomic int counter = 0; e non avrei bisogno di usare il mutex
int counter = 0;   // variabile condivisa

void *worker(void *arg)
{
    for (int i = 0; i < 100000; i++) {
        pthread_mutex_lock(&my_mutex); // blocca il my_mutex prima di accedere a counter
        counter++;   // <-- NON atomico
        pthread_mutex_unlock(&my_mutex); // sblocca il my_mutex dopo aver modificato counter
    }
    return NULL;
}

int main(void)
{
    pthread_t t1, t2;

    pthread_create(&t1, NULL, worker, NULL);
    pthread_create(&t2, NULL, worker, NULL);

    pthread_join(t1, NULL);
    pthread_join(t2, NULL);

    printf("counter = %d\n", counter);
    pthread:pthread_mutex_destroy(&my_mutex); // distrugge il my_mutex
    return 0;
}
