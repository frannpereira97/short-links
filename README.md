# short-links
Proyecto creado para realizar links mas cortos y aprender como funcionaba Go y Gorilla/Mux en general.

El desarollo Front-End esta realizado unicamente para probar la funcionalidad, hecho con templates de Boostrap. Las proximas actualizaciones van a estar enfocadas en separar el proyecto en Front-End y Back-End utilizando React.js

La idea del proyecto es poder utilizar un dominio corto o empresarial para poder realizar redirecciones a la pagina deseada por mas larga que sea.

## Por Mejorar/Agregar/Corregir
Este proyecto no fue creado para utilizar en produccion, por lo que le quedan muchas cosas que se podrian mejorar, sientanse libres de modificar el codigo a gusto y dejar recomendaciones.

### Mejoras
* **Seguridad** - En cuanto a seguridad, actualmente se basa en JWT para realizar las consultas (API) se podria agregar cookies o otro metodo de verificacion

* **Front** - Como aclare arriba el front esta utilizando templates modificadas de boostrap por lo que no es responsive por completo, mi idea es pasarlo a la libreria de ReactJs en un futuro

* **Formularios** - Editar los formularios para que restringan el tipo de datos que introducen los usuarios, asi como modificar los Select a eleccion (popular con los datos necesarios)

* **Debug** - Falta debugear el codigo y eliminar codigo repetido

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
Esto nos configurara para que se levante la aplicacion en:
```
localhost:4000\
```

### Kubernetes
Proximamente instrucciones para levantarlo en Kubernetes

## Construido con

Todavia en construccion

* [Golang](https://go.dev/doc/) - 
* [Javascript](hhttps://devdocs.io/javascript/) - 

## Autores

* **Franco Pereira** - *short-links* - [frannpereira97](https://github.com/frannpereira97)

