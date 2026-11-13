package services

import (
	"context"
	"encoding/base64"
	"errors"
	"os"
	"slices"

	"github.com/astaxie/beego"
	"github.com/udistrital/autenticacion_mid/helpers"
	"github.com/udistrital/autenticacion_mid/models"
	"github.com/udistrital/autenticacion_mid/scim"
)

var secretPrefix = beego.AppConfig.String("secretPrefix")

func GetInfoByEmail(ctx context.Context, m *models.Token) (u *models.UserInfo, err error) {
	userRoles := []string{}
	EstudianteInfo, err := helpers.GetCodeByEmailStudentService(ctx, m.Email)

	if err != nil {
		return nil, err
	}

	if len(EstudianteInfo.EstudianteCollection.Estudiante) <= 0 {
		return nil, errors.New("email no registrado")
	}

	userRoles = append(userRoles, "ESTUDIANTE")

	u = &models.UserInfo{
		Codigo: EstudianteInfo.EstudianteCollection.Estudiante[0].Codigo,
		Estado: EstudianteInfo.EstudianteCollection.Estudiante[0].Estado,
		Email:  m.Email,
		Rol:    userRoles,
	}

	return u, nil
}

// GetInfoUsuario devuelve los datos, roles y grupos del usuario en IS 7.3, sin los
// grupos con roles, junto con sus datos de estudiante.
func GetInfoUsuario(ctx context.Context, user models.UserName) (*models.InfoUsuario, map[string]any) {
	usuario, err := scim.GetGruposUsuario(ctx, user.User)
	if err != nil {
		outputError := map[string]any{"Function": "FuncionalidadMidController:userRol", "Error": err}
		return nil, outputError
	}

	if usuario == nil {
		outputError := map[string]any{"Function": "FuncionalidadMidController:userRol", "Error": "Usuario no registrado"}
		return nil, outputError
	}
	usuario.Grupos = filtrarGruposConRoles(ctx, usuario.Grupos)

	estudiante, err := helpers.GetEstudiante(ctx, usuario.Email)
	if err != nil {
		outputError := map[string]any{"Function": "FuncionalidadMidController:userRol", "Error": err}
		return nil, outputError
	}

	return mergeUsuarioEstudiante(*usuario, estudiante), nil
}

// GetInfoDocumento devuelve los datos, roles y grupos del usuario en IS 7.3 por
// documento, sin los grupos con roles, junto con sus datos de estudiante.
func GetInfoDocumento(ctx context.Context, user models.Documento) (*models.InfoUsuario, map[string]any) {
	usuario, err := scim.GetInfoByDocumento(ctx, user.Numero)
	if err != nil {
		outputError := map[string]any{"Function": "FuncionalidadMidController:GetInfoDocumento", "Error": err}
		return nil, outputError
	}

	if usuario == nil {
		outputError := map[string]any{"Function": "FuncionalidadMidController:GetInfoDocumento", "Error": "Usuario no registrado"}
		return nil, outputError
	}
	usuario.Grupos = filtrarGruposConRoles(ctx, usuario.Grupos)

	estudiante, err := helpers.GetEstudiante(ctx, usuario.Email)
	if err != nil {
		outputError := map[string]any{"Function": "FuncionalidadMidController:GetInfoDocumento", "Error": err}
		return nil, outputError
	}

	return mergeUsuarioEstudiante(*usuario, estudiante), nil
}

// mergeUsuarioEstudiante construye los datos del usuario, agregando el rol
// ESTUDIANTE, el código y el estado si estudiante no es nil.
func mergeUsuarioEstudiante(usuario scim.Usuario, estudiante *models.Estudiante) *models.InfoUsuario {
	info := &models.InfoUsuario{
		Role:               append(slices.Clone(usuario.Roles), usuario.Grupos...),
		DocumentoCompuesto: usuario.DocumentoCompuesto,
		Documento:          usuario.Documento,
		Email:              usuario.Email,
		FamilyName:         usuario.FamilyName,
	}

	if estudiante != nil {
		info.Role = append(info.Role, "ESTUDIANTE")
		info.Codigo = estudiante.Codigo
		info.Estado = estudiante.Estado
	}
	return info
}

func GetClientAuth(ctx context.Context, req models.ClientAuthRequestBody) (response models.ClientCredentialsResponse, outputError map[string]any) {
	decodedClientId, err := base64.StdEncoding.DecodeString(req.ClienteId)
	if err != nil {
		outputError = map[string]any{"Function": "TokenController:GetClientAuth", "Error": err}
		return
	}

	secret := os.Getenv(secretPrefix + string(decodedClientId))
	if secret == "" {
		outputError = map[string]any{"Function": "TokenController:GetClientAuth", "Error": "No se pudo generar la autorización para el cliente indicado"}
		return
	}

	response, err = scim.ClientCredentialsRequest(ctx, string(decodedClientId), secret, "")
	if err != nil {
		outputError = map[string]any{"Function": "TokenController:GetClientAuth", "Error": err}
	}

	return
}
