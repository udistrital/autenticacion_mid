package models

type Periodo struct {
	Id                int    `json:"Id"`
	Activo            bool   `json:"Activo"`
	FechaCreacion     string `json:"FechaCreacion"`
	FechaModificacion string `json:"FechaModificacion"`
	FechaInicio       string `json:"FechaInicio"`
	FechaFin          string `json:"FechaFin"`
	Finalizado        bool   `json:"Finalizado"`
	Usuario           struct {
		Id                int    `json:"Id"`
		Activo            bool   `json:"Activo"`
		FechaCreacion     string `json:"FechaCreacion"`
		FechaModificacion string `json:"FechaModificacion"`
		Documento         string `json:"Documento"`
	} `json:"UsuarioId"`
	Rol struct {
		Id                 int    `json:"Id"`
		Activo             bool   `json:"Activo"`
		FechaCreacion      string `json:"FechaCreacion"`
		FechaModificacion  string `json:"FechaModificacion"`
		Nombre             string `json:"Nombre"`
		SistemaInformacion struct {
			Id                int    `json:"Id"`
			Activo            bool   `json:"Activo"`
			FechaCreacion     string `json:"FechaCreacion"`
			FechaModificacion string `json:"FechaModificacion"`
			Nombre            string `json:"Nombre"`
			Descripcion       string `json:"Descripcion"`
		} `json:"SistemaInformacionId"`
	} `json:"RolId"`
}

type Response struct {
	Data    []Periodo `json:"Data"`
	Message string    `json:"Message"`
	Status  int       `json:"Status"`
	Success bool      `json:"Success"`
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
}
