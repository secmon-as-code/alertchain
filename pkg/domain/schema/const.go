package schema

const (
	edgeAlertToAttrs = "attributes"
	edgeAttrsToAnn   = "annotations"
	edgeAlertToRef   = "references"
	edgeAlertToJob   = "jobs"

	edgeJobToActionLog = "action_logs"
	edgeJobToAlert     = "alert"

	edgeActionLogToJob = "job"
	edgeAttrToAlert    = "alert"
	edgeAnnToAttr      = "attribute"
)
