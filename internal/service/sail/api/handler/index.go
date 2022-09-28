package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/HYY-yu/sail/internal/service/sail/api/svc"
	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type IndexHandler struct {
	projectGroupSvc *svc.ProjectGroupSvc
}

func NewIndexHandler(
	projectGroupSvc *svc.ProjectGroupSvc,
) *IndexHandler {
	return &IndexHandler{
		projectGroupSvc: projectGroupSvc,
	}
}

func (h *IndexHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func (h *IndexHandler) Group(c *gin.Context) {
	c.HTML(http.StatusOK, "group.html", gin.H{})
}

func (h *IndexHandler) GroupAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "group_add.html", gin.H{})
}

func (h *IndexHandler) GroupEdit(c *gin.Context) {
	name := c.Query("name")
	id := c.Query("id")

	c.HTML(http.StatusOK, "group_edit.html", gin.H{"Name": name, "ID": id})
}

func (h *IndexHandler) StaffAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "staff_add.html", gin.H{})
}

func (h *IndexHandler) StaffEdit(c *gin.Context) {
	name := c.Query("name")
	id := c.Query("id")

	c.HTML(http.StatusOK, "staff_edit.html", gin.H{"Name": name, "ID": id})
}

func (h *IndexHandler) StaffDelGrant(c *gin.Context) {
	name := c.Query("name")
	id := c.Query("id")

	c.HTML(http.StatusOK, "staff_del_grant.html", gin.H{"Name": name, "ID": id})
}

func (h *IndexHandler) StaffGrant(c *gin.Context) {
	id := c.Query("id")
	name := c.Query("name")

	projectGroups := h.projectGroupSvc.SimpleList()

	roles := model.RoleAdmin.AllRole()

	var rds []struct {
		ID   int
		Name string
	}
	for i := range roles {
		rds = append(rds, struct {
			ID   int
			Name string
		}{ID: int(roles[i]), Name: roles[i].String()})
	}

	c.HTML(http.StatusOK, "staff_grant.html", gin.H{
		"Name":  name,
		"ID":    id,
		"PGArr": projectGroups,
		"RDS":   rds,
	})
}

func (h *IndexHandler) Staff(c *gin.Context) {
	c.HTML(http.StatusOK, "staff.html", gin.H{})
}

func (h *IndexHandler) Namespace(c *gin.Context) {
	projectGroups := h.projectGroupSvc.SimpleList()

	c.HTML(http.StatusOK, "namespace.html", gin.H{
		"PGArr": projectGroups,
	})
}

func (h *IndexHandler) NamespaceEdit(c *gin.Context) {
	name := c.Query("name")
	id := c.Query("id")
	realTime := c.Query("real_time")
	rb, _ := strconv.ParseBool(realTime)
	check1, check2 := "", ""
	if rb {
		check1 = "checked"
	} else {
		check2 = "checked"
	}

	c.HTML(
		http.StatusOK,
		"namespace_edit.html",
		gin.H{
			"Name":   name,
			"ID":     id,
			"Check1": check1,
			"Check2": check2,
		},
	)
}

func (h *IndexHandler) NamespaceAdd(c *gin.Context) {
	projectGroups := h.projectGroupSvc.SimpleList()

	c.HTML(http.StatusOK, "namespace_add.html", gin.H{
		"PGArr": projectGroups,
	})
}
