package main

import (
	"encoding/json"
	"io"
	"os"
	"runtime"
)

type Config struct {
	SyntaxHighlightingStyle string `json:"syntax_highlighting_style"`
}

func (c *Config) JSON() ([]byte, error) {
	return json.MarshalIndent(c, " ", "  ")
}

func (c *Config) Write(w io.Writer) (int, error) {
	b, err := c.JSON()
	if err != nil {
		return 0, err
	}
	return w.Write(b)
}

func init() {
	home := homeDir()

	if _, err := os.Stat(home + "/.config"); os.IsNotExist(err) {
		os.Mkdir(home+"/.config", 0755)
	}

	c := &Config{
		SyntaxHighlightingStyle: "monokai",
	}

	f, err := os.Create(home + "/.config/xcat/config.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = c.Write(f)
	if err != nil {
		panic(err)
	}

}

func homeDir() string {
	home, _ := os.UserHomeDir()
	if home == "" {
		switch goos := runtime.GOOS; goos {
		case "windows":
			home = os.Getenv("USERPROFILE")
		case "darwin", "linux":
			home = os.Getenv("HOME")
		default:
			home = "."
		}
	}
	return home
}

func loadConfig() (*Config, error) {
	home := homeDir()

	f, err := os.Open(home + "/.config/xcat/config.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := &Config{}
	err = json.NewDecoder(f).Decode(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
