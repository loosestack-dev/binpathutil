package path

import (
	"fmt"
	"os"
)

func getEnv(env string) (string, error) {
	value, ok := os.LookupEnv(env)
	if !ok {
		return "", fmt.Errorf("environment variable %q is not set", env)
	}
	return value, nil
}

func GetEnvPath() (string, error) {
	value, err := getEnv("PATH")
	if err != nil {
		return "", err
	}

	return value, nil
}
