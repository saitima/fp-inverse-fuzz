#include <stdio.h>
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include "libmontinversego.h"

int hex2bin(char *source_str, char *dest_buffer)
{
  char *line = source_str;
  char *data = line;
  int offset;
  int read_byte;
  int data_len = 0;

  while (sscanf(data, " %02x%n", &read_byte, &offset) == 1)
  {
    dest_buffer[data_len++] = read_byte;
    data += offset;
  }
  return data_len;
}

int main()
{
  char buf[] = "";  
  // char str[] = "080000f98889454fc308000001fe45623da1"; 
  char str[] = "01f301e3";
  int o_len;
  int len = hex2bin(str, buf);
  char* o_buff = (char *)malloc(64*sizeof(char));
  o_buff = c_perform_inverse(buf, len, &o_len);  
  if(o_len == 0)
  {
    printf("no inverse\n");
    return 1;
  }
  printf("o len %d\n", o_len);
  printf("inverse: %s\n", o_buff);
  printf("success!");
  return 0;
}