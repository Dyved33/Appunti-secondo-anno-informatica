    #include "my_socket.h"

int main()
{
    int iRet = EXIT_SUCCESS, fd_client; //stato di terminazione del programma e file descriptor del socket client
    struct sockaddr_in sa;//struttura dati (internet socket address) per configuarare indirizzo IP e porta del server a cui connettersi

    // Ignora SIGINT, SIGQUIT e SIGTERM
    struct sigaction sig_sa;//struttura per modificare il comportamento dei segnali (settata per modificare il comportamento dei segnali scelti)
    memset(&sig_sa, 0, sizeof(sig_sa));
    sig_sa.sa_handler = SIG_IGN;
    sigaction(SIGINT, &sig_sa, NULL);
    sigaction(SIGQUIT, &sig_sa, NULL);
    sigaction(SIGTERM, &sig_sa, NULL);

    if ((fd_client = socket(iDomain, SOCK_STREAM, 0)) != -1)//creazione di endpoint per la comunicazione con specifica sul protocollo e dominio (TCP e IPv4) 
    {
        memset(&sa, '\0', sizeof(sa));

        sa.sin_family = iDomain;
        sa.sin_port = htons(MY_PORT);//conversione del numero di porta da host byte order a network byte order (big-endian)
        sa.sin_addr.s_addr = inet_addr(MY_IP);

        if (connect(fd_client, (struct sockaddr *)&sa, sizeof(struct sockaddr_in)) != -1)//connessione del socket client al server specificato da sa
        {
            printf("Connessione al server avvenuta con successo!\n");
            
            char buffer[256];
            while (1)
            {
                printf("Cmd>");
                fflush(stdout);//svuota esplicitamente il buffer dello standard output per assicurarsi che il prompt venga mostrato immediatamente a schermo (previene la perdita di dati in caso di crash improvviso del programma)

                if (fgets(buffer, sizeof(buffer), stdin) == NULL)
                {
                    //premendo Ctrl+D (EOF) durante la lettura da stdin, fgets restituisce NULL, quindi in questo caso si imposta il buffer a "exit" per terminare il client in modo ordinato
                    strcpy(buffer, "exit");
                }

                //per come funziona il comando fgets, potrebbe essere presente un carattere di nuova linea ('\n') alla fine del buffer che va rimosso per evitare problemi nell'invio del comando 
                size_t len = strlen(buffer);
                if (len > 0 && buffer[len - 1] == '\n')
                {
                    buffer[len - 1] = '\0';
                }

                // se l'utente non ha digitato nulla, ripeti il ciclo aspettando un comando valido 
                if (strlen(buffer) == 0)
                {
                    continue;
                }

                // scrive i byte del comando nel socket del client verso l'esecutore. Se ritorna -1 c'è un errore di rete
                if (write(fd_client, buffer, strlen(buffer)) == -1)
                {
                    perror("Errore nell'invio del comando");
                    iRet = EXIT_FAILURE;
                    break;
                }

                // se il comando è "exit", il client termina
                if (strcmp(buffer, "exit") == 0)
                {
                    printf("Uscita...\n");
                    break;
                }

                // siccome la funzione read è bloccante, ho bisogno di sapere quanti byte aspettarmi prima di ricevere l'output
                int out_size = 0;
                ssize_t bytes_recv = read(fd_client, &out_size, sizeof(out_size));
                if (bytes_recv <= 0)
                {
                    printf("Connessione chiusa dall'esecutore remoto.\n");
                    break;
                }

                if (out_size == -1)
                {
                    printf("Errore sul server nell'esecuzione del comando.\n");
                    continue;
                }

                // qui si riceve l'output effettivo del comando 
                if (out_size > 0)
                {
                    char *response = malloc(out_size + 1);
                    if (response == NULL)
                    {
                        perror("Errore nell'allocazione della memoria");
                        iRet = EXIT_FAILURE;
                        break;
                    }
                    int total_received = 0;
                    while (total_received < out_size)
                    {
                        int n = read(fd_client, response + total_received, out_size - total_received);
                        //riceve i byte dell'output del comando e li memorizza in response partendo dall'indirizzo response + total_received per evitare di sovrascrivere i dati già ricevuti. Il terzo parametro mi indica quanti byte ancora devo ricevere 
                        if (n <= 0)
                        {
                            printf("Connessione interrotta durante la ricezione dell'output \n");
                            break;
                        }
                        total_received += n;
                    }
                    
                    response[total_received] = '\0';//chiudo la stringa 
                    printf("%s\n", response);
                    free(response);
                }
                else
                {
                    // output vuoto
                    printf("Il comando non ha prodotto output\n");
                }
            }
            close(fd_client);
        }
        else
        {
            perror("Errore in connect");
            iRet = EXIT_FAILURE;
        }
    }
    else
    {
        perror("Errore in socket");
        iRet = EXIT_FAILURE;
    }

    return iRet;
}