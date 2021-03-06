package postgres

import (
	"github.com/dfuse-io/solana-go"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) GetValidatorByVotePK(key solana.PublicKey) (*dmodels.ValidatorView, error) {
	validator := &dmodels.ValidatorView{}
	err := db.Table("public.material_validator_data_view as validators").Where("id = ?", key.String()).First(validator).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return validator, err
}

func (db *DB) GetValidator(validatorID string, epoch uint64) (*dmodels.ValidatorView, error) {
	validator := &dmodels.ValidatorView{}
	DB := db.Table("public.validator_view_current_data as validators")
	if epoch == 10 {
		DB = db.Table("public.material_validator_data_view as validators")
	}
	err := DB.Where("id = ?", validatorID).First(validator).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return validator, err
}

func (db *DB) GetValidators(condition *ValidatorCondition, epoch uint64) ([]*dmodels.ValidatorView, error) {
	validators := make([]*dmodels.ValidatorView, 0)
	DB := db.Table("public.validator_view_current_data as validators")
	if epoch == 10 {
		DB = db.Table("public.material_validator_data_view as validators")
	}
	return validators, withValidatorCondition(DB, condition).Select("validators.*").Find(&validators).Error
}

func (db *DB) GetValidatorCount(condition *ValidatorCondition, epoch uint64) (int64, error) {
	i := int64(0)
	DB := db.Table("public.validator_view_current_data as validators")
	if epoch == 10 {
		DB = db.Table("public.material_validator_data_view as validators")
	}
	return i, withValidatorCondition(DB, condition).Count(&i).Error
}

func withValidatorCondition(db *gorm.DB, condition *ValidatorCondition) *gorm.DB {
	if condition == nil {
		return db
	}

	if condition.Condition != nil {
		if condition.Condition.Name != "" {
			db = db.Where(`validators.name ilike ?`, "%"+condition.Condition.Name+"%")
			condition.Condition.Name = ""
		}

		if len(condition.ValidatorIDs) > 0 {
			db = db.Where(`validators.id in (?)`, condition.ValidatorIDs)
		}
	}

	if len(condition.Epochs) > 0 {
		db = db.Where(`validators.epoch in (?)`, condition.Epochs)
	}

	if len(condition.PoolDataIDs) > 0 {
		db = db.Joins("INNER JOIN pool_validator_data pvd ON validators.id = pvd.validator_id")
		db = db.Where(`pvd.pool_data_id in (?)`, condition.PoolDataIDs)
	}

	db = withCond(db, condition.Condition)

	if condition.Sort != nil {
		return sortValidator(db, condition.Sort.ValidatorSort, condition.Sort.Desc)
	}

	return db
}

func sortValidator(db *gorm.DB, sort ValidatorSortType, desc bool) *gorm.DB {
	switch sort {
	case ValidatorAPY:
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "validators.apy",
					},
					Desc: desc,
				},
			},
		})
	case ValidatorStake:
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "validators.active_stake",
					},
					Desc: desc,
				},
			},
		})
	case ValidatorFee:
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "validators.fee",
					},
					Desc: desc,
				},
			},
		})
	case ValidatorScore:
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "validators.score",
					},
					Desc: desc,
				},
			},
		})
	case ValidatorSkippedSlot:
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "validators.skipped_slots",
					},
					Desc: desc,
				},
			},
		})
	case ValidatorDataCenter:
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "validators.data_center",
					},
					Desc: desc,
				},
			},
		})
	case StakingAccounts:
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "validators.staking_accounts",
					},
					Desc: desc,
				},
			},
		})
	}
	return db
}

func (db *DB) UpdateValidators(validators ...*dmodels.Validator) error {
	return db.Save(&validators).Error
}

func (db *DB) UpdateValidatorsData(data ...*dmodels.ValidatorData) error {
	return db.Save(&data).Error
}
