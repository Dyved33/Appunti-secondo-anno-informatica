#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdbool.h>
#include <errno.h>
#include <sys/socket.h>
#include <sys/un.h>

int main()
{
    int iRet = EXIT_SUCCESS, fd_s, iReadRet;
    struct sockaddr_un sa, sa_client;

    if ((fd_s = socket(AF_UNIX, SOCK_STREAM, 0)) != -1)
    {
        sa.sun_family = AF_UNIX;
        strcpy(sa.sun_path, "/tmp/my_socket");
        if (connect(fd_s, (struct sockaddr *)&sa, sizeof(struct sockaddr_un)) != -1)
        {
            char p[128];
            if ((iReadRet = read(fd_s, p, sizeof(p))) != -1)
            {
                p[iReadRet] = '\0';
                printf("Messaggio ricevuto: %s\n", p);
            }
            else
            {
                perror("Errore in read");
                iRet = EXIT_FAILURE;
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