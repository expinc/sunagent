package grimoire

import (
	"expinc/sunagent/common"
	"fmt"
	"io/fs"
	"io/ioutil"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Grimoire interface {
	SetArcane(name string, timeout time.Duration) error
	GetArcane(name string) (Arcane, error)
	RemoveArcane(name string) error
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

func (grimoire *grimoireImpl) RemoveArcane(name string) error {
	_, ok := grimoire.arcanes[name]
	if ok {
		delete(grimoire.arcanes, name)
		return nil
	} else {
		return common.NewError(common.ErrorNotFound, fmt.Sprintf("No arcane \"%s\" in grimoire", name))
	}
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

func Grimoire2Yaml(grimoire Grimoire) (yamlBytes []byte, err error) {
	grimoireImpl := grimoire.(*grimoireImpl)
	grimoireStruct := GrimoireStruct{
		Arcanes: make(map[string]ArcaneStruct),
	}
	for key, val := range grimoireImpl.arcanes {
		arcane := val.(*arcaneImpl)
		arcaneStruct := ArcaneStruct{
			Timeout: uint(arcane.timeout / time.Second),
			Spells:  make(map[string]SpellStruct),
		}

		for key2, val2 := range arcane.spells {
			spell := val2.(*spellImpl)
			spellStruct := SpellStruct{
				Args: strings.Join(spell.args, " "),
			}
			arcaneStruct.Spells[key2] = spellStruct
		}
		grimoireStruct.Arcanes[key] = arcaneStruct
	}

	yamlBytes, err = yaml.Marshal(grimoireStruct)
	return
}

func WriteGrimioreToYamlFile(grimoire Grimoire, filePath string) error {
	// Transform to yaml
	yamlBytes, err := Grimoire2Yaml(grimoire)
	if nil != err {
		return err
	}

	// Write to file
	err = ioutil.WriteFile(filePath, yamlBytes, fs.ModePerm)
	return err
}
