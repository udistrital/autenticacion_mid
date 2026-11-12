package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/astaxie/beego"
)

// fakeUsuariosIS starts a fake IS that returns client_credentials tokens and
// answers SCIM2 user searches by email or username. Usernames are compared
// case-insensitively, as IS does.
func fakeUsuariosIS(t *testing.T, usuarios map[string]string, correos map[string][]string) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/token" {
			json.NewEncoder(w).Encode(map[string]any{"access_token": "t", "expires_in": 3600})
			return
		}
		filter := r.URL.Query().Get("filter")
		valor := strings.Trim(filter[strings.Index(filter, `"`):], `"`)
		var ids []string
		if strings.HasPrefix(filter, "emails") {
			ids = correos[valor]
		} else {
			for id, nombre := range usuarios {
				if strings.EqualFold(nombre, valor) {
					ids = append(ids, id)
				}
			}
		}
		var res []map[string]any
		for _, id := range ids {
			// familyName is the username, so tests can tell which account was returned.
			res = append(res, map[string]any{"id": id, "userName": usuarios[id], "name": map[string]string{"familyName": usuarios[id]}})
		}
		json.NewEncoder(w).Encode(map[string]any{"totalResults": len(res), "Resources": res})
	}))
	t.Cleanup(srv.Close)
	beego.AppConfig.Set("Wso2IsService", srv.URL)
	tokenMu.Lock()
	credentials = &clientCredentials{ClientID: "id", ClientSecret: "secret"}
	token = ""
	tokenMu.Unlock()
}

// TestGetGruposUsuarioISInstitucional checks that an institutional email with no
// account is looked up by its username, and that other domains and shared
// emails are not.
func TestGetGruposUsuarioISInstitucional(t *testing.T) {
	fakeUsuariosIS(t,
		map[string]string{"u1": "dsparedest", "u2": "dup1", "u3": "dup2", "u4": "dup"},
		map[string][]string{"dup@udistrital.edu.co": {"u2", "u3"}})
	ctx := context.Background()
	encontrado := func(usuario string) bool {
		u, err := GetGruposUsuarioIS(ctx, usuario)
		if err != nil {
			t.Fatalf("%s: %v", usuario, err)
		}
		return u != nil
	}

	if !encontrado("dsparedest@udistrital.edu.co") || !encontrado("DSparedest@UDistrital.edu.co") {
		t.Error("institutional email without account: not found by username")
	}
	if encontrado("dsparedest@gmail.com") {
		t.Error("another domain was mapped to the institutional username")
	}
	if encontrado("dup@udistrital.edu.co") {
		t.Error("a shared email fell back to a username")
	}
}

func TestUsuarioInstitucionalNombre(t *testing.T) {
	for in, want := range map[string]string{
		"dsparedest@udistrital.edu.co": "dsparedest",
		"x@UDISTRITAL.EDU.CO":          "x",
		"x@gmail.com":                  "",
		"x@correo.udistrital.edu.co":   "",
		"@udistrital.edu.co":           "",
		"dsparedest":                   "",
	} {
		if got, ok := usuarioInstitucional(in); got != want || ok != (want != "") {
			t.Errorf("usuarioInstitucional(%q) = %q, %v; want %q", in, got, ok, want)
		}
	}
}

// TestGetGruposUsuarioISDesempateToken checks that when several accounts share
// an email, the caller's own account is used, and none is used otherwise.
func TestGetGruposUsuarioISDesempateToken(t *testing.T) {
	fakeUsuariosIS(t,
		map[string]string{"9ef63ad4-0000-4000-8000-000000000001": "ana", "9ef63ad4-0000-4000-8000-000000000002": "ana.old", "u3": "otro"},
		map[string][]string{"ana@gmail.com": {"9ef63ad4-0000-4000-8000-000000000001", "9ef63ad4-0000-4000-8000-000000000002"}})

	elegido := func(sub string) string {
		ctx := context.Background()
		if sub != "" {
			ctx = context.WithValue(ctx, auditoriaUserKey, sub)
		}
		u, err := GetGruposUsuarioIS(ctx, "ana@gmail.com")
		if err != nil {
			t.Fatalf("sub %q: %v", sub, err)
		}
		if u == nil {
			return ""
		}
		return u.FamilyName
	}

	for sub, want := range map[string]string{
		"9ef63ad4-0000-4000-8000-000000000002": "ana.old", // the account ID
		"ANA":                                  "ana",     // the username, in any case
		"ana@carbon.super":                     "ana",     // the username with the tenant
		"PRIMARY/ana.old@carbon.super":         "ana.old", // the username with the user store and tenant
		"otro":                                 "",        // another user's account
		"f0f0f0f0-federated-login":             "",        // a federated login
		"":                                     "",        // no user, e.g. a client_credentials token
	} {
		if got := elegido(sub); got != want {
			t.Errorf("sub %q: got %q, want %q", sub, got, want)
		}
	}
}

// parcheGrupo is a SCIM PATCH that fakeGruposIS received.
type parcheGrupo struct {
	grupo string
	body  map[string]any
}

// fakeGruposIS starts a fake IS that returns client_credentials tokens, has a
// single group, PLANEACION (ID g1), and records the PATCHes it receives.
func fakeGruposIS(t *testing.T) *parcheGrupo {
	t.Helper()
	parche := &parcheGrupo{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oauth2/token":
			json.NewEncoder(w).Encode(map[string]any{"access_token": "t", "expires_in": 3600})
		case r.URL.Path == "/scim2/Groups" && r.Method == http.MethodGet:
			var res []map[string]any
			if r.URL.Query().Get("filter") == `displayName eq "PLANEACION"` {
				res = append(res, map[string]any{"id": "g1", "displayName": "PLANEACION"})
			}
			json.NewEncoder(w).Encode(map[string]any{"totalResults": len(res), "Resources": res})
		case strings.HasPrefix(r.URL.Path, "/scim2/Groups/") && r.Method == http.MethodPatch:
			parche.grupo = strings.TrimPrefix(r.URL.Path, "/scim2/Groups/")
			json.NewDecoder(r.Body).Decode(&parche.body)
			json.NewEncoder(w).Encode(map[string]any{"id": parche.grupo})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	beego.AppConfig.Set("Wso2IsService", srv.URL)
	tokenMu.Lock()
	credentials = &clientCredentials{ClientID: "id", ClientSecret: "secret"}
	token = ""
	tokenMu.Unlock()
	return parche
}

// operacion returns the single operation of the PATCH, failing the test if
// there is not exactly one.
func (p *parcheGrupo) operacion(t *testing.T) map[string]any {
	t.Helper()
	if p.grupo != "g1" {
		t.Errorf("patched group %q, want g1", p.grupo)
	}
	ops, _ := p.body["Operations"].([]any)
	if len(ops) != 1 {
		t.Fatalf("patch operations: %v", p.body)
	}
	return ops[0].(map[string]any)
}

// TestAgregarUsuarioGrupoIS checks that the user is added to the group with
// the given name through a SCIM PATCH, and that an unknown group fails.
func TestAgregarUsuarioGrupoIS(t *testing.T) {
	parche := fakeGruposIS(t)

	if err := AgregarUsuarioGrupoIS(context.Background(), "u1", "PLANEACION"); err != nil {
		t.Fatalf("existing group: %v", err)
	}
	op := parche.operacion(t)
	miembros, _ := op["value"].(map[string]any)["members"].([]any)
	if op["op"] != "add" || len(miembros) != 1 || miembros[0].(map[string]any)["value"] != "u1" {
		t.Errorf("patch operation: %v", op)
	}

	if err := AgregarUsuarioGrupoIS(context.Background(), "u1", "NO_EXISTE"); !errors.Is(err, ErrNoExiste) {
		t.Errorf("unknown group: got %v, want ErrNoExiste", err)
	}
}

// TestQuitarUsuarioGrupoIS checks that the user is removed from the group with
// the given name through a SCIM PATCH, and that an unknown group fails.
func TestQuitarUsuarioGrupoIS(t *testing.T) {
	parche := fakeGruposIS(t)

	if err := QuitarUsuarioGrupoIS(context.Background(), "u1", "PLANEACION"); err != nil {
		t.Fatalf("existing group: %v", err)
	}
	op := parche.operacion(t)
	if op["op"] != "remove" || op["path"] != `members[value eq "u1"]` {
		t.Errorf("patch operation: %v", op)
	}

	if err := QuitarUsuarioGrupoIS(context.Background(), "u1", "NO_EXISTE"); !errors.Is(err, ErrNoExiste) {
		t.Errorf("unknown group: got %v, want ErrNoExiste", err)
	}
}
