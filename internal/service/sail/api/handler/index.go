package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/HYY-yu/sail/internal/service/sail/api/svc"
	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type IndexHandler struct {
	projectGroupSvc *svc.ProjectGroupSvc
	namespaceSvc    *svc.NamespaceSvc
	configSvc       *svc.ConfigSvc
}

func NewIndexHandler(
	projectGroupSvc *svc.ProjectGroupSvc,
	namespaceSvc *svc.NamespaceSvc,
	configSvc *svc.ConfigSvc,
) *IndexHandler {
	return &IndexHandler{
		projectGroupSvc: projectGroupSvc,
		namespaceSvc:    namespaceSvc,
		configSvc:       configSvc,
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

func (h *IndexHandler) Project(c *gin.Context) {
	projectGroups := h.projectGroupSvc.SimpleList()

	c.HTML(http.StatusOK, "project.html", gin.H{
		"PGArr": projectGroups,
	})
}

func (h *IndexHandler) ProjectAdd(c *gin.Context) {
	projectGroups := h.projectGroupSvc.SimpleList()

	c.HTML(http.StatusOK, "project_add.html", gin.H{
		"PGArr": projectGroups,
	})
}

func (h *IndexHandler) ProjectEdit(c *gin.Context) {
	name := c.Query("name")
	id := c.Query("id")
	projectGroups := h.projectGroupSvc.SimpleList()

	c.HTML(http.StatusOK, "project_edit.html", gin.H{
		"Name":  name,
		"ID":    id,
		"PGArr": projectGroups,
	})
}

func (h *IndexHandler) Config(c *gin.Context) {
	c.HTML(http.StatusOK, "config.html", gin.H{})
}

func (h *IndexHandler) ConfigAdd(c *gin.Context) {
	projectId := c.Query("project_id")
	projectGroupId := c.Query("project_group_id")

	NSArr := h.namespaceSvc.SimpleList(gconv.Int(projectGroupId))
	ac := model.ConfigTypeYaml.AllConfigType()
	ConfigNSMap := h.configSvc.SimplePublicTree(gconv.Int(projectGroupId))

	c.HTML(http.StatusOK, "config_add.html", gin.H{
		"projectId":      projectId,
		"projectGroupId": projectGroupId,
		"NSArr":          NSArr,
		"AConfigType":    ac,
		"ConfigNSMap":    ConfigNSMap,
	})
}

func (h *IndexHandler) ConfigMeta(c *gin.Context) {
	projectId := c.Query("project_id")
	projectGroupId := c.Query("project_group_id")
	NSArr := h.namespaceSvc.SimpleList(gconv.Int(projectGroupId))

	c.HTML(http.StatusOK, "config_meta.html", gin.H{
		"projectId":      projectId,
		"projectGroupId": projectGroupId,
		"NSArr":          NSArr,
	})
}

func (h *IndexHandler) Public(c *gin.Context) {
	projectGroups := h.projectGroupSvc.SimpleList()

	c.HTML(http.StatusOK, "public.html", gin.H{
		"PGArr": projectGroups,
	})
}

func (h *IndexHandler) PublicAdd(c *gin.Context) {
	projectGroupId := c.Query("project_group_id")
	NSArr := h.namespaceSvc.SimpleList(gconv.Int(projectGroupId))
	ac := model.ConfigTypeYaml.AllConfigType()

	c.HTML(http.StatusOK, "public_add.html", gin.H{
		"projectGroupId": projectGroupId,
		"NSArr":          NSArr,
		"AConfigType":    ac,
	})
}

func (h *IndexHandler) PublicHistory(c *gin.Context) {
	ConfigID := c.Query("config_id")

	c.HTML(http.StatusOK, "public_history.html", gin.H{
		"configID": ConfigID,
	})
}

func (h *IndexHandler) PublicHistoryDiff(c *gin.Context) {
	ConfigID := c.Query("config_id")
	reversion := c.Query("reversion")

	c.HTML(http.StatusOK, "public_diff.html", gin.H{
		"configID":  ConfigID,
		"reversion": reversion,
	})
}
