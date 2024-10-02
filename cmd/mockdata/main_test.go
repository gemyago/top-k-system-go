package main

import (
	"os"
	"path"
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
	t.Run("generate-item-ids", func(t *testing.T) {
		t.Run("should generate item ids", func(t *testing.T) {
			rootCmd := setupCommands()
			randomFileName := faker.UUIDHyphenated()
			require.NoError(t, os.Chdir(path.Join("..", "..")))
			rootCmd.SetArgs([]string{"generate-item-ids", "-n", "1", "-o", randomFileName, "--logs-file", "../../test.log"})
			require.NoError(t, rootCmd.Execute())
			// check if file exists
			filePath := path.Join("tmp", "blobs", randomFileName)
			_, err := os.Stat(filePath)
			require.NoError(t, err)
			require.NoError(t, os.Remove(filePath))
		})
	})
}
