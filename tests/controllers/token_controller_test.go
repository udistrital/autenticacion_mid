package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/udistrital/autenticacion_mid/models"
)

func TestEmailToken(t *testing.T) {
	endpoint := "http://localhost:8082/v1/token/emailToken"
	contentType := "application/json"
	body := models.Token{
		Email: "ccmendezt@udistrital.edu.co",
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
			t.Error("Error en EmailToken, Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("EmailToken Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error EmailToken:", err.Error())
		t.Fail()
	}
}

func TestUserRol(t *testing.T) {
	endpoint := "http://localhost:8082/v1/token/userRol"
	contentType := "application/json"
	body := models.UserName{
		User: "ccmendezt@udistrital.edu.co",
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
			t.Error("Error en UserRol, Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("UserRol Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error UserRol: ", err.Error())
		t.Fail()
	}
}
