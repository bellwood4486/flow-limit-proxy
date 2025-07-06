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
				MaxConns: 10,
			},
			wantErr: false,
		},
		{
			name: "valid config with limit",
			args: []string{"cmd", "-limit=5", "8080:9090"},
			want: &Config{
				FromPort: 8080,
				ToPort:   9090,
				MaxConns: 5,
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
			
			if got.MaxConns != tt.want.MaxConns {
				t.Errorf("Expected MaxConns %d, got %d", tt.want.MaxConns, got.MaxConns)
			}
		})
	}
}


