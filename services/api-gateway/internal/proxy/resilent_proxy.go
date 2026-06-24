package proxy

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/sony/gobreaker"
)

const (
	maxRetries     = 3
	baseRetryDelay = 100 * time.Millisecond
	timeout30s     = 30 * time.Second
	timeout60s     = 60 * time.Second
	retryIn        = 30
)

type resilientProxy struct {
	cb          *gobreaker.CircuitBreaker
	proxyFn     fiber.Handler
	serviceName string
}

func newFullpathProxy(targetURL, serviceName string) *resilientProxy {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        serviceName,
		MaxRequests: maxRetries,
		Timeout:     timeout60s,
		Interval:    timeout30s,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= maxRetries-1
		},
	})

	proxyHandler := proxy.Balancer(proxy.Config{
		Servers: []string{targetURL},
		ModifyRequest: func(c *fiber.Ctx) error {
			c.Request().Header.Add("X-Real-IP", c.IP())
			return nil
		},
		ModifyResponse: func(c *fiber.Ctx) error {
			c.Response().Header.Del(fiber.HeaderServer)
			return nil
		},
		Timeout: timeout60s,
	})

	return &resilientProxy{
		cb:          cb,
		proxyFn:     proxyHandler,
		serviceName: serviceName,
	}
}

func (rp *resilientProxy) fallbackResponse(c *fiber.Ctx, err error) error {
	status := fiber.StatusServiceUnavailable
	if errors.Is(err, gobreaker.ErrOpenState) {
		status = fiber.StatusTooManyRequests
	}

	return c.Status(status).JSON(fiber.Map{
		"error":    true,
		"msg":      "Service temporarily unavailable",
		"service":  rp.serviceName,
		"retry_in": retryIn,
		"fallback": true,
	})
}

func (rp *resilientProxy) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// для запроса с телом не делать retry
		if c.Request().Header.ContentLength() > 0 {
			_, err := rp.cb.Execute(func() (interface{}, error) {
				return nil, rp.proxyFn(c)
			})
			if err == nil {
				return nil
			}
			return rp.fallbackResponse(c, err)
		}

		// Для GET/HEAD делать retry
		var lastErr error
		for i := range maxRetries {
			_, err := rp.cb.Execute(func() (interface{}, error) {
				c.Locals("attempt", i+1)
				return nil, rp.proxyFn(c)
			})

			if err == nil {
				return nil
			}

			lastErr = err
			if i < maxRetries-1 {
				delay := baseRetryDelay * time.Duration(1+i*i) // при базе 100 = 100,200,500ms
				time.Sleep(delay)
			}
		}
		return rp.fallbackResponse(c, lastErr)
	}
}

// newStrippedProxy — удаляет /api/v1 (APIPrefix + APIVersion) из пути.
func newStrippedProxy(targetURL, serviceName string) *resilientProxy {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        serviceName,
		MaxRequests: maxRetries,
		Timeout:     timeout60s,
		Interval:    timeout30s,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= maxRetries-1
		},
	})

	proxyHandler := proxy.Balancer(proxy.Config{
		Servers: []string{targetURL},
		ModifyRequest: func(c *fiber.Ctx) error {
			c.Request().Header.Add("X-Real-IP", c.IP())

			originalPath := c.OriginalURL()
			stripPrefix := APIPrefix + APIVersion
			stripPrefixLen := len(stripPrefix)
			if len(originalPath) >= stripPrefixLen && originalPath[:7] == stripPrefix {
				newPath := originalPath[7:] // /api/v1/files -> /files
				if newPath == "" {
					newPath = "/"
				}
				c.Request().SetRequestURI(newPath)
			}
			return nil
		},
		ModifyResponse: func(c *fiber.Ctx) error {
			c.Response().Header.Del(fiber.HeaderServer)
			return nil
		},
		Timeout: timeout60s,
	})

	return &resilientProxy{
		cb:          cb,
		proxyFn:     proxyHandler,
		serviceName: serviceName,
	}
}
