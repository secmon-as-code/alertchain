package server

/*
func getAction(c *gin.Context) {
	uc := ctxUsecase(c)
	s := c.Param("id")
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		c.Error(goerr.Wrap(types.ErrInvalidInput.Wrap(err), "invalid action log id").With("id", s))
		return
	}

	actionLog, err := uc.GetActionLog(types.NewContextWith(c), int(id))
	if err != nil {
		c.Error(err)
		return
	}

	resp(c, http.StatusOK, actionLog)
}

type postActionBody struct {
	ActionID string `json:"action_id"`
	AttrID   int    `json:"attr_id"`
}

func postAction(c *gin.Context) {
	uc := ctxUsecase(c)
	var body postActionBody
	if err := c.Bind(&body); err != nil {
		c.Error(err)
		return
	}

	actionLog, err := uc.ExecuteAction(types.NewContextWith(c), body.ActionID, body.AttrID)
	if err != nil {
		c.Error(err)
		return
	}

	resp(c, http.StatusOK, actionLog)
}
*/
