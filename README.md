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
Para la peersistencia usada en este proyecto se ha decidido usar JSON. Se han creado dos JSON, uno para guardar los usuarios y otro para guardar los archivos.

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

### 2.5. paquete cmd
