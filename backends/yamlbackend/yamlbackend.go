package yamlbackend

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/akmyazilim/translator"
	yaml "gopkg.in/yaml.v2"
)

// YamlBackend struct
type YamlBackend struct {
	files        []string
	translations []*translator.Translation
}

// Add yaml file
func (backend *YamlBackend) Add(locale string, value interface{}, scopes []string) {
	switch v := value.(type) {
	case yaml.MapSlice:
		for _, s := range v {
			backend.Add(locale, s.Value, append(scopes, s.Key.(string)))
		}
	default:
		//spew.Dump(v)
		var translation = &translator.Translation{
			Locale: locale,
			Key:    strings.Join(scopes, "."),
			Value:  v.(string),
		}
		backend.translations = append(backend.translations, translation)
	}
	return
}

// Load return translations
func (backend *YamlBackend) Load() (translations []*translator.Translation) {
	//spew.Dump(backend.files)
	for _, file := range backend.files {
		fName := filepath.Base(file)
		extName := filepath.Ext(file)
		locale := fName[:len(fName)-len(extName)]
		if content, err := ioutil.ReadFile(file); err == nil {
			var slice yaml.MapSlice
			if err = yaml.Unmarshal(content, &slice); err == nil {
				backend.Add(locale, slice, []string{})
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return backend.translations
}

// Delete translation
func (backend *YamlBackend) Delete(trans *translator.Translation) {

}

// Save translation
func (backend *YamlBackend) Save(trans *translator.Translation) {

}

// New Yaml Backend struct
func New(paths ...string) *YamlBackend {
	backend := &YamlBackend{}
	for _, p := range paths {
		if file, err := os.Open(p); err == nil {
			defer file.Close()
			if fileInfo, err := file.Stat(); err == nil {
				if fileInfo.IsDir() {
					yamlFiles, _ := filepath.Glob(path.Join(p, "*.yaml"))
					backend.files = append(backend.files, yamlFiles...)
					ymlFiles, _ := filepath.Glob(path.Join(p, "*.yml"))
					backend.files = append(backend.files, ymlFiles...)
				} else if fileInfo.Mode().IsRegular() {
					backend.files = append(backend.files, p)
				}
			}
		}
	}
	return backend
}
