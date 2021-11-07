package alertchain

import (
	"os"

	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/zlog"
)

type Option func(chain *Chain)

func OptLogLevel(logLevel string) Option {
	return func(chain *Chain) {
		if err := chain.logger.SetLogLevel(logLevel); err != nil {
			chain.logger.With("given", logLevel).Warn("ignore invalid log level")
		}
	}
}

func OptLogFormat(format string) Option {
	return func(chain *Chain) {
		switch format {
		case "console":
			chain.logger.Emitter = zlog.NewWriterWith(zlog.NewConsoleFormatter(), os.Stdout)
		case "json":
			chain.logger.Emitter = zlog.NewWriterWith(zlog.NewJsonFormatter(), os.Stdout)
		default:
			chain.logger.With("given", format).Warn("ignore invalid logger format")
		}

	}
}

func OptDBConfig(dbType, dbConfig string) Option {
	return func(chain *Chain) {
		chain.config.DB = types.DBConfig{
			Type:   dbType,
			Config: dbConfig,
		}
	}
}

func OptDB(db db.Interface) Option {
	return func(chain *Chain) {
		chain.db = db
	}
}

func OptJobs(jobs ...*Job) Option {
	return func(chain *Chain) {
		chain.jobs = jobs
	}
}

func OptSources(src ...Source) Option {
	return func(chain *Chain) {
		chain.sources = src
	}
}

func OptAction(actions ...Action) Option {
	return func(chain *Chain) {
		chain.actions = actions
	}
}
