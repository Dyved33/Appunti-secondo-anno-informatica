#include "my_socket.h"

volatile sig_atomic_t terminate_flag = 0; // flag per indicare la richiesta di terminazione del server (il tipo garantisce letture/scritture atomiche anche se interrotte da segnali)
pid_t children[MAX_CLIENTS];              // per tenere traccia dei PID dei processi figli attivi
int child_count = 0;

// handler per il segnale SIGTERM (per uscire dal ciclo di accept si importa la flag in questione a 1)
void handle_sigterm(int sig)
{
    (void)sig;
    terminate_flag = 1;
}

// elimina un PID quando il relativo processo termina
void remove_child_pid(pid_t pid)
{
    for (int i = 0; i < child_count; i++)
    {
        if (children[i] == pid)
        {
            for (int j = i; j < child_count - 1; j++) // ricompatta l'array dei PID per non avere "buchi" dopo la rimozione
            {
                children[j] = children[j + 1];
            }
            child_count--;
            break;
        }
    }
}

// handler per il segnale asincrono SIGCHLD (quando un figlio cambia stato o termina, non posso prevedere quando)
void handle_sigchld(int sig)
{
    (void)sig;
    int saved_errno = errno; // CONTROLLA CHE SERVA EFFETTIVAMENTE (di salvare lo stato di errno per ripristinarlo alla fine)
    pid_t pid;
    while ((pid = waitpid(-1, NULL, WNOHANG)) > 0) // -1 indica "attendi qualsiasi figlio", WNOHANG indica di non bloccare se non ci sono figli terminati, ma restituire 0
    {
        remove_child_pid(pid);
    }
    // il ciclo serve perché se più figli terminano contemporaneamente, il kernel potrebbe accorpare i segnali  SIGCHLD  inviandone solo uno al padre. Con il ciclo recuperiamo tutti gli zombie pendenti
    errno = saved_errno;
}

int main()
{
    int iRet = EXIT_SUCCESS, fd_server, fd_new; // socket di ascolto del server e socket temporaneo per la connessione con un client
    socklen_t client_len, sockaddr_len;
    struct sockaddr_in sa_in, sa_client_in;
    struct sigaction sig_sa;

    // ignora SIGINT e SIGQUIT (impedisce che il server termini premendo  Ctrl+C  o  Ctrl+\  sul terminale associato)
    memset(&sig_sa, 0, sizeof(sig_sa));
    sig_sa.sa_handler = SIG_IGN;
    sigaction(SIGINT, &sig_sa, NULL);
    sigaction(SIGQUIT, &sig_sa, NULL);

    // gestione di SIGTERM (terminazione controllata)
    struct sigaction sa_term;
    memset(&sa_term, 0, sizeof(sa_term));
    sa_term.sa_handler = handle_sigterm;
    sigaction(SIGTERM, &sa_term, NULL);

    // gestione di SIGCHLD (per evitare processi zombie)
    struct sigaction sa_chld;
    memset(&sa_chld, 0, sizeof(sa_chld));
    sa_chld.sa_handler = handle_sigchld;
    sa_chld.sa_flags = SA_RESTART; // indica al sistema operativo di riavviare automaticamente alcune chiamate di sistema bloccanti (come la  accept) se vengono interrotte dall'arrivo di questo segnale
    sigaction(SIGCHLD, &sa_chld, NULL);

    if ((fd_server = socket(iDomain, SOCK_STREAM, 0)) != -1)
    {
        memset(&sa_in, '\0', sizeof(struct sockaddr_in));
        memset(&sa_client_in, '\0', sizeof(struct sockaddr_in));

        sa_in.sin_family = iDomain;
        sa_in.sin_port = htons(MY_PORT);
        sa_in.sin_addr.s_addr = inet_addr(MY_IP);
        sockaddr_len = sizeof(struct sockaddr_in);

        if (bind(fd_server, (struct sockaddr *)&sa_in, sockaddr_len) != -1)
        {
            if (listen(fd_server, 16) != -1) // al massimo 16 connessioni in attesa
            {
                printf("Server avviato e in ascolto su %s:%d...\n", MY_IP, MY_PORT);
                while (!terminate_flag)
                {
                    client_len = sizeof(struct sockaddr_in);
                    fd_new = accept(fd_server, (struct sockaddr *)&sa_client_in, &client_len);
                    if (fd_new != -1)
                    {
                        printf("Connessione avvenuta...\n");

                        // Blocchiamo temporaneamente SIGCHLD per evitare race conditions sull'array children
                        sigset_t block_mask, old_mask;
                        sigemptyset(&block_mask);
                        sigaddset(&block_mask, SIGCHLD);
                        sigprocmask(SIG_BLOCK, &block_mask, &old_mask);
                        // DA CAPIRE SE QUESTA PARTE RELATIVA ALLA GESTIONE DEI CHILDREN SERVA EFFETTIVAMENTE
                        if (child_count >= MAX_CLIENTS)
                        {
                            printf("Raggiunto il numero massimo di client connessi (%d).\n", MAX_CLIENTS);
                            close(fd_new);
                            sigprocmask(SIG_SETMASK, &old_mask, NULL); // SE TOGLI IL RESTO RICORDATI DI PUNTI COME QUESTO
                            continue;
                        }

                        pid_t pid = fork();
                        if (pid == 0)
                        { // Processo figlio

                            sigprocmask(SIG_SETMASK, &old_mask, NULL); // ripristina i segnali di default/ignora appropriati per il figlio

                            char buffer[16];
                            snprintf(buffer, sizeof(buffer), "%d", fd_new); // converte il file descriptor in stringa per passarlo come argomento a execl

                            close(fd_server); // il figlio non usa il socket di ascolto (per non interferire con l'accettazione di nuovi client)

                            execl("./executor", "./executor", buffer, (char *)NULL);
                            // se execl ha successo non ritorna mai
                            perror("Errore in execl");
                            _exit(EXIT_FAILURE);
                        }
                        else if (pid > 0)
                        { // Processo padre (server)

                            printf("Processo figlio esecutore creato, pid = %d\n", pid);
                            children[child_count++] = pid;
                            close(fd_new); // Il padre non usa il socket del client
                            sigprocmask(SIG_SETMASK, &old_mask, NULL);// CONTROLLA CHE SERVA EFFETTIVAMENTE (di ripristinare la maschera dei segnali dopo aver modificato l'array dei figli)    
                        }
                        else
                        {
                            perror("Errore in fork");
                            close(fd_new);
                            sigprocmask(SIG_SETMASK, &old_mask, NULL);
                        }
                    }
                    else
                    {
                        if (errno == EINTR)// se interrotti dal segnale, il ciclo while controllerà se terminate_flag è impostato 
                        {                            
                            continue;
                        }
                        else
                        {
                            perror("Errore in accept");
                            iRet = EXIT_FAILURE;
                            break;
                        }
                    }
                }
            }
            else
            {
                perror("Errore in listen");
                iRet = EXIT_FAILURE;
            }
        }
        else
        {
            perror("Errore in bind");
            iRet = EXIT_FAILURE;
        }

        // se il server esce a causa di SIGTERM, notifica i figli
        if (terminate_flag)
        {
            printf("\nServer in chiusura per SIGTERM: chiudo %d esecutori...\n", child_count);
            // disabilita l'handler SIGCHLD per evitare interferenze durante l'attesa finale
            signal(SIGCHLD, SIG_DFL);

            for (int i = 0; i < child_count; ++i)
            {
                kill(children[i], SIGTERM);// invia un segnale  SIGTERM  a tutti i figli esecutori attivi per forzarli a chiudere le connessioni e terminare
            }

            // attende la terminazione di tutti i processi esecutori
            while (child_count > 0)
            {
                int status;
                pid_t w = wait(&status);// si blocca in attesa che ciascun figlio termini definitivamente, liberando le risorse del sistema operativo
                if (w > 0)
                {
                    printf("Processo esecutore pid %d terminato.\n", w);
                    remove_child_pid(w);
                }
                else
                {
                    if (errno == ECHILD)
                        break;
                }
            }
            printf("Tutti gli esecutori sono terminati. Spegnimento server.\n");
        }

        close(fd_server);
    }
    else
    {
        perror("Errore in socket");
        iRet = EXIT_FAILURE;
    }

    return iRet;
}
