package gates

import (
	"fmt"
	"strconv"

	"github.com/lujjjh/gates/syntax"
)

type compiler struct {
	program *Program
}

func (c *compiler) compileLit(l *syntax.Lit) {
	switch l.Kind {
	case syntax.NUMBER:
		var v Value
		i, err := strconv.ParseInt(l.Value, 0, 64)
		if err == nil {
			v = intNumber(i)
		} else {
			f, _ := strconv.ParseFloat(l.Value, 64)
			v = floatNumber(f)
		}
		c.program.emit(load(c.program.defineLit(v)))
	case syntax.STRING:
		s, _ := strconv.Unquote(l.Value)
		c.program.emit(load(c.program.defineLit(String(s))))
	case syntax.BOOL:
		c.program.emit(load(c.program.defineLit(Bool(l.Value == "true"))))
	default:
		panic(fmt.Errorf("unknown token type %v", l.Kind))
	}
}

func (c *compiler) compileUnaryExpr(e *syntax.UnaryExpr) {
	c.compileExpr(e.X)
	switch e.Op {
	case syntax.ADD:
		c.program.emit(plus)
	case syntax.SUB:
		c.program.emit(neg)
	case syntax.NOT:
		c.program.emit(not)
	default:
		panic(fmt.Errorf("unknown unary operator: %s", e.Op))
	}
}

func (c *compiler) compileBinaryExpr(e *syntax.BinaryExpr) {
	switch e.Op {
	case syntax.ADD:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(add)
	case syntax.SUB:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(sub)
	case syntax.MUL:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(mul)
	case syntax.QUO:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(div)
	case syntax.REM:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(mod)
	case syntax.AND:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(and)
	case syntax.OR:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(or)
	case syntax.XOR:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(xor)
	case syntax.SHL:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(shl)
	case syntax.SHR:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(shr)
	case syntax.LAND:
		c.compileExpr(e.X)
		j := len(c.program.code)
		c.program.emit(nil, pop)
		c.compileExpr(e.Y)
		c.program.code[j] = jneq1(len(c.program.code) - j)
	case syntax.LOR:
		c.compileExpr(e.X)
		j := len(c.program.code)
		c.program.emit(nil, pop)
		c.compileExpr(e.Y)
		c.program.code[j] = jeq1(len(c.program.code) - j)
	case syntax.EQL:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(eq)
	case syntax.LSS:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(lt)
	case syntax.GTR:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(gt)
	case syntax.NEQ:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(neq)
	case syntax.LEQ:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(lte)
	case syntax.GEQ:
		c.compileExpr(e.X)
		c.compileExpr(e.Y)
		c.program.emit(gte)
	default:
		panic(fmt.Errorf("unknown binary operator: %s", e.Op))
	}
}

func (c *compiler) compileExpr(e syntax.Expr) {
	switch e := e.(type) {
	case *syntax.Lit:
		c.compileLit(e)
	case *syntax.UnaryExpr:
		c.compileUnaryExpr(e)
	case *syntax.BinaryExpr:
		c.compileBinaryExpr(e)
	case *syntax.ParenExpr:
		c.compileExpr(e.X)
	default:
		panic(fmt.Errorf("unknown expression type: %T", e))
	}
}

func (c *compiler) compile(e syntax.Expr) {
	c.compileExpr(e)
	c.program.emit(halt)
}