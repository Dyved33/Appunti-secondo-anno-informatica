#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>
#include <unistd.h>

int main()
{
    int iRet = EXIT_SUCCESS, fd;
    char sBuf[] = "Ciao bro!\n";

    fd = open("pippo.txt", O_RDWR | O_CREAT | O_TRUNC, S_IRUSR | S_IWUSR);
    if (fd >= 0)
    {
        if (write(fd, sBuf, strlen(sBuf)) == strlen(sBuf))
        {
            printf("Ho scritto nel file pippo.txt, fd = %d\n", fd);
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
