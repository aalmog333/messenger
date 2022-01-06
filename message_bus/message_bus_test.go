package message_bus

import (
	"context"
	"github.com/riid/messenger/bus"
	"github.com/riid/messenger/envelope"
	"github.com/riid/messenger/middleware"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestMessageBus_Dispatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	e := envelope.FromMessage("test message")

	var b *messageBus
	handlerCalled := false
	m := middleware.HandleFunc(func(hCtx context.Context, hb bus.Bus, he envelope.Envelope) {
		handlerCalled = true
		assert.Same(t, ctx, hCtx)
		assert.Same(t, b, hb)
		assert.Same(t, e, he)
		<-time.After(1 * time.Second)
	})

	b = New(m, 1, 4)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		_ = b.Run(ctx)
		wg.Done()
	}()

	b.Dispatch(ctx, e)
	b.Dispatch(ctx, e)
	b.Dispatch(ctx, e)
	cancel()

	wg.Wait()

	assert.True(t, handlerCalled)
}

func TestIdentityNext_should_always_return_the_passed_envelope(t *testing.T) {
	ctx := context.Background()
	e := envelope.FromMessage("test message")
	res := identityNext(ctx, e)

	assert.Same(t, e, res)
}
