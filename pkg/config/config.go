/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Credentials struct {
	Profile string `yaml:"profile"`
	Host    string `yaml:"host"`
	User    string `yaml:"user"`
	Port    string `yaml:"port"`
	KeyPath string `yaml:"key_path"`
}

type Config struct {
	Profiles map[string]Credentials `yaml:"profiles"`
}

func InitDefaultProfile() (*Credentials, error) {
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDirPath, ".fnd.yaml")

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	config := &Config{}
	decoder := yaml.NewDecoder(f)

	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	defaultCreds := config.Profiles["Default"]

	return &defaultCreds, nil

}

func CreateConfigFile() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	_, err = os.Stat(filepath.Join(homeDir, ".fnd.yaml"))

	if os.IsNotExist(err) {
		_, e := os.Create(filepath.Join(homeDir, ".fnd.yaml"))
		if e != nil {
			return e
		}

		// slog.Info("Config file created")
	}

	return err
}

func GetSshCredentials() (*Credentials, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Profile name [********me]: ")
	profile, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	profile = strings.TrimSpace(profile)

	if len(profile) == 0 {
		profile = "Default"
	}

	var host string
	fmt.Print("Host [********me]: ")
	host, err = reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	host = strings.TrimSpace(host)

	var user string
	fmt.Print("User [********me]: ")
	user, err = reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	user = strings.TrimSpace(user)

	var port string
	fmt.Print("Port [********me]: ")
	port, err = reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	port = strings.TrimSpace(port)

	if len(port) == 0 {
		port = "22"
	}

	var keyPath string
	fmt.Print("Key Path [********me]: ")
	keyPath, err = reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	keyPath = strings.TrimSpace(keyPath)

	creds := Credentials{
		Profile: profile,
		User:    user,
		Host:    host,
		Port:    port,
		KeyPath: keyPath,
	}
	return &creds, nil
}

func UpdateConfig(credentials *Credentials) error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homedir, ".fnd.yaml")

	f, e := os.Open(configPath)
	if e != nil {
		return e
	}
	defer f.Close()

	config := &Config{}
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(config); err != nil && err.Error() != "EOF" {
		return err
	}

	if config.Profiles == nil {
		config.Profiles = make(map[string]Credentials)
	}

	config.Profiles[credentials.Profile] = *credentials

	outFile, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	encoder := yaml.NewEncoder(outFile)
	if err = encoder.Encode(config); err != nil {
		return err
	}

	fmt.Println("Config Successfully Updated")

	return nil
}
