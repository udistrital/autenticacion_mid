package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

// ? Inputs structures
// Token structure
type Token struct {
	Email string `json:"email"`
}

// UserName structure
type UserName struct {
	User string `json:"user"`
}

// Body UpdateRol structure
type UpdateRol struct {
	User   string `json:"user"`
	Rol  string `json:"rol"`
}

type UpdateRol2 struct {
	User   string `json:"user"`
	Roles  []string `json:"rol"`
}

//Body UpdatePerfil structure
type UpdatePerfil struct {
	UmAttrValue string `json:"um_attr_value"`
	UmId        int    `json:"um_id"`
}

// PostUsuarioRol structure
type PostUsuarioRol struct {
	um_role_id int `json:"um_role_id"`
	um_user_id int `json:"um_role_id"`
}

// ? Request response structures

type ResUpdatePerfil struct {
	Perfiles struct {
		Perfil []struct {
			UmId string `json:"um_id"`
		} `json:"perfil"`
	} `json:"perfiles"`
}

type ResUsuarioRoles struct {
	Usuario struct {
		Roles []struct {
			UmId     string `json:"um_id"`
			UmUserId string `json:"um_user_id"`
			UmRoleId string `json:"um_role_id"`
		} `json:"Roles"`
	} `json:"Usuario"`
}

type ResUserId struct {
	Usuarios struct {
		Usuario []struct {
			Id string `json:"Id"`
		} `json:"usuario"`
	} `json:"Usuarios"`
}

type ResRolId struct {
	Roles struct {
		Rol []struct {
			Id string `json:"id"`
		} `json:"Rol"`
	} `json:"Roles"`
}

type ResPerfilUsuario struct {
	Perfiles struct {
		Perfil []struct {
			UmId        string `json:"um_id"`
			UmAttrName  string `json:"um_attr_name"`
			UmAttrValue string `json:"um_attr_value"`
		} `json:"Perfil"`
	} `json:"Perfiles"`
}

// ? Outputs structures
// resUpdateRol structure
type resUpdateRol struct {
	Data    map[string]interface{} `json:"Data`
	Success bool                   `json:"Success"`
	Status  int                    `json:"Status"`
	Message string                 `json:"Message"`
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
			for _, v := range RolesUsuario.Usuario.Atributos {
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

				// fmt.Println(k, v)
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
			// fmt.Println("Payload: ", payload)
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
func AddRol(user UpdateRol) map[string]interface{} {
	var Uid ResUserId
	var Rid ResRolId
	var notExist = true
	var userName UserName
	userName.User = user.User
	urlUsuario := beego.AppConfig.String("Wso2Service") + "usuario"
	urlGetProfile := beego.AppConfig.String("Wso2Service") + "perfil"
	urlRol := beego.AppConfig.String("Wso2Service") + "rol"
	urlUpdateProfile := beego.AppConfig.String("Wso2UserService") + "updateperfil"
	urlAddProfile := beego.AppConfig.String("Wso2Service") + "perfil"
	postRolUsuario := beego.AppConfig.String("Wso2UserService") + "usuario_rol"
	var respuesta map[string]interface{}

	// ? Lanzar error para probar el manejo de errores del Controller
	// panic(map[string]interface{}{
	// 	"funcion": "TokenController.AddRol",
	// 	"err":     "El usuario no existe",
	// 	"status":  "400",
	// })

	if responseData, err := GetRolesByUser(userName); err == nil {
		for i := range responseData.Role {
			if responseData.Role[i] == user.Rol {
				beego.Info("YA EXISTE: ", responseData.Role[i])
				notExist = false
				returnData := map[string]interface{}{
					"InfoUser": map[string]interface{}{
						"Role": responseData.Role,
					},
				}
				respuesta = returnData
			}
		}
	}

	if notExist {
		var (
			m                 PostUsuarioRol
			res               map[string]map[string]interface{}
			updateProfileBody UpdatePerfil
			resUpdateProfile  map[string]map[string]interface{}
			addProfileBody    map[string]interface{}
			resAddProfile     map[string]map[string]interface{}
			resProfile        ResPerfilUsuario
			body              map[string]interface{}
		)
		beego.Info("URL: ", urlUsuario+"?user="+user.User)
		requestUsuario := httplib.Get(urlUsuario + "?user=" + user.User)
		requestUsuario.Header("Accept", "application/json")

		if err := requestUsuario.ToJSON(&Uid); err == nil {
			if len(Uid.Usuarios.Usuario) > 0 {
				beego.Info("User: ", (Uid.Usuarios.Usuario[0].Id))
				beego.Info("URL Rol: ", urlRol+"?rol="+user.Rol)
				requestRol := httplib.Get(urlRol + "?rol=" + user.Rol)
				requestRol.Header("Accept", "application/json")

				if err := requestRol.ToJSON(&Rid); err == nil {
					if len(Rid.Roles.Rol) > 0 {
						beego.Info("Rol: ", (Rid.Roles.Rol[0].Id))
						roleId, err := strconv.Atoi(Rid.Roles.Rol[0].Id)
						if err == nil {
							m.um_role_id = roleId
						}
						idUsuario, err := strconv.ParseInt((Uid.Usuarios.Usuario[0].Id), 10, 32)
						if err == nil {
							m.um_user_id = int(idUsuario)
						}
						sendRolUser := httplib.Post(postRolUsuario)
						beego.Info(sendRolUser)
						body = map[string]interface{}{
							"um_role_id": m.um_role_id,
							"um_user_id": m.um_user_id,
						}
						sendRolUser.Header("Accept", "application/json")
						sendRolUser.Header("Content-Type", "application/json")
						sendRolUser.JSONBody(body)
						if err := sendRolUser.ToJSON(&res); err == nil {
							if res["user_role"]["id"] != nil {
								beego.Info("Se añadió el rol!")
							} else {
								beego.Info("Ya existe asignación del rol!")
							}
							beego.Info("URL Get Profile: ", urlGetProfile+"?um_user_id="+strconv.Itoa(m.um_user_id))
							requestProfile := httplib.Get(urlGetProfile + "?um_user_id=" + strconv.Itoa(m.um_user_id))
							requestProfile.Header("Accept", "application/json")
							if err := requestProfile.ToJSON(&resProfile); err == nil {
								if len(resProfile.Perfiles.Perfil) == 0 {
									// post profile if not exist
									addProfile := httplib.Post(urlAddProfile)
									addProfile.Header("Accept", "application/json")
									addProfile.Header("Content-Type", "application/json")
									addProfileBody = map[string]interface{}{
										"um_attr_value": "Internal/everyone," + user.Rol, // adding rol to update profile
										"um_user_id":    m.um_user_id,
									}
									addProfile.JSONBody(addProfileBody)
									if err := addProfile.ToJSON(&resAddProfile); err == nil {
										if resUpdateProfile["perfiles"]["perfil"] != nil {
											beego.Info(resUpdateProfile["perfiles"]["perfil"])
											addProfileBody2 := map[string]interface{}{
												"post_addperfil": map[string]interface{}{
													"um_attr_value": "Internal/everyone," + user.Rol,
													"um_user_id":    m.um_user_id,
												},
											}
											respuesta = addProfileBody2
										} else {
											beego.Info("No se puede actualizar perfil !", resAddProfile)
											panic(map[string]interface{}{"Function": "FuncionalidadMidController:UpdateProfile", "Error": "No se puede actualizar perfil ! " + err.Error()})
										}
									} else {
										beego.Info("No se puede actualizar perfil !")
										panic(map[string]interface{}{"Function": "FuncionalidadMidController:UpdateProfile", "Error": "No se puede actualizar perfil ! " + err.Error()})
									}
								} else {
									if resProfile.Perfiles.Perfil[0].UmAttrValue != "" &&
										resProfile.Perfiles.Perfil[0].UmId != "" {
										// put profile if exist
										beego.Info("perfil actual", resProfile.Perfiles.Perfil[0].UmAttrValue)
										beego.Info("ID perfil actual", resProfile.Perfiles.Perfil[0].UmId)
										str := fmt.Sprintf("%v", resProfile.Perfiles.Perfil[0].UmAttrValue)
										rolesUsuario := strings.Split(str, ",")
										rolExiste := false

										for _, role := range rolesUsuario {
											fmt.Println("role: ", role)
											if role == user.Rol {
												rolExiste = true
												break
											}
										}
										fmt.Println("str: ", str)
										if !rolExiste { // if not exist profile
											updateProfile := httplib.Put(urlUpdateProfile)
											updateProfile.Header("Accept", "application/json")
											updateProfile.Header("Content-Type", "application/json")
											UmIdProfile, errIdProfile := strconv.Atoi(fmt.Sprintf("%v", resProfile.Perfiles.Perfil[0].UmId))
											if errIdProfile == nil {
												updateProfileBody.UmAttrValue = user.Rol + "," + str
												updateProfileBody.UmId = UmIdProfile
											}
											updateProfile.JSONBody(updateProfileBody)
											fmt.Println("updateProfileBody: ", updateProfileBody)
											if err := updateProfile.ToJSON(&resUpdateProfile); err == nil {
												fmt.Println("resUpdateProfile: ", resUpdateProfile)
												if resUpdateProfile["perfiles"]["perfil"] != nil {
													beego.Info(resUpdateProfile["perfiles"]["perfil"])
													updateProfileBody2 := map[string]interface{}{
														"put_updateperfil": map[string]interface{}{
															"um_attr_value": user.Rol + "," + str,
															"um_id":         UmIdProfile,
														},
													}
													respuesta = updateProfileBody2
												} else {
													beego.Info("No se puede actualizar perfil !")
													panic(map[string]interface{}{"Function": "FuncionalidadMidController:UpdateProfile", "Error": err})
												}
											} else {
												beego.Info("No se puede actualizar perfil, error en la peticion!")
												panic(map[string]interface{}{"Function": "FuncionalidadMidController:UpdateProfile", "Error": "No se puede actualizar perfil, error en la peticion! " + err.Error()})
											}
										} else {
											panic(map[string]interface{}{"Function": "FuncionalidadMidController:UpdateProfile", "Error": "El rol ya existe en el perfil!"})
										}
									}
								}
							} else {
								beego.Info("No se puede actualizar perfil !")
								panic(map[string]interface{}{"Function": "FuncionalidadMidController:UpdateProfile", "Error": err})
							}
						} else {
							beego.Info("No se puede actualizar el rol!")
							panic(map[string]interface{}{"Function": "FuncionalidadMidController:addRol No se puede actualizar rol !", "Error": err})
						}
					} else {
						beego.Info("El rol no existe !")
						panic(map[string]interface{}{"Function": "FuncionalidadMidController:addRol", "Error": "El rol no existe! " + err.Error()})
					}
				} else {
					beego.Info("R: ", requestUsuario)
					panic(map[string]interface{}{"Function": "FuncionalidadMidController:addRol", "Error": err})
				}
			} else {
				beego.Info("El usuario no existe !")
				panic(map[string]interface{}{"Function": "FuncionalidadMidController:addRol", "Error": "El usuario no existe! " + err.Error()})
			}
		} else {
			beego.Info("R: ", requestUsuario)
			panic(map[string]interface{}{"Function": "FuncionalidadMidController:addRol", "Error": err})
		}
	}
	return respuesta
}

func DeleteRol(user UpdateRol) map[string]interface{} {
	var updatePerfilBody UpdatePerfil
	var resUserId ResUserId
	var resRolId ResRolId
	var resUpdateProfile ResUpdatePerfil
	var resUsuarioRoles ResUsuarioRoles
	var perfilUsuario ResPerfilUsuario
	var userName UserName
	
	var rolEliminar string

	userName.User = user.User
	urlUsuario := beego.AppConfig.String("Wso2Service") + "usuario"
	urlGetProfile := beego.AppConfig.String("Wso2Service") + "perfil"
	urlRol := beego.AppConfig.String("Wso2Service") + "rol"
	urlDeleteUsuarioRol := beego.AppConfig.String("Wso2UserService") + "deleteusuariorol/"
	urlUpdateProfile := beego.AppConfig.String("Wso2UserService") + "updateperfil"

	if responseData, err := GetRolesByUser(userName); err == nil {
		for i := range responseData.Role {
			if responseData.Role[i] == user.Rol {
				rolEliminar = responseData.Role[i]
				beego.Info("URL Obtener UserId: ", urlUsuario+"?user="+user.User)
				requestUsuario := httplib.Get(urlUsuario + "?user=" + user.User)
				requestUsuario.Header("Accept", "application/json")
				// ? Lanzar error para probar el manejo de errores del Controller
				// panic(map[string]interface{}{
				// 	"funcion": "TokenController.DeleteRol",
				// 	"err":     "El usuario ya tiene asignado este rol",
				// 	"status":  "400",
				// })
				if err := requestUsuario.ToJSON(&resUserId); err == nil {
					if len(resUserId.Usuarios.Usuario) > 0 {
						beego.Info("User ENCONTRADO: ", (resUserId.Usuarios.Usuario[0].Id))
						requestRol := httplib.Get(urlRol + "?rol=" + user.Rol)
						beego.Info("URL Obtener RolId: ", urlRol+"?rol="+user.Rol)
						requestRol.Header("Accept", "application/json")
						if err := requestRol.ToJSON(&resRolId); err == nil {
							if len(resRolId.Roles.Rol) > 0 {
								beego.Info("Rol ENCONTRADO: ", (resRolId.Roles.Rol[0].Id))
								urlGetUsuarioRol := beego.AppConfig.String("Wso2Service") + "get_usuario_rol?usuario=" + resUserId.Usuarios.Usuario[0].Id
								beego.Info("URL Obtener ResUsuarioRoles (Rompimiento): ", urlGetUsuarioRol)
								requestUsuarioRol := httplib.Get(urlGetUsuarioRol)
								requestUsuarioRol.Header("Accept", "application/json")
								if err := requestUsuarioRol.ToJSON(&resUsuarioRoles); err == nil {
									if len(resUsuarioRoles.Usuario.Roles) > 0 {
										bandera := false
										for i := range resUsuarioRoles.Usuario.Roles {
											if resUsuarioRoles.Usuario.Roles[i].UmRoleId == resRolId.Roles.Rol[0].Id {
												// Role found
												bandera = true
												requestDeleteUsuarioRol := httplib.Delete(urlDeleteUsuarioRol + resUsuarioRoles.Usuario.Roles[i].UmId)
												if err := requestDeleteUsuarioRol.ToJSON(&resUsuarioRoles); err == nil {
													// Get perfil usuario
													requestPerfilUsuario := httplib.Get(urlGetProfile + "?um_user_id=" + resUserId.Usuarios.Usuario[0].Id)
													requestPerfilUsuario.Header("Accept", "application/json")
													if err := requestPerfilUsuario.ToJSON(&perfilUsuario); err == nil {
														if len(perfilUsuario.Perfiles.Perfil) > 0 {
															// perfil exists
															perfilActualizado := strings.Split(perfilUsuario.Perfiles.Perfil[0].UmAttrValue, ",")
															for i := range perfilActualizado {
																if perfilActualizado[i] == rolEliminar {
																	perfilActualizado = append(perfilActualizado[:i], perfilActualizado[i+1:]...)
																	break
																}
															}
															strPerfil := strings.Join(perfilActualizado, ",") // Convirtiendo array a string
															updateProfile := httplib.Put(urlUpdateProfile)
															updateProfile.Header("Accept", "application/json")
															updateProfile.Header("Content-Type", "application/json")
															updatePerfilBody.UmAttrValue = strPerfil
															updatePerfilBody.UmId, _ = strconv.Atoi(perfilUsuario.Perfiles.Perfil[0].UmId)
															updateProfile.JSONBody(updatePerfilBody)
															if err := updateProfile.ToJSON(&resUpdateProfile); err == nil {
																if resUpdateProfile.Perfiles.Perfil != nil {
																	beego.Info("PERFIL ACTUALIZADO: ", resUpdateProfile.Perfiles.Perfil)
																	response := map[string]interface{}{
																		"put_updateperfil": map[string]interface{}{
																			"um_attr_value": strPerfil,
																			"um_id": updatePerfilBody.UmId,
																		},
																	}
																	return response
																} else {
																	beego.Info("No se puede actualizar perfil!")
																	panic(map[string]interface{}{"funcion": "TokenController.DeleteRol", "err": "No se puede actualizar perfil", "status": "400"})
																}
															} else {
																panic(map[string]interface{}{"funcion": "TokenController.DeleteRol", "err": "Error al actualizar el perfil del usuario", "status": "400"})
															}
														}
													} else {
														panic(map[string]interface{}{"funcion": "TokenController.DeleteRol", "err": "Error al obtener el perfil del usuario", "status": "400"})
													}
												} else {
													panic(map[string]interface{}{"funcion": "TokenController.DeleteRol", "err": "Error al eliminar el rol del usuario", "status": "400"})
												}
											}
										}
										if !bandera { // Caso en que el rol no existe en la tabla de rompimiento pero igual se debe validar que no exista en el perfil
											// Get perfil usuario
											requestPerfilUsuario := httplib.Get(urlGetProfile + "?um_user_id=" + resUserId.Usuarios.Usuario[0].Id)
											requestPerfilUsuario.Header("Accept", "application/json")
											if err := requestPerfilUsuario.ToJSON(&perfilUsuario); err == nil {
												if len(perfilUsuario.Perfiles.Perfil) > 0 {
													// perfil exists
													perfilActualizado := strings.Split(perfilUsuario.Perfiles.Perfil[0].UmAttrValue, ",")
													for i := range perfilActualizado {
														if perfilActualizado[i] == rolEliminar {
															perfilActualizado = append(perfilActualizado[:i], perfilActualizado[i+1:]...)
															break
														}
													}
													strPerfil := strings.Join(perfilActualizado, ",") // Convirtiendo array a string
													updateProfile := httplib.Put(urlUpdateProfile)
													updateProfile.Header("Accept", "application/json")
													updateProfile.Header("Content-Type", "application/json")
													updatePerfilBody.UmAttrValue = strPerfil
													updatePerfilBody.UmId, _ = strconv.Atoi(perfilUsuario.Perfiles.Perfil[0].UmId)
													updateProfile.JSONBody(updatePerfilBody)
													if err := updateProfile.ToJSON(&resUpdateProfile); err == nil {
														if resUpdateProfile.Perfiles.Perfil != nil {
															beego.Info("PERFIL ACTUALIZADO: ", resUpdateProfile.Perfiles.Perfil)
															response := map[string]interface{}{
																"put_updateperfil": map[string]interface{}{
																	"um_attr_value": strPerfil,
																	"um_id": updatePerfilBody.UmId,
																},
															}
															return response
														} else {
															beego.Info("No se puede actualizar perfil!")
															panic(map[string]interface{}{"funcion": "TokenController.DeleteRol", "err": "No se puede actualizar perfil", "status": "400"})
														}
													} else {
														panic(map[string]interface{}{"funcion": "TokenController.DeleteRol", "err": "Error al actualizar el perfil del usuario", "status": "400"})
													}
												}
											} else {
												panic(map[string]interface{}{"funcion": "TokenController.DeleteRol", "err": "Error al obtener el perfil del usuario", "status": "400"})
											}
										}
									}
								} else {
									panic(map[string]interface{}{"funcion": "TokenController.DeleteRol", "err": "Error al obtener los roles del usuario", "status": "400"})
								}
							}
						} else {
							panic(map[string]interface{}{"funcion": "TokenController.DeleteRol", "err": "Error al obtener el Id del rol", "status": "400"})
						}
					}
				} else {
					panic(map[string]interface{}{"funcion": "TokenController.DeleteRol", "err": "Error al obtener el Id del usuario", "status": "400"})
				}
			}
		}
		if rolEliminar == "" {
			panic(map[string]interface{}{"funcion": "TokenController.DeleteRol", "err": "El usuario no tiene asignado el rol", "status": "400"})
		}
	}
	return nil
}