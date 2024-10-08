package gitlab

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eadydb/nebulae/pkg/consts"
)

// scanDirectory scans a directory for files with specific names or extensions
func scanDirectory(dir string, targets []string) (map[string]bool, error) {
	result := make(map[string]bool)
	for _, target := range targets {
		result[target] = false
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		fileName := info.Name()
		fileExt := strings.ToLower(filepath.Ext(fileName))

		for target := range result {
			if strings.HasPrefix(target, ".") {
				// Check for file extension
				if fileExt == strings.ToLower(target) {
					result[target] = true
				}
			} else {
				// Check for exact file name
				if strings.ToLower(fileName) == strings.ToLower(target) {
					result[target] = true
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error scanning directory: %w", err)
	}

	return result, nil
}

func (g *Gitlab) GetProjectLanguage(dir string) (string, error) {
	languages := make([]string, 0)
	for language, targets := range consts.CommonLanguages {
		result, err := scanDirectory(dir, targets)
		if err != nil {
			return "", err
		}
		for _, found := range result {
			if found {
				languages = append(languages, language)
				break
			}
		}
	}
	if len(languages) > 0 {
		return strings.Join(languages, ", "), nil
	}

	return "", errors.New("not found language")
}

func (g *Gitlab) UpdateRepositoryLanguage(id int, dir string) error {
	language, err := g.GetProjectLanguage(dir)
	if err != nil {
		return err
	}
	return newRepositoryService(g.Ctx).UpdateRepository(language, dir, id)
}
