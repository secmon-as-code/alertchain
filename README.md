# AlertChain

AlertChain is simple & programmable SOAR (Security Orchestration, Automation and Response) platform.

## Concept

![Concept](https://user-images.githubusercontent.com/605953/130339742-4aba4f88-1b1d-4b48-8323-0dce5f8a85fc.jpg)

AlertChain has `Chain` consisting of multiple steps, `Stage`. Also Each `Stage` has one ore more `Task` to annotate and evaluate a security alert and execute workflow.

A security engineer can build own chain as a Go plugin and AlertChain import the plugin. AlertChain executes a `Chain` provided by the plugin step by step. A security engineer can not only use pre-defined Task but also write own original procedure.

## Installation

Build binary with `go >= 1.16` and `npm >= 7.18.1`.

```sh
% git clone https://github.com/m-mizutani/alertchain.git
% cd alertchain
% make
% cp alertchain /path/to/bin
```

## Usage

### 1) Create your workflow in your repository

- Workflow

```go
package main

import (
	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/types"
)

type myEvaluator struct{}

func (x *myEvaluator) Name() string                              { return "myEvaluator" }
func (x *myEvaluator) Execute(alert *alertchain.Alert) error {
	if alert.Title == "Something wrong" {
		alert.Severity = types.SevAffected
	}
	if err := alert.Commit(); err != nil {
		return err
	}
	return nil
}

func Chain() *alertchain.Chain {
	return &alertchain.Chain{
		Stages: []alertchain.Tasks{
			{&myEvaluator{}},
		},
	}
}
```

Then, compile it as plugin.

```sh
$ go build -buildmode=plugin -o mychain .
```

Finally, run `alertchain` with your plugin.

```sh
$ alertchain -c mychain
```
