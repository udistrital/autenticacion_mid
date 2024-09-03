# autenticacion_mid

API Mid para el sistema de autenticación Universidad Distrital.

Actualmente la API valida:

- Código de los estudiantes por usuario
- Validación de usuario por correo institucional
- Validación de usuario por documento de identidad
- Adición de roles a usuarios
- Eliminación de roles a usuarios

## Especificaciones Técnicas

### Tecnologías Implementadas y Versiones

- Golang
- BeeGo

## Variables de Entorno

```
  AUTENTICACION_MID_HTTP_PORT: [Puerto de ejecución API]
  CODE_BY_EMAIL_STUDENT_URL: [Servicio WSO2 de consulta de estudiantes]
  WSO2_AUTH_SERVICE: [Servicio WSO2 de Autenticación]
  WSO2_USER_SERVICE: [Servicio WSO2 de Usuarios]
  HISTORICO_ROLES_SERVICE: [Servicio para gestión del histórico de roles de los usuarios]
  TERCEROS_SERVICE: [Servicio de terceros]
```

**NOTA:** Las variables se pueden ver en el fichero conf/app.conf y están identificadas con AUTENTICACION_MID_HTTP_PORT...


## Ejecución del proyecto
```
#1. Obtener el repositorio con Go
go get github.com/udistrital/autenticacion_mid

#2. Moverse a la carpeta del repositorio
cd $GOPATH/src/github.com/udistrital/autenticacion_mid

# 3. Moverse a la rama **develop**
git pull origin develop && git checkout develop

# 4. alimentar todas las variables de entorno que utiliza el proyecto.
AUTENTICACION_MID_HTTP_PORT=8080 CODE_BY_EMAIL_STUDENT_URL=some_value bee run
```


## Ejecución Pruebas

### Pruebas Unitarias

#### TokenController

- **TestEmailToken:** <span style="color: #4cc61e;"><b>Test OK</b></span>
![TestEmailToken](tests/Unit%20Test/TestEmailToken.png)

- **TestUserRol:** <span style="color: #4cc61e;"><b>Test OK</b></span>
![TestUserRol](tests/Unit%20Test/TestUserRol.png)

- **TestDocumentoToken:** <span style="color: #4cc61e;"><b>Test OK</b></span>
![TestDocumentoToken](tests/Unit%20Test/TestDocumentoToken.png)

#### RolController

- **TestAddRol:** <span style="color: #4cc61e;"><b>Test OK</b></span>
![TestAddRol](tests/Unit%20Test/TestAddRol.png)

- **TestRemoveRol:** <span style="color: #4cc61e;"><b>Test OK</b></span>
![TestRemoveRol](tests/Unit%20Test/TestRemoveRol.png)

## Licencia

This file is part of autenticacion_mid.

autenticacion_mid is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

planeacion_mid is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with planeacion_mid. If not, see https://www.gnu.org/licenses/.

## Estado CI

| Develop | Release 0.0.1 | Master |
| -- | -- | -- |
| [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/catalogo_elementos_crud/status.svg?ref=refs/heads/develop)](https://hubci.portaloas.udistrital.edu.co/udistrital/catalogo_elementos_crud/) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/catalogo_elementos_crud/status.svg?ref=refs/heads/release/0.0.1)](https://hubci.portaloas.udistrital.edu.co/udistrital/catalogo_elementos_crud/) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/catalogo_elementos_crud/status.svg?ref=refs/heads/master)](https://hubci.portaloas.udistrital.edu.co/udistrital/catalogo_elementos_crud/) |
