package switchorder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/token"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Config holds configuration for the switchorder analyzer.
type Config struct {
	Order       string        `json:"order"`
	Comparators []Comparator  `json:"comparators"`
	DefaultLast bool          `json:"default-last"`
	Autofix     AutofixConfig `json:"autofix"`
}

// Comparator defines how to compare case values of a given type.
type Comparator struct {
	Type       string `json:"type"`
	IgnoreCase bool   `json:"ignore-case"`
}

// AutofixConfig controls the behavior of suggested fixes.
type AutofixConfig struct {
	Enabled          bool `json:"enabled"`
	AllowFallthrough bool `json:"allow-fallthrough"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		Order: "asc",
		Comparators: []Comparator{
			{Type: "numeric"},
			{Type: "alphabetical", IgnoreCase: true},
		},
		DefaultLast: true,
		Autofix: AutofixConfig{
			Enabled:          true,
			AllowFallthrough: false,
		},
	}
}

var (
	defaultCfg = DefaultConfig()
	Analyzer   = NewWithConfig(&defaultCfg)
)

// New is required for golangci-lint compatibility.
func New(conf any) ([]*analysis.Analyzer, error) {
	cfg := DefaultConfig()
	if conf != nil {
		data, err := json.Marshal(conf)
		if err != nil {
			return nil, fmt.Errorf("switchorder: marshal config: %w", err)
		}
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("switchorder: unmarshal config: %w", err)
		}
	}
	return []*analysis.Analyzer{NewWithConfig(&cfg)}, nil
}

// NewWithConfig creates an analyzer with the given configuration pointer.
// Flags registered on the returned analyzer's FlagSet write directly into cfg,
// so their values are visible to Run without any extra wiring.
func NewWithConfig(cfg *Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "switchorder",
		Doc:  "checks that switch case statements are in alphabetical or numerical order",
		Run:  makeRun(cfg),
	}
}

func makeRun(cfg *Config) func(*analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		// Snapshot config once per Run call, after flags have been parsed.
		c := *cfg
		for _, file := range pass.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				switchStmt, ok := n.(*ast.SwitchStmt)
				if !ok {
					return true
				}
				processSwitch(pass, switchStmt, c)
				return true
			})
		}
		return nil, nil
	}
}

type valueKind int

const (
	kindNumeric valueKind = iota
	kindAlphabetical
)

type valueInfo struct {
	expr ast.Expr
	lit  *ast.BasicLit
	val  constant.Value
	kind valueKind
}

type caseClauseInfo struct {
	clause *ast.CaseClause
	values []valueInfo
}

func getValue(expr ast.Expr) (*ast.BasicLit, constant.Value, valueKind, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		switch e.Kind {
		case token.STRING:
			return e, constant.MakeFromLiteral(e.Value, e.Kind, 0), kindAlphabetical, true
		case token.INT, token.FLOAT, token.CHAR:
			return e, constant.MakeFromLiteral(e.Value, e.Kind, 0), kindNumeric, true
		}
	case *ast.UnaryExpr:
		if bl, ok := e.X.(*ast.BasicLit); ok {
			if bl.Kind == token.INT || bl.Kind == token.FLOAT {
				val := constant.MakeFromLiteral(bl.Value, bl.Kind, 0)
				if e.Op == token.SUB {
					val = constant.UnaryOp(token.SUB, val, 0)
				}
				return bl, val, kindNumeric, true
			}
		}
	}
	return nil, nil, 0, false
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

func compareValues(a, b valueInfo, cfg Config) int {
	for _, comp := range cfg.Comparators {
		switch comp.Type {
		case "numeric":
			if a.kind == kindNumeric && b.kind == kindNumeric {
				if constant.Compare(a.val, token.LSS, b.val) {
					return -1
				}
				if constant.Compare(a.val, token.GTR, b.val) {
					return 1
				}
				return 0
			}
		case "alphabetical":
			if a.kind == kindAlphabetical && b.kind == kindAlphabetical {
				aStr := constant.StringVal(a.val)
				bStr := constant.StringVal(b.val)
				if comp.IgnoreCase {
					aStr = strings.ToLower(aStr)
					bStr = strings.ToLower(bStr)
				}
				return strings.Compare(aStr, bStr)
			}
		}
	}
	// Fallback for unrecognized types
	if constant.Compare(a.val, token.LSS, b.val) {
		return -1
	}
	if constant.Compare(a.val, token.GTR, b.val) {
		return 1
	}
	return 0
}

func valueLess(a, b valueInfo, cfg Config) bool {
	cmp := compareValues(a, b, cfg)
	if cfg.Order == "desc" {
		return cmp > 0
	}
	return cmp < 0
}

func groupLess(a, b []caseClauseInfo, cfg Config) bool {
	aIsDefault := len(a[0].values) == 0
	bIsDefault := len(b[0].values) == 0

	if cfg.DefaultLast {
		if bIsDefault {
			return !aIsDefault
		}
		if aIsDefault {
			return false
		}
	} else {
		if aIsDefault {
			return !bIsDefault
		}
		if bIsDefault {
			return false
		}
	}

	return valueLess(a[0].values[0], b[0].values[0], cfg)
}

func groupLabel(g []caseClauseInfo) string {
	if len(g[0].values) == 0 {
		return "default"
	}
	return getLitValueString(g[0].values[0].expr)
}

func isSorted(values []valueInfo, cfg Config) bool {
	for i := 1; i < len(values); i++ {
		if valueLess(values[i], values[i-1], cfg) {
			return false
		}
	}
	return true
}

func sortValues(values []valueInfo, cfg Config) {
	sort.Slice(values, func(i, j int) bool {
		return valueLess(values[i], values[j], cfg)
	})
}

func isGroupsSorted(groups [][]caseClauseInfo, cfg Config) bool {
	for i := 1; i < len(groups); i++ {
		if groupLess(groups[i], groups[i-1], cfg) {
			return false
		}
	}
	return true
}

func sortGroups(groups [][]caseClauseInfo, cfg Config) {
	sort.SliceStable(groups, func(i, j int) bool {
		return groupLess(groups[i], groups[j], cfg)
	})
}

func hasFallthroughGroups(groups [][]caseClauseInfo) bool {
	for _, g := range groups {
		if len(g) > 1 {
			return true
		}
	}
	return false
}

func processSwitch(pass *analysis.Pass, sw *ast.SwitchStmt, cfg Config) {
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
				lit, val, kind, ok := getValue(expr)
				if !ok {
					return
				}
				values = append(values, valueInfo{
					expr: expr,
					lit:  lit,
					val:  val,
					kind: kind,
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
		if !isSorted(cases[i].values, cfg) {
			changed = true
			for j := 1; j < len(cases[i].values); j++ {
				if valueLess(cases[i].values[j], cases[i].values[j-1], cfg) {
					diagnostics = append(diagnostics, analysis.Diagnostic{
						Pos:     cases[i].values[j].expr.Pos(),
						End:     cases[i].values[j].expr.End(),
						Message: fmt.Sprintf("case value %s should come before %s", getLitValueString(cases[i].values[j].expr), getLitValueString(cases[i].values[j-1].expr)),
					})
				}
			}
			sortValues(cases[i].values, cfg)
			newList := make([]ast.Expr, len(cases[i].values))
			for j, v := range cases[i].values {
				newList[j] = v.expr
			}
			cases[i].clause.List = newList
		}
	}

	// --- Check and Fix across groups ---
	if !isGroupsSorted(groups, cfg) {
		changed = true
		for i := 1; i < len(groups); i++ {
			if groupLess(groups[i], groups[i-1], cfg) {
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
		sortGroups(groups, cfg)
	}

	if !changed {
		return
	}

	// --- Build suggested fix ---
	var fix *analysis.SuggestedFix
	if cfg.Autofix.Enabled && (!hasFallthroughGroups(groups) || cfg.Autofix.AllowFallthrough) {
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
		if err := format.Node(&buf, pass.Fset, &ast.BlockStmt{List: newList}); err == nil {
			s := buf.String()
			// s is "{\n\tcase ...\n\tcase ...\n}"
			lines := strings.Split(s, "\n")
			if len(lines) >= 2 {
				// content will be the lines between { and }
				content := strings.Join(lines[1:len(lines)-1], "\n")
				fix = &analysis.SuggestedFix{
					Message: "reorder switch cases",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     sw.Body.Lbrace + 1,
							End:     sw.Body.Rbrace,
							NewText: []byte("\n" + content + "\n"),
						},
					},
				}
			}
		}
	}

	for _, d := range diagnostics {
		if fix != nil {
			d.SuggestedFixes = []analysis.SuggestedFix{*fix}
		}
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
