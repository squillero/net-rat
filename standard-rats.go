//        ____()()     NetRat v0.1 -- Yet another Go experiment
//       /      @@     Copyright (c) 2024 by Giovanni Squillero
// `~~~~~\_;m__m._>o   Distributed under 0BSD (see LICENSE)

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PublicIpStandard struct {
	Ip        string
	Source    string
	Timestamp time.Time
	reliable  bool
}

func (ni PublicIpStandard) String() string {
	return ni.Ip
}

func (ni PublicIpStandard) GetType() IpType {
	if !ni.reliable {
		return IllegalIP
	} else {
		return PublicIP
	}
}

func genericFetchJson(out chan IpInfo, url, tag string) {
	result, err := http.Get(url)
	if err != nil {
		return
	}
	raw, err := io.ReadAll(result.Body)
	if err != nil {
		return
	}
	var cooked map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &cooked); err != nil {
		fmt.Println(err)
	}
	info := PublicIpStandard{
		Ip:     cooked[tag].(string),
		Source: url,
	}
	out <- info
}
