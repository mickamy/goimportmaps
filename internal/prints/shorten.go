package prints

import (
	"strings"
)

func Shorten(pkgPath, modulePath string) string {
	if strings.HasPrefix(pkgPath, modulePath) {
		rel := strings.TrimPrefix(pkgPath, modulePath)
		rel = strings.TrimPrefix(rel, "/")
		return rel
	}
	return pkgPath
}
