package models

type Token struct {
	Email string `json:"email"`
}

type Documento struct {
	Numero string `json:"numero"`
}

type UserName struct {
	User string `json:"user" valid:"Required"`
}

type InfoUsuario struct {
	Role               []string `json:"role"`
	Documento          string   `json:"documento"`
	DocumentoCompuesto string   `json:"documento_compuesto"`
	Email              string   `json:"email"`
	FamilyName         string   `json:"FamilyName"`
	Codigo             string   `json:"Codigo"`
	Estado             string   `json:"Estado"`
}

type EstudianteInfo struct {
	EstudianteCollection struct {
		Estudiante []Estudiante `json:"estudiante"`
	} `json:"estudianteCollection"`
}

type Estudiante struct {
	Codigo string `json:"codigo"`
	Estado string `json:"estado"`
}

type Usuario struct {
	Roles              []string
	FamilyName         string
	Email              string
	Documento          string
	DocumentoCompuesto string
}

// AtributosToken structure
type AtributosToken struct {
	Usuario *UsuarioData
}

type UsuarioData struct {
	Atributos []Atributo
}

type Atributo struct {
	Atributo string `json:"atributo"`
	Valor    string `json:"valor"`
}

type UserInfo struct {
	Codigo string   `json:"Codigo"`
	Estado string   `json:"Estado"`
	Email  string   `json:"email"`
	Rol    []string `json:"rol"`
}

type ClientCredentialsResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type ClientAuthRequestBody struct {
	Documento string `json:"documento"`
	ClienteId string `json:"clienteId"`
}
