package main

import (
	"time"

	"github.com/stonelgh/log"
)

func sample() {
	for {
		if log.GetLevel() >= probeSample10.Level {
			log.Logln("sample:10", probeSample10.Level)
		}
		time.Sleep(3 * time.Second)
	}
}

var probeSample10 = &log.Probe{Package: "github.com/stonelgh/log/main", Module: "sample", Line: 10, Level: log.LvlInfo, Sample: "sample:10", Msg: "sample:10"}

func init() {
	log.RegProbe("github.com/stonelgh/log/main/sample:10", probeSample10)
}
