package helpers

import (
	"errors"
	"slices"
	"strings"
)

var ErrGrupoNoGestionable = errors.New("el grupo no se puede gestionar")

var gruposIS = []string{"admin"}

func GrupoGestionable(grupo string) bool {
	return !esRolIS(grupo) && !esGrupoIS(grupo)
}

func esRolIS(nombre string) bool {
	for _, prefijo := range []string{"Internal/", "Application/"} {
		if len(nombre) >= len(prefijo) && strings.EqualFold(nombre[:len(prefijo)], prefijo) {
			return true
		}
	}
	return false
}

func esGrupoIS(grupo string) bool {
	return slices.ContainsFunc(gruposIS, func(g string) bool { return strings.EqualFold(g, grupo) })
}
