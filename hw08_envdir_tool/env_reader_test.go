package main

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	envGolden := Environment{
		"BAR":   {Value: "bar"},
		"EMPTY": {},
		"FOO": {Value: `   foo
with new line`},
		"HELLO": {Value: `"hello"`},
		"UNSET": {NeedRemove: true},
	}

	t.Run("success", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		require.NoError(t, err)
		require.Equal(t, envGolden, env)
	})

	t.Run("skip irregular files", func(t *testing.T) {
		dir, err := os.MkdirTemp("testdata/env", "testdir-*")
		require.NoError(t, err)
		defer os.RemoveAll(dir) // clean up

		env, err := ReadDir("testdata/env")
		require.NoError(t, err)
		require.Equal(t, envGolden, env)
	})

	t.Run("fail read env dir", func(t *testing.T) {
		_, err := ReadDir("testdata/env_fail")
		require.ErrorIs(t, err, ErrInvalidEnvDir)
	})

	t.Run("incorrect filename", func(t *testing.T) {
		f, err := os.CreateTemp("testdata/env", "test=*")
		require.NoError(t, err)
		err = f.Close()
		require.NoError(t, err)
		defer os.Remove(f.Name())

		env, err := ReadDir("testdata/env")
		require.Error(t, err)
		require.Equal(t, envGolden, env)
	})

	t.Run("unable read file", func(t *testing.T) {
		f, err := os.OpenFile("testdata/env/test_127890", os.O_CREATE|os.O_WRONLY, 0o222)
		require.NoError(t, err)
		err = f.Close()
		require.NoError(t, err)
		defer os.Remove(f.Name())

		env, err := ReadDir("testdata/env")
		var pathError *fs.PathError
		require.ErrorAs(t, err, &pathError)
		require.Equal(t, envGolden, env)
	})
}
