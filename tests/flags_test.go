package tests

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		wantOutput    string
		outputPattern *regexp.Regexp
		wantExitCode  int
	}{
		{
			name:         "succeed running main with one arg",
			args:         []string{"-V"},
			wantOutput:   "dyndns version (devel) (0001-01-01)\n",
			wantExitCode: 0,
		},
		{
			name:         "fail running main with two args",
			args:         []string{"--version"},
			wantOutput:   "dyndns version (devel) (0001-01-01)\n",
			wantExitCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, exitCode, err := collector.RunBinary(binPath, "TestBincoverRunMain", []string{}, tt.args)
			require.NoError(t, err)

			if tt.outputPattern != nil {
				require.Regexp(t, tt.outputPattern, output)
			} else {
				require.Equal(t, tt.wantOutput, output)
			}
			require.Equal(t, tt.wantExitCode, exitCode)
		})
	}
}

func TestRequiredFlags(t *testing.T) {
	env := []string{
		"GANDI_TOKEN=xxx",
		"DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/xxx",
	}

	tests := []struct {
		name          string
		args          string
		env           []string
		wantOutput    string
		outputPattern *regexp.Regexp
		wantExitCode  int
	}{
		{
			name:         "missing GANDI_TOKEN env var",
			args:         "--domain example.com --record yolo",
			env:          []string{"DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/xxx"},
			wantOutput:   "error: required environment variable GANDI_TOKEN is empty or missing\n",
			wantExitCode: 1,
		},
		{
			name:         "missing DISCORD_WEBHOOK_URL env var",
			args:         "--domain example.com --record yolo",
			env:          []string{"GANDI_TOKEN=xxx"},
			wantOutput:   "error: required environment variable DISCORD_WEBHOOK_URL is empty or missing\n",
			wantExitCode: 1,
		},
		{
			name:         "missing --record flag",
			args:         "--domain example.com",
			env:          env,
			wantOutput:   "error: required flag --record is missing\n",
			wantExitCode: 1,
		},
		{
			name:         "missing --record flag",
			args:         "--record www",
			env:          env,
			wantOutput:   "error: required flag --domain is missing\n",
			wantExitCode: 1,
		},
		{
			name:          "no args",
			args:          "",
			env:           env,
			outputPattern: regexp.MustCompile("Usage:.*"),
			wantExitCode:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, exitCode, err := collector.RunBinary(binPath, "TestBincoverRunMain", tt.env, strings.Split(tt.args, " "))
			require.NoError(t, err)

			if tt.outputPattern != nil {
				require.Regexp(t, tt.outputPattern, output)
			} else {
				require.Equal(t, tt.wantOutput, output)
			}
			require.Equal(t, tt.wantExitCode, exitCode)
		})
	}
}
