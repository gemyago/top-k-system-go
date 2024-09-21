package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	t.Run("http", func(t *testing.T) {
		t.Run("should initialize http app", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SetArgs([]string{"http", "--noop"})
			require.NotPanics(t, func() {
				executeRootCommand(rootCmd)
			})
		})
	})
}
