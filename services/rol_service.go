package services

import (
	"errors"
	"strconv"
	"strings"

	"github.com/udistrital/autenticacion_mid/helpers"
	"github.com/udistrital/autenticacion_mid/models"
)

// AddRol ...
func AddRol(user models.UpdateRol) (map[string]interface{}, error) {
	var userName models.UserName
	var respuesta map[string]interface{}
	var rolEncontrado = false
	userName.User = user.User

	RolesUsuario, err1 := GetRolesByUser(userName)

	if err1 != nil {
		return nil, errors.New("Error al obtener los roles del usuario en rol_service.AddRol: " + err1["Error"].(string))
	}

	for i := range RolesUsuario.Role {
		if RolesUsuario.Role[i] == user.Rol {
			// El rol ya existe en el usuario
			rolEncontrado = true
			returnData := map[string]interface{}{
				"InfoUser": map[string]interface{}{
					"Role": RolesUsuario.Role,
				},
			}
			respuesta = returnData
		}
	}

	if rolEncontrado { // Retorna rol que ya existe
		return respuesta, nil
	}

	// Add role to user
	User, err := helpers.GetUsuario(user.User)
	if err != nil {
		return nil, err
	}

	Rol, err := helpers.GetRol(user.Rol)
	if err != nil {
		return nil, err
	}

	idRol, err := strconv.Atoi(Rol.Roles.Rol[0].Id)
	if err != nil {
		return nil, errors.New("Error al convertir el Id del rol a int")
	}

	idUsuario, err := strconv.Atoi(User.Usuarios.Usuario[0].Id)
	if err != nil {
		return nil, errors.New("Error al convertir el Id del usuario a int")
	}

	_, err = helpers.PostUsuarioRol(idUsuario, idRol)
	if err != nil {
		return nil, err
	}

	PerfilUsuario, err := helpers.GetPerfilUsuario(idUsuario)

	if err != nil {
		return nil, err
	}

	if len(PerfilUsuario.Perfiles.Perfil) == 0 {
		// Post profile if not exist

		valuePerfil := "Internal/everyone," + user.Rol
		PostPerfil, err := helpers.PostPerfilUsuario(idUsuario, valuePerfil)
		if PostPerfil != nil {
			return nil, err
		}

		if PostPerfil["perfiles"]["perfil"] == nil {
			return nil, errors.New("No se puede actualizar perfil! ")
		}

		// Estructura de respuesta
		res := map[string]interface{}{
			"post_addperfil": map[string]interface{}{
				"um_attr_value": "Internal/everyone," + user.Rol,
				"um_user_id":    idUsuario,
			},
		}
		respuesta = res
	} else {
		if PerfilUsuario.Perfiles.Perfil[0].UmAttrValue != "" &&
			PerfilUsuario.Perfiles.Perfil[0].UmId != "" {
			// Update profile if exist

			StringPerfil := PerfilUsuario.Perfiles.Perfil[0].UmAttrValue
			IdPerfil := PerfilUsuario.Perfiles.Perfil[0].UmId
			UmIdProfile, err := strconv.Atoi(IdPerfil)

			if err != nil {
				return nil, errors.New("Error al convertir el Id del perfil a int")
			}

			rolesUsuario := strings.Split(StringPerfil, ",")
			rolExiste := false

			for _, role := range rolesUsuario {
				if role == user.Rol {
					rolExiste = true
					break
				}
			}

			if rolExiste { // Si el perfil ya tiene el rol
				return nil, errors.New("El rol ya existe en el perfil del usuario")
			}

			valuePerfil := user.Rol + "," + StringPerfil
			UpdatePerfil, err := helpers.UpdatePerfilUsuario(UmIdProfile, valuePerfil)

			if err != nil {
				return nil, err
			}

			if UpdatePerfil.Perfiles.Perfil == nil {
				return nil, errors.New("No se puede actualizar perfil!")
			}

			// Estructura de respuesta
			res := map[string]interface{}{
				"put_updateperfil": map[string]interface{}{
					"um_attr_value": user.Rol + "," + StringPerfil,
					"um_id":         UmIdProfile,
				},
			}
			respuesta = res
		}
	}
	return respuesta, nil
}

func RemoveRol(user models.UpdateRol) (map[string]interface{}, error) {
	var userName models.UserName
	var rolEncontrado = false
	var rolEliminar string
	userName.User = user.User

	// Obtener Roles del Usuario
	RolesUsuario, err1 := GetRolesByUser(userName)
	if err1 != nil {
		return nil, errors.New("Error al obtener los roles del usuario en rol_service.RemoveRol: " + err1["Error"].(string))
	}

	for i := range RolesUsuario.Role {
		if RolesUsuario.Role[i] == user.Rol {

			rolEncontrado = true
			rolEliminar = RolesUsuario.Role[i]

			// Obtener Usuario por email
			User, err := helpers.GetUsuario(user.User)
			if err != nil {
				return nil, err
			}

			idUsuario, err := strconv.Atoi(User.Usuarios.Usuario[0].Id)
			if err != nil {
				return nil, errors.New("Error al convertir el Id del usuario a int")
			}

			// Obtener Rol por nombre
			Rol, err := helpers.GetRol(user.Rol)
			if err != nil {
				return nil, err
			}

			// Obtener UsuarioRol (Tabla de Rompimiento)
			UsuarioRol, err := helpers.GetUsuarioRol(User.Usuarios.Usuario[0].Id)

			if err != nil {
				return nil, err
			}

			for i := range UsuarioRol.Usuario.Roles {
				if UsuarioRol.Usuario.Roles[i].UmRoleId == Rol.Roles.Rol[0].Id {
					// Role found

					// Eliminar UsuarioRol (Tabla de Rompimiento)
					_, err := helpers.DeleteUsuarioRol(UsuarioRol.Usuario.Roles[i].UmId)

					if err != nil {
						return nil, err
					}
				}
			}

			// Obtener Perfil de Usuario por idUsuario
			PerfilUsuario, err := helpers.GetPerfilUsuario(idUsuario)
			IdPerfil := PerfilUsuario.Perfiles.Perfil[0].UmId
			UmIdProfile, err := strconv.Atoi(IdPerfil)

			if err != nil {
				return nil, errors.New("Error al convertir el Id del perfil a int")
			}

			if len(PerfilUsuario.Perfiles.Perfil) > 0 { // El perfil existe

				// Obtener string para actualizar perfil
				strPerfil := helpers.ObtenerStringPerfil(PerfilUsuario.Perfiles.Perfil[0].UmAttrValue, rolEliminar)

				// Actualizar perfil
				UpdatePerfil, err := helpers.UpdatePerfilUsuario(UmIdProfile, strPerfil)
				if err != nil {
					return nil, err
				}

				response := map[string]interface{}{
					"put_updateperfil": map[string]interface{}{
						"um_attr_value": strPerfil,
						"um_id":         UpdatePerfil.Perfiles.Perfil[0].UmId,
					},
				}
				return response, nil
			}
		}
	}
	if !rolEncontrado {
		return nil, errors.New("El usuario no tiene asignado el rol")
	}
	return nil, errors.New("Error en la petición")
}
