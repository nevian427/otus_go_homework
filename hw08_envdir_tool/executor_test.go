package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func clearEnv(env Environment) error {
	for k := range env {
		if err := os.Unsetenv(k); err != nil {
			return err
		}
	}
	return nil
}

func checkEnv(env Environment) error {
	for k, v := range env {
		val, ok := os.LookupEnv(k)
		if ok == v.NeedRemove {
			return fmt.Errorf("variable '%s' must not exists", k)
		}
		if val != v.Value {
			return fmt.Errorf("variable '%s' has incorrect value '%s' vs '%s'", k, v.Value, val)
		}
	}
	return nil
}

func TestRunCmd(t *testing.T) {
	env, err := ReadDir("testdata/env")
	require.NoError(t, err)

	t.Run("success exec", func(t *testing.T) {
		err := clearEnv(env)
		require.NoError(t, err)

		res := RunCmd([]string{"bash", "-c", "echo"}, env)
		require.Equal(t, 0, res)
	})

	t.Run("fail exec", func(t *testing.T) {
		err := clearEnv(env)
		require.NoError(t, err)

		res := RunCmd([]string{"bash", "-c", "exit 2"}, env)
		require.Equal(t, 2, res)

		err = checkEnv(env)
		require.NoError(t, err)
	})
}
