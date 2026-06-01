package domain

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

type EnvPath struct {
	Entries []string
}

func NewEnvPath(rawPath string) *EnvPath {
	if rawPath == "" {
		return &EnvPath{Entries: []string{}}
	}
	return &EnvPath{
		Entries: strings.Split(rawPath, string(os.PathListSeparator)),
	}
}

func (p *EnvPath) Contains(entry string) bool {
	return slices.Contains(p.Entries, entry)
}

func (p *EnvPath) ContainsMatch(re *regexp.Regexp) bool {
	return slices.ContainsFunc(p.Entries, re.MatchString)
}

func (p *EnvPath) Remove(entry string) error {
	i := slices.Index(p.Entries, entry)
	if i == -1 {
		return fmt.Errorf("entry %q not found in PATH", entry)
	}
	p.Entries = slices.Delete(p.Entries, i, i+1)
	return nil
}

// RemoveMatch removes the first entry matching the expression; errors if none match.
func (p *EnvPath) RemoveMatch(re *regexp.Regexp) error {
	i := slices.IndexFunc(p.Entries, re.MatchString)
	if i == -1 {
		return fmt.Errorf("no entry matching %q found in PATH", re.String())
	}
	p.Entries = slices.Delete(p.Entries, i, i+1)
	return nil
}

func (p *EnvPath) Prepend(entry string) {
	p.Entries = slices.Insert(p.Entries, 0, entry)
}

func (p *EnvPath) Append(entry string) {
	p.Entries = append(p.Entries, entry)
}

func (p *EnvPath) String() string {
	return strings.Join(p.Entries, string(os.PathListSeparator))
}
