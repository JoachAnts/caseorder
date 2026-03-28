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

func getValue(expr ast.Expr) (*ast.BasicLit, constant.Value, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING || e.Kind == token.INT || e.Kind == token.FLOAT || e.Kind == token.CHAR {
			return e, constant.MakeFromLiteral(e.Value, e.Kind, 0), true
		}
	case *ast.UnaryExpr:
		if bl, ok := e.X.(*ast.BasicLit); ok {
			if bl.Kind == token.INT || bl.Kind == token.FLOAT {
				val := constant.MakeFromLiteral(bl.Value, bl.Kind, 0)
				if e.Op == token.SUB {
					val = constant.UnaryOp(token.SUB, val, 0)
				}
				return bl, val, true
			}
		}
	}
	return nil, nil, false
}

func getLitValueString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return e.Value
	case *ast.UnaryExpr:
		if bl, ok := e.X.(*ast.BasicLit); ok {
			return fmt.Sprintf("%s%s", e.Op, bl.Value)
		}
	}
	return ""
}

func processSwitch(pass *analysis.Pass, sw *ast.SwitchStmt) {
	var cases []caseClauseInfo

	// --- Collect cases ---
	for _, stmt := range sw.Body.List {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok {
			return
		}

		var values []valueInfo
		if cc.List != nil {
			for _, expr := range cc.List {
				lit, val, ok := getValue(expr)
				if !ok {
					return
				}
				values = append(values, valueInfo{
					expr: expr,
					lit:  lit,
					val:  val,
				})
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

	// --- Divide into groups connected by fallthrough ---
	var groups [][]caseClauseInfo
	var currentGroup []caseClauseInfo
	for i, c := range cases {
		currentGroup = append(currentGroup, c)
		if !endsWithFallthrough(c.clause) || i == len(cases)-1 {
			groups = append(groups, currentGroup)
			currentGroup = nil
		}
	}

	changed := false
	var diagnostics []analysis.Diagnostic

	// --- Check and Fix within each case (sorting values) ---
	for i := range cases {
		if len(cases[i].values) < 2 {
			continue
		}
		if !isSorted(cases[i].values) {
			changed = true
			for j := 1; j < len(cases[i].values); j++ {
				if constant.Compare(cases[i].values[j-1].val, token.GTR, cases[i].values[j].val) {
					diagnostics = append(diagnostics, analysis.Diagnostic{
						Pos:     cases[i].values[j].expr.Pos(),
						End:     cases[i].values[j].expr.End(),
						Message: fmt.Sprintf("case value %s should come before %s", getLitValueString(cases[i].values[j].expr), getLitValueString(cases[i].values[j-1].expr)),
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

	// --- Check and Fix across groups ---
	if !isGroupsSorted(groups) {
		changed = true
		for i := 1; i < len(groups); i++ {
			if isLess(groups[i], groups[i-1]) {
				label := groupLabel(groups[i])
				prevLabel := groupLabel(groups[i-1])
				msg := fmt.Sprintf("case %s should come before %s", label, prevLabel)
				diagnostics = append(diagnostics, analysis.Diagnostic{
					Pos:     groups[i][0].clause.Pos(),
					End:     groups[i][0].clause.End(),
					Message: msg,
				})
			}
		}
		sortGroups(groups)
	}

	if !changed {
		return
	}

	// --- Build new body ---
	var newList []ast.Stmt
	for _, g := range groups {
		for _, c := range g {
			// Clone to reset position and avoid sparse formatting in go/format
			cloned := *c.clause
			cloned.Case = token.NoPos
			newList = append(newList, &cloned)
		}
	}

	var buf bytes.Buffer
	// Format a dummy block to get standard indentation and spacing
	if err := format.Node(&buf, pass.Fset, &ast.BlockStmt{List: newList}); err != nil {
		return
	}

	s := buf.String()
	// s is "{\n\tcase ...\n\tcase ...\n}"
	lines := strings.Split(s, "\n")
	if len(lines) < 2 {
		return
	}
	// content will be the lines between { and }
	content := strings.Join(lines[1:len(lines)-1], "\n")

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

func endsWithFallthrough(cc *ast.CaseClause) bool {
	if len(cc.Body) == 0 {
		return false
	}
	last := cc.Body[len(cc.Body)-1]
	br, ok := last.(*ast.BranchStmt)
	return ok && br.Tok == token.FALLTHROUGH
}

func isLess(a, b []caseClauseInfo) bool {
	// A group starting with default is considered "larger" than any group with values.
	if len(b[0].values) == 0 {
		return len(a[0].values) > 0
	}
	if len(a[0].values) == 0 {
		return false
	}
	return constant.Compare(a[0].values[0].val, token.LSS, b[0].values[0].val)
}

func groupLabel(g []caseClauseInfo) string {
	if len(g[0].values) == 0 {
		return "default"
	}
	return getLitValueString(g[0].values[0].expr)
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

func isGroupsSorted(groups [][]caseClauseInfo) bool {
	for i := 1; i < len(groups); i++ {
		if isLess(groups[i], groups[i-1]) {
			return false
		}
	}
	return true
}

func sortGroups(groups [][]caseClauseInfo) {
	sort.SliceStable(groups, func(i, j int) bool {
		return isLess(groups[i], groups[j])
	})
}
