%{
    package main
    const (
        AstBlockUnfin       = 10
        AstBlock            = 11
        AstQCircleUnfin     = 12
        AstQCircle          = 13
        AstQuoteUnfin       = 14
        AstQuote            = 15
        AstTicklitUnfin     = 16
        AstTicklit          = 17
        AstInterpolUnfin    = 18
        AstInterpol         = 19
        AstToplevel         = 20
        AstMacro            = 21
        AstLfnArgs          = 22
        AstLfnRetType       = 23
        AstLfnBody          = 24
        AstEnd              = 25
    )
    /* this is a "union" between
        token and nodes. all of
        the actual token constants are enormous,
        so if token.kind < AstEnd, we know it is an
        []AstNode instead. */
    type AstNode struct {
        token LexToken
        nodes []AstNode
    }
%}

%token TIdentifier
%token TConst TStructUnionEnumClass TFunctionForWhileJsclass TJsOnly
%token TMinusGt
%token TEq
%token TLSquare TRSquare
%token TOther

%token TLBrace TRBrace
%token TQuote TDollarLBrace TTick
%token TLCircle TRCircle
%token TSemi

%token TMacroStart TMacroEnd

%union {
    node AstNode
}

%type <node> TIdentifier TConst TStructUnionEnumClass TFunctionForWhileJsclass TJsOnly TMinusGt
%type <node> TEq TLSquare TRSquare TOther
%type <node> TLBrace TRBrace TQuote TDollarLBrace TTick TLCircle TRCircle
%type <node> TSemi
%type <node> TMacroStart TMacroEnd

%type <node> harmless block_unfinished block quote_harmless quote_circle_unfinished quote_circle quote_unfinished quote
%type <node> ticklit_unfinished ticklit interpol_unfinished interpol toplevel result
%type <node> tl_interpol_unfinished tl_interpol

%type <node> block_unfinished_harmless all_interpol_harmless toplevel_harmless
%type <node> macro_harmless macro_unfinished macro

%start result

/* note: compile with
    go run golang.org/x/tools/cmd/goyacc -o ast.go -p ast ast.y */

/* note: I manually change the parser parameter
    astLexer in the generated code to
    the custom ParserLexer defined in parser.go 
    */

%%

harmless
    : TIdentifier
    | TConst
    | TStructUnionEnumClass
    | TFunctionForWhileJsclass
    | TJsOnly
    | TMinusGt
    | TEq
    | TLSquare
    | TRSquare
    | TOther
    | TSemi
    ;

block_unfinished_harmless
    : harmless
    | TLCircle
    | TRCircle
    | quote
    | ticklit
    | interpol
    | block
    ;
block_unfinished
    : TLBrace { astlex.beginBlock(); $$ = AstNode{LexToken{0, 0, AstBlockUnfin}, []AstNode{$1}} }
    | block_unfinished block_unfinished_harmless { $1.nodes = append($1.nodes, $2); $$ = $1 }
    ;
block
    : block_unfinished TRBrace { astlex.endBlock(); $$ = AstNode{LexToken{0, 0, AstBlock}, append($1.nodes, $2)} }
    ;

quote_harmless
    : harmless
    | TLBrace
    | TRBrace
    | TQuote
    | TTick
    | interpol
    ;
quote_circle_unfinished
    : TLCircle { $$ = AstNode{LexToken{0, 0, AstQCircleUnfin}, []AstNode{$1}}}
    | quote_circle_unfinished quote_harmless { $1.nodes = append($1.nodes, $2); $$ = $1 }
    | quote_circle_unfinished quote_circle { $1.nodes = append($1.nodes, $2.nodes...); $$ = $1 }
    ;
quote_circle
    : quote_circle_unfinished TRCircle

quote_unfinished
    : TQuote TLCircle
        { $$ = AstNode{LexToken{0, 0, AstQuoteUnfin}, []AstNode{$1, $2}} }
    | quote_unfinished quote_harmless
        { $1.nodes = append($1.nodes, $2); $$ = $1 }
    | quote_unfinished quote_circle
        { $1.nodes = append($1.nodes, $2.nodes...); $$ = $1 }
    ;
quote
    : quote_unfinished TRCircle { $$ = AstNode{LexToken{0, 0, AstQuote}, append($1.nodes, $2)} }
    ;

all_interpol_harmless
    : harmless
    | TLCircle
    | TRCircle
    | quote
    | ticklit
    | interpol
    | block
    ;

tl_interpol_unfinished
    : TDollarLBrace { astlex.beginBlock(); astlex.outTicks(); $$ = AstNode{LexToken{0, 0, AstInterpolUnfin}, []AstNode{$1}} }
    | tl_interpol_unfinished all_interpol_harmless  { $1.nodes = append($1.nodes, $2); $$ = $1 }
    ;
tl_interpol
    : tl_interpol_unfinished TRBrace { astlex.endBlock(); astlex.inTicks(); $$ = AstNode{LexToken{0, 0, AstInterpol}, append($1.nodes, $2)} }
    ;
ticklit_unfinished
    : TTick { astlex.inTicks(); $$ = AstNode{LexToken{0, 0, AstTicklitUnfin}, []AstNode{$1}} }
    | ticklit_unfinished TOther        { $1.nodes = append($1.nodes, $2); $$ = $1 }
    | ticklit_unfinished tl_interpol   { $1.nodes = append($1.nodes, $2); $$ = $1 }
    ;
ticklit
    : ticklit_unfinished TTick
        { astlex.outTicks(); $$ = AstNode{LexToken{0, 0, AstTicklit}, append($1.nodes, $2)} }
    ;

interpol_unfinished
    : TDollarLBrace { astlex.beginBlock(); $$ = AstNode{LexToken{0, 0, AstInterpolUnfin}, []AstNode{$1}} }
    | interpol_unfinished all_interpol_harmless  { $1.nodes = append($1.nodes, $2); $$ = $1 }
    ;
interpol
    : interpol_unfinished TRBrace { astlex.endBlock(); $$ = AstNode{LexToken{0, 0, AstInterpol}, append($1.nodes, $2)} }
    ;

macro_harmless
    : harmless
    | TLBrace
    | TRBrace
    | TLCircle
    | TRCircle
    | TQuote
    | TTick
    | interpol
    ;
macro_unfinished
    : TMacroStart { $$ = AstNode{LexToken{0, 0, AstMacro}, []AstNode{$1}} }
    | macro_unfinished macro_harmless { $1.nodes = append($1.nodes, $2); $$ = $1 }
    ;
macro
    : macro_unfinished TMacroEnd

toplevel_harmless
    : harmless
    | TLCircle
    | TRCircle
    | quote
    | ticklit
    | interpol
    | block
    | macro
    ;
toplevel
    : toplevel_harmless  { $$ = AstNode{LexToken{0, 0, AstToplevel}, []AstNode{$1}} }
    | toplevel toplevel_harmless { $1.nodes = append($1.nodes, $2); $$ = $1 }
    ;

result
    : toplevel { astlex.receiveResult($1); return 0 }
    ;