#include <sys/wait.h>
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <fcntl.h>

#define MY_FIFO "myfifo"

int main()
{
    int iRet = EXIT_SUCCESS, fd;
    char sBuf[128];

    //access verifica se posso accedere al file di cui ho passato il pathname, F_OK verifica se il file esiste, R_OK verifica se posso leggere il file, 
    //W_OK verifica se posso scrivere sul file, X_OK verifica se posso eseguire il file
    //restituisce 0 se il file esiste e -1 se non esiste o se non ho i permessi necessari per accedervi
    if (access(MY_FIFO, F_OK) != 0)
    {
        if (mkfifo(MY_FIFO, S_IRUSR | S_IWUSR) != 0)//crea una fifo con i permessi di lettura e scrittura per l'utente (con -1 restituisce un errore    )
        {
            perror("Error in mkfifo\n");
            return EXIT_FAILURE;
        }
    }

    if (fd = open(MY_FIFO, O_WRONLY) != -1)
    {
        printf("Open in write mode\n");
        sleep(5);
        sprintf(sBuf, "I wrote in the pipe\n");
        if (write(fd, sBuf, sizeof(sBuf)) != -1)
        {
            printf("Wrote in the pipe\n");
        }
        else
        {
            perror("Error in write");
            iRet = EXIT_FAILURE;
        }
        close(fd);
    }
    else
    {
        perror("Error in open");
        iRet = EXIT_FAILURE;
    }

    printf("Write exiting...\n");
    return iRet;
}