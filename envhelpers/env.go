package envhelpers

import "strings"

// SetEnvVars sets vars on the form key=value in the oldVars slice.
func SetEnvVars(oldVars *[]string, keyValues ...string) {
	for i := 0; i < len(keyValues); i += 2 {
		setEnvVar(oldVars, keyValues[i], keyValues[i+1])
	}
}

// SplitEnvVar splits an env var into key and value.
func SplitEnvVar(v string) (string, string) {
	name, value, _ := strings.Cut(v, "=")
	return name, value
}

func setEnvVar(vars *[]string, key, value string) {
	for i := range *vars {
		if strings.HasPrefix((*vars)[i], key+"=") {
			(*vars)[i] = key + "=" + value
			return
		}
	}
	// New var.
	*vars = append(*vars, key+"="+value)
}
