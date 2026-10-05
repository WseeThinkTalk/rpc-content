package breaker

import (
	"context"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
)

// ExecuteWithFallback 使用 go-zero Google SRE 弹性熔断器保护受保护资源的调用
// 当 action 发生连续超时、错误或熔断器开启时，快速执行 fallback 降级策略
func ExecuteWithFallback[T any](
	ctx context.Context,
	name string,
	action func() (T, error),
	fallback func(err error) (T, error),
) (T, error) {
	var result T

	b := breaker.GetBreaker(name)
	err := b.DoCtx(ctx, func() error {
		res, err := action()
		if err != nil {
			return err
		}
		result = res
		return nil
	})

	if err != nil {
		logx.WithContext(ctx).Slowf("[Breaker:%s] action failed or breaker tripped: %v, triggering fallback", name, err)
		if fallback != nil {
			return fallback(err)
		}
		var zero T
		return zero, err
	}

	return result, nil
}
