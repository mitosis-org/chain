package testutil

import (
	"crypto/ecdsa"
	"encoding/hex"
	"strings"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/ethereum/go-ethereum/crypto"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/mitosis-org/chain/x/evmvalidator/keeper"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/stretchr/testify/suite"
)

// TestKeeper is a minimal working keeper for testing
type TestKeeper struct {
	Keeper     *keeper.Keeper
	Ctx        sdk.Context
	Cdc        codec.Codec
	StoreKey   storetypes.StoreKey
	MockSlash  *MockSlashingKeeper
	MockEvmEng *MockEvmEngineKeeper
}

// GenerateSecp256k1Key generates a new secp256k1 private key and returns the private key, compressed pubkey, and eth address
func GenerateSecp256k1Key() (*ecdsa.PrivateKey, []byte, mitotypes.EthAddress) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	// Get the compressed public key
	compressedPubKey := crypto.CompressPubkey(&privateKey.PublicKey)

	// Get the Ethereum address
	addr := crypto.PubkeyToAddress(privateKey.PublicKey)
	ethAddr := mitotypes.EthAddress(addr)

	return privateKey, compressedPubKey, ethAddr
}

// PubkeyToConsAddr converts a pubkey to a consensus address
func PubkeyToConsAddr(pubkey []byte) sdk.ConsAddress {
	pubKey := &secp256k1.PubKey{Key: pubkey}

	// Get consensus key
	cpk, err := cryptocodec.ToCmtPubKeyInterface(pubKey)
	if err != nil {
		panic(err)
	}

	// Get consensus address from consensus key
	return sdk.ConsAddress(cpk.Address())
}

// HexToCompressedPubkey converts a hex string to compressed pubkey bytes
func HexToCompressedPubkey(hexPubkey string) []byte {
	hexPubkey = strings.TrimPrefix(hexPubkey, "0x")
	bz, err := hex.DecodeString(hexPubkey)
	if err != nil {
		panic(err)
	}
	return bz
}

// NewTestKeeper returns a test keeper with minimal working dependencies
func NewTestKeeper(s *suite.Suite) TestKeeper {
	// Create store keys
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), storeKey, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	// Create mock keepers
	mockSlash := &MockSlashingKeeper{}
	mockEvmEng := &MockEvmEngineKeeper{}

	// Create keeper
	k := keeper.NewKeeper(
		encCfg.Codec,
		storeKey,
		address.NewBech32Codec("mitovaloper"),
		address.NewBech32Codec("mitovalcons"),
		"evmgov",
	)

	k.SetSlashingKeeper(mockSlash)
	k.SetEvmEngineKeeper(mockEvmEng)

	return TestKeeper{
		Keeper:     k,
		Ctx:        ctx,
		Cdc:        encCfg.Codec,
		StoreKey:   storeKey,
		MockSlash:  mockSlash,
		MockEvmEng: mockEvmEng,
	}
}

// SetupDefaultTestParams sets up test parameters with commonly used default values
func (tk TestKeeper) SetupDefaultTestParams() types.Params {
	return tk.SetupTestParams(
		types.Params{
			MaxValidators:    100,
			MaxLeverageRatio: math.LegacyNewDec(10),
			MinVotingPower:   1,
			WithdrawalLimit:  10,
		},
	)
}

// SetupTestParams sets up standard test parameters for tests
func (tk TestKeeper) SetupTestParams(params types.Params) types.Params {
	err := tk.Keeper.SetParams(tk.Ctx, params)
	if err != nil {
		panic(err)
	}
	return params
}

// RegisterTestValidator is a helper function to register a validator for testing
func (tk TestKeeper) RegisterTestValidator(collateral math.Uint, extraVotingPower math.Uint, jailed bool) types.Validator {
	_, pubkey, valAddr := GenerateSecp256k1Key()
	err := tk.Keeper.RegisterValidator(tk.Ctx, valAddr, pubkey, collateral, extraVotingPower, jailed)
	if err != nil {
		panic(err)
	}

	validator, found := tk.Keeper.GetValidator(tk.Ctx, valAddr)
	if !found {
		panic("validator not found after registration")
	}
	return validator
}
