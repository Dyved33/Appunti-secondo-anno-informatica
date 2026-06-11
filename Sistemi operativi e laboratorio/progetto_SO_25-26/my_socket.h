#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdbool.h>
#include <string.h>
#include <errno.h>
#include <signal.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/wait.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <arpa/inet.h>

#define MY_PORT 8080
#define MY_IP "127.0.0.1"
#define MAX_CLIENTS 16

static const int iDomain = AF_INET;
// dominio di comunicazione, AF_INET per socket di rete, AF_UNIX per socket locali