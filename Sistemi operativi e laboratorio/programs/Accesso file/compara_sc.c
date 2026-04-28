#include <fcntl.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdio.h>

int main(int argc, char *argv[])
{
	int iRet = EXIT_SUCCESS;
	
	if(argc >= 4) {
		long lCount = atol(argv[2]), lSize = atol(argv[3]);
		char sBuffer[lSize];
		int fd;
		
		memset(sBuffer, '1', lSize - 1);
		sBuffer[lSize - 1] = '\n';
		printf("Using System Calls, %ld scritture di %ld bytes...\n", lCount, lSize);
		
		if((fd = open(argv[1], O_WRONLY|O_CREAT|O_TRUNC, S_IRUSR|S_IWUSR|S_IRGRP|S_IWGRP)) == -1) {
			perror("Open error");
			iRet = EXIT_FAILURE;
		}
		else {
			for(int i = 1; i <= lCount; i++) {
				if(write(fd, sBuffer, sizeof(sBuffer)) != sizeof(sBuffer)) {
					perror("Write error");
					iRet = EXIT_FAILURE;
					break;
				}
			}
			close(fd);
		}
	}
	else {
		fprintf(stderr, "Usage: %s   file_to_write  write_count   write_size\n", argv[0]);
		iRet = EXIT_FAILURE;
	}
	
	return iRet;
}
