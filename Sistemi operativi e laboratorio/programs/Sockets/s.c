#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdbool.h>
#include <string.h>
#include <errno.h>
#include <signal.h>
#include "my_socket.h"

void Handler(int iSignum) {
	printf("Ricevuto signal \"%s\"\n", strsignal(iSignum));
	return;
}

int main() {
	int iRet = EXIT_SUCCESS, fd_s, fd_c, client_len, sockaddr_len;
	struct sockaddr_un sa, sa_client;
	struct sockaddr_in sa_in, sa_client_in;
	struct sockaddr *pst_sockaddr, *pst_sockaddr_client;
	struct sigaction sig_sa;
	
	memset(&sig_sa, '\0', sizeof(struct sigaction));
	sig_sa.sa_handler = Handler;
	sigaction(SIGINT,  &sig_sa, NULL);
	sigaction(SIGUSR1, &sig_sa, NULL);
	sigaction(SIGPIPE, &sig_sa, NULL);
	
	if((fd_s = socket(iDomain, SOCK_STREAM, 0)) != -1) {
		if(iDomain == AF_UNIX) {
			sa.sun_family = AF_UNIX;
			strcpy(sa.sun_path, "/tmp/my_socket");
			sockaddr_len = sizeof(struct sockaddr_un);
			pst_sockaddr = (struct sockaddr *)&sa;
			pst_sockaddr_client = (struct sockaddr *)&sa_client;
		}
		else {
			sa_in.sin_family = AF_INET;
			sa_in.sin_port = htons(MY_PORT);
			sa_in.sin_addr.s_addr = inet_addr(MY_IP);
			sockaddr_len = sizeof(struct sockaddr_in);
			pst_sockaddr = (struct sockaddr *)&sa_in;
			pst_sockaddr_client = (struct sockaddr *)&sa_client_in;
		}
		//if(bind(fd_s, (struct sockaddr *)&sa, sizeof(struct sockaddr_un)) != -1) {
		if(bind(fd_s, pst_sockaddr, sockaddr_len) != -1) {
			if(listen(fd_s, 16) != -1) {
				while(true) {
					//client_len = sizeof(struct sockaddr_un);
					if((fd_c = accept(fd_s, pst_sockaddr_client, &sockaddr_len)) != -1) {
						printf("Connessione avvenuta, fd_c = %d\n", fd_c);
						char p[] = "Ciao dal server";
						//sleep(1);
						if(write(fd_c , p, strlen(p)) == -1){
							if(errno == EPIPE) {
								printf("SIGPIPE ricevuto...\n");
							}
							else {
								perror("Errore in write");
								iRet = EXIT_FAILURE;
								break;
							}
						}
						close(fd_c);
					}
					else {
						if(errno == EINTR) {
							printf("Accept interrotta, terminazione...\n");
							break;
						}
						else {
							perror("Errore in accept");
							iRet = EXIT_FAILURE;
							break;
						}
					}
				}				
			}
			else {
				perror("Errore in listen");
				iRet = EXIT_FAILURE;
			}	
		}
		else {
			perror("Errore in bind");
			iRet = EXIT_FAILURE;
		}
		
		close(fd_s);
		if(iDomain == AF_UNIX) unlink("/tmp/my_socket");
	}
	else {
		perror("Errore in socket");
		iRet = EXIT_FAILURE;
	}
	
	return iRet;
}
