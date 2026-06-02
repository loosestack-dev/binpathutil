package usecases

import (
	"fmt"

	"github.com/loosestack-dev/binpathutil/internal/path/domain"
)

func prepend(element string, ignoreIfPresent bool, getPathDependency func() (string, error)) (string, error) {
	rawEnvValue, err := getPathDependency()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve PATH env variable: %w", err)
	}

	env := domain.NewEnvPath(rawEnvValue)
	if env.Contains(element) && ignoreIfPresent {
		return env.String(), nil
	} else {
		env.Prepend(element)
	}

	return env.String(), nil
}

func Prepend(element string, getPathDependency func() (string, error)) (string, error) {
	return prepend(element, false, getPathDependency)
}

func PrependIfAbsent(element string, getPathDependency func() (string, error)) (string, error) {
	return prepend(element, true, getPathDependency)
}
