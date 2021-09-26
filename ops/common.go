package ops

import (
	"expinc/sunagent/common"
	"expinc/sunagent/grimoire"
	"fmt"
	"path"
)

var opsGrimoire grimoire.Grimoire

func ReloadGrimoireFromFile() error {
	var err error
	grimoireFile := fmt.Sprintf("%s.yaml", nodeInfo.OsType)
	grimoirePath := path.Join(common.CurrentDir, "grimoires", grimoireFile)
	opsGrimoire, err = grimoire.NewGrimoireFromYamlFile(grimoirePath)
	return err
}
