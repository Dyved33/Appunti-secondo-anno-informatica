#include <sys/socket.h>
#include <sys/un.h>
#include <netinet/in.h>
#include <arpa/inet.h>

#define MY_PORT		5152
#define MY_IP		"10.0.2.15"

int iDomain = AF_INET;