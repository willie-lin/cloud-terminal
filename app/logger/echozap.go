package logger

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func ZapLogger(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			req := c.Request()
			res := c.Response()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			fields := []zapcore.Field{
				zap.String("request_id", id),
				zap.String("remote_ip", c.RealIP()),
				zap.String("host", req.Host),
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.String("user_agent", req.UserAgent()),
			}

			err := next(c)

			fields = append(fields,
				zap.Int("status", res.Status),
				zap.Int64("size", res.Size),
				zap.String("time", time.Since(start).String()),
			)

			n := res.Status
			msg := "Request"

			switch {
			case err != nil:
				msg = "Request error"
				fields = append(fields, zap.Error(err))
				log.Error(msg, fields...)
				return err
			case n >= 500:
				log.Error(msg, fields...)
			case n >= 400:
				log.Warn(msg, fields...)
			case n >= 300:
				log.Info(msg, fields...)
			default:
				log.Info(msg, fields...)
			}

			return nil
		}
	}
}
