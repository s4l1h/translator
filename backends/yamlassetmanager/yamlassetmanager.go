s4l1hpackage yamlassetmanager

import (
	"path/filepath"
	"strings"

	"github.com/akmyazilim/assetmanager"
	"github.com/akmyazilim/translator"
	yaml "gopkg.in/yaml.v2"
)

// YamlAssetManager struct
type YamlAssetManager struct {
	manager      *assetmanager.AssetManager
	translations []*translator.Translation
}

// Add yaml file
func (backend *YamlAssetManager) Add(locale string, value interface{}, scopes []string) {
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
func (backend *YamlAssetManager) Load() (translations []*translator.Translation) {
	//spew.Dump(backend.files)
	var err error
	for file, content := range backend.manager.GetAll() {
		fName := filepath.Base(file)
		extName := filepath.Ext(file)
		locale := fName[:len(fName)-len(extName)]
		var slice yaml.MapSlice
		if err = yaml.Unmarshal(content, &slice); err == nil {
			backend.Add(locale, slice, []string{})
		} else {
			panic(err)
		}
	}
	return backend.translations
}

// Delete translation
func (backend *YamlAssetManager) Delete(trans *translator.Translation) {

}

// Save translation
func (backend *YamlAssetManager) Save(trans *translator.Translation) {

}

// New Yaml Asset Manager Backend struct
func New(manager *assetmanager.AssetManager) *YamlAssetManager {
	return &YamlAssetManager{
		manager: manager,
	}
}
