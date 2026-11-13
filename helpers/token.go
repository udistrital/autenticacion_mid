package helpers

import (
	"context"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/autenticacion_mid/models"
	"github.com/udistrital/utils_oas/request"
)

var getCodeByEmailStudentService = beego.AppConfig.String("GetCodeByEmailStudentService")

func GetCodeByEmailStudentService(ctx context.Context, email string) (models.EstudianteInfo, error) {
	var EstudianteInfo models.EstudianteInfo
	urlGetCodeByEmailStudentService := getCodeByEmailStudentService + email
	_, err := request.GetWithContext(ctx, urlGetCodeByEmailStudentService, &EstudianteInfo)

	if err != nil {
		logs.Info("Error", err)
		return models.EstudianteInfo{}, fmt.Errorf("error al obtener el código del estudiante %s", email)
	}

	return EstudianteInfo, nil
}

// GetEstudiante devuelve el estudiante con el correo, o nil si no hay ninguno.
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
