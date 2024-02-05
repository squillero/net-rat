//        ____()()     NetRat v0.1 -- Yet another Go experiment
//       /      @@     Copyright (c) 2024 by Giovanni Squillero
// `~~~~~\_;m__m._>o   Distributed under 0BSD (see LICENSE)

package main

import "time"

type Snitch interface {
	String() string
	IsReliable() bool
	FullInfo() string
}

type NetInfo struct {
	Ip        string
	Source    string
	Timestamp time.Time
	reliable  bool
}

func (ni NetInfo) String() string {
	return ni.Ip
}
func (ni NetInfo) IsReliable() bool {
	return ni.reliable
}
func (ni NetInfo) FullInfo() string {
	return "<not implemented>"
}
