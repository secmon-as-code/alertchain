package policy

type LocalClient struct {
}

func NewLocalClient(path string) (*LocalClient, error) {
	return &LocalClient{}, nil
}

func (x *LocalClient) Eval(in interface{}, out interface{}) error {
	panic("not implemented") // TODO: Implement
}
