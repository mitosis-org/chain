package validation

import (
	"testing"

	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func setupValidOnlineFlags(cmd *cobra.Command) {
	_ = flags.StringFlagWithValue(cmd, "rpc-url", "http://localhost:8545", "")
	_ = flags.StringFlagWithValue(cmd, "validator-manager-contract-addr", "0x123...", "")
}

func setupValidOfflineFlags(cmd *cobra.Command) {
	_ = flags.StringFlagWithValue(cmd, "chain-id", "1", "")
	_ = flags.StringFlagWithValue(cmd, "gas-limit", "21000", "")
	_ = flags.StringFlagWithValue(cmd, "gas-price", "20000000000", "")
	_ = flags.StringFlagWithValue(cmd, "nonce", "0", "")
	_ = flags.StringFlagWithValue(cmd, "contract-fee", "1000000000000000000", "")
}

func TestValidateCreateTxFlagGroups(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*cobra.Command)
		wantErr bool
	}{
		{
			name: "valid flags with private key and online mode",
			setup: func(cmd *cobra.Command) {
				_ = flags.StringFlagWithValue(cmd, "private-key", "some-private-key", "")
				_ = flags.StringFlagWithValue(cmd, "signed", "signed-tx", "")
				setupValidOnlineFlags(cmd)
			},
			wantErr: false,
		},
		{
			name: "valid flags with keyfile and online mode",
			setup: func(cmd *cobra.Command) {
				_ = flags.StringFlagWithValue(cmd, "keyfile", "path/to/keyfile", "")
				_ = flags.StringFlagWithValue(cmd, "keyfile-password", "password123", "")
				_ = flags.StringFlagWithValue(cmd, "signed", "signed-tx", "")
				setupValidOnlineFlags(cmd)
			},
			wantErr: false,
		},
		{
			name: "valid flags with private key and offline mode",
			setup: func(cmd *cobra.Command) {
				_ = flags.StringFlagWithValue(cmd, "private-key", "some-private-key", "")
				_ = flags.StringFlagWithValue(cmd, "signed", "signed-tx", "")
				setupValidOfflineFlags(cmd)
			},
			wantErr: false,
		},
		{
			name:    "missing all flags",
			setup:   func(cmd *cobra.Command) {},
			wantErr: true,
		},
		{
			name: "missing network mode flags",
			setup: func(cmd *cobra.Command) {
				_ = flags.StringFlagWithValue(cmd, "private-key", "some-private-key", "")
				_ = flags.StringFlagWithValue(cmd, "signed", "signed-tx", "")
			},
			wantErr: true,
		},
		{
			name: "conflicting signing methods",
			setup: func(cmd *cobra.Command) {
				_ = flags.StringFlagWithValue(cmd, "private-key", "some-private-key", "")
				_ = flags.StringFlagWithValue(cmd, "keyfile", "path/to/keyfile", "")
				_ = flags.StringFlagWithValue(cmd, "keyfile-password", "password123", "")
				_ = flags.StringFlagWithValue(cmd, "signed", "signed-tx", "")
				setupValidOnlineFlags(cmd)
			},
			wantErr: true,
		},
		{
			name: "conflicting network modes",
			setup: func(cmd *cobra.Command) {
				_ = flags.StringFlagWithValue(cmd, "private-key", "some-private-key", "")
				_ = flags.StringFlagWithValue(cmd, "signed", "signed-tx", "")
				setupValidOnlineFlags(cmd)
				setupValidOfflineFlags(cmd)
			},
			wantErr: true,
		},
		{
			name: "incomplete offline mode flags",
			setup: func(cmd *cobra.Command) {
				_ = flags.StringFlagWithValue(cmd, "private-key", "some-private-key", "")
				_ = flags.StringFlagWithValue(cmd, "signed", "signed-tx", "")
				_ = flags.StringFlagWithValue(cmd, "chain-id", "1", "")
				// Missing other offline mode flags
			},
			wantErr: true,
		},
		{
			name: "incomplete online mode flags",
			setup: func(cmd *cobra.Command) {
				_ = flags.StringFlagWithValue(cmd, "rpc-url", "http://localhost:8545", "")
				// Missing validator-manager-contract-addr
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			if tt.setup != nil {
				tt.setup(cmd)
			}
			err := ValidateCreateTxFlagGroups(cmd)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSendTxFlagGroups(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*cobra.Command)
		wantErr bool
	}{
		{
			name: "valid flags with private key",
			setup: func(cmd *cobra.Command) {
				_ = flags.StringFlagWithValue(cmd, "private-key", "some-private-key", "")
				_ = flags.StringFlagWithValue(cmd, "signed", "signed-tx", "")
				setupValidOnlineFlags(cmd)
			},
			wantErr: false,
		},
		{
			name: "valid flags with keyfile",
			setup: func(cmd *cobra.Command) {
				_ = flags.StringFlagWithValue(cmd, "keyfile", "path/to/keyfile", "")
				_ = flags.StringFlagWithValue(cmd, "keyfile-password", "password123", "")
				_ = flags.StringFlagWithValue(cmd, "signed", "signed-tx", "")
				setupValidOnlineFlags(cmd)
			},
			wantErr: false,
		},
		{
			name:    "missing all flags",
			setup:   func(cmd *cobra.Command) {},
			wantErr: true,
		},
		{
			name: "missing online mode flags",
			setup: func(cmd *cobra.Command) {
				_ = flags.StringFlagWithValue(cmd, "private-key", "some-private-key", "")
				_ = flags.StringFlagWithValue(cmd, "signed", "signed-tx", "")
			},
			wantErr: true,
		},
		{
			name: "missing validator manager contract address",
			setup: func(cmd *cobra.Command) {
				_ = flags.StringFlagWithValue(cmd, "private-key", "some-private-key", "")
				_ = flags.StringFlagWithValue(cmd, "signed", "signed-tx", "")
				_ = flags.StringFlagWithValue(cmd, "rpc-url", "http://localhost:8545", "")
			},
			wantErr: true,
		},
		{
			name: "conflicting signing methods",
			setup: func(cmd *cobra.Command) {
				_ = flags.StringFlagWithValue(cmd, "private-key", "some-private-key", "")
				_ = flags.StringFlagWithValue(cmd, "keyfile", "path/to/keyfile", "")
				_ = flags.StringFlagWithValue(cmd, "keyfile-password", "password123", "")
				_ = flags.StringFlagWithValue(cmd, "signed", "signed-tx", "")
				setupValidOnlineFlags(cmd)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			if tt.setup != nil {
				tt.setup(cmd)
			}
			err := ValidateSendTxFlagGroups(cmd)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
