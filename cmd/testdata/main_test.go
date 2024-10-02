package main

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	t.Run("send-test-events", func(t *testing.T) {
		t.Run("should send events in noop mode", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SetArgs([]string{"send-test-events", "--noop", "--logs-file", "../../test.log"})
			require.NoError(t, rootCmd.Execute())
		})
		t.Run("should fail if bad log level", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SilenceErrors = true
			rootCmd.SilenceUsage = true
			rootCmd.SetArgs([]string{"send-test-events", "--noop", "-l", faker.Word(), "--logs-file", "../../test.log"})
			assert.Error(t, rootCmd.Execute())
		})
		t.Run("should fail if bad env level", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SilenceErrors = true
			rootCmd.SilenceUsage = true
			rootCmd.SetArgs([]string{"send-test-events", "--noop", "--env", faker.Word(), "--logs-file", "../../test.log"})
			assert.Error(t, rootCmd.Execute())
		})
	})
}
