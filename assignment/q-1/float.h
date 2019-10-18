#ifndef H_FLOAT
#define H_FLOAT

// Compute the minimum representable number at
// machine accuracy (using this arch's 'long double')
// NOTE: No guarantees given by C99 standard as to what
// precision long double corresponds to.
unsigned long long int float_acc();

#endif
