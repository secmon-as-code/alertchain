package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

func getAlerts(c *gin.Context) {
	uc := ctxUsecase(c)

	alerts, err := uc.GetAlerts()
	if err != nil {
		c.Error(err)
		return
	}

	resp(c, http.StatusOK, alerts)
}

func getAlert(c *gin.Context) {
	alertID := c.Param("id")
	uc := ctxUsecase(c)

	alert, err := uc.GetAlert(types.AlertID(alertID))
	if err != nil {
		c.Error(err)
		return
	}

	resp(c, http.StatusOK, alert)
}

func postAlert(c *gin.Context) {
	var alert ent.Alert

	if err := c.BindJSON(&alert); err != nil {
		c.Error(types.ErrInvalidInput.Wrap(err))
		return
	}

	uc := ctxUsecase(c)
	newAlert, err := uc.RecvAlert(&alert)
	if err != nil {
		c.Error(err)
		return
	}

	resp(c, http.StatusCreated, newAlert)
}
