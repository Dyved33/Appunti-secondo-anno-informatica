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
		FILE *fd;
		
		memset(sBuffer, '1', lSize - 1);
		sBuffer[lSize - 1] = '\n';
		printf("Using Standard I/O Library, %ld scritture di %ld bytes...\n", lCount, lSize);
		
		if((fd = fopen(argv[1], "w")) == NULL) {
			perror("Open error");
			iRet = EXIT_FAILURE;
		}
		else {
			for(int i = 1; i <= lCount; i++) {
				if(fwrite(sBuffer, sizeof(sBuffer), 1, fd) != 1) {
					perror("Write error");
					iRet = EXIT_FAILURE;
					break;
				}
			}
			fclose(fd);
		}
	}
	else {
		fprintf(stderr, "Usage: %s   file_to_write  write_count   write_size\n", argv[0]);
		iRet = EXIT_FAILURE;
	}
	
	return iRet;
}
