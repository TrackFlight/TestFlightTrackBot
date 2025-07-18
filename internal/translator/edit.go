package translator

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type Action int

var editLangMutex = &sync.Mutex{}

const (
	Add Action = iota
	Remove
	Update
)

func Edit(key, value string, action Action) (bool, error) {
	editLangMutex.Lock()
	defer editLangMutex.Unlock()

	exports, err := loadLocaleFile[string](DefaultLanguage)
	if err != nil {
		return false, err
	}

	if action == Remove {
		if _, exists := exports[key]; !exists {
			return false, nil
		}
		delete(exports, key)
	} else {
		if _, exists := exports[key]; exists && action == Add {
			return false, nil
		} else if action == Update && !exists {
			return false, nil
		}
		exports[key] = value
	}

	if action == Update || action == Remove {
		err = filepath.Walk(LocalePath, func(dir string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			extension := filepath.Ext(info.Name())
			lang := info.Name()[0 : len(info.Name())-len(extension)]
			if info.IsDir() || extension != ".json" {
				return nil
			}
			if lang == DefaultLanguage {
				return nil
			}
			tmpData, err := loadLocaleFile[string](lang)
			if err != nil {
				return err
			}
			if _, exists := tmpData[key]; exists {
				delete(tmpData, key)
			}
			return writeLocaleFile(lang, tmpData)
		})
		if err != nil {
			return false, fmt.Errorf("failed to walk through locale files: %w", err)
		}
	}

	if err = writeLocaleFile(DefaultLanguage, exports); err != nil {
		return false, err
	}
	return true, nil
}

func SearchVar(content string) (map[string]string, error) {
	editLangMutex.Lock()
	defer editLangMutex.Unlock()
	exports, err := loadLocaleFile[string](DefaultLanguage)
	if err != nil {
		return nil, err
	}
	content = strings.ToUpper(content)
	result := make(map[string]string)
	for key, value := range exports {
		keyU := strings.ToUpper(key)
		valueU := strings.ToUpper(value)
		if strings.Contains(keyU, content) || strings.Contains(valueU, content) {
			result[key] = value
		}
	}
	return result, nil
}

func writeLocaleFile(locale string, data map[string]string) error {
	filePath := path.Join(LocalePath, fmt.Sprintf("%s.json", locale))
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	if err = os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}
	return nil
}
