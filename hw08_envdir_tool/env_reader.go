package main

import (
	"bytes"
	"errors"
	"os"
	"strings"
)

type Environment map[string]EnvValue

var (
	ErrInvalidEnvFile = errors.New("env file name can not contain '='")
	ErrInvalidEnvDir  = errors.New("can't read env direcrory")
)

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, ErrInvalidEnvDir
	}
	env := make(Environment, len(entries))
	errs := []error{}

	for i := range entries {
		if !entries[i].Type().IsRegular() {
			continue
		}
		if strings.Contains(entries[i].Name(), "=") {
			errs = append(errs, ErrInvalidEnvFile)
			continue
		}
		envfile, err := os.ReadFile(dir + "/" + entries[i].Name())
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if len(envfile) == 0 {
			env[entries[i].Name()] = EnvValue{NeedRemove: true}
			continue
		}
		if pos := bytes.IndexByte(envfile, 0x0a); pos > -1 {
			envfile = envfile[:pos]
		}
		envfile = bytes.ReplaceAll(envfile, []byte{0x00}, []byte{0x0a})
		envfile = bytes.TrimRight(envfile, "\t ")
		env[entries[i].Name()] = EnvValue{Value: string(envfile)}
	}

	return env, errors.Join(errs...)
}
