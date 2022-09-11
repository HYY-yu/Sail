package api

import (
	"github.com/HYY-yu/seckill.pkg/core"
)

func (s *Server) Route(c *Handlers, engine core.Engine) {
	v1Group := engine.Group("/v1")
	{
		{
			g := v1Group.Group("/project_group")
			g.Use(core.WrapAuthHandler(s.HTTPMiddles.Jwt))
			g.GET("/list", c.projectGroupHandler.List)
			g.POST("/add", c.projectGroupHandler.Add)
			g.POST("/edit", c.projectGroupHandler.Edit)
			g.POST("/del", c.projectGroupHandler.Del)
		}

		{
			g := v1Group.Group("/staff")
			g.Use(core.WrapAuthHandler(s.HTTPMiddles.Jwt))
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

		v1Group.POST("/login", c.loginHandler.Login)
		v1Group.POST("/login/refresh", c.loginHandler.RefreshToken)
	}
}
