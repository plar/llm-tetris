package main

import (
	"encoding/json"
	"os"
)

type KeyBindings struct {
	Left     string
	Right    string
	Down     string
	Rotate   string
	HardDrop string
	Pause    string
}

type Config struct {
	KeyBindings    KeyBindings
	SimpleRotation bool
	HighScore      int
}

func NewConfig() *Config {
	return &Config{
		KeyBindings: KeyBindings{
			Left:     "Left",
			Right:    "Right",
			Down:     "Down",
			Rotate:   "Up",
			HardDrop: " ",
			Pause:    "P",
		},
		SimpleRotation: true, // Default is simple rotation
		HighScore:      0,
	}
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := NewConfig()
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) SaveConfig(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}
