package models

// UpdateRol es el cuerpo de /rol/add y /rol/remove.
type UpdateRol struct {
	User string `json:"user"`
	Rol  string `json:"rol"`
}

// ResRol es la respuesta de /rol/add y /rol/remove.
type ResRol struct {
	InfoUser InfoUserRoles `json:"InfoUser"`
}

type InfoUserRoles struct {
	Role []string `json:"Role"`
}
