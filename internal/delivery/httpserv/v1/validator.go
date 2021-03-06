package v1

import (
	"errors"
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/delivery/httpserv/tools"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// GetPoolValidators godoc
// @Summary RestAPI
// @Schemes
// @Description This list with pool's validators.
// @Tags validatorData
// @Param pname path string true "Name of the pool with strict observance of the case." default(Eversol)
// @Param vname query string false "The name of the validatorData without strict observance of the case."
// @Param epoch query number false "Epoch aggregation." Enums(1, 10) default(10)
// @Param sort query string false "sort param" Enums(apy, pool stake, stake, fee, score, skipped slot, data center) default(apy)
// @Param desc query bool false "desc" default(true)
// @Param offset query number true "offset for aggregation" default(0)
// @Param limit query number true "limit for aggregation" default(10)
// @Accept json
// @Produce json
// @Success 200 {object} tools.ResponseArrayData{data=[]validatorData} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /pool-validators/{pname} [get]
func (h *Handler) GetPoolValidators(ctx *gin.Context) (interface{}, error) {
	name := ctx.Param("pname")
	q := struct {
		Name   string `form:"vname"`
		Epoch  uint64 `form:"epoch,default=10"`
		Sort   string `form:"sort,default=apy"`
		Desc   bool   `form:"desc,default=true"`
		Offset uint64 `form:"offset,default=0"`
		Limit  uint64 `form:"limit,default=10"`
	}{}
	if err := ctx.ShouldBind(&q); err != nil {
		return nil, tools.NewStatus(http.StatusBadRequest, err)
	}

	resp, amount, err := h.svc.GetPoolValidators(name, q.Name, q.Sort, q.Desc, q.Epoch, q.Limit, q.Offset)
	if err != nil {
		h.log.Error("API GetPoolData", zap.Error(err))
		if errors.Is(err, postgres.ErrorRecordNotFounded) {
			return nil, tools.NewStatus(http.StatusBadRequest, fmt.Errorf("%s pool not found", name))
		}
		return nil, tools.NewStatus(http.StatusInternalServerError, err)
	}

	arr := make([]*validatorData, len(resp))
	for i, v := range resp {
		arr[i] = (&validatorData{}).Set(v)
	}

	return tools.ResponseArrayData{
		Data: arr,
		MetaData: &tools.MetaData{
			Offset:      q.Offset,
			Limit:       q.Limit,
			TotalAmount: amount,
		},
	}, nil
}

// GetAllValidators godoc
// @Summary RestAPI
// @Schemes
// @Description This list with all Solana's validators.
// @Tags validatorData
// @Param name query string false "The name of the validatorData without strict observance of the case."
// @Param epoch query number false "Epoch aggregation." Enums(1, 10) default(10)
// @Param epochs query []number false "Epochs for filter."
// @Param sort query string false "sort param" Enums(apy, stake, fee, score, skipped slot, data center, staking accounts) default(apy)
// @Param desc query bool false "desc" default(true)
// @Param offset query number true "offset for aggregation" default(0)
// @Param limit query number true "limit for aggregation" default(10)
// @Accept json
// @Produce json
// @Success 200 {object} tools.ResponseArrayData{data=[]validator} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /validators [get]
func (h *Handler) GetAllValidators(ctx *gin.Context) (interface{}, error) {
	q := struct {
		Name   string   `form:"name"`
		Epoch  uint64   `form:"epoch,default=10"`
		Epochs []uint64 `form:"epochs"`
		Sort   string   `form:"sort,default=apy"`
		Desc   bool     `form:"desc,default=true"`
		Offset uint64   `form:"offset,default=0"`
		Limit  uint64   `form:"limit,default=10"`
	}{}
	if err := ctx.ShouldBind(&q); err != nil {
		return nil, tools.NewStatus(http.StatusBadRequest, err)
	}

	resp, amount, err := h.svc.GetAllValidators(q.Name, q.Sort, q.Desc, q.Epoch, q.Epochs, q.Limit, q.Offset)
	if err != nil {
		return nil, tools.NewStatus(http.StatusInternalServerError, err)
	}

	arr := make([]*validator, len(resp))
	for i, v := range resp {
		arr[i] = (&validator{}).Set(v)
	}

	return tools.ResponseArrayData{
		Data: arr,
		MetaData: &tools.MetaData{
			Offset:      q.Offset,
			Limit:       q.Limit,
			TotalAmount: amount,
		},
	}, nil
}

type validator struct {
	Name             string  `json:"name"`
	Delinquent       bool    `json:"delinquent"`
	Image            string  `json:"image"`
	NodePK           string  `json:"node_pk"`
	APY              float64 `json:"apy"`
	VotePK           string  `json:"vote_pk"`
	TotalActiveStake float64 `json:"total_active_stake"`
	StakingAccounts  uint64  `json:"staking_accounts"`
	Fee              float64 `json:"fee"`
	Score            int64   `json:"score"`
	SkippedSlots     float64 `json:"skipped_slots"`
	DataCenter       string  `json:"data_center"`
	Epoch            uint64  `json:"epoch"`
}

func (v *validator) Set(validator *smodels.Validator) *validator {
	v.NodePK = validator.NodePK
	v.Name = validator.Name
	v.Delinquent = validator.Delinquent
	v.Image = validator.Image
	v.APY, _ = validator.APY.Float64()
	v.VotePK = validator.VotePK
	v.TotalActiveStake, _ = validator.TotalActiveStake.Float64()
	v.StakingAccounts = validator.StakingAccounts
	v.Fee, _ = validator.Fee.Float64()
	v.Score = validator.Score
	v.SkippedSlots, _ = validator.SkippedSlots.Float64()
	v.DataCenter = validator.DataCenter
	v.Epoch = validator.Epoch
	return v
}

type validatorData struct {
	Name             string  `json:"name"`
	Image            string  `json:"image"`
	NodePK           string  `json:"node_pk"`
	APY              float64 `json:"apy"`
	VotePK           string  `json:"vote_pk"`
	PoolActiveStake  float64 `json:"pool_active_stake"`
	TotalActiveStake float64 `json:"total_active_stake"`
	StakingAccounts  uint64  `json:"staking_accounts"`
	Fee              float64 `json:"fee"`
	Score            int64   `json:"score"`
	SkippedSlots     float64 `json:"skipped_slots"`
	DataCenter       string  `json:"data_center"`
}

func (v *validatorData) Set(validator *smodels.PoolValidatorData) *validatorData {
	v.NodePK = validator.NodePK
	v.Name = validator.Name
	v.Image = validator.Image
	v.APY, _ = validator.APY.Float64()
	v.VotePK = validator.VotePK
	v.PoolActiveStake, _ = validator.PoolActiveStake.Float64()
	v.TotalActiveStake, _ = validator.TotalActiveStake.Float64()
	v.StakingAccounts = validator.StakingAccounts
	v.Fee, _ = validator.Fee.Float64()
	v.Score = validator.Score
	v.SkippedSlots, _ = validator.SkippedSlots.Float64()
	v.DataCenter = validator.DataCenter

	return v
}
