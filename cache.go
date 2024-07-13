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

const INFO_TIMEOUT = 5 * time.Minute

var CacheFile = filepath.Join(os.TempDir(), "netrat.json")

type CacheError struct {
	croak string
}

func (e CacheError) Error() string {
	return e.croak
}

func CacheSave(info NetInfo) error {
	slog.Debug("Writing cache", "file", CacheFile)

	var cache []IpInfo
	for _, v := range info.ips {
		cache = append(cache, v)
	}
	file, _ := json.MarshalIndent(cache, "", "    ")
	err := os.WriteFile(CacheFile, file, 0600)
	return err
}

func CacheLoad() (NetInfo, error) {
	var err error
	var data []byte
	var cache []IpInfo
	if data, err = os.ReadFile(CacheFile); err != nil {
		return NewNetInfo(), err
	}
	if err = json.Unmarshal(data, &cache); err != nil {
		return NewNetInfo(), err
	}
	info := NewNetInfo()
	slog.Debug("Using cache", "file", CacheFile)
	for _, v := range cache {
		info.add(v)
	}

	return info, nil
}
