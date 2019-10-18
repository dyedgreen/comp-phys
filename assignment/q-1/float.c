#include "float.h"

unsigned long long int float_acc() {
  unsigned long long int size = 0;
  long double eps = 1;
  while (eps > 0) {
    eps /= 2;
    size ++;
  }
  size --;
  return size;
}
