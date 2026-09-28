package main

import (
	_ "github.com/udistrital/autenticacion_mid/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	apistatus "github.com/udistrital/utils_oas/apiStatusLib"
	auditoria "github.com/udistrital/utils_oas/auditoria"
	"github.com/udistrital/utils_oas/customerror"
	"github.com/udistrital/utils_oas/security"
	"github.com/udistrital/utils_oas/xray"
)

func main() {
	allowedOrigins := []string{"*.udistrital.edu.co"}
	if beego.BConfig.RunMode == beego.DEV {
		allowedOrigins = []string{"*"}
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"User-Agent",
			"X-Amzn-Trace-Id"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	apistatus.Init()
	auditoria.InitMiddleware()
	security.SetSecurityHeaders()
	xray.Init()

	beego.ErrorController(&customerror.CustomErrorController{})
	beego.Run()
}
