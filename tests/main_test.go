//go:build !testbincover
// +build !testbincover

package tests

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/confluentinc/bincover"
)

const binPath = "../dist/dyndns.test"

var collector *bincover.CoverageCollector

func TestMain(m *testing.M) {
	flag.Parse()

	dockerComposeUp()
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

	if dev := os.Getenv("DEV"); dev == "" {
		dockerComposeDown()
	}

	os.Exit(code)
}

func dockerComposeUp() {
	cmd := exec.Command("docker", "compose", "up", "-d")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("err: %s, output: %s", err, output)
	}
}

func dockerComposeDown() {
	cmd := exec.Command("docker", "compose", "down")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("err: %s, output: %s", err, output)
	}
}

func buildTestBinary() {
	buildTestCmd := exec.Command("../scripts/build-test-binary.sh")
	output, err := buildTestCmd.CombinedOutput()
	if err != nil {
		log.Fatalf("err: %s, output: %s", err, output)
	}
}

func pushMock(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}

	url := "http://localhost:8081/mocks?reset=true"
	res, err := http.Post(url, "content-type: application/x-yaml", f)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("error %d while performing POST %s", res.StatusCode, url)
	}

	return err
}
