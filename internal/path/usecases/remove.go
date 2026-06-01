package usecases

import (
	"fmt"
	"regexp"

	"binpathutil/internal/path/domain"
)

func RemoveIfPresent(element string, getPathDependency func() (string, error)) (string, error) {
	rawEnvValue, err := getPathDependency()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve PATH env variable: %w", err)
	}

	env := domain.NewEnvPath(rawEnvValue)
	_ = env.Remove(element)
	return env.String(), nil
}

func Remove(element string, getPathDependency func() (string, error)) (string, error) {
	rawEnvValue, err := getPathDependency()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve PATH env variable: %w", err)
	}

	env := domain.NewEnvPath(rawEnvValue)
	errRemove := env.Remove(element)
	if errRemove != nil {
		return "", fmt.Errorf("unable to remove %s from PATH: %w", element, errRemove)
	} else {
		return env.String(), nil
	}
}

func RemoveAllOccurence(element string, getPathDependency func() (string, error)) (string, error) {
	rawEnvValue, err := getPathDependency()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve PATH env variable: %w", err)
	}

	env := domain.NewEnvPath(rawEnvValue)
	errRemove := env.Remove(element)

	for errRemove == nil {
		errRemove = env.Remove(element)
	}

	return env.String(), nil
}

func RemoveRegex(expression string, getPathDependency func() (string, error)) (string, error) {
	rawEnvValue, err := getPathDependency()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve PATH env variable: %w", err)
	}

	re, err := regexp.Compile(expression)
	if err != nil {
		return "", fmt.Errorf("invalid regular expression %q: %w", expression, err)
	}

	env := domain.NewEnvPath(rawEnvValue)
	err = env.RemoveMatch(re)
	if err != nil {
		return "", fmt.Errorf("unable to remove entries matching %q from PATH: %w", expression, err)
	}
	return env.String(), nil
}

func RemoveIfPresentRegex(expression string, getPathDependency func() (string, error)) (string, error) {
	rawEnvValue, err := getPathDependency()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve PATH env variable: %w", err)
	}

	re, err := regexp.Compile(expression)
	if err != nil {
		return "", fmt.Errorf("invalid regular expression %q: %w", expression, err)
	}

	env := domain.NewEnvPath(rawEnvValue)
	_ = env.RemoveMatch(re)
	return env.String(), nil
}

func RemoveAllOccurenceRegex(expression string, getPathDependency func() (string, error)) (string, error) {
	rawEnvValue, err := getPathDependency()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve PATH env variable: %w", err)
	}

	re, err := regexp.Compile(expression)
	if err != nil {
		return "", fmt.Errorf("invalid regular expression %q: %w", expression, err)
	}

	env := domain.NewEnvPath(rawEnvValue)
	errRemove := env.RemoveMatch(re)

	for errRemove == nil {
		errRemove = env.RemoveMatch(re)
	}

	return env.String(), nil
}
