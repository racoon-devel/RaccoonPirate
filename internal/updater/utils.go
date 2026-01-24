package updater

import (
	"errors"
	"strings"

	"github.com/blang/semver"
)

func ParseVersion(version string) (semver.Version, error) {
	if !strings.HasPrefix(version, "v") {
		return semver.Version{}, errors.New("version must begin from v")
	}
	return semver.Parse(version)
}
