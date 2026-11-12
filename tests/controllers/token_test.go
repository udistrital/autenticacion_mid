package controllers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/udistrital/autenticacion_mid/models"
)

func TestEmailToken(t *testing.T) {
	resp := postJSON(t, "/token/emailToken", models.Token{Email: correoEstudiante})
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserRol(t *testing.T) {
	resp := postJSON(t, "/token/userRol", models.UserName{User: usuarioPrueba})
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDocumentoToken(t *testing.T) {
	resp := postJSON(t, "/token/documentoToken", models.Documento{Numero: documentoPrueba})
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]any
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
}
