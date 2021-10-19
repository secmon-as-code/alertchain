package alertchain

import (
	"context"
	"sync"
	"time"

	"github.com/m-mizutani/alertchain/pkg/utils"
)

type Handler func(ctx context.Context, alert *Alert) error

type Source interface {
	Name() string
	Run(handler Handler) error
}

func (x *Chain) StartSources() {
	wg := &sync.WaitGroup{}
	x.startSources(wg)
	wg.Wait()
}

func (x *Chain) StartSourcesAsync() {
	x.startSources(nil)
}

func (x *Chain) startSources(wg *sync.WaitGroup) {
	if err := x.diagnosis(); err != nil {
		logger.With("err", err).Error(err.Error())
		panic(err)
	}

	handler := func(ctx context.Context, alert *Alert) error {
		_, err := x.Execute(ctx, alert)
		return err
	}
	for i := range x.Sources {
		if wg != nil {
			wg.Add(1)
		}

		go func(src Source) {
			if wg != nil {
				defer wg.Done()
			}

			for {
				if err := src.Run(handler); err != nil {
					utils.HandleError(err)
				} else {
					break
				}

				time.Sleep(time.Second * 3)
			}
		}(x.Sources[i])
	}
}
