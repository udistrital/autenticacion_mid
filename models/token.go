package models

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

// Token structure
type Token struct {
	Email string `json:"email"`
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

// UserInfo structure
type UserInfo struct {
	Codigo string   `json:"Codigo"`
	Estado string   `json:"Estado"`
	Email  string   `json:"Email"`
	Rol    []string `json:"Rol"`
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
			return nil, errors.New("Email no registrado")
		}
	}
	fmt.Println(err)
	return nil, err
}
