package lualin

import (
	"regexp"

	"github.com/yuin/gopher-lua/ast"
)

type RuleLevel int

const (
	Error RuleLevel = iota
	Warning
)

type Rule interface {
	Validate(l *Lualin, stmt ast.Stmt) error
	Level() RuleLevel
}

type RuleFunc func(l *Lualin, stmt ast.Stmt) error

func (f RuleFunc) Validate(l *Lualin, stmt ast.Stmt) error {
	return f(l, stmt)
}

func matchWhiteList(wl []*regexp.Regexp, str string) bool {
	if wl == nil {
		return false
	}

	for _, w := range wl {
		if w.MatchString(str) {
			return true
		}
	}
	return false
}
