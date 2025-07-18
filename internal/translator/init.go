package translator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

var langPacks = make(map[string]map[Key]string)

const LocalePath = "locales"
const DefaultLanguage = "en"

var RegexPlaceholder = regexp.MustCompile(`\{\{\.(\w+)}}`)

func init() {
	err := filepath.Walk(LocalePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			lang := extractLang(path)
			if lang == "" {
				return fmt.Errorf("invalid language file: %s", path)
			}
			localeFile, err := loadLocaleFile[Key](lang)
			if err != nil {
				return err
			}
			langPacks[lang] = localeFile
		}
		return nil
	})
	if err != nil {
		log.Fatalf("failed loading translations: %v", err)
	}
}

func loadLocaleFile[T comparable](locale string) (map[T]string, error) {
	filePath := path.Join(LocalePath, fmt.Sprintf("%s.json", locale))
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	var exports map[T]string
	err = json.Unmarshal(file, &exports)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return exports, nil
}
