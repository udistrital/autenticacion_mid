package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/autenticacion_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
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

// GetRol ...
// @Title GetRol
// @Description Recibe el usuario y devuelve el rol
// @Param	email		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.RolesUsuario
// @Failure 403 is empty
// @router /userRol [post]
func (c *TokenController) GetRol() {
	var (
		v   models.Token
		res map[string]interface{}
	)
	res = make(map[string]interface{}, 0)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		beego.Info("V: ", v)
		if response, err := models.GetRolesByUser(&v); err == nil {
			mapa1, _ := formatdata.ToMap(response, "json")
			beego.Info(mapa1)
			for k, v := range mapa1 {
				res[k] = v
			}
		} else {
			fmt.Println(err)
			c.Data["system"] = err
		}
	} else {
		fmt.Println(err)
		c.Data["system"] = err
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		beego.Info("V: ", v)
		if response, err := models.GetInfoByEmail(&v); err == nil {
			mapa2, _ := formatdata.ToMap(response, "json")
			for k, v := range mapa2 {
				res[k] = v
			}

		} else {
			fmt.Println(err)
			c.Data["system"] = err
		}
	} else {
		fmt.Println(err)
		c.Data["system"] = err
	}
	c.Data["json"] = res
	c.ServeJSON()

}
