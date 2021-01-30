package controller

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

type Matcher struct{}

func (matcher *Matcher) Compile(rule string) (CompiledRule, error) {
	prog := CompiledRule{}
	prg, err := expr.Compile(rule, expr.AsBool())
	if err != nil {
		return prog, fmt.Errorf("compile a rule: "+rule+": %w", err)
	}
	prog.prg = prg
	return prog, nil
}

type CompiledRule struct {
	prg *vm.Program
}

func (cr *CompiledRule) Match(rsc interface{}) (bool, error) {
	output, err := expr.Run(cr.prg, rsc)
	if err != nil {
		return false, fmt.Errorf("evaluate an expression with params: %w", err)
	}
	if f, ok := output.(bool); !ok || !f {
		return false, nil
	}
	return true, nil
}
