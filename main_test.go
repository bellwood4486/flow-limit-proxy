package main

import (
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