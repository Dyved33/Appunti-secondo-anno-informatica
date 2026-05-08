#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdbool.h>
#include <errno.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <sys/socket.h>
#include <signal.h>

int main()
{
    int iRet = EXIT_SUCCESS, fd_s, fd_c, client_len;
    struct sockaddr_un sa, sa_client;
    struct sigaction sa;

    memset(&sa, '\0', sizeof(struct sigaction));
    sa.sa_handler = Handler;
    sigaction(SIGCHLD, &sa, NULL);

    if ((fd_s = socket(AF_UNIX, SOCK_STREAM, 0)) != -1)
    {
        sa.sun_family = AF_UNIX;
        strcpy(sa.sun_path, "/tmp/my_socket");
        if (bind(fd_s, (struct sockaddr *)&sa, sizeof(sa)) != -1)
        {
            if (listen(fd_s, 16) != -1)
            {
                while (true)
                {
                    client_len = sizeof(struct sockaddr_un);
                    if (fd_c = accept(fd_s, (struct sockaddr *)&sa_client, &client_len) != -1)
                    {
                        printf("Connessione avvenuta , fd_c = %d\n", fd_c);
                        char p[] = "Ciao dal server";
                        if (write(fd_c, p, sizeof(p)) == -1)
                        {
                            perror("Errore in write");
                            iRet = EXIT_FAILURE;
                        }
                        close(fd_c);
                    }
                        else
                        {
                            perror("Errore in accept");
                            iRet = EXIT_FAILURE;
                        }
                    }
                    else
                    {
                        perror("Errore in listen");
                        iRet = EXIT_FAILURE;
                    }
                }
            }
            else
            {
                perror("Errore in bind");
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
}