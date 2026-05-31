package usecases

import (
	"fmt"

	"binpathutil/internal/path/domain"
)

func Contains(element string, getPathDependency func() (string, error)) (bool, error) {
	rawEnvValue, err := getPathDependency()
	if err != nil {
		return false, fmt.Errorf("unable to retrieve PATH env variable: %w", err)
	}

	env := domain.NewEnvPath(rawEnvValue)
	return env.Contains(element), nil
}
