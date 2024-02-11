//        ____()()     NetRat v0.1 -- Yet another Go experiment
//       /      @@     Copyright (c) 2024 by Giovanni Squillero
// `~~~~~\_;m__m._>o   Distributed under 0BSD (see LICENSE)

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// https://api.ipify.org?format=json
// http://api.ipify.org?format=json
// http://ipinfo.io
// https://freeipapi.com/api/json
// http://ipv4.iplocation.net
// https://am.i.mullvad.net/json
// https://airvpn.org/api/whatismyip

// -- https://iplocation.io/

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
	return "ip: " + ni.Ip
}

func (ni PublicIpInfo) IsReliable() bool {
	return ni.reliable
}

func (ni PublicIpInfo) FullInfo() string {
	return "<not implemented>"
}

func getMullvad(out chan PublicIpInfo) {
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

	info := PublicIpInfo{
		Ip:     mullvadInfo["ip"].(string),
		Geo:    mullvadInfo["country"].(string),
		Source: "Mullvad",
	}
	log.Println(info)
	out <- info
}

func getAirVPN(out chan PublicIpInfo) {
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

	info := PublicIpInfo{
		Ip:     mullvadInfo["ip"].(string),
		Geo:    (mullvadInfo["geo"].(map[string]interface{}))["name"].(string),
		Source: "AirVPN",
	}
	log.Println(info)
	out <- info

}

func eager() {
	ch := make(chan PublicIpInfo, 1)

	// Create a context with a timeout of 1 seconds
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	// Start the doSomething function
	go getMullvad(ch)
	go getAirVPN(ch)

	select {
	case <-ctxTimeout.Done():
		fmt.Printf("Context cancelled: %v\n", ctxTimeout.Err())
	case result := <-ch:
		fmt.Printf("Received: %s\n", result)
	}
}
