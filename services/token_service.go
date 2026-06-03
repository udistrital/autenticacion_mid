package services

import (
	"context"
	"encoding/base64"
	"errors"
	"os"

	"github.com/udistrital/autenticacion_mid/helpers"
	"github.com/udistrital/autenticacion_mid/models"
)

// GetInfoByEmail ...
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

// GetRolesByUser
func GetRolesByUser(ctx context.Context, user models.UserName) (*models.Payload, map[string]interface{}) {
	userRoles := []string{}

	RolesUsuario, err := helpers.GetRolesUsuario(ctx, user.User)
	if err != nil {
		outputError := map[string]interface{}{"Function": "FuncionalidadMidController:userRol", "Error": err}
		return nil, outputError
	}

	if RolesUsuario.Usuario != nil && len(RolesUsuario.Usuario.Atributos) > 0 {
		payload, err := helpers.GetPayload(ctx, userRoles, RolesUsuario)
		if err != nil {
			outputError := map[string]interface{}{"Function": "FuncionalidadMidController:userRol", "Error": err}
			return nil, outputError
		}
		return payload, nil
	}

	outputError := map[string]interface{}{"Function": "FuncionalidadMidController:userRol", "Error": "Usuario no registrado"}
	return nil, outputError
}

// GetInfoDocumento
func GetInfoDocumento(ctx context.Context, user models.Documento) (*models.Payload, map[string]interface{}) {
	userRoles := []string{}

	RolesUsuario, err := helpers.GetInfoByDocumentoService(ctx, user.Numero)
	if err != nil {
		outputError := map[string]interface{}{"Function": "FuncionalidadMidController:GetInfoDocumento", "Error": err}
		return nil, outputError
	}

	if len(RolesUsuario.Usuario.Atributos) > 0 {
		payload, err := helpers.GetPayload(ctx, userRoles, RolesUsuario)
		if err != nil {
			outputError := map[string]interface{}{"Function": "FuncionalidadMidController:GetInfoDocumento", "Error": err}
			return nil, outputError
		}
		return payload, nil
	}

	outputError := map[string]interface{}{"Function": "FuncionalidadMidController:GetInfoDocumento", "Error": "Usuario no registrado"}
	return nil, outputError
}

// GetClientAuth
func GetClientAuth(req models.ClientAuthRequestBody) (response models.ClientCredentialsResponse, outputError map[string]interface{}) {

	if req.ClienteId == "" || req.Documento == "" {
		outputError = map[string]interface{}{"Function": "TokenController:GetClientAuth", "Error": "Debe indicar un cliente y un documento"}
	}

	decodedClientId, err := base64.StdEncoding.DecodeString(req.ClienteId)
	if err != nil {
		outputError = map[string]interface{}{"Function": "TokenController:GetClientAuth", "Error": err}
	}

	secret := os.Getenv("SECRET_CLIENT_" + string(decodedClientId))
	if secret == "" {
		outputError = map[string]interface{}{"Function": "TokenController:GetClientAuth", "Error": "No se pudo generar la autorización para el cliente indicado"}
		return
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(string(decodedClientId) + ":" + secret))
	response, err = helpers.ClientCredentialsRequest(encoded)
	if err != nil {
		outputError = map[string]interface{}{"Function": "TokenController:GetClientAuth", "Error": err}
	}

	return
}
