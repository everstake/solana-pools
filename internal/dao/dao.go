package dao

import (
	"fmt"
	"github.com/dfuse-io/solana-go"
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	uuid "github.com/satori/go.uuid"
	"time"
)

type (
	DAO interface {
		Postgres
	}
	Postgres interface {
		GetPool(name string) (dmodels.Pool, error)
		GetPoolCount(*postgres.Condition) (int64, error)
		GetLastPoolData(poolID uuid.UUID, condition *postgres.Condition) (*dmodels.PoolData, error)
		GetLastEpochPoolData(PoolID uuid.UUID, currentEpoch uint64) (*dmodels.PoolData, error)
		GetPools(*postgres.Condition) ([]dmodels.Pool, error)
		UpdatePoolData(*dmodels.PoolData) error
		GetPoolStatistic(poolID uuid.UUID, aggregate postgres.Aggregate, from time.Time, to time.Time) ([]*dmodels.PoolData, error)

		GetValidatorCount(condition *postgres.PoolValidatorDataCondition) (int64, error)
		GetValidatorByVotePK(key solana.PublicKey) (*dmodels.Validator, error)
		GetValidator(validatorID string) (*dmodels.Validator, error)
		GetValidators(condition *postgres.Condition) ([]*dmodels.Validator, error)
		UpdateValidators(validators ...*dmodels.Validator) error
		GetPoolValidatorData(poolDataID uuid.UUID) ([]*dmodels.PoolValidatorData, error)
		CreatePoolValidatorData(pools ...*dmodels.PoolValidatorData) error
		DeleteValidators(poolID uuid.UUID) error
	}
	Imp struct {
		*postgres.DB
	}
)

func NewDAO(cfg config.Env) (d DAO, err error) {
	p, err := postgres.NewDB(cfg.PostgresDSN)
	if err != nil {
		return d, fmt.Errorf("postgres.NewDB: %s", err.Error())
	}

	return &Imp{
		p,
	}, nil
}