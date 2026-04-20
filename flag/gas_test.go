package flag

import (
	"testing"
)

func TestGasValue_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint64
		wantErr bool
	}{
		// Raw wei values
		{
			name:  "raw wei integer",
			input: "1000000000",
			want:  1000000000,
		},
		{
			name:  "raw wei zero",
			input: "0",
			want:  0,
		},
		{
			name:  "raw wei large value",
			input: "100000000000000000",
			want:  100000000000000000,
		},
		// Gwei values
		{
			name:  "1 gwei",
			input: "1gwei",
			want:  1e9,
		},
		{
			name:  "100 gwei",
			input: "100gwei",
			want:  100e9,
		},
		{
			name:  "gwei with space",
			input: "100 gwei",
			want:  100e9,
		},
		{
			name:  "decimal gwei",
			input: "0.1gwei",
			want:  1e8,
		},
		{
			name:  "decimal gwei 2.5",
			input: "2.5gwei",
			want:  25e8,
		},
		{
			name:  "small decimal gwei",
			input: "0.001gwei",
			want:  1e6,
		},
		// Ether values
		{
			name:  "1 ether",
			input: "1ether",
			want:  1e18,
		},
		{
			name:  "0.001 ether",
			input: "0.001ether",
			want:  1e15,
		},
		{
			name:  "1 eth alias",
			input: "1eth",
			want:  1e18,
		},
		// Case insensitivity
		{
			name:  "uppercase GWEI",
			input: "100GWEI",
			want:  100e9,
		},
		{
			name:  "mixed case Gwei",
			input: "100Gwei",
			want:  100e9,
		},
		{
			name:  "uppercase ETHER",
			input: "1ETHER",
			want:  1e18,
		},
		{
			name:  "uppercase ETH",
			input: "1ETH",
			want:  1e18,
		},
		// Wei explicit
		{
			name:  "explicit wei",
			input: "1000wei",
			want:  1000,
		},
		// Error cases
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid unit",
			input:   "100foo",
			wantErr: true,
		},
		{
			name:    "invalid number",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "negative value",
			input:   "-100gwei",
			wantErr: true,
		},
		{
			name:    "overflow",
			input:   "1000000ether",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var val uint64
			gv := &GasValue{Val: &val}
			err := gv.Set(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GasValue.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && val != tt.want {
				t.Errorf("GasValue.Set() = %v, want %v", val, tt.want)
			}
		})
	}
}

func TestGasValue_String(t *testing.T) {
	tests := []struct {
		name string
		val  uint64
		want string
	}{
		{
			name: "zero",
			val:  0,
			want: "0",
		},
		{
			name: "1 gwei",
			val:  1e9,
			want: "1000000000",
		},
		{
			name: "100 gwei",
			val:  100e9,
			want: "100000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := tt.val
			gv := &GasValue{Val: &val}
			if got := gv.String(); got != tt.want {
				t.Errorf("GasValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGasValue_String_Nil(t *testing.T) {
	gv := &GasValue{Val: nil}
	if got := gv.String(); got != "0" {
		t.Errorf("GasValue.String() with nil = %v, want 0", got)
	}
}

func TestGasValue_Type(t *testing.T) {
	gv := &GasValue{}
	if got := gv.Type(); got != "gas" {
		t.Errorf("GasValue.Type() = %v, want gas", got)
	}
}

func TestParseGasUnit(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint64
		wantErr bool
	}{
		{
			name:  "100gwei",
			input: "100gwei",
			want:  100e9,
		},
		{
			name:  "raw number",
			input: "1000000000",
			want:  1000000000,
		},
		{
			name:    "invalid",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGasUnit(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGasUnit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseGasUnit() = %v, want %v", got, tt.want)
			}
		})
	}
}
