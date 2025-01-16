package remote

import (
	"os"
	"path/filepath"
)

const tokenFileName = ".raccoon-pirate-token"

func getPossibleTokenLocations() []string {
	locations := []string{}

	curDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		locations = append(locations, filepath.Join(curDir, tokenFileName))
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		locations = append(locations, filepath.Join(homeDir, tokenFileName))
	}

	cfgDir, err := os.UserConfigDir()
	if err == nil {
		locations = append(locations, filepath.Join(cfgDir, tokenFileName))
	}

	locations = append(locations, filepath.Join(os.TempDir(), tokenFileName))
	return locations
}

func tryReadToken() (token string, ok bool) {
	locations := getPossibleTokenLocations()
	for _, loc := range locations {
		content, err := os.ReadFile(loc)
		if err == nil {
			ok = true
			token = string(content)
			return
		}
	}

	return
}

func tryWriteToken(token string) bool {
	locations := getPossibleTokenLocations()
	for _, loc := range locations {
		err := os.WriteFile(loc, []byte(token), 0600)
		if err == nil {
			return true
		}
	}
	return false
}
