package httpserv

import (
	"fmt"
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/docs"
	"github.com/everstake/solana-pools/internal/delivery/httpserv/tools"
	v1 "github.com/everstake/solana-pools/internal/delivery/httpserv/v1"
	"github.com/everstake/solana-pools/internal/services"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"time"
)

type (
	API struct {
		cfg config.Env
		svc services.Service
		log *zap.Logger
		v1  *v1.Handler
	}
)

func NewAPI(cfg config.Env, svc services.Service, log *zap.Logger) (api *API, err error) {
	return &API{
		cfg: cfg,
		svc: svc,
		log: log,
		v1:  v1.New(svc, log),
	}, nil
}

func (api *API) Run() error {
	logger := ginzap.Ginzap(api.log, time.RFC3339, true)
	router := gin.New()
	router.Use(logger)
	docs.SwaggerInfo.BasePath = "/v1"
	v1g := router.Group("/v1")
	v1g.GET("/pools", tools.Must(api.v1.GetPools))
	v1g.GET("/pool/:name", tools.WSMust(api.v1.GetPool, time.Second*30))
	v1g.GET("/pool-statistic", tools.Must(api.v1.GetPoolsStatistic))
	v1g.GET("/pools-statistic", tools.WSMust(api.v1.GetTotalPoolsStatistic, time.Second*30))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	api.log.Info("Starting API server", zap.Uint64("port", api.cfg.HttpPort))
	return router.Run(fmt.Sprintf(":%d", api.cfg.HttpPort))
}