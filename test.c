#include <stdio.h>
typedef struct {
    int v;
} Integer;
int main(void) {
    const char* str = "\n        multi-line strings in C!\n    ";
    printf("%s\n", str);
    Integer i = (Integer){.v=10005};
    printf("JS objects -> C99 struct literals: %d\n", i.v);
}

