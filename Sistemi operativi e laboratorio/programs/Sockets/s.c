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
	int iRet = EXIT_SUCCESS, fd_s, fd_c, sockaddr_len, iReadRet;
	struct sockaddr_un sa, sa_client;
	struct sockaddr_in sa_in, sa_client_in;
	struct sockaddr *pst_sockaddr, *pst_sockaddr_client;
	struct sigaction sig_sa;
	
	memset(&sig_sa, '\0', sizeof(struct sigaction));
	sig_sa.sa_handler = Handler;
	sigaction(SIGINT,  &sig_sa, NULL);
	sigaction(SIGPIPE, &sig_sa, NULL);
	
	if((fd_s = socket(iDomain, SOCK_STREAM, 0)) != -1) {
		if(iDomain == AF_UNIX) {
			sa.sun_family 		= AF_UNIX;
			strcpy(sa.sun_path, MY_LOCAL_SOCKET);		
			pst_sockaddr 		= (struct sockaddr *)&sa;
			pst_sockaddr_client	= (struct sockaddr *)&sa_client;
			sockaddr_len = 		sizeof(struct sockaddr_un);
		}
		else {
			sa_in.sin_family 		= AF_INET;
			sa_in.sin_port			= htons(MY_PORT);
			sa_in.sin_addr.s_addr	= inet_addr(MY_IP);
			pst_sockaddr 			= (struct sockaddr *)&sa_in;
			pst_sockaddr_client 	= (struct sockaddr *)&sa_client_in;			
			sockaddr_len 			= sizeof(struct sockaddr_in);
		}
		if(bind(fd_s, pst_sockaddr, sockaddr_len) != -1) {
			if(listen(fd_s, 16) != -1) {
				while(true) {
					if((fd_c = accept(fd_s, pst_sockaddr_client, &sockaddr_len)) != -1) {
						printf("Connessione avvenuta, fd_c = %d\n", fd_c);
						if(iDomain == AF_INET) {
							printf("Indirizzo del client %s, porta %d\n",
									inet_ntoa(((struct sockaddr_in *)pst_sockaddr_client)->sin_addr),
									ntohs(((struct sockaddr_in *)pst_sockaddr_client)->sin_port));
						}
						if(!bClientEoF) {  // vedi my_socket.h
							char p[] = MSG;
							if(bServerSignalPipe && iDomain == AF_UNIX) sleep(1); // vedi my_socket.h
							if(write(fd_c, p, strlen(p)) != -1){
								char p[128];
								if((iReadRet = read(fd_c, p, sizeof(p) - 1)) > 0) {
									p[iReadRet] = '\0';
									if(strcmp(ACK, p) == 0) printf("Client Acked...\n", p);
								}
								else {
									if(iReadRet == 0) {
										printf("Il client ha chiuso la connessione!\n");
									}
									else {
										perror("Errore in read");
										iRet = EXIT_FAILURE;
									}
								}
							}
							else {
								if(errno == EPIPE) {
									// Il client ha chiuso la connessione e, se utilizziamo
									// AF_UNIX, il processo riceve SIGPIPE. Con AF_INET
									// il comportamento è diverso
									printf("SIGPIPE ricevuto...\n");
								}
								else {
									perror("Errore in write");
									iRet = EXIT_FAILURE;
									break;
								}
							}
						}
						close(fd_c);
					}
					else {
						if(errno == EINTR) {
							// E' stato ricevuto un signal per cui è stato definito un handler.
							// accept() è stata interrotta, ha ritornato -1 ed errno
							// è stata impostata ad EINTR
							printf("Accept interrotta da handler...\n");
						}
						else {
							perror("Errore in accept");
							iRet = EXIT_FAILURE;
							break;
						}
					}
				} // end while(true)		
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