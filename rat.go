//        ____()()     NetRat v0.1 -- Yet another Go experiment
//       /      @@     Copyright (c) 2024 by Giovanni Squillero
// `~~~~~\_;m__m._>o   Distributed under 0BSD (see LICENSE)

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Rat interface {
	String() string
	IsReliable() bool
	FullInfo() string
}

type PublicIpInfo struct {
	Ip        string
	Geo       string
	Source    string
	Timestamp time.Time
	reliable  bool
}

func (ni PublicIpInfo) String() string {
	return ni.Ip
}

func (ni PublicIpInfo) IsReliable() bool {
	return ni.reliable
}

func (ni PublicIpInfo) FullInfo() string {
	return "<not implemented>"
}

func getMullvad() {
	const url string = "https://am.i.mullvad.net/json"

	result, err := http.Get(url)
	if err != nil {
		return
	}
	cooked, err := io.ReadAll(result.Body)
	if err != nil {
		return
	}

	var mullvadInfo map[string]interface{}
	if err := json.Unmarshal([]byte(cooked), &mullvadInfo); err != nil {
		// print out if error is not nil
		fmt.Println(err)
	}

	//log.Println(string(cooked))
	log.Println(mullvadInfo)
}

func getAirVPN() {
	const url string = "https://airvpn.org/api/whatismyip/"

	result, err := http.Get(url)
	if err != nil {
		return
	}
	cooked, err := io.ReadAll(result.Body)
	if err != nil {
		return
	}

	var mullvadInfo map[string]interface{}
	if err := json.Unmarshal([]byte(cooked), &mullvadInfo); err != nil {
		// print out if error is not nil
		fmt.Println(err)
	}

	//log.Println(string(cooked))
	log.Println(mullvadInfo)
}
