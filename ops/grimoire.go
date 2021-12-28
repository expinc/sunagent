package ops

import (
	"context"
	"expinc/sunagent/common"
	"expinc/sunagent/grimoire"
	"expinc/sunagent/log"
	"expinc/sunagent/util"
	"fmt"
	"path/filepath"
)

var (
	opsGrimoire grimoire.Grimoire
)

func ReloadGrimoire(ctx context.Context) error {
	var err error
	grimoirePath := filepath.Join(common.CurrentDir, "grimoires", fmt.Sprintf("%s.yaml", nodeInfo.OsType))
	log.InfoCtx(ctx, "Reloading grimoire from "+grimoirePath)
	opsGrimoire, err = grimoire.NewGrimoireFromYamlFile(grimoirePath)
	util.LogErrorIfNotNilCtx(ctx, err)
	return err
}

func castGrimoireArcaneContext(ctx context.Context, arcaneName string, args ...string) (output []byte, err error) {
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

	output, err = spell.CastContext(ctx, args...)
	return
}

func castGrimoireArcane(arcaneName string, args ...string) (output []byte, err error) {
	output, err = castGrimoireArcaneContext(context.Background(), arcaneName, args...)
	return
}

func GetGrimoireAsYaml(ctx context.Context, osType string) (output []byte, err error) {
	log.InfoCtx(ctx, "Getting grimoire as yaml of OS type "+osType)
	if "default" == osType {
		output, err = grimoire.Grimoire2Yaml(opsGrimoire)
		util.LogErrorIfNotNilCtx(ctx, err)
	} else {
		grimoirePath := filepath.Join(common.CurrentDir, "grimoires", fmt.Sprintf("%s.yaml", osType))
		var theGrimoire grimoire.Grimoire
		theGrimoire, err = grimoire.NewGrimoireFromYamlFile(grimoirePath)
		if nil == err {
			output, err = grimoire.Grimoire2Yaml(theGrimoire)
		}
		util.LogErrorIfNotNilCtx(ctx, err)
	}
	return
}
