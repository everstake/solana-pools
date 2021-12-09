package services

import (
	"errors"
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/cache"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/services/smodels"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"sync"
	"time"
)

func (s *Imp) GetPool(name string) (*smodels.PoolDetails, error) {
	pd, err := s.cache.GetPool(name)
	if err != nil && !errors.Is(err, cache.KeyWasNotFound) {
		return nil, err
	}
	if pd != nil {
		return pd, nil
	}

	dPool, err := s.dao.GetPool(name)
	if err != nil {
		return nil, fmt.Errorf("dao.GetPool: %s", err.Error())
	}
	dLastPoolData, err := s.dao.GetLastPoolData(dPool.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("dao.GetPoolData: %s", err.Error())
	}
	dValidators, err := s.dao.GetPoolValidatorData(dPool.ID)
	if err != nil {
		return nil, fmt.Errorf("dao.GetPoolValidatorData: %s", err.Error())
	}
	validatorsS := make([]*smodels.Validator, len(dValidators))
	validatorsD := make([]*dmodels.Validator, len(dValidators))
	for i, v := range dValidators {
		validatorsD[i], err = s.dao.GetValidator(v.ValidatorID)
		if err != nil {
			return nil, fmt.Errorf("dao.GetValidator(%s): %w", err)
		}
		validatorsS[i] = (&smodels.Validator{}).Set(v.ActiveStake, validatorsD[i])
	}

	Pool := (&smodels.Pool{}).Set(dLastPoolData, &dPool, validatorsD)

	pd = &smodels.PoolDetails{
		Pool:       *Pool,
		Validators: validatorsS,
	}

	s.cache.SetPool(pd, time.Second*30)

	return pd, nil
}

func (s *Imp) GetPools(name string, limit uint64, offset uint64) ([]*smodels.PoolDetails, error) {
	dPools, err := s.dao.GetPools(&postgres.Condition{
		Name: name,
		Pagination: postgres.Pagination{
			Limit:  limit,
			Offset: offset,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("dao.GetPool: %w", err)
	}
	if len(dPools) == 0 {
		return nil, nil
	}
	pools := make([]*smodels.PoolDetails, len(dPools))
	for i, v1 := range dPools {
		pools[i] = &smodels.PoolDetails{
			Pool: smodels.Pool{
				Address: v1.Address,
				Name:    v1.Name,
			},
		}

		dLastPoolData, err := s.dao.GetLastPoolData(v1.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("dao.GetLastPoolData: %w", err)
		}

		dValidators, err := s.dao.GetPoolValidatorData(dLastPoolData.ID)
		if err != nil {
			return nil, fmt.Errorf("dao.GetPoolValidatorData: %w", err)
		}

		validatorsS := make([]*smodels.Validator, len(dValidators))
		validatorsD := make([]*dmodels.Validator, len(dValidators))
		for i, v2 := range dValidators {
			validatorsD[i], err = s.dao.GetValidator(v2.ValidatorID)
			if err != nil {
				return nil, fmt.Errorf("dao.GetValidator: %w", err)
			}
			validatorsS[i] = (&smodels.Validator{}).Set(v2.ActiveStake, validatorsD[i])
		}

		pools[i].Set(dLastPoolData, &v1, validatorsD)

		pools[i].Validators = validatorsS
	}

	return pools, nil
}

func (s *Imp) GetPoolsCurrentStatistic() (*smodels.Statistic, error) {
	stat, err := s.cache.GetCurrentStatistic()
	if err != nil && !errors.Is(err, cache.KeyWasNotFound) {
		return nil, err
	}
	if stat != nil {
		return stat, nil
	}

	dPools, err := s.dao.GetPools(&postgres.Condition{Network: postgres.MainNet})
	if err != nil {
		return nil, fmt.Errorf("dao.GetPool: %w", err)
	}
	if len(dPools) == 0 {
		return nil, nil
	}
	stat = &smodels.Statistic{}

	once := sync.Once{}
	pools := make([]*smodels.PoolDetails, len(dPools))
	var ActiveStakeSum, UnstakeSum uint64
	for i, v1 := range dPools {
		pools[i] = &smodels.PoolDetails{
			Pool: smodels.Pool{
				Address: v1.Address,
				Name:    v1.Name,
			},
		}

		dLastPoolData, err := s.dao.GetLastPoolData(v1.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("dao.GetLastPoolData: %w", err)
		}

		dValidators, err := s.dao.GetPoolValidatorData(dLastPoolData.ID)
		if err != nil {
			return nil, fmt.Errorf("dao.GetValidators: %w", err)
		}

		validatorsD := make([]*dmodels.Validator, len(dValidators))
		for i, v2 := range dValidators {
			validatorsD[i], err = s.dao.GetValidator(v2.ValidatorID)
			if err != nil {
				return nil, fmt.Errorf("dao.GetValidator: %w", err)
			}
		}

		pools[i].Set(dLastPoolData, &v1, validatorsD)

		once.Do(func() {
			stat.MINScore = pools[i].AVGScore
			stat.MAXScore = pools[i].AVGScore
		})

		if pools[i].AVGScore > stat.MAXScore {
			stat.MAXScore = pools[i].AVGScore
		}
		if pools[i].AVGScore < stat.MINScore {
			stat.MINScore = pools[i].AVGScore
		}

		ActiveStakeSum += dLastPoolData.ActiveStake
		UnstakeSum += dLastPoolData.UnstakeLiquidity
		stat.AVGSkippedSlots = stat.AVGSkippedSlots.Add(pools[i].AVGSkippedSlots)
		stat.AVGScore += pools[i].AVGScore
		stat.Delinquent = stat.Delinquent.Add(pools[i].Delinquent)
	}

	stat.ActiveStake.SetLamports(ActiveStakeSum)
	stat.UnstakeLiquidity.SetLamports(UnstakeSum)
	stat.AVGSkippedSlots = stat.AVGSkippedSlots.Div(decimal.NewFromInt(int64(len(dPools))))
	stat.AVGScore /= int64(len(dPools))

	s.cache.SetCurrentStatistic(stat, time.Second*30)

	return stat, nil
}

func (s *Imp) GetPoolsStatistic(name string, aggregate string, from time.Time, to time.Time) ([]*smodels.Pool, error) {
	pool, err := s.dao.GetPool(name)
	if err != nil {
		return nil, err
	}

	a, err := s.dao.GetPoolStatistic(pool.ID, postgres.SearchAggregate(aggregate), from, to)
	if err != nil {
		return nil, err
	}

	data := make([]*smodels.Pool, len(a))
	for i, v := range a {
		data[i] = (&smodels.Pool{}).Set(v, &pool, nil)
		data[i].ValidatorCount, err = s.dao.GetValidatorCount(&postgres.PoolValidatorDataCondition{
			PoolDataIDs: []uuid.UUID{
				v.ID,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (s *Imp) GetPoolCount() (int64, error) {
	i, err := s.dao.GetPoolCount(nil)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (s *Imp) GetNetworkAPY() (float64, error) {
	d, err := s.cache.GetAPY()
	if err != nil {
		return 0, err
	}

	f, _ := d.Float64()

	return f, nil
}

func (s *Imp) GetUSD() (float64, error) {
	d, err := s.cache.GetPrice()
	if err != nil {
		return 0, err
	}

	f, _ := d.Float64()

	return f, nil
}