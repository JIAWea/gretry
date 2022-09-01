package gretry

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestBase(t *testing.T) {
	r := New(WithRecovery(), WithBaseDelay(100*time.Millisecond))
	err := r.EnsureRetryTimes(10, func() error {
		fmt.Println(time.Now())
		return CreateRetryErr(errors.New("haha"))
	})

	t.Log(err)
}

func TestBackoff(t *testing.T) {
	bo := &Backoff{
		MinDelay: time.Second,
		MaxDelay: time.Second * 5,
		Factor:   2,
	}
	r := New(WithRecovery(), WithBaseDelay(time.Second), WithBackoff(bo))
	err := r.EnsureRetryTimes(5, func() error {
		fmt.Println(time.Now())
		return CreateRetryErr(errors.New("haha"))
	})
	t.Log(err)
}

func TestPanic(t *testing.T) {
	r := New(WithRecovery())
	err := r.Ensure(func() error {
		panic("haha")
		return nil
	})
	t.Log(err)
}

func TestContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	r := New(WithCtx(ctx), WithBaseDelay(time.Second))
	err := r.Ensure(func() error {
		t.Log(time.Now())
		return CreateRetryErrMsg("haha")
	})

	if !reflect.DeepEqual(err, ctx.Err()) {
		t.Errorf("got: %v, expect: %v", err, ctx.Err())
	}
}

func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second*2, func() {
		cancel()
	})

	r := New(WithCtx(ctx))
	err := r.Ensure(func() error {
		t.Log(time.Now())
		return CreateRetryErrMsg("haha")
	})

	if !reflect.DeepEqual(err, ctx.Err()) {
		t.Errorf("got: %v, expect: %v", err, ctx.Err())
	}
}
