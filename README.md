# GO Mock Service

Este servicio permite desplegar mocks fácilmente configurables definidos a través de un archivo .json.

## Instrucciones de uso

### Sin build

En un servidor/pc que tenga docker instalado, correr los siguientes comandos:

- ./build_image.sh
- ./run_container.sh

Esto dejará corriendo el mock en el puerto especificado en el archivo .env y con el basepath indicado en el archivo config.json de configuración.
