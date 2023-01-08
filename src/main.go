

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

/*
func debugNode(n *AstNode, input string) {
	put := func(s string) {
		fmt.Print(" " + s + " ")
	}
	cont := func() {
		for i := range n.nodes {
			debugNode(&n.nodes[i], input)
		}
	}
	lit := func() {
		put(input[n.token.begin:n.token.end])
	}
	lst := func(name string) {
		put(name + ":")
		cont()
		put(":" + name + " end")
	}
	tok := func(name string) {
		put("[" + name)
		lit()
		put("]")
	}
	switch n.token.kind {
	case AstBlockUnfin:
		lst("unfinished block")
		break
	case AstBlock:
		lst("block")
		break
	case AstQuoteUnfin:
		lst("unfinished quote")
		break
	case AstQuote:
		lst("quote")
		break
	case AstTicklitUnfin:
		lst("unfinished ticklit")
		break
	case AstTicklit:
		lst("ticklit")
		break
	case AstInterpolUnfin:
		lst("unfinished interpol")
		break
	case AstInterpol:
		lst("interpol")
		break
	case AstToplevel:
		lst("toplevel")
		break
	case AstMacro:
		lst("macro")
		break
	case AstLfnArgs:
		lst("lfn args")
		break
	case AstLfnRetType:
		lst("lfn ret type")
		break
	case AstLfnBody:
		lst("lfn body")
		break
	case TMacroStart:
		tok("macro start")
		break
	case TIdentifier:
		tok("identifier")
		break
	case TConst:
		tok("const")
		break
	case TStructUnionEnum:
		tok("seu")
		break
	case TJsOnly:
		tok("jsonly")
		break
	case TMinusGt:
		tok("minusgt")
		break
	case TEq:
		tok("eq")
		break
	case TLSquare:
		tok("lsquare")
		break
	case TRSquare:
		tok("rsquare")
		break
	case TOther:
		tok("other")
		break
	case TLBrace:
		tok("lbrace")
		break
	case TRBrace:
		tok("rbrace")
		break
	case TQuote:
		tok("quote")
		break
	case TDollarLBrace:
		tok("dollarlbrace")
		break
	case TTick:
		tok("tick")
		break
	case TLCircle:
		tok("lcircle")
		break
	case TRCircle:
		tok("rcircle")
		break
	case TSemi:
		tok("semi")
		break
	}
}*/

func main() {
	if len(os.Args) != 7 ||
		os.Args[1] != "-i" ||
		os.Args[3] != "-o" ||
		os.Args[5] != "-n" {
		fmt.Println("usage: jc -i <infile> -o <outfile> -n <js interpreter>")
		os.Exit(1)
	}
	astErrorVerbose = true
	infile := os.Args[2]
	outfile := os.Args[4]
	interpreter := os.Args[6]
	input_bytes, err := os.ReadFile(infile)
	if err != nil {
		fmt.Printf("can't read file %s: %v\n", infile, err)
		os.Exit(1)
	}
	input := string(input_bytes)
	root, estr := parseIntoAst(input)
	if estr != "" {
		fmt.Printf("parsing error: %s\n", estr)
		os.Exit(1)
	}
	//debugNode(&root, input)
	result := codegenPerform(&root, input, outfile)
	dotidx := strings.LastIndexByte(infile, '.')
	var tmpfile string
	if dotidx != -1 {
		tmpfile = infile[0:dotidx] + "_" + infile[dotidx+1:] + "_gen.js"
	} else {
		tmpfile = infile + "_gen.js"
	}
	err = os.WriteFile(tmpfile, result, 0o644)
	if err != nil {
		fmt.Printf("can't write to %s: %v\n", tmpfile, err)
		os.Exit(1)
	}
	_ = os.Remove(outfile)
	cmd := exec.Command(interpreter, tmpfile)
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		fmt.Printf("can't start interpreter %s: %v\n", interpreter, err)
		os.Exit(1)
	}
	_ = cmd.Wait()
	_, err = os.Lstat(outfile)
	if err != nil {
		os.Exit(1)
	} else {
		_ = os.Remove(tmpfile)
	}
}
