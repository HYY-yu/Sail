package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndexHandler struct {
}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func (h *IndexHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func (h *IndexHandler) Group(c *gin.Context) {
	c.HTML(http.StatusOK, "group.html", gin.H{})
}

func (h *IndexHandler) Staff(c *gin.Context) {
	c.HTML(http.StatusOK, "staff.html", gin.H{})
}
