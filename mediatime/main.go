package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/zencoder/bolt-fmp4/helpers"
)

var (
	fromTimescale uint64
	toTimescale   uint64
)

func main() {
	flag.Uint64Var(&fromTimescale, "from", 90000, "timescale to convert from")
	flag.Uint64Var(&toTimescale, "to", 1000, "timescale to convert to")
	flag.Parse()

	timeStr := flag.Arg(0)

	if timeStr == "" {
		println("must provide time as an argument")
		os.Exit(1)
	}

	time, err := strconv.ParseUint(timeStr, 10, 64)
	if err != nil {
		panic(err)
	}

	res := helpers.ConvertTimeRound(fromTimescale, toTimescale, time)

	println(res)
}
