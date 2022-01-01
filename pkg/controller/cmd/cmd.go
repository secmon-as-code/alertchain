package cmd

type Cmd struct{}

func New() *Cmd {
	return &Cmd{}
}

func (x *Cmd) Run(argv []string) error {
	return nil
}
