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
	FromPort uint
	ToPort   uint
	Limit    int64
}

// NewConfig creates a new Config with validation
func NewConfig(fromPort, toPort int, limit int64) (*Config, error) {
	if err := validatePort(fromPort); err != nil {
		return nil, fmt.Errorf("invalid fromPort: %w", err)
	}
	
	if err := validatePort(toPort); err != nil {
		return nil, fmt.Errorf("invalid toPort: %w", err)
	}
	
	return &Config{
		FromPort: uint(fromPort),
		ToPort:   uint(toPort),
		Limit:    limit,
	}, nil
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

	if err := ListenProxy(config.FromPort, config.ToPort, config.Limit); err != nil {
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

	return NewConfig(from, to, *limit)
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
