package grimoire

import (
	"expinc/sunagent/command"
	"expinc/sunagent/common"
	"strings"
	"time"
)

// a specific command
type Spell interface {
	Cast(args ...string) ([]byte, error)
}

type spellImpl struct {
	program string
	args    []string
	timeout time.Duration
}

func newSpell(program string, args []string, timeout time.Duration) (spell Spell, err error) {
	err = nil
	if "" == strings.TrimSpace(program) {
		err = common.NewError(common.ErrorInvalidParameter, "No program specified for the spell")
		return
	}

	spell = &spellImpl{
		program: program,
		args:    args,
		timeout: timeout,
	}
	return
}

func (spell *spellImpl) Cast(args ...string) (output []byte, err error) {
	// {} in the spell will be replaced by the args in order
	// to specify {} literally in the args, use {{}}
	// example 1: "echo {}" with arg "hello" will get "echo hello"
	// example 2: "echo {{}}" will get "echo {}"
	// example 3: "echo {}" with no arg will get "echo "
	actualArgs := make([]string, len(spell.args), len(spell.args))
	i := 0
	for j := 0; j < len(actualArgs); j++ {
		if "{}" == spell.args[j] {
			if i < len(args) {
				actualArgs[j] = args[i]
				i++
			} else {
				actualArgs[j] = ""
			}
		} else if "{{}}" == spell.args[j] {
			actualArgs[j] = "{}"
		} else {
			actualArgs[j] = spell.args[j]
		}
	}

	output, err = command.CheckCombinedOutput(spell.program, actualArgs, spell.timeout)
	return
}
