package scim

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/autenticacion_mid/models"
	"github.com/udistrital/utils_oas/ssm"
)

// Llamadas a WSO2 IS 7.3: operaciones sobre usuarios y grupos con su API SCIM2,
// y tokens client_credentials.

const (
	customSchema     = "urn:scim:schemas:extension:custom:User"
	scimScopes       = "internal_user_mgt_list internal_user_mgt_view internal_group_mgt_view internal_group_mgt_update"
	auditoriaUserKey = "user"
)

var wso2ISBasePath = beego.AppConfig.String("Wso2IsService")
var wso2IsSecret = beego.AppConfig.String("Wso2IsSecret")
var wso2IsClient = beego.AppConfig.String("Wso2IsClient")

type scimUser struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
	Name     struct {
		FamilyName string `json:"familyName"`
	} `json:"name"`
	Emails []struct {
		Value   string `json:"value"`
		Primary bool   `json:"primary"`
	} `json:"emails"`
	Groups []struct {
		Display string `json:"display"`
	} `json:"groups"`
	Roles []struct {
		Display      string `json:"display"`
		AudienceType string `json:"audienceType"`
	} `json:"roles"`
	Custom struct {
		Documento          string `json:"documento"`
		DocumentoCompuesto string `json:"documento_compuesto"`
	} `json:"urn:scim:schemas:extension:custom:User"`
}

type scimList struct {
	TotalResults int        `json:"totalResults"`
	Resources    []scimUser `json:"Resources"`
}

type scimGroup struct {
	ID    string `json:"id"`
	Roles []struct {
		Display string `json:"display"`
	} `json:"roles"`
}

type scimGroupList struct {
	Resources []scimGroup `json:"Resources"`
}

type scimPatch struct {
	Schemas    []string        `json:"schemas"`
	Operations []scimOperacion `json:"Operations"`
}

type scimOperacion struct {
	Op    string        `json:"op"`
	Path  string        `json:"path,omitempty"`
	Value *scimMiembros `json:"value,omitempty"`
}

type scimMiembros struct {
	Members []scimMiembro `json:"members"`
}

type scimMiembro struct {
	Value   string `json:"value"`
	Display string `json:"display,omitempty"`
}

type clientCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// Usuario es un usuario de IS: sus datos, sus roles y sus grupos.
type Usuario struct {
	ID                 string
	UserName           string
	Roles              []string
	Grupos             []string
	FamilyName         string
	Email              string
	Documento          string
	DocumentoCompuesto string
}

var (
	ErrUsuarioNoExiste = errors.New("el usuario no existe")
	ErrGrupoNoExiste   = errors.New("el grupo no existe")
)

var (
	tokenMu     sync.Mutex
	credentials *clientCredentials
	token       string
	tokenExpiry time.Time
	isClient    = &http.Client{Timeout: 30 * time.Second}
)

// GetGruposUsuario devuelve los datos y grupos del usuario en IS 7.3 por
// correo o nombre de usuario, o nil si no coincide una sola cuenta. Si
// varias cuentas comparten el correo, se usa la del usuario del token.
func GetGruposUsuario(ctx context.Context, usuario string) (*Usuario, error) {
	atributo := "userName"
	if strings.Contains(usuario, "@") {
		atributo = "emails"
	}
	usuarios, err := buscarUsuarios(ctx, atributo, usuario)
	if err != nil {
		return nil, fmt.Errorf("error al obtener los grupos del usuario %s: %w", usuario, err)
	}

	if len(usuarios) == 0 {
		logs.Warning("usuario no registrado %s", usuario)
		return nil, nil
	}

	cuenta := &usuarios[0]
	if len(usuarios) > 1 {
		sub, _ := ctx.Value(auditoriaUserKey).(string)
		if sub == "" {
			logs.Warning("%d cuentas tienen el email %s y el token no tiene usuario: usuario no registrado", len(usuarios), usuario)
			return nil, nil
		}
		if cuenta = cuentaDelToken(sub, usuarios); cuenta == nil {
			logs.Warning("%d cuentas tienen el email %s y ninguna es la del token (%s): usuario no registrado", len(usuarios), usuario, sub)
			return nil, nil
		}
		logs.Info("%d cuentas tienen el email %s: se usa %s (%s), la del token", len(usuarios), usuario, cuenta.UserName, cuenta.ID)
	}

	return toUsuario(*cuenta), nil
}

// GetInfoByDocumento devuelve los datos y grupos del usuario en IS 7.3 por
// documento, o nil si no coincide una sola cuenta.
func GetInfoByDocumento(ctx context.Context, documento string) (*Usuario, error) {
	usuarios, err := buscarUsuarios(ctx, customSchema+":documento", documento)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la información del documento %s: %w", documento, err)
	}

	if len(usuarios) == 0 {
		return nil, nil
	}

	if len(usuarios) > 1 {
		logs.Warning("%d cuentas tienen el documento %s: usuario no registrado", len(usuarios), documento)
		return nil, nil
	}

	return toUsuario(usuarios[0]), nil
}

// AgregarUsuarioGrupo agrega el usuario al grupo de IS 7.3. Falla si no hay
// un solo grupo con ese nombre.
func AgregarUsuarioGrupo(ctx context.Context, usuario *Usuario, grupo string) error {
	operacion := scimOperacion{
		Op:    "add",
		Value: &scimMiembros{Members: []scimMiembro{{Value: usuario.ID, Display: usuario.UserName}}},
	}
	if err := modificarGrupo(ctx, grupo, operacion); err != nil {
		return fmt.Errorf("error al agregar el usuario al grupo %s: %w", grupo, err)
	}
	return nil
}

// QuitarUsuarioGrupo quita el usuario con ID usuarioID del grupo de IS 7.3.
// Falla si no hay un solo grupo con ese nombre.
func QuitarUsuarioGrupo(ctx context.Context, usuarioID, grupo string) error {
	operacion := scimOperacion{
		Op:   "remove",
		Path: "members[value eq " + scimString(usuarioID) + "]",
	}
	if err := modificarGrupo(ctx, grupo, operacion); err != nil {
		return fmt.Errorf("error al quitar el usuario del grupo %s: %w", grupo, err)
	}
	return nil
}

func ClientCredentialsRequest(ctx context.Context, clientID, clientSecret, scope string) (models.ClientCredentialsResponse, error) {
	form := url.Values{"grant_type": {"client_credentials"}}
	if scope != "" {
		form.Set("scope", scope)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, wso2ISBasePath+"/oauth2/token",
		strings.NewReader(form.Encode()))
	if err != nil {
		return models.ClientCredentialsResponse{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(clientID, clientSecret)

	resp, err := isClient.Do(req)
	if err != nil {
		return models.ClientCredentialsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.ClientCredentialsResponse{}, fmt.Errorf("IS respondió HTTP %d al pedir el token", resp.StatusCode)
	}

	var response models.ClientCredentialsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	return response, err
}

// cuentaDelToken devuelve la cuenta de usuarios a la que pertenece el token con
// el sub dado, o nil. El sub es el ID de la cuenta o su nombre de usuario, al que
// IS puede anteponer PRIMARY/ y agregar @carbon.super. Los inicios de sesión
// federados no coinciden con ninguna cuenta.
func cuentaDelToken(sub string, usuarios []scimUser) *scimUser {
	nombre := sub
	if i := strings.Index(nombre, "/"); i >= 0 && strings.EqualFold(nombre[:i], "PRIMARY") {
		nombre = nombre[i+1:]
	}

	if strings.HasSuffix(strings.ToLower(nombre), "@carbon.super") {
		nombre = nombre[:len(nombre)-len("@carbon.super")]
	}

	for i := range usuarios {
		u := &usuarios[i]
		if (u.ID != "" && strings.EqualFold(sub, u.ID)) || strings.EqualFold(nombre, u.UserName) {
			return u
		}
	}

	return nil
}

// GrupoConRoles indica si el grupo tiene roles asignados. Falla con
// ErrGrupoNoExiste si no hay un solo grupo con ese nombre.
func GrupoConRoles(ctx context.Context, grupo string) (bool, error) {
	g, err := buscarGrupo(ctx, grupo)
	if err != nil {
		return false, err
	}
	return len(g.Roles) > 0, nil
}

func modificarGrupo(ctx context.Context, grupo string, operacion scimOperacion) error {
	g, err := buscarGrupo(ctx, grupo)
	if err != nil {
		return err
	}
	body := scimPatch{
		Schemas:    []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
		Operations: []scimOperacion{operacion},
	}
	endpoint := wso2ISBasePath + "/scim2/Groups/" + url.PathEscape(g.ID)

	return callWSO2IS(ctx, http.MethodPatch, endpoint, body, nil)
}

// buscarGrupo busca el grupo por nombre en IS 7.3 y lo devuelve con sus roles.
// Falla si no hay uno solo.
func buscarGrupo(ctx context.Context, nombre string) (scimGroup, error) {
	query := url.Values{}
	query.Set("filter", "displayName eq "+scimString(nombre))
	query.Set("attributes", "id,displayName,roles")
	endpoint := wso2ISBasePath + "/scim2/Groups?" + query.Encode()

	var list scimGroupList
	if err := callWSO2IS(ctx, http.MethodGet, endpoint, nil, &list); err != nil {
		return scimGroup{}, fmt.Errorf("error al obtener el grupo %s: %w", nombre, err)
	}
	if len(list.Resources) != 1 {
		return scimGroup{}, fmt.Errorf("%w: %s", ErrGrupoNoExiste, nombre)
	}
	return list.Resources[0], nil
}

func buscarUsuarios(ctx context.Context, atributo, valor string) ([]scimUser, error) {
	query := url.Values{}
	query.Set("filter", fmt.Sprintf("%s eq %s", atributo, scimString(valor)))
	query.Set("attributes", "id,userName,name.familyName,emails,groups,roles,"+customSchema)
	endpoint := wso2ISBasePath + "/scim2/Users?" + query.Encode()

	var list scimList
	if err := callWSO2IS(ctx, http.MethodGet, endpoint, nil, &list); err != nil {
		return nil, err
	}
	return list.Resources, nil
}

func toUsuario(u scimUser) *Usuario {
	roles := []string{"Internal/everyone"}
	for _, r := range u.Roles {
		if r.AudienceType == "organization" && r.Display == "selfsignup" {
			roles = append(roles, "Internal/selfsignup")
			break
		}
	}
	var grupos []string
	for _, g := range u.Groups {
		grupos = append(grupos, g.Display)
	}
	mail := ""
	for _, e := range u.Emails {
		if mail == "" || e.Primary {
			mail = e.Value
		}
	}

	return &Usuario{
		ID:                 u.ID,
		UserName:           u.UserName,
		Roles:              roles,
		Grupos:             grupos,
		FamilyName:         u.Name.FamilyName,
		Email:              mail,
		Documento:          u.Custom.Documento,
		DocumentoCompuesto: u.Custom.DocumentoCompuesto,
	}
}

func scimString(v string) string {
	return `"` + strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(v) + `"`
}

// callWSO2IS envía payload como JSON a IS 7.3 y decodifica la respuesta en
// target; ambos son opcionales. Si IS responde HTTP 401, reintenta una vez con un
// token nuevo. Si la respuesta no es exitosa, el error incluye su estado y cuerpo.
func callWSO2IS(ctx context.Context, method, endpoint string, payload, target any) error {
	var data []byte
	if payload != nil {
		var err error
		if data, err = json.Marshal(payload); err != nil {
			return err
		}
	}

	resp, err := enviarWSO2IS(ctx, method, endpoint, data)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		// El token pudo ser revocado antes de expirar.
		resp.Body.Close()
		resetToken()
		if resp, err = enviarWSO2IS(ctx, method, endpoint, data); err != nil {
			return err
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		resetToken()
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("IS respondió HTTP %d: %s", resp.StatusCode, string(b))
	}

	if target != nil {
		return json.NewDecoder(resp.Body).Decode(target)
	}
	return nil
}

// enviarWSO2IS envía una petición a IS 7.3 con el token de la aplicación y
// devuelve la respuesta. Crea el cuerpo desde data en cada llamada para poder
// reenviar la petición. data es opcional.
func enviarWSO2IS(ctx context.Context, method, endpoint string, data []byte) (*http.Response, error) {
	bearer, err := getToken(ctx)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if data != nil {
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("Accept", "application/scim+json")
	if data != nil {
		req.Header.Set("Content-Type", "application/scim+json")
	}

	return isClient.Do(req)
}

func resetToken() {
	tokenMu.Lock()
	defer tokenMu.Unlock()
	token = ""
}

// getToken devuelve el token client_credentials, reutilizándolo hasta un minuto
// antes de que expire.
func getToken(ctx context.Context) (string, error) {
	tokenMu.Lock()
	defer tokenMu.Unlock()

	if token != "" && time.Now().Before(tokenExpiry) {
		return token, nil
	}

	if credentials == nil {
		parameterBasePath := "/" + beego.AppConfig.String("parameterStore") + "/" + beego.AppConfig.String("appname") + "/"
		secret, err := ssm.GetValueFromParameterStore(ctx, parameterBasePath+wso2IsSecret)
		if err != nil {
			logs.Error("Error getting Wso2IsSecret from parameter store: ", err)
			return "", fmt.Errorf("no se pudo leer el secreto de la aplicación de IS: %w", err)
		}

		client, err := ssm.GetValueFromParameterStore(ctx, parameterBasePath+wso2IsClient)
		if err != nil {
			logs.Error("Error getting Wso2IsClient from parameter store: ", err)
			return "", fmt.Errorf("no se pudo leer el client_id de la aplicación de IS: %w", err)
		}

		credentials = &clientCredentials{
			ClientID:     client,
			ClientSecret: secret,
		}
	}

	body, err := ClientCredentialsRequest(ctx, credentials.ClientID, credentials.ClientSecret, scimScopes)
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el token de IS: %w", err)
	}
	if body.AccessToken == "" {
		return "", fmt.Errorf("respuesta de token de IS inválida")
	}

	token = body.AccessToken
	tokenExpiry = time.Now().Add(time.Duration(body.ExpiresIn)*time.Second - time.Minute)
	return token, nil
}
