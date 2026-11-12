package helpers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/udistrital/autenticacion_mid/models"
)

// fakeIS inicia un IS falso que entrega el token client_credentials "t" y
// atiende con mux las demás peticiones hechas con ese token, y dirige a él las
// peticiones a IS hasta que termine la prueba.
func fakeIS(t *testing.T, mux *http.ServeMux) *httptest.Server {
	t.Helper()
	raiz := http.NewServeMux()
	raiz.HandleFunc("POST /oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		responderJSON(t, w, map[string]any{"access_token": "t", "expires_in": 3600})
	})
	raiz.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer t" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		mux.ServeHTTP(w, r)
	})
	srv := httptest.NewServer(raiz)
	t.Cleanup(srv.Close)
	usarIS(t, srv.URL)
	tokenMu.Lock()
	credentials = &clientCredentials{ClientID: "id", ClientSecret: "secret"}
	token = ""
	tokenMu.Unlock()
	return srv
}

// fakeUsuariosIS inicia un IS falso que responde las búsquedas de usuarios
// SCIM2. usuarios asocia el ID de cada cuenta con su nombre de usuario; filtros
// asocia cada otro filtro SCIM con los IDs de las cuentas que coinciden. Los
// nombres de usuario se comparan sin distinguir mayúsculas, como en IS.
func fakeUsuariosIS(t *testing.T, usuarios map[string]string, filtros map[string][]string) {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /scim2/Users", func(w http.ResponseWriter, r *http.Request) {
		filter := r.URL.Query().Get("filter")
		valor := strings.Trim(filter[strings.Index(filter, `"`):], `"`)
		var ids []string
		if !strings.HasPrefix(filter, "userName") {
			ids = filtros[filter]
		} else {
			for id, nombre := range usuarios {
				if strings.EqualFold(nombre, valor) {
					ids = append(ids, id)
				}
			}
		}
		var res []map[string]any
		for _, id := range ids {
			// familyName es el nombre de usuario, para saber qué cuenta se devolvió.
			res = append(res, map[string]any{"id": id, "userName": usuarios[id], "name": map[string]string{"familyName": usuarios[id]}})
		}
		responderJSON(t, w, map[string]any{"totalResults": len(res), "Resources": res})
	})
	fakeIS(t, mux)
}

// TestGetGruposUsuarioISDesempateToken comprueba que si varias cuentas comparten
// un correo se usa la del usuario del token, y ninguna en otro caso.
func TestGetGruposUsuarioISDesempateToken(t *testing.T) {
	fakeUsuariosIS(t,
		map[string]string{"9ef63ad4-0000-4000-8000-000000000001": "ana", "9ef63ad4-0000-4000-8000-000000000002": "ana.old", "u3": "otro"},
		map[string][]string{`emails eq "ana@gmail.com"`: {"9ef63ad4-0000-4000-8000-000000000001", "9ef63ad4-0000-4000-8000-000000000002"}})

	elegido := func(sub string) string {
		ctx := context.Background()
		if sub != "" {
			ctx = context.WithValue(ctx, auditoriaUserKey, sub)
		}
		u, err := GetGruposUsuarioIS(ctx, "ana@gmail.com")
		require.NoError(t, err, "sub %q", sub)
		if u == nil {
			return ""
		}
		return u.FamilyName
	}

	for sub, want := range map[string]string{
		"9ef63ad4-0000-4000-8000-000000000002": "ana.old", // el ID de la cuenta
		"ANA":                                  "ana",     // el nombre de usuario, sin importar mayúsculas
		"ana@carbon.super":                     "ana",     // el nombre de usuario con el tenant
		"PRIMARY/ana.old@carbon.super":         "ana.old", // el nombre de usuario con el user store y el tenant
		"otro":                                 "",        // la cuenta de otro usuario
		"f0f0f0f0-federated-login":             "",        // un inicio de sesión federado
		"":                                     "",        // sin usuario, p. ej. un token client_credentials
	} {
		assert.Equal(t, want, elegido(sub), "sub %q", sub)
	}
}

func TestGetInfoByDocumentoIS(t *testing.T) {
	fakeUsuariosIS(t,
		map[string]string{"u1": "ana", "u2": "dup1", "u3": "dup2"},
		map[string][]string{
			customSchema + `:documento eq "100"`: {"u1"},
			customSchema + `:documento eq "200"`: {"u2", "u3"},
		})

	for documento, want := range map[string]string{
		"100": "ana", // una sola cuenta
		"200": "",    // varias cuentas
		"300": "",    // ninguna cuenta
	} {
		u, err := GetInfoByDocumentoIS(context.Background(), documento)
		require.NoError(t, err, "documento %s", documento)
		got := ""
		if u != nil {
			got = u.FamilyName
		}
		assert.Equal(t, want, got, "documento %s", documento)
	}
}

// parcheGrupo es un PATCH de SCIM que recibió fakeGruposIS.
type parcheGrupo struct {
	grupo string
	body  map[string]any
}

// fakeGruposIS inicia un IS falso que tiene un solo grupo, PLANEACION (ID g1),
// y registra los PATCH que recibe.
func fakeGruposIS(t *testing.T) *parcheGrupo {
	t.Helper()
	parche := &parcheGrupo{}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /scim2/Groups", func(w http.ResponseWriter, r *http.Request) {
		var res []map[string]any
		if r.URL.Query().Get("filter") == `displayName eq "PLANEACION"` {
			res = append(res, map[string]any{"id": "g1", "displayName": "PLANEACION"})
		}
		responderJSON(t, w, map[string]any{"totalResults": len(res), "Resources": res})
	})
	mux.HandleFunc("PATCH /scim2/Groups/{id}", func(w http.ResponseWriter, r *http.Request) {
		parche.grupo = r.PathValue("id")
		if err := json.NewDecoder(r.Body).Decode(&parche.body); err != nil {
			t.Errorf("cuerpo del PATCH: %v", err)
		}
		responderJSON(t, w, map[string]any{"id": parche.grupo})
	})
	fakeIS(t, mux)
	return parche
}

// operacion devuelve la única operación del PATCH, y hace fallar la prueba si el
// PATCH no fue a g1 o no tiene exactamente una operación.
func (p *parcheGrupo) operacion(t *testing.T) map[string]any {
	t.Helper()
	assert.Equal(t, "g1", p.grupo)
	ops, _ := p.body["Operations"].([]any)
	require.Len(t, ops, 1, "operaciones del PATCH: %v", p.body)
	return ops[0].(map[string]any)
}

func TestAgregarUsuarioGrupoIS(t *testing.T) {
	parche := fakeGruposIS(t)

	usuario := &models.UsuarioIS{ID: "u1", UserName: "ana"}
	require.NoError(t, AgregarUsuarioGrupoIS(context.Background(), usuario, "PLANEACION"))
	assert.Equal(t, map[string]any{
		"op":    "add",
		"value": map[string]any{"members": []any{map[string]any{"value": "u1", "display": "ana"}}},
	}, parche.operacion(t))

	assert.ErrorIs(t, AgregarUsuarioGrupoIS(context.Background(), usuario, "NO_EXISTE"), ErrGrupoNoExiste)
}

func TestQuitarUsuarioGrupoIS(t *testing.T) {
	parche := fakeGruposIS(t)

	require.NoError(t, QuitarUsuarioGrupoIS(context.Background(), "u1", "PLANEACION"))
	assert.Equal(t, map[string]any{
		"op":   "remove",
		"path": `members[value eq "u1"]`,
	}, parche.operacion(t))

	assert.ErrorIs(t, QuitarUsuarioGrupoIS(context.Background(), "u1", "NO_EXISTE"), ErrGrupoNoExiste)
}

func TestClientCredentialsRequest(t *testing.T) {
	good := base64.StdEncoding.EncodeToString([]byte("cliente:secreto"))
	var scope string
	mux := http.NewServeMux()
	mux.HandleFunc("POST /oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("formulario del token: %v", err)
		}
		scope = r.PostForm.Get("scope")
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Basic "+good || r.PostForm.Get("grant_type") != "client_credentials" {
			w.WriteHeader(http.StatusUnauthorized)
			responderJSON(t, w, json.RawMessage(`{"error":"invalid_client","error_description":"A valid OAuth client could not be found"}`))
			return
		}
		responderJSON(t, w, json.RawMessage(`{"access_token":"tok","scope":"default","token_type":"Bearer","expires_in":3600}`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	usarIS(t, srv.URL)
	ctx := context.Background()

	res, err := ClientCredentialsRequest(ctx, "cliente", "secreto", "")
	require.NoError(t, err)
	assert.Equal(t, models.ClientCredentialsResponse{AccessToken: "tok", ExpiresIn: 3600, Scope: "default", TokenType: "Bearer"}, res)
	assert.Empty(t, scope)

	_, err = ClientCredentialsRequest(ctx, "cliente", "secreto", "a b")
	require.NoError(t, err)
	assert.Equal(t, "a b", scope)

	res, err = ClientCredentialsRequest(ctx, "cliente", "otro", "")
	assert.Error(t, err)
	assert.Empty(t, res.AccessToken)
}

func TestCallWSO2ISErrorDetail(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /scim2/Groups/g1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		responderJSON(t, w, json.RawMessage(`{"schemas":["urn:ietf:params:scim:api:messages:2.0:Error"],"detail":"falta el display","status":"400"}`))
	})
	srv := fakeIS(t, mux)

	err := callWSO2IS(context.Background(), http.MethodPatch, srv.URL+"/scim2/Groups/g1", map[string]any{}, nil)
	require.Error(t, err)
	assert.ErrorContains(t, err, "IS respondió HTTP 400: ")
	assert.ErrorContains(t, err, `"detail":"falta el display"`)
}

// TestCallWSO2ISReintento comprueba que si IS rechaza el token guardado, la
// petición se reintenta una vez con un token nuevo, y que no se reintenta más.
func TestCallWSO2ISReintento(t *testing.T) {
	var tokens, peticiones int
	valido := "nuevo"
	mux := http.NewServeMux()
	mux.HandleFunc("POST /oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		tokens++
		responderJSON(t, w, map[string]any{"access_token": "nuevo", "expires_in": 3600})
	})
	mux.HandleFunc("GET /scim2/Users", func(w http.ResponseWriter, r *http.Request) {
		peticiones++
		if r.Header.Get("Authorization") != "Bearer "+valido {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		responderJSON(t, w, json.RawMessage(`{"Resources":[]}`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	usarIS(t, srv.URL)
	conToken := func(guardado string) {
		tokenMu.Lock()
		credentials = &clientCredentials{ClientID: "id", ClientSecret: "secret"}
		token, tokenExpiry = guardado, time.Now().Add(time.Hour)
		tokenMu.Unlock()
		tokens, peticiones = 0, 0
	}

	conToken("revocado")
	var list scimList
	require.NoError(t, callWSO2IS(context.Background(), http.MethodGet, srv.URL+"/scim2/Users", nil, &list))
	assert.Equal(t, 1, tokens, "tokens pedidos")
	assert.Equal(t, 2, peticiones, "peticiones a IS")

	conToken("revocado")
	valido = "ninguno"
	assert.ErrorContains(t, callWSO2IS(context.Background(), http.MethodGet, srv.URL+"/scim2/Users", nil, &list), "IS respondió HTTP 401")
	assert.Equal(t, 2, peticiones, "peticiones a IS")
}

// responderJSON escribe v como JSON en la respuesta de un IS falso, y hace fallar
// la prueba si no puede.
func responderJSON(t *testing.T, w http.ResponseWriter, v any) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Errorf("respuesta del IS falso: %v", err)
	}
}

func TestToUsuarioGrupos(t *testing.T) {
	var u scimUser
	require.NoError(t, json.Unmarshal([]byte(`{
		"groups": [{"display": "PLANEACION"}, {"display": "Internal/admin"}, {"display": "Application/app"}, {"display": "application/otra"}, {"display": "ADMIN_CAMPUS"}],
		"roles": [{"display": "selfsignup", "audienceType": "organization"}, {"display": "admin", "audienceType": "organization"}]
	}`), &u))
	usuario := toUsuario(u)
	assert.Equal(t, []string{"Internal/everyone", "Internal/selfsignup"}, usuario.Roles)
	assert.Equal(t, []string{"PLANEACION", "ADMIN_CAMPUS"}, usuario.Grupos)
}

// usarIS dirige las peticiones a IS a url hasta que termine la prueba.
func usarIS(t *testing.T, url string) {
	t.Helper()
	anterior := wso2ISBasePath
	wso2ISBasePath = url
	t.Cleanup(func() { wso2ISBasePath = anterior })
}
