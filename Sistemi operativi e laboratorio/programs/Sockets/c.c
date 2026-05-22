#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <errno.h>
#include <signal.h>
#include "my_socket.h"

void Handler(int iSignum) {
	printf("Ricevuto signal \"%s\"\n", strsignal(iSignum));
	return;
}

int main() {
	int iRet = EXIT_SUCCESS, fd_s, iReadRet, sockaddr_len;
	struct sockaddr_un sa;
	struct sockaddr_in sa_in;
	struct sockaddr *pst_sockaddr;
	struct sigaction sig_sa;
	
	memset(&sig_sa, '\0', sizeof(struct sigaction));
	sig_sa.sa_handler = Handler;
	sigaction(SIGPIPE, &sig_sa, NULL);
	
	if((fd_s = socket(iDomain, SOCK_STREAM, 0)) != -1) {
		if(iDomain == AF_UNIX) {
			sa.sun_family 		= AF_UNIX;
			strcpy(sa.sun_path, MY_LOCAL_SOCKET);		
			pst_sockaddr 		= (struct sockaddr *)&sa;
			sockaddr_len = 		sizeof(struct sockaddr_un);
		}
		else {
			sa_in.sin_family 		= AF_INET;
			sa_in.sin_port			= htons(MY_PORT);
			sa_in.sin_addr.s_addr	= inet_addr(MY_IP);
			pst_sockaddr 			= (struct sockaddr *)&sa_in;
			sockaddr_len 			= sizeof(struct sockaddr_in);
		}
		if(connect(fd_s, pst_sockaddr, sockaddr_len) != -1) {
			if(bServerSignalPipe && iDomain == AF_UNIX) {  // vedi my_socket.h
				close(fd_s);
			}
			else {
				char p[128];
				if(bClientEoF) sleep(1); // vedi my_socket.h
				if((iReadRet = read(fd_s, p, sizeof(p) - 1)) > 0) {
					p[iReadRet] = '\0';
					printf("Ho ricevuto dal server \"%s\"\n", p);
					strcpy(p, ACK);
					if(write(fd_s, p, strlen(p)) == -1) {
						if(errno == EPIPE) {
							printf("SIGPIPE ricevuto...\n");
						}
						else {
							perror("Errore in write");
							iRet = EXIT_FAILURE;
						}
					}
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