package switchorder

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "switchorder",
	Doc:  "checks that switch case statements are in lexicographical order",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switchStmt, ok := n.(*ast.SwitchStmt)
			if !ok {
				return true
			}

			var prev string

			for _, stmt := range switchStmt.Body.List {
				caseClause, ok := stmt.(*ast.CaseClause)
				if !ok || len(caseClause.List) == 0 {
					continue
				}

				bl, ok := caseClause.List[0].(*ast.BasicLit)

				current := bl.Value

				if prev != "" && current < prev {
					pass.Reportf(caseClause.Pos(),
						"case %s should come before %s", current, prev)
				}

				prev = current
			}

			return true
		})
	}

	return nil, nil
}
