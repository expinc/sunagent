package grimoire

import (
	"expinc/sunagent/command"
	"expinc/sunagent/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArcane_Ordinary(t *testing.T) {
	_, err := NewArcane("arcane", command.DefaultTimeout)
	assert.Equal(t, nil, err)
}

func TestNewArcane_NoName(t *testing.T) {
	_, err := NewArcane(" \t\r\n", command.DefaultTimeout)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, common.ErrorInvalidParameter, err.(common.Error).Code())
}

func TestSetSpell_Ordinary(t *testing.T) {
	arcane, _ := NewArcane("arcane", command.DefaultTimeout)
	err := arcane.SetSpell("spell1", "echo hello")
	assert.Equal(t, nil, err)
	err = arcane.SetSpell("spell2", "date")
	assert.Equal(t, nil, err)
	err = arcane.SetSpell("spell3", "ls")
	assert.Equal(t, nil, err)

	var spell Spell
	spell, err = arcane.GetSpell("spell1")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, spell)
	spell, err = arcane.GetSpell("spell2")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, spell)
	spell, err = arcane.GetSpell("spell3")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, spell)
}

func TestSetSpell_ReplaceExisting(t *testing.T) {
	arcane, _ := NewArcane("arcane", command.DefaultTimeout)
	arcane.SetSpell("spell1", "echo hello")
	arcane.SetSpell("spell2", "date")
	arcane.SetSpell("spell3", "ls")

	// assert original spell
	spell, err := arcane.GetSpell("spell2")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, spell)
	impl, ok := spell.(*spellImpl)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(impl.args))
	assert.Equal(t, "date", impl.args[0])

	// update spell and assert new
	arcane.SetSpell("spell2", "top")
	spell, err = arcane.GetSpell("spell2")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, spell)
	impl, ok = spell.(*spellImpl)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(impl.args))
	assert.Equal(t, "top", impl.args[0])
}

func TestSetSpell_Empty(t *testing.T) {
	arcane, _ := NewArcane("arcane", command.DefaultTimeout)
	err := arcane.SetSpell("spell1", " \t\r\n")
	assert.Equal(t, common.ErrorInvalidParameter, err.(common.Error).Code())
}

func TestGetSpell_Ordinary(t *testing.T) {
	arcane, _ := NewArcane("arcane", command.DefaultTimeout)
	arcane.SetSpell("spell1", "echo hello")
	arcane.SetSpell("spell2", "date")
	arcane.SetSpell("spell3", "")

	spell, err := arcane.GetSpell("spell2")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, spell)
	impl, ok := spell.(*spellImpl)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(impl.args))
	assert.Equal(t, "date", impl.args[0])
}

func TestGetSpell_NonExisting(t *testing.T) {
	arcane, _ := NewArcane("arcane", command.DefaultTimeout)
	arcane.SetSpell("spell1", "echo hello")
	arcane.SetSpell("spell2", "date")
	arcane.SetSpell("spell3", "")

	spell, err := arcane.GetSpell("nonexist")
	assert.Equal(t, common.ErrorNotFound, err.(common.Error).Code())
	assert.Equal(t, nil, spell)
}
