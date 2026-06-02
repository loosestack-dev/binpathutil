package usecases

import (
	"fmt"

	"github.com/loosestack-dev/binpathutil/internal/path/domain"
)

func appendElement(element string, ignoreIfPresent bool, getPathDependency func() (string, error)) (string, error) {
	rawEnvValue, err := getPathDependency()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve PATH env variable: %w", err)
	}

	env := domain.NewEnvPath(rawEnvValue)
	if env.Contains(element) && ignoreIfPresent {
		return env.String(), nil
	} else {
		env.Append(element)
	}

	return env.String(), nil
}

func Append(element string, getPathDependency func() (string, error)) (string, error) {
	return appendElement(element, false, getPathDependency)
}

func AppendIfAbsent(element string, getPathDependency func() (string, error)) (string, error) {
	return appendElement(element, true, getPathDependency)
}
