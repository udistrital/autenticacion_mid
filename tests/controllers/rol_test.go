package controllers_test

import (
	"encoding/json"
	"net/http"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/udistrital/autenticacion_mid/models"
)

func TestAddRol(t *testing.T) {
	conRol(t, false)
	resp := postJSON(t, "/rol/add", models.UpdateRol{User: usuarioPrueba, Rol: rolPrueba})
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRemoveRol(t *testing.T) {
	conRol(t, true)
	resp := postJSON(t, "/rol/remove", models.UpdateRol{User: usuarioPrueba, Rol: rolPrueba})
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// conRol deja al usuario de prueba con el rol, o sin él, durante la prueba, y al
// terminar le devuelve el estado que tenía.
func conRol(t *testing.T, tener bool) {
	t.Helper()
	tenia := tieneRol(t)
	if tenia != tener {
		cambiarRol(t, tener)
	}
	t.Cleanup(func() {
		if tieneRol(t) != tenia {
			cambiarRol(t, tenia)
		}
	})
}

// tieneRol indica si el usuario de prueba tiene el rol, según /token/userRol.
func tieneRol(t *testing.T) bool {
	t.Helper()
	resp := postJSON(t, "/token/userRol", models.UserName{User: usuarioPrueba})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var info models.InfoUsuario
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&info))
	return slices.Contains(info.Role, rolPrueba)
}

func cambiarRol(t *testing.T, tener bool) {
	t.Helper()
	ruta := "/rol/remove"
	if tener {
		ruta = "/rol/add"
	}
	resp := postJSON(t, ruta, models.UpdateRol{User: usuarioPrueba, Rol: rolPrueba})
	require.Equal(t, http.StatusOK, resp.StatusCode, "%s para preparar la prueba", ruta)
}
