package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

type Request struct {
	MaxContextLength int     `json:"max_context_length"`
	MaxLength        int     `json:"max_length"`
	Prompt           string  `json:"prompt"`
	Quiet            bool    `json:"quiet"`
	RepPen           float64 `json:"rep_pen"`
	RepPenRange      int     `json:"rep_pen_range"`
	RepPenSlope      float64 `json:"rep_pen_slope"`
	Temperature      float64 `json:"temperature"`
	Tfs              int     `json:"tfs"`
	TopA             int     `json:"top_a"`
	TopK             int     `json:"top_k"`
	TopP             float64 `json:"top_p"`
	Typical          int     `json:"typical"`
}

func chat_response(c *gin.Context) {
	reqs, _ := c.GetQuery("response")

	fmt.Println("Request: ", c.Request.Body)
	// fmt.Println("Response: ", resp)

	url := "https://herbal-pmc-allowance-cognitive.trycloudflare.com/api/v1/generate" // replace with your API endpoint

	data := Request{
		MaxContextLength: 2048,
		MaxLength:        100,
		Prompt:           reqs,
		Quiet:            false,
		RepPen:           1.1,
		RepPenRange:      256,
		RepPenSlope:      1,
		Temperature:      0.5,
		Tfs:              1,
		TopA:             0,
		TopK:             100,
		TopP:             0.9,
		Typical:          1,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	spew.Dump(string(body))
	c.JSON(http.StatusOK, gin.H{
		"response": string(body),
	})
}

func main() {
	router := gin.Default()
	router.POST("/chat", chat_response)
	router.Run(":8085")

}
