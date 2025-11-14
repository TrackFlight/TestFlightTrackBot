package translator

import (
	"fmt"
	"maps"
	"os"
	"slices"

	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
)

type Translator struct {
	language string
}

func New(language string) *Translator {
	if !IsSupported(language) {
		language = DefaultLanguage
	}
	return &Translator{
		language: language,
	}
}

func (t *Translator) Language() string {
	return t.language
}

func (t *Translator) T(key Key) string {
	return t.TWithData(key, nil)
}

func (t *Translator) TWithData(key Key, values map[string]string) string {
	tmpl, ok := langPacks[t.language][key]
	if !ok {
		tmpl = langPacks[DefaultLanguage][key]
	}
	return RegexPlaceholder.ReplaceAllStringFunc(tmpl, func(match string) string {
		parts := RegexPlaceholder.FindStringSubmatch(match)
		subKey := parts[1]
		if val, exists := values[subKey]; exists {
			return val
		}
		return parts[0]
	})
}

func (t *Translator) TWithDataCountable(key Key, values map[string]string, count int) string {
	firstCheckKey := Key(fmt.Sprintf("%s_%s", key, pluralCategory(t.language, count)))
	otherCheckKey := Key(fmt.Sprintf("%s_OTHER", key))
	if _, ok := langPacks[t.language][firstCheckKey]; ok {
		key = firstCheckKey
	} else if _, ok = langPacks[t.language][otherCheckKey]; ok {
		key = otherCheckKey
	}
	tmpl, ok := langPacks[t.language][key]
	if !ok {
		tmpl = langPacks[DefaultLanguage][key]
	}
	return RegexPlaceholder.ReplaceAllStringFunc(tmpl, func(match string) string {
		parts := RegexPlaceholder.FindStringSubmatch(match)
		subKey := parts[1]
		if val, exists := values[subKey]; exists {
			return val
		}
		return parts[0]
	})
}

func TAll(key Key) []string {
	var translations []string
	for _, lang := range SupportedLanguages() {
		translations = append(translations, New(lang).T(key))
	}
	return translations
}

func TKeys() []Key {
	return slices.Collect(maps.Keys(langPacks[DefaultLanguage]))
}

func LangPack(lang string) map[Key]string {
	if !IsSupported(lang) {
		lang = DefaultLanguage
	}
	return langPacks[lang]
}

func SupportedLanguages() []string {
	return slices.Collect(maps.Keys(langPacks))
}

func IsSupported(lang string) bool {
	if _, exists := langPacks[lang]; exists {
		return true
	}
	return false
}

func pluralCategory(langCode string, count int) string {
	tag := language.Make(langCode)
	form := plural.Cardinal.MatchPlural(tag, count, 0, 0, 0, 0)
	switch form {
	case plural.Zero:
		return "ZERO"
	case plural.One:
		return "ONE"
	case plural.Two:
		return "TWO"
	case plural.Few:
		return "FEW"
	case plural.Many:
		return "MANY"
	default:
		return "OTHER"
	}
}

func extractLang(path string) (langTag string) {
	formatStartIdx := -1
	for i := len(path) - 1; i >= 0; i-- {
		c := path[i]
		if os.IsPathSeparator(c) {
			if formatStartIdx != -1 {
				langTag = path[i+1 : formatStartIdx]
			}
			return
		}
		if path[i] == '.' {
			if formatStartIdx != -1 {
				langTag = path[i+1 : formatStartIdx]
				return
			}
			formatStartIdx = i
		}
	}
	if formatStartIdx != -1 {
		langTag = path[:formatStartIdx]
	}
	return
}
