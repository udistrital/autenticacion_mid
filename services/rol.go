package services

import (
	"context"
	"fmt"
	"slices"

	"github.com/udistrital/autenticacion_mid/helpers"
	"github.com/udistrital/autenticacion_mid/models"
)

// AddRol adds the rol to the user in IS 7.3 and returns the user's roles.
// yaAsignado is true, and nothing changes, when the user already has it.
func AddRol(ctx context.Context, user models.UpdateRol) (roles []string, yaAsignado bool, err error) {
	usuario, err := helpers.GetGruposUsuarioIS(ctx, user.User)
	if err != nil {
		return nil, false, err
	}
	if usuario == nil {
		return nil, false, fmt.Errorf("El usuario %s %w", user.User, helpers.ErrNoExiste)
	}

	if slices.Contains(usuario.Roles, user.Rol) {
		return usuario.Roles, true, nil
	}

	if err := helpers.AgregarUsuarioGrupoIS(ctx, usuario.ID, user.Rol); err != nil {
		return nil, false, err
	}

	return append(usuario.Roles, user.Rol), false, nil
}

// RemoveRol removes the rol from the user in IS 7.3 and returns the user's
// roles. noAsignado is true, and nothing changes, when the user doesn't have it.
func RemoveRol(ctx context.Context, user models.UpdateRol) (roles []string, noAsignado bool, err error) {
	usuario, err := helpers.GetGruposUsuarioIS(ctx, user.User)
	if err != nil {
		return nil, false, err
	}
	if usuario == nil {
		return nil, false, fmt.Errorf("El usuario %s %w", user.User, helpers.ErrNoExiste)
	}

	if !slices.Contains(usuario.Roles, user.Rol) {
		return usuario.Roles, true, nil
	}

	if err := helpers.QuitarUsuarioGrupoIS(ctx, usuario.ID, user.Rol); err != nil {
		return nil, false, err
	}

	return slices.DeleteFunc(usuario.Roles, func(r string) bool { return r == user.Rol }), false, nil
}
