package config

import (
	"encoding/xml"
	"os"
)

func LoadXmlConfig(filename string, v interface{}) error {
	if contents, err := os.ReadFile(filename); err != nil {
		return err
	} else {
		if err = xml.Unmarshal(contents, v); err != nil {
			return err
		}
		return nil
	}
}

func SaveXmlConfig(filename string, v interface{}) error {
	if contents, err := xml.Marshal(v); err != nil {
		return err
	} else {
		if err = os.WriteFile(filename, contents, 0644); err != nil {
			return err
		}
		return nil
	}
}
