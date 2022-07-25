package api

import (
	"github.com/HYY-yu/seckill.pkg/core"
)

func (s *Server) Route(c *Handlers, engine core.Engine) {
	v1Group := engine.Group("/v1")
	{
		v1Group.Use(core.WrapAuthHandler(s.HTTPMiddles.Jwt))

	}
}