package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type bodyLogReader struct {
	io.ReadCloser
	body *bytes.Buffer
}

func (r *bodyLogReader) Read(p []byte) (n int, err error) {
	n, err = r.ReadCloser.Read(p)
	r.body.Write(p)
	return
}

func setupRouter() *gin.Engine {

	router := gin.New()

	router.Use(Logger())
	router.StaticFile("/", "static/index.html")
	router.Static("/static", "static")
	router.GET("/api/generate", generateFactsHandle)
	router.POST("/api/find", findHandle)
	return router

}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		blr := &bodyLogReader{body: bytes.NewBufferString(""), ReadCloser: c.Request.Body}
		c.Request.Body = blr
		c.Set("rid", rand.Uint64())
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		t := time.Now()

		c.Next()

		latency := time.Since(t)
		fmt.Printf("%s[%d] - req: %s rps: %s\n", c.Request.RequestURI, latency.Nanoseconds()/1e6, blr.body, blw.body)
	}
}
