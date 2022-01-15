package policy

type RemoteClient struct {
}

func NewRemoteClient(path string) (*RemoteClient, error) {
	return &RemoteClient{}, nil
}

func (x *RemoteClient) Eval(in interface{}, out interface{}) error {
	panic("not implemented") // TODO: Implement
}
