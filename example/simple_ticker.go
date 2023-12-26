package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/Speshl/go-sbus"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(ctx)

	sBus, _ := sbus.NewSBus(
		"/dev/ttyAMA0",
		true,
		true,
		nil,
	)

	group.Go(func() error {
		defer cancel()
		return sBus.Start(ctx)
	})

	group.Go(func() error {
		logTicker := time.NewTicker(1 * time.Second)
		finishTicker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-logTicker.C:
				rxFrame := sBus.GetReadFrame()
				sBus.SetWriteFrame(rxFrame)
				slog.Info("echoing frame", "frame", rxFrame)
			case <-finishTicker.C:
				cancel()
			}
		}
	})

	err := group.Wait()
	if err != nil {
		slog.Error("failed echoing sbus", "error", err)
	}
}
