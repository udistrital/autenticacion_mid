package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrupoNoGestionable(t *testing.T) {
	_, _, err := AgregarGrupo(context.Background(), "user@udistrital.edu.co", "Internal/otro")
	assert.ErrorIs(t, err, ErrGrupoNoGestionable)

	_, _, err = QuitarGrupo(context.Background(), "user@udistrital.edu.co", "Internal/otro")
	assert.ErrorIs(t, err, ErrGrupoNoGestionable)
}

func TestEsRolInternoOAplicacion(t *testing.T) {
	for nombre, esRol := range map[string]bool{
		"Internal/everyone":   true,
		"Internal/selfsignup": true,
		"internal/otro":       true,
		"Application/app":     true,
		"Internal":            false,
		"GRUPO":               false,
	} {
		assert.Equal(t, esRol, esRolInternoOAplicacion(nombre), nombre)
	}
}

// conGruposFalsos reemplaza la consulta a IS: los grupos de conRoles tienen
// roles, los de sinRoles no y los demás fallan. Devuelve las consultas hechas.
func conGruposFalsos(t *testing.T, conRoles, sinRoles []string) map[string]int {
	t.Helper()
	consultas := map[string]int{}
	anterior := consultarGrupoConRoles
	consultarGrupoConRoles = func(_ context.Context, grupo string) (bool, error) {
		consultas[grupo]++
		for _, g := range conRoles {
			if g == grupo {
				return true, nil
			}
		}
		for _, g := range sinRoles {
			if g == grupo {
				return false, nil
			}
		}
		return false, errors.New("el grupo no existe")
	}
	require.NoError(t, cacheGrupos.ClearAll())
	t.Cleanup(func() { consultarGrupoConRoles = anterior })
	return consultas
}

// TestGrupoConRolesCache comprueba que la respuesta de IS se guarda, y que los
// errores no.
func TestGrupoConRolesCache(t *testing.T) {
	consultas := conGruposFalsos(t, []string{"GRUPO_CON_ROLES"}, []string{"GRUPO"})
	ctx := context.Background()

	for range 2 {
		conRoles, err := grupoConRoles(ctx, "GRUPO_CON_ROLES")
		require.NoError(t, err)
		assert.True(t, conRoles)

		conRoles, err = grupoConRoles(ctx, "GRUPO")
		require.NoError(t, err)
		assert.False(t, conRoles)

		_, err = grupoConRoles(ctx, "NO_EXISTE")
		assert.Error(t, err)
	}
	assert.Equal(t, map[string]int{"GRUPO_CON_ROLES": 1, "GRUPO": 1, "NO_EXISTE": 2}, consultas)
}

// TestFiltrarGruposConRoles comprueba que se quitan los roles internos o de
// aplicación, los grupos con roles y los que no se pueden consultar.
func TestFiltrarGruposConRoles(t *testing.T) {
	conGruposFalsos(t, []string{"GRUPO_CON_ROLES"}, []string{"GRUPO", "OTRO_GRUPO"})

	grupos := filtrarGruposConRoles(context.Background(), []string{"GRUPO", "Internal/everyone", "GRUPO_CON_ROLES", "NO_EXISTE", "Application/app", "OTRO_GRUPO"})
	assert.Equal(t, []string{"GRUPO", "OTRO_GRUPO"}, grupos)
}
