package services

import (
	"fmt"

	"github.com/udistrital/autenticacion_mid/helpers"
	"github.com/udistrital/autenticacion_mid/models"
)

func GetPeriodoInfo(documento string, query map[string]string, limit int64, offset int64) (map[string]any, error) {

	infoDocumento, err := helpers.GetInfoByDocumentoService(documento)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener la información del documento: %v", err)
	}

	var correo string
	if infoDocumento.Usuario != nil {
		_, _, correo, _, _ = helpers.MapAtributos(infoDocumento)
		if err != nil {
			return nil, fmt.Errorf("Error al obtener la información del correo: %v", err)
		}

	}

	periodoUsuario, err := helpers.GetPeriodoUsuario(documento, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener los periodos del usuario: %v", err)
	}

	terceroInfo, err := helpers.GetTerceroInfo(documento)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener la información del tercero: %v", err)
	}

	var periodosRolUsuario []models.PeriodoRolUsuario
	for _, periodo := range periodoUsuario.Data {
		periodoRolUsuario := models.PeriodoRolUsuario{
			Nombre:       terceroInfo.Tercero.NombreCompleto,
			Documento:    terceroInfo.Identificacion.Numero,
			Correo:       correo,
			RolUsuario:   periodo.Rol.Nombre,
			Estado:       periodo.Activo,
			FechaInicial: periodo.FechaInicio,
			FechaFinal:   periodo.FechaFin,
			Finalizado:   periodo.Finalizado,
			IdPeriodo:    int(periodo.Id),
			IdTercero:    int(terceroInfo.Tercero.Id),
		}
		periodosRolUsuario = append(periodosRolUsuario, periodoRolUsuario)
	}

	response := map[string]any{
		"Data":     periodosRolUsuario,
		"Metadata": periodoUsuario.Metadata,
	}

	return response, nil
}

func GetAllPeriodosRoles(query map[string]string, limit int64, offset int64) (map[string]any, error) {
	var response map[string]any
	periodosResponse, err := helpers.GetAllPeriodos(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener todos los periodos: %v", err)
	}

	var periodoRolUsuario []models.PeriodoRolUsuario
	var errores []string

	for _, periodos := range periodosResponse.Data {
		terceroInfo, err := helpers.GetTerceroInfo(periodos.Usuario.Documento)
		if err != nil {

			periodoRolUsuario = append(periodoRolUsuario, models.PeriodoRolUsuario{
				Nombre:       "No encontrado",
				Documento:    periodos.Usuario.Documento,
				Correo:       "No encontrado",
				RolUsuario:   periodos.Rol.Nombre,
				Estado:       periodos.Activo,
				FechaInicial: periodos.FechaInicio,
				FechaFinal:   periodos.FechaFin,
				Finalizado:   periodos.Finalizado,
				IdPeriodo:    int(periodos.Id),
				IdTercero:    int(terceroInfo.Tercero.Id),
			})

			errores = append(errores, fmt.Sprintf("Error al obtener la información del tercero con documento %s ", periodos.Usuario.Documento))
			continue
		}

		infoDocumento, err := helpers.GetInfoByDocumentoService(periodos.Usuario.Documento)
		if err != nil {
			errores = append(errores, fmt.Sprintf("Error al obtener la información del documento  %s ", periodos.Usuario.Documento))
			continue
		}
		var correo string
		if infoDocumento.Usuario != nil {
			_, _, correo, _, _ = helpers.MapAtributos(infoDocumento)
			if err != nil {
				return nil, fmt.Errorf("Error al obtener la información del correo: %v", err)
			}
		}

		periodoRolUsuario = append(periodoRolUsuario, models.PeriodoRolUsuario{
			Nombre:       terceroInfo.Tercero.NombreCompleto,
			Documento:    terceroInfo.Identificacion.Numero,
			Correo:       correo,
			RolUsuario:   periodos.Rol.Nombre,
			Estado:       periodos.Activo,
			FechaInicial: periodos.FechaInicio,
			FechaFinal:   periodos.FechaFin,
			Finalizado:   periodos.Finalizado,
			IdPeriodo:    int(periodos.Id),
			IdTercero:    int(terceroInfo.Tercero.Id),
		})
	}

	if len(errores) > 0 {
		return map[string]any{
			"Data":     periodoRolUsuario,
			"Metadata": periodosResponse.Metadata,
			"Errores":    errores,
		}, nil
	}

	response = map[string]interface{}{
		"Data":     periodoRolUsuario,
		"Metadata": periodosResponse.Metadata,
	}

	return response, nil
}
