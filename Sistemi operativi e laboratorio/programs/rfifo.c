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

    if (fd = open(MY_FIFO, O_RDONLY) != -1)
    {
        int iReadCount;
        if ((iReadCount = read(fd, sBuf, sizeof(sBuf)) - 1) != -1)
        {
            printf("Open in read mode\n");
            if (iReadCount == 0)
            {
                printf("Every descriptor closed\n");
            }
            else
            {
                sBuf[iReadCount] = '\0';
                printf("Read from the pipe: \"%s\"\n", sBuf);
            }
        }
        else
        {
            perror("Error in read");
            iRet = EXIT_FAILURE;
        }
        close(fd);
    }
    else
    {
        perror("Error in open");
        iRet = EXIT_FAILURE;
    }

    printf("Read exiting...\n");
    return iRet;
}