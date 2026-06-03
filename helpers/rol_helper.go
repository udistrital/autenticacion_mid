package helpers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/autenticacion_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// Get_Methods

func GetUsuario(ctx context.Context, email string) (models.ResUserId, error) {
	var UserId models.ResUserId
	urlUsuario := beego.AppConfig.String("Wso2Service") + "usuario"

	if _, err := request.GetWithContext(ctx, urlUsuario+"?user="+email, &UserId); err != nil {
		return models.ResUserId{}, fmt.Errorf("Error al obtener el usuario en rol_helper.GetUsuario: %v", err)
	}

	if len(UserId.Usuarios.Usuario) <= 0 {
		return models.ResUserId{}, fmt.Errorf("El usuario %s no existe", email)
	}

	return UserId, nil
}

func GetRol(ctx context.Context, nombreRol string) (models.ResRolId, error) {
	var RolId models.ResRolId

	urlRol := beego.AppConfig.String("Wso2Service") + "rol"

	if _, err := request.GetWithContext(ctx, urlRol+"?rol="+nombreRol, &RolId); err != nil {
		return models.ResRolId{}, fmt.Errorf("Error al obtener el rol en rol_helper.GetRol: %v", err)
	}

	if len(RolId.Roles.Rol) <= 0 {
		return models.ResRolId{}, fmt.Errorf("El rol %s no existe", nombreRol)
	}

	return RolId, nil
}

func GetPerfilUsuario(ctx context.Context, idUsuario int) (models.ResPerfilUsuario, error) {
	var Perfil models.ResPerfilUsuario
	urlGetPerfilUsuario := beego.AppConfig.String("Wso2Service") + "perfil"

	if _, err := request.GetWithContext(ctx, urlGetPerfilUsuario+"?um_user_id="+strconv.Itoa(idUsuario), &Perfil); err != nil {
		return models.ResPerfilUsuario{}, fmt.Errorf("Error al obtener el perfil en rol_helper.GetPerfilUsuario: %v", err)
	}

	return Perfil, nil
}

func GetUsuarioRol(ctx context.Context, idUsuario string) (models.ResGetUsuarioRoles, error) {
	var respuesta models.ResGetUsuarioRoles
	urlGetUsuarioRol := beego.AppConfig.String("Wso2Service") + "get_usuario_rol?usuario=" + idUsuario

	if _, err := request.GetWithContext(ctx, urlGetUsuarioRol, &respuesta); err != nil {
		return models.ResGetUsuarioRoles{}, fmt.Errorf("Error al obtener los roles del usuario en rol_helper.GetUsuarioRol: %v", err)
	}

	if len(respuesta.Usuario.Roles) <= 0 {
		return models.ResGetUsuarioRoles{}, fmt.Errorf("El usuario %s no tiene roles asignados", idUsuario)
	}

	return respuesta, nil
}

// Post_Methods

func PostUsuarioRol(ctx context.Context, idUsuario int, idRol int) (models.ResUpdateUsuarioRol, error) {
	var respuesta models.ResUpdateUsuarioRol
	urlPostUsuarioRol := beego.AppConfig.String("Wso2UserService") + "usuario_rol"

	body := models.UpdateUsuarioRol{
		UmRoleId: idRol,
		UmUserId: idUsuario,
	}

	if _, err := request.PostWithContext(ctx, urlPostUsuarioRol, body, &respuesta); err != nil {
		return models.ResUpdateUsuarioRol{}, fmt.Errorf("Error al enviar la solicitud POST en rol_helper.PostUsuarioRol: %v", err)
	}

	return respuesta, nil
}

func PostPerfilUsuario(ctx context.Context, idUsuario int, value string) (map[string]map[string]interface{}, error) {
	var respuesta map[string]map[string]interface{}
	urlAddProfile := beego.AppConfig.String("Wso2UserService") + "addperfil"

	body := models.AddPerfil{
		UmUserId:    idUsuario,
		UmAttrValue: value,
	}

	if _, err := request.PostWithContext(ctx, urlAddProfile, body, &respuesta); err != nil {
		return nil, fmt.Errorf("No se pudo añadir perfil, error en la petición en rol_helper.PostPerfilUsuario: %v", err)
	}

	return respuesta, nil
}

// Put_Methods

func UpdatePerfilUsuario(ctx context.Context, idPerfil int, value string) (models.ResUpdatePerfil, error) {
	var respuesta models.ResUpdatePerfil
	urlUpdateProfile := beego.AppConfig.String("Wso2UserService") + "updateperfil"

	body := models.UpdatePerfil{
		UmId:        idPerfil,
		UmAttrValue: value,
	}

	if _, err := request.PutWithContext(ctx, urlUpdateProfile, body, &respuesta); err != nil {
		return models.ResUpdatePerfil{}, fmt.Errorf("No se pudo actualizar el perfil, error en la petición en rol_helper.UpdatePerfilUsuario: %v", err)
	}

	return respuesta, nil
}

// Delete_Methods

func DeleteUsuarioRol(ctx context.Context, UmId string) (models.ResDeleteUsuarioRol, error) {
	var respuesta models.ResDeleteUsuarioRol
	urlDeleteUsuarioRol := beego.AppConfig.String("Wso2UserService") + "deleteusuariorol/"

	if _, err := request.DeleteWithContext(ctx, urlDeleteUsuarioRol+UmId, &respuesta); err != nil {
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
