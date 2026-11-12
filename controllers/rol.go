package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"github.com/udistrital/autenticacion_mid/helpers"
	"github.com/udistrital/autenticacion_mid/models"
	"github.com/udistrital/autenticacion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// RolController gestiona los roles de los usuarios y sus periodos.
type RolController struct {
	beego.Controller
}

func (c *RolController) URLMapping() {
	c.Mapping("AddRol", c.AddRol)
	c.Mapping("RemoveRol", c.RemoveRol)
	c.Mapping("GetPeriodoInfo", c.GetPeriodoInfo)
	c.Mapping("GetAllPeriodos", c.GetAllPeriodos)
}

// @Title AddRol
// @Description Asigna un rol a un usuario y devuelve sus roles.
// @Param	body	body	models.UpdateRol	true	"Usuario (correo o nombre de usuario) y rol a asignar"
// @Success 200 {object} requestresponse.APIResponse
// @Failure 400 El cuerpo de la petición no es válido o el rol no se puede gestionar
// @Failure 404 El usuario o el rol no existe
// @Failure 409 El rol ya está asignado al usuario
// @Failure 500 Error interno
// @router /add [post]
func (c *RolController) AddRol() {
	defer errorhandler.HandlePanic(&c.Controller)

	var v models.UpdateRol

	ctx := c.Ctx.Request.Context()

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		c.ServeJSON()
		return
	}

	if usuario, yaAsignado, err := services.AgregarGrupo(ctx, v.User, v.Rol); err != nil {
		logs.Error(err)
		status, mensaje := 500, "Error al agregar el rol al usuario"
		switch {
		case errors.Is(err, helpers.ErrGrupoNoGestionable):
			status, mensaje = 400, fmt.Sprintf("el rol %s no se puede gestionar", v.Rol)
		case errors.Is(err, helpers.ErrUsuarioNoExiste):
			status, mensaje = 404, fmt.Sprintf("el usuario %s no existe", v.User)
		case errors.Is(err, helpers.ErrGrupoNoExiste):
			status, mensaje = 404, fmt.Sprintf("el rol %s no existe", v.Rol)
		}
		c.Ctx.Output.SetStatus(status)
		c.Data["json"] = requestresponse.APIResponseDTO(false, status, nil, mensaje)
	} else {
		response := models.ResRol{InfoUser: models.InfoUserRoles{Role: rolesUsuario(usuario)}}
		if yaAsignado {
			c.Ctx.Output.SetStatus(409)
			c.Data["json"] = requestresponse.APIResponseDTO(false, 409, response, "El rol ya está asignado al usuario")
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = requestresponse.APIResponseDTO(true, 200, response)
		}
	}

	c.ServeJSON()
}

// @Title RemoveRol
// @Description Quita un rol a un usuario y devuelve sus roles.
// @Param	body	body	models.UpdateRol	true	"Usuario (correo o nombre de usuario) y rol a quitar"
// @Success 200 {object} requestresponse.APIResponse
// @Failure 400 El cuerpo de la petición no es válido o el rol no se puede gestionar
// @Failure 404 El usuario o el rol no existe
// @Failure 409 El usuario no tiene asignado el rol
// @Failure 500 Error interno
// @router /remove [post]
func (c *RolController) RemoveRol() {
	defer errorhandler.HandlePanic(&c.Controller)

	var v models.UpdateRol

	ctx := c.Ctx.Request.Context()

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		c.ServeJSON()
		return
	}

	if usuario, noAsignado, err := services.QuitarGrupo(ctx, v.User, v.Rol); err != nil {
		logs.Error(err)
		status, mensaje := 500, "Error al quitar el rol al usuario"
		switch {
		case errors.Is(err, helpers.ErrGrupoNoGestionable):
			status, mensaje = 400, fmt.Sprintf("el rol %s no se puede gestionar", v.Rol)
		case errors.Is(err, helpers.ErrUsuarioNoExiste):
			status, mensaje = 404, fmt.Sprintf("el usuario %s no existe", v.User)
		case errors.Is(err, helpers.ErrGrupoNoExiste):
			status, mensaje = 404, fmt.Sprintf("el rol %s no existe", v.Rol)
		}
		c.Ctx.Output.SetStatus(status)
		c.Data["json"] = requestresponse.APIResponseDTO(false, status, nil, mensaje)
	} else {
		response := models.ResRol{InfoUser: models.InfoUserRoles{Role: rolesUsuario(usuario)}}
		if noAsignado {
			c.Ctx.Output.SetStatus(409)
			c.Data["json"] = requestresponse.APIResponseDTO(false, 409, response, "El usuario no tiene asignado el rol")
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = requestresponse.APIResponseDTO(true, 200, response)
		}
	}

	c.ServeJSON()
}

// @Title GetPeriodoInfo
// @Description Lista los periodos de roles de un usuario.
// @Param	documento	path	string	true	"Documento del usuario"
// @Param	query	query	string	false	"Filtros con el formato campo:valor,campo:valor"
// @Param	limit	query	string	false	"Cantidad máxima de resultados"
// @Param	offset	query	string	false	"Posición del primer resultado"
// @Success 200 {object} []models.MetadataResponse
// @Failure 400 Los parámetros de consulta no son válidos
// @Failure 404 No se pudieron obtener los periodos o el tercero
// @Failure 500 Error interno
// @router /user/:documento/periods [get]
func (c *RolController) GetPeriodoInfo() {
	defer errorhandler.HandlePanic(&c.Controller)

	ctx := c.Ctx.Request.Context()
	documento := c.Ctx.Input.Param(":documento")

	params, query, err := c.paramsPeriodos()
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		c.ServeJSON()
		return
	}

	periodsData, err := services.GetPeriodoInfo(ctx, documento, query, params.Limit, params.Offset)
	if err != nil {
		logs.Error(err)
		status := 404
		if errors.Is(err, services.ErrConsultaUsuario) {
			status = 500
		}
		c.Ctx.Output.SetStatus(status)
		c.Data["json"] = requestresponse.APIResponseDTO(false, status, nil, "Error al obtener los periodos del usuario")
	} else {
		data := periodsData["Data"].([]models.PeriodoRolUsuario)
		metadata := periodsData["Metadata"].(map[string]any)
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseMetadataDTO(true, 200, data, metadata)
	}

	c.ServeJSON()
}

// @Title GetAllPeriodos
// @Description Lista los periodos de roles de todos los usuarios. Los usuarios que no se pudieron consultar se indican en Errores.
// @Param	query	query	string	false	"Filtros con el formato campo:valor,campo:valor"
// @Param	RolId	query	string	false	"IDs de rol separados por comas, p. ej. 1,2,3"
// @Param	limit	query	string	false	"Cantidad máxima de resultados"
// @Param	offset	query	string	false	"Posición del primer resultado"
// @Success 200 {object} []models.MetadataResponse
// @Failure 400 Los parámetros de consulta no son válidos
// @Failure 404 No se pudieron obtener los periodos
// @router /periods [get]
func (c *RolController) GetAllPeriodos() {
	defer errorhandler.HandlePanic(&c.Controller)

	ctx := c.Ctx.Request.Context()

	params, query, err := c.paramsPeriodos()
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		c.ServeJSON()
		return
	}
	if params.RolId != "" {
		query["RolId"] = params.RolId
	}

	periodsData, err := services.GetAllPeriodosRoles(ctx, query, params.Limit, params.Offset)
	if err != nil {
		logs.Error(err)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, "Error al obtener los periodos")
	} else {
		data := periodsData["Data"].([]models.PeriodoRolUsuario)
		metadata := periodsData["Metadata"].(map[string]any)
		errores := periodsData["Errores"]
		if errores != nil {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = requestresponse.APIResponseMetadataDTO(true, 200, data, metadata, errores)
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = requestresponse.APIResponseMetadataDTO(true, 200, data, metadata)
		}
	}
	c.ServeJSON()
}

// paramsPeriodos lee los parámetros de consulta de los periodos y devuelve los
// filtros de Query como mapa, o un error para el cliente si no son válidos.
func (c *RolController) paramsPeriodos() (params models.ParamsPeriodos, query map[string]string, err error) {
	if err := c.ParseForm(&params); err != nil {
		return params, nil, errors.New("limit y offset deben ser números enteros")
	}
	valid := validation.Validation{}
	if ok, _ := valid.Valid(&params); !ok {
		return params, nil, errors.New("el parámetro query debe tener el formato campo:valor,campo:valor")
	}

	query = make(map[string]string)
	if params.Query != "" {
		for _, cond := range strings.Split(params.Query, ",") {
			kv := strings.SplitN(cond, ":", 2)
			query[kv[0]] = kv[1]
		}
	}
	return params, query, nil
}

func rolesUsuario(u *models.UsuarioIS) []string {
	return append(slices.Clone(u.Roles), u.Grupos...)
}
