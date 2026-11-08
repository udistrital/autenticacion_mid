package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/astaxie/beego"
)

// fakeUsuariosIS answers client_credentials tokens and SCIM2 user searches
// (emails eq / userName eq, usernames compared case-insensitively as IS does).
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
			res = append(res, map[string]any{"id": id, "userName": usuarios[id]})
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

// An institutional email no account has is looked up by its username (Outlook
// logins); other domains and shared emails are not.
func TestGetRolesUsuarioISInstitucional(t *testing.T) {
	fakeUsuariosIS(t,
		map[string]string{"u1": "dsparedest", "u2": "dup1", "u3": "dup2", "u4": "dup"},
		map[string][]string{"dup@udistrital.edu.co": {"u2", "u3"}})
	ctx := context.Background()
	encontrado := func(usuario string) bool {
		at, err := GetRolesUsuarioIS(ctx, usuario)
		if err != nil {
			t.Fatalf("%s: %v", usuario, err)
		}
		return len(at.Usuario.Atributos) > 0
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
