package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func generateFactsHandle(c *gin.Context) {
	n, err := strconv.Atoi(c.Query("n"))
	if err != nil {
		errorHandle(http.StatusBadRequest, c, err)
		return
	}
	m, err := strconv.Atoi(c.Query("m"))
	if err != nil {
		errorHandle(http.StatusBadRequest, c, err)
		return
	}
	costs := make([][]float64, n)
	demands := make([]float64, m)
	supplies := make([]float64, n)
	for i := 0; i < n; i++ {
		costs[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			costs[i][j] = 1 + float64(rand.Intn(8))
		}
	}
	for i := 0; i < m; i++ {
		demands[i] = (1 + float64(rand.Intn(8))) * 10
	}
	for i := 0; i < n; i++ {
		supplies[i] = (1 + float64(rand.Intn(8))) * 10
	}

	facts := Facts{Costs: costs, Demands: demands, Supplies: supplies}
	b, err := facts.MarshalBinary()
	if err != nil {
		errorHandle(http.StatusInternalServerError, c, err)
		return
	}
	c.Writer.Write(b)
	return
}

func findHandle(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errorHandle(http.StatusBadRequest, c, err)
		return
	}

	var facts Facts
	if err := facts.UnmarshalBinary(b); err != nil {
		errorHandle(http.StatusBadRequest, c, err)
		return
	}
	t, err := NewFactsWrapper(facts)
	if err != nil {
		errorHandle(http.StatusBadRequest, c, err)
		return
	}
	t.Find()
	fmt.Println("Optimal dicision:")
	fmt.Println(t)
	fmt.Println()
	di := t.decision
	br, err := di.MarshalBinary()
	if err != nil {
		errorHandle(http.StatusInternalServerError, c, err)
		return
	}
	c.Writer.Write(br)
	return
}

func corezoid(c *gin.Context) {
	c.String(200, "ok")
	return
}

func errorHandle(code int, c *gin.Context, err error) {
	Error.Printf("%v", err)
	c.String(code, fmt.Sprintf("%v", err))
}
