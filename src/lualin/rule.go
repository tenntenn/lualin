package lualin

import (
	"fmt"
	"regexp"

	"github.com/yuin/gopher-lua/ast"
)

type Rule interface {
	Validate(l *Lualin, stmt ast.Stmt) error
}

type RuleFunc func(l *Lualin, stmt ast.Stmt) error

func (f RuleFunc) Validate(l *Lualin, stmt ast.Stmt) error {
	return f(l, stmt)
}

type VarName struct {
	Regexp *regexp.Regexp
}

func (v *VarName) Validate(l *Lualin, stmt ast.Stmt) error {

	errs := []*LintError{}
	switch stmt.(type) {
	case *ast.AssignStmt:
		s, _ := stmt.(*ast.AssignStmt)
		for _, lh := range s.Lhs {
			if le, ok := lh.(*ast.IdentExpr); ok && !v.Regexp.MatchString(le.Value) {
				errs = append(errs, &LintError{
					Line:    s.Line(),
					Message: fmt.Sprintf("%s is invalid name", le.Value),
				})
			}
		}
	case *ast.LocalAssignStmt:
		s, _ := stmt.(*ast.LocalAssignStmt)
		for _, name := range s.Names {
			if !v.Regexp.MatchString(name) {
				errs = append(errs, &LintError{
					Line:    s.Line(),
					Message: fmt.Sprintf("%s is invalid var name", name),
				})
			}
		}
	}

	if len(errs) >= 0 {
		return LintErrors(errs)
	}
	return nil
}

type FuncName struct {
	Regexp *regexp.Regexp
}

func (v *FuncName) Validate(l *Lualin, stmt ast.Stmt) error {

	switch stmt.(type) {
	case *ast.FuncDefStmt:
		s, _ := stmt.(*ast.FuncDefStmt)
		if le, ok := s.Name.Func.(*ast.IdentExpr); ok && !v.Regexp.MatchString(le.Value) {
			return &LintError{
				Line:    s.Line(),
				Message: fmt.Sprintf("%s is invalid func name", le.Value),
			}
		}
	}

	return nil
}