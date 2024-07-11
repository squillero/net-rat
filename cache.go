//        ____()()     NetRat v0.1
//       /      @@     ~~~~~~~~~~~
// `~~~~~\_;m__m._>o   A Go hack
//
// Coded in July 2024, between Italy and Australia (34,138 km).
// Copyright © 2024 Giovanni Squillero <giovanni.squillero@polito.it>
// Released for educational and academic purposes under 0BSD.

package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const CACHE_TIMEOUT = 5 * time.Minute

var CacheFile = filepath.Join(os.TempDir(), "netrat.json")

type CacheError struct {
	croak string
}

func (e CacheError) Error() string {
	return e.croak
}

func CacheSave(info [2]IpInfo) error {
	slog.Debug("Writing cache", "file", CacheFile)

	file, _ := json.MarshalIndent(info, "", "    ")
	err := os.WriteFile(CacheFile, file, 0600)
	return err
}

func CacheLoad() ([2]IpInfo, error) {
	var info [2]IpInfo
	var err error
	data, err := os.ReadFile(CacheFile)
	if err != nil {
		return [2]IpInfo{}, err
	}
	err = json.Unmarshal(data, &info)
	if err != nil {
		return [2]IpInfo{}, err
	}
	return info, nil
}
