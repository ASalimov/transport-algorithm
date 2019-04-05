package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func generateFactsHandle(c *gin.Context) {
	n, err := strconv.Atoi(c.Query("n"))
	if err != nil {
		return
	}
	m, err := strconv.Atoi(c.Query("m"))
	if err != nil {
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
		fmt.Println("err = ", err)
	}
	c.Writer.Write(b)
	return
}

func findHandle(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("internal error: %v", err))
		return
	}

	var facts Facts
	if err := facts.UnmarshalBinary(b); err != nil {
		fmt.Println("err = ", err)
	}
	t, err := NewFactsWrapper(facts)
	if err != nil {
		fmt.Println("err = ", err)
	}
	t.Find()
	di := t.decision
	for i := range di.Volume {
		for j := range di.Volume[i] {
			if di.Volume[i][j] == empty {
				di.Volume[i][j] = 0
			}
		}
	}
	br, err := di.MarshalBinary()
	if err != nil {
		fmt.Println("err = ", err)
	}
	c.Writer.Write(br)
	return
}
