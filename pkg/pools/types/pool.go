package types

import (
	"github.com/dfuse-io/solana-go"
)

var EmptyAddress = solana.MustPublicKeyFromBase58("11111111111111111111111111111111")

const (
	ParrotPool   = "Parrot"
	MarinadePool = "Marinade"
	SolidoPool   = "Solido"
	EverSOL      = "Eversol"
	JPool        = "JPool"
	Socean       = "Socean"
	DaoPool      = "daoPool"
)

type (
	Pool struct {
		Address          solana.PublicKey
		APY              float64
		Epoch            uint64
		SolanaStake      uint64
		TotalTokenSupply uint64
		TotalLamports    uint64
		UnstakeLiquidity uint64
		DepositFee       float64
		WithdrawalFee    float64
		RewardsFee       float64
		Validators       []PoolValidator
	}
	PoolValidator struct {
		ActiveStake uint64
		NodePK      solana.PublicKey
		VotePK      solana.PublicKey
	}
)
