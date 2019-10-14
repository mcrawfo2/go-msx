package stats

import "strings"

func Name(prefix, api, param string) string {
	pathParts := []string{
		prefix,
		api,
	}

	if param != "" {
		pathParts = append(pathParts, param)
	}

	return strings.Join(pathParts, ".")
}
