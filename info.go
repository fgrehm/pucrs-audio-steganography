package main

import (
	"io/ioutil"
	"encoding/json"
)

type stegInfo struct {
	Filename string
}

func readInfo(path string) (*stegInfo, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	info := &stegInfo{}
	if err := json.Unmarshal(data, info); err != nil {
		return nil, err
	}

	return info, nil
}

func writeInfo(sourcePath, payloadPath, infoPath, filename string) error {
	info := stegInfo{Filename: filename}
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(infoPath, data, 0644)
}
