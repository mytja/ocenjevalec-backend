package main

import (
	"HTTP-boilerplate/ast"
	"encoding/json"
	"fmt"
)

func main() {
	solutionText := "AND(NOT(A), B)"
	submissionText := "AND(a, b)"

	subS := ast.MinifyString(submissionText)
	sub, err := ast.BuildAST(subS)
	if err != nil {
		fmt.Println("AST build failed", err.Error())
		return
	}
	j, err := json.MarshalIndent(sub, "", "\t")
	if err != nil {
		fmt.Println("JSON marshal failed", err.Error())
		return
	}
	fmt.Println("Contestant solution AST", string(j))

	solS := ast.MinifyString(solutionText)
	sol, err := ast.BuildAST(solS)
	if err != nil {
		fmt.Println("AST solution build failed", err.Error())
		return
	}

	solution, err := ast.TestDigitalSolution(sub, sol)
	if err != nil {
		fmt.Println("Digital solution testing failed", err.Error())
		return
	}

	if solution.Verdict == "AC" && solS != subS {
		fmt.Printf("Applying PARTIAL verdict! Contestant: %s, Judge: %s.\n", subS, solS)
		solution.Verdict = "PARTIAL"
	}

	fmt.Println("Judging complete!")
	fmt.Println("Verdict:", solution.Verdict)
	fmt.Printf("Correct test cases/wrong test cases: %d/%d\n", solution.CorrectTestCases, solution.WrongTestCases)
	fmt.Println("Evaluation log:")
	fmt.Printf(solution.EvaluationLog)
}
