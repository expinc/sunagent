package grimoire

import (
	"expinc/sunagent/command"
	"expinc/sunagent/common"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var (
	ordinaryGrimoire string
	invalidYaml      string
	invalidGrimoire  string
	invalidArcane    string
	invalidSpell     string
)

func init() {
	ordinaryGrimoire = filepath.Join(os.Getenv("TEST_ARTIFACT_DIR"), "functionality", "data", "grimoires", "ordinary.yaml")
	invalidYaml = filepath.Join(os.Getenv("TEST_ARTIFACT_DIR"), "functionality", "data", "grimoires", "invalid-yaml.yaml")
	invalidGrimoire = filepath.Join(os.Getenv("TEST_ARTIFACT_DIR"), "functionality", "data", "grimoires", "invalid-grimoire.yaml")
	invalidArcane = filepath.Join(os.Getenv("TEST_ARTIFACT_DIR"), "functionality", "data", "grimoires", "invalid-arcane.yaml")
	invalidSpell = filepath.Join(os.Getenv("TEST_ARTIFACT_DIR"), "functionality", "data", "grimoires", "invalid-spell.yaml")
}

func TestGrimoire_SetArcane_Ordinary(t *testing.T) {
	grimoire := grimoireImpl{
		arcanes: make(map[string]Arcane),
	}
	err := grimoire.SetArcane("arcane1", command.DefaultTimeout)
	assert.Equal(t, nil, err)
	err = grimoire.SetArcane("arcane2", command.DefaultTimeout)
	assert.Equal(t, nil, err)
	err = grimoire.SetArcane("arcane3", command.DefaultTimeout)
	assert.Equal(t, nil, err)
}

func TestGrimoire_SetArcane_NoName(t *testing.T) {
	grimoire := grimoireImpl{
		arcanes: make(map[string]Arcane),
	}
	err := grimoire.SetArcane(" \t\r\n", command.DefaultTimeout)
	assert.Equal(t, common.ErrorInvalidParameter, err.(common.Error).Code())
}

func TestGrimoire_GetArcane_Ordinary(t *testing.T) {
	grimoire := grimoireImpl{
		arcanes: make(map[string]Arcane),
	}
	grimoire.SetArcane("arcane1", command.DefaultTimeout)
	grimoire.SetArcane("arcane2", command.DefaultTimeout)
	grimoire.SetArcane("arcane3", command.DefaultTimeout)

	arcane, err := grimoire.GetArcane("arcane1")
	assert.Equal(t, nil, err)
	assert.Equal(t, "arcane1", arcane.(*arcaneImpl).name)
	arcane, err = grimoire.GetArcane("arcane2")
	assert.Equal(t, nil, err)
	assert.Equal(t, "arcane2", arcane.(*arcaneImpl).name)
	arcane, err = grimoire.GetArcane("arcane3")
	assert.Equal(t, nil, err)
	assert.Equal(t, "arcane3", arcane.(*arcaneImpl).name)
}

func TestGrimoire_GetArcane_NotExist(t *testing.T) {
	grimoire := grimoireImpl{
		arcanes: make(map[string]Arcane),
	}
	grimoire.SetArcane("arcane1", command.DefaultTimeout)
	grimoire.SetArcane("arcane2", command.DefaultTimeout)
	grimoire.SetArcane("arcane3", command.DefaultTimeout)

	arcane, err := grimoire.GetArcane("nonexist")
	assert.Equal(t, common.ErrorNotFound, err.(common.Error).Code())
	assert.Equal(t, nil, arcane)
}

func TestGrimoire_SetArcane_ReplaceExisting(t *testing.T) {
	grimoire := grimoireImpl{
		arcanes: make(map[string]Arcane),
	}
	grimoire.SetArcane("arcane1", command.DefaultTimeout)
	grimoire.SetArcane("arcane2", command.DefaultTimeout)
	grimoire.SetArcane("arcane3", command.DefaultTimeout)

	arcane, err := grimoire.GetArcane("arcane2")
	assert.Equal(t, nil, err)
	assert.Equal(t, command.DefaultTimeout, arcane.(*arcaneImpl).timeout)

	err = grimoire.SetArcane("arcane2", 9*time.Second)
	assert.Equal(t, nil, err)

	arcane, err = grimoire.GetArcane("arcane2")
	assert.Equal(t, nil, err)
	assert.Equal(t, 9*time.Second, arcane.(*arcaneImpl).timeout)
}

func TestNewGrimoireFromYamlBytes_Ordinary(t *testing.T) {
	grimoire, err := NewGrimoireFromYamlFile(ordinaryGrimoire)
	assert.Equal(t, nil, err)

	var arcane Arcane
	var spell Spell
	arcane, err = grimoire.GetArcane("get-package")
	assert.Equal(t, time.Second*60, arcane.(*arcaneImpl).timeout)
	spell, err = arcane.GetSpell("suse")
	assert.Equal(t, []string{"rpm", "-qi", "{}"}, spell.(*spellImpl).args)

	arcane, err = grimoire.GetArcane("install-package")
	assert.Equal(t, time.Second*600, arcane.(*arcaneImpl).timeout)
	spell, err = arcane.GetSpell("suse")
	assert.Equal(t, []string{"zypper", "-n", "install", "{}"}, spell.(*spellImpl).args)
}

func TestNewGrimoireFromYamlBytes_EmptyBytes(t *testing.T) {
	grimoire, err := NewGrimoireFromYamlBytes([]byte{})
	assert.Equal(t, nil, err)
	countArcanes := len(grimoire.(*grimoireImpl).arcanes)
	assert.Equal(t, 0, countArcanes)
}

func TestNewGrimoireFromYamlBytes_InvalidYaml(t *testing.T) {
	_, err := NewGrimoireFromYamlFile(invalidYaml)
	_, ok := err.(*yaml.TypeError)
	assert.Equal(t, true, ok)
}

func TestNewGrimoireFromYamlBytes_InvalidGrimoire(t *testing.T) {
	_, err := NewGrimoireFromYamlFile(invalidGrimoire)
	_, ok := err.(*yaml.TypeError)
	assert.Equal(t, true, ok)
}

func TestNewGrimoireFromYamlBytes_InvalidArcane(t *testing.T) {
	_, err := NewGrimoireFromYamlFile(invalidArcane)
	_, ok := err.(*yaml.TypeError)
	assert.Equal(t, true, ok)
}

func TestNewGrimoireFromYamlBytes_InvalidSpell(t *testing.T) {
	_, err := NewGrimoireFromYamlFile(invalidSpell)
	_, ok := err.(*yaml.TypeError)
	assert.Equal(t, true, ok)
}
