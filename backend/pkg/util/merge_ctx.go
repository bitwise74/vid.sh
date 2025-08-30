package util

import "context"

func MergeContexts(ctx1, ctx2 context.Context) (context.Context, context.CancelFunc) {
	merged, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-ctx1.Done():
			if ctx1.Err() != context.Canceled && ctx1.Err() != context.DeadlineExceeded {
				cancel()
			}
		case <-ctx2.Done():
			if ctx2.Err() != context.Canceled && ctx2.Err() != context.DeadlineExceeded {
				cancel()
			}
		case <-merged.Done():
		}
	}()

	return merged, cancel
}
