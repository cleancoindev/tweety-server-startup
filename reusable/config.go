package reusable

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func get_config_path(app string) string {
	home := os.Getenv("HOME")
	dir := filepath.Join(home, ".config")
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
		dir = filepath.Join(home, "Application Data")
	}
	_, err := os.Stat(dir)
	if err != nil {
		if os.Mkdir(dir, 0700) != nil {
			log.Fatal("failed to create directory:", err)
		}
	}
	file := filepath.Join(dir, app+".json")
	return file
}

func GetConfig(app string) map[string]string {
	config := map[string]string{}
	b, err := ioutil.ReadFile(get_config_path(app))
	if err != nil {
		// No config file found
	} else {
		err = json.Unmarshal(b, &config)
		if err != nil {
			log.Fatal("could not unmarhal ", app, ".json:", err)
		}
	}
	return config
}

func SetConfig(app string, config map[string]string) {
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal("failed to store file:", err)
	}
	err = ioutil.WriteFile(get_config_path(app), b, 0700)
	if err != nil {
		log.Fatal("failed to store file:", err)
	}
}
