package civov2

import "strings"

func mergeHostPath(host, path string) string {
	return strings.TrimSuffix(host, "/") + "/" + strings.TrimPrefix(path, "/")
}
