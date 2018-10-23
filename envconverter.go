package envconv

import "os"

type envConverter struct {
	getEnv func(string) string
}

func (c *envConverter) PrependToFlags(argv []string, env map[string]string) []string {
	envFlags := c.FlagsFromEnv(env)
	return append(envFlags, argv...)
}

func (c *envConverter) FlagsFromEnv(env map[string]string) []string {
	var flags []string
	for flagName, envName := range env {
		v := c.getEnv(envName)
		if v != "" {
			flags = append(flags, flagName, v)
		}
	}
	return flags
}

var defaultEnvConverter = envConverter{getEnv: os.Getenv}

func PrependToFlags(argv []string, env map[string]string) []string {
	return defaultEnvConverter.PrependToFlags(argv, env)
}

func PrependToOSargs(env map[string]string) []string {
	return defaultEnvConverter.PrependToFlags(os.Args[1:], env)
}

func FlagsFromEnv(env map[string]string) []string {
	return defaultEnvConverter.FlagsFromEnv(env)
}
