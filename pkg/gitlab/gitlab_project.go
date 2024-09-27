package gitlab

import (
	"errors"

	"github.com/eadydb/nebulae/pkg/walk"
)

var LanguageMap = map[string]string{
	"pom.xml":          "Java",
	"main.go":          "Go",
	"requirements.txt": "Python",
	"package.json":     "JavaScript",
	"build.gradle":     "Java",
}

func GetProjectLanguage(dir string) (string, error) {
	for key, language := range LanguageMap {
		if pomPaths, err := walk.From(dir).CollectFilterPaths(key); err == nil {
			if len(pomPaths) > 0 {
				return language, nil
			}
		}
	}
	return "", errors.New("not found language")
}
