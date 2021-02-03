package controllers

import (
	"encoding/json"
	"fmt"

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
// @Failure 404 not found resource
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

// GetRol ...
// @Title GetRol
// @Description Recibe el usuario y devuelve información detallada del usuario
// @Param	body	body 	models.UserName  	true		"Usuario registrado en wso2"
// @Success 200 {object} models.Payload
// @Failure 404 not found resource
// @router /userRol [post]
func (c *TokenController) GetRol() {
	var (
		v models.UserName
	)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		beego.Info("V: ", v)
		if response, err := models.GetRolesByUser(v); err == nil {
			c.Data["json"] = response
		} else {
			fmt.Println("error: ", err)
			c.Data["system"] = err
			c.Abort("400")
		}
	} else {
		fmt.Println("error: ", err)
		c.Data["system"] = err
		c.Abort("400")
	}

	c.ServeJSON()
}

// AddRol ...
// @Title AddRol
// @Description Recibe el usuario y el rol
// @Param	body	body 	models.UpdateRol  true	"Usuario registrado en wso2, rol en wso2"
// @Success 200 {object} models.Payload
// @Failure 404 not found resource
// @router /addRol [post]
func (c *TokenController) AddRol() {
	var (
		v models.UpdateRol
	)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		beego.Info("V: ", v)
		if response, err := models.AddRol(v); err == nil {
			c.Data["json"] = response
		} else {
			fmt.Println("error: ", err)
			c.Data["system"] = err
			c.Abort("400")
		}
	} else {
		fmt.Println("error: ", err)
		c.Data["system"] = err
		c.Abort("400")
	}

	c.ServeJSON()

}
