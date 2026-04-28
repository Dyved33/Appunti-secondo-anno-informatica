#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>
#include <unistd.h>

int main(int argc, char *argv[])
{
    int iRet = EXIT_SUCCESS, fd, iReadRet;
    char sBuf[] = "Ciao!\n";
    char sReadBuf[128];

    if (argc < 2)
    {
        printf("Dammi il nome del file\n");
        return EXIT_FAILURE;
    }

    fd = open(argv[1], O_RDWR | O_CREAT | O_TRUNC, S_IRUSR | S_IWUSR);
    if (fd >= 0)
    {
        if (write(fd, sBuf, strlen(sBuf)) == strlen(sBuf))
        {
            printf("Ho scritto nel file %s, fd = %d\n", argv[1], fd);

            if (lseek(fd, 0, SEEK_SET) != -1)
            {
                if ((iReadRet = read(fd, sReadBuf, sizeof(sReadBuf) - 1)) != -1) // sReadBuf è l'indirizzo al primo byte del buffer
                {
                    char sWriteBuf[] = "Questa proviene da write()\n";
                    int iCount = strlen(sWriteBuf);
                    sReadBuf[iReadRet] = '\0';
                    dup2(fd,1); // duplica fd su STDOUT_FILENO, ora invece di scrivere su STDOUT_FILENO scrive su fd
                    printf("Ho riletto per %d byte: %s", iReadRet, sReadBuf);//printf scrive sul descrittore 1
                    /* write(STDOUT_FILENO, sWriteBuf, iCount); */
                }
                else
                {
                    perror("Errore in read()");
                    iRet = EXIT_FAILURE;
                }
            }
            else
            {
                perror("Errore in lseek()");
                iRet = EXIT_FAILURE;
            }
            close(fd);
        }
        else
        {
            printf("errno = %d\n", errno);
            perror("Errore in write()");
            iRet = EXIT_FAILURE;
        }
    }
    else
    {
        printf("errno = %d\n", errno);
        perror("Errore in open()"); /* stampa la stringa a video: descrizione errore*/
        iRet = EXIT_FAILURE;
    }
    return iRet;
}