package grimoire

import (
	"expinc/sunagent/command"
	"os"
	"path/filepath"
	"testing"

	"github.com/shirou/gopsutil/host"
	"github.com/stretchr/testify/assert"
)

var echoScript string
var hostInfo *host.InfoStat

func init() {
	hostInfo, _ = host.Info()
	echoScript = filepath.Join(os.Getenv("TEST_ARTIFACT_DIR"), "functionality", "exe", "echo.py")
}

func TestCast_NoArg(t *testing.T) {
	var (
		spell Spell
		err   error
	)
	if "windows" == hostInfo.OS {
		spell, err = newSpell([]string{"tasklist"}, command.DefaultTimeout)
	} else {
		spell, err = newSpell([]string{"echo"}, command.DefaultTimeout)
	}
	assert.Equal(t, nil, err)

	output, err := spell.Cast()
	assert.Equal(t, nil, err)
	if "windows" == hostInfo.OS {
		// the output would be very long
		// here we only check the beginning content
		assert.Equal(t, "\r\nImage Name", string(output)[:12])
	} else {
		assert.Equal(t, "\n", string(output))
	}
}

func TestCast_WithArgs(t *testing.T) {
	spell, err := newSpell([]string{"python3", echoScript, "{}", "{}"}, command.DefaultTimeout)
	assert.Equal(t, nil, err)

	// three positional args. only specify first two
	output, err := spell.Cast("hello", "world")
	assert.Equal(t, nil, err)
	if "windows" == hostInfo.OS {
		assert.Equal(t, "hello world\r\n", string(output))
	} else {
		assert.Equal(t, "hello world\n", string(output))
	}
}

func TestCast_WithBraces(t *testing.T) {
	spell, err := newSpell([]string{"python3", echoScript, "{{}}", "}{"}, command.DefaultTimeout)
	assert.Equal(t, nil, err)

	output, err := spell.Cast()
	assert.Equal(t, nil, err)
	if "windows" == hostInfo.OS {
		assert.Equal(t, "{} }{\r\n", string(output))
	} else {
		assert.Equal(t, "{} }{\n", string(output))
	}
}
