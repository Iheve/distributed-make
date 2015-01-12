#include<stdio.h>
#include<stdlib.h>
#include<math.h>

int main(int argc, char **argv) {

	if (argc != 3) {
		fprintf(stderr, "usage : %s début fin\n", argv[0]);
		exit(1);
	}

	int debut = atoi(argv[1]);
	int fin = atoi(argv[2]);
	int i;
	for(i = debut; i <= fin ; i++) {
		int j;
		char i_premier = 1;
		if (i % 2 == 0) {
			i_premier = 0;
		} else {
			for(j = 3 ; j < sqrt(i) ; j+=2) {
				if (i % j == 0) {
					i_premier = 0;
					break;
				}
			}
		}
		if (i_premier) printf("%d\n", i);
	}
	exit(0);
}
