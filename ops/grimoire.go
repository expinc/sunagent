package ops

import (
	"context"
	"expinc/sunagent/common"
	"expinc/sunagent/grimoire"
	"expinc/sunagent/log"
	"expinc/sunagent/util"
	"fmt"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	GrimoireFolder string
	opsGrimoire    grimoire.Grimoire
)

func ReloadGrimoire(ctx context.Context) error {
	var err error
	grimoirePath := filepath.Join(GrimoireFolder, fmt.Sprintf("%s.yaml", nodeInfo.OsType))
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
		grimoirePath := filepath.Join(GrimoireFolder, fmt.Sprintf("%s.yaml", osType))
		var theGrimoire grimoire.Grimoire
		theGrimoire, err = grimoire.NewGrimoireFromYamlFile(grimoirePath)
		if nil == err {
			output, err = grimoire.Grimoire2Yaml(theGrimoire)
		}
		util.LogErrorIfNotNilCtx(ctx, err)
	}
	return
}

func SetGrimoireArcane(ctx context.Context, osType string, arcaneName string, yamlContent []byte) error {
	log.InfoCtx(ctx, fmt.Sprintf("Setting grimoire arcane %s of OS type %s", arcaneName, osType))

	// Read origin grimoire
	isDefault := false
	if "default" == osType {
		osType = nodeInfo.OsType
		isDefault = true
	}
	grimoirePath := filepath.Join(GrimoireFolder, fmt.Sprintf("%s.yaml", osType))
	theGrimoire, err := grimoire.NewGrimoireFromYamlFile(grimoirePath)
	if nil != err {
		log.ErrorCtx(ctx, err)
		return err
	}

	// Deserialize content
	var arcaneStruct grimoire.ArcaneStruct
	err = yaml.UnmarshalStrict(yamlContent, &arcaneStruct)
	if nil != err {
		log.ErrorCtx(ctx, err)
		return err
	}

	// Set arcane
	err = theGrimoire.SetArcane(arcaneName, time.Second*time.Duration(arcaneStruct.Timeout))
	if nil != err {
		log.ErrorCtx(ctx, err)
		return err
	}
	arcane, _ := theGrimoire.GetArcane(arcaneName)
	for spellIndex, spellStruct := range arcaneStruct.Spells {
		err = arcane.SetSpell(spellIndex, spellStruct.Args)
		if nil != err {
			log.ErrorCtx(ctx, err)
			return err
		}
	}

	// Write to file
	err = grimoire.WriteGrimioreToYamlFile(theGrimoire, grimoirePath)
	if nil != err {
		log.ErrorCtx(ctx, err)
		return err
	}

	// Final step
	if isDefault {
		opsGrimoire = theGrimoire
	}
	return nil
}

func RemoveGrimoireArcane(ctx context.Context, osType string, arcaneName string) error {
	log.InfoCtx(ctx, fmt.Sprintf("Removing grimoire arcane %s of OS type %s", arcaneName, osType))

	// Read origin grimoire
	isDefault := false
	if "default" == osType {
		osType = nodeInfo.OsType
		isDefault = true
	}
	grimoirePath := filepath.Join(GrimoireFolder, fmt.Sprintf("%s.yaml", osType))
	theGrimoire, err := grimoire.NewGrimoireFromYamlFile(grimoirePath)
	if nil != err {
		log.ErrorCtx(ctx, err)
		return err
	}

	// Remove arcane
	err = theGrimoire.RemoveArcane(arcaneName)
	if nil != err {
		log.ErrorCtx(ctx, err)
		return err
	}

	// Write to file
	err = grimoire.WriteGrimioreToYamlFile(theGrimoire, grimoirePath)
	if nil != err {
		log.ErrorCtx(ctx, err)
		return err
	}

	// Final step
	if isDefault {
		opsGrimoire = theGrimoire
	}
	return nil
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
