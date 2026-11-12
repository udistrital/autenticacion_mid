package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"github.com/udistrital/autenticacion_mid/models"
	"github.com/udistrital/autenticacion_mid/services"
)

// TokenController entrega la información de los usuarios para sus tokens.
type TokenController struct {
	beego.Controller
}

func (c *TokenController) URLMapping() {
	c.Mapping("GetEmail", c.GetEmail)
	c.Mapping("GetRol", c.GetRol)
	c.Mapping("GetDocumento", c.GetDocumento)
	c.Mapping("ClientAuth", c.ClientAuth)
}

// @Title GetEmail
// @Description Devuelve el código, el estado y el rol ESTUDIANTE del estudiante con el correo.
// @Param	body	body	models.Token	true	"Correo del estudiante"
// @Success 200 {object} models.UserInfo
// @Failure 400 El cuerpo no es válido o el correo no es de un estudiante
// @router /emailToken [post]
func (c *TokenController) GetEmail() {

	var v models.Token

	ctx := c.Ctx.Request.Context()

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if response, err := services.GetInfoByEmail(ctx, &v); err == nil {
			c.Data["json"] = response
		} else {
			logs.Error(err)
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Abort("400")
	}

	c.ServeJSON()
}

// @Title GetRol
// @Description Devuelve los datos y roles del usuario, incluidos sus datos de estudiante.
// @Param	body	body	models.UserName	true	"Correo o nombre de usuario"
// @Success 200 {object} models.InfoUsuario
// @Failure 400 El cuerpo no es válido o no indica el usuario
// @Failure 500 El usuario no está registrado o no se pudo consultar
// @router /userRol [post]
func (c *TokenController) GetRol() { // ? No se usa el APIResponseDTO debido a que este endpoint se utiliza por muchos servicios en producción con el formato de respuesta antiguo

	var v models.UserName

	ctx := c.Ctx.Request.Context()

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Data["system"] = err
		c.Abort("400")
	}

	valid := validation.Validation{}
	if ok, _ := valid.Valid(&v); !ok {
		c.Data["system"] = valid.Errors
		c.Abort("400")
	}

	if response, err := services.GetInfoUsuario(ctx, v); err == nil {
		c.Data["json"] = response
	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("500")
	}

	c.ServeJSON()
}

// @Title GetDocumento
// @Description Devuelve los datos y roles del usuario con el documento, incluidos sus datos de estudiante.
// @Param	body	body	models.Documento	true	"Documento del usuario"
// @Success 200 {object} models.InfoUsuario
// @Failure 400 El cuerpo no es válido, el usuario no está registrado o no se pudo consultar
// @router /documentoToken [post]
func (c *TokenController) GetDocumento() {

	var v models.Documento

	ctx := c.Ctx.Request.Context()

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if response, err := services.GetInfoDocumento(ctx, v); err == nil {
			c.Data["json"] = response
		} else {
			logs.Error(err)
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Abort("400")
	}

	c.ServeJSON()
}

// @Title ClientAuth
// @Description Devuelve un token client_credentials para el cliente.
// @Param	body	body	models.ClientAuthRequestBody	true	"ID del cliente en base64 y documento del usuario que pide el token"
// @Success 200 {object} models.ClientCredentialsResponse
// @Failure 400 El cuerpo no es válido, el cliente no está configurado o no se pudo obtener el token
// @router /clientAuth [post]
func (c *TokenController) ClientAuth() {

	var body models.ClientAuthRequestBody

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	if err != nil {
		logs.Error(err)
		c.Abort("400")
	}

	valid := validation.Validation{}
	if ok, _ := valid.Valid(&body); !ok {
		logs.Error(valid.Errors)
		c.Abort("400")
	}

	response, error := services.GetClientAuth(c.Ctx.Request.Context(), body)
	if error != nil {
		logs.Error(error)
		c.Abort("400")
	} else {
		c.Data["json"] = response
	}

	c.ServeJSON()
}
