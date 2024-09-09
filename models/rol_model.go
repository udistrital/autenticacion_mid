package models

// ? Inputs structures

// Body UpdateRol structure
type UpdateRol struct {
	User string `json:"user"`
	Rol  string `json:"rol"`
}

//Body UpdatePerfil structure
type UpdatePerfil struct {
	UmId        int    `json:"um_id"`
	UmAttrValue string `json:"um_attr_value"`
}

//Body AddPerfil structure
type AddPerfil struct {
	UmUserId    int    `json:"um_user_id"`
	UmAttrValue string `json:"um_attr_value"`
}

// Body UpdateUsuarioRol structure tabla_rompimiento
type UpdateUsuarioRol struct {
	UmRoleId int `json:"um_role_id"`
	UmUserId int `json:"um_user_id"`
}

// ? Request response structures

// Response ResUpdatePerfil structure
type ResUpdatePerfil struct {
	Perfiles struct {
		Perfil []struct {
			UmId string `json:"um_id"`
		} `json:"perfil"`
	} `json:"perfiles"`
}

// Response ResUsuarioRoles structure
type ResGetUsuarioRoles struct {
	Usuario struct {
		Roles []struct {
			UmId     string `json:"um_id"`
			UmUserId string `json:"um_user_id"`
			UmRoleId string `json:"um_role_id"`
		} `json:"Roles"`
	} `json:"Usuario"`
}

// Response ResUpdateUsuarioRol structure
type ResUpdateUsuarioRol struct {
	Id string `json:"id"`
}

// Response ResUserId structure
type ResUserId struct {
	Usuarios struct {
		Usuario []struct {
			Id string `json:"Id"`
		} `json:"usuario"`
	} `json:"Usuarios"`
}

// Response ResRolId structure
type ResRolId struct {
	Roles struct {
		Rol []struct {
			Id string `json:"id"`
		} `json:"Rol"`
	} `json:"Roles"`
}

// Response ResPerfilUsuario structure
type ResPerfilUsuario struct {
	Perfiles struct {
		Perfil []struct {
			UmId        string `json:"um_id"`
			UmAttrName  string `json:"um_attr_name"`
			UmAttrValue string `json:"um_attr_value"`
		} `json:"Perfil"`
	} `json:"Perfiles"`
}

// Response DeleteUsuarioRol structure
type ResDeleteUsuarioRol struct {
	UserRole *struct {
		Id string `json:"id"`
	} `json:"user_role"`
}
