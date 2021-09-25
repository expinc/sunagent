package grimoire

import (
	"expinc/sunagent/common"
	"fmt"
	"strings"
	"time"
)

// a group of spells for a specific purpose
type Arcane interface {
	SetSpell(index string, args string) error
	GetSpell(index string) (Spell, error)
}

type arcaneImpl struct {
	name    string
	program string
	timeout time.Duration
	spells  map[string]Spell
}

func (arcane *arcaneImpl) SetSpell(index string, args string) error {
	spell, err := newSpell(arcane.program, strings.Fields(args), arcane.timeout)
	if nil == err {
		arcane.spells[index] = spell
	}
	return err
}

func (arcane *arcaneImpl) GetSpell(index string) (spell Spell, err error) {
	var ok bool
	spell, ok = arcane.spells[index]
	if ok {
		err = nil
	} else {
		err = common.NewError(common.ErrorNotFound, fmt.Sprintf("No spell \"%s\" in arcane \"%s\"", index, arcane.name))
	}
	return
}

func NewArcane(name string, program string, timeout time.Duration) (arcane Arcane, err error) {
	err = nil
	if "" == strings.TrimSpace(name) {
		err = common.NewError(common.ErrorInvalidParameter, "No name specified for the arcane")
		return
	}
	if "" == strings.TrimSpace(program) {
		err = common.NewError(common.ErrorInvalidParameter, "No program specified for the arcane")
		return
	}

	arcane = &arcaneImpl{
		name:    name,
		program: program,
		timeout: timeout,
		spells:  make(map[string]Spell),
	}
	return
}
