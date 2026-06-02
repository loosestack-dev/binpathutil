package usecases

import (
	"fmt"
	"regexp"

	"github.com/loosestack-dev/binpathutil/internal/path/domain"
)

func Contains(element string, getPathDependency func() (string, error)) (bool, error) {
	rawEnvValue, err := getPathDependency()
	if err != nil {
		return false, fmt.Errorf("unable to retrieve PATH env variable: %w", err)
	}

	env := domain.NewEnvPath(rawEnvValue)
	return env.Contains(element), nil
}

func ContainsRegex(expression string, getPathDependency func() (string, error)) (bool, error) {
	rawEnvValue, err := getPathDependency()
	if err != nil {
		return false, fmt.Errorf("unable to retrieve PATH env variable: %w", err)
	}

	re, err := regexp.Compile(expression)
	if err != nil {
		return false, fmt.Errorf("invalid regular expression %q: %w", expression, err)
	}

	env := domain.NewEnvPath(rawEnvValue)
	return env.ContainsMatch(re), nil
}
