/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// corsMiddleware solves cors
func corsMiddleware() gin.HandlerFunc {
	mwCORS := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 24 * time.Hour,
	})
	return mwCORS
}

// ginLogger defines gin logger
func ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// begin time
		start := time.Now()
		// handle request
		c.Next()
		// end time
		end := time.Now()
		latency := end.Sub(start)

		path := c.Request.URL.Path

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		logger.Infof("[LEDGERDATA REFINER]| %3d | %13v | %15s | %s  \"%s\"",
			statusCode,
			latency,
			clientIP,
			method, path,
		)
	}
}
