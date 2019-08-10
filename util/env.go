package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// RegisterEnvFromFile - register env data from file
func RegisterEnvFromFile(envFile string) error {
	envMap, err := ParseEnvFromFile(envFile)
	if err != nil {
		return err
	}

	for key, value := range envMap {
		os.Setenv(key, value)
	}
	return nil
}

// ParseEnvFromFile - parse env file and returns a map explaining
// all the details
func ParseEnvFromFile(envFile string) (map[string]string, error) {
	var envMap = map[string]string{}

	envs, err := ioutil.ReadFile(envFile)
	if err != nil {
		return envMap, fmt.Errorf("read env file error: read file(%s) failed", envFile)
	}
	re, _ := regexp.Compile("^([A-Z_]+)=(.*)$")

	var strEnvs = strings.Split(string(envs), "\n")
	for _, strEnv := range strEnvs {
		kvPair := re.FindStringSubmatch(strEnv)
		envMap[kvPair[1]] = kvPair[2]
	}

	return envMap, nil
}
