package controllers

import (
	"net/url"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/stretchr/testify/assert"
	"github.com/udistrital/autenticacion_mid/models"
)

func TestParamsPeriodos(t *testing.T) {
	for consulta, valido := range map[string]bool{
		"":                             true,
		"query=activo:true":            true,
		"query=activo:true,rol:1":      true,
		"query=fecha:2026-01-01T10:00": true,
		"query=activo":                 false,
		"query=activo:true,":           false,
		"query=:true":                  false,
		"limit=10&offset=5":            true,
		"limit=diez":                   false,
		"offset=x":                     false,
	} {
		form, err := url.ParseQuery(consulta)
		assert.NoError(t, err)
		var p models.ParamsPeriodos
		ok := beego.ParseForm(form, &p) == nil
		if ok {
			valid := validation.Validation{}
			ok, err = valid.Valid(&p)
			assert.NoError(t, err, consulta)
		}
		assert.Equal(t, valido, ok, consulta)
	}
}
