package ipfinder

import (
	"flag"
	"testing"
)

var liveboxFlag bool

func init() {
	flag.BoolVar(&liveboxFlag, "livebox", liveboxFlag, "Run tests requiring a Livebox")
}

func TestLiveboxNoEmptyIPs(t *testing.T) {
	if liveboxFlag == false {
		t.Skip("Skipping testing requiring a Livebox. Use the --livebox flag to run those tests.")
	}

	resolvedIPs, err := Livebox()
	if err != nil {
		t.Fatal(err)
	}

	if resolvedIPs.V4 == nil && resolvedIPs.V6 == nil {
		t.Fatal("IP V4 & V6 are nil")
	}

	if resolvedIPs.V4.String() == "" {
		t.Fatal("IP V4 is empty")
	}
}
