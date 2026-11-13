package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/require"
)

// baseURL es la URL de la API en ejecución, con el puerto que usa la API según
// conf/app.conf, o el de beego por defecto si está vacío.
var baseURL string

// Valores de las pruebas, que se definen con variables de entorno.
var (
	usuarioPrueba    string // TEST_USUARIO: correo o nombre de un usuario de IS
	rolPrueba        string // TEST_ROL: rol que se asigna y se quita al usuario
	correoEstudiante string // TEST_CORREO_ESTUDIANTE: correo de un estudiante
	documentoPrueba  string // TEST_DOCUMENTO: documento de un usuario de IS
)

func TestMain(m *testing.M) {
	if err := beego.LoadAppConfig("ini", "../../conf/app.conf"); err != nil {
		log.Fatal(err)
	}
	baseURL = fmt.Sprintf("http://localhost:%d/v1", beego.BConfig.Listen.HTTPPort)

	var faltantes []string
	for nombre, valor := range map[string]*string{
		"TEST_USUARIO":           &usuarioPrueba,
		"TEST_ROL":               &rolPrueba,
		"TEST_CORREO_ESTUDIANTE": &correoEstudiante,
		"TEST_DOCUMENTO":         &documentoPrueba,
	} {
		if *valor = os.Getenv(nombre); *valor == "" {
			faltantes = append(faltantes, nombre)
		}
	}
	if len(faltantes) > 0 {
		slices.Sort(faltantes)
		log.Fatalf("defina las variables de entorno %s", strings.Join(faltantes, ", "))
	}
	os.Exit(m.Run())
}

// postJSON envía body como JSON a la ruta de la API y devuelve la respuesta, que
// se cierra al terminar la prueba.
func postJSON(t *testing.T, ruta string, body any) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)
	resp, err := http.Post(baseURL+ruta, "application/json", bytes.NewReader(b))
	require.NoError(t, err)
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}
