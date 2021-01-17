package controllers

import (
	"github.com/revel/revel"
	"leanote/app/service"
	"strings"
)

type MinioR struct {
	BaseController
}

func (c MinioR) GetResource() revel.Result {
	mc, _ := service.GetMinioClient()
	object := c.Request.URL.Path
	url := mc.PresignedGetObject(strings.TrimLeft(object, "/"), 3600 * 2)
	return c.Redirect(url)
}