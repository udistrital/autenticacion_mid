package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/udistrital/autenticacion_mid/models"
)

// Get_Methods

func GetUsuario(email string) (models.ResUserId, error) {
	var UserId models.ResUserId
	urlUsuario := beego.AppConfig.String("Wso2Service") + "usuario"
	requestUsuario := httplib.Get(urlUsuario + "?user=" + email)
	requestUsuario.Header("Accept", "application/json")

	if err := requestUsuario.ToJSON(&UserId); err != nil {
		return models.ResUserId{}, fmt.Errorf("Error al obtener el usuario en rol_helper.GetUsuario: %v", err)
	}

	if len(UserId.Usuarios.Usuario) <= 0 {
		return models.ResUserId{}, fmt.Errorf("El usuario %s no existe", email)
	}

	return UserId, nil
}

func GetRol(nombreRol string) (models.ResRolId, error) {
	var RolId models.ResRolId

	urlRol := beego.AppConfig.String("Wso2Service") + "rol"
	requestRol := httplib.Get(urlRol + "?rol=" + nombreRol)
	requestRol.Header("Accept", "application/json")

	if err := requestRol.ToJSON(&RolId); err != nil {
		return models.ResRolId{}, fmt.Errorf("Error al obtener el rol en rol_helper.GetRol: %v", err)
	}

	if len(RolId.Roles.Rol) <= 0 {
		return models.ResRolId{}, fmt.Errorf("El rol %s no existe", nombreRol)
	}

	return RolId, nil
}

func GetPerfilUsuario(idUsuario int) (models.ResPerfilUsuario, error) {
	var Perfil models.ResPerfilUsuario
	urlGetPerfilUsuario := beego.AppConfig.String("Wso2Service") + "perfil"

	requestPerfil := httplib.Get(urlGetPerfilUsuario + "?um_user_id=" + strconv.Itoa(idUsuario))
	requestPerfil.Header("Accept", "application/json")

	if err := requestPerfil.ToJSON(&Perfil); err != nil {
		return models.ResPerfilUsuario{}, fmt.Errorf("Error al obtener el perfil en rol_helper.GetPerfilUsuario: %v", err)
	}

	return Perfil, nil
}

func GetUsuarioRol(idUsuario string) (models.ResGetUsuarioRoles, error) {
	var respuesta models.ResGetUsuarioRoles
	urlGetUsuarioRol := beego.AppConfig.String("Wso2Service") + "get_usuario_rol?usuario=" + idUsuario

	requestUsuarioRol := httplib.Get(urlGetUsuarioRol)
	requestUsuarioRol.Header("Accept", "application/json")

	if err := requestUsuarioRol.ToJSON(&respuesta); err != nil {
		return models.ResGetUsuarioRoles{}, fmt.Errorf("Error al obtener los roles del usuario en rol_helper.GetUsuarioRol: %v", err)
	}

	if len(respuesta.Usuario.Roles) <= 0 {
		return models.ResGetUsuarioRoles{}, fmt.Errorf("El usuario %s no tiene roles asignados", idUsuario)
	}

	return respuesta, nil
}

// Post_Methods

func PostUsuarioRol(idUsuario int, idRol int) (models.ResUpdateUsuarioRol, error) {
	var respuesta models.ResUpdateUsuarioRol
	urlPostUsuarioRol := beego.AppConfig.String("Wso2UserService") + "usuario_rol"

	sendUsuarioRol := httplib.Post(urlPostUsuarioRol)
	body := models.UpdateUsuarioRol{
		UmRoleId: idRol,
		UmUserId: idUsuario,
	}

	sendUsuarioRol.Header("Accept", "application/json")
	sendUsuarioRol.Header("Content-Type", "application/json")
	sendUsuarioRol.JSONBody(body)

	if err := sendUsuarioRol.ToJSON(&respuesta); err != nil {
		return models.ResUpdateUsuarioRol{}, fmt.Errorf("Error al enviar la solicitud POST en rol_helper.PostUsuarioRol: %v", err)
	}

	return respuesta, nil
}

func PostPerfilUsuario(idUsuario int, value string) (map[string]map[string]interface{}, error) {
	var respuesta map[string]map[string]interface{}
	urlAddProfile := beego.AppConfig.String("Wso2UserService") + "addperfil"

	body := models.AddPerfil{
		UmUserId:    idUsuario,
		UmAttrValue: value,
	}

	addProfile := httplib.Post(urlAddProfile)
	addProfile.Header("Accept", "application/json")
	addProfile.Header("Content-Type", "application/json")
	addProfile.JSONBody(body)

	if err := addProfile.ToJSON(&respuesta); err != nil {
		return nil, fmt.Errorf("No se pudo añadir perfil, error en la petición en rol_helper.PostPerfilUsuario: %v", err)
	}

	return respuesta, nil
}

// Put_Methods

func UpdatePerfilUsuario(idPerfil int, value string) (models.ResUpdatePerfil, error) {
	var respuesta models.ResUpdatePerfil
	urlUpdateProfile := beego.AppConfig.String("Wso2UserService") + "updateperfil"

	body := models.UpdatePerfil{
		UmId:        idPerfil,
		UmAttrValue: value,
	}

	updateProfile := httplib.Put(urlUpdateProfile)
	updateProfile.Header("Accept", "application/json")
	updateProfile.Header("Content-Type", "application/json")
	updateProfile.JSONBody(body)

	if err := updateProfile.ToJSON(&respuesta); err != nil {
		return models.ResUpdatePerfil{}, fmt.Errorf("No se pudo actualizar el perfil, error en la petición en rol_helper.UpdatePerfilUsuario: %v", err)
	}

	return respuesta, nil
}

// Delete_Methods

func DeleteUsuarioRol(UmId string) (models.ResDeleteUsuarioRol, error) {
	var respuesta models.ResDeleteUsuarioRol
	urlDeleteUsuarioRol := beego.AppConfig.String("Wso2UserService") + "deleteusuariorol/"

	deleteUsuarioRol := httplib.Delete(urlDeleteUsuarioRol + UmId)
	deleteUsuarioRol.Header("Accept", "application/json")

	if err := deleteUsuarioRol.ToJSON(&respuesta); err != nil {
		return models.ResDeleteUsuarioRol{}, fmt.Errorf("Error al eliminar el usuario rol en rol_helper.DeleteUsuarioRol: %v", err)
	}

	return respuesta, nil
}

// Other_Methods_Helpers

func ObtenerStringPerfil(perfilValue string, rolEliminar string) string {
	perfilActualizado := strings.Split(perfilValue, ",")

	for i := range perfilActualizado {
		if perfilActualizado[i] == rolEliminar {
			perfilActualizado = append(perfilActualizado[:i], perfilActualizado[i+1:]...)
			break
		}
	}

	strPerfil := strings.Join(perfilActualizado, ",") // Convirtiendo array a string

	return strPerfil
}
