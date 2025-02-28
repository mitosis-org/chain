package types

// Event types for evmvalidator
const (
	EventTypeRegisterValidator  = "register_validator"
	EventTypeDepositCollateral  = "deposit_collateral"
	EventTypeWithdrawCollateral = "withdraw_collateral"
	EventTypeUnjailValidator    = "unjail_validator"
	EventTypeUpdateVotingPower  = "update_voting_power"
	EventTypeJailValidator      = "jail_validator"
	EventTypeSlashValidator     = "slash_validator"
	EventTypeUpdateParams       = "update_params"
	EventTypeWithdrawalMatured  = "withdrawal_matured"

	// Attributes
	AttributeKeyPubkey           = "pubkey"
	AttributeKeyCollateral       = "collateral"
	AttributeKeyExtraVotingPower = "extra_voting_power"
	AttributeKeyVotingPower      = "voting_power"
	AttributeKeyAmount           = "amount"
	AttributeKeyReceiver         = "receiver"
	AttributeKeyReceivesAt       = "receives_at"
	AttributeKeySlashFraction    = "slash_fraction"
	AttributeKeyInfractionHeight = "infraction_height"
	AttributeKeyInfractionPower  = "infraction_power"
	AttributeKeyReason           = "reason"
)
