package api

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m-mizutani/alertchain/frontend"
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/pkg/errors"
)

const (
	contextUsecase      = "usecase"
	contextRequestIDKey = "requestID"
	cookieTokenName     = "token"
	cookieReferrerName  = "referrer"
)

type errorResponse struct {
	Error  string                 `json:"error"`
	Values map[string]interface{} `json:"values,omitempty"`
}

func respError(c *gin.Context, code int, err error) {
	c.JSON(code, errorResponse{Error: err.Error()})
}

func New(uc interfaces.Usecase) *gin.Engine {
	engine := gin.Default()

	engine.Use(func(c *gin.Context) {
		reqID := uuid.New().String()

		c.Set(contextRequestIDKey, reqID)
		c.Set(contextUsecase, uc)
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
	engine.GET("/bundle.js", getBundleJS)

	r := engine.Group("/api/v1")
	r.GET("/alert", getAlert)

	return engine
}

type cache struct {
	data []byte
	eTag string
}
type cacheMap map[string]*cache

var assetCache = cacheMap{}

func init() {
	assets := frontend.Assets()

	indexHTML, err := assets.ReadFile("dist/index.html")
	if err != nil {
		panic("Open dist/index.html: " + err.Error())
	}
	bundleJS, err := assets.ReadFile("dist/bundle.js")
	if err != nil {
		panic("Open dist/bundle.js: " + err.Error())
	}

	assetCache["index.html"] = &cache{
		data: indexHTML,
		eTag: fmt.Sprintf("%x", md5.Sum(indexHTML)),
	}

	assetCache["bundle.js"] = &cache{
		data: bundleJS,
		eTag: fmt.Sprintf("%x", md5.Sum(bundleJS)),
	}
}

func handleAsset(ctx *gin.Context, fname, contentType string) {
	c, ok := assetCache[fname]
	if !ok {
		respError(ctx, http.StatusNotFound, errors.New("Not found"))
		return
	}

	ctx.Header("Cache-Control", "public, max-age=31536000")
	ctx.Header("ETag", c.eTag)

	if match := ctx.GetHeader("If-None-Match"); match != "" {
		if strings.Contains(match, c.eTag) {
			ctx.Status(http.StatusNotModified)
			return
		}
	}

	ctx.Data(http.StatusOK, contentType, c.data)
}

// Assets
func getIndex(c *gin.Context) {
	handleAsset(c, "index.html", "text/html")
}

func getBundleJS(c *gin.Context) {
	handleAsset(c, "bundle.js", "application/javascript")
}
