package main

const js_prelude string = `
/* beginning of JC prelude */

if (typeof JC_INTERNAL_OUTFILE !== "string") {
    throw new Error();
}

class JcInternalQuote {
    constructor(sections) {
        this.sections = sections;
    }
    add(other) {
        if (!(other instanceof JcInternalQuote)) {
            console.error("got", other);
            throw new Error(
                "bad argument 1 to quote::add: "+
                "one can only add a quote to a quote"
            );
        }
        var mysec = this.sections;
        var othersec = other.sections;
        for (var item of othersec) {
            mysec.push(item);
        }
    }
    plus(other) {
        if (!(other instanceof JcInternalQuote)) {
            console.error("got", other);
            throw new Error(
                "bad argument 1 to quote::plus: "+
                "one can only concatenate two quotes"
            );
        }
        var mysec = this.sections;
        var othersec = other.sections;
        var newsec = [];
        for (var item of mysec) {
            newsec.push(item);
        }
        for (var item of othersec) {
            newsec.push(item);
        }
        return new JcInternalQuote(newsec);
    }
}

const JC_INTERNAL_TOPLEVEL = [];
var JC_INTERNAL_DENO = false;
try {
    var item = Deno.writeTextFileSync;
    if (item) {
        JC_INTERNAL_DENO = true;
    }
} catch (_e) {}

function strtoq(inp) {
    if (typeof inp !== "string") {
        console.error("expected string, got", inp);
        throw new Error("bad argument 1 to strtoq");
    }
    return new JcInternalQuote([inp]);
}

var JC_INTERNAL_NEXTSYM = 1;
function gensym() {
    return new JcInternalQuote(["jc_sym"+(JC_INTERNAL_NEXTSYM++)]);
}

class JcInternalFunction {
    constructor(args, ret_type, body) {
        this.args = args;
        this.ret_type = ret_type;
        this.body = body;
        this.num = 0;
    }
}

const JC_INTERNAL_ENCODER = new TextEncoder();
function jcInternalToCLiteral(str) {
    const src = JC_INTERNAL_ENCODER.encode(str);
    var res = "\"";
    for (const ch of src) {
        switch (ch) {
        case 0x22:
            res += "\\\"";
            break;
        case 0x5c:
            res += "\\";
            break;
        case 0x07:
            res += "\\a";
            break;
        case 0x08:
            res += "\\b";
            break;
        case 0x0c:
            res += "\\f";
            break;
        case 0x0a:
            res += "\\n";
            break;
        case 0x0d:
            res += "\\r";
            break;
        case 0x09:
            res += "\\t";
            break;
        case 0x0b:
            res += "\\v";
            break;
        default:
            if (ch >= 32 && ch < 127) {
                res += String.fromCharCode(ch);
            } else {
                res += "\\";
                res += ch.toString(8).padStart(3, "0");
            }
        }
    }
    res += "\"";
    return res;
}

var JC_INTERNAL_OUTPUT = "";
var JC_INTERNAL_NEXT_FN = 1;
async function jcInternalWrite(obj) {
    if (typeof obj === "string") {
        return jcInternalToCLiteral(obj);
    } else if (typeof obj === "object") {
        while (obj instanceof Promise) {
            obj = await obj;
        }
        if (Array.isArray(obj)) {
            if (obj.length === 0) {
                return "{}";
            }
            var str = "{";
            str += await jcInternalWrite(obj[0]);
            for (var i = 1; i < obj.length; i++) {
                str += ", ";
                str += await jcInternalWrite(obj[i]);    
            }
            str += "}";
            return str;
        } else if (obj instanceof JcInternalQuote) {
            var str = "";
            for (const item of obj.sections) {
                if (typeof item === "string") {
                    str += item;
                } else {
                    str += await jcInternalWrite(item.obj);
                }
            }
            return str;
        } else if (obj instanceof JcInternalFunction) {
            if (obj.num === 0) {
                var str = "\n\n";
                str += await jcInternalWrite(obj.ret_type);
                str += " jc_lfn";
                str += (obj.num = (JC_INTERNAL_NEXT_FN++));
                str += await jcInternalWrite(obj.args);
                str += " ";
                str += await jcInternalWrite(obj.body);
                str += "\n\n";
                JC_INTERNAL_OUTPUT += str;
            }
            return "jc_lfn"+obj.num;
        } else {
            var str = "{";
            var first = true;
            for (var key in obj) {
                var value = obj[key];
                if (typeof value === "function") {
                    continue;
                }
                if (first) {
                    first = false;
                } else {
                    str += ", ";
                }
                str += ".";
                str += key;
                str += "=";
                str += await jcInternalWrite(value);
            }
            str += "}";
            return str;
        }
    } else if (typeof obj === "bigint") {
        return obj.toString(10);
    } else if (typeof obj === "number") {
        return obj.toString(10);
    } else {
        console.error("got", obj);
        throw new Error("don't know how to convert given object into C");
    }
}

async function jcInternalFinish() {
    for (const item of JC_INTERNAL_TOPLEVEL) {
        var str = await jcInternalWrite(item);
        JC_INTERNAL_OUTPUT += str;
    }
    JC_INTERNAL_OUTPUT += "\n";
    if (JC_INTERNAL_DENO) {
        Deno.writeTextFileSync(
            JC_INTERNAL_OUTFILE,
            JC_INTERNAL_OUTPUT
        );
    } else {
        require("fs").writeFileSync(
            JC_INTERNAL_OUTFILE,
            JC_INTERNAL_OUTPUT
        );
    }
}

/* end of JC prelude */
`

type Codegen struct {
	sb    []byte
	input string
    inqtable []byte
}

func pushStringLiteral(src string, sb []byte) []byte {
	sb = append(sb, '"')
	for i := range src {
		ch := src[i]
		switch ch {
		case 0x22:
			sb = append(sb, '\\', '"')
			break
		case 0x5c:
			sb = append(sb, '\\', '\\')
			break
		case 0x07:
			sb = append(sb, '\\', 'a')
			break
		case 0x08:
			sb = append(sb, '\\', 'b')
			break
		case 0x0c:
			sb = append(sb, '\\', 'f')
			break
		case 0x0a:
			sb = append(sb, '\\', 'n')
			break
		case 0x0d:
			sb = append(sb, '\\', 'r')
			break
		case 0x09:
			sb = append(sb, '\\', 't')
			break
		case 0x0b:
			sb = append(sb, '\\', 'v')
			break
		default:
			if ch >= 32 && ch < 127 {
				sb = append(sb, ch)
			} else {
				ch3 := ch & 8
				ch >>= 3
				ch2 := ch & 8
				ch >>= 3
				ch1 := ch & 8
				sb = append(sb, '\\', ch1, ch2, ch3)
			}
		}
	}
	sb = append(sb, '"')
	return sb
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

func (cg *Codegen) generateForC(nodes []AstNode) {
	iend_prev := startOf(&nodes[0])
	for i := range nodes {
		kind := nodes[i].token.kind
		if kind > AstEnd {
			continue
		}
		istart := startOf(&nodes[i])
		if iend_prev != istart {
			cg.sb = pushStringLiteral(
				cg.input[iend_prev:istart],
				cg.sb,
			)
			cg.sb = append(cg.sb, ',', '\n')
		}
		iend_prev = endOf(&nodes[i])
		if kind == AstBlock {
			cg.generateForC(nodes[i].nodes)
		} else {
			cg.sb = append(cg.sb, "{obj: ("...)
			cg.generateForJs(nodes[i].nodes[1 : len(nodes[i].nodes)-1])
			cg.sb = append(cg.sb, ')', '}', ',', '\n')
		}
	}
	iend := endOf(&nodes[len(nodes)-1])
	if iend_prev != iend {
		cg.sb = pushStringLiteral(
			cg.input[iend_prev:iend],
			cg.sb,
		)
		cg.sb = append(cg.sb, ',', '\n')
	}
}

func (cg *Codegen) generateForJs(nodes []AstNode) {
	iend_prev := startOf(&nodes[0])
	for i := range nodes {
		kind := nodes[i].token.kind
		if kind > AstEnd {
			continue
		}
		istart := startOf(&nodes[i])
		if iend_prev != istart {
			cg.sb = append(cg.sb,
				cg.input[iend_prev:istart]...,
			)
		}
		iend_prev = endOf(&nodes[i])
		if kind == AstBlock {
			cg.generateForJs(nodes[i].nodes)
		} else if kind == AstMacro {
			cg.sb = append(cg.sb, "\nJC_INTERNAL_TOPLEVEL.push(new JcInternalQuote([\n"...)
			cg.generateForC(nodes[i].nodes)
			cg.sb = append(cg.sb, "\n\"\\n\"\n]));\n"...)
		} else if kind == AstLfnArgs {
			cg.sb = append(cg.sb, "new JcInternalFunction(\nnew JcInternalQuote([\n"...)
			cg.generateForC(nodes[i].nodes)
			cg.sb = append(cg.sb, "]),\nnew JcInternalQuote([\n"...)
			iend_prev = startOf(&nodes[i+1])
		} else if kind == AstLfnRetType {
			cg.generateForC(nodes[i].nodes)
			cg.sb = append(cg.sb, "]),\nnew JcInternalQuote([\n"...)
		} else if kind == AstLfnBody {
			cg.generateForC(nodes[i].nodes)
			cg.sb = append(cg.sb, "]))"...)
		} else if kind == AstTicklit {
			tnodes := nodes[i].nodes
			for j := range tnodes {
				kind = tnodes[j].token.kind
				if kind == AstInterpol {
					cg.generateForJs(tnodes[j].nodes)
				} else {
					cg.sb = append(cg.sb, cg.input[tnodes[j].token.begin:tnodes[j].token.end]...)
				}
			}
		} else if kind == AstQuote {
			qnodes := nodes[i].nodes
            if len(qnodes) == 3 {
                cg.sb = append(cg.sb, "new JcInternalQuote([])"...)
            } else {
                cg.sb = append(cg.sb, "new JcInternalQuote([\n"...)
                lparen_end := qnodes[1].token.end
                quote_start := startOf(&qnodes[2])
                if lparen_end < quote_start {
                    cg.sb = pushStringLiteral(cg.input[lparen_end:quote_start], cg.sb)
                    cg.sb = append(cg.sb, ",\n"...)
                }
                cg.generateForC(qnodes[2:len(qnodes)-1])
                quote_end := endOf(&qnodes[len(qnodes)-2])
                rparen_start := qnodes[len(qnodes)-1].token.begin
                if quote_end < rparen_start {
                    cg.sb = pushStringLiteral(cg.input[quote_end:rparen_start], cg.sb)
                    cg.sb = append(cg.sb, "\n])"...)
                } else {
                    cg.sb = append(cg.sb, "])"...)
                }
            }
		}
	}
	iend := endOf(&nodes[len(nodes)-1])
	if iend_prev != iend {
		cg.sb = append(cg.sb, cg.input[iend_prev:iend]...)
	}
}

func codegenPerform(root *AstNode, input string, outfile string) []byte {
	cg := Codegen{
		sb:    make([]byte, 0, 16384),
		input: input,
	}
	cg.sb = append(cg.sb, "const JC_INTERNAL_OUTFILE = "...)
	cg.sb = pushStringLiteral(outfile, cg.sb)
	cg.sb = append(cg.sb, ";\n"...)
	cg.sb = append(cg.sb, js_prelude...)
	cg.generateForJs(root.nodes)
    cg.sb = append(cg.sb, "\njcInternalFinish();\n"...)
	return cg.sb
}
