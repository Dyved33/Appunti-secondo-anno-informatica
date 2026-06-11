#include "my_socket.h"

// se l'esecutore riceve un SIGTERM dal server, termina con successo in modo pulito
void handle_sigterm(int sig)
{
    (void)sig;
    exit(EXIT_SUCCESS);
}

int main(int argc, char *argv[])
{
    // verifico che il descrittore del client sia stato passato come argomento
    if (argc < 2)
    {
        fprintf(stderr, "Uso: %s <client_fd>\n", argv[0]); // per scrivere su strerr dati formattati (simile a printf ma richiede un puntatore allo stream come argomento)
        exit(EXIT_FAILURE);
    }

    int client_fd = atoi(argv[1]); // converte la stringa passata dal server come argomento della linea di comando in un intero (client_fd)

    // Ignora SIGINT e SIGQUIT
    struct sigaction sig_sa;
    memset(&sig_sa, 0, sizeof(sig_sa));
    sig_sa.sa_handler = SIG_IGN;
    sigaction(SIGINT, &sig_sa, NULL);
    sigaction(SIGQUIT, &sig_sa, NULL);

    // gestione di SIGTERM
    struct sigaction sa_term;
    memset(&sa_term, 0, sizeof(sa_term));
    sa_term.sa_handler = handle_sigterm;
    sigaction(SIGTERM, &sa_term, NULL);

    char command[256];

    while (1)
    {
        memset(command, 0, sizeof(command)); // pulisce il buffer del comando

        // attendecomando dal client
        ssize_t bytes_read = recv(client_fd, command, sizeof(command) - 1, 0);

        if (bytes_read <= 0) // se il client chiude la connessione o si verifica un errore, termina
        {
            break;
        }

        command[bytes_read] = '\0'; // chiude la stringa del comando

        size_t len = strlen(command);
        if (len > 0 && command[len - 1] == '\n') // rimozione di eventuali caratteri di a capo
        {
            command[len - 1] = '\0';
        }

        // se riceve il comando "exit", esce dal loop e termina
        if (strcmp(command, "exit") == 0)
        {
            break;
        }

        // creazione pipe per catturare l'output del comando
        int pipefd[2];
        if (pipe(pipefd) == -1)
        {
            perror("Errore creazione pipe");
            continue;
        }

        pid_t pid = fork();
        if (pid == 0)
        { // processo figlio

            close(pipefd[0]); // chiudo la pipe all'estremità di lettura nel figlio (il figlio scriverà solo sulla pipe)

            // qualsiasi output del comando non andrà a video, ma entrerà nella pipe
            dup2(pipefd[1], STDOUT_FILENO); // ridireziona lo Standard Output del comando nel descrittore di scrittura della pipe
            dup2(pipefd[1], STDERR_FILENO); // fa lo stesso per lo Standard Error
            close(pipefd[1]);               // chìiude il vecchio descrittore dato che è stato duplicato

            // ripristina i comportamenti di default per i segnali (così che l'eseguibile richiamato possa reagire correttamente a eventuali segnali)
            struct sigaction sa_default;
            memset(&sa_default, 0, sizeof(sa_default));
            sa_default.sa_handler = SIG_DFL;
            sigaction(SIGINT, &sa_default, NULL);
            sigaction(SIGQUIT, &sa_default, NULL);

            // parsing degli argomenti
            char *args[64];
            int arg_count = 0;
            char *token = strtok(command, " ");
            // strtok spezza la stringa usando lo spazio come delimitatore, riempiendo il vettore di puntatori a stringhe args e terminandolo con  NULL  (necessario per l'exec)
            while (token != NULL && arg_count < 63)
            {
                args[arg_count++] = token;
                token = strtok(NULL, " ");
            }
            args[arg_count] = NULL;

            if (arg_count > 0)
            {
                execvp(args[0], args); // cerca il comando in  args[0]  all'interno dei percorsi di sistema (variabile d'ambiente  PATH) ed esegue il programma passandogli gli argomenti in  args
                // se execvp ritorna, c'è stato un errore
                perror("Errore execvp");
            }
            exit(EXIT_FAILURE);
        }
        else if (pid > 0)
        { // processo padre

            close(pipefd[1]);// chiudo la pipe all'estremità di scrittura nel padre (il padre leggerà solo dalla pipe)

            // alloca un buffer per leggere l'output del comando
            char output_buffer[65536]; 
            memset(output_buffer, 0, sizeof(output_buffer));

            ssize_t total_bytes = 0;
            ssize_t n;

            // legge dalla pipe l'output prodotto dal comando fino a quando il processo figlio non ha finito (EOF sulla pipe)
            while ((n = read(pipefd[0], output_buffer + total_bytes, sizeof(output_buffer) - total_bytes - 1)) > 0)
            {
                total_bytes += n;
                if (total_bytes >= (ssize_t)sizeof(output_buffer) - 1)
                {
                    break;
                }
            }
            close(pipefd[0]);
            
            waitpid(pid, NULL, 0);// l'esecutore attende la terminazione del figlio che ha eseguito il comando per ripulirlo (per evitare processi zombie)

            // invia prima la dimensione dell'output (come int) al client
            int out_size = (int)total_bytes;
            send(client_fd, &out_size, sizeof(out_size), 0);

            // invia l'output effettivo al client
            if (out_size > 0)
            {
                send(client_fd, output_buffer, out_size, 0);
            }
        }
        else
        {
            perror("Errore fork esecutore");
            int err_size = -1;
            send(client_fd, &err_size, sizeof(err_size), 0);
        }
    }

    close(client_fd);
    return EXIT_SUCCESS;
}