#include <sys/socket.h>
#include <sys/un.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <arpa/inet.h>
#include <stdbool.h>

#define MY_LOCAL_SOCKET		"/tmp/my_socket"
#define MY_PORT				5152
#define MY_IP				"10.0.2.15"

#define ACK					"Ack"
#define MSG					"Ciao dal server"

int iDomain = AF_INET;

/**************************************************************************/
// Se attivate questo flag ed utilizzate i socket AF_UNIX il server riceverà
// SIGPIPE all'atto di scrivere nel socket. Questo perché il client ha già
// chiuso la connessione. Per AF_INET il comportamento è diverso
bool bServerSignalPipe = false;
/**************************************************************************/

/**************************************************************************/
// Se attivate questo flag il client si accorgerà che il server ha chiuso
// la connessione (read() ritorna 0) all'atto di leggere dal socket
bool bClientEoF = false;
/**************************************************************************/



