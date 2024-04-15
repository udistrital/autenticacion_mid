package services

import (
	"errors"
	"strings"

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

// GetRolesByUser ...
func GetRolesByUser(user models.UserName) (roles *models.Payload, outputError error) {
	var familyName string
	var documento string
	var mail string
	var documentoCompuesto string
	userRoles := []string{}

	RolesUsuario, err := helpers.GetRolesUsuario(user.User)

	if err != nil {
		return nil, err
	}

	if len(RolesUsuario.Usuario.Atributos) > 0 {
		for _, v := range RolesUsuario.Usuario.Atributos {
			switch v.Atributo {
			case "role":
				roles := strings.Split(v.Valor, ",")
				for _, v := range roles {
					userRoles = append(userRoles, v)
				}
			case "sn":
				familyName = v.Valor
			case "documento":
				documento = v.Valor

			case "documento_compuesto":
				documentoCompuesto = v.Valor
			case "mail":
				mail = v.Valor
			}

			// fmt.Println(k, v)
		}
		payload := &models.Payload{
			Role:               userRoles,
			DocumentoCompuesto: documentoCompuesto,
			Documento:          documento,
			Email:              mail,
			FamilyName:         familyName,
		}

		EstudianteInfo, err := helpers.GetCodeByEmailStudentService(mail)

		if err != nil {
			return nil, err
		}

		if len(EstudianteInfo.EstudianteCollection.Estudiante) > 0 {
			userRoles = append(userRoles, "ESTUDIANTE")
			payload.Codigo = EstudianteInfo.EstudianteCollection.Estudiante[0].Codigo
			payload.Estado = EstudianteInfo.EstudianteCollection.Estudiante[0].Estado
			payload.Role = userRoles
		}

		return payload, nil

	} else {
		return nil, errors.New("Usuario no registrado")
	}
}
