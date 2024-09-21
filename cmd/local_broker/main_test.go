package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	t.Run("send-test-events", func(t *testing.T) {
		t.Run("should initialize deps app", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SetArgs([]string{"send-test-events", "--noop", "--logs-file", "../../test.log"})
			require.NotPanics(t, func() {
				executeRootCommand(rootCmd)
			})
		})
	})
}
