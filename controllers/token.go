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
	c.Mapping("GetRol", c.GetRol)
	c.Mapping("AddRol", c.AddRol)
	c.Mapping("DeleteRol", c.DeleteRol)
}

// GetEmail ...
// @Title GetEmail
// @Description Recibe el correo electrónico del usuario desde la autenticación
// @Param	body body 	models.Token	true		"The key for staticblock"
// @Success 200 {object} models.UserInfo
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
// @Param	body body 	models.UpdateRol	true "Usuario y roles a adicionar"
// @Success 200 {object} models.resUpdateRol
// @Failure 404 not found resource
// @router /addRol [post]
func (c *TokenController) AddRol() {
	defer func() {
		if err := recover(); err != nil {
				localError := err.(map[string]interface{})
				c.Data["json"] = map[string]interface{}{
						"Success": false,
						"Status":  localError["status"].(string),
						"Message": localError["err"].(string),
						"Data":    nil,
				}
				c.ServeJSON()
		}
	}()
	var (
		v models.UpdateRol
	)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		beego.Info("V: ", v)
		if response := models.AddRol(v); err == nil {
				// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": response}
			if response["InfoUser"] != nil {
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "200", "Message": "El usuario ya tiene el rol " + v.Rol + " asignado", "Data": response}
			} else {
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": response}
			}
		} else {
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "404", "Message": "Unsuccessful", "Data": nil}
			c.Abort("404")
		}
	} else {
		fmt.Println("error: ", err)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Unsuccessful", "Data": nil}
		c.Abort("400")
	}
	c.ServeJSON()
}

// DeleteRol ...
// @Title DeleteRol
// @Description Recibe el usuario y el rol
// @Param	body body 	models.UpdateRol	true "Usuario y rol a eliminar"
// @Success 200 {object} models.resUpdateRol
// @Failure 404 not found resource
// @router /deleteRol [post]
func (c *TokenController) DeleteRol() {
	defer func() {
		if err := recover(); err != nil {
				localError := err.(map[string]interface{})
				c.Data["json"] = map[string]interface{}{
						"Success": false,
						"Status":  localError["status"].(string),
						"Message": localError["err"].(string),
						"Data":    nil,
				}
				c.ServeJSON()
		}
	}()
	var (
		v models.UpdateRol
	)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		beego.Info("V: ", v)
		if response := models.DeleteRol(v); err == nil {
			if response["InfoUser"] != nil {
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "200", "Message": "El usuario ya tiene el rol " + v.Rol + " desvinculado", "Data": response}
			} else {
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": response}
			}
		} else {
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Unsuccessful", "Data": nil}
			c.Abort("404")
		}
	} else {
		fmt.Println("error: ", err)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Unsuccessful", "Data": nil}
		c.Abort("400")
	}
	c.ServeJSON()
}
