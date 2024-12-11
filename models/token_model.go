package models

// ? Inputs structures
// Token structure
type Token struct {
	Email string `json:"email"`
}

// Documento structure
type Documento struct {
	Numero string `json:"numero"`
}

// UserName structure
type UserName struct {
	User string `json:"user"`
}

// Payload structure
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

// AtributosToken structure
type AtributosToken struct {
	Usuario *struct {
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
