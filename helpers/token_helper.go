package helpers

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/autenticacion_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// Get_Methods

func GetRolesUsuario(email string) (models.AtributosToken, error) {
	var RolesUsuario models.AtributosToken
	urlGetRolesUsuario := beego.AppConfig.String("AutenticacionCrudService") + "roles?usuario=" + email

	err := request.GetJson(urlGetRolesUsuario, &RolesUsuario)
	if err != nil {
		logs.Info("Error", err)
		return models.AtributosToken{}, fmt.Errorf("Error al obtener los roles del usuario %s", email)
	}

	return RolesUsuario, nil
}

func GetCodeByEmailStudentService(email string) (models.EstudianteInfo, error) {
	var EstudianteInfo models.EstudianteInfo
	urlGetCodeByEmailStudentService := beego.AppConfig.String("GetCodeByEmailStudentService") + email
	err := request.GetXml(urlGetCodeByEmailStudentService, &EstudianteInfo)

	if err != nil {
		logs.Info("Error", err)
		return models.EstudianteInfo{}, fmt.Errorf("Error al obtener el código del estudiante %s", email)
	}

	return EstudianteInfo, nil
}

func GetInfoByDocumentoService(documento string) (models.AtributosToken, error) {
	var RolesUsuario models.AtributosToken
	urlGetInfoByDocumentoService := httplib.Get(beego.AppConfig.String("Wso2Service") + "usuario_documento?documento=" + documento)
	urlGetInfoByDocumentoService.Header("Accept", "application/json")

	err := urlGetInfoByDocumentoService.ToJSON(&RolesUsuario)
	if err != nil {
		return models.AtributosToken{}, fmt.Errorf("Error al obtener la información del documento %s", documento)
	}

	return RolesUsuario, nil

}

func GetPayload(userRoles []string, RolesUsuario models.AtributosToken) (*models.Payload, error) {
	familyName, documento, mail, documentoCompuesto, roles := MapAtributos(RolesUsuario)

	userRoles = append(userRoles, roles...)
	payload := &models.Payload{
		Role:               userRoles,
		DocumentoCompuesto: documentoCompuesto,
		Documento:          documento,
		Email:              mail,
		FamilyName:         familyName,
	}

	EstudianteInfo, err := GetCodeByEmailStudentService(mail)
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
}

func MapAtributos(RolesUsuario models.AtributosToken) (string, string, string, string, []string) {
	var familyName, documento, mail, documentoCompuesto string
	var roles []string

	for _, v := range RolesUsuario.Usuario.Atributos {
		switch v.Atributo {
		case "role":
			roles = strings.Split(v.Valor, ",")
		case "sn":
			familyName = v.Valor
		case "documento":
			documento = v.Valor
		case "documento_compuesto":
			documentoCompuesto = v.Valor
		case "mail":
			mail = v.Valor
		}
	}

	return familyName, documento, mail, documentoCompuesto, roles
}

func ClientCredentialsRequest(payload string) (response models.ClientCredentialsResponse, err error) {
	request := httplib.Post(beego.AppConfig.String("Wso2AuthService") + "apioas/token")
	request.Header("Content-Type", "application/x-www-form-urlencoded")
	request.Header("Authorization", "Basic "+payload)
	request.Body("grant_type=client_credentials")
	err = request.ToJSON(&response)

	return
}
