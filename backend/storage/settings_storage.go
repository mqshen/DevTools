package storage

import (
	"devtools/backend/types"
	"gopkg.in/yaml.v3"
	"sync"
)

type SettingsStorage struct {
	storage *localStorage
	mutex   sync.Mutex
}

func NewSettings() *SettingsStorage {
	return &SettingsStorage{
		storage: NewLocalStore("settings.yaml"),
	}
}

func (s *SettingsStorage) defaultSettings() types.Settings {
	return types.Settings{
		IPInfoToken: "",
	}
}

func (s *SettingsStorage) GetIPInfoToken() string {
	settings := s.getSettings()
	return settings.IPInfoToken
}

func (s *SettingsStorage) getSettings() (ret types.Settings) {
	b, err := s.storage.Load()
	ret = s.defaultSettings()
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(b, &ret); err != nil {
		ret = s.defaultSettings()
		return
	}
	return
}
