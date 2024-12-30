package helpers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/udistrital/autenticacion_mid/models"
)

func GetPeriodoUsuario(documento string, query map[string]string, limit int64, offset int64) (models.MetadataResponse, error) {
	var PeriodoUsuario models.MetadataResponse
	var queryParams []string
	for k, v := range query {
		queryParams = append(queryParams, k+"="+v)
	}
	queryString := strings.Join(queryParams, "&")

	urlGetPeriodoUsuario := httplib.Get(beego.AppConfig.String("HistoricoRolesCrudService") + "usuarios/" + documento + "/periodos?" + queryString + "&limit=" + strconv.Itoa(int(limit)) + "&offset=" + strconv.Itoa(int(offset)))
	urlGetPeriodoUsuario.Header("Accept", "application/json")

	err := urlGetPeriodoUsuario.ToJSON(&PeriodoUsuario)
	if err != nil {

		return models.MetadataResponse{}, fmt.Errorf("error al obtener los roles del usuario %s", documento)
	}

	return PeriodoUsuario, nil
}

func GetTerceroInfo(documento string) (models.TerceroInfo, error) {
	var TerceroInfo []models.TerceroInfo
	urlGetTerceroInfo := httplib.Get(beego.AppConfig.String("TercerosService") + "tercero/identificacion?query=" + documento)

	urlGetTerceroInfo.Header("Accept", "application/json")

	err := urlGetTerceroInfo.ToJSON(&TerceroInfo)
	if err != nil {
		return models.TerceroInfo{}, fmt.Errorf("error al obtener la información del tercero con documento %s", documento)
	}

	if len(TerceroInfo) == 0 {
		return models.TerceroInfo{}, fmt.Errorf("no se encontró información del tercero con documento %s", documento)
	}

	return TerceroInfo[0], nil
}

func GetAllPeriodos(query map[string]string, limit int64, offset int64) (models.MetadataResponse, error) {
	var response models.MetadataResponse
	var queryParams []string
	for k, v := range query {
		queryParams = append(queryParams, k+"="+v)
	}
	queryString := strings.Join(queryParams, "&")

	urlGetPeriodos := httplib.Get(beego.AppConfig.String("HistoricoRolesCrudService") +
		"periodos-rol-usuarios?" + queryString + "&limit=" + strconv.Itoa(int(limit)) + "&offset=" + strconv.Itoa(int(offset)))
	urlGetPeriodos.Header("Accept", "application/json")

	err := urlGetPeriodos.ToJSON(&response)
	if err != nil {
		log.Println("Error en la solicitud HTTP:", err)
		return models.MetadataResponse{}, fmt.Errorf("error al obtener los periodos")
	}

	return response, nil
}
