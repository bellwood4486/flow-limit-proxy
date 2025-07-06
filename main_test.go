package main

import (
	"flag"
	"os"
	"testing"
)

func TestPortParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantFrom int
		wantTo   int
		wantErr  bool
	}{
		{
			name:     "valid ports",
			input:    "8080:9090",
			wantFrom: 8080,
			wantTo:   9090,
			wantErr:  false,
		},
		{
			name:    "invalid format - no colon",
			input:   "8080",
			wantErr: true,
		},
		{
			name:    "invalid format - too many parts",
			input:   "8080:9090:1010",
			wantErr: true,
		},
		{
			name:    "invalid from port",
			input:   "invalid:9090",
			wantErr: true,
		},
		{
			name:    "invalid to port",
			input:   "8080:invalid",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, to, err := parsePortString(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				return
			}
			
			if from != tt.wantFrom {
				t.Errorf("Expected from port %d, got %d", tt.wantFrom, from)
			}
			
			if to != tt.wantTo {
				t.Errorf("Expected to port %d, got %d", tt.wantTo, to)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{
			name:    "valid port 1",
			port:    1,
			wantErr: false,
		},
		{
			name:    "valid port 80",
			port:    80,
			wantErr: false,
		},
		{
			name:    "valid port 8080",
			port:    8080,
			wantErr: false,
		},
		{
			name:    "valid port 65535",
			port:    65535,
			wantErr: false,
		},
		{
			name:    "invalid port 0",
			port:    0,
			wantErr: true,
		},
		{
			name:    "invalid port -1",
			port:    -1,
			wantErr: true,
		},
		{
			name:    "invalid port 65536",
			port:    65536,
			wantErr: true,
		},
		{
			name:    "invalid port 100000",
			port:    100000,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePort(tt.port)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for port %d, but got none", tt.port)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for port %d: %v", tt.port, err)
				}
			}
		})
	}
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		want     *Config
		wantErr  bool
	}{
		{
			name: "valid config",
			args: []string{"cmd", "8080:9090"},
			want: &Config{
				FromPort: 8080,
				ToPort:   9090,
				Limit:    10,
			},
			wantErr: false,
		},
		{
			name: "valid config with limit",
			args: []string{"cmd", "-limit=5", "8080:9090"},
			want: &Config{
				FromPort: 8080,
				ToPort:   9090,
				Limit:    5,
			},
			wantErr: false,
		},
		{
			name:    "invalid port format",
			args:    []string{"cmd", "8080"},
			wantErr: true,
		},
		{
			name:    "invalid port range",
			args:    []string{"cmd", "0:8080"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags for each test
			flag.CommandLine = flag.NewFlagSet(tt.args[0], flag.ContinueOnError)
			
			// Set os.Args for this test
			oldArgs := os.Args
			os.Args = tt.args
			defer func() { os.Args = oldArgs }()
			
			got, err := parseArgs()
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for args %v, but got none", tt.args)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error for args %v: %v", tt.args, err)
				return
			}
			
			if got.FromPort != tt.want.FromPort {
				t.Errorf("Expected FromPort %d, got %d", tt.want.FromPort, got.FromPort)
			}
			
			if got.ToPort != tt.want.ToPort {
				t.Errorf("Expected ToPort %d, got %d", tt.want.ToPort, got.ToPort)
			}
			
			if got.Limit != tt.want.Limit {
				t.Errorf("Expected Limit %d, got %d", tt.want.Limit, got.Limit)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name     string
		fromPort int
		toPort   int
		limit    int64
		want     *Config
		wantErr  bool
	}{
		{
			name:     "valid config",
			fromPort: 8080,
			toPort:   9090,
			limit:    10,
			want: &Config{
				FromPort: 8080,
				ToPort:   9090,
				Limit:    10,
			},
			wantErr: false,
		},
		{
			name:     "valid edge case ports",
			fromPort: 1,
			toPort:   65535,
			limit:    100,
			want: &Config{
				FromPort: 1,
				ToPort:   65535,
				Limit:    100,
			},
			wantErr: false,
		},
		{
			name:     "invalid fromPort - too low",
			fromPort: 0,
			toPort:   8080,
			limit:    10,
			wantErr:  true,
		},
		{
			name:     "invalid fromPort - too high",
			fromPort: 65536,
			toPort:   8080,
			limit:    10,
			wantErr:  true,
		},
		{
			name:     "invalid toPort - negative",
			fromPort: 8080,
			toPort:   -1,
			limit:    10,
			wantErr:  true,
		},
		{
			name:     "invalid toPort - too high",
			fromPort: 8080,
			toPort:   100000,
			limit:    10,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.fromPort, tt.toPort, tt.limit)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for NewConfig(%d, %d, %d), but got none", 
						tt.fromPort, tt.toPort, tt.limit)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error for NewConfig(%d, %d, %d): %v", 
					tt.fromPort, tt.toPort, tt.limit, err)
				return
			}
			
			if got.FromPort != tt.want.FromPort {
				t.Errorf("Expected FromPort %d, got %d", tt.want.FromPort, got.FromPort)
			}
			
			if got.ToPort != tt.want.ToPort {
				t.Errorf("Expected ToPort %d, got %d", tt.want.ToPort, got.ToPort)
			}
			
			if got.Limit != tt.want.Limit {
				t.Errorf("Expected Limit %d, got %d", tt.want.Limit, got.Limit)
			}
		})
	}
}
