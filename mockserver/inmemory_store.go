package mockserver

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/suzuki-shunsuke/go-graylog"
)

type InMemoryStore struct {
	users             map[string]graylog.User                  `json:"users"`
	roles             map[string]graylog.Role                  `json:"roles"`
	inputs            map[string]graylog.Input                 `json:"inputs"`
	indexSets         map[string]graylog.IndexSet              `json:"index_sets"`
	defaultIndexSetID string                                   `json:"default_index_set_id"`
	indexSetStats     map[string]graylog.IndexSetStats         `json:"index_set_stats"`
	streams           map[string]graylog.Stream                `json:"streams"`
	streamRules       map[string]map[string]graylog.StreamRule `json:"stream_rules"`
	dataPath          string                                   `json:"-"`
}

func NewInMemoryStore(dataPath string) Store {
	return &InMemoryStore{
		roles:         map[string]graylog.Role{},
		users:         map[string]graylog.User{},
		inputs:        map[string]graylog.Input{},
		indexSets:     map[string]graylog.IndexSet{},
		indexSetStats: map[string]graylog.IndexSetStats{},
		streams:       map[string]graylog.Stream{},
		streamRules:   map[string]map[string]graylog.StreamRule{},
		dataPath:      dataPath,
	}
}

// Save writes Mock Server's data in a file for persistence.
func (store *InMemoryStore) Save() error {
	if store.dataPath == "" {
		return nil
	}
	b, err := json.Marshal(store)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(store.dataPath, b, 0600)
}

// Load reads Mock Server's data from a file.
func (store *InMemoryStore) Load() error {
	if store.dataPath == "" {
		return nil
	}
	if _, err := os.Stat(store.dataPath); err != nil {
		return nil
	}
	b, err := ioutil.ReadFile(store.dataPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, store)
}
