package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/autenticacion_mid/controllers:RolController"] = append(beego.GlobalControllerRouter["github.com/udistrital/autenticacion_mid/controllers:RolController"],
        beego.ControllerComments{
            Method: "AddRol",
            Router: "/add",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/autenticacion_mid/controllers:RolController"] = append(beego.GlobalControllerRouter["github.com/udistrital/autenticacion_mid/controllers:RolController"],
        beego.ControllerComments{
            Method: "GetAllPeriodos",
            Router: "/periods",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/autenticacion_mid/controllers:RolController"] = append(beego.GlobalControllerRouter["github.com/udistrital/autenticacion_mid/controllers:RolController"],
        beego.ControllerComments{
            Method: "RemoveRol",
            Router: "/remove",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/autenticacion_mid/controllers:RolController"] = append(beego.GlobalControllerRouter["github.com/udistrital/autenticacion_mid/controllers:RolController"],
        beego.ControllerComments{
            Method: "GetPeriodoInfo",
            Router: "/user/:documento/periods",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/autenticacion_mid/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/udistrital/autenticacion_mid/controllers:TokenController"],
        beego.ControllerComments{
            Method: "GetDocumento",
            Router: "/documentoToken",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/autenticacion_mid/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/udistrital/autenticacion_mid/controllers:TokenController"],
        beego.ControllerComments{
            Method: "GetEmail",
            Router: "/emailToken",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/autenticacion_mid/controllers:TokenController"] = append(beego.GlobalControllerRouter["github.com/udistrital/autenticacion_mid/controllers:TokenController"],
        beego.ControllerComments{
            Method: "GetRol",
            Router: "/userRol",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
