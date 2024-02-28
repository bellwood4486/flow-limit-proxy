package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usages:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  %s [-limit=<number>] <fromPort>:<toPort>\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}
	log.SetPrefix("[flproxy] ")
}

func main() {
	// limit
	limit := flag.Int64("limit", 10, "concurrent transfer limit")
	flag.Parse()

	// port
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	ports := strings.Split(flag.Arg(0), ":")
	if len(ports) != 2 {
		flag.Usage()
		os.Exit(1)
	}
	var from, to int
	var err error
	if from, err = strconv.Atoi(ports[0]); err != nil {
		log.Fatalf("invalid fromPort: %v\n", err)
	}
	if to, err = strconv.Atoi(ports[1]); err != nil {
		log.Fatalf("invalid toPort: %v\n", err)
	}

	log.SetPrefix(fmt.Sprintf("[flproxy(%d->%d)] ", from, to))

	if err := ListenProxy(uint(from), uint(to), int64(*limit)); err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
}
