package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/udistrital/autenticacion_mid/models"
)

func TestAddRol(t *testing.T) {
	endpoint := "http://localhost:8082/v1/rol/add"
	contentType := "application/json"
	body := models.UpdateRol{
		User: "ccmendezt@udistrital.edu.co",
		Rol:  "PLANEACION",
	}

	// Convertir el cuerpo a formato JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Error("Error al convertir el cuerpo a JSON:", err.Error())
		t.Fail()
		return
	}

	// Crear un io.Reader con los datos JSON
	bodyReader := bytes.NewReader(bodyBytes)

	if response, err := http.Post(endpoint, contentType, bodyReader); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error en AddRol, Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("AddRol Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error AddRol:", err.Error())
		t.Fail()
	}
}

func TestRemoveRol(t *testing.T) {
	endpoint := "http://localhost:8082/v1/rol/remove"
	contentType := "application/json"
	body := models.UpdateRol{
		User: "ccmendezt@udistrital.edu.co",
		Rol:  "PLANEACION",
	}

	// Convertir el cuerpo a formato JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Error("Error al convertir el cuerpo a JSON:", err.Error())
		t.Fail()
		return
	}

	// Crear un io.Reader con los datos JSON
	bodyReader := bytes.NewReader(bodyBytes)

	if response, err := http.Post(endpoint, contentType, bodyReader); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error en AddRol, Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("AddRol Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error AddRol:", err.Error())
		t.Fail()
	}
}
