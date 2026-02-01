package retry

import (
	"context"
	"time"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RetryInterceptor struct {
	maxAttempts int
	baseDelay   time.Duration
}

func NewRetryInterceptor(maxAttempts int, baseDelay time.Duration) *RetryInterceptor {
	return &RetryInterceptor{
		maxAttempts: maxAttempts,
		baseDelay:   baseDelay,
	}
}

func (r *RetryInterceptor) Interceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var err error
	multiplier := 1
	for attempt := 1; attempt <= r.maxAttempts; attempt++ {
		err = invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			return nil
		}

		st, ok := status.FromError(err)
		if !ok {
			return err
		}

		if st.Code() != codes.Unavailable &&
			st.Code() != codes.ResourceExhausted &&
			st.Code() != codes.DeadlineExceeded &&
			st.Code() != codes.Internal &&
			st.Code() != codes.Aborted {
			return err
		}

		metrics.GatewayRetriesTotal.WithLabelValues(method, st.Code().String()).Inc()

		if attempt < r.maxAttempts {
			delay := r.baseDelay * time.Duration(multiplier)
			multiplier *= 2
			time.Sleep(delay)
		}
	}

	return err
}
