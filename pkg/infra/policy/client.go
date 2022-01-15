package policy

type Client interface {
	Eval(in interface{}, out interface{}) error
}
