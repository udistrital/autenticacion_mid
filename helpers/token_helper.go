package helpers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/udistrital/autenticacion_mid/models"
)

// Get_Methods

func GetRolesUsuario(email string) (models.AtributosToken, error) {
	var RolesUsuario models.AtributosToken
	urlGetRolesUsuario := httplib.Get(beego.AppConfig.String("Wso2Service") + "roles?usuario=" + email)
	urlGetRolesUsuario.Header("Accept", "application/json")

	err := urlGetRolesUsuario.ToJSON(&RolesUsuario)
	if err != nil {
		return models.AtributosToken{}, fmt.Errorf("Error al obtener los roles del usuario %s", email)
	}

	return RolesUsuario, nil
}

func GetCodeByEmailStudentService(email string) (models.EstudianteInfo, error) {
	var EstudianteInfo models.EstudianteInfo
	urlGetCodeByEmailStudentService := httplib.Get(beego.AppConfig.String("GetCodeByEmailStudentService") + email)
	urlGetCodeByEmailStudentService.Header("Accept", "application/json")

	err := urlGetCodeByEmailStudentService.ToJSON(&EstudianteInfo)
	if err != nil {
		return models.EstudianteInfo{}, fmt.Errorf("Error al obtener el código del estudiante %s", email)
	}

	return EstudianteInfo, nil
}
