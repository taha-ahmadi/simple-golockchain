package golockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type folderConfig struct {
	LastHash []byte
}

type folderStore struct {
	root string

	config *folderConfig

	configPath string
}

func (fs *folderStore) Load(hash []byte) (*Block, error) {
	path := filepath.Join(fs.root, fmt.Sprintf("%x.json", hash))
	var b Block
	if err := readJSON(path, &b); err != nil {
		return nil, fmt.Errorf("read JOSN file failed: %w", err)
	}

	return &b, nil
}

func (fs *folderStore) Append(block *Block) (*Block, error) {
	if !bytes.Equal(fs.config.LastHash, block.PrevHash) {
		return nil, fmt.Errorf("store is out of sync")
	}

	path := filepath.Join(fs.root, fmt.Sprintf("%x.json", block.Hash))
	if err := writeJSON(path, block); err != nil {
		return nil, fmt.Errorf("write JSON file failed: %w", err)
	}

	fs.config.LastHash = block.Hash
	if err := writeJSON(fs.configPath, fs.config); err != nil {
		return nil, fmt.Errorf("write configuration file failed: %w", err)
	}

	return block, nil
}

func (fs *folderStore) LastHash() ([]byte, error) {
	if len(fs.config.LastHash) == 0 {
		return nil, ErrNoInitialized
	}

	return fs.config.LastHash, nil
}

func readJSON(path string, v interface{}) error {
	fl, err := os.Open(path)
	defer fl.Close()
	if err != nil {
		return fmt.Errorf("failed to open the file %s: %v", path, err)
	}
	decoder := json.NewDecoder(fl)

	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("failed to decode json: %v", err)
	}

	return nil
}

func writeJSON(path string, v interface{}) error {
	fl, err := os.Create(path)
	defer fl.Close()
	if err != nil {
		return fmt.Errorf("failed to create the file %s: %v", path, err)
	}
	encoder := json.NewEncoder(fl)
	encoder.SetIndent("", " ")

	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("failed to encode json: %v", err)
	}

	return nil
}

// NewFolderStore is persistent storage that stores blocks in the
// JSON file, each block is in the one file and there is a config file
// that keeps the last hash of the blockchain
func NewFolderStore(root string) Store {
	fs := &folderStore{
		root:       root,
		config:     &folderConfig{},
		configPath: filepath.Join(root+"/config", "config.json"),
	}

	if err := readJSON(fs.configPath, fs.config); err != nil {
		log.Print("failed to read config")
		fs.config.LastHash = nil
	}

	return fs
}
