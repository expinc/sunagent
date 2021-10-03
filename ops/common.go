package ops

import (
	"expinc/sunagent/common"
	"expinc/sunagent/grimoire"
	"fmt"
	"path/filepath"
)

var opsGrimoire grimoire.Grimoire

func ReloadGrimoireFromFile() error {
	var err error
	grimoireFile := fmt.Sprintf("%s.yaml", nodeInfo.OsType)
	grimoirePath := filepath.Join(common.CurrentDir, "grimoires", grimoireFile)
	opsGrimoire, err = grimoire.NewGrimoireFromYamlFile(grimoirePath)
	return err
}

func castGrimoireArcane(arcaneName string, args ...string) (output []byte, err error) {
	var arcane grimoire.Arcane
	arcane, err = opsGrimoire.GetArcane(arcaneName)
	if nil != err {
		return
	}

	var spell grimoire.Spell
	spell, err = arcane.GetSpell(nodeInfo.OsFamily)
	if nil != err {
		return
	}

	output, err = spell.Cast(args...)
	return
}
