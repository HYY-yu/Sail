package api

import (
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/gin-gonic/gin"
)

func (s *Server) RouteHTML(c *Handlers, staticEngine *gin.Engine) {
	templateGroup := staticEngine.Group("/ui")

	templateGroup.GET("/login", c.loginHandler.LoginHTML)
	templateGroup.GET("/index", c.indexHandler.Index)
	templateGroup.GET("/group", c.indexHandler.Group)
	templateGroup.GET("/group/add", c.indexHandler.GroupAdd)
	templateGroup.GET("/group/edit", c.indexHandler.GroupEdit)
	templateGroup.GET("/staff/add", c.indexHandler.StaffAdd)
	templateGroup.GET("/staff/edit", c.indexHandler.StaffEdit)
	templateGroup.GET("/staff/grant", c.indexHandler.StaffGrant)
	templateGroup.GET("/staff/del_grant", c.indexHandler.StaffDelGrant)
	templateGroup.GET("/staff", c.indexHandler.Staff)
	templateGroup.GET("/namespace", c.indexHandler.Namespace)
	templateGroup.GET("/namespace/add", c.indexHandler.NamespaceAdd)
	templateGroup.GET("/namespace/edit", c.indexHandler.NamespaceEdit)
	templateGroup.GET("/project", c.indexHandler.Project)
	templateGroup.GET("/project/add", c.indexHandler.ProjectAdd)
	templateGroup.GET("/project/edit", c.indexHandler.ProjectEdit)
	templateGroup.GET("/config", c.indexHandler.Config)
	templateGroup.GET("/config/add", c.indexHandler.ConfigAdd)
	templateGroup.GET("/config/meta", c.indexHandler.ConfigMeta)
	templateGroup.GET("/public", c.indexHandler.Public)
	templateGroup.GET("/public/add", c.indexHandler.PublicAdd)
	templateGroup.GET("/public/history", c.indexHandler.PublicHistory)
	templateGroup.GET("/public/history_diff", c.indexHandler.PublicHistoryDiff)

}

func (s *Server) Route(c *Handlers, engine core.Engine) {
	v1Group := engine.Group("/v1")
	{
		{
			g := v1Group.Group("/project_group")
			g.Use(core.WrapAuthHandler(s.HTTPMiddles.Jwt))
			g.Use(c.staffHandler.MiddlewareStaffGroup)
			g.GET("/list", c.projectGroupHandler.List)
			g.POST("/add", c.projectGroupHandler.Add)
			g.POST("/edit", c.projectGroupHandler.Edit)
			g.POST("/del", c.projectGroupHandler.Del)
		}

		{
			g := v1Group.Group("/project")
			g.Use(core.WrapAuthHandler(s.HTTPMiddles.Jwt))
			g.Use(c.staffHandler.MiddlewareStaffGroup)
			g.GET("/list", c.projectHandler.List)
			g.POST("/add", c.projectHandler.Add)
			g.POST("/edit", c.projectHandler.Edit)
			g.POST("/del", c.projectHandler.Del)
		}

		{
			g := v1Group.Group("/namespace")
			g.Use(core.WrapAuthHandler(s.HTTPMiddles.Jwt))
			g.Use(c.staffHandler.MiddlewareStaffGroup)
			g.GET("/list", c.namespaceHandler.List)
			g.POST("/add", c.namespaceHandler.Add)
			g.POST("/edit", c.namespaceHandler.Edit)
			g.POST("/del", c.namespaceHandler.Del)
		}

		{
			g := v1Group.Group("/config")
			g.Use(core.WrapAuthHandler(s.HTTPMiddles.Jwt))
			g.Use(c.staffHandler.MiddlewareStaffGroup)
			g.GET("/tree", c.configHandler.Tree)
			g.GET("/meta", c.configHandler.MetaConfig)
			g.GET("/info", c.configHandler.Info)
			g.GET("/history", c.configHandler.History)
			g.GET("/history_info", c.configHandler.HistoryInfo)
			g.POST("/rollback", c.configHandler.Rollback)
			g.POST("/add", c.configHandler.Add)
			g.POST("/edit", c.configHandler.Edit)
			g.POST("/del", c.configHandler.Del)
			g.POST("/copy", c.configHandler.Copy)
		}

		{
			g := v1Group.Group("/staff")
			g.Use(core.WrapAuthHandler(s.HTTPMiddles.Jwt))
			g.Use(c.staffHandler.MiddlewareStaffGroup)
			g.GET("/list", c.staffHandler.List)
			g.POST("/add", c.staffHandler.Add)
			g.POST("/edit", c.staffHandler.Edit)
			g.POST("/del", c.staffHandler.Del)
			g.POST("/grant", c.staffHandler.Grant)
			g.POST("/del_grant", c.staffHandler.DelGrant)
		}

		{
			g := v1Group.Group("/login")
			g.Use(core.WrapAuthHandler(s.HTTPMiddles.Jwt))

			g.GET("/login_out", c.loginHandler.LoginOut)
			g.POST("/new_pass", c.loginHandler.ChangePassword)
		}

		{
			g := v1Group.Group("/env_test")
			g.Use(core.WrapAuthHandler(s.HTTPMiddles.Jwt))

			g.GET("/create", c.envTestHandler.CreateData)
			g.GET("/clean", c.envTestHandler.CleanData)
		}

		v1Group.POST("/login", c.loginHandler.Login)
		v1Group.POST("/login/refresh", c.loginHandler.RefreshToken)
	}
}
