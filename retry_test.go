package gretry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	r := New(WithRecovery(), WithBaseDelay(100*time.Millisecond))
	err := r.EnsureRetryTimes(10, func() error {
		fmt.Println(time.Now())
		return Retriable(errors.New("haha"))
	})
	assert.ErrorContains(t, err, "haha")
}

func TestBackoff(t *testing.T) {
	bo := &Backoff{
		MinDelay: time.Millisecond * 10,
		MaxDelay: time.Second,
		Factor:   2,
	}
	r := New(WithRecovery(), WithBaseDelay(100*time.Millisecond), WithBackoff(bo))
	err := r.EnsureRetryTimes(10, func() error {
		fmt.Println(time.Now())
		return Retriable(errors.New("haha"))
	})
	assert.ErrorContains(t, err, "haha")
}

func TestPanic(t *testing.T) {
	r := New(WithRecovery())
	err := r.Ensure(func() error {
		panic("haha")
		return nil
	})
	assert.ErrorContains(t, err, "haha")
}

func TestContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r := New(WithCtx(ctx))
	err := r.Ensure(func() error {
		t.Log(time.Now())
		return RetriableMesg("haha")
	})
	assert.Equal(t, err, ctx.Err())
}

func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second*2, func() {
		cancel()
	})

	r := New(WithCtx(ctx))
	err := r.Ensure(func() error {
		t.Log(time.Now())
		return RetriableMesg("haha")
	})
	assert.Equal(t, err, ctx.Err())
}
