package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

func getAlerts(c *gin.Context) {
	uc := ctxUsecase(c)

	alerts, err := uc.GetAlerts(types.WrapContext(c))
	if err != nil {
		c.Error(err)
		return
	}

	resp(c, http.StatusOK, alerts)
}

func getAlert(c *gin.Context) {
	alertID := c.Param("id")
	uc := ctxUsecase(c)

	ctx := types.WrapContext(c)
	alert, err := uc.GetAlert(ctx, types.AlertID(alertID))
	if err != nil {
		c.Error(err)
		return
	}

	resp(c, http.StatusOK, alert)
}

func postAlert(c *gin.Context) {
	var alert alertchain.Alert

	if err := c.BindJSON(&alert); err != nil {
		c.Error(types.ErrInvalidInput.Wrap(err))
		return
	}

	uc := ctxUsecase(c)
	attrs := make([]*ent.Attribute, len(alert.Attributes))
	for i := range alert.Attributes {
		attrs[i] = &alert.Attributes[i].Attribute
	}
	ctx := types.WrapContext(c)
	newAlert, err := uc.HandleAlert(ctx, &alert.Alert, attrs)
	if err != nil {
		c.Error(goerr.Wrap(err).With("alert", alert))
		return
	}

	resp(c, http.StatusCreated, newAlert)
}
