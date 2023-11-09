# short-links
Proyecto creado para realizar links mas cortos y aprender como funcionaba Go y Gorilla/Mux en general.

El desarollo Front-End esta realizado unicamente para probar la funcionalidad, hecho con templates de Boostrap. Las proximas actualizaciones van a estar enfocadas en separar el proyecto en Front-End y Back-End utilizando React.js

## Instrucciones

Para correr este proyecto debemos clonar el repositorio y tenemos la opcion de correr el codigo de forma manual si ya tenemos una base configurada.

Tambien podemos levantarlo con Docker modificando las variables de entorno ubicadas en .env o adjuntandolas al momento de correr el contenedor por separado.

Y si queremos solamente correr un comando y que funcione podemos utilizar: docker-compose up

Esto nos va a levantar automaticamente un contenedor de base de datos PostgreSQL y un contenedor con la aplicacion. 

Mira **Instalación** para saber cómo desplegar el proyecto paso por paso.


## Instalación

Primero debemos ejecutar el siguiente comando para clonar el repositorio:

```
git clone https://github.com/frannpereira97/short-links.git
```

Luego tenemos distintas formas de correr la aplicacion (Docker o Local):

### Local

Debemos asegurarnos de tener una base de datos PostgreSQL configurada, luego modificamos el archivo .env con los datos necesarios para la conexion.

Pasamos en una terminal con Go instalado y nos dirigimos a **.short-links/api** donde tenemos el archivo main.go y lo ejecutamos

```
go init
go run .
```
Y nos dirigimos hacia el dominio que hayamos especificado en las variables de entorno, por default:
```
localhost:4000/
```

### Docker

Para realizar la instalacion mediante Docker podemos optar por unicamente correr la aplicacion dirigiendonos hacia la carpeta **.short-links/api** y corriendo el contenedor con el siguiente comando:
```
docker run -e DB_ADDR="" \
           -e DB_PORT="" \
           -e DB_PASS="" \
           -e DB_USER="" \
           -e DB_NAME="" \
           -e DOMAIN="" \
           --name nombre_contenedor \
           imagen_docker
```

Si no tenemos base de datos configurada podemos simplemente dirigirnos al directorio **.short-links** utilizar el comando:
```
docker-compose up
```
## Construido con

Todavia en construccion

* [Golang](https://go.dev/doc/) - 
* [Javascript](hhttps://devdocs.io/javascript/) - 

## Autores

* **Franco Pereira** - *short-links* - [frannpereira97](https://github.com/frannpereira97)


