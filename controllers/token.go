package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/autenticacion_mid/models"
)

// TokenController operations for Token
type TokenController struct {
	beego.Controller
}

// URLMapping ...
func (c *TokenController) URLMapping() {
	c.Mapping("GetEmail", c.GetEmail)
}

// GetEmail ...
// @Title GetEmail
// @Description Recibe el correo electrónico del usuario desde la autenticación
// @Param	email		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Token
// @Failure 403 is empty
// @router /emailToken [post]
func (c *TokenController) GetEmail() {
	var v models.Token
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		beego.Info("V: ", v)
		if response, err := models.GetInfoByEmail(&v); err == nil {
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
