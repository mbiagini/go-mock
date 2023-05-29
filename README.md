# GO Mock Service

Este servicio permite desplegar mocks fácilmente configurables definidos a través de un archivo .json.

## Instrucciones de uso

### Ejecutables

En el directorio base del repositorio están los ejecutables zerver_lin.exe y zerver_win.exe. Cada uno puede ser ejecutado desde el SO correspondiente (ambos fueron compilados para una arquitectura de 64-bit).

### Compilación

Si se realizó alguna modificación al código (a cargo de cada desarrollador), se puede volver a compilar para generar los ejecutables corriendo el script "server_compile.sh" que se encuentra en el directorio principal.

### Docker

En el directorio "docker" se encuentran todos los archivos necesarios para hacer el build de la imagen Docker (script "build_image.sh") y, luego, levantar un container que utilice la misma (script "run_container.sh"). Dependiendo de dónde vaya a utilizar estos archivos (localmente o en un servidor), se debe configurar el archivo ".env" acordemente.
