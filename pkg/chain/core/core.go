package core

import (
	"os"
	"strings"
	"time"

	"github.com/m-mizutani/alertchain/pkg/action"
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/memory"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
)

type Core struct {
	alertPolicy  *policy.Client
	actionPolicy *policy.Client
	dbClient     interfaces.Database
	timeout      time.Duration

	scenarioLogger interfaces.ScenarioLogger
	actionMock     interfaces.ActionMock
	actionMap      map[types.ActionName]interfaces.RunAction

	disableAction bool
	enablePrint   bool
	maxSequences  int

	now func() time.Time
	env interfaces.Env
}

func New(options ...Option) *Core {
	c := &Core{
		dbClient:       memory.New(),
		timeout:        5 * time.Minute,
		actionMap:      action.Map(),
		scenarioLogger: &dummyScenarioLogger{},
		maxSequences:   types.DefaultMaxSequences,
		now:            time.Now,
		env:            Env,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func (x *Core) DBClient() interfaces.Database             { return x.dbClient }
func (x *Core) Timeout() time.Duration                    { return x.timeout }
func (x *Core) ScenarioLogger() interfaces.ScenarioLogger { return x.scenarioLogger }
func (x *Core) ActionMock() interfaces.ActionMock         { return x.actionMock }

func (x *Core) DisableAction() bool { return x.disableAction }
func (x *Core) MaxSequences() int   { return x.maxSequences }
func (x *Core) Now() time.Time      { return x.now() }
func (x *Core) Env() types.EnvVars  { return x.env() }

type Option func(c *Core)

func WithPolicyAlert(p *policy.Client) Option {
	return func(c *Core) {
		c.alertPolicy = p
	}
}

func WithPolicyAction(p *policy.Client) Option {
	return func(c *Core) {
		c.actionPolicy = p
	}
}

func WithDisableAction() Option {
	return func(c *Core) {
		c.disableAction = true
	}
}

func WithEnablePrint() Option {
	return func(c *Core) {
		c.enablePrint = true
	}
}

func WithExtraAction(name types.ActionName, action interfaces.RunAction) Option {
	return func(c *Core) {
		if _, ok := c.actionMap[name]; ok {
			panic("action name is already registered: " + name)
		}
		c.actionMap[name] = action
	}
}

func WithActionMock(mock interfaces.ActionMock) Option {
	return func(c *Core) {
		c.actionMock = mock
	}
}

func WithScenarioLogger(logger interfaces.ScenarioLogger) Option {
	return func(c *Core) {
		c.scenarioLogger = logger
	}
}

func WithEnv(f interfaces.Env) Option {
	return func(c *Core) {
		c.env = f
	}
}

func WithDatabase(db interfaces.Database) Option {
	return func(c *Core) {
		c.dbClient = db
	}
}

func Env() types.EnvVars {
	vars := types.EnvVars{}
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		vars[types.EnvVarName(pair[0])] = types.EnvVarValue(pair[1])
	}
	return vars
}
