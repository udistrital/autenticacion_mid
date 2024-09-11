package controllers

import (
	"encoding/json"
	"errors"
	"strings"

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
	c.Mapping("GetPeriodoInfo", c.GetPeriodoInfo)
	c.Mapping("GetAllPeriodos", c.GetAllPeriodos)
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

// GetPeriodoInfo ...
// @Title GetPeriodoInfo
// @Description Obtiene los periodos de roles de un usuario por su documento
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	documento	path 	string	true	"Documento del usuario"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} []models.PeriodoRolUsuario
// @Failure 404 not found resource
// @router /user/:documento/periods [get]
func (c *RolController) GetPeriodoInfo() {
	defer errorhandler.HandlePanic(&c.Controller)

	var query = make(map[string]string)
	var limit int64
	var offset int64

	documento := c.Ctx.Input.Param(":documento")

	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}

	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}

	periods, err := services.GetPeriodoInfo(documento, query, limit, offset)
	if err != nil {
		beego.Error(err)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	} else {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, periods)
	}

	c.ServeJSON()
}

// GetAllPeriodos ...
// @Title GetAllPeriodos
// @Description Obtiene los periodos de todos los usuarios
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} []models.PeriodoRolUsuario
// @Failure 404 not found resource
// @router /periods [get]
func (c *RolController) GetAllPeriodos() {
	defer errorhandler.HandlePanic(&c.Controller)

	var query = make(map[string]string)
	var limit int64
	var offset int64

	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}

	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}

	response, err := services.GetAllPeriodosRoles(query, limit, offset)
	if err != nil {
		beego.Error(err)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	} else {
		data := response["Data"].([]models.PeriodoRolUsuario)
		metadata := response["Metadata"].(map[string]interface{})
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseMetadataDTO(true, 200, data, metadata)
	}
	c.ServeJSON()
}
