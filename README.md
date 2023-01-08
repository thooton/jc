# jc
`jc` is a tool that enables the meta-programming of C with JavaScript.

```cpp
// Specialize exponentiation for a power
function powFor(pwr) {
    // A quote is a piece of C code expressed as a JS value.
    var expr = quote(1);
    // Add " * v" to the expression `pwr` times.
    for (var i = 0; i < pwr; i++) {
        expr.add(quote( * v));
    }

    // Create a C function as a JS value.
    const pow_fn = (int v) -> static int {
        // JS interpolation: insert a JS value into C code.
        return ${expr};
    };
    // Return it.
    return pow_fn;
}

const pow7 = powFor(7);
/*
static int jc_lfn1(int v) {
    return 1 * v * v * v * v * v * v * v;
}
*/

#include <stdio.h>
int main(void) {
    printf("got %d\n", ${pow7}(2));
    return 0;
}
```
```sh
> jc -i example.jc -o example.c -n node && gcc example.c && ./a.out
got 128
```
# Examples
## Compute static data at compile-time
```cpp
var fib = [];
var prev = 0;
var cur = 1;
for (var i = 0; i < 40; i++) {
    var next = prev + cur;
    fib.push(cur);
    prev = cur;
    cur = next;
}

const int FIB_SEQ[${fib.length}] = ${fib};
#include <stdio.h>
int main(void) {
    for (int i = 0; i < ${fib.length}; i++) {
        printf("%d ", FIB_SEQ[i]);
    }
    return 0;
}
```
```
1 1 2 3 5 8 13 21 34 55 89 144 233 377 610 987 1597 2584 4181 6765 10946 17711 28657 46368 75025 121393 196418 317811 514229 832040 1346269 2178309 3524578 5702887 9227465 14930352 24157817 39088169 63245986 102334155
```

## Compile domain-specific languages
```cpp
#include <stdio.h>
#include <string.h>

function compileBf(code, tapesz) {
    var body = quote();
    var loopstack = [];
    for (const ch of code) {
        switch (ch) {
        case '>':
            body.add(quote(
                i++;
            ));
            break;
        case '<':
            body.add(quote(
                i--;
            ));
            break;
        case '+':
            body.add(quote(
                tape[i]++;
            ));
            break;
        case '-':
            body.add(quote(
                tape[i]--;
            ));
            break;
        case '.':
            body.add(quote(
                putchar(tape[i]);
            ));
            break;
        case ',':
            body.add(quote(
                tape[i] = getchar();
            ));
            break;
        case '[':
            var loop = {
                // generate a unique set of symbols for the label
                before: gensym(),
                after: gensym(),
            };
            loopstack.push(loop);
            body.add(quote(
                ${loop.before}:
                if (tape[i] == 0) {
                    goto ${loop.after};
                }
            ));
            break;
        case ']':
            var loop = loopstack.pop();
            if (!loop) {
                throw new Error("unmatched loop end");
            }
            body.add(quote(
                goto ${loop.before};
                ${loop.after}:
            ));
            break;
        }
    }
    if (loopstack.length > 0) {
        throw new Error("unmatched loop start");
    }
    return (void) -> static void {
        int tape[${tapesz}];
        int i = 0;
        memset(tape, 0, sizeof(tape));
        ${body}
    };
}

const add3 = compileBf("+[,+++.]", 10);
/*  
static void jc_lfn1(void) {
    int tape[10];
    int i = 0;
    memset(tape, 0, sizeof(tape));
    tape[i]++;
jc_sym1:
    if (tape[i] == 0) {
        goto jc_sym2;
    }
    tape[i] = getchar();
    tape[i]++;
    tape[i]++;
    tape[i]++;
    putchar(tape[i]);
    goto jc_sym1;
jc_sym2:
    ;
}
*/

int main(void) {
    ${add3}();
}
```
```
> a
d
> w
z
```

## Compile-time `printf`
```cpp
#include <string.h>
#include <stdio.h>
var mkprintf_memory = {};
function mkprintf(fmt) {
    if (mkprintf_memory[fmt]) {
        return mkprintf_memory[fmt];
    }
    var specifiers = [];
    var start_index = 0;
    var in_spec = false;
    for (var i = 0; i < fmt.length; i++) {
        const ch = fmt[i];
        if (!in_spec) {
            if (ch == '%') {
                const literal = fmt.substring(
                    start_index, i
                );
                specifiers.push(["lit", literal]);
                in_spec = true;
            }
            continue;
        }
        switch (ch) {
        case '%':
            specifiers.push(["lit", "%"]);
            break;
        case 's':
            specifiers.push(["str"]);
            break;
        case 'd':
            specifiers.push(["num"]);
            break;
        default:
            throw new Error("fmt specifier "+ch+" unimplemented");
        }
        start_index = i+1;
        in_spec = false;
    }
    if (start_index != fmt.length) {
        const literal = fmt.substring(
            start_index
        );
        specifiers.push(["lit", literal]);
    }
    var args = quote();
    var first_arg = true;
    var body = quote();
    for (var spec of specifiers) {
        var kind = spec[0];
        if (kind == "lit") {
            const literal = spec[1];
            body.add(quote(
                fwrite(${literal}, 1, ${literal.length}, stdout);
            ));
        } else if (kind == "num") {
            const arg = gensym();
            if (!first_arg) {
                args.add(quote(, ));
            } else {
                first_arg = false;
            }
            args.add(quote(int ${arg}));
            body.add(quote({
                char buf[32];
                int idx = sizeof(buf);
                int neg = ${arg} < 0;
                if (neg) {
                    ${arg} = -${arg};
                }
                do {
                    buf[--idx] = '0'+(${arg} % 10);
                    ${arg} /= 10;
                } while (${arg} > 0);
                if (neg) {
                    buf[--idx] = '-';
                }
                fwrite(buf+idx, 1, sizeof(buf)-idx, stdout);
            }));
        } else if (kind == "str") {
            if (!first_arg) {
                args.add(quote(, ));
            } else {
                first_arg = false;
            }
            const arg = gensym();
            args.add(quote(const char* ${arg}));
            body.add(quote(
                fwrite(${arg}, 1, strlen(${arg}), stdout);
            ));
        } else {
            throw new Error("internal mkprintf error");
        }
    }
    const printfn = (${args}) -> static void {
        ${body}
    };
    mkprintf_memory[fmt] = printfn;
    return printfn;
}

int main(void) {
    ${mkprintf("Hello, %s!\n")}("world");
    ${mkprintf("%s has %d %s.\n")}("He", 20, "stones");
}
```
```
Hello, world!
He has 20 stones.
```

## Compile-time regular expressions
```cpp
const spawn = require("child_process").spawn;
var mkmatches_memory = {};
function mkmatches(reg) {
    if (mkmatches_memory[reg]) {
        return mkmatches_memory[reg];
    }
    /* by using promises, we ensure
        mkmatches calls will be
        executed concurrently by the JS
        event loop during compilation */
    const promise = new Promise((resolve, reject) => {
        const re2c = spawn("re2c", ["-"]);
        var output = "";
        var errmsg = "";
        re2c.on("error", function(err) {
            reject("can't spawn re2c: "+err);
        });
        re2c.stdin.end(`
            /*!re2c
                re2c:yyfill:enable = 0;
                re2c:define:YYCTYPE = char;
                
                * { return 0; }
                ${reg} { return 1; }
            */
        `);
        re2c.stdout.on("data", function(chunk) {
            output += chunk;
        });
        re2c.stderr.on("data", function(chunk) {
            errmsg += chunk;
        });
        re2c.on("close", function(code) {
            if (code !== 0) {
                reject(
                    "re2c exited with code "
                    +code+" and msg "+errmsg
                );
                return;
            }
            resolve((const char* input) -> static int {
                const char* YYCURSOR = input;
                const char* YYMARKER;
                // convert a string to a quote
                ${strtoq(output)}
            });
        });
    });
    mkmatches_memory[reg] = promise;
    return promise;
}

#include <stdio.h>
#define matches1 ${mkmatches(`[A-Z]+ "test"`)}
#define matches2 ${mkmatches(`[a-z]+ "TEST"`)}

int main(void) {
    printf("regex 1: %d %d\n", 
        matches1("hellotest"), 
        matches1("HELLOtest")
    );
    printf("regex 2: %d %d\n",
        matches2("goodbyeTEST"),
        matches2("goodbyeTESt")
    );
}
```
```
regex 1: 0 1
regex 2: 1 0
```

## Generics
```cpp
#include <stdlib.h>
function mkvec(T) {
    return quote(
        typedef struct {
            ${T}* ptr;
            long len, cap;
        } vec_${T};
        static void vec_init_${T}(vec_${T}* v) {
            v->ptr = NULL;
            v->len = v->cap = 0;
        }
        static void vec_push_back_${T}(vec_${T}* v, ${T} val) {
            if (v->len == v->cap) {
                if (v->cap == 0) {
                    v->cap = 1;
                } else {
                    v->cap *= 2;
                }
                v->ptr = realloc(v->ptr, v->cap);
            }
            v->ptr[v->len++] = val;
        }
        static ${T} vec_pop_back_${T}(vec_${T}* v) {
            return v->ptr[--v->len];
        }
        static ${T} vec_pop_front_${T}(vec_${T}* v) {
            ${T} item = v->ptr[0];
            for (long i = 1; i < v->len; i++) {
                v->ptr[i-1] = v->ptr[i];
            }
            v->len--;
            return item;
        }
    );
}

${mkvec(quote(int))};

#include <stdio.h>
int main(void) {
    vec_int v;
    vec_init_int(&v);
    vec_push_back_int(&v, 5);
    vec_push_back_int(&v, 10);
    vec_push_back_int(&v, 15);
    int a = vec_pop_back_int(&v);
    int b = vec_pop_front_int(&v);
    int c = vec_pop_back_int(&v);
    printf("%d %d %d", a, b, c);
}
```
```
15 5 10
```

## Conditional compilation
```cpp

const util = require("util");
const execFile = util.promisify(
        require("child_process").execFile
);
const fs = require("fs/promises");
const cc = (async() => {
    if (process.env.CC) {
        return process.env.CC;
    }
    try {
        await execFile("gcc", ["-v"]);
        return "gcc";
    } catch (_err) {}
    try {
        await execFile("clang", ["-v"]);
        return "clang";
    } catch (_err) {}
    try {
        await execFile("cc", ["-v"]);
        return "cc";
    } catch (_err) {}
    throw new Error("can't find C compiler:"+
        "supply one by setting env variable CC");
})();
const has_std_atomics = (async() => {
    const tmpdir = await fs.mkdtemp("atomicjh");
    const tmpfile = tmpdir+"/test.c";
    await fs.writeFile(tmpfile, `
        #include <stdatomic.h>
        int main(void) {
            _Atomic(int) _a;
            return 0;
        }
    `);
    var supported = true;
    try {
        await execFile(await cc, [
                tmpfile, "-o", 
                tmpdir+"/a.out", "-O0"
        ]);
    } catch (_err) {
        supported = false;
    }
    await fs.rm(tmpdir, {recursive: true, force: true});
    return supported;
})();
${(async() => {
    if (await has_std_atomics) {
        return quote(
            #include <stdatomic.h>
            #define Relaxed memory_order_relaxed
            #define Acquire memory_order_acquire
            #define Release memory_order_release
            #define AcqRel memory_order_acq_rel
            #define SeqCst memory_order_seq_cst
        );
    } else {
        return quote(
            #define Relaxed __ATOMIC_RELAXED
            #define Acquire __ATOMIC_ACQUIRE
            #define Release __ATOMIC_RELEASE
            #define AcqRel  __ATOMIC_ACQ_REL
            #define SeqCst  __ATOMIC_SEQ_CST
        );
    }
})()};

function atomicStdFor(T) {
    return quote(
        typedef struct {
            _Atomic(${T}) repr;
        } atomic_${T};
        static inline atomic_${T} atomic_new_${T}(
            ${T} v
        ) {
            atomic_${T} a;
        #ifdef atomic_init
            atomic_init(&a.repr, v);
        #else
            a.repr = ATOMIC_VAR_INIT(v);
        #endif
            return a;
        }
        static inline ${T} atomic_load_${T}(
            const atomic_${T}* a,
            memory_order order
        ) {
            atomic_${T}* a_nonconst = (atomic_${T}*)a;
            return atomic_load_explicit(
                &a_nonconst->repr, order
            );
        }
        static inline void atomic_store_${T}(
            atomic_${T}* a, ${T} v,
            memory_order order
        ) {
            return atomic_store_explicit(
                &a->repr, v, order
            );
        }
        /* ... */
    );
}

function atomicGccFor(T) {
    return quote(
        typedef struct {
            ${T} repr;
        } atomic_${T};
        static inline atomic_${T} atomic_new_${T}(
            ${T} v
        ) {
            atomic_${T} a;
            a.repr = v;
            return a;
        }
        static inline ${T} atomic_load_${T}(
            const atomic_${T}* a,
            int order
        ) {
            ${T} res;
            __atomic_load(&a->repr, &res, order);
            return res;
        }
        static inline void atomic_store_${T}(
            atomic_${T}* a, ${T} v,
            int order
        ) {
            __atomic_store(&a->repr, &v, order);
        }
        /* ... */
    )
}

async function atomicFor(T) {
    if (await has_std_atomics) {
        return atomicStdFor(T);
    } else {
        return atomicGccFor(T);
    }
}

#include <stdio.h>
#include <stdint.h>
typedef uint32_t u32;
${atomicFor(quote(u32))};

int main(void) {
    atomic_u32 i = atomic_new_u32(79);
    printf("i = %d\n",
        (int)atomic_load_u32(&i, Relaxed)
    );
    atomic_store_u32(&i, 95, SeqCst);
    printf("i = %d\n",
        (int)atomic_load_u32(&i, AcqRel)
    );
}
```
```
i = 79
i = 95
```

## Miscellaneous
```cpp
#include <stdio.h>

typedef struct {
    int v;
} Integer;

int main(void) {
    const char* str = ${`
        multi-line strings in C!
    `};
    printf("%s\n", str);
    Integer i = (Integer)${{v: 10005}};
    printf("JS objects -> C99 struct literals: %d\n", i.v);
}
```
```

        multi-line strings in C!

JS objects -> C99 struct literals: 10005
```

# Compilation
Go 1.13 or newer is required.
```
git clone https://github.com/thooton/jc
cd jc/src
go build .
```

# Todo list
- Module system: one can `inquire` another JC file, which consists of including all C in that file, as well as returning the value of the other file's JC_EXPORTS global variable. `const atomic = inquire "atomic.jh";`
- C++ support: when a flag is passed to the program, `class` should mean C++, not JS class. Create a keyword `jclass` to represent JS classes. Alternatively, do it the other way around, and have the `cppclass` keyword.

# Contributing
Any and all contributions are greatly appreciated. Just send a PR :)

# License
MIT
