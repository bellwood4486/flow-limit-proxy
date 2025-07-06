package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Config holds the application configuration
type Config struct {
	FromPort int
	ToPort   int
	Limit    int64
}

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
	config, err := parseArgs()
	if err != nil {
		log.Fatalf("configuration error: %v\n", err)
	}

	log.SetPrefix(fmt.Sprintf("[flproxy(%d->%d)] ", config.FromPort, config.ToPort))

	if err := ListenProxy(uint(config.FromPort), uint(config.ToPort), config.Limit); err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
}

// parseArgs parses command line arguments and returns configuration
func parseArgs() (*Config, error) {
	limit := flag.Int64("limit", 10, "concurrent transfer limit")
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	from, to, err := parsePortString(flag.Arg(0))
	if err != nil {
		return nil, fmt.Errorf("invalid port format: %w", err)
	}

	if err := validatePort(from); err != nil {
		return nil, fmt.Errorf("invalid fromPort: %w", err)
	}

	if err := validatePort(to); err != nil {
		return nil, fmt.Errorf("invalid toPort: %w", err)
	}

	return &Config{
		FromPort: from,
		ToPort:   to,
		Limit:    *limit,
	}, nil
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
