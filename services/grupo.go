package services

import (
	"context"
	"fmt"
	"slices"

	"github.com/udistrital/autenticacion_mid/helpers"
	"github.com/udistrital/autenticacion_mid/models"
)

func AgregarGrupo(ctx context.Context, usuario, grupo string) (u *models.UsuarioIS, yaAsignado bool, err error) {
	if !helpers.GrupoGestionable(grupo) {
		return nil, false, fmt.Errorf("%w: %s", helpers.ErrGrupoNoGestionable, grupo)
	}

	u, err = helpers.GetGruposUsuarioIS(ctx, usuario)
	if err != nil {
		return nil, false, err
	}
	if u == nil {
		return nil, false, fmt.Errorf("%w: %s", helpers.ErrUsuarioNoExiste, usuario)
	}

	if slices.Contains(u.Grupos, grupo) {
		return u, true, nil
	}

	if err := helpers.AgregarUsuarioGrupoIS(ctx, u, grupo); err != nil {
		return nil, false, err
	}

	u.Grupos = append(u.Grupos, grupo)
	return u, false, nil
}

func QuitarGrupo(ctx context.Context, usuario, grupo string) (u *models.UsuarioIS, noAsignado bool, err error) {
	if !helpers.GrupoGestionable(grupo) {
		return nil, false, fmt.Errorf("%w: %s", helpers.ErrGrupoNoGestionable, grupo)
	}

	u, err = helpers.GetGruposUsuarioIS(ctx, usuario)
	if err != nil {
		return nil, false, err
	}
	if u == nil {
		return nil, false, fmt.Errorf("%w: %s", helpers.ErrUsuarioNoExiste, usuario)
	}

	if !slices.Contains(u.Grupos, grupo) {
		return u, true, nil
	}

	if err := helpers.QuitarUsuarioGrupoIS(ctx, u.ID, grupo); err != nil {
		return nil, false, err
	}

	u.Grupos = slices.DeleteFunc(u.Grupos, func(g string) bool { return g == grupo })
	return u, false, nil
}
