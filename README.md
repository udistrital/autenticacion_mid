# autenticacion_mid

API MID del sistema de autenticación de la Universidad Distrital. Consulta los usuarios en WSO2 Identity Server 7.3 y:

- Entrega los datos y roles de un usuario para su token, por correo, nombre de usuario o documento.
- Entrega el código y el estado de un estudiante por su correo.
- Asigna y quita roles a los usuarios.
- Lista los periodos de roles de los usuarios.
- Entrega tokens client_credentials a los clientes configurados.

## Endpoints

Todos bajo `/v1`. La documentación completa está en `swagger/`.

| Método | Ruta | Descripción |
| -- | -- | -- |
| POST | `/token/userRol` | Datos y roles del usuario por correo o nombre de usuario |
| POST | `/token/documentoToken` | Datos y roles del usuario por documento |
| POST | `/token/emailToken` | Código y estado del estudiante por correo |
| POST | `/token/clientAuth` | Token client_credentials para un cliente |
| POST | `/rol/add` | Asigna un rol a un usuario |
| POST | `/rol/remove` | Quita un rol a un usuario |
| GET | `/rol/user/{documento}/periods` | Periodos de roles de un usuario |
| GET | `/rol/periods` | Periodos de roles de todos los usuarios |

Los roles de un usuario son `Internal/everyone`, `Internal/selfsignup` si tiene ese rol en IS, sus grupos en IS y `ESTUDIANTE` si su correo es de un estudiante. Asignar o quitar un rol agrega o quita al usuario del grupo con ese nombre.

Los grupos que tienen roles de IS asignados (de la Consola o de una aplicación) dan accesos que se administran en IS: no se devuelven como roles del usuario, y no se pueden asignar ni quitar, como los roles `Internal/*` y `Application/*` (el MID responde 400).

El MID guarda en memoria durante 10 minutos si cada grupo tiene roles, así que un cambio en los roles de un grupo puede tardar hasta 10 minutos en reflejarse. Si no puede consultar un grupo, lo omite de los roles del usuario y lo vuelve a consultar en la siguiente petición.

## Especificaciones Técnicas

- Go 1.26
- Beego 1.12
- WSO2 Identity Server 7.3 (API SCIM2)

## Variables de Entorno

```
AUTENTICACION_MID_HTTP_PORT: [Puerto de ejecución de la API]
WSO2_IS_URL:                 [URL base de WSO2 Identity Server]
CODE_BY_EMAIL_STUDENT_URL:   [Servicio de consulta de estudiantes por correo]
HISTORICO_ROLES_SERVICE:     [Servicio del histórico de roles de los usuarios]
TERCEROS_SERVICE:            [Servicio de terceros]
```

Las variables están en `conf/app.conf`.

## Ejecución del proyecto

```
# 1. Obtener el repositorio
git clone https://github.com/udistrital/autenticacion_mid.git
cd autenticacion_mid

# 2. Moverse a la rama develop
git checkout develop

# 3. Definir las variables de entorno y ejecutar
AUTENTICACION_MID_HTTP_PORT=8080 WSO2_IS_URL=... bee run
```

## Ejecución Pruebas

### Pruebas Unitarias

```
go test ./scim/ ./services/ ./controllers/
```

### Pruebas de Integración

`tests/controllers` prueba los endpoints contra la API en ejecución en `localhost`, con el puerto de `conf/app.conf`. Las pruebas de `/rol` asignan y quitan un rol a un usuario, y al terminar le devuelven el estado que tenía:

```
TEST_USUARIO=usuario@udistrital.edu.co TEST_ROL=ROL TEST_CORREO_ESTUDIANTE=estudiante@udistrital.edu.co TEST_DOCUMENTO=123456 go test -count=1 ./tests/...
```

## Licencia

This file is part of autenticacion_mid.

autenticacion_mid is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

autenticacion_mid is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with autenticacion_mid. If not, see https://www.gnu.org/licenses/.

## Estado CI

| Develop | Release 0.2.1 | Master |
| -- | -- | -- |
| [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/autenticacion_mid/status.svg?ref=refs/heads/develop)](https://hubci.portaloas.udistrital.edu.co/udistrital/autenticacion_mid/) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/autenticacion_mid/status.svg?ref=refs/heads/release/0.2.1)](https://hubci.portaloas.udistrital.edu.co/udistrital/autenticacion_mid/) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/autenticacion_mid/status.svg?ref=refs/heads/master)](https://hubci.portaloas.udistrital.edu.co/udistrital/autenticacion_mid/) |
