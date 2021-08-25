package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m-mizutani/alertchain/pkg/usecase"
	"github.com/pkg/errors"
)

type Server struct {
	usecase usecase.Interface
	addr    string
	port    uint64
}

func New(uc usecase.Interface, addr string, port uint64) *Server {
	return &Server{
		usecase: uc,
		addr:    addr,
		port:    port,
	}
}

func (x *Server) Run() error {
	if mode, ok := os.LookupEnv("GIN_MODE"); ok {
		gin.SetMode(mode)
	} else {
		gin.SetMode("release") // Set release as default
	}

	engine := gin.Default()
	engine.Use(func(c *gin.Context) {
		reqID := uuid.New().String()

		c.Set(contextRequestIDKey, reqID)
		c.Set(contextUsecase, x.usecase)
		c.Next()
	})

	engine.Use(func(c *gin.Context) {
		c.Next()

		if ginError := c.Errors.Last(); ginError != nil {
			if err := errors.Cause(ginError); err != nil {
				respError(c, http.StatusInternalServerError, ginError)
			} else {
				respError(c, http.StatusInternalServerError, ginError)
			}
		}
	})

	engine.GET("/", getIndex)
	engine.GET("/alert", getIndex)
	engine.GET("/alert/:id", getIndex)
	engine.GET("/bundle.js", getBundleJS)

	r := engine.Group("/api/v1")
	r.GET("/alert", getAlerts)
	r.GET("/alert/:id", getAlert)
	r.POST("/alert", postAlert)

	if err := engine.Run(fmt.Sprintf("%s:%d", x.addr, x.port)); err != nil {
		return err
	}
	return nil
}

const (
	contextUsecase      = "usecase"
	contextRequestIDKey = "requestID"
	cookieTokenName     = "token"
	cookieReferrerName  = "referrer"
)

func ctxUsecase(c *gin.Context) usecase.Interface {
	val, ok := c.Get(contextUsecase)
	if !ok {
		panic("No usecase saved in gin.Context")
	}
	uc, ok := val.(usecase.Interface)
	if !ok {
		panic("Can not cast value in gin.Context to usecase.Interface")
	}
	return uc
}

type errorResponse struct {
	Error  string                 `json:"error"`
	Values map[string]interface{} `json:"values,omitempty"`
}

func respError(c *gin.Context, code int, err error) {
	c.JSON(code, errorResponse{Error: err.Error()})
}

type dataResponse struct {
	Data interface{} `json:"data"`
}

func resp(c *gin.Context, code int, data interface{}) {
	c.JSON(code, dataResponse{Data: data})
}
