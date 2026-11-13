package services

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/autenticacion_mid/scim"
)

var ErrGrupoNoGestionable = errors.New("el grupo no se puede gestionar")

func AgregarGrupo(ctx context.Context, usuario, grupo string) (u *scim.Usuario, yaAsignado bool, err error) {
	if esRolInternoOAplicacion(grupo) {
		return nil, false, fmt.Errorf("%w: %s", ErrGrupoNoGestionable, grupo)
	}

	u, err = scim.GetGruposUsuario(ctx, usuario)
	if err != nil {
		return nil, false, err
	}
	if u == nil {
		return nil, false, fmt.Errorf("%w: %s", scim.ErrUsuarioNoExiste, usuario)
	}
	u.Grupos = filtrarGruposConRoles(ctx, u.Grupos)
	conRoles, err := grupoConRoles(ctx, grupo)
	if err != nil {
		return nil, false, err
	}
	if conRoles {
		return nil, false, fmt.Errorf("%w: %s tiene roles", ErrGrupoNoGestionable, grupo)
	}

	if slices.Contains(u.Grupos, grupo) {
		return u, true, nil
	}

	if err := scim.AgregarUsuarioGrupo(ctx, u, grupo); err != nil {
		return nil, false, err
	}

	u.Grupos = append(u.Grupos, grupo)
	return u, false, nil
}

func QuitarGrupo(ctx context.Context, usuario, grupo string) (u *scim.Usuario, noAsignado bool, err error) {
	if esRolInternoOAplicacion(grupo) {
		return nil, false, fmt.Errorf("%w: %s", ErrGrupoNoGestionable, grupo)
	}

	u, err = scim.GetGruposUsuario(ctx, usuario)
	if err != nil {
		return nil, false, err
	}
	if u == nil {
		return nil, false, fmt.Errorf("%w: %s", scim.ErrUsuarioNoExiste, usuario)
	}
	u.Grupos = filtrarGruposConRoles(ctx, u.Grupos)
	conRoles, err := grupoConRoles(ctx, grupo)
	if err != nil {
		return nil, false, err
	}
	if conRoles {
		return nil, false, fmt.Errorf("%w: %s tiene roles", ErrGrupoNoGestionable, grupo)
	}

	if !slices.Contains(u.Grupos, grupo) {
		return u, true, nil
	}

	if err := scim.QuitarUsuarioGrupo(ctx, u.ID, grupo); err != nil {
		return nil, false, err
	}

	u.Grupos = slices.DeleteFunc(u.Grupos, func(g string) bool { return g == grupo })
	return u, false, nil
}

// duracionCacheGrupos es el tiempo que se guarda si un grupo tiene roles.
const duracionCacheGrupos = 10 * time.Minute

var cacheGrupos = nuevoCacheGrupos()

// consultarGrupoConRoles consulta en IS si un grupo tiene roles; las pruebas la
// reemplazan.
var consultarGrupoConRoles = scim.GrupoConRoles

func nuevoCacheGrupos() cache.Cache {
	c, err := cache.NewCache("memory", `{"interval":180}`)
	if err != nil {
		panic(fmt.Sprintf("no se pudo crear el cache de grupos: %v", err))
	}
	return c
}

// grupoConRoles indica si el grupo tiene roles asignados, y guarda la respuesta
// de IS durante duracionCacheGrupos. Los errores no se guardan.
func grupoConRoles(ctx context.Context, grupo string) (bool, error) {
	if conRoles, ok := cacheGrupos.Get(grupo).(bool); ok {
		return conRoles, nil
	}

	conRoles, err := consultarGrupoConRoles(ctx, grupo)
	if err != nil {
		return false, err
	}

	if err := cacheGrupos.Put(grupo, conRoles, duracionCacheGrupos); err != nil {
		logs.Warning("no se pudo guardar en el cache el grupo %s: %v", grupo, err)
	}
	return conRoles, nil
}

// filtrarGruposConRoles quita los roles internos o de aplicación y los grupos que
// tienen roles asignados: esos accesos se administran en IS y no son roles de las
// aplicaciones. También quita, y registra, los grupos que no se pueden consultar.
func filtrarGruposConRoles(ctx context.Context, grupos []string) []string {
	var sinRoles []string
	for _, g := range grupos {
		if esRolInternoOAplicacion(g) {
			continue
		}
		conRoles, err := grupoConRoles(ctx, g)
		if err != nil {
			logs.Error("se omite el grupo %s: %v", g, err)
			continue
		}
		if !conRoles {
			sinRoles = append(sinRoles, g)
		}
	}
	return sinRoles
}

func esRolInternoOAplicacion(nombre string) bool {
	for _, prefijo := range []string{"Internal/", "Application/"} {
		if len(nombre) >= len(prefijo) && strings.EqualFold(nombre[:len(prefijo)], prefijo) {
			return true
		}
	}
	return false
}
