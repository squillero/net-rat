//        ____()()     NetRat v0.1
//       /      @@     ~~~~~~~~~~~
// `~~~~~\_;m__m._>o   Yet another little Go experiment
//
// Copyright © 2024 Giovanni Squillero <giovanni.squillero@polito.it>
// This code is being released for educational and academic purposes.

package main

import (
	"context"
	"fmt"
	"log"
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

type IpType int

const (
	IllegalIP IpType = iota
	PrivateIP
	PublicIP
)

type IpInfo interface {
	GetType() IpType
	String() string
}

func eager() {
	ch := make(chan IpInfo, 1)

	// Create a context with a timeout of 1 seconds
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	// Start the doSomething function
	go getMullvad(ch)
	go getAirVPN(ch)

	select {
	case <-ctxTimeout.Done():
		log.Fatalf("Context cancelled: %v\n", ctxTimeout.Err())
	case result := <-ch:
		fmt.Printf("Received: %s\n", result)
	}
}
