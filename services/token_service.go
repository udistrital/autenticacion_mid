package services

import (
	"errors"
	
	"github.com/udistrital/autenticacion_mid/helpers"
	"github.com/udistrital/autenticacion_mid/models"
)

// GetInfoByEmail ...
func GetInfoByEmail(m *models.Token) (u *models.UserInfo, err error) {
	userRoles := []string{}
	EstudianteInfo, err := helpers.GetCodeByEmailStudentService(m.Email)

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
func GetRolesByUser(user models.UserName) (*models.Payload, map[string]interface{}) {
	userRoles := []string{}

	RolesUsuario, err := helpers.GetRolesUsuario(user.User)
	if err != nil {
		outputError := map[string]interface{}{"Function": "FuncionalidadMidController:userRol", "Error": err}
		return nil, outputError
	}

	if len(RolesUsuario.Usuario.Atributos) > 0 {
		payload, err := helpers.GetPayload(userRoles, RolesUsuario)
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
func GetInfoDocumento(user models.Documento) (*models.Payload, map[string]interface{}) {
	userRoles := []string{}

	RolesUsuario, err := helpers.GetInfoByDocumentoService(user.Numero)
	if err != nil {
		outputError := map[string]interface{}{"Function": "FuncionalidadMidController:GetInfoDocumento", "Error": err}
		return nil, outputError
	}

	if len(RolesUsuario.Usuario.Atributos) > 0 {
		payload, err := helpers.GetPayload(userRoles, RolesUsuario)
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

	response, err := helpers.ClientCredentialsRequest(req.ClienteId)
	outputError = map[string]interface{}{"Function": "TokenController:GetClientAuth", "Error": err}

	return
}
