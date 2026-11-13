// @APIVersion 0.2.1
// @Title autenticacion_mid
// @Description API MID de autenticación: entrega los datos y roles de los usuarios para sus tokens, y gestiona sus roles y periodos de roles.
// @License GNU GPL v3
// @LicenseUrl https://www.gnu.org/licenses/gpl-3.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/autenticacion_mid/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/token",
			beego.NSInclude(
				&controllers.TokenController{},
			),
		),
		beego.NSNamespace("/rol",
			beego.NSInclude(
				&controllers.RolController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
