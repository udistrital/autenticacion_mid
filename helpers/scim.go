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

// User lookups in IS 7.3's SCIM2 API, using the application's client credentials.

const (
	customSchema = "urn:scim:schemas:extension:custom:User"
	// Institutional emails are <username>@udistrital.edu.co.
	dominioInstitucional = "udistrital.edu.co"
	scimScopes           = "internal_user_mgt_list internal_user_mgt_view"
	// Context key request.GetWithContext sends as the Authorization header.
	authorizationKey = "Authorization"
	// Context key where utils_oas/auditoria stores the caller's token sub.
	auditoriaUserKey = "user"
)

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

// GetGruposUsuarioIS returns the user's info and groups from IS 7.3 by email or
// username, or nil if no single account matches. An institutional email with no
// account is looked up by its username, and when several accounts share an
// email, the caller's own account is used.
func GetGruposUsuarioIS(ctx context.Context, usuario string) (*models.Usuario, error) {
	atributo := "userName"
	if strings.Contains(usuario, "@") {
		atributo = "emails"
	}
	usuarios, err := buscarUsuarios(ctx, atributo, usuario)
	if nombre, ok := usuarioInstitucional(usuario); err == nil && len(usuarios) == 0 && ok {
		usuarios, err = buscarUsuarios(ctx, "userName", nombre)
		if err == nil && len(usuarios) == 1 {
			logs.Info("Ninguna cuenta tiene el email %s: se usa la cuenta %s (%s) por su nombre de usuario",
				usuario, usuarios[0].UserName, usuarios[0].ID)
		}
	}
	if err != nil {
		logs.Error("Error al obtener los roles del usuario %s: %w", usuario, err)
		return nil, fmt.Errorf("Error al obtener los roles del usuario %s: %w", usuario, err)
	}
	if len(usuarios) > 1 {
		if u := cuentaDelToken(ctx, usuarios); u != nil {
			logs.Info("%d cuentas tienen el email %s: se usa %s (%s), la del token",
				len(usuarios), usuario, u.UserName, u.ID)
			usuarios = []scimUser{*u}
		} else {
			sub, _ := ctx.Value(auditoriaUserKey).(string)
			if sub == "" {
				sub = "sin usuario"
			}
			logs.Warning("%d cuentas tienen el email %s y ninguna es la del token (%s): usuario no registrado",
				len(usuarios), usuario, sub)
		}
	}
	if len(usuarios) != 1 {
		return nil, nil
	}
	return toUsuario(usuarios[0]), nil
}

// cuentaDelToken returns the account in usuarios that the caller's token belongs
// to, or nil. The token's sub is the account's ID or its username, which IS may
// prefix with PRIMARY/ and suffix with @carbon.super. Federated logins and
// client_credentials tokens match no account.
func cuentaDelToken(ctx context.Context, usuarios []scimUser) *scimUser {
	sub, _ := ctx.Value(auditoriaUserKey).(string)
	if sub == "" {
		return nil
	}
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

// usuarioInstitucional returns the username of an institutional email, the part
// before the @. It returns false for other domains, whose part before the @ may
// be someone else's username.
func usuarioInstitucional(email string) (string, bool) {
	i := strings.LastIndex(email, "@")
	if i <= 0 || !strings.EqualFold(email[i+1:], dominioInstitucional) {
		return "", false
	}
	return email[:i], true
}

// buscarUsuarios returns the users whose atributo equals valor.
func buscarUsuarios(ctx context.Context, atributo, valor string) ([]scimUser, error) {
	bearer, err := getToken(ctx)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("filter", fmt.Sprintf("%s eq %s", atributo, scimString(valor)))
	query.Set("attributes", "id,userName,name.familyName,emails,groups,roles,"+customSchema)
	endpoint := beego.AppConfig.String("Wso2IsService") + "/scim2/Users?" + query.Encode()

	var list scimList
	status, err := request.GetWithContext(context.WithValue(ctx, authorizationKey, "Bearer "+bearer), endpoint, &list)
	if status == http.StatusUnauthorized {
		resetToken()
	}
	if err != nil {
		return nil, err
	}
	return list.Resources, nil
}

// toUsuario maps a SCIM user to a models.Usuario. Its roles are
// Internal/everyone, Internal/selfsignup if the user has IS's selfsignup role,
// and the user's groups.
func toUsuario(u scimUser) *models.Usuario {
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

	return &models.Usuario{
		Roles:              roles,
		FamilyName:         u.Name.FamilyName,
		Email:              mail,
		Documento:          u.Custom.Documento,
		DocumentoCompuesto: u.Custom.DocumentoCompuesto,
	}
}

func scimString(v string) string {
	return `"` + strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(v) + `"`
}

func resetToken() {
	tokenMu.Lock()
	defer tokenMu.Unlock()
	token = ""
}

// getToken returns the client_credentials token, reusing it until a minute
// before it expires.
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
