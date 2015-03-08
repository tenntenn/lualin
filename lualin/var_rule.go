package lualin

import (
	"fmt"
	"regexp"

	"github.com/yuin/gopher-lua/ast"
)

type NoGlobalVarRule struct {
	WhiteList []*regexp.Regexp
	RuleLevel RuleLevel
}

func (r *NoGlobalVarRule) validate(s *ast.AssignStmt, le *ast.IdentExpr) *LintError {

	for _, v := range r.WhiteList {
		if v.MatchString(le.Value) {
			return nil
		}
	}

	return &LintError{
		Rule:    r,
		Line:    s.Line(),
		Message: fmt.Sprintf("%s is invalid global var", le.Value),
	}
}

func (r *NoGlobalVarRule) Validate(l *Lualin, stmt ast.Stmt) error {

	errs := []*LintError{}
	switch stmt.(type) {
	case *ast.AssignStmt:
		s, _ := stmt.(*ast.AssignStmt)

		for _, lh := range s.Lhs {
			if le, ok := lh.(*ast.IdentExpr); ok {
				err := r.validate(s, le)
				if err != nil {
					errs = append(errs, err)
				}
			}
		}
	}

	if len(errs) >= 0 {
		return LintErrors(errs)
	}
	return nil
}

func (r *NoGlobalVarRule) Level() RuleLevel {
	return r.RuleLevel
}
