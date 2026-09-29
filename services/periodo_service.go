package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/udistrital/autenticacion_mid/helpers"
	"github.com/udistrital/autenticacion_mid/models"
)

// maxConsultasConcurrentes limita las llamadas simultáneas a Terceros y WSO2
const maxConsultasConcurrentes = 10

func GetPeriodoInfo(ctx context.Context, documento string, query map[string]string, limit int64, offset int64) (map[string]any, error) {

	infoDocumento, err := helpers.GetInfoByDocumentoService(ctx, documento)
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

	periodoUsuario, err := helpers.GetPeriodoUsuario(ctx, documento, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener los periodos del usuario: %v", err)
	}

	terceroInfo, err := helpers.GetTerceroInfo(ctx, documento)
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

// consultarPersonas consulta Terceros y WSO2 una sola vez por documento, en paralelo.
// WSO2 solo se consulta si el tercero existe, igual que en el flujo secuencial original.
func consultarPersonas(ctx context.Context, documentos []string) map[string]models.InfoPersona {
	personas := make(map[string]models.InfoPersona, len(documentos))
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConsultasConcurrentes)

	for _, documento := range documentos {
		wg.Add(1)
		go func(documento string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			var info models.InfoPersona
			info.Tercero, info.ErrTercero = helpers.GetTerceroInfo(ctx, documento)
			if info.ErrTercero == nil {
				var infoDocumento models.AtributosToken
				infoDocumento, info.ErrWso2 = helpers.GetInfoByDocumentoService(ctx, documento)
				if info.ErrWso2 == nil && infoDocumento.Usuario != nil {
					_, _, info.Correo, _, _ = helpers.MapAtributos(infoDocumento)
				}
			}

			mu.Lock()
			personas[documento] = info
			mu.Unlock()
		}(documento)
	}
	wg.Wait()

	return personas
}

func GetAllPeriodosRoles(ctx context.Context, query map[string]string, limit int64, offset int64) (map[string]any, error) {
	var response map[string]any
	periodosResponse, err := helpers.GetAllPeriodos(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener todos los periodos: %v", err)
	}

	// Un usuario puede tener varios periodos: se consulta una vez por documento
	var documentos []string
	vistos := make(map[string]bool)
	for _, periodos := range periodosResponse.Data {
		if !vistos[periodos.Usuario.Documento] {
			vistos[periodos.Usuario.Documento] = true
			documentos = append(documentos, periodos.Usuario.Documento)
		}
	}
	personas := consultarPersonas(ctx, documentos)

	var periodoRolUsuario []models.PeriodoRolUsuario
	var errores []string

	for _, periodos := range periodosResponse.Data {
		persona := personas[periodos.Usuario.Documento]
		if persona.ErrTercero != nil {

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
				IdTercero:    int(persona.Tercero.Tercero.Id),
			})

			errores = append(errores, fmt.Sprintf("Error al obtener la información del tercero con documento %s ", periodos.Usuario.Documento))
			continue
		}

		if persona.ErrWso2 != nil {
			errores = append(errores, fmt.Sprintf("Error al obtener la información del documento  %s ", periodos.Usuario.Documento))
			continue
		}

		terceroInfo := persona.Tercero
		periodoRolUsuario = append(periodoRolUsuario, models.PeriodoRolUsuario{
			Nombre:       terceroInfo.Tercero.NombreCompleto,
			Documento:    terceroInfo.Identificacion.Numero,
			Correo:       persona.Correo,
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
			"Errores":  errores,
		}, nil
	}

	response = map[string]interface{}{
		"Data":     periodoRolUsuario,
		"Metadata": periodosResponse.Metadata,
	}

	return response, nil
}
