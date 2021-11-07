package alertchain

var NewAlert = newAlert

var NewDefault = newDefault

func (x *Chain) Jobs() Jobs {
	return x.jobs
}
