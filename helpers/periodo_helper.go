package helpers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/autenticacion_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func GetPeriodoUsuario(ctx context.Context, documento string, query map[string]string, limit int64, offset int64) (models.MetadataResponse, error) {
	var PeriodoUsuario models.MetadataResponse
	var queryParams []string
	for k, v := range query {
		queryParams = append(queryParams, k+"="+v)
	}
	queryString := strings.Join(queryParams, "&")

	urlGetPeriodoUsuario := beego.AppConfig.String("HistoricoRolesCrudService") + "usuarios/" + documento + "/periodos?" + queryString + "&limit=" + strconv.Itoa(int(limit)) + "&offset=" + strconv.Itoa(int(offset))

	_, err := request.GetWithContext(ctx, urlGetPeriodoUsuario, &PeriodoUsuario)
	if err != nil {

		return models.MetadataResponse{}, fmt.Errorf("error al obtener los roles del usuario %s", documento)
	}

	return PeriodoUsuario, nil
}

func GetTerceroInfo(ctx context.Context, documento string) (models.TerceroInfo, error) {
	var TerceroInfo []models.TerceroInfo
	urlGetTerceroInfo := beego.AppConfig.String("TercerosService") + "tercero/identificacion?query=" + documento

	_, err := request.GetWithContext(ctx, urlGetTerceroInfo, &TerceroInfo)
	if err != nil {
		return models.TerceroInfo{}, fmt.Errorf("error al obtener la información del tercero con documento %s", documento)
	}

	if len(TerceroInfo) == 0 {
		return models.TerceroInfo{}, fmt.Errorf("no se encontró información del tercero con documento %s", documento)
	}

	return TerceroInfo[0], nil
}

func GetAllPeriodos(ctx context.Context, query map[string]string, limit int64, offset int64) (models.MetadataResponse, error) {
	var response models.MetadataResponse
	var queryParams []string
	for k, v := range query {
		queryParams = append(queryParams, k+"="+v)
	}
	queryString := strings.Join(queryParams, "&")

	urlGetPeriodos := beego.AppConfig.String("HistoricoRolesCrudService") +
		"periodos-rol-usuarios?" + queryString + "&limit=" + strconv.Itoa(int(limit)) + "&offset=" + strconv.Itoa(int(offset))

	_, err := request.GetWithContext(ctx, urlGetPeriodos, &response)
	if err != nil {
		log.Println("Error en la solicitud HTTP:", err)
		return models.MetadataResponse{}, fmt.Errorf("error al obtener los periodos")
	}

	return response, nil
}
