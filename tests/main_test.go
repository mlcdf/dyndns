//go:build !testbincover
// +build !testbincover

package tests

import (
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/confluentinc/bincover"
	"go.mlcdf.fr/dyndns/tests/smockertest"
)

const binPath = "../dist/dyndns.test"

var collector *bincover.CoverageCollector

func TestMain(m *testing.M) {
	smocker := smockertest.MustStart()
	buildTestBinary()

	collector = bincover.NewCoverageCollector("../dist/coverage.out", true)

	err := collector.Setup()
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	code := m.Run()

	err = collector.TearDown()
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	err = os.Remove(binPath)
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	smocker.Nuke()

	os.Exit(code)
}

func buildTestBinary() {
	buildTestCmd := exec.Command("../scripts/build-test-binary.sh")
	output, err := buildTestCmd.CombinedOutput()
	if err != nil {
		log.Fatalf("err: %s, output: %s", err, output)
	}
}
