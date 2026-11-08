package helpers

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego"
)

func TestClientCredentialsRequest(t *testing.T) {
	good := base64.StdEncoding.EncodeToString([]byte("cliente:secreto"))
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth2/token" || r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		r.ParseForm()
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Basic "+good || r.PostForm.Get("grant_type") != "client_credentials" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"invalid_client","error_description":"A valid OAuth client could not be found"}`))
			return
		}
		w.Write([]byte(`{"access_token":"tok","scope":"default","token_type":"Bearer","expires_in":3600}`))
	}))
	defer srv.Close()
	beego.AppConfig.Set("Wso2IsService", srv.URL)

	res, err := ClientCredentialsRequest(good)
	if err != nil || res.AccessToken != "tok" || res.ExpiresIn != 3600 || res.TokenType != "Bearer" {
		t.Fatalf("valid client: %+v, %v", res, err)
	}

	res, err = ClientCredentialsRequest(base64.StdEncoding.EncodeToString([]byte("cliente:otro")))
	if err == nil || res.AccessToken != "" {
		t.Fatalf("rejected client: want an error and no token, got %+v, %v", res, err)
	}
}
