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
	from, to, err := parsePortString(flag.Arg(0))
	if err != nil {
		log.Fatalf("invalid port format: %v\n", err)
	}
	if err := validatePort(from); err != nil {
		log.Fatalf("invalid fromPort: %v\n", err)
	}
	if err := validatePort(to); err != nil {
		log.Fatalf("invalid toPort: %v\n", err)
	}

	log.SetPrefix(fmt.Sprintf("[flproxy(%d->%d)] ", from, to))

	if err := ListenProxy(uint(from), uint(to), int64(*limit)); err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
}

// parsePortString parses a port string in format "from:to" and returns the port numbers
func parsePortString(portStr string) (int, int, error) {
	ports := strings.Split(portStr, ":")
	if len(ports) != 2 {
		return 0, 0, fmt.Errorf("invalid port format, expected 'from:to', got '%s'", portStr)
	}
	
	from, err := strconv.Atoi(ports[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid fromPort '%s': %w", ports[0], err)
	}
	
	to, err := strconv.Atoi(ports[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid toPort '%s': %w", ports[1], err)
	}
	
	return from, to, nil
}

// validatePort validates that a port number is within the valid range
func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", port)
	}
	return nil
}
