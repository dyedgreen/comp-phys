#ifndef H_FLOAT
#define H_FLOAT

// Compute the machine accuracy by computing the
// accuracy of the `long double` type, which should
// map to the best floating point natively supported
// (although there is no guarantees in the standard).
unsigned long long int float_acc();

#endif
