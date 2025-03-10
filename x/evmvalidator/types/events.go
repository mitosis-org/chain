package types

// Event types for evmvalidator
const (
	EventTypeRegisterValidator      = "register_validator"
	EventTypeDepositCollateral      = "deposit_collateral"
	EventTypeWithdrawCollateral     = "withdraw_collateral"
	EventTypeUnjailValidator        = "unjail_validator"
	EventTypeUpdateExtraVotingPower = "update_extra_voting_power"
	EventTypeUpdateVotingPower      = "update_voting_power"
	EventTypeJailValidator          = "jail_validator"
	EventTypeSlashValidator         = "slash_validator"
	EventTypeUpdateParams           = "update_params"
	EventTypeWithdrawalMatured      = "withdrawal_matured"

	// Attributes
	AttributeKeyValAddr             = "val_addr"
	AttributeKeyPubkey              = "pubkey"
	AttributeKeyCollateral          = "collateral"
	AttributeKeyExtraVotingPower    = "extra_voting_power"
	AttributeKeyOldExtraVotingPower = "old_extra_voting_power"
	AttributeKeyVotingPower         = "voting_power"
	AttributeKeyOldVotingPower      = "old_voting_power"
	AttributeKeyJailed              = "jailed"
	AttributeKeyAmount              = "amount"
	AttributeKeyReceiver            = "receiver"
	AttributeKeyMaturesAt           = "matures_at"
	AttributeKeySlashFraction       = "slash_fraction"
	AttributeKeyInfractionHeight    = "infraction_height"
	AttributeKeyInfractionPower     = "infraction_power"
	AttributeKeyReason              = "reason"
)
