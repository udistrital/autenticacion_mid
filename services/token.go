package services

import (
	"context"
	"encoding/base64"
	"errors"
	"os"

	"github.com/udistrital/autenticacion_mid/helpers"
	"github.com/udistrital/autenticacion_mid/models"
)

func GetInfoByEmail(ctx context.Context, m *models.Token) (u *models.UserInfo, err error) {
	userRoles := []string{}
	EstudianteInfo, err := helpers.GetCodeByEmailStudentService(ctx, m.Email)

	if err != nil {
		return nil, err
	}

	if len(EstudianteInfo.EstudianteCollection.Estudiante) <= 0 {
		return nil, errors.New("Email no registrado")
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

// GetRolesByUser returns the user's info and roles from IS 7.3, merged with their student data.
func GetRolesByUser(ctx context.Context, user models.UserName) (*models.InfoUsuario, map[string]any) {
	usuario, err := helpers.GetGruposUsuarioIS(ctx, user.User)
	if err != nil {
		outputError := map[string]any{"Function": "FuncionalidadMidController:userRol", "Error": err}
		return nil, outputError
	}

	if usuario == nil {
		outputError := map[string]any{"Function": "FuncionalidadMidController:userRol", "Error": "Usuario no registrado"}
		return nil, outputError
	}

	estudiante, err := helpers.GetEstudiante(ctx, usuario.Email)
	if err != nil {
		outputError := map[string]any{"Function": "FuncionalidadMidController:userRol", "Error": err}
		return nil, outputError
	}

	return mergeUsuarioEstudiante(*usuario, estudiante), nil
}

// GetInfoDocumento returns the user's info and roles from EI by documento, merged with their student data.
func GetInfoDocumento(ctx context.Context, user models.Documento) (*models.InfoUsuario, map[string]any) {
	RolesUsuario, err := helpers.GetInfoByDocumentoService(ctx, user.Numero)
	if err != nil {
		outputError := map[string]any{"Function": "FuncionalidadMidController:GetInfoDocumento", "Error": err}
		return nil, outputError
	}

	if RolesUsuario.Usuario == nil || len(RolesUsuario.Usuario.Atributos) == 0 {
		outputError := map[string]any{"Function": "FuncionalidadMidController:GetInfoDocumento", "Error": "Usuario no registrado"}
		return nil, outputError
	}

	familyName, documento, mail, documentoCompuesto, roles := helpers.MapAtributos(RolesUsuario)

	estudiante, err := helpers.GetEstudiante(ctx, mail)
	if err != nil {
		outputError := map[string]any{"Function": "FuncionalidadMidController:GetInfoDocumento", "Error": err}
		return nil, outputError
	}

	return mergeUsuarioEstudiante(models.Usuario{
		Roles:              roles,
		FamilyName:         familyName,
		Email:              mail,
		Documento:          documento,
		DocumentoCompuesto: documentoCompuesto,
	}, estudiante), nil
}

// mergeUsuarioEstudiante builds the user's info, adding the ESTUDIANTE role, código and estado when estudiante is not nil.
func mergeUsuarioEstudiante(usuario models.Usuario, estudiante *models.Estudiante) *models.InfoUsuario {
	info := &models.InfoUsuario{
		Role:               usuario.Roles,
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

func GetClientAuth(req models.ClientAuthRequestBody) (response models.ClientCredentialsResponse, outputError map[string]any) {

	if req.ClienteId == "" || req.Documento == "" {
		outputError = map[string]any{"Function": "TokenController:GetClientAuth", "Error": "Debe indicar un cliente y un documento"}
	}

	decodedClientId, err := base64.StdEncoding.DecodeString(req.ClienteId)
	if err != nil {
		outputError = map[string]any{"Function": "TokenController:GetClientAuth", "Error": err}
	}

	secret := os.Getenv("SECRET_CLIENT_" + string(decodedClientId))
	if secret == "" {
		outputError = map[string]any{"Function": "TokenController:GetClientAuth", "Error": "No se pudo generar la autorización para el cliente indicado"}
		return
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(string(decodedClientId) + ":" + secret))
	response, err = helpers.ClientCredentialsRequest(encoded)
	if err != nil {
		outputError = map[string]any{"Function": "TokenController:GetClientAuth", "Error": err}
	}

	return
}
