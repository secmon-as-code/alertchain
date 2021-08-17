package alertchain

type Chain struct {
	Stages    []*Stage
	Arbitrary []Task
}

type Stage struct {
	Tasks []Task
}

func (x *Chain) NewStage() *Stage {
	stage := &Stage{}
	x.Stages = append(x.Stages, stage)
	return stage
}

func (x *Stage) AddTask(task Task) {
	x.Tasks = append(x.Tasks, task)
}
