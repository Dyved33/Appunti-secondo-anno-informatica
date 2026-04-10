#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/wait.h>
#include <fcntl.h>
#include <errno.h>
#include <sys/stat.h>

#define READ_END 0
#define WRITE_END 1

int main()
{
    int iRet = EXIT_SUCCESS, pipefd[2];
    pid_t pid_cat, pid_sort;

    if (pipe(pipefd) == 0)
    {
        pid_cat = fork(); // ho creato una copia quasi uguale del processo parent
        if (pid_cat == 0)
        {
            // cat
            char *sArgs[] = {"cat", "/etc/passwd", NULL};
            close(pipefd[READ_END]); // chiudo la pipe in lettura perché mi serve solo per scrivere con cat
            dup2(pipefd[WRITE_END], STDOUT_FILENO); // ridireziono l'output standard di cat alla pipe
            dup2(pipefd[WRITE_END], STDERR_FILENO); // ridireziono l'output di errore standard di cat alla pipe
            execve("/usr/bin/cat", sArgs, NULL);
            perror("Errore nell'esecuzione di cat");
            return EXIT_FAILURE;
        }
        else
        {
            // parent
            pid_sort = fork();
            if (pid_sort == 0)
            {
                // sort
                char *sArgs[] = {"sort", NULL};
                close(pipefd[WRITE_END]); // chiudo la pipe in scrittura 
                dup2(pipefd[READ_END], STDIN_FILENO); // ridireziono l'input standard di sort alla pipe
                int fd = open("pippo.out", O_WRONLY | O_CREAT | O_TRUNC, S_IRUSR | S_IWUSR); // apro il file pippo.out in scrittura
                dup2(fd, STDOUT_FILENO); // ridireziono l'output standard di sort al file pippo.out
                dup2(fd, STDERR_FILENO); // ridireziono l'output di errore standard di sort al file pippo.out
                execve("/usr/bin/sort", sArgs, NULL);
                perror("Errore nell'esecuzione di sort");
                return EXIT_FAILURE;
            }
            else
            {
                // parent
                int wstatus;
                close(pipefd[READ_END]);
                close(pipefd[WRITE_END]);
                //al parent non serve più la pipe (quindi nessuno dei due endpoint)
                waitpid(pid_cat, &wstatus,0); // aspetto che cat finisca
                printf("Il processo cat è terminato con stato %d\n", WEXITSTATUS(wstatus));
                waitpid(pid_sort, &wstatus,0); // aspetto che sort finisca
                printf("Il processo sort è terminato con stato %d\n", WEXITSTATUS(wstatus));
            }
        }
    }
    else
    {
        perror("Errore nella creazione della pipe");
        iRet = EXIT_FAILURE;
    }

    return iRet;
}