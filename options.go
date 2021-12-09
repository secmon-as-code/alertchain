package alertchain

import (
	"net/http"
	"os"
	"strings"

	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/zlog"
)

type Option func(chain *Chain) error

func WithLogLevel(logLevel string) Option {
	return func(chain *Chain) error {
		if err := chain.logger.SetLogLevel(logLevel); err != nil {
			chain.logger.With("given", logLevel).Warn("ignore invalid log level")
		}
		return nil
	}
}

func WithLogFormat(format string) Option {
	return func(chain *Chain) error {
		switch format {
		case "console":
			chain.logger.Emitter = zlog.NewWriterWith(zlog.NewConsoleFormatter(), os.Stdout)
		case "json":
			chain.logger.Emitter = zlog.NewWriterWith(zlog.NewJsonFormatter(), os.Stdout)
		default:
			chain.logger.With("given", format).Warn("ignore invalid logger format")
		}
		return nil
	}
}

func WithDBConfig(dbType, dbConfig string) Option {
	return func(chain *Chain) error {
		chain.config.DB = types.DBConfig{
			Type:   dbType,
			Config: dbConfig,
		}
		return nil
	}
}

func WithDB(db db.Interface) Option {
	return func(chain *Chain) error {
		chain.db = db
		return nil
	}
}

func WithAPI(addr, url string, fallback http.Handler) Option {
	return func(chain *Chain) error {
		chain.apiURL = strings.TrimRight(url, "/")
		chain.api = newAPIServer(addr, chain.db, fallback, chain.logger)
		return nil
	}
}

func WithJobs(jobs ...*Job) Option {
	return func(chain *Chain) error {
		chain.jobs = append(chain.jobs, jobs...)
		return nil
	}
}

func WithSources(src ...Source) Option {
	return func(chain *Chain) error {
		chain.sources = append(chain.sources, src...)
		return nil
	}
}

func WithAction(actions ...Action) Option {
	return func(chain *Chain) error {
		chain.actions = append(chain.actions, actions...)
		return nil
	}
}
