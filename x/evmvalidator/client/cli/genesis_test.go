package cli

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test data - valid secp256k1 compressed public key (33 bytes)
const (
	validHexPubkey    = "03107dd702ec9618b9f928a308f4fb6719ac6f4e21d667bde8dde2291b2d3375d7"
	validBase64Pubkey = "AxB91wLslhi5+SijCPT7Zxmsb04h1me96N3iKRstM3XX"

	// Invalid test cases
	invalidShortHex = "03107dd702ec9618b9f928a308f4fb6719ac6f4e21d667bde8dde2291b2d3375"     // 32 bytes
	invalidLongHex  = "03107dd702ec9618b9f928a308f4fb6719ac6f4e21d667bde8dde2291b2d3375d7aa" // 34 bytes
	invalidBase64   = "InvalidBase64String!!!"
	invalidHexChars = "03107dd702ec9618b9f928a308f4fb6719ac6f4e21d667bde8dde2291b2d3375gx" // contains 'g' and 'x'
	emptyString     = ""
)

func TestParsePubkey(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		{
			name:        "valid hex pubkey",
			input:       validHexPubkey,
			expectError: false,
			description: "should parse valid 66-character hex string",
		},
		{
			name:        "valid hex pubkey with 0x prefix",
			input:       "0x" + validHexPubkey,
			expectError: false,
			description: "should parse valid hex string with 0x prefix",
		},
		{
			name:        "valid base64 pubkey",
			input:       validBase64Pubkey,
			expectError: false,
			description: "should parse valid base64 encoded pubkey",
		},
		{
			name:        "uppercase hex pubkey",
			input:       "03107DD702EC9618B9F928A308F4FB6719AC6F4E21D667BDE8DDE2291B2D3375D7",
			expectError: false,
			description: "should parse uppercase hex string",
		},
		{
			name:        "mixed case hex pubkey",
			input:       "03107Dd702Ec9618B9f928A308F4fB6719aC6f4E21d667BdE8DdE2291b2D3375D7",
			expectError: false,
			description: "should parse mixed case hex string",
		},
		{
			name:        "short hex pubkey",
			input:       invalidShortHex,
			expectError: true,
			description: "should fail for pubkey shorter than 33 bytes",
		},
		{
			name:        "long hex pubkey",
			input:       invalidLongHex,
			expectError: true,
			description: "should fail for pubkey longer than 33 bytes",
		},
		{
			name:        "invalid hex characters",
			input:       invalidHexChars,
			expectError: true,
			description: "should fail for invalid hex characters",
		},
		{
			name:        "invalid base64",
			input:       invalidBase64,
			expectError: true,
			description: "should fail for invalid base64 string",
		},
		{
			name:        "empty string",
			input:       emptyString,
			expectError: true,
			description: "should fail for empty string",
		},
		{
			name:        "short hex without 0x",
			input:       "03107dd7",
			expectError: true,
			description: "should fail for short hex string",
		},
		{
			name:        "valid length but invalid format",
			input:       "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
			expectError: true,
			description: "should fail for 66-char string that's not valid hex or base64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pubkey, err := parsePubkey(tt.input)

			if tt.expectError {
				require.Error(t, err, tt.description)
				require.Nil(t, pubkey, "pubkey should be nil on error")
			} else {
				require.NoError(t, err, tt.description)
				require.NotNil(t, pubkey, "pubkey should not be nil on success")
				require.Len(t, pubkey, 33, "pubkey should be 33 bytes")

				// Verify the pubkey is the expected value
				expectedPubkey, _ := hex.DecodeString(validHexPubkey)
				require.Equal(t, expectedPubkey, pubkey, "parsed pubkey should match expected value")
			}
		})
	}
}

func TestParsePubkey_HexVsBase64Consistency(t *testing.T) {
	// Test that the same pubkey parsed from hex and base64 produces identical results
	hexPubkey, err := parsePubkey(validHexPubkey)
	require.NoError(t, err)

	base64Pubkey, err := parsePubkey(validBase64Pubkey)
	require.NoError(t, err)

	require.Equal(t, hexPubkey, base64Pubkey, "hex and base64 parsing should produce identical results")
}

func TestParsePubkey_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		{
			name:        "32 byte hex (uncompressed without prefix)",
			input:       "107dd702ec9618b9f928a308f4fb6719ac6f4e21d667bde8dde2291b2d3375d7", // 32 bytes (64 chars)
			expectError: true,
			description: "should fail for 32-byte hex (missing compression prefix)",
		},
		{
			name:        "valid base64 but wrong length",
			input:       base64.StdEncoding.EncodeToString([]byte("this is exactly 32 bytes long!!")), // 32 bytes
			expectError: true,
			description: "should fail for valid base64 but wrong pubkey length",
		},
		{
			name:        "hex with spaces",
			input:       "03 10 7d d7 02 ec 96 18 b9 f9 28 a3 08 f4 fb 67 19 ac 6f 4e 21 d6 67 bd e8 dd e2 29 1b 2d 33 75 d7",
			expectError: true,
			description: "should fail for hex with spaces",
		},
		{
			name:        "base64 with padding issues",
			input:       validBase64Pubkey + "=", // Extra padding
			expectError: true,
			description: "should handle base64 padding correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pubkey, err := parsePubkey(tt.input)

			if tt.expectError {
				require.Error(t, err, tt.description)
				require.Nil(t, pubkey)
			} else {
				require.NoError(t, err, tt.description)
				require.NotNil(t, pubkey)
				require.Len(t, pubkey, 33)
			}
		})
	}
}

func TestIsHexString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid lowercase hex",
			input:    "0123456789abcdef",
			expected: true,
		},
		{
			name:     "valid uppercase hex",
			input:    "0123456789ABCDEF",
			expected: true,
		},
		{
			name:     "valid mixed case hex",
			input:    "0123456789AbCdEf",
			expected: true,
		},
		{
			name:     "invalid character g",
			input:    "0123456789abcdeg",
			expected: false,
		},
		{
			name:     "invalid character z",
			input:    "0123456789abcdez",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: true, // empty string is technically valid hex
		},
		{
			name:     "single valid hex char",
			input:    "a",
			expected: true,
		},
		{
			name:     "single invalid char",
			input:    "g",
			expected: false,
		},
		{
			name:     "special characters",
			input:    "012345!@#$",
			expected: false,
		},
		{
			name:     "numbers only",
			input:    "0123456789",
			expected: true,
		},
		{
			name:     "letters only lowercase",
			input:    "abcdef",
			expected: true,
		},
		{
			name:     "letters only uppercase",
			input:    "ABCDEF",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isHexString(tt.input)
			require.Equal(t, tt.expected, result, "isHexString result should match expected value")
		})
	}
}

func TestParsePubkey_RealWorldExamples(t *testing.T) {
	// Test with real-world examples from different sources
	realWorldTests := []struct {
		name         string
		hexFormat    string
		base64Format string
		description  string
	}{
		{
			name:         "test pubkey 1",
			hexFormat:    "03107dd702ec9618b9f928a308f4fb6719ac6f4e21d667bde8dde2291b2d3375d7",
			base64Format: "AxB91wLslhi5+SijCPT7Zxmsb04h1me96N3iKRstM3XX",
			description:  "main test pubkey used in examples",
		},
		{
			name:         "test pubkey 2",
			hexFormat:    "02a98478cf8213c7fea5a328d89675b5b544fb0c677893690b88473aa3aac0f3ec",
			base64Format: "AqmEeM+CE8f+paMo2JZ1tbVE+wxneJNpC4hHOqOqwPPs",
			description:  "another valid secp256k1 compressed pubkey",
		},
	}

	for _, tt := range realWorldTests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse hex format
			hexPubkey, err := parsePubkey(tt.hexFormat)
			require.NoError(t, err, "hex format should parse successfully")
			require.Len(t, hexPubkey, 33, "hex pubkey should be 33 bytes")

			// Parse base64 format
			base64Pubkey, err := parsePubkey(tt.base64Format)
			require.NoError(t, err, "base64 format should parse successfully")
			require.Len(t, base64Pubkey, 33, "base64 pubkey should be 33 bytes")

			// They should be identical
			require.Equal(t, hexPubkey, base64Pubkey, "hex and base64 should produce identical pubkeys")

			// Verify round-trip conversion
			hexStr := hex.EncodeToString(hexPubkey)
			base64Str := base64.StdEncoding.EncodeToString(base64Pubkey)

			require.Equal(t, tt.hexFormat, hexStr, "hex round-trip should match")
			require.Equal(t, tt.base64Format, base64Str, "base64 round-trip should match")
		})
	}
}

func TestParsePubkey_Performance(t *testing.T) {
	// Simple performance test to ensure parsing is efficient
	const iterations = 1000

	// Test hex parsing performance
	t.Run("hex parsing performance", func(t *testing.T) {
		for i := 0; i < iterations; i++ {
			_, err := parsePubkey(validHexPubkey)
			require.NoError(t, err)
		}
	})

	// Test base64 parsing performance
	t.Run("base64 parsing performance", func(t *testing.T) {
		for i := 0; i < iterations; i++ {
			_, err := parsePubkey(validBase64Pubkey)
			require.NoError(t, err)
		}
	})
}

func BenchmarkParsePubkey_Hex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = parsePubkey(validHexPubkey)
	}
}

func BenchmarkParsePubkey_Base64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = parsePubkey(validBase64Pubkey)
	}
}

func BenchmarkIsHexString(b *testing.B) {
	testStr := validHexPubkey
	for i := 0; i < b.N; i++ {
		_ = isHexString(testStr)
	}
}

// Additional integration tests can be added here when needed.
// For now, we focus on testing the core parsing functionality.
