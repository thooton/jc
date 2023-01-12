package main

import (
	"strconv"
	"strings"
	"github.com/evanw/esbuild/pkg/api"
)

type ParserLexer struct {
	result AstNode
	input  []byte
	err    string
	//parser_ref  astParser
	cursor      uint
	prev_tok    uint32
	block_level uint32
	in_ticks    bool
	in_macro    bool
}

func plNew(input []byte /*, parser_ref astParser*/) ParserLexer {
	return ParserLexer{
		result: AstNode{LexToken{0, 0, 0}, nil},
		input:  input,
		err:    "",
		//parser_ref:  parser_ref,
		cursor:      0,
		prev_tok:    0,
		block_level: 0,
		in_ticks:    false,
		in_macro:    false,
	}
}

/*
func (pl *ParserLexer) ensureNoLookahead(ctx string) {
	yychar := pl.parser_ref.Lookahead()
	if yychar != -1 {
		fmt.Printf("parser looked ahead "+ctx+" and found %d, this should not happen\n", yychar)
		os.Exit(1)
	}
}*/

/*
	once an lbrace, rbrace, start btick, or end btick

appears in one of these respective situations,
it should be immediately used to reduce.
*/
func (pl *ParserLexer) beginBlock() {
	//pl.ensureNoLookahead("into block")
	pl.block_level++
}
func (pl *ParserLexer) endBlock() {
	//pl.ensureNoLookahead("after block")
	if pl.block_level != 0 {
		pl.block_level--
	}
}
func (pl *ParserLexer) inTicks() {
	//pl.ensureNoLookahead("into tick literal")
	pl.in_ticks = true
}
func (pl *ParserLexer) outTicks() {
	//pl.ensureNoLookahead("past tick literal")
	pl.in_ticks = false
}

func (pl *ParserLexer) receiveResult(result AstNode) {
	pl.result = result
}
func (pl *ParserLexer) Lex(lval *astSymType) int {
	var token LexToken
	for {
		if pl.in_ticks {
			token = lexTickLit(pl.input, pl.cursor)
		} else {
			token = lexNextToken(pl.input, pl.cursor)
		}
		if token.kind == TEof {
			if pl.in_macro {
				pl.in_macro = false
				lval.node = AstNode{LexToken{
					uint32(pl.cursor), uint32(pl.cursor), TMacroEnd,
				}, nil}
				return TMacroEnd
			}
			return 0
		} else if token.kind == TInvalid {
			pl.err = "invalid token"
			return 0
		} else if pl.in_ticks {
			pl.cursor = uint(token.end)
			lval.node = AstNode{token, nil}
			return int(token.kind)
		} else if token.kind == TPLSlash &&
		pl.prev_tok != TPLInc &&
		pl.prev_tok != TPLDec &&
		pl.prev_tok != TPLValue &&
		pl.prev_tok != TRBrace &&
		pl.prev_tok != TRCircle &&
		pl.prev_tok != TRSquare && 
		pl.prev_tok != TIdentifier {
			token = lexRegexpLit(pl.input, uint(token.begin))
			pl.prev_tok = TPLValue
			pl.cursor = uint(token.end)
			lval.node = AstNode{token, nil}
			return int(token.kind)
		} else if token.kind == TPLPound &&
			(pl.prev_tok == TPLNewline || pl.prev_tok == 0) &&
			pl.block_level == 0 &&
			!pl.in_macro {
			pl.prev_tok = TPLPound
			pl.cursor = uint(token.end)
			token.kind = TMacroStart
			lval.node = AstNode{token, nil}
			pl.in_macro = true
			return TMacroStart
		}
		pl.cursor = uint(token.end)
		pl.prev_tok = token.kind
		if token.kind == TPLNewline {
			if pl.in_macro {
				token.kind = TMacroEnd
				pl.in_macro = false
			} else {
				continue
			}
		} else if token.kind > TPLBefore &&
			token.kind < TPLAfter {
			token.kind = TOther
		}
		lval.node = AstNode{token, nil}
		return int(token.kind)
	}
}
func (pl *ParserLexer) Error(s string) {
	pl.err = s
}

func parseInitialAst(input []byte) (AstNode, string) {
	parser := astNewParser()
	pl := plNew(input /*, parser*/)
	result := parser.Parse(&pl)
	if pl.err == "" && result != 0 {
		pl.err = "got result " + strconv.Itoa(result)
	}
	if pl.err != "" {
		return AstNode{LexToken{0, 0, 0}, nil}, pl.err
	} else {
		return pl.result, ""
	}
}

const U32_MAX uint32 = 4294967295

func isIndisputablyJs(stmt []AstNode) bool {
	ch1 := stmt[0].token.kind
	/* C declaration can't start with special characters
	   (we allow '[' for C23 attributes)
	*/
	if ch1 == TLSquare {
		if len(stmt) < 2 {
			return true
		}
		ch1 = stmt[1].token.kind
		if ch1 != TLSquare {
			return true
		}
	} else if ch1 != TConst &&
		ch1 != TStructUnionEnumClass &&
		ch1 != TIdentifier &&

		ch1 != AstInterpol {
		return true
	}
	if len(stmt) >= 2 {
		ch2 := stmt[1].token.kind
		/* object destructuring */
		if (ch1 == TConst || ch1 == TLCircle) &&
			ch2 == TLBrace {
			return true
		}
		/* "const identifier =" */
		if len(stmt) >= 3 {
			ch3 := stmt[2].token.kind
			if ch1 == TConst &&
				ch2 == TIdentifier &&
				ch3 == TEq {
				return true
			}
		}
	}
	return false
}

func getJsBlockStmtEnd(nodes []AstNode) uint32 {
	var ch uint32
	var circle_nesting uint32 = 0
	for i := range nodes {
		ch = nodes[i].token.kind
		if ch == TLCircle {
			circle_nesting++
		} else if ch == TRCircle {
			if circle_nesting > 0 {
				circle_nesting--
			}
		} else if ch == AstBlock &&
			circle_nesting == 0 {
			return uint32(i + 1)
		}
	}
	return U32_MAX
}

func nextMaybeJsStmt(nodes []AstNode) (uint32, uint32) {
	var ch uint32
	var offset uint32 = 0
	for {
		found_meat := false
		for i := range nodes {
			ch = nodes[i].token.kind
			if ch != AstBlock && ch != AstMacro &&
				ch != TSemi {
				found_meat = true
				nodes = nodes[i:]
				offset += uint32(i)
				break
			}
		}
		/*fmt.Println("\nFOUND MEAT!!!!!!")
		debugNode(&AstNode{
			token: LexToken{0, 1, AstBlock},
			nodes: nodes,
		}, input)
		fmt.Println("\nMEAT END !!!!!!")*/
		if !found_meat || len(nodes) < 2 {
			return U32_MAX, U32_MAX
		}
		ch = nodes[0].token.kind
		if ch == TFunctionForWhileJsclass ||
		nodes[1].token.kind == TFunctionForWhileJsclass {
			jend := getJsBlockStmtEnd(nodes)
			if jend == U32_MAX {
				//fmt.Println("BAD JEND!!!!")
				return U32_MAX, U32_MAX
			}
			nodes = nodes[jend:]
			offset += jend
			continue
		}
		var stmt_end uint32 = uint32(len(nodes))
		for i := range nodes {
			ch = nodes[i].token.kind
			if ch == TSemi {
				stmt_end = uint32(i + 1)
				break
			}
		}
		if isIndisputablyJs(nodes[0:stmt_end]) {
			//fmt.Println("BUT THE MEAT IS JS!!!!")
			nodes = nodes[stmt_end:]
			offset += stmt_end
			continue
		}
		return offset, offset + stmt_end
	}
}

func startOf(node *AstNode) uint32 {
	for node.token.kind < AstEnd {
		node = &node.nodes[0]
	}
	return node.token.begin
}

func endOf(node *AstNode) uint32 {
	for node.token.kind < AstEnd {
		node = &node.nodes[len(node.nodes)-1]
	}
	return node.token.end
}

/*
	We don't allow unnecessary parentheses

in C declarations as they create ambiguity.
For example:

	"console(logger[5])[10] = value"

JS -> console fn is called with value of logger array at index 5,

	then the item at index 10 of the returned array is assigned to value

C  -> logger is an array console[5][10], it is assigned to value

	which is a macro that expands to an initializer list

In contrast, "console logger[5][10] = value" is unambiguously C,
and so is "console (*logger[5])[10] = value".
*/
func isValidJs(js_maybe []AstNode, input []byte) bool {
	/* cheap trick */
	tokens_builder := strings.Builder{}
	for i := range js_maybe {
		node := &js_maybe[i]
		_, _ = tokens_builder.Write(input[startOf(node):endOf(node)])
		_ = tokens_builder.WriteByte(' ')
	}
	tokens := tokens_builder.String()
	options := api.TransformOptions{}
	result := api.Transform(tokens, options)
	return len(result.Errors) == 0
}

func tryExtractC(nodes []AstNode, input []byte) uint32 {
	var ch uint32
	var stmt_end uint32 = uint32(len(nodes))
	var circle_nesting uint32 = 0
	has_circles := false
	found_seu := false
	for i := range nodes {
		ch = nodes[i].token.kind
		/* has a JS keyword */
		if ch == TJsOnly {
			return U32_MAX
		} else if ch == AstMacro {
			return uint32(i)
		} else if ch == TLCircle {
			has_circles = true	
			circle_nesting++
		} else if ch == TRCircle {
			if circle_nesting > 0 {
				circle_nesting--
			} else {
				return U32_MAX
			}
		} else if ch == TStructUnionEnumClass {
			found_seu = true
		} else if ch == AstBlock &&
		circle_nesting == 0 {
			// consume possible trailing semicolon post-block
			// const int values[1] = {0};
			// (nodes will always be semicolon-terminated)
			if (i+2) == len(nodes) {
				break
			} else {
				stmt_end = uint32(i+1)
				break
			}
		}
	}
	/* we don't allow struct declarations with parentheses
	   as they are ambiguous:
	   - valid struct declaration
	       struct identifier macro(identifier) { ... }
	   - valid function declaration
	       struct identifier name(arg_list) { ... }
	   because macros can be generated dynamically,
	   we cannot tell the difference.
	*/
	if found_seu && !has_circles {
		return uint32(len(nodes))
	}

	stmt := nodes[:stmt_end]
	for i := range stmt {
		ch = stmt[i].token.kind
		/* has interpolated JS,
		   or struct/enum/union keyword */
		if ch == AstInterpol ||
			ch == TStructUnionEnumClass {
			return stmt_end
		}
	}
	for i := range stmt {
		ch = stmt[i].token.kind
		if ch == TEq {
			break
			/* dot operator in type specifier
			or declarator impossible */
		} else if ch == '.' {
			return U32_MAX
		}
	}
	if isValidJs(nodes, input) {
		return U32_MAX
	}
	return stmt_end
}

func nextCSection(root []AstNode, input []byte) (uint32, uint32) {
	var offset uint32 = 0
	for {
		jstart, jend := nextMaybeJsStmt(root)
		//fmt.Printf("jstart is %d\n", jstart)
		if jstart == U32_MAX {
			return U32_MAX, U32_MAX
		}
		cend := tryExtractC(root[jstart:jend], input)
		if cend == U32_MAX {
			root = root[jend:]
			offset += jend
			continue
		}

		return offset + jstart, offset + jstart + cend
	}
}

/*
func getBeginOf(node *AstNode) uint32 {
	for node.token.kind < AstEnd {
		node = &node.nodes[0]
	}
	return node.token.begin
}

func getEndOf(node *AstNode) uint32 {
	for node.token.kind < AstEnd {
		node = &node.nodes[len(node.nodes)-1]
	}
	return node.token.end
}*/

func reviseTLC(root []AstNode, input []byte) []AstNode {
	new_root := make([]AstNode, 0, 4)
	for {
		cstart, cend := nextCSection(root, input)
		if cstart == U32_MAX {
			break
		} else {
			//fmt.Printf("c section is [%d, %d)\n", cstart, cend)
		}
		new_root = append(new_root, root[0:cstart]...)
		new_root = append(new_root, AstNode{
			token: LexToken{
				begin: 0,
				end:   0,
				kind:  AstMacro,
			},
			nodes: root[cstart:cend],
		})
		root = root[cend:]
	}
	new_root = append(new_root, root...)
	return new_root
}

func findLastLfn(nodes []AstNode) (uint32, uint32, uint32) {
	var prev uint32 = 0
	var cur uint32 = 0
	for i := len(nodes) - 1; i >= 0; i-- {
		cur = prev
		prev = nodes[i].token.kind
		if prev == TRCircle &&
			cur == TMinusGt {
			var circle_nesting uint32 = 1
			mgt := i + 1
			for i--; i >= 0; i-- {
				var kind uint32 = nodes[i].token.kind
				if kind == TRCircle {
					circle_nesting++
				} else if kind == TLCircle {
					circle_nesting--
					if circle_nesting == 0 {
						break
					}
				}
			}
			if circle_nesting != 0 {
				break
			}
			start := i
			circle_nesting = 0
			found_block := false
			for i = mgt + 1; i < len(nodes); i++ {
				var kind uint32 = nodes[i].token.kind
				if kind == TLCircle {
					circle_nesting++
				} else if kind == TRCircle {
					if circle_nesting == 0 {
						break
					}
					circle_nesting--
				} else if kind == AstBlock &&
					(mgt+1) != i {
					found_block = true
					break
				}
			}
			if !found_block {
				i = start
				continue
			}
			return uint32(start), uint32(mgt), uint32(i + 1)
		}
	}
	return U32_MAX, U32_MAX, U32_MAX
}

func replaceLfnMaybe(root []AstNode) []AstNode {
	type Lfn struct {
		start, mgt, end uint32
	}
	lfn_locs := make([]Lfn, 0, 4)
	nodes := root
	for {
		lstart, lmgt, lend := findLastLfn(nodes)
		if lstart == U32_MAX {
			break
		}
		lfn_locs = append(lfn_locs, Lfn{lstart, lmgt, lend})
		nodes = nodes[0:lstart]
	}
	if len(lfn_locs) == 0 {
		return nil
	}
	replaced := make([]AstNode, 0, len(root)-(3*len(lfn_locs)))
	var prev_end uint32 = 0
	for i := len(lfn_locs) - 1; i >= 0; i-- {
		l := lfn_locs[i]
		replaced = append(replaced, root[prev_end:l.start]...)
		replaced = append(replaced, AstNode{
			token: LexToken{0, 0, AstLfnArgs},
			nodes: root[l.start:l.mgt],
		})
		replaced = append(replaced, AstNode{
			token: LexToken{0, 0, AstLfnRetType},
			nodes: root[l.mgt+1 : l.end-1],
		})
		body := root[l.end-1]
		body.token.kind = AstLfnBody
		replaced = append(replaced, body)
		prev_end = l.end
	}
	replaced = append(replaced, root[prev_end:]...)
	return replaced
}

func reviseWithinLfn(root []AstNode) {
	for i := range root {
		kind := root[i].token.kind
		if kind == AstBlock {
			reviseWithinLfn(root[i].nodes)
		} else if kind == AstInterpol {
			replaced := reviseLfn(root[i].nodes)
			if replaced != nil {
				root[i].nodes = replaced
			}
		} else {
			return
		}
	}
}

func reviseLfn(root []AstNode) []AstNode {
	replaced := replaceLfnMaybe(root)
	changed := false
	if replaced != nil {
		changed = true
		root = replaced
	}
	for i := range root {
		kind := root[i].token.kind
		if kind > AstEnd {
			continue
		}
		if kind == AstLfnArgs ||
			kind == AstLfnRetType {
			reviseWithinLfn(root[i].nodes)
		} else {
			replaced = reviseLfn(root[i].nodes)
			if replaced != nil {
				root[i].nodes = replaced
			}
		}
	}
	if changed {
		return root
	} else {
		return nil
	}
}

func parseIntoAst(input []byte) (AstNode, string) {
	root, estr := parseInitialAst(input)
	if estr != "" {
		return AstNode{LexToken{0, 0, 0}, nil}, estr
	}
	root.nodes = reviseTLC(root.nodes, input)
	revised := reviseLfn(root.nodes)
	if revised != nil {
		root.nodes = revised
	}
	return root, ""
}
