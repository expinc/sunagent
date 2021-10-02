package grimoire

import (
	"expinc/sunagent/common"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Grimoire interface {
	SetArcane(name string, timeout time.Duration) error
	GetArcane(name string) (Arcane, error)
}

type grimoireImpl struct {
	arcanes map[string]Arcane
}

func (grimoire *grimoireImpl) SetArcane(name string, timeout time.Duration) error {
	if "" == strings.TrimSpace(name) {
		return common.NewError(common.ErrorInvalidParameter, "No name specified for the arcane")
	}

	arcane, err := NewArcane(name, timeout)
	if nil == err {
		grimoire.arcanes[name] = arcane
	}
	return err
}

func (grimoire *grimoireImpl) GetArcane(name string) (arcane Arcane, err error) {
	var ok bool
	arcane, ok = grimoire.arcanes[name]
	if ok {
		err = nil
	} else {
		err = common.NewError(common.ErrorNotFound, fmt.Sprintf("No arcane \"%s\" in grimoire", name))
	}
	return
}

type SpellStruct struct {
	Args string `yaml:"args"`
}

type ArcaneStruct struct {
	Timeout uint                   `yaml:"timeout"`
	Spells  map[string]SpellStruct `yaml:"spells"`
}

type GrimoireStruct struct {
	Arcanes map[string]ArcaneStruct `yaml:"arcanes"`
}

func NewGrimoireFromYamlFile(path string) (grimoire Grimoire, err error) {
	var content []byte
	content, err = ioutil.ReadFile(path)
	if nil != err {
		return
	}

	grimoire, err = NewGrimoireFromYamlBytes(content)
	return
}

func NewGrimoireFromYamlBytes(bytes []byte) (grimoire Grimoire, err error) {
	var grimoireStruct GrimoireStruct
	err = yaml.UnmarshalStrict(bytes, &grimoireStruct)
	if nil != err {
		return
	}

	grimoire = &grimoireImpl{
		arcanes: make(map[string]Arcane),
	}
	for arcaneName, arcaneStruct := range grimoireStruct.Arcanes {
		err = grimoire.SetArcane(arcaneName, time.Second*time.Duration(arcaneStruct.Timeout))
		if nil != err {
			return
		}

		arcane, _ := grimoire.GetArcane(arcaneName)
		for spellIndex, spellStruct := range arcaneStruct.Spells {
			err = arcane.SetSpell(spellIndex, spellStruct.Args)
			if nil != err {
				return
			}
		}
	}
	return
}
