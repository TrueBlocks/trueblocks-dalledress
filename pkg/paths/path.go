package paths

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

// GetConfigDir returns the operating system's configuration folder for the current user.
// If the folder does not exist, it is created.
func GetConfigDir() (string, error) {
	if configPath, err := os.UserConfigDir(); err != nil {
		return "", err
	} else {
		path := filepath.Join(configPath, "TrueBlocks/browse")
		file.EstablishFolder(path)
		return path, nil
	}
}

// GetCacheDir returns the operating system's cache folder for the current user.
// If the folder does not exist, it is created.
func GetCacheDir() (string, error) {
	if cachePath, err := os.UserCacheDir(); err != nil {
		return "", err
	} else {
		path := filepath.Join(cachePath, "TrueBlocks/browse")
		file.EstablishFolder(path)
		return path, nil
	}
}

var ErrDocFolderNotFound = errors.New("no documents folder found")

// GetDocumentsDir returns the operating system's documents folder for the current user.
func GetDocumentsDir() (string, error) {
	if userPath, err := os.UserHomeDir(); err != nil {
		return "", err
	} else {
		path := filepath.Join(userPath, "Documents")
		if !file.FolderExists(path) {
			return "", ErrDocFolderNotFound
		}
		return path, nil
	}
}
