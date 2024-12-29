# PRACTICA 4 LUIS BENITO LÓPEZ

## Explicación del proyecto

La practica consiste en el despliegue de un servicio web en Docker, este servicio web se encuentra dividio en 3 partes:

- lbroker.duckdns.org: que se encarga de recibir todas las peticiones que lleguen, procesarlas y redirigirlas a una de las otras dos partes. 
- lauth.duckdns.org: se encarga de la autenticacion y registro de los usuarios, ademas de la generacion y validacion de tokens de acceso.
- lfile.duckdns.org: se encarga de la subida, descarga, actualizacio, borrado y listado de archivos.

Aparte de los servicios explicados, encontramos dos nodos mas, uno jump y otro work, estos tienen la finalidad de conectarse a la red de los servicios por ssh. La maquina jump sirve para hacer de intermediaria entre el usuario y la maquina work, ya que se tienen que conectar primero a work. Encontramos dos usuarios, dev y op. op tiene acceso a todas las maquinas por ssh mientras que dev solo a work.

## Instalación

### Despliegue de los servicios
En el proyecto, encontramos un Makefile que nos facilita la instalacion y despliegue de los servicios. Para la compilacion de los servicios y la creacion de las imagenes docker, ejecutamos el siguiente comando:

```bash
make images
```

Para la creacion de las redes y el despliegue de los servicios, ejecutamos el siguiente comando:

```bash
make container
```
Para que el la resolucion del dominio de broker funcione bien, es necesario añadir la siguiente linea al archivo /etc/hosts:

```bash
172.17.0.2 lbroker.duckdns.org
```
Para la deetencion de los servicios y la eliminacion de las redes, ejecutamos el siguiente comando:

```bash
make remove
```
Para el borrado de los binarios y de las imagenes, ejecutamos el siguiente comando:

```bash
make clean
```
### Instalacion del ssh
Para el ssh, he configurado un archivo de hosts para acceder a las maquinas de manera mas sencilla:

```bash
sudo cp assets/config ~/.ssh/config
```
El archivo consiste en la definicion de los hosts y la clave para entrar, para ello es necesario tener las claves que estan en la carpeta keys en la carpeta .ssh



## Ejecucion de los test
Para ejecutar los test, es necesario crear un entorno virtual de python e instalar las librerias necesarias:

```bash
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

Una vez instaladas, en el Makefile existe una regla para ejecutarlos, es necesario tener el entorno virtual activado:

```bash
make run_test
```

Estos test sirven para probar la api en funcionamiento