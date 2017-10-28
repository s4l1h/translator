package translator

import (
	"strings"

	"github.com/theplant/cldr"
)

// Translation  object
type Translation struct {
	Locale string
	Key    string
	Value  string
}

// New return Translator object
func New(Backends ...Backend) *Translator {
	translator := &Translator{Data: make(map[string]map[string]string)}
	for i := len(Backends); i > 0; i-- {
		backend := Backends[i-1]
		translator.AddBackend(backend)
	}

	return translator
}

// Translator object
type Translator struct {
	Data        map[string]map[string]string
	DefaultLang string
}

// AddBackend to translator object
func (t *Translator) AddBackend(backend Backend) {
	for _, translation := range backend.Load() {
		t.Add(translation)
	}
}

// Add translation to translator object
func (t *Translator) Add(trans *Translation) {
	if _, ok := t.Data[trans.Locale]; !ok {
		t.Data[trans.Locale] = make(map[string]string)
	}
	t.Data[trans.Locale][trans.Key] = trans.Value
}

// Get translation
func (t *Translator) Get(trans Translation) string {
	if val, ok := t.Data[trans.Locale]; ok {
		if res, okk := val[trans.Key]; okk {
			return res
		}
	}
	return ""
}

// GetKey translation with locale and key
func (t *Translator) GetKey(locale string, key string) string {
	return t.Get(Translation{Locale: locale, Key: key})
}

// GetKeys translation with locale and search pattern
func (t *Translator) GetKeys(locale string, pattern string, disallow string) map[string]string {
	result := make(map[string]string)
	if val, ok := t.Data[locale]; ok {
		for k, v := range val {
			if strings.Contains(k, pattern) && !strings.Contains(k, disallow) {
				// logrus.Warn(pattern, k)
				result[k] = v
			}
		}
	}
	return result
}

// T return translation
func (t *Translator) T(locale, key string, args ...interface{}) string {

	// logrus.Warnf("locale %s key %s args %v", locale, key, args)

	// spew.Dump(args)
	value := t.GetKey(locale, key)
	if value == "" {
		return key // return translate key because no val
		// logrus.Warnf("locale %s key %s args %v", locale, key, args)
		// logrus.Warnf("value %s", value)
	}

	str, err := cldr.Parse(locale, value, args...)
	if err == nil {
		value = str
	}

	return str
}
