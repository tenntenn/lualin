package lualin

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/yuin/gopher-lua/ast"
)

var DefaultRules = []Rule{
	&LocalVarNameRule{
		Regexp:   regexp.MustCompile("^[a-z_][a-z0-9_]*$"),
		FuncSkip: true,
	},
	&GlobalVarNameRule{
		Regexp: regexp.MustCompile("^[A-Z_][A-Z0-9_]*$"),
	},
	&FuncNameRule{
		Regexp: regexp.MustCompile("^[a-z]+([A-Z][a-z0-9]+)*$"),
	},
	&NoGlobalVarRule{},
}

type LintError struct {
	Rule    Rule
	Line    int
	Message string
}

func (err *LintError) Error() string {
	return fmt.Sprintf("%d: %s", err.Line, err.Message)
}

type LintErrors []*LintError

func (errs LintErrors) Error() string {
	msgs := make([]string, 0, len(errs))
	for _, err := range errs {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "\n")
}

func Lint(w io.Writer, chunk []ast.Stmt) error {
	l := NewLualin(w, DefaultRules)
	return l.Lint(chunk)
}

type Scope struct{}

type Lualin struct {
	writer io.Writer
	rules  []Rule
}

func NewLualin(w io.Writer, r []Rule) *Lualin {
	return &Lualin{
		writer: w,
		rules:  r,
	}
}

func (l *Lualin) Lint(chunk []ast.Stmt) error {
	errs := []*LintError{}
	for _, stmt := range chunk {
		err := l.lint(stmt)
		if err == nil {
			continue
		}

		switch err.(type) {
		case *LintError:
			lintErr, _ := err.(*LintError)
			errs = append(errs, lintErr)
		case LintErrors:
			lintErrs, _ := err.(LintErrors)
			errs = append(errs, lintErrs...)
		default:
			return err
		}
	}

	if len(errs) > 0 {
		return LintErrors(errs)
	}
	return nil
}

func (l *Lualin) lint(stmt ast.Stmt) error {

	errs := []*LintError{}

	for _, r := range l.rules {
		err := r.Validate(l, stmt)

		if err == nil {
			continue
		}

		switch err.(type) {
		case *LintError:
			lintErr, _ := err.(*LintError)
			errs = append(errs, lintErr)
		case LintErrors:
			lintErrs, _ := err.(LintErrors)
			errs = append(errs, lintErrs...)
		default:
			return err
		}
	}

	var stmts []ast.Stmt
	switch stmt.(type) {
	case *ast.DoBlockStmt:
		s, _ := stmt.(*ast.DoBlockStmt)
		stmts = s.Stmts
	case *ast.GenericForStmt:
		s, _ := stmt.(*ast.GenericForStmt)
		stmts = s.Stmts
	case *ast.NumberForStmt:
		s, _ := stmt.(*ast.NumberForStmt)
		stmts = s.Stmts
	case *ast.RepeatStmt:
		s, _ := stmt.(*ast.RepeatStmt)
		stmts = s.Stmts
	case *ast.WhileStmt:
		s, _ := stmt.(*ast.WhileStmt)
		stmts = s.Stmts
	case *ast.FuncDefStmt:
		s, _ := stmt.(*ast.FuncDefStmt)
		stmts = s.Func.Stmts
	}

	if stmts != nil {
		err := l.Lint(stmts)

		if err != nil {
			switch err.(type) {
			case *LintError:
				lintErr, _ := err.(*LintError)
				errs = append(errs, lintErr)
			case LintErrors:
				lintErrs, _ := err.(LintErrors)
				errs = append(errs, lintErrs...)
			default:
				return err
			}
		}
	}

	if len(errs) > 0 {
		return LintErrors(errs)
	}
	return nil
}
