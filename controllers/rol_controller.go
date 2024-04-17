package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/autenticacion_mid/models"
	"github.com/udistrital/autenticacion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// RolController operations for Rol
type RolController struct {
	beego.Controller
}

// URLMapping ...
func (c *RolController) URLMapping() {
	c.Mapping("AddRol", c.AddRol)
	c.Mapping("RemoveRol", c.RemoveRol)
}

// AddRol ...
// @Title AddRol
// @Description Recibe el usuario y el rol
// @Param	body body 	models.UpdateRol	true "Usuario y roles a adicionar"
// @Success 200 {object} models.ResponseDTO
// @Failure 400 El rol ya está asignado al usuario
// @Failure 404 not found resource
// @router /add [post]
func (c *RolController) AddRol() {
	defer errorhandler.HandlePanic(&c.Controller)

	var (
		v models.UpdateRol
	)

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		beego.Error(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	if response, err := services.AddRol(v); err != nil {
		beego.Error(err)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	} else {
		if response["InfoUser"] != nil { // El rol ya existe en el usuario
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = requestresponse.APIResponseDTO(false, 400, response, "El rol ya está asignado al usuario")
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = requestresponse.APIResponseDTO(true, 200, response)
		}
	}

	c.ServeJSON()
}

// RemoveRol ...
// @Title RemoveRol
// @Description Recibe el usuario y el rol
// @Param	body body 	models.UpdateRol	true "Usuario y rol a remover"
// @Success 200 {object} models.ResponseDTO
// @Failure 404 not found resource
// @router /remove [post]
func (c *RolController) RemoveRol() {
	defer errorhandler.HandlePanic(&c.Controller)

	var (
		v models.UpdateRol
	)

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		beego.Error(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	if response, err := services.RemoveRol(v); err != nil {
		beego.Error(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, err.Error())
	} else {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, response)
	}
	c.ServeJSON()
}
