#include <errno.h>
#include <semaphore.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <unistd.h>
#include <assert.h>

sem_t sem;

#define handle_error(msg)  \
   do                      \
   {                       \
      perror(msg);         \
      _exit(EXIT_FAILURE); \
   } while (0)

static void handler(int sig)
{
   write(STDOUT_FILENO, "sem_post() from handler\n", 24);
   if (sem_post(&sem) == -1)
   {
      write(STDERR_FILENO, "sem_post() failed\n", 18); // uso write e non printf perché printf non è async-signal-safe
      _exit(EXIT_FAILURE);
   }
}

int main(int argc, char *argv[])
{
   struct sigaction sa;
   struct timespec ts;
   int s;

   if (argc != 3)
   {
      fprintf(stderr, "Usage: %s <alarm-secs> <wait-secs>\n",
              argv[0]);
      _exit(EXIT_FAILURE);
   }

   if (sem_init(&sem, 0, 0) == -1)
      handle_error("sem_init");

   /* Establish SIGALRM handler; set alarm timer using argv[1]. */

   sa.sa_handler = handler;
   sigemptyset(&sa.sa_mask);
   sa.sa_flags = 0;
   if (sigaction(SIGALRM, &sa, NULL) == -1)
      handle_error("sigaction");

   alarm(atoi(argv[1])); // alarm dopo argv[1] secondi invia SIGALRM al processo, che viene gestito dall'handler definito sopra

   /* Calculate relative interval as current time plus
     number of seconds given argv[2]. */

   if (clock_gettime(CLOCK_REALTIME, &ts) == -1)
      handle_error("clock_gettime");

   ts.tv_sec += atoi(argv[2]);

   printf("%s() about to call sem_timedwait()\n", __func__);
   while ((s = sem_timedwait(&sem, &ts)) == -1 && errno == EINTR)
      // prova a fare sem_timedwait, che si blocca finché non viene postato il semaforo o finché non scade il timer. Se viene interrotto da un signal handler (nel nostro caso SIGALRM)
      // allora sem_timedwait ritorna -1 con errno settata a EINTR
      continue; /* Restart if interrupted by handler. */
   /* Check what happened. */

   if (s == -1)
   {
      if (errno == ETIMEDOUT)
         printf("sem_timedwait() timed out\n");
      else
         perror("sem_timedwait");
   }
   else
      printf("sem_timedwait() succeeded\n");

   _exit((s == 0) ? EXIT_SUCCESS : EXIT_FAILURE);
}
