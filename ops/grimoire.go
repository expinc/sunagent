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

func CastGrimoireArcaneContext(ctx context.Context, arcaneName string, args ...string) (output []byte, err error) {
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
	output, err = CastGrimoireArcaneContext(context.Background(), arcaneName, args...)
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

type CastArcaneJob struct {
	jobBase

	cancelFunc context.CancelFunc
	canceled   bool
}

func (job *CastArcaneJob) execute() error {
	// Get parameters
	arcaneName := job.jobBase.params["arcaneName"].(string)
	args := job.jobBase.params["args"].([]string)

	// Execute casting
	output, err := CastGrimoireArcaneContext(job.jobBase.ctx, arcaneName, args...)
	if job.canceled {
		// If the job is canceled, it should not be considered as failed job.
		// Therefore override the error as nil if any
		err = nil
	}

	// Render result
	result := CombinedScriptResult{
		Output: string(output),
	}
	if nil != err {
		errMsg := fmt.Sprintf("Execute script failed: err=%s, output=%s", err.Error(), output)
		err = common.NewError(common.ErrorUnexpected, errMsg)
	}
	job.jobBase.getInfo().Result = result
	return err
}

func (job *CastArcaneJob) cancel() {
	job.canceled = true
	job.getInfo().Status = JobStatusCanceled
	job.cancelFunc()
}

func (job *CastArcaneJob) dispose() {}
