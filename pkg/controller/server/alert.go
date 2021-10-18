package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

func postAlert(c *gin.Context) {
	var alert alertchain.Alert

	if err := c.BindJSON(&alert); err != nil {
		c.Error(types.ErrInvalidInput.Wrap(err))
		return
	}

	newAlert, err := ctxChain(c).Execute(c, &alert)
	if err != nil {
		c.Error(goerr.Wrap(err).With("alert", alert))
		return
	}

	resp(c, http.StatusCreated, newAlert)
}
