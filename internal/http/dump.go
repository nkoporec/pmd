package http

import (
	"encoding/json"

	"github.com/dgraph-io/ristretto"
	"github.com/gin-gonic/gin"
)

type RequestData struct {
	Payload   string `json:"payload"`
	Callstack   []CallstackData `json:"callstack"`
	File      string `json:"file"`
	Line      string `json:"line"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
}

type CallstackData struct {
	File   string `json:"file"`
	Line      string `json:"line"`
	Function      string `json:"function"`
}

func dump(c *gin.Context) {
	cache := c.MustGet("cache").(*ristretto.Cache)
	messages := c.MustGet("messages").(chan interface{})

	var data *RequestData

	// Request needs to have:
	// payload - The data that the client dumped.
	// callstack - The full callstack if applicable.
	// file - File path.
	// line - Line number.
	// type - Debug type (PHP, JS, GO...)
	// timestamp
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	// Validate that the correct data is send.
	requestIsValid := validateRequest(data)
	if !requestIsValid {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Wrong request data.",
		})
		return
	}

	// Get old data.
	var dumpData []*RequestData
	oldData, found := cache.Get("breakpoints")
	if found {
		dumpData = oldData.([]*RequestData)
	}

	dumpData = append(dumpData, data)
	cache.Set("breakpoints", dumpData, 1)
	messages <- dumpData
}

func validateRequest(data *RequestData) bool {
	// dumb way to do it, but it works.
	if len(data.Payload) <= 0 {
		return false
	}

	if len(data.File) <= 0 {
		return false
	}

	if len(data.Type) <= 0 {
		return false
	}

	return true
}
