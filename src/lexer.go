// Code generated by re2c 3.0 on Sun Jan  8 14:13:22 2023, DO NOT EDIT.
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
    
{
	var yych byte
	yyaccept := 0
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00:
		goto yy1
	case 0x01,0x02,0x03,0x04,0x05,0x06,0x07,0x08,'\t':
		fallthrough
	case '\v','\f','\r',0x0E,0x0F,0x10,0x11,0x12,0x13,0x14,0x15,0x16,0x17,0x18,0x19,0x1A,0x1B,0x1C,0x1D,0x1E,0x1F,' ':
		fallthrough
	case 0x85:
		goto yy2
	case '\n':
		goto yy4
	case '!':
		goto yy5
	case '"':
		goto yy7
	case '#':
		goto yy9
	case '$':
		goto yy10
	case '%':
		fallthrough
	case '^':
		goto yy11
	case '&':
		goto yy12
	case '\'':
		goto yy13
	case '(':
		goto yy15
	case ')':
		goto yy16
	case '*':
		goto yy17
	case '+':
		goto yy18
	case ',':
		fallthrough
	case '~':
		goto yy19
	case '-':
		goto yy20
	case '.':
		goto yy21
	case '/':
		goto yy22
	case '0':
		goto yy24
	case '1','2','3','4','5','6','7','8','9':
		goto yy25
	case ':':
		goto yy27
	case ';':
		goto yy28
	case '<':
		goto yy29
	case '=':
		goto yy30
	case '>':
		goto yy32
	case '?':
		goto yy33
	case '@':
		fallthrough
	case 0x7F:
		goto yy34
	case 'L':
		fallthrough
	case 'U':
		goto yy38
	case '[':
		goto yy39
	case '\\':
		goto yy40
	case ']':
		goto yy41
	case '`':
		goto yy42
	case 'a':
		goto yy43
	case 'c':
		goto yy44
	case 'e':
		goto yy45
	case 'f':
		goto yy46
	case 'l':
		goto yy47
	case 'q':
		goto yy48
	case 's':
		goto yy49
	case 'u':
		goto yy50
	case 'v':
		goto yy51
	case 'w':
		goto yy52
	case '{':
		goto yy53
	case '|':
		goto yy54
	case '}':
		goto yy55
	default:
		goto yy35
	}
yy1:
	cursor++
	{
            return LexToken{0, 0, TEof}
        }
yy2:
	yyaccept = 0
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x01,0x02,0x03,0x04,0x05,0x06,0x07,0x08,'\t':
		fallthrough
	case '\v','\f','\r',0x0E,0x0F,0x10,0x11,0x12,0x13,0x14,0x15,0x16,0x17,0x18,0x19,0x1A,0x1B,0x1C,0x1D,0x1E,0x1F,' ':
		fallthrough
	case 0x85:
		goto yy2
	case '/':
		goto yy56
	case '\\':
		goto yy58
	default:
		goto yy3
	}
yy3:
	{
            began = cursor
            goto lex_start
        }
yy4:
	cursor++
	{
            return LexToken{uint32(began), uint32(cursor), TPLNewline}
        }
yy5:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '=':
		goto yy59
	default:
		goto yy6
	}
yy6:
	{
            return LexToken{uint32(began), uint32(cursor), TOther}
        }
yy7:
	yyaccept = 1
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00:
		fallthrough
	case '\n':
		goto yy8
	default:
		goto yy62
	}
yy8:
	{
            return LexToken{0, 0, TInvalid}
        }
yy9:
	cursor++
	{
            return LexToken{uint32(began), uint32(cursor), TPLPound}
        }
yy10:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '{':
		goto yy65
	default:
		goto yy8
	}
yy11:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '=':
		goto yy66
	default:
		goto yy6
	}
yy12:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '&':
		fallthrough
	case '=':
		goto yy66
	default:
		goto yy6
	}
yy13:
	yyaccept = 2
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00:
		fallthrough
	case '\n':
		goto yy14
	case '\'':
		goto yy25
	case '.':
		goto yy69
	case '0','1','2','3','4','5','6','7','8','9':
		goto yy13
	case 'E':
		fallthrough
	case 'e':
		goto yy70
	case 'L':
		fallthrough
	case 'U':
		fallthrough
	case 'Z':
		fallthrough
	case 'l':
		fallthrough
	case 'u':
		fallthrough
	case 'z':
		goto yy71
	case '\\':
		goto yy72
	case 'n':
		goto yy73
	default:
		goto yy67
	}
yy14:
	{
            return LexToken{uint32(began), uint32(cursor), TPLValue}
        }
yy15:
	cursor++
	{
            return LexToken{uint32(began), uint32(cursor), TLCircle}
        }
yy16:
	cursor++
	{
            return LexToken{uint32(began), uint32(cursor), TRCircle}
        }
yy17:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '*':
		goto yy59
	case '=':
		goto yy66
	default:
		goto yy6
	}
yy18:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '+':
		goto yy74
	case '=':
		goto yy66
	default:
		goto yy6
	}
yy19:
	cursor++
	goto yy6
yy20:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '-':
		goto yy75
	case '=':
		goto yy66
	case '>':
		goto yy76
	default:
		goto yy6
	}
yy21:
	yyaccept = 3
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1','2','3','4','5','6','7','8','9':
		goto yy77
	case '.':
		goto yy78
	default:
		goto yy6
	}
yy22:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '*':
		goto yy79
	case '/':
		goto yy80
	case '=':
		goto yy66
	default:
		goto yy23
	}
yy23:
	{
            return LexToken{uint32(began), uint32(cursor), TPLSlash}
        }
yy24:
	yyaccept = 2
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'B':
		fallthrough
	case 'b':
		goto yy81
	case 'O':
		fallthrough
	case 'o':
		goto yy84
	case 'X':
		fallthrough
	case 'x':
		goto yy85
	default:
		goto yy26
	}
yy25:
	yyaccept = 2
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
yy26:
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1','2','3','4','5','6','7','8','9':
		goto yy25
	case '.':
		goto yy77
	case 'E':
		fallthrough
	case 'e':
		goto yy82
	case 'L':
		fallthrough
	case 'U':
		fallthrough
	case 'Z':
		fallthrough
	case 'l':
		fallthrough
	case 'u':
		fallthrough
	case 'z':
		goto yy83
	case 'n':
		goto yy63
	default:
		goto yy14
	}
yy27:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case ':':
		goto yy66
	default:
		goto yy6
	}
yy28:
	cursor++
	{
            return LexToken{uint32(began), uint32(cursor), TSemi}
        }
yy29:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '<':
		goto yy59
	case '=':
		goto yy86
	default:
		goto yy6
	}
yy30:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '=':
		goto yy59
	default:
		goto yy31
	}
yy31:
	{
            return LexToken{uint32(began), uint32(cursor), TEq}
        }
yy32:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '=':
		goto yy59
	case '>':
		goto yy87
	default:
		goto yy6
	}
yy33:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '?':
		goto yy59
	default:
		goto yy6
	}
yy34:
	cursor++
	goto yy8
yy35:
	cursor++
	yych = peek_next(input, cursor)
yy36:
	switch (yych) {
	case 0x00,0x01,0x02,0x03,0x04,0x05,0x06,0x07,0x08,'\t','\n','\v','\f','\r',0x0E,0x0F,0x10,0x11,0x12,0x13,0x14,0x15,0x16,0x17,0x18,0x19,0x1A,0x1B,0x1C,0x1D,0x1E,0x1F,' ','!','"','#','$','%','&','\'','(',')','*','+',',','-','.','/':
		fallthrough
	case ':',';','<','=','>','?','@':
		fallthrough
	case '[','\\',']','^':
		fallthrough
	case '`':
		fallthrough
	case '{','|','}','~',0x7F:
		fallthrough
	case 0x85:
		goto yy37
	default:
		goto yy35
	}
yy37:
	{
            return LexToken{uint32(began), uint32(cursor), TIdentifier}
        }
yy38:
	yyaccept = 4
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case '"':
		goto yy61
	case '\'':
		goto yy67
	default:
		goto yy36
	}
yy39:
	cursor++
	{
            return LexToken{uint32(began), uint32(cursor), TLSquare}
        }
yy40:
	yyaccept = 1
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\n':
		goto yy2
	case '\r':
		goto yy88
	default:
		goto yy8
	}
yy41:
	cursor++
	{
            return LexToken{uint32(began), uint32(cursor), TRSquare}
        }
yy42:
	cursor++
	{
            return LexToken{uint32(began), uint32(cursor), TTick}
        }
yy43:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 's':
		goto yy89
	default:
		goto yy36
	}
yy44:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'l':
		goto yy90
	case 'o':
		goto yy91
	default:
		goto yy36
	}
yy45:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'n':
		goto yy92
	default:
		goto yy36
	}
yy46:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'o':
		goto yy93
	case 'u':
		goto yy94
	default:
		goto yy36
	}
yy47:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'e':
		goto yy95
	default:
		goto yy36
	}
yy48:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'u':
		goto yy96
	default:
		goto yy36
	}
yy49:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 't':
		goto yy97
	default:
		goto yy36
	}
yy50:
	yyaccept = 4
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case '"':
		goto yy61
	case '\'':
		goto yy67
	case '8':
		goto yy38
	case 'n':
		goto yy98
	default:
		goto yy36
	}
yy51:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'a':
		goto yy99
	default:
		goto yy36
	}
yy52:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'h':
		goto yy100
	default:
		goto yy36
	}
yy53:
	cursor++
	{
            return LexToken{uint32(began), uint32(cursor), TLBrace}
        }
yy54:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '=':
		fallthrough
	case '|':
		goto yy66
	default:
		goto yy6
	}
yy55:
	cursor++
	{
            return LexToken{uint32(began), uint32(cursor), TRBrace}
        }
yy56:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '/':
		goto yy80
	default:
		goto yy57
	}
yy57:
	cursor = marker
	switch (yyaccept) {
	case 0:
		goto yy3
	case 1:
		goto yy8
	case 2:
		goto yy14
	case 3:
		goto yy6
	default:
		goto yy37
	}
yy58:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\n':
		goto yy2
	case '\r':
		goto yy88
	default:
		goto yy57
	}
yy59:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '=':
		goto yy66
	default:
		goto yy60
	}
yy60:
	{
            return LexToken{uint32(began), uint32(cursor), TOther}
        }
yy61:
	cursor++
	yych = peek_next(input, cursor)
yy62:
	switch (yych) {
	case 0x00:
		fallthrough
	case '\n':
		goto yy57
	case '"':
		goto yy63
	case '\\':
		goto yy64
	default:
		goto yy61
	}
yy63:
	cursor++
	goto yy14
yy64:
	cursor++
	yych = peek_next(input, cursor)
	if (yych <= 0x00) {
		goto yy57
	}
	goto yy61
yy65:
	cursor++
	{
            return LexToken{uint32(began), uint32(cursor), TDollarLBrace}
        }
yy66:
	cursor++
	goto yy60
yy67:
	cursor++
	yych = peek_next(input, cursor)
yy68:
	switch (yych) {
	case 0x00:
		fallthrough
	case '\n':
		goto yy57
	case '\'':
		goto yy63
	case '\\':
		goto yy72
	default:
		goto yy67
	}
yy69:
	yyaccept = 2
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00:
		fallthrough
	case '\n':
		goto yy14
	case '\'':
		goto yy77
	case '0','1','2','3','4','5','6','7','8','9':
		goto yy69
	case 'B':
		goto yy101
	case 'E':
		fallthrough
	case 'e':
		goto yy70
	case 'F':
		fallthrough
	case 'f':
		goto yy102
	case 'L':
		fallthrough
	case 'l':
		goto yy73
	case '\\':
		goto yy72
	case 'b':
		goto yy103
	default:
		goto yy67
	}
yy70:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		goto yy104
	case '+':
		fallthrough
	case '-':
		goto yy105
	case '0','1','2','3','4','5','6','7','8','9':
		goto yy106
	default:
		goto yy68
	}
yy71:
	yyaccept = 2
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00:
		fallthrough
	case '\n':
		goto yy14
	case '\'':
		goto yy63
	case 'L':
		fallthrough
	case 'U':
		fallthrough
	case 'Z':
		fallthrough
	case 'l':
		fallthrough
	case 'u':
		fallthrough
	case 'z':
		goto yy71
	case '\\':
		goto yy72
	default:
		goto yy67
	}
yy72:
	cursor++
	yych = peek_next(input, cursor)
	if (yych <= 0x00) {
		goto yy57
	}
	goto yy67
yy73:
	yyaccept = 2
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00:
		fallthrough
	case '\n':
		goto yy14
	default:
		goto yy68
	}
yy74:
	cursor++
	{
            return LexToken{uint32(began), uint32(cursor), TPLInc}
        }
yy75:
	cursor++
	{
            return LexToken{uint32(began), uint32(cursor), TPLDec}
        }
yy76:
	cursor++
	{
            return LexToken{uint32(began), uint32(cursor), TMinusGt}
        }
yy77:
	yyaccept = 2
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1','2','3','4','5','6','7','8','9':
		goto yy77
	case 'B':
		goto yy107
	case 'E':
		fallthrough
	case 'e':
		goto yy82
	case 'F':
		fallthrough
	case 'f':
		goto yy108
	case 'L':
		fallthrough
	case 'l':
		goto yy63
	case 'b':
		goto yy109
	default:
		goto yy14
	}
yy78:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '.':
		goto yy66
	default:
		goto yy57
	}
yy79:
	cursor++
	{
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
yy80:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00:
		fallthrough
	case '\n':
		goto yy3
	case '\\':
		goto yy110
	default:
		goto yy80
	}
yy81:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1':
		goto yy111
	default:
		goto yy57
	}
yy82:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1','2','3','4','5','6','7','8','9':
		goto yy104
	case '+':
		fallthrough
	case '-':
		goto yy112
	default:
		goto yy57
	}
yy83:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'L':
		fallthrough
	case 'U':
		fallthrough
	case 'Z':
		fallthrough
	case 'l':
		fallthrough
	case 'u':
		fallthrough
	case 'z':
		goto yy83
	default:
		goto yy14
	}
yy84:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1','2','3','4','5','6','7','8','9':
		goto yy113
	default:
		goto yy57
	}
yy85:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1','2','3','4','5','6','7','8','9':
		fallthrough
	case 'A','B','C','D','E','F':
		fallthrough
	case 'a','b','c','d','e','f':
		goto yy114
	case '.':
		goto yy115
	default:
		goto yy57
	}
yy86:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '=','>':
		goto yy66
	default:
		goto yy60
	}
yy87:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '=':
		goto yy66
	case '>':
		goto yy59
	default:
		goto yy60
	}
yy88:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\n':
		goto yy2
	default:
		goto yy57
	}
yy89:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'y':
		goto yy116
	default:
		goto yy36
	}
yy90:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'a':
		goto yy117
	default:
		goto yy36
	}
yy91:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'n':
		goto yy118
	default:
		goto yy36
	}
yy92:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'u':
		goto yy119
	default:
		goto yy36
	}
yy93:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'r':
		goto yy120
	default:
		goto yy36
	}
yy94:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'n':
		goto yy122
	default:
		goto yy36
	}
yy95:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 't':
		goto yy123
	default:
		goto yy36
	}
yy96:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'o':
		goto yy125
	default:
		goto yy36
	}
yy97:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'r':
		goto yy126
	default:
		goto yy36
	}
yy98:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'd':
		goto yy127
	case 'i':
		goto yy128
	default:
		goto yy36
	}
yy99:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'r':
		goto yy123
	default:
		goto yy36
	}
yy100:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'i':
		goto yy129
	default:
		goto yy36
	}
yy101:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'F':
		goto yy130
	default:
		goto yy68
	}
yy102:
	yyaccept = 2
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00:
		fallthrough
	case '\n':
		goto yy14
	case '1':
		goto yy131
	case '3':
		goto yy132
	case '6':
		goto yy133
	default:
		goto yy68
	}
yy103:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'f':
		goto yy130
	default:
		goto yy68
	}
yy104:
	yyaccept = 2
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1','2','3','4','5','6','7','8','9':
		goto yy104
	case 'B':
		goto yy107
	case 'F':
		fallthrough
	case 'f':
		goto yy108
	case 'L':
		fallthrough
	case 'l':
		goto yy63
	case 'b':
		goto yy109
	default:
		goto yy14
	}
yy105:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		goto yy104
	case '0','1','2','3','4','5','6','7','8','9':
		goto yy106
	default:
		goto yy68
	}
yy106:
	yyaccept = 2
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00:
		fallthrough
	case '\n':
		goto yy14
	case '\'':
		goto yy104
	case '0','1','2','3','4','5','6','7','8','9':
		goto yy106
	case 'B':
		goto yy101
	case 'F':
		fallthrough
	case 'f':
		goto yy102
	case 'L':
		fallthrough
	case 'l':
		goto yy73
	case '\\':
		goto yy72
	case 'b':
		goto yy103
	default:
		goto yy67
	}
yy107:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'F':
		goto yy134
	default:
		goto yy57
	}
yy108:
	yyaccept = 2
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case '1':
		goto yy135
	case '3':
		goto yy136
	case '6':
		goto yy137
	default:
		goto yy14
	}
yy109:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'f':
		goto yy134
	default:
		goto yy57
	}
yy110:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00:
		goto yy3
	case '\n':
		goto yy2
	case '\r':
		goto yy138
	case '\\':
		goto yy110
	default:
		goto yy80
	}
yy111:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1':
		goto yy111
	case 'L':
		fallthrough
	case 'U':
		fallthrough
	case 'Z':
		fallthrough
	case 'l':
		fallthrough
	case 'u':
		fallthrough
	case 'z':
		goto yy83
	case 'n':
		goto yy63
	default:
		goto yy14
	}
yy112:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1','2','3','4','5','6','7','8','9':
		goto yy104
	default:
		goto yy57
	}
yy113:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1','2','3','4','5','6','7','8','9':
		goto yy113
	case 'L':
		fallthrough
	case 'U':
		fallthrough
	case 'Z':
		fallthrough
	case 'l':
		fallthrough
	case 'u':
		fallthrough
	case 'z':
		goto yy83
	case 'n':
		goto yy63
	default:
		goto yy14
	}
yy114:
	yyaccept = 2
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1','2','3','4','5','6','7','8','9':
		fallthrough
	case 'A','B','C','D','E','F':
		fallthrough
	case 'a','b','c','d','e','f':
		goto yy114
	case '.':
		goto yy139
	case 'L':
		fallthrough
	case 'U':
		fallthrough
	case 'Z':
		fallthrough
	case 'l':
		fallthrough
	case 'u':
		fallthrough
	case 'z':
		goto yy83
	case 'P':
		fallthrough
	case 'p':
		goto yy141
	case 'n':
		goto yy63
	default:
		goto yy14
	}
yy115:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'P':
		fallthrough
	case 'p':
		goto yy57
	default:
		goto yy140
	}
yy116:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'n':
		goto yy142
	default:
		goto yy36
	}
yy117:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 's':
		goto yy143
	default:
		goto yy36
	}
yy118:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 's':
		goto yy144
	default:
		goto yy36
	}
yy119:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'm':
		goto yy145
	default:
		goto yy36
	}
yy120:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00,0x01,0x02,0x03,0x04,0x05,0x06,0x07,0x08,'\t','\n','\v','\f','\r',0x0E,0x0F,0x10,0x11,0x12,0x13,0x14,0x15,0x16,0x17,0x18,0x19,0x1A,0x1B,0x1C,0x1D,0x1E,0x1F,' ','!','"','#','$','%','&','\'','(',')','*','+',',','-','.','/':
		fallthrough
	case ':',';','<','=','>','?','@':
		fallthrough
	case '[','\\',']','^':
		fallthrough
	case '`':
		fallthrough
	case '{','|','}','~',0x7F:
		fallthrough
	case 0x85:
		goto yy121
	default:
		goto yy35
	}
yy121:
	{
            return LexToken{uint32(began), uint32(cursor), TFunctionForWhileClass}
        }
yy122:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'c':
		goto yy147
	default:
		goto yy36
	}
yy123:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00,0x01,0x02,0x03,0x04,0x05,0x06,0x07,0x08,'\t','\n','\v','\f','\r',0x0E,0x0F,0x10,0x11,0x12,0x13,0x14,0x15,0x16,0x17,0x18,0x19,0x1A,0x1B,0x1C,0x1D,0x1E,0x1F,' ','!','"','#','$','%','&','\'','(',')','*','+',',','-','.','/':
		fallthrough
	case ':',';','<','=','>','?','@':
		fallthrough
	case '[','\\',']','^':
		fallthrough
	case '`':
		fallthrough
	case '{','|','}','~',0x7F:
		fallthrough
	case 0x85:
		goto yy124
	default:
		goto yy35
	}
yy124:
	{
            return LexToken{uint32(began), uint32(cursor), TJsOnly}
        }
yy125:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 't':
		goto yy148
	default:
		goto yy36
	}
yy126:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'u':
		goto yy149
	default:
		goto yy36
	}
yy127:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'e':
		goto yy150
	default:
		goto yy36
	}
yy128:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'o':
		goto yy151
	default:
		goto yy36
	}
yy129:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'l':
		goto yy152
	default:
		goto yy36
	}
yy130:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '1':
		goto yy153
	default:
		goto yy68
	}
yy131:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '2':
		goto yy154
	case '6':
		goto yy73
	default:
		goto yy68
	}
yy132:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '2':
		goto yy73
	default:
		goto yy68
	}
yy133:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '4':
		goto yy73
	default:
		goto yy68
	}
yy134:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '1':
		goto yy155
	default:
		goto yy57
	}
yy135:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '2':
		goto yy156
	case '6':
		goto yy63
	default:
		goto yy57
	}
yy136:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '2':
		goto yy63
	default:
		goto yy57
	}
yy137:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '4':
		goto yy63
	default:
		goto yy57
	}
yy138:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00:
		goto yy3
	case '\n':
		goto yy2
	case '\\':
		goto yy110
	default:
		goto yy80
	}
yy139:
	cursor++
	yych = peek_next(input, cursor)
yy140:
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1','2','3','4','5','6','7','8','9':
		fallthrough
	case 'A','B','C','D','E','F':
		fallthrough
	case 'a','b','c','d','e','f':
		goto yy139
	case 'P':
		fallthrough
	case 'p':
		goto yy141
	default:
		goto yy57
	}
yy141:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1','2','3','4','5','6','7','8','9':
		fallthrough
	case 'A','B','C','D','E','F':
		fallthrough
	case 'a','b','c','d','e','f':
		goto yy157
	case '+':
		fallthrough
	case '-':
		goto yy158
	default:
		goto yy57
	}
yy142:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'c':
		goto yy123
	default:
		goto yy36
	}
yy143:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 's':
		goto yy120
	default:
		goto yy36
	}
yy144:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 't':
		goto yy159
	default:
		goto yy36
	}
yy145:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00,0x01,0x02,0x03,0x04,0x05,0x06,0x07,0x08,'\t','\n','\v','\f','\r',0x0E,0x0F,0x10,0x11,0x12,0x13,0x14,0x15,0x16,0x17,0x18,0x19,0x1A,0x1B,0x1C,0x1D,0x1E,0x1F,' ','!','"','#','$','%','&','\'','(',')','*','+',',','-','.','/':
		fallthrough
	case ':',';','<','=','>','?','@':
		fallthrough
	case '[','\\',']','^':
		fallthrough
	case '`':
		fallthrough
	case '{','|','}','~',0x7F:
		fallthrough
	case 0x85:
		goto yy146
	default:
		goto yy35
	}
yy146:
	{
            return LexToken{uint32(began), uint32(cursor), TStructUnionEnum}
        }
yy147:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 't':
		goto yy161
	default:
		goto yy36
	}
yy148:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'e':
		goto yy162
	default:
		goto yy36
	}
yy149:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'c':
		goto yy164
	default:
		goto yy36
	}
yy150:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'f':
		goto yy165
	default:
		goto yy36
	}
yy151:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'n':
		goto yy145
	default:
		goto yy36
	}
yy152:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'e':
		goto yy120
	default:
		goto yy36
	}
yy153:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '6':
		goto yy73
	default:
		goto yy68
	}
yy154:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '8':
		goto yy73
	default:
		goto yy68
	}
yy155:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '6':
		goto yy63
	default:
		goto yy57
	}
yy156:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '8':
		goto yy63
	default:
		goto yy57
	}
yy157:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1','2','3','4','5','6','7','8','9':
		fallthrough
	case 'A','B','C','D','E','F':
		fallthrough
	case 'a','b','c','d','e','f':
		goto yy157
	case 'L':
		fallthrough
	case 'l':
		goto yy63
	default:
		goto yy14
	}
yy158:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case '\'':
		fallthrough
	case '0','1','2','3','4','5','6','7','8','9':
		fallthrough
	case 'A','B','C','D','E','F':
		fallthrough
	case 'a','b','c','d','e','f':
		goto yy157
	default:
		goto yy57
	}
yy159:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00,0x01,0x02,0x03,0x04,0x05,0x06,0x07,0x08,'\t','\n','\v','\f','\r',0x0E,0x0F,0x10,0x11,0x12,0x13,0x14,0x15,0x16,0x17,0x18,0x19,0x1A,0x1B,0x1C,0x1D,0x1E,0x1F,' ','!','"','#','$','%','&','\'','(',')','*','+',',','-','.','/':
		fallthrough
	case ':',';','<','=','>','?','@':
		fallthrough
	case '[','\\',']','^':
		fallthrough
	case '`':
		fallthrough
	case '{','|','}','~',0x7F:
		fallthrough
	case 0x85:
		goto yy160
	default:
		goto yy35
	}
yy160:
	{ 
            return LexToken{uint32(began), uint32(cursor), TConst} 
        }
yy161:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'i':
		goto yy166
	default:
		goto yy36
	}
yy162:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00,0x01,0x02,0x03,0x04,0x05,0x06,0x07,0x08,'\t','\n','\v','\f','\r',0x0E,0x0F,0x10,0x11,0x12,0x13,0x14,0x15,0x16,0x17,0x18,0x19,0x1A,0x1B,0x1C,0x1D,0x1E,0x1F,' ','!','"','#','$','%','&','\'','(',')','*','+',',','-','.','/':
		fallthrough
	case ':',';','<','=','>','?','@':
		fallthrough
	case '[','\\',']','^':
		fallthrough
	case '`':
		fallthrough
	case '{','|','}','~',0x7F:
		fallthrough
	case 0x85:
		goto yy163
	default:
		goto yy35
	}
yy163:
	{
            return LexToken{uint32(began), uint32(cursor), TQuote}
        }
yy164:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 't':
		goto yy145
	default:
		goto yy36
	}
yy165:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'i':
		goto yy167
	default:
		goto yy36
	}
yy166:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'o':
		goto yy168
	default:
		goto yy36
	}
yy167:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'n':
		goto yy169
	default:
		goto yy36
	}
yy168:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'n':
		goto yy120
	default:
		goto yy36
	}
yy169:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'e':
		goto yy170
	default:
		goto yy36
	}
yy170:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'd':
		goto yy123
	default:
		goto yy36
	}
}

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
    
{
	var yych byte
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00:
		goto yy172
	case '/':
		goto yy175
	default:
		goto yy173
	}
yy172:
	cursor++
	{
            return LexToken{0, 0, TEof}
        }
yy173:
	cursor++
yy174:
	{
            return LexToken{0, 0, TInvalid}
        }
yy175:
	cursor++
	marker = cursor
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00:
		fallthrough
	case '\n':
		fallthrough
	case '/':
		goto yy174
	default:
		goto yy177
	}
yy176:
	cursor++
	yych = peek_next(input, cursor)
yy177:
	switch (yych) {
	case 0x00:
		fallthrough
	case '\n':
		goto yy178
	case '/':
		goto yy181
	case '[':
		goto yy179
	case '\\':
		goto yy180
	default:
		goto yy176
	}
yy178:
	cursor = marker
	goto yy174
yy179:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 0x00:
		fallthrough
	case '\n':
		fallthrough
	case '/':
		goto yy178
	case '\\':
		goto yy183
	case ']':
		goto yy176
	default:
		goto yy179
	}
yy180:
	cursor++
	yych = peek_next(input, cursor)
	if (yych <= 0x00) {
		goto yy178
	}
	goto yy176
yy181:
	cursor++
	yych = peek_next(input, cursor)
	switch (yych) {
	case 'd':
		fallthrough
	case 'g':
		fallthrough
	case 'i':
		fallthrough
	case 'm':
		fallthrough
	case 's':
		fallthrough
	case 'u','v':
		fallthrough
	case 'y':
		goto yy181
	default:
		goto yy182
	}
yy182:
	{
            return LexToken{uint32(began), uint32(cursor), TOther}
        }
yy183:
	cursor++
	yych = peek_next(input, cursor)
	if (yych <= 0x00) {
		goto yy178
	}
	goto yy179
}

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