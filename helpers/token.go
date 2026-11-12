package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/autenticacion_mid/models"
	"github.com/udistrital/utils_oas/request"
)

var clientCredentialsClient = &http.Client{Timeout: 30 * time.Second}

func GetCodeByEmailStudentService(ctx context.Context, email string) (models.EstudianteInfo, error) {
	var EstudianteInfo models.EstudianteInfo
	urlGetCodeByEmailStudentService := beego.AppConfig.String("GetCodeByEmailStudentService") + email
	_, err := request.GetWithContext(ctx, urlGetCodeByEmailStudentService, &EstudianteInfo)

	if err != nil {
		logs.Info("Error", err)
		return models.EstudianteInfo{}, fmt.Errorf("Error al obtener el código del estudiante %s", email)
	}

	return EstudianteInfo, nil
}

// GetInfoByDocumentoService reads the user from EI (documentoToken, periods).
func GetInfoByDocumentoService(ctx context.Context, documento string) (models.AtributosToken, error) {
	var RolesUsuario models.AtributosToken
	urlGetInfoByDocumentoService := beego.AppConfig.String("Wso2Service") + "usuario_documento?documento=" + documento

	_, err := request.GetWithContext(ctx, urlGetInfoByDocumentoService, &RolesUsuario)
	if err != nil {
		return models.AtributosToken{}, fmt.Errorf("Error al obtener la información del documento %s", documento)
	}

	return RolesUsuario, nil

}

// GetEstudiante returns the student with the email, or nil if there is none.
func GetEstudiante(ctx context.Context, email string) (*models.Estudiante, error) {
	EstudianteInfo, err := GetCodeByEmailStudentService(ctx, email)
	if err != nil {
		return nil, err
	}

	if len(EstudianteInfo.EstudianteCollection.Estudiante) == 0 {
		return nil, nil
	}

	return &EstudianteInfo.EstudianteCollection.Estudiante[0], nil
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

func ClientCredentialsRequest(payload string) (models.ClientCredentialsResponse, error) {
	form := url.Values{"grant_type": {"client_credentials"}}
	req, err := http.NewRequest(http.MethodPost, beego.AppConfig.String("Wso2IsService")+"/oauth2/token",
		strings.NewReader(form.Encode()))
	if err != nil {
		return models.ClientCredentialsResponse{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+payload)

	resp, err := clientCredentialsClient.Do(req)
	if err != nil {
		return models.ClientCredentialsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.ClientCredentialsResponse{}, fmt.Errorf("IS respondió HTTP %d al pedir el token", resp.StatusCode)
	}

	var response models.ClientCredentialsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	return response, err
}
