package controller

import (
	"blog/internal/domain/blog"
	routing "github.com/qiangxue/fasthttp-routing"
)

type blogController struct {
	router  *routing.Router
	service *blog.Service
}

func NewBlogController(router *routing.Router, service *blog.Service) *blogController {
	return &blogController{
		router:  router,
		service: service,
	}
}
