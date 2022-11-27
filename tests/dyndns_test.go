package tests

import (
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mlcdf.fr/dyndns/tests/smockertest"
)

func TestDyndns(t *testing.T) {
	env := []string{
		"GANDI_TOKEN=xxx",
		"DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/xxx",
	}

	tests := []struct {
		name          string
		args          string
		mockfile      string
		wantOutput    string
		outputPattern *regexp.Regexp
		wantExitCode  int
	}{
		{
			name:         "up to date",
			args:         "--domain example.com --record www",
			mockfile:     "mocks/up-to-date.yaml",
			wantExitCode: 0,
		},
		{
			name:         "up to date always notify",
			args:         "--domain example.com --record www --always-notify",
			mockfile:     "mocks/up-to-date-always-notify.yaml",
			wantExitCode: 0,
		},
		{
			name:         "update ipv4",
			args:         "--domain example.com --record www --always-notify",
			mockfile:     "mocks/update-ipv4.yaml",
			wantExitCode: 0,
		},
		{
			name:         "update ipv6",
			args:         "--domain example.com --record www --always-notify",
			mockfile:     "mocks/update-ipv6.yaml",
			wantExitCode: 0,
		},
		{
			name:         "update both",
			args:         "--domain example.com --record www --always-notify",
			mockfile:     "mocks/update-both.yaml",
			wantExitCode: 0,
		},
		{
			name:         "update both with ttl 1337",
			args:         "--domain example.com --record www --always-notify --ttl 1337",
			mockfile:     "mocks/update-both-with-ttl.yaml",
			wantExitCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockfile != "" {
				err := smockertest.PushMock(tt.mockfile)
				require.NoError(t, err)
			}

			output, exitCode, err := collector.RunBinary(binPath, "TestBincoverRunMain", env, strings.Split(tt.args, " "))
			require.NoError(t, err)

			if tt.outputPattern != nil {
				require.Regexp(t, tt.outputPattern, output)
			} else if tt.wantOutput != "" {
				require.Equal(t, tt.wantOutput, output)
			}

			log.Println(output)
			require.Equal(t, tt.wantExitCode, exitCode)
		})
	}
}
