package models

type Periodo struct {
	Id                int     `json:"Id"`
	Activo            bool    `json:"Activo"`
	FechaCreacion     string  `json:"FechaCreacion"`
	FechaModificacion string  `json:"FechaModificacion"`
	FechaInicio       string  `json:"FechaInicio"`
	FechaFin          string  `json:"FechaFin"`
	Finalizado        bool    `json:"Finalizado"`
	Usuario           Usuario `json:"UsuarioId"`
	Rol               Rol     `json:"RolId"`
}

type Usuario struct {
	Id                int    `json:"Id"`
	Activo            bool   `json:"Activo"`
	FechaCreacion     string `json:"FechaCreacion"`
	FechaModificacion string `json:"FechaModificacion"`
	Documento         string `json:"Documento"`
}

type Rol struct {
	Id                 int    `json:"Id"`
	Activo             bool   `json:"Activo"`
	FechaCreacion      string `json:"FechaCreacion"`
	FechaModificacion  string `json:"FechaModificacion"`
	Nombre             string `json:"Nombre"`
	SistemaInformacion any    `json:"SistemaInformacionId"`
}

type Response struct {
	Data    []Periodo `json:"Data"`
	Message string    `json:"Message"`
	Status  int       `json:"Status"`
	Success bool      `json:"Success"`
}

type MetadataResponse struct {
	Data     []Periodo `json:"Data"`
	Message  string    `json:"Message"`
	Status   int       `json:"Status"`
	Success  bool      `json:"Success"`
	Metadata any       `json:"Metadata"`
}

type TerceroInfo struct {
	Tercero struct {
		Id             int    `json:"Id"`
		NombreCompleto string `json:"NombreCompleto"`
	} `json:"Tercero"`
	Identificacion struct {
		Numero             string `json:"Numero"`
		DigitoVerificacion int    `json:"DigitoVerificacion"`
		TipoDocumentoId    struct {
			CodigoAbreviacion string `json:"CodigoAbreviacion"`
		} `json:"TipoDocumentoId"`
	} `json:"Identificacion"`
}

type PeriodoRolUsuario struct {
	Nombre       string `json:"nombre"`
	Documento    string `json:"documento"`
	Correo       string `json:"correo"`
	RolUsuario   string `json:"rol_usuario"`
	Estado       bool   `json:"estado"`
	FechaInicial string `json:"fecha_inicial"`
	FechaFinal   string `json:"fecha_final"`
	Finalizado   bool   `json:"finalizado"`
	IdPeriodo    int    `json:"id_periodo"`
	IdTercero    int    `json:"id_tercero"`
}

// InfoPersona agrupa lo consultado en Terceros y en IS 7.3 para un documento.
type InfoPersona struct {
	Tercero    TerceroInfo
	ErrTercero error
	Correo     string
	ErrWso2    error
}

// ParamsPeriodos son los parámetros de consulta, todos opcionales, de
// /rol/periods y /rol/user/{documento}/periods. Query tiene el formato
// campo:valor,campo:valor; RolId solo aplica a /rol/periods.
type ParamsPeriodos struct {
	Query  string `form:"query" valid:"Match(/^([^:,]+:[^,]*(,[^:,]+:[^,]*)*)?$/)"`
	RolId  string `form:"RolId"`
	Limit  int64  `form:"limit"`
	Offset int64  `form:"offset"`
}
