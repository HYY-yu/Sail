package api

import (
	"github.com/HYY-yu/seckill.pkg/core"
)

func (s *Server) Route(c *Handlers, engine core.Engine) {
	v1Group := engine.Group("/v1")
	{
		//v1Group.Use(core.WrapAuthHandler(s.HTTPMiddles.Jwt))
		{
			projectGroupGroup := v1Group.Group("/project_group")
			projectGroupGroup.GET("/list", c.projectGroupHandler.List)
			projectGroupGroup.POST("/add", c.projectGroupHandler.Add)
			projectGroupGroup.POST("/edit", c.projectGroupHandler.Edit)
			projectGroupGroup.POST("/del", c.projectGroupHandler.Del)
		}

		{
			projectGroupGroup := v1Group.Group("/staff")
			projectGroupGroup.GET("/list", c.staffHandler.List)
			projectGroupGroup.POST("/add", c.staffHandler.Add)
			projectGroupGroup.POST("/edit", c.staffHandler.Edit)
			projectGroupGroup.POST("/del", c.staffHandler.Del)
			projectGroupGroup.POST("/grant", c.staffHandler.Grant)
			projectGroupGroup.POST("/del_grant", c.staffHandler.DelGrant)
		}
	}
}
