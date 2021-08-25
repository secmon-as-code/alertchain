package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

func getAlerts(c *gin.Context) {
	uc := ctxUsecase(c)

	alerts, err := uc.GetAlerts(c)
	if err != nil {
		c.Error(err)
		return
	}

	newAlerts := make([]*alertchain.Alert, len(alerts))
	for i, alert := range alerts {
		newAlerts[i] = alertchain.NewAlert(alert, nil)
	}
	resp(c, http.StatusOK, newAlerts)
}

func getAlert(c *gin.Context) {
	alertID := c.Param("id")
	uc := ctxUsecase(c)

	alert, err := uc.GetAlert(c, types.AlertID(alertID))
	if err != nil {
		c.Error(err)
		return
	}

	resp(c, http.StatusOK, alertchain.NewAlert(alert, nil))
}

func postAlert(c *gin.Context) {
	var alert alertchain.Alert

	if err := c.BindJSON(&alert); err != nil {
		c.Error(types.ErrInvalidInput.Wrap(err))
		return
	}

	uc := ctxUsecase(c)
	newAlert, err := uc.RecvAlert(c, &alert)
	if err != nil {
		c.Error(goerr.Wrap(err).With("alert", alert))
		return
	}

	resp(c, http.StatusCreated, newAlert)
}
