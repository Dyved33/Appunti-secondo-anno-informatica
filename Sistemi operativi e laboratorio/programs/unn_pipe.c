#include <sys/wait.h>
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>

#define READ_END 0
#define WRITE_END 1

int main()
{
    pid_t pid;
    int pipefd[2], iRet = EXIT_SUCCESS;
    char sBuf[128];
    /* fork a child process */
    if (pipe(pipefd) != -1)
    {
        pid = fork();
        if (pid == 0)
        {
            close(pipefd[READ_END]);
            printf("Child: my pid is %d\n", getpid());
            sprintf(sBuf, "Hello from child process\n");
            if (write(pipefd[WRITE_END], sBuf, sizeof(sBuf)) != -1)
            {
                printf("Child: wrote to pipe\n");
            }
            else
            {
                fprintf(stderr, "Child: Failed to write to pipe\n");
                iRet = EXIT_FAILURE;
            }
            close(pipefd[WRITE_END]);
            return iRet;
        }
        else
        { /* parent process */
            /* parent will wait for the child to complete */
            close(pipefd[WRITE_END]);
            wait(NULL);
            int iReadCount;
            if ((iReadCount = read(pipefd[READ_END], sBuf, sizeof(sBuf) - 1)) >= 0)
            {
                sBuf[iReadCount] = '\0';
                printf("Parent: read from pipe: \"%s\"\n", sBuf);
            }
            else
            {
                if (iReadCount == 0)
                {
                    fprintf(stderr, "Parent: No data to read from pipe\n");
                    iRet = EXIT_FAILURE;
                }
                else
                {
                    perror("Parent: Failed to read from pipe\n");
                    iRet = EXIT_FAILURE;
                }
            }
            close(pipefd[READ_END]);
        }
    }
    else /* parent success */
    {
        fprintf(stderr, "Pipe Failed\n");
        iRet = EXIT_FAILURE;
    }
    return iRet;
}