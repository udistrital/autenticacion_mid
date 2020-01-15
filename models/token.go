package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

// Token structure
type Token struct {
	Email string `json:"email"`
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

type returnInfo struct {
	InfoUserInfo     UserInfo ``
	InfoRolesUsuario RolesUsuario
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
func GetRolesByUser(user string) (roles *Payload, outputError map[string]interface{}) {
	var RolesUsuario AtributosToken
	var estudianteInfo EstudianteInfo
	var familyName string
	var documento string
	var mail string
	var documentoCompuesto string
	userRoles := []string{}

	fmt.Println(beego.AppConfig.String("GetRoleByUser") + user)

	r := httplib.Get(beego.AppConfig.String("GetRoleByUser") + user)
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
	// if err := request.GetJson(beego.AppConfig.String("GetRoleByUser")+m.User, RolesUsuario); err == nil {

	// }

}
