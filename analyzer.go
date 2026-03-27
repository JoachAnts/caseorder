package switchorder

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/token"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "switchorder",
	Doc:  "checks that switch case statements are in alphabetical or numerical order",
	Run:  run,
}

// New is required for golangci-lint compatibility.
func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{Analyzer}, nil
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switchStmt, ok := n.(*ast.SwitchStmt)
			if !ok {
				return true
			}

			processSwitch(pass, switchStmt)
			return true
		})
	}
	return nil, nil
}

type valueInfo struct {
	expr ast.Expr
	lit  *ast.BasicLit
	val  constant.Value
}

type caseClauseInfo struct {
	clause *ast.CaseClause
	values []valueInfo
}

func processSwitch(pass *analysis.Pass, sw *ast.SwitchStmt) {
	var cases []caseClauseInfo
	var defaultCase *ast.CaseClause

	// --- Collect cases ---
	for _, stmt := range sw.Body.List {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok {
			return
		}

		// default case
		if cc.List == nil {
			defaultCase = cc
			continue
		}

		var values []valueInfo
		for _, expr := range cc.List {
			bl, ok := expr.(*ast.BasicLit)
			if !ok || (bl.Kind != token.STRING && bl.Kind != token.INT && bl.Kind != token.FLOAT && bl.Kind != token.CHAR) {
				return
			}
			values = append(values, valueInfo{
				expr: expr,
				lit:  bl,
				val:  constant.MakeFromLiteral(bl.Value, bl.Kind, 0),
			})
		}

		// Skip if fallthrough exists (unsafe to reorder)
		for _, stmt := range cc.Body {
			if br, ok := stmt.(*ast.BranchStmt); ok && br.Tok == token.FALLTHROUGH {
				return
			}
		}

		cases = append(cases, caseClauseInfo{
			clause: cc,
			values: values,
		})
	}

	if len(cases) == 0 {
		return
	}

	changed := false
	var diagnostics []analysis.Diagnostic

	// --- Check and Fix within each case ---
	for i := range cases {
		if len(cases[i].values) < 2 {
			continue
		}
		if !isSorted(cases[i].values) {
			changed = true
			// Report first out-of-order value
			for j := 1; j < len(cases[i].values); j++ {
				if constant.Compare(cases[i].values[j-1].val, token.GTR, cases[i].values[j].val) {
					diagnostics = append(diagnostics, analysis.Diagnostic{
						Pos:     cases[i].values[j].expr.Pos(),
						End:     cases[i].values[j].expr.End(),
						Message: fmt.Sprintf("case value %s should come before %s", cases[i].values[j].lit.Value, cases[i].values[j-1].lit.Value),
					})
				}
			}
			sortValues(cases[i].values)
			newList := make([]ast.Expr, len(cases[i].values))
			for j, v := range cases[i].values {
				newList[j] = v.expr
			}
			cases[i].clause.List = newList
		}
	}

	// --- Check and Fix across cases ---
	if !isCasesSorted(cases) {
		changed = true
		for i := 1; i < len(cases); i++ {
			if constant.Compare(cases[i-1].values[0].val, token.GTR, cases[i].values[0].val) {
				diagnostics = append(diagnostics, analysis.Diagnostic{
					Pos:     cases[i].clause.Pos(),
					End:     cases[i].clause.End(),
					Message: fmt.Sprintf("case %s should come before %s", cases[i].values[0].lit.Value, cases[i-1].values[0].lit.Value),
				})
			}
		}
		sortCases(cases)
	}

	if !changed {
		return
	}

	// --- Build new body ---
	var parts []string
	for _, c := range cases {
		var buf bytes.Buffer
		if err := format.Node(&buf, pass.Fset, c.clause); err == nil {
			parts = append(parts, buf.String())
		}
	}
	if defaultCase != nil {
		var buf bytes.Buffer
		if err := format.Node(&buf, pass.Fset, defaultCase); err == nil {
			parts = append(parts, buf.String())
		}
	}
	content := strings.Join(parts, "\n")

	fix := analysis.SuggestedFix{
		Message: "reorder switch cases",
		TextEdits: []analysis.TextEdit{
			{
				Pos:     sw.Body.Lbrace + 1,
				End:     sw.Body.Rbrace,
				NewText: []byte("\n" + content + "\n"),
			},
		},
	}

	for _, d := range diagnostics {
		d.SuggestedFixes = []analysis.SuggestedFix{fix}
		pass.Report(d)
	}
}

func isSorted(values []valueInfo) bool {
	for i := 1; i < len(values); i++ {
		if constant.Compare(values[i-1].val, token.GTR, values[i].val) {
			return false
		}
	}
	return true
}

func sortValues(values []valueInfo) {
	sort.Slice(values, func(i, j int) bool {
		return constant.Compare(values[i].val, token.LSS, values[j].val)
	})
}

func isCasesSorted(cases []caseClauseInfo) bool {
	for i := 1; i < len(cases); i++ {
		if constant.Compare(cases[i-1].values[0].val, token.GTR, cases[i].values[0].val) {
			return false
		}
	}
	return true
}

func sortCases(cases []caseClauseInfo) {
	sort.Slice(cases, func(i, j int) bool {
		return constant.Compare(cases[i].values[0].val, token.LSS, cases[j].values[0].val)
	})
}
