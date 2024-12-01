# PRACTICA 3
## 1. PUESTA EN MARCHA: DEFINICION DEL FRAMEWORK A UTILIZAR Y CREACION DE LA ESTRUCTURA DEL PROYECTO
### 1.1. DEFINICION DEL FRAMEWORK A UTILIZAR

Para la realización de la Api que se pide para la practica 3, en el lenguaje Go, se ha optado por usar Gin, que tiene una interfaz amigable y fácil de usar.

### 1.2. ESTRUCTURA DEL PROYECTO
El proyecto se ha estructurado de la siguiente manera:
- api -> contiene el codigo de las funciones de la api
- cmd -> contiene el codigo de la aplicacion principal
- database -> contiene la base de datos
- models -> contiene los modelos de la base de datos
- repository -> contiene el codigo de las funciones de la base de datos

### 1.3. CREACION DE LA BASE DE DATOS
Para la peersistencia usada en este proyecto se ha decidido usar JSON. Se han creado dos JSON, uno para guardar los usuarios y otro para guardar los archivos. En el config.toml de la carpeta config se especifica el directorio donde se guardaran los JSON.por defecto esta en database.

## 2. EXPLICACION DE LA API
### 2.1 paquete repository
#### 2.1.1. user_json.go
En este archivo encontramos las funciones que se encargan de guardar los usuarios en la persistencia, ademas de las funciones que se encargan de obtener los usuarios de la persistencia.
#### 2.1.2. json_json.go
Ese archivo tiene las funciones de crear los archivos json, obtenerlos, actualizarlos, borrarlos y listar los archivos de un usuario

### 2.2 paquete models
#### 2.2.1. user.go
Tenemos la estructura de un usuario, que contiene un nombre de usuario y su contraseña
#### 2.2.2. json.go
La estrcutra json solo cuenta con el contenido del archivo
#### 2.2.3. token.go
Esta es la estrcutura del token, cuenta con un codigo, el usuario al que pertenece y su tiempo de expiracion
#### 2.2.4 File.go
Aqui se define la estructura del archivo, que cuenta con su id y el json

### 2.3 paquete constants
En este paquete definimos las constantes que se usan en el proyecto

### 2.4. paquete api
En este paquete tenemos las funciones que definen la capa de negocio de la aplicacion
#### 2.4.1. authentication.go
En este archivo encontramos las funciones que se encargan de la autenticacion de los usuarios
#### 2.4.2. json_storage.go
En este archivo encontramos las funciones que se encargan de la persistencia de los archivos

### 2.5. paquete cmd
Encontramos el main, aqui se define la api y los distintos endpoints

## 3. EJECUCION
Para ejecutar la api, se ejecuta el siguiente comando:
```bash
go run cmd/main.go
```
## 4. TESTING
Se han realizado test a dos niveles, a nivel de http y a nivel de https:
### 4.1. Test de http
Para estos test se ha usado pytest, por ello se ha creado un entorno virtual con el siguiente comando:
```bash
python3 -m venv .venv
```
Luego se instala con pip las librerias pytest y requests. Y ejecutamos el comando:
```bash
pytest test/test_http.py
```
### 4.2. Test de https
Lo que hemos usado esta vez, es la extension de visual studio code thunder client, se han definido una operacion por cada endpoint y se puede ir cambiando cabeceras y body para ver como se comporta la api. Para usarlos, se tiene que tener instalada la extension e importar el archuvo json que se encuentra en la carpeta test.
