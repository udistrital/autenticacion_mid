package models

// ? Inputs structures

// Body UpdateRol structure
type UpdateRol struct {
	User string `json:"user"`
	Rol  string `json:"rol"`
}

// ? Output structures

// Response ResRol structure
type ResRol struct {
	InfoUser InfoUserRoles `json:"InfoUser"`
}

// InfoUserRoles structure
type InfoUserRoles struct {
	Role []string `json:"Role"`
}
