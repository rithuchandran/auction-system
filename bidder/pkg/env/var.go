package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Vars struct {
	missing   []string // names of mandatory environment variable that are missing
	malformed []string // errors describing malformed environment variable values
}

func (vars *Vars) Mandatory(key string) string {
	val := os.Getenv(key)

	if val == "" {
		vars.missing = append(vars.missing, key)
		return ""
	}

	return val
}

func (vars *Vars) MandatoryInt(key string) int {
	valStr := vars.Mandatory(key)

	val, err := strconv.Atoi(valStr)
	if err != nil {
		vars.malformed = append(vars.malformed, fmt.Sprintf("mandatory %s (value=%q) is not a number", key, valStr))
		return 0
	}

	return val
}

func (vars Vars) Optional(key, fallback string) string {
	val := os.Getenv(key)

	if val == "" {
		return fallback
	}

	return val
}

func (vars *Vars) OptionalInt(key string, fallback int) int {
	valStr := os.Getenv(key)

	if valStr == "" {
		return fallback
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		vars.malformed = append(vars.malformed, fmt.Sprintf("optional %s (value=%q) is not a number", key, valStr))
		return fallback
	}

	return val
}

func (vars Vars) Error() error {
	if len(vars.missing) > 0 {
		return fmt.Errorf("missing mandatory configuration: %s", strings.Join(vars.missing, ", "))
	}

	if len(vars.malformed) > 0 {
		return fmt.Errorf("malformed configuration: %s", strings.Join(vars.malformed, "; "))
	}

	return nil
}
