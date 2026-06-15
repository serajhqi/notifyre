package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type keyEntry struct {
	Key   string `yaml:"key"`
	Label string `yaml:"label"`
}

type keysFile struct {
	Keys []keyEntry `yaml:"keys"`
}

type KeyStore struct {
	lookup map[string]string // key → label
}

func LoadKeys(path string) (*KeyStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read keys file: %w", err)
	}
	var kf keysFile
	if err := yaml.Unmarshal(data, &kf); err != nil {
		return nil, fmt.Errorf("parse keys file: %w", err)
	}
	if len(kf.Keys) == 0 {
		return nil, fmt.Errorf("keys file has no entries")
	}
	store := &KeyStore{lookup: make(map[string]string, len(kf.Keys))}
	for _, e := range kf.Keys {
		if e.Key == "" {
			return nil, fmt.Errorf("keys file contains entry with empty key")
		}
		store.lookup[e.Key] = e.Label
	}
	return store, nil
}

func (s *KeyStore) Valid(key string) bool {
	_, ok := s.lookup[key]
	return ok
}

func (s *KeyStore) Label(key string) string {
	return s.lookup[key]
}
