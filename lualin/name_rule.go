package lualin

import (
	"fmt"
	"regexp"

	"github.com/yuin/gopher-lua/ast"
)

type LocalVarNameRule struct {
	Regexp    *regexp.Regexp
	FuncSkip  bool
	WhiteList []*regexp.Regexp
	RuleLevel RuleLevel
}

func (r *LocalVarNameRule) Validate(l *Lualin, stmt ast.Stmt) error {

	errs := []*LintError{}
	switch stmt.(type) {
	case *ast.LocalAssignStmt:
		s, _ := stmt.(*ast.LocalAssignStmt)
		for i, name := range s.Names {
			if _, isfunc := s.Exprs[i].(*ast.FunctionExpr); r.FuncSkip && isfunc {
				continue
			}

			if matchWhiteList(r.WhiteList, name) {
				continue
			}

			if !r.Regexp.MatchString(name) {
				errs = append(errs, &LintError{
					Rule:    r,
					Line:    s.Line(),
					Message: fmt.Sprintf("%s is invalid local var name", name),
				})
			}
		}
	}

	if len(errs) >= 0 {
		return LintErrors(errs)
	}
	return nil
}

func (r *LocalVarNameRule) Level() RuleLevel {
	return r.RuleLevel
}

type GlobalVarNameRule struct {
	Regexp    *regexp.Regexp
	WhiteList []*regexp.Regexp
	RuleLevel RuleLevel
}

func (r *GlobalVarNameRule) Validate(l *Lualin, stmt ast.Stmt) error {

	errs := []*LintError{}
	switch stmt.(type) {
	case *ast.AssignStmt:
		s, _ := stmt.(*ast.AssignStmt)
		for _, lh := range s.Lhs {
			if le, ok := lh.(*ast.IdentExpr); ok && !r.Regexp.MatchString(le.Value) {
				if matchWhiteList(r.WhiteList, le.Value) {
					continue
				}

				errs = append(errs, &LintError{
					Rule:    r,
					Line:    s.Line(),
					Message: fmt.Sprintf("%s is invalid global var name", le.Value),
				})
			}
		}
	}

	if len(errs) >= 0 {
		return LintErrors(errs)
	}
	return nil
}

func (r *GlobalVarNameRule) Level() RuleLevel {
	return r.RuleLevel
}

type FuncNameRule struct {
	Regexp    *regexp.Regexp
	WhiteList []*regexp.Regexp
	RuleLevel RuleLevel
}

func (r *FuncNameRule) Validate(l *Lualin, stmt ast.Stmt) error {

	switch stmt.(type) {
	case *ast.FuncDefStmt:
		s, _ := stmt.(*ast.FuncDefStmt)
		if le, ok := s.Name.Func.(*ast.IdentExpr); ok && !r.Regexp.MatchString(le.Value) {

			if matchWhiteList(r.WhiteList, le.Value) {
				return nil
			}

			return &LintError{
				Rule:    r,
				Line:    s.Line(),
				Message: fmt.Sprintf("%s is invalid func name", le.Value),
			}
		}
	}

	return nil
}

func (r *FuncNameRule) Level() RuleLevel {
	return r.RuleLevel
}
