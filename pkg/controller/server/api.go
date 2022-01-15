package model

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/zlog"
)

type apiServer struct {
	addr string

	engine *gin.Engine
}

func newAPIServer(addr string, db *db.Client, fallback http.Handler, logger *zlog.Logger) *apiServer {
	return &apiServer{
		addr:   addr,
		engine: newAPIEngine(db, fallback, logger),
	}
}

func (x *apiServer) Run() error {
	if err := x.engine.Run(x.addr); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}

type apiResponse struct {
	Message string `json:"message,omitempty"`
}

const (
	ctxKeyDB     = "db"
	ctxKeyLogger = "logger"
)

func ctxSetDB(c *gin.Context, db *db.Client) {
	c.Set(ctxKeyDB, db)
}

func ctxGetDB(c *gin.Context) *db.Client {
	obj, ok := c.Get(ctxKeyDB)
	if !ok {
		panic("DB is not found in gin.Context")
	}

	db, ok := obj.(*db.Client)
	if !ok {
		panic("DB object in gin.Context is not *db.Client")
	}

	return db
}

func ctxSetLogger(c *gin.Context, logger *zlog.Logger) {
	c.Set(ctxKeyLogger, logger)
}

func ctxGetLogger(c *gin.Context) *zlog.Logger {
	obj, ok := c.Get(ctxKeyLogger)
	if !ok {
		panic("Logger is not found in gin.Context")
	}

	logger, ok := obj.(*zlog.Logger)
	if !ok {
		panic("Logger object in gin.Context is not zlog.Logger")
	}

	return logger
}

func newAPIEngine(db *db.Client, fallback http.Handler, logger *zlog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.NoRoute(func(c *gin.Context) {
		if fallback != nil {
			fallback.ServeHTTP(c.Writer, c.Request)
		} else {
			c.JSON(http.StatusNotFound, apiResponse{Message: "not found"})
		}
	})
	engine.Use(func(c *gin.Context) {
		ctxSetLogger(c, logger)
		ctxSetDB(c, db)
		c.Next()
	})
	api := engine.Group("/api/v1")
	api.GET("/alert", getAlerts)
	api.GET("/alert/:id", getAlert)

	api.POST("/alert", func(c *gin.Context) {
		c.JSON(200, struct{}{})
	})

	return engine
}

func getAlert(c *gin.Context) {
	id := types.AlertID(c.Param("id"))
	ctx := types.NewContextWith(c, ctxGetLogger(c))

	resp, err := ctxGetDB(c).GetAlert(ctx, id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if resp == nil {
		c.JSON(http.StatusNotFound, apiResponse{Message: "not found"})
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func queryToInt(c *gin.Context, name string, defaultValue int) int {
	if q := c.Query(name); q != "" {
		v, err := strconv.ParseUint(q, 10, 64)
		if err != nil {
			_ = c.Error(err)
			return -1
		}
		return int(v)
	}

	return defaultValue
}

func getAlerts(c *gin.Context) {
	ctx := types.NewContextWith(c, ctxGetLogger(c))

	offset := queryToInt(c, "offset", 0)
	limit := queryToInt(c, "limit", 10)
	if offset < 0 || limit < 0 {
		return
	}

	resp, err := ctxGetDB(c).GetAlerts(ctx, offset, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
