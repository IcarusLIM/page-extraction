package main

import (
	extract "github.com/Ghamster0/page-extraction/src/extractor"
	"github.com/gin-gonic/gin"
)

type Request struct {
	Url      string            `json:"url" binding:"required"`
	Html     string            `json:"html" binding:"required"`
	Template *extract.Template `json:"template" binding:"required"`
}

func main() {
	r := gin.Default()
	r.POST("/content", func(c *gin.Context) {
		var req *Request = &Request{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(400, gin.H{
				"err": "Invalid Request body",
			})
		}
		c.JSON(200, extract.Extract(req.Url, req.Html, req.Template))
	})
	r.Run()
}
