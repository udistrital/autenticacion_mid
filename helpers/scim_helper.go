package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/autenticacion_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/ssm"
)

// User lookups in WSO2 IS 7.3's SCIM2 API, as an M2M application.

const (
	customSchema = "urn:scim:schemas:extension:custom:User"
	scimScopes   = "internal_user_mgt_list internal_user_mgt_view"
	// request.GetWithContext sends this context value as the Authorization header.
	authorizationKey = "Authorization"
)

type scimUser struct {
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

type clientCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

var (
	tokenMu     sync.Mutex
	credentials *clientCredentials
	token       string
	tokenExpiry time.Time
	tokenClient = &http.Client{Timeout: 30 * time.Second}
)

// GetRolesUsuarioIS looks the user up by email when usuario has an @, else by
// username.
func GetRolesUsuarioIS(ctx context.Context, usuario string) (models.AtributosToken, error) {
	atributo := "userName"
	if strings.Contains(usuario, "@") {
		atributo = "emails"
	}
	atributos, err := buscarUsuario(ctx, atributo, usuario)
	if err != nil {
		logs.Error("Error al obtener los roles del usuario %s: %w", usuario, err)
		return models.AtributosToken{}, fmt.Errorf("Error al obtener los roles del usuario %s: %w", usuario, err)
	}
	return atributos, nil
}

// GetInfoByDocumentoIS is not wired yet: documentoToken and periods use GetInfoByDocumentoService.
func GetInfoByDocumentoIS(ctx context.Context, documento string) (models.AtributosToken, error) {
	atributos, err := buscarUsuario(ctx, customSchema+":documento", documento)
	if err != nil {
		return models.AtributosToken{}, fmt.Errorf("Error al obtener la información del documento %s: %w", documento, err)
	}
	return atributos, nil
}

// buscarUsuario returns no attributes ("Usuario no registrado") unless exactly
// one user matches.
func buscarUsuario(ctx context.Context, atributo, valor string) (models.AtributosToken, error) {
	vacio := models.AtributosToken{Usuario: &models.UsuarioData{}}

	bearer, err := getToken(ctx)
	if err != nil {
		return vacio, err
	}

	query := url.Values{}
	query.Set("filter", fmt.Sprintf("%s eq %s", atributo, scimString(valor)))
	query.Set("attributes", "userName,name.familyName,emails,groups,roles,"+customSchema)
	endpoint := beego.AppConfig.String("Wso2IsService") + "/scim2/Users?" + query.Encode()

	var list scimList
	status, err := request.GetWithContext(context.WithValue(ctx, authorizationKey, "Bearer "+bearer), endpoint, &list)
	if status == http.StatusUnauthorized {
		resetToken()
	}
	if err != nil {
		return vacio, err
	}
	if len(list.Resources) != 1 {
		return vacio, nil
	}
	return toAtributos(list.Resources[0]), nil
}

// toAtributos returns the rows vista_usuarios_atributos returned; role is
// Internal/everyone (implicit in IS 5.3), Internal/selfsignup if the user has
// IS's selfsignup role, then the user's groups.
func toAtributos(u scimUser) models.AtributosToken {
	roles := []string{"Internal/everyone"}
	for _, r := range u.Roles {
		if r.AudienceType == "organization" && r.Display == "selfsignup" {
			roles = append(roles, "Internal/selfsignup")
			break
		}
	}
	for _, g := range u.Groups {
		roles = append(roles, g.Display)
	}
	mail := ""
	for _, e := range u.Emails {
		if mail == "" || e.Primary {
			mail = e.Value
		}
	}

	data := &models.UsuarioData{}
	for _, a := range []models.Atributo{
		{Atributo: "role", Valor: strings.Join(roles, ",")},
		{Atributo: "sn", Valor: u.Name.FamilyName},
		{Atributo: "mail", Valor: mail},
		{Atributo: "documento", Valor: u.Custom.Documento},
		{Atributo: "documento_compuesto", Valor: u.Custom.DocumentoCompuesto},
	} {
		if a.Valor != "" {
			data.Atributos = append(data.Atributos, a)
		}
	}
	return models.AtributosToken{Usuario: data}
}

func scimString(v string) string {
	return `"` + strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(v) + `"`
}

func resetToken() {
	tokenMu.Lock()
	defer tokenMu.Unlock()
	token = ""
}

// getToken reuses the token until a minute before it expires.
func getToken(ctx context.Context) (string, error) {
	tokenMu.Lock()
	defer tokenMu.Unlock()

	if token != "" && time.Now().Before(tokenExpiry) {
		return token, nil
	}

	if credentials == nil {
		parameterBasePath := "/" + beego.AppConfig.String("parameterStore") + "/" + beego.AppConfig.String("appname") + "/"
		secret, err := ssm.GetValueFromParameterStore(ctx, parameterBasePath+beego.AppConfig.String("Wso2IsSecret"))
		if err != nil {
			logs.Error("Error getting Wso2IsSecret from parameter store: ", err)
			return "", fmt.Errorf("no se pudo leer el secreto de la aplicación de IS: %w", err)
		}

		client, err := ssm.GetValueFromParameterStore(ctx, parameterBasePath+beego.AppConfig.String("Wso2IsClient"))
		if err != nil {
			logs.Error("Error getting Wso2IsClient from parameter store: ", err)
			return "", fmt.Errorf("no se pudo leer el client_id de la aplicación de IS: %w", err)
		}

		credentials = &clientCredentials{
			ClientID:     client,
			ClientSecret: secret,
		}
	}

	form := url.Values{"grant_type": {"client_credentials"}, "scope": {scimScopes}}
	endpoint := beego.AppConfig.String("Wso2IsService") + "/oauth2/token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(credentials.ClientID, credentials.ClientSecret)

	resp, err := tokenClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el token de IS: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("no se pudo obtener el token de IS: HTTP %d", resp.StatusCode)
	}

	var body struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil || body.AccessToken == "" {
		return "", fmt.Errorf("respuesta de token de IS inválida")
	}

	token = body.AccessToken
	tokenExpiry = time.Now().Add(time.Duration(body.ExpiresIn)*time.Second - time.Minute)
	return token, nil
}
