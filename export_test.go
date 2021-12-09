package alertchain

var (
	NewAlert     = newAlert
	NewAPIEngine = newAPIEngine
)

func (x *Chain) Jobs() Jobs {
	return x.jobs
}
