package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

// Token structure
type Token struct {
	Email string `json:"email"`
}

// UserName structure
type UserName struct {
	User string `json:"user"`
}

//Payload structure
type Payload struct {
	Role               []string `json:"role"`
	Documento          string   `json:"documento"`
	DocumentoCompuesto string   `json:"documento_compuesto"`
	Email              string   `json:"email"`
	FamilyName         string   `json:"FamilyName"`
	Codigo             string   `json:"Codigo"`
	Estado             string   `json:"Estado"`
}

// EstudianteInfo structure
type EstudianteInfo struct {
	EstudianteCollection struct {
		Estudiante []struct {
			Codigo string `json:"codigo"`
			Estado string `json:"estado"`
		} `json:"estudiante"`
	} `json:"estudianteCollection"`
}

//RolesUsuario structure
type RolesUsuario struct {
	Usuario struct {
		Roles []struct {
			Rol string `json:"rol"`
		} `json:"Roles"`
	} `json:"Usuario"`
}

//AtributosToken structure
type AtributosToken struct {
	Usuario struct {
		Atributos []struct {
			Atributo string `json:"atributo"`
			Valor    string `json:"valor"`
		} `json:"Atributos"`
	} `json:"Usuario"`
}

// UserInfo structure
type UserInfo struct {
	Codigo string   `json:"Codigo"`
	Estado string   `json:"Estado"`
	Email  string   `json:"email"`
	Rol    []string `json:"rol"`
}

// returnInfo structure

type returnInfo struct {
	InfoUserInfo     UserInfo ``
	InfoRolesUsuario RolesUsuario
}

//UserId structure
type UserId struct {
	Usuarios struct {
		Usuario []struct {
			Id string `json:"um_id"`
		} `json:"usuario"`
	} `json:"Usuarios"`
}

type RolId struct {
	Roles struct {
		Rol []struct {
			Id int `json:"um_id"`
		} `json:"Rol"`
	} `json:"Roles"`
}

type UpdateRol struct {
	User string `json:"user"`
	Rol  string `json:"rol"`
}

type PostUsuarioRol struct {
	_post_usuario_rol struct {
		um_role_id int   `json:"um_role_id"`
		um_user_id int64 `json:"um_user_id"`
	} `json:"_post_usuario_rol"`
}

// GetInfoByEmail ...
func GetInfoByEmail(m *Token) (u *UserInfo, err error) {
	var estudianteInfo EstudianteInfo
	userRoles := []string{}
	r := httplib.Get(beego.AppConfig.String("GetCodeByEmailStudentService") + m.Email)
	r.Header("Accept", "application/json")
	if err = r.ToJSON(&estudianteInfo); err == nil {
		if estudianteInfo.EstudianteCollection.Estudiante != nil {
			userRoles = append(userRoles, "ESTUDIANTE")
			u := &UserInfo{
				Codigo: estudianteInfo.EstudianteCollection.Estudiante[0].Codigo,
				Estado: estudianteInfo.EstudianteCollection.Estudiante[0].Estado,
				Email:  m.Email,
				Rol:    userRoles,
			}
			return u, nil
		} else {
			beego.Info(err)
			return nil, errors.New("Email no registrado")
		}
	}
	fmt.Println(err)
	return nil, err
}

// GetRolesByUser ...
func GetRolesByUser(user UserName) (roles *Payload, outputError map[string]interface{}) {
	var RolesUsuario AtributosToken
	var estudianteInfo EstudianteInfo
	var familyName string
	var documento string
	var mail string
	var documentoCompuesto string
	userRoles := []string{}

	beego.Info("URL: ", beego.AppConfig.String("Wso2Service")+"roles/"+user.User)
	r := httplib.Get(beego.AppConfig.String("Wso2Service") + "roles/" + user.User)
	r.Header("Accept", "application/json")
	if err := r.ToJSON(&RolesUsuario); err == nil {
		if RolesUsuario.Usuario.Atributos != nil {
			for k, v := range RolesUsuario.Usuario.Atributos {
				switch v.Atributo {
				case "role":
					roles := strings.Split(v.Valor, ",")
					for _, v := range roles {
						userRoles = append(userRoles, v)
					}
				case "sn":
					familyName = v.Valor
				case "documento":
					documento = v.Valor

				case "documento_compuesto":
					documentoCompuesto = v.Valor
				case "mail":
					mail = v.Valor
				}

				fmt.Println(k, v)
			}
			payload := &Payload{
				Role:               userRoles,
				DocumentoCompuesto: documentoCompuesto,
				Documento:          documento,
				Email:              mail,
				FamilyName:         familyName,
			}

			r2 := httplib.Get(beego.AppConfig.String("GetCodeByEmailStudentService") + mail)
			r2.Header("Accept", "application/json")
			if err := r2.ToJSON(&estudianteInfo); err == nil {
				if estudianteInfo.EstudianteCollection.Estudiante != nil {
					userRoles = append(userRoles, "ESTUDIANTE")
					payload.Codigo = estudianteInfo.EstudianteCollection.Estudiante[0].Codigo
					payload.Estado = estudianteInfo.EstudianteCollection.Estudiante[0].Estado
					payload.Role = userRoles
				}
			}

			return payload, nil
		} else {
			outputError = map[string]interface{}{"Function": "FuncionalidadMidController:userRol", "Error": "usuario no registrado"}

			return nil, outputError
		}

	} else {
		outputError = map[string]interface{}{"Function": "FuncionalidadMidController:userRol", "Error": err}
		return nil, outputError
	}
}

// AddRol ...
func AddRol(user UpdateRol) (roles *Payload, outputError map[string]interface{}) {
	var Uid UserId
	var Rid RolId
	var notExist = true
	var userName UserName
	userName.User = user.User
	urlUsuario := beego.AppConfig.String("Wso2Service") + "usuario/"
	urlRol := beego.AppConfig.String("Wso2Service") + "rol/"
	postRolSuario := beego.AppConfig.String("Wso2Service") + "usuario_rol/"

	if responseData, err := GetRolesByUser(userName); err == nil {
		for i := range responseData.Role {
			if responseData.Role[i] == user.Rol {
				beego.Info("YA EXISTE: ", responseData.Role[i])
				notExist = false
			}
		}
	}
	if notExist {
		var m PostUsuarioRol
		var res map[string]interface{}
		beego.Info("URL: ", urlUsuario+user.User)
		requestUsuario := httplib.Get(urlUsuario + user.User)
		requestUsuario.Header("Accept", "application/json")
		if err := requestUsuario.ToJSON(&Uid); err == nil {
			if len(Uid.Usuarios.Usuario) > 0 {
				beego.Info("User: ", (Uid.Usuarios.Usuario))
				beego.Info("URL: ", urlRol+user.Rol)
				requestRol := httplib.Get(urlRol + user.Rol)
				requestRol.Header("Accept", "application/json")
				if err := requestRol.ToJSON(&Rid); err == nil {
					if len(Rid.Roles.Rol) > 0 {
						beego.Info("Rol: ", (Rid.Roles.Rol))
						m._post_usuario_rol.um_role_id = Rid.Roles.Rol[0].Id
						idUsuario, err := strconv.ParseInt((Uid.Usuarios.Usuario[0].Id), 10, 32)
						if err == nil {
							m._post_usuario_rol.um_user_id = idUsuario
						}
						sendRolUser := httplib.Post(postRolSuario)
						sendRolUser.Header("Accept", "application/json")
						sendRolUser.Body(m)
						if err := sendRolUser.ToJSON(&res); err == nil {
							beego.Info("insertado ...", res)
						} else {
							beego.Info("No se puede actualizar el rol!")
							outputError = map[string]interface{}{"Function": "FuncionalidadMidController:UpdateRol", "Error": err}
							return nil, outputError

						}
					} else {
						beego.Info("El rol no existe !")
						outputError = map[string]interface{}{"Function": "FuncionalidadMidController:UpdateRol", "Error": err}
						return nil, outputError
					}
				} else {
					beego.Info("R: ", requestUsuario)
					outputError = map[string]interface{}{"Function": "FuncionalidadMidController:UpdateRol", "Error": err}
					return nil, outputError
				}
			} else {
				beego.Info("El usuario no existe !")
				outputError = map[string]interface{}{"Function": "FuncionalidadMidController:UpdateRol", "Error": err}
				return nil, outputError
			}
			return roles, nil
		} else {
			beego.Info("R: ", requestUsuario)
			outputError = map[string]interface{}{"Function": "FuncionalidadMidController:UpdateRol", "Error": err}
			return nil, outputError
		}
	} else {
		return nil, outputError
	}
	// if err := request.GetJson(beego.AppConfig.String("GetRoleByUser")+m.User, RolesUsuario); err == nil {
	// }
}
