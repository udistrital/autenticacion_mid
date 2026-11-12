package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrupoGestionable(t *testing.T) {
	for grupo, gestionable := range map[string]bool{
		"PLANEACION":        true,
		"admin_campus":      true,
		"Internal/everyone": false,
		"admin":             false,
	} {
		assert.Equal(t, gestionable, GrupoGestionable(grupo), grupo)
	}
}

func TestEsRolIS(t *testing.T) {
	for nombre, esRol := range map[string]bool{
		"Internal/everyone":   true,
		"Internal/selfsignup": true,
		"internal/admin":      true,
		"Application/app":     true,
		"Internal":            false,
		"PLANEACION":          false,
	} {
		assert.Equal(t, esRol, esRolIS(nombre), nombre)
	}
}

func TestEsGrupoIS(t *testing.T) {
	for grupo, esGrupo := range map[string]bool{
		"admin":        true,
		"ADMIN":        true,
		"admin_campus": false,
		"PLANEACION":   false,
	} {
		assert.Equal(t, esGrupo, esGrupoIS(grupo), grupo)
	}
}
