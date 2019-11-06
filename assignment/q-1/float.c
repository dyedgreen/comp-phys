#include "float.h"

unsigned long long int float_acc() {
  // Compute eps in constant time using the next
  // float after 1. (Note, this works if the long
  // double is a IEEE standard float.)
  long double eps = 1.0l;
  long long int * iEps = (long long int *)&eps;
  *iEps += 1;
  eps -= 1.0l;
  // Convert size to power of two, ruining our O(1)
  // runtime...
  unsigned long long int size = 0;
  while (eps < 1l) {
    eps *= 2;
    size ++;
  }
  return size;
}
