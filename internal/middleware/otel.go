package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Tracing is a lightweight request-tracing middleware.
// In production, replace the log.Printf calls with your OpenTelemetry exporter.
func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate a simple trace ID (replace with otel.Tracer in production)
		traceID := fmt.Sprintf("%d", start.UnixNano())
		c.Set("trace_id", traceID)
		c.Header("X-Trace-Id", traceID)

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// Structured log — wire up to your OTEL collector here
		fmt.Printf("[trace] id=%s method=%s path=%s status=%d latency=%s ip=%s\n",
			traceID,
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			latency.String(),
			c.ClientIP(),
		)
	}
}
