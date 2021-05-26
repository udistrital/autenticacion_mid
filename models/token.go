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
		um_role_id int `json:"um_role_id"`
		um_user_id int `json:"um_role_id"`
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

	beego.Info("URL: ", beego.AppConfig.String("Wso2Service")+"roles?usuario="+user.User)
	r := httplib.Get(beego.AppConfig.String("Wso2Service") + "roles?usuario=" + user.User)
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
func AddRol(user UpdateRol) (roleSuccess *map[string]map[string]interface{}, outputError map[string]interface{}) {
	var Uid UserId
	var Rid RolId
	var notExist = true
	var userName UserName
	userName.User = user.User
	urlUsuario := beego.AppConfig.String("Wso2Service") + "usuario/"
	urlGetProfile := beego.AppConfig.String("Wso2Service") + "perfil/"
	urlRol := beego.AppConfig.String("Wso2Service") + "rol/"
	urlUpdateProfile := beego.AppConfig.String("Wso2Service") + "updateperfil"
	urlAddProfile := beego.AppConfig.String("Wso2Service") + "perfil"
	postRolSuario := beego.AppConfig.String("Wso2Service") + "usuario_rol"

	if responseData, err := GetRolesByUser(userName); err == nil {
		for i := range responseData.Role {
			if responseData.Role[i] == user.Rol {
				beego.Info("YA EXISTE: ", responseData.Role[i])
				notExist = false
				returnData := map[string]map[string]interface{}{
					"InfoUser": {
						"Role":      responseData.Role,
						"Codigo":    responseData.Codigo,
						"Estado":    responseData.Estado,
						"Documento": responseData.Documento,
					},
				}
				return &returnData, nil
			}
		}
	}
	if notExist {
		var (
			m                 PostUsuarioRol
			res               map[string]map[string]interface{}
			updateProfileBody map[string]map[string]interface{}
			resUpdateProfile  map[string]map[string]interface{}
			addProfileBody    map[string]map[string]interface{}
			resAddProfile     map[string]map[string]interface{}
			resProfile        map[string]map[string][]map[string]interface{}
			body              map[string]map[string]interface{}
		)
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
							m._post_usuario_rol.um_user_id = int(idUsuario)
						}
						sendRolUser := httplib.Post(postRolSuario)
						beego.Info(sendRolUser)
						body = map[string]map[string]interface{}{
							"_post_usuario_rol": {
								"um_user_id": m._post_usuario_rol.um_user_id,
								"um_role_id": m._post_usuario_rol.um_role_id,
							},
						}
						sendRolUser.Header("Accept", "application/json")
						sendRolUser.Header("Content-Type", "application/json")
						sendRolUser.JSONBody(body)
						if err := sendRolUser.ToJSON(&res); err == nil {
							if res["Fault"]["faultstring"] != nil {
								beego.Info("Ya existe asignación del rol!")
							} else {
								beego.Info("Se añadió!")
							}
							requestProfile := httplib.Get(urlGetProfile + strconv.Itoa(m._post_usuario_rol.um_user_id))
							requestProfile.Header("Accept", "application/json")
							if err := requestProfile.ToJSON(&resProfile); err == nil {
								if len(resProfile["Perfiles"]["Perfil"]) == 0 {
									// post profile if not exist
									addProfile := httplib.Post(urlAddProfile)
									addProfile.Header("Accept", "application/json")
									addProfile.Header("Content-Type", "application/json")
									addProfileBody = map[string]map[string]interface{}{
										"_post_perfil": {
											"um_attr_value": "Internal/everyone," + user.Rol, // adding rol to update profile
											"um_user_id":    m._post_usuario_rol.um_user_id,
										},
									}
									addProfile.JSONBody(addProfileBody)
									if err := addProfile.ToJSON(&resAddProfile); err == nil {
										if resUpdateProfile["perfiles"]["perfil"] != nil {
											beego.Info(resUpdateProfile["perfiles"]["perfil"])
											return &updateProfileBody, nil
										} else {
											beego.Info("No se puede actualizar perfil !", resAddProfile)
											outputError = map[string]interface{}{"Function": "FuncionalidadMidController:UpdateProfile", "Error": err}
											return nil, outputError
										}
									} else {
										beego.Info("No se puede actualizar perfil !")
										outputError = map[string]interface{}{"Function": "FuncionalidadMidController:UpdateProfile", "Error": err}
										return nil, outputError
									}
								} else {
									if resProfile["Perfiles"]["Perfil"][0]["um_attr_value"] != nil &&
										resProfile["Perfiles"]["Perfil"][0]["um_id"] != nil {
										// put profile if exist
										beego.Info("perfil actual", resProfile["Perfiles"]["Perfil"][0]["um_attr_value"])
										beego.Info("ID perfil actual", resProfile["Perfiles"]["Perfil"][0]["um_id"])
										str := fmt.Sprintf("%v", resProfile["Perfiles"]["Perfil"][0]["um_attr_value"])
										if strings.Index(str, user.Rol) == -1 { // if not exist profile
											updateProfile := httplib.Put(urlUpdateProfile)
											updateProfile.Header("Accept", "application/json")
											updateProfile.Header("Content-Type", "application/json")
											UmIdProfile, errIdProfile := strconv.Atoi(fmt.Sprintf("%v", resProfile["Perfiles"]["Perfil"][0]["um_id"]))
											if errIdProfile == nil {
												updateProfileBody = map[string]map[string]interface{}{
													"_put_updateperfil": {
														"um_attr_value": user.Rol + "," + str, // adding rol to update profile
														"um_id":         UmIdProfile,
													},
												}
											}
											updateProfile.JSONBody(updateProfileBody)
											if err := updateProfile.ToJSON(&resUpdateProfile); err == nil {
												if resUpdateProfile["perfiles"]["perfil"] != nil {
													beego.Info(resUpdateProfile["perfiles"]["perfil"])
													return &updateProfileBody, nil
												} else {
													beego.Info("No se puede actualizar perfil !")
													outputError = map[string]interface{}{"Function": "FuncionalidadMidController:UpdateProfile", "Error": err}
													return nil, outputError
												}
											} else {
												beego.Info("No se puede actualizar perfil !")
												outputError = map[string]interface{}{"Function": "FuncionalidadMidController:UpdateProfile", "Error": err}
												return nil, outputError
											}
										} else {
											returnData := map[string]map[string]interface{}{
												"InfoUser": {
													"Role": strings.Split(str, ","),
												},
											}
											return &returnData, nil
										}
									}
								}
							} else {
								beego.Info("No se puede actualizar perfil !")
								outputError = map[string]interface{}{"Function": "FuncionalidadMidController:addRol No se puede actualizar perfil !", "Error": err}
								return nil, outputError
							}
						} else {
							beego.Info("No se puede actualizar el rol!")
							outputError = map[string]interface{}{"Function": "FuncionalidadMidController:addRol No se puede actualizar rol !", "Error": err}
							return nil, outputError
						}
					} else {
						beego.Info("El rol no existe !")
						outputError = map[string]interface{}{"Function": "FuncionalidadMidController:addRol el rol no existe!", "Error": err}
						return nil, outputError
					}
				} else {
					beego.Info("R: ", requestUsuario)
					outputError = map[string]interface{}{"Function": "FuncionalidadMidController: R", "Error": err}
					return nil, outputError
				}
			} else {
				beego.Info("El usuario no existe !")
				outputError = map[string]interface{}{"Function": "FuncionalidadMidController:addRol", "Error": err}
				return nil, outputError
			}
			return roleSuccess, nil
		} else {
			beego.Info("R: ", requestUsuario)
			outputError = map[string]interface{}{"Function": "FuncionalidadMidController:addRol", "Error": err}
			return nil, outputError
		}
	} else {
		return nil, outputError
	}
	// if err := request.GetJson(beego.AppConfig.String("GetRoleByUser")+m.User, RolesUsuario); err == nil {
	// }
}
