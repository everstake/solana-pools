package v1

import (
	"github.com/everstake/solana-pools/internal/delivery/httpserv/tools"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetCoins godoc
// @Summary RestAPI
// @Schemes
// @Description The information on tokens with the specified search parameters.
// @Tags coin
// @Accept json
// @Produce json
// @Param offset query number true "offset for aggregation" default(0)
// @Param limit query number true "limit for aggregation" default(10)
// @Param name query string false "coin name"
// @Success 200 {object} tools.ResponseArrayData{data=[]coin} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /coins [get]
func (h *Handler) GetCoins(ctx *gin.Context) (interface{}, error) {
	q := struct {
		Name   string `form:"name"`
		Limit  uint64 `form:"limit,default=0"`
		Offset uint64 `form:"offset,default=10"`
	}{}
	if err := ctx.ShouldBind(&q); err != nil {
		return nil, tools.NewStatus(http.StatusBadRequest, err)
	}

	scoins, count, err := h.svc.GetCoins(q.Name, q.Limit, q.Offset)
	if err != nil {
		return nil, tools.NewStatus(http.StatusInternalServerError, err)
	}

	coins := make([]*coin, len(scoins))
	for i, c := range scoins {
		coins[i] = (&coin{}).Set(c)
	}

	return tools.ResponseArrayData{
		Data: coins,
		MetaData: &tools.MetaData{
			Offset:      q.Offset,
			Limit:       q.Limit,
			TotalAmount: count,
		},
	}, nil
}

// GetPoolsCoins godoc
// @Summary RestAPI
// @Schemes
// @Description The information about pool tokens with the specified search parameters.
// @Tags coin
// @Accept json
// @Produce json
// @Param name query string false "coin name"
// @Param sort query string false "sort param" Enums(price, name) default(price)
// @Param desc query bool false "desc" default(true)
// @Param offset query number true "offset for aggregation" default(0)
// @Param limit query number true "limit for aggregation" default(10)
// @Success 200 {object} tools.ResponseArrayData{data=[]coin} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /pool-coins [get]
func (h *Handler) GetPoolsCoins(ctx *gin.Context) (interface{}, error) {
	q := struct {
		Name   string `form:"name"`
		Sort   string `form:"sort,default=price"`
		Desc   bool   `form:"desc,default=true"`
		Offset uint64 `form:"offset,default=0"`
		Limit  uint64 `form:"limit,default=10"`
	}{}
	if err := ctx.ShouldBind(&q); err != nil {
		return nil, tools.NewStatus(http.StatusBadRequest, err)
	}

	scoins, count, err := h.svc.GetPoolCoins(q.Name, q.Sort, q.Desc, q.Limit, q.Offset)
	if err != nil {
		return nil, tools.NewStatus(http.StatusInternalServerError, err)
	}

	coins := make([]*coin, len(scoins))
	for i, c := range scoins {
		coins[i] = (&coin{}).Set(c)
	}

	return tools.ResponseArrayData{
		Data: coins,
		MetaData: &tools.MetaData{
			Offset:      q.Offset,
			Limit:       q.Limit,
			TotalAmount: count,
		},
	}, nil
}

type liquidityPool struct {
	Name  string `json:"name"`
	About string `json:"about"`
	Image string `json:"image"`
	URL   string `json:"url"`
}

func (lp *liquidityPool) Set(pool *smodels.LiquidityPool) *liquidityPool {
	lp.Name = pool.Name
	lp.About = pool.About
	lp.URL = pool.URL
	lp.Image = pool.Image
	return lp
}

type deFi struct {
	BuyCoin       *coin          `json:"buy_coin"`
	LiquidityPool *liquidityPool `json:"liquidity_pool"`
	Liquidity     float64        `json:"liquidity"`
	APY           float64        `json:"apy"`
}

func (f *deFi) Set(defi *smodels.DeFi, buyCoin *coin, liquidityPool *liquidityPool) *deFi {
	f.APY, _ = defi.APY.Float64()
	f.LiquidityPool = liquidityPool
	f.BuyCoin = buyCoin
	f.Liquidity = defi.Liquidity
	return f
}

type coin struct {
	Name       string             `json:"name"`
	Address    string             `json:"address"`
	USD        float64            `json:"usd"`
	ThumbImage string             `json:"thumb_image"`
	SmallImage string             `json:"small_image"`
	LargeImage string             `json:"large_image"`
	DeFi       map[string][]*deFi `json:"de_fi,omitempty"`
}

func (c *coin) Set(coinM *smodels.Coin) *coin {
	c.USD = coinM.USD
	c.ThumbImage = coinM.ThumbImage
	c.SmallImage = coinM.SmallImage
	c.LargeImage = coinM.LargeImage
	c.Name = coinM.Name
	c.Address = coinM.Address
	if coinM.DeFi != nil {
		c.DeFi = make(map[string][]*deFi)
		for _, fi := range coinM.DeFi {
			_, ok := c.DeFi[fi.BuyCoin.Name]
			if !ok {
				defi := make([]*deFi, 0)
				defi = append(defi, (&deFi{}).Set(fi, (&coin{}).Set(fi.BuyCoin), (&liquidityPool{}).Set(fi.LiquidityPool)))
				c.DeFi[fi.BuyCoin.Name] = defi
				continue
			}

			c.DeFi[fi.BuyCoin.Name] = append(c.DeFi[fi.BuyCoin.Name], (&deFi{}).Set(fi, (&coin{}).Set(fi.BuyCoin), (&liquidityPool{}).Set(fi.LiquidityPool)))
		}
	}
	return c
}
