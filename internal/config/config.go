package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	cfg_file_path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(cfg_file_path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	cfg := Config{}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	config_dir := filepath.Join(homedir, configFileName)

	return config_dir, nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName

	return write(*cfg)
}

func write(cfg Config) error {
	cfg_file_path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(cfg_file_path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)

	if err != nil {
		return err
	}

	return nil

}
