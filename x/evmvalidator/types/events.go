package types

// Event types for evmvalidator
const (
	EventTypeRegisterValidator           = "register_validator"
	EventTypeDepositCollateral           = "deposit_collateral"
	EventTypeWithdrawCollateral          = "withdraw_collateral"
	EventTypeTransferCollateralOwnership = "transfer_collateral_ownership"
	EventTypeUnjailValidator             = "unjail_validator"
	EventTypeUpdateExtraVotingPower      = "update_extra_voting_power"
	EventTypeUpdateVotingPower           = "update_voting_power"
	EventTypeJailValidator               = "jail_validator"
	EventTypeSlashValidator              = "slash_validator"
	EventTypeWithdrawalMatured           = "withdrawal_matured"

	// Attributes
	AttributeKeyValAddr             = "val_addr"
	AttributeKeyPubkey              = "pubkey"
	AttributeKeyCollateral          = "collateral"
	AttributeKeyCollateralShares    = "collateral_shares"
	AttributeKeyCollateralOwner     = "collateral_owner"
	AttributeKeyCollateralNewOwner  = "collateral_new_owner"
	AttributeKeyExtraVotingPower    = "extra_voting_power"
	AttributeKeyOldExtraVotingPower = "old_extra_voting_power"
	AttributeKeyVotingPower         = "voting_power"
	AttributeKeyOldVotingPower      = "old_voting_power"
	AttributeKeyJailed              = "jailed"
	AttributeKeyWithdrawalID        = "withdrawal_id"
	AttributeKeyAmount              = "amount"
	AttributeKeyShares              = "shares"
	AttributeKeyReceiver            = "receiver"
	AttributeKeyMaturesAt           = "matures_at"
	AttributeKeySlashFraction       = "slash_fraction"
	AttributeKeyInfractionHeight    = "infraction_height"
	AttributeKeyInfractionPower     = "infraction_power"
	AttributeKeyReason              = "reason"
)
