package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/udistrital/autenticacion_mid/helpers"
)

func TestGrupoNoGestionable(t *testing.T) {
	_, _, err := AgregarGrupo(context.Background(), "user@udistrital.edu.co", "Internal/admin")
	assert.ErrorIs(t, err, helpers.ErrGrupoNoGestionable)

	_, _, err = QuitarGrupo(context.Background(), "user@udistrital.edu.co", "Internal/admin")
	assert.ErrorIs(t, err, helpers.ErrGrupoNoGestionable)
}
