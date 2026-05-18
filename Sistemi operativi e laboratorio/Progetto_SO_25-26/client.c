#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdbool.h>
#include <string.h>
#include <errno.h>
#include "my_socket.h"

int main() {
	int iRet = EXIT_SUCCESS, fd_s, iReadRet;
	struct sockaddr_un sa;
	
	if((fd_s = socket(AF_UNIX, SOCK_STREAM, 0)) != -1) {
		sa.sun_family = AF_UNIX;
		strcpy(sa.sun_path, "/tmp/my_socket");		
		if(connect(fd_s, (struct sockaddr *)&sa, sizeof(struct sockaddr_un)) != -1) {
			//close(fd_s);
			char p[128];
			sleep(1);
			if((iReadRet = read(fd_s, p, sizeof(p) - 1)) > 0) {
				p[iReadRet] = '\0';
				printf("Ho ricevuto dal server \"%s\"\n", p);
			}
			else {
				if(iReadRet == 0) {
					printf("Il server ha chiuso la connessione!\n");
				}
				else {
					perror("Errore in read");
					iRet = EXIT_FAILURE;
				}
			}
		}
		else {
			perror("Errore in connect");
			iRet = EXIT_FAILURE;
		}
		close(fd_s);
	}
	else {
		perror("Errore in socket");
		iRet = EXIT_FAILURE;
	}
	
	return iRet;
}
