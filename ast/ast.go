package ast

import (
	"errors"
	"fmt"
	"hash/fnv"
	"math"
	"strings"
)

const (
	NOT   = iota
	AND   = iota
	NAND  = iota
	OR    = iota
	NOR   = iota
	XOR   = iota
	XNOR  = iota
	INPUT = iota
)

type AST struct {
	Type       int  `json:"type"`
	Input      rune `json:"input"`
	SubEntity1 *AST `json:"sub_entity_1"`
	SubEntity2 *AST `json:"sub_entity_2"`
}

func recursivelyBuildAST(s *string) (*AST, error) {
	runes := []rune(*s)
	if len(runes) == 0 {
		return nil, errors.New("AST is empty")
	}

	if runes[len(runes)-1] != ')' {
		if len(runes) != 1 {
			return nil, errors.New("AST is invalid. Invalid ending")
		}
		r := runes[0]
		if r < 'A' || r > 'Z' {
			return nil, errors.New("AST is invalid. Invalid ASCII character")
		}
		return &AST{
			Type:       INPUT,
			Input:      r,
			SubEntity1: nil,
			SubEntity2: nil,
		}, nil
	}

	levels := 0
	t := ""
	value := ""
	values := make([]string, 0)

	for _, v := range runes {
		if v == '(' {
			levels++
			if levels == 1 {
				continue
			}
		}
		if v == ')' {
			levels--
			if levels == 0 {
				continue
			}
			if levels < 0 {
				return nil, errors.New(fmt.Sprintf("AST is invalid. Invalid bracket level count %d/0", levels))
			}
		}
		if v == ',' && levels == 1 {
			values = append(values, value)
			value = ""
			continue
		}
		if levels == 0 {
			t += string(v)
			continue
		}
		value += string(v)
	}

	if levels != 0 {
		return nil, errors.New(fmt.Sprintf("AST is invalid. Invalid bracket level count %d/0", levels))
	}

	if value != "" {
		values = append(values, value)
	}

	if t == "NOT" && len(values) != 1 {
		return nil, errors.New(fmt.Sprintf("AST is invalid. Invalid value count for type NOT: %d/1", len(values)))
	}
	if t != "NOT" && len(values) != 2 {
		return nil, errors.New(fmt.Sprintf("AST is invalid. Invalid value count for type %s: %d/2", t, len(values)))
	}

	recursive1, err := recursivelyBuildAST(&values[0])
	if err != nil {
		return nil, err
	}

	if t == "NOT" {
		return &AST{
			Type:       NOT,
			SubEntity1: recursive1,
			SubEntity2: nil,
		}, nil
	}

	recursive2, err := recursivelyBuildAST(&values[1])
	if err != nil {
		return nil, err
	}

	if t == "AND" {
		return &AST{
			Type:       AND,
			SubEntity1: recursive1,
			SubEntity2: recursive2,
		}, nil
	} else if t == "NAND" {
		return &AST{
			Type:       NAND,
			SubEntity1: recursive1,
			SubEntity2: recursive2,
		}, nil
	} else if t == "OR" {
		return &AST{
			Type:       OR,
			SubEntity1: recursive1,
			SubEntity2: recursive2,
		}, nil
	} else if t == "NOR" {
		return &AST{
			Type:       NOR,
			SubEntity1: recursive1,
			SubEntity2: recursive2,
		}, nil
	} else if t == "XOR" {
		return &AST{
			Type:       XOR,
			SubEntity1: recursive1,
			SubEntity2: recursive2,
		}, nil
	} else if t == "XNOR" {
		return &AST{
			Type:       XNOR,
			SubEntity1: recursive1,
			SubEntity2: recursive2,
		}, nil
	}

	return nil, errors.New(fmt.Sprintf("AST is invalid. Invalid type %s", t))
}

func ASTLength(ast *AST, current int) int {
	if ast.Type == INPUT {
		return current
	}
	if ast.SubEntity1 == nil {
		return current
	}
	if ast.Type == NOT {
		return ASTLength(ast.SubEntity1, current+1)
	}
	if ast.SubEntity2 == nil {
		return ASTLength(ast.SubEntity1, current+1)
	}
	return current + ASTLength(ast.SubEntity1, 0) + ASTLength(ast.SubEntity2, 0) + 1
}

func VerifyAgainstBase(ast *AST) bool {
	if ast.Type == INPUT {
		return true
	}
	if ast.SubEntity1 == nil {
		return true
	}
	if ast.Type == NOT {
		return VerifyAgainstBase(ast.SubEntity1)
	}
	if ast.SubEntity2 == nil {
		return VerifyAgainstBase(ast.SubEntity1)
	}
	if ast.Type == AND || ast.Type == OR {
		return VerifyAgainstBase(ast.SubEntity1) && VerifyAgainstBase(ast.SubEntity2)
	}
	return false
}

func MinifyString(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ToUpper(s)
	return s
}

func BuildAST(s string) (*AST, error) {
	return recursivelyBuildAST(&s)
}

type DigitalSolutionEvaluation struct {
	CorrectTestCases int
	WrongTestCases   int
	EvaluationLog    string
	Verdict          string
}

func recursiveBuildInputs(a *AST, m *map[rune]bool) *map[rune]bool {
	if a.Type == INPUT {
		l := *m
		l[a.Input] = true
		return &l
	}
	m = recursiveBuildInputs(a.SubEntity1, m)
	if a.Type == NOT {
		return m
	}
	m = recursiveBuildInputs(a.SubEntity2, m)
	return m
}

func BuildInputs(a *AST) *[]rune {
	m := make(map[rune]bool)
	recursiveBuildInputs(a, &m)
	r := make([]rune, 0)
	for i := range m {
		r = append(r, i)
	}
	return &r
}

func evaluate(a *AST, m *map[rune]bool) (bool, error) {
	if a.Type == INPUT {
		c, exists := (*m)[a.Input]
		if !exists {
			return false, errors.New(fmt.Sprintf("Rune %s (%d) doesn't exist amongst values", string(a.Input), a.Input))
		}
		return c, nil
	}

	sube1, err := evaluate(a.SubEntity1, m)
	if err != nil {
		return false, err
	}

	if a.Type == NOT {
		return NOTOperation(sube1), nil
	}

	sube2, err := evaluate(a.SubEntity2, m)
	if err != nil {
		return false, err
	}

	if a.Type == AND {
		return ANDOperation(sube1, sube2), nil
	} else if a.Type == NAND {
		return NANDOperation(sube1, sube2), nil
	} else if a.Type == OR {
		return OROperation(sube1, sube2), nil
	} else if a.Type == NOR {
		return NOROperation(sube1, sube2), nil
	} else if a.Type == XOR {
		return XOROperation(sube1, sube2), nil
	} else if a.Type == XNOR {
		return XNOROperation(sube1, sube2), nil
	}

	return false, errors.New(fmt.Sprintf("Invalid type %d", a.Type))
}

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

func HashAST(a *AST) int {
	if a.Type == INPUT {
		return int(a.Input)
	}
	if a.Type == NOT {
		return hash(fmt.Sprintf("NOT(%d)", HashAST(a.SubEntity1)))
	}

	t := ""
	if a.Type == AND {
		t = "AND"
	} else if a.Type == NAND {
		t = "NAND"
	} else if a.Type == OR {
		t = "OR"
	} else if a.Type == NOR {
		t = "NOR"
	} else if a.Type == XOR {
		t = "XOR"
	} else if a.Type == XNOR {
		t = "XNOR"
	}

	h1 := HashAST(a.SubEntity1)
	h2 := HashAST(a.SubEntity2)
	if h1 < h2 {
		return hash(fmt.Sprintf("%s(%d,%d)", t, h1, h2))
	}
	return hash(fmt.Sprintf("%s(%d,%d)", t, h2, h1))
}

func boolToInt(a bool) int {
	if a {
		return 1
	}
	return 0
}

func TestDigitalSolution(aSubmission *AST, aSolution *AST) (*DigitalSolutionEvaluation, error) {
	inputs := *BuildInputs(aSolution)
	m := make(map[rune]bool)
	dse := DigitalSolutionEvaluation{
		CorrectTestCases: 0,
		WrongTestCases:   0,
		EvaluationLog:    "",
		Verdict:          "",
	}
	for n := 0; n < int(math.Pow(2, float64(len(inputs)))); n++ {
		f := fmt.Sprint(len(inputs))
		if len(inputs) < 10 {
			f = "0" + f
		}
		format := fmt.Sprintf("%%%sb", f)
		ns := fmt.Sprintf(format, int64(n))
		for i, v := range []rune(ns) {
			m[inputs[i]] = v != '0'
		}
		sol, err := evaluate(aSolution, &m)
		if err != nil {
			dse.EvaluationLog += fmt.Sprintf("Solution evaluation failed on test case %d. Error: %s.\n", n+1, err.Error())
			dse.Verdict = "SOL_RTE" // Solution Runtime error
			return &dse, nil
		}
		sub, err := evaluate(aSubmission, &m)
		if err != nil {
			dse.EvaluationLog += fmt.Sprintf("Submission evaluation failed on test case %d. Error: %s.\n", n+1, err.Error())
			dse.Verdict = "RTE"
			continue
		}
		if sub != sol {
			dse.WrongTestCases++
			dse.EvaluationLog += fmt.Sprintf("Wrong answer on test case %d (%s)! Contestant: %d, Judge: %d.\n", n+1, ns, boolToInt(sub), boolToInt(sol))
			dse.Verdict = "WA"
			continue
		}
		dse.EvaluationLog += fmt.Sprintf("Correct answer on test case %d (%s)! Contestant: %d, Judge: %d.\n", n+1, ns, boolToInt(sub), boolToInt(sol))
		dse.CorrectTestCases++
	}
	if dse.Verdict == "" {
		dse.Verdict = "AC"
	}
	return &dse, nil
}
