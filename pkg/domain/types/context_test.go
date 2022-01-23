package types_test

import (
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestContextTimeout(t *testing.T) {
	ctx := types.NewContext()
	timeout, cancel := ctx.WithTimeout(time.Millisecond * 100)
	defer cancel()

	var count int
	done := make(chan *struct{})
	go func() {
		for count < 1000000 {
			count++
			time.Sleep(time.Millisecond)
		}
		close(done)
	}()

	var exitByDone, exitByTimeout bool
	select {
	case <-done:
		exitByDone = true
	case <-timeout.Done():
		exitByTimeout = true
	}

	assert.False(t, exitByDone)
	assert.True(t, exitByTimeout)
	assert.Greater(t, count, 0)
}

func TestContextLogger(t *testing.T) {
	log := utils.Logger.Log()
	ctx := types.NewContext(types.WithLogger(log))
	assert.Equal(t, log, ctx.Log())
}
