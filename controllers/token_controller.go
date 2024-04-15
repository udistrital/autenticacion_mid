package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/autenticacion_mid/models"
	"github.com/udistrital/autenticacion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// TokenController operations for Token
type TokenController struct {
	beego.Controller
}

// URLMapping ...
func (c *TokenController) URLMapping() {
	c.Mapping("GetEmail", c.GetEmail)
	c.Mapping("GetRol", c.GetRol)
}

// GetEmail ...
// @Title GetEmail
// @Description Recibe el correo electrónico del usuario desde la autenticación
// @Param	body body 	models.Token	true		"The key for staticblock"
// @Success 200 {object} models.UserInfo
// @Failure 404 not found resource
// @router /emailToken [post]
func (c *TokenController) GetEmail() {
	defer errorhandler.HandlePanic(&c.Controller)

	var (
		v models.Token
	)

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		beego.Error(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 400, nil, err.Error())
	}

	if response, err := services.GetInfoByEmail(&v); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = response // ? No se usa el APIResponseDTO debido a que este endpoint se utiliza por muchos servicios en producción con el formato de respuesta antiguo
	} else {
		beego.Error(err)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// GetRol ...
// @Title GetRol
// @Description Recibe el usuario y devuelve información detallada del usuario
// @Param	body	body 	models.UserName  	true		"Usuario registrado en wso2"
// @Success 200 {object} models.Payload
// @Failure 404 not found resource
// @router /userRol [post]
func (c *TokenController) GetRol() {
	defer errorhandler.HandlePanic(&c.Controller)
	var (
		v models.UserName
	)

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		beego.Error(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	if response, err := services.GetRolesByUser(v); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = response // ? No se usa el APIResponseDTO debido a que este endpoint se utiliza por muchos servicios en producción con el formato de respuesta antiguo
	} else {
		beego.Error(err)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}
