package grill

import (
	"context"
	"fmt"
	"sync"
)

type LifeCycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func StartAll(ctx context.Context, grills ...LifeCycle) error {
	startFn := func(ctx context.Context, lc LifeCycle) error {
		return lc.Start(ctx)
	}
	return doAll(ctx, startFn, grills...)
}

func StopAll(ctx context.Context, grills ...LifeCycle) error {
	stopFn := func(ctx context.Context, lc LifeCycle) error {
		return lc.Stop(ctx)
	}
	return doAll(ctx, stopFn, grills...)
}

func doAll(ctx context.Context, fn func(ctx context.Context, lc LifeCycle) error, grills ...LifeCycle) error {
	wg := sync.WaitGroup{}
	wg.Add(len(grills))
	errChan := make(chan error, len(grills))
	for _, grill := range grills {
		go func(g LifeCycle, wg *sync.WaitGroup) {
			defer wg.Done()
			if err := fn(ctx, g); err != nil {
				errChan <- err
			}
		}(grill, &wg)
	}
	wg.Wait()

	var errors []string
	for err := range errChan {
		errors = append(errors, err.Error())
	}

	return fmt.Errorf("%v", errors)
}
