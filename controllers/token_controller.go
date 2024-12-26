package controllers

import (
	"encoding/json"
	
	"github.com/astaxie/beego"
	"github.com/udistrital/autenticacion_mid/models"
	"github.com/udistrital/autenticacion_mid/services"
)

// TokenController operations for Token
type TokenController struct {
	beego.Controller
}

// URLMapping ...
func (c *TokenController) URLMapping() {
	c.Mapping("GetEmail", c.GetEmail)
	c.Mapping("GetRol", c.GetRol)
	c.Mapping("GetDocumento", c.GetDocumento)
	c.Mapping("ClientAuth", c.ClientAuth)
}

// GetEmail ...
// @Title GetEmail
// @Description Recibe el correo electrónico del usuario desde la autenticación
// @Param	body body 	models.Token	true		"The key for staticblock"
// @Success 200 {object} models.UserInfo
// @Failure 404 not found resource
// @router /emailToken [post]
func (c *TokenController) GetEmail() { // ? No se usa el APIResponseDTO debido a que este endpoint se utiliza por muchos servicios en producción con el formato de respuesta antiguo
	
	var (
		v models.Token
	)

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if response, err := services.GetInfoByEmail(&v); err == nil {
			c.Data["json"] = response
		} else {
			beego.Error(err)
			c.Abort("400")
		}
	} else {
		beego.Error(err)
		c.Abort("400")
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
func (c *TokenController) GetRol() { // ? No se usa el APIResponseDTO debido a que este endpoint se utiliza por muchos servicios en producción con el formato de respuesta antiguo
	
	var (
		v models.UserName
	)

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if response, err := services.GetRolesByUser(v); err == nil {
			c.Data["json"] = response
		} else {
			c.Data["system"] = err
			c.Abort("400")
		}
	} else {
		c.Data["system"] = err
		c.Abort("400")
	}

	c.ServeJSON()
}
// GetDocumento ...
// @Title GetDocumento
// @Description Recibe el documento y devuelve información detallada del usuario
// @Param	body	body 	models.Documento  	true		"Documento del usuario"
// @Success 200 {object} models.Payload
// @Failure 404 not found resource
// @router /documentoToken [post]
func (c *TokenController) GetDocumento() { 
	
	var (
		v models.Documento
	)
	
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if response, err := services.GetInfoDocumento(v); err == nil {
			c.Data["json"] = response
		} else {
			beego.Error(err)
			c.Abort("400")
		}
	} else {
		beego.Error(err)
		c.Abort("400")
	}

	c.ServeJSON()
}

// ClientAuth ...
// @Title ClientAuth
// @Description Recibe el id del cliente y el número de documento del usuario que solicita el token
// @Param	body	body	models.ClientAuthRequestBody	true	"ClienteId en base64 y Número de Documento del usuario"
// @Success 200 {object} models.ClientAuthRequestBody
// @Failure 404 not found resource
// @router /clientAuth [post]
func (c *TokenController) ClientAuth() {

	var body models.ClientAuthRequestBody

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	if err != nil {
		beego.Error(err)
		c.Abort("400")
	}

	response, error := services.GetClientAuth(body)
	if error != nil {
		beego.Error(error)
		c.Abort("400")
	} else {
		c.Data["json"] = response
	}

	c.ServeJSON()
}
