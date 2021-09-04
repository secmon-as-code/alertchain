# AlertChain

AlertChain is simple & programmable SOAR (Security Orchestration, Automation and Response) platform.

## Concept

![Concept](https://user-images.githubusercontent.com/605953/132087986-318968c1-6e0a-419b-9013-be50c2b93930.jpg)

AlertChain has `Chain` consisting of multiple steps, `Stage`. Also Each `Stage` has one ore more `Task` to annotate and evaluate a security alert and execute workflow.

A security engineer can build own chain as a Go plugin and AlertChain import the plugin. AlertChain executes a `Chain` provided by the plugin step by step. A security engineer can use not only pre-defined Task but also original procedure written by own.

## Installation

Build binary with `go >= 1.16` and `npm >= 7.18.1`.

```sh
% git clone https://github.com/m-mizutani/alertchain.git
% cd alertchain
% make
% cp alertchain /path/to/bin
```

## Usage

### 0) Setup your own repository

```sh
$ mkdir your-chain
$ cd your-chain
$ go mod init your-chain
% go get github.com/m-mizutani/alertchain
```

Code examples are available in [./examples/simple](./examples/simple)

### 1) Create your tasks in your repository

At first, create your structure that has `Task` interface. `Task` is a minimum unit of automated workflow. Example tasks are in [task.go](./examples/simple/task.go). One of tasks is `Evaluator` to check if "Suspicious Login" alert comes from internal network or not in following.

```go
type Evaluator struct{}

func (x *Evaluator) Name() string { return "Evaluation" }

func (x *Evaluator) Execute(ctx *types.Context, alert *alertchain.Alert) error {
	if alert.Title == "Suspicious Login" {
		attrs := alert.Attributes.FindByKey("srcAddr").FindByType(types.AttrIPAddr)
		if len(attrs) != 1 {
			return nil // Attribute not found
		}

		addr := net.ParseIP(attrs[0].Value)
		_, internal, _ := net.ParseCIDR("10.1.0.0/16")
		if internal.Contains(addr) {
			alert.UpdateSeverity(types.SevSafe)
		}
	}

	return nil
}
```

### 2) Create your chain to describe workflow.

In [chain.go](./examples/simple/chain.go), prepare a function `Chain()` to setup and return `*alertchain.Chain` object.

`Chain` has a list of `Job`.  A set of `Task`. `Job` is executed one by one and `Task`s in `Job` are executed concurrently. Next Job is kicked after exiting all `Task` |

```go
func Chain() (*alertchain.Chain, error) {
	return &alertchain.Chain{
		Jobs: []*alertchain.Job{
			{
				Tasks: []alertchain.Task{&Evaluator{}},
			},
			{
				Tasks: []alertchain.Task{&CreateTicket{}}
			},
		},
	}, nil
}

```

### 3) Run with your plugin

Then, compile it as plugin.

```sh
$ go build -buildmode=plugin -o mychain.so .
```

Finally, run `alertchain` with your plugin.

```sh
$ alertchain -c mychain.so
```

After that, http://localhost:9080 will be opened.

# License

MIT License