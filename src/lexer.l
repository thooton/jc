package main

const (
    TEof        = 0
    TInvalid    = 1
    TPLBefore   = 2
    TPLNewline  = 3
    TPLValue    = 4
    TPLInc      = 5
    TPLDec      = 6
    TPLPound    = 7
    TPLSlash    = 8
    TPLAfter    = 9
)

type LexToken struct {
    begin, end, kind uint32
}

func lexNextToken(input string, cursor uint) LexToken {
	var began uint = cursor
    var marker uint = 0
    peek_next := func(str string, i uint) byte {
        if i < uint(len(str)) {
            return str[i]
        } else {
            return 0
        }
    }
lex_start:
    /*!re2c
        re2c:yyfill:enable = 0;
        re2c:define:YYCTYPE = byte;
        re2c:define:YYPEEK = "peek_next(input, cursor)";
        re2c:define:YYSKIP = "cursor++";
        re2c:define:YYBACKUP = "marker = cursor";
        re2c:define:YYRESTORE = "cursor = marker";
        re2c:tags = 0;
        
        * {
            return LexToken{0, 0, TInvalid}
        }
        [\x00] {
            return LexToken{0, 0, TEof}
        }
        mcm = "/" "*";
        mcm {
            for ;; {
                yych = peek_next(input, cursor)
                cursor++
                if yych == 0 {
                    return LexToken{0, 0, TEof}
                } else if yych == '*' {
                    yych = peek_next(input, cursor)
                    if yych == 0 {
                        return LexToken{0, 0, TEof}
                    } else if yych == '/' {
                        cursor++;
                        began = cursor
                        goto lex_start
                    }
                }
            }
        }
        scm = "//" [^\n\x00]*;
        wsp = ([\x01-\x09\x0B-\x20\x85] | scm | "\\\r\n" | "\\\n")+;
        wsp {
            began = cursor
            goto lex_start
        }
        
        "const" { 
            return LexToken{uint32(began), uint32(cursor), TConst} 
        }
        ("struct" | "union" | "enum") {
            return LexToken{uint32(began), uint32(cursor), TStructUnionEnum}
        }
        "quote" {
            return LexToken{uint32(began), uint32(cursor), TQuote}
        }
        "function" | "for" | "while" | "class" {
            return LexToken{uint32(began), uint32(cursor), TFunctionForWhileClass}
        }
        ("async" | "let" | "var" | "undefined") {
            return LexToken{uint32(began), uint32(cursor), TJsOnly}
        }
        "->" {
            return LexToken{uint32(began), uint32(cursor), TMinusGt}
        }
        "=" {
            return LexToken{uint32(began), uint32(cursor), TEq}
        }
        "[" {
            return LexToken{uint32(began), uint32(cursor), TLSquare}
        }
        "]" {
            return LexToken{uint32(began), uint32(cursor), TRSquare}
        }
        [!>.<+\-*%&|^?:,~] {
            return LexToken{uint32(began), uint32(cursor), TOther}
        }
        "{" {
            return LexToken{uint32(began), uint32(cursor), TLBrace}
        }
        "}" {
            return LexToken{uint32(began), uint32(cursor), TRBrace}
        }
        "quote" {
            return LexToken{uint32(began), uint32(cursor), TQuote}
        }
        "${" {
            return LexToken{uint32(began), uint32(cursor), TDollarLBrace}
        }
        "`" {
            return LexToken{uint32(began), uint32(cursor), TTick}
        }
        "(" {
            return LexToken{uint32(began), uint32(cursor), TLCircle}
        }
        ")" {
            return LexToken{uint32(began), uint32(cursor), TRCircle}
        }
        ";" {
            return LexToken{uint32(began), uint32(cursor), TSemi}
        }
        "\n" {
            return LexToken{uint32(began), uint32(cursor), TPLNewline}
        }

        id_start = [A-Za-z_\x80-\x84\x86-\xFF];
        id_continue = [A-Za-z0-9_\x80-\x84\x86-\xFF];
        identifier = id_start id_continue*;
        identifier {
            return LexToken{uint32(began), uint32(cursor), TIdentifier}
        }

        escapedchr = "\\" [^\x00];
        strpfx = "L"|("u" "8"?)|"U";
        dquotelit = strpfx? '"' (escapedchr | [^"\n\\\x00])* '"';
        squotelit = strpfx? "'" (escapedchr | [^'\n\\\x00])* "'";

        digit_sequence = [0-9']+;
        hex_digit_sequence = [0-9A-Fa-f']+;
        bin_digit_sequence = [01']+;

        declit = ("0" [oO])? digit_sequence;
        hexlit = ("0" [xX]) hex_digit_sequence;
        binlit = ("0" [bB]) bin_digit_sequence;
        intsfx = "n"|([lLuUzZ]+);
        intlit = (declit|hexlit|binlit) intsfx?;

        decimal_exponent = [eE] [+-]? digit_sequence;
        hex_exponent = [pP] [+-]? hex_digit_sequence;
        floatsfx = [flFL]|"f16"|"f32"|"f64"|"f128"|"bf16"|"F16"|"F32"|"F64"|"F128"|"BF16";
        floatlit = (
            (digit_sequence decimal_exponent)
            | (
                (
                    (digit_sequence '.')
                    | (digit_sequence? '.' digit_sequence)
                )
                decimal_exponent?
            )
            | ("0" [xX]
                (
                    (hex_digit_sequence '.'?)
                    | (hex_digit_sequence? '.' hex_digit_sequence)
                )
                hex_exponent
            )
        ) floatsfx?;

        value = dquotelit | squotelit
            | intlit | floatlit;
        value {
            return LexToken{uint32(began), uint32(cursor), TPLValue}
        }

        "++" {
            return LexToken{uint32(began), uint32(cursor), TPLInc}
        }
        "--" {
            return LexToken{uint32(began), uint32(cursor), TPLDec}
        }
        "#" {
            return LexToken{uint32(began), uint32(cursor), TPLPound}
        }
        "/" {
            return LexToken{uint32(began), uint32(cursor), TPLSlash}
        }

        coperator = ">>" | "<<" | "&&" | "||"
            | "==" | "!=" | ">=" | "<=" | "+=" | "-=" | "*="
            | "/=" | "%=" | "&=" | "|=" | "^="
            | ">>=" | "<<=" | "::" | "..." | "<=>";

        jsoperator = "===" | "!==" | ">==" | "<=="
        | ">>>" | ">>>=" | "??" | "**" | "**=" | "??=";
        
        other = coperator | jsoperator;
        other {
            return LexToken{uint32(began), uint32(cursor), TOther}
        }
    */
}

func lexRegexpLit(input string, cursor uint) LexToken {
    var began uint = cursor
    var marker uint = 0
    peek_next := func(str string, i uint) byte {
        if i < uint(len(str)) {
            return str[i]
        } else {
            return 0
        }
    }
    /*!re2c
        * {
            return LexToken{0, 0, TInvalid}
        }
        [\x00] {
            return LexToken{0, 0, TEof}
        }
        regchr = [^/\[\n\\\x00] | escapedchr;
        regseq = "[" ([^/\]\n\\\x00] | escapedchr)* "]";
        regexp = "/" (regchr | regseq)+ "/" [dgimsuvy]*;
        regexp {
            return LexToken{uint32(began), uint32(cursor), TOther}
        }
    */
}

func lexTickLit(input string, cursor uint) LexToken {
    var began uint = cursor
    var ch byte
    for ;; {
        if cursor < uint(len(input)) {
            ch = input[cursor]
            cursor++
        } else {
            return LexToken{0, 0, TInvalid}
        }
        if ch == '\\' {
            if cursor < uint(len(input)) {
                cursor++ 
            } else {
                return LexToken{0, 0, TInvalid}
            }
        } else if ch == '`' {
            if (began+1) == cursor {
                return LexToken{uint32(began), uint32(cursor), TTick}
            } else {
                return LexToken{uint32(began), uint32(cursor-1), TOther}
            }
        } else if (ch == '$') {
            if cursor < uint(len(input)) {
                ch = input[cursor]
                cursor++
            } else {
                return LexToken{0, 0, TInvalid}
            }
            if (ch == '{') {
                if (began+2) == cursor {
                    return LexToken{uint32(began), uint32(cursor), TDollarLBrace}
                } else {
                    return LexToken{uint32(began), uint32(cursor-2), TOther}
                }
            }
        }
    }
}