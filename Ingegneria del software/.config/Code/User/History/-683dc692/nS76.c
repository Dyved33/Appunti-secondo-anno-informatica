#import <stdio.h>
#import <stdlib.h>
#import <string.h>

void fill_even(char *s, int arr[], int sLen, int nR)
{
}

void fill_odd(char *s, int arr[], int sLen, int nR)
{
    int step = 2 * (nR - 1);

    for (int i = 0; i < nR; ++i)
    {
        if (i == 0 || i == nR - 1) /* first and last row has the same increment */
        {
            for (int j = i; j < sLen; j += 2 * (nR - 1))
            {
                arr[j] = s[j];
            }
        }
        else
        {

            for (int j = i; j < sLen; j += step)
            {
                arr[j] = s[j];
            }
            step -= 2;
        }
    }
}

char *convert(char *s, int numRows)
{
    int sLen = strlen(s);
    int arr[sLen];

    numRows % 2 == 0 ? fill_even(char *s, int arr[], int sLen, int nR) : fill_odd(char *s, int arr[], int sLen, int nR);

    return arr;
}