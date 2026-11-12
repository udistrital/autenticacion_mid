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
var getCodeByEmailStudentService = beego.AppConfig.String("GetCodeByEmailStudentService")

func GetCodeByEmailStudentService(ctx context.Context, email string) (models.EstudianteInfo, error) {
	var EstudianteInfo models.EstudianteInfo
	urlGetCodeByEmailStudentService := getCodeByEmailStudentService + email
	_, err := request.GetWithContext(ctx, urlGetCodeByEmailStudentService, &EstudianteInfo)

	if err != nil {
		logs.Info("Error", err)
		return models.EstudianteInfo{}, fmt.Errorf("Error al obtener el código del estudiante %s", email)
	}

	return EstudianteInfo, nil
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

func ClientCredentialsRequest(payload string) (models.ClientCredentialsResponse, error) {
	form := url.Values{"grant_type": {"client_credentials"}}
	req, err := http.NewRequest(http.MethodPost, wso2ISBasePath+"/oauth2/token",
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
