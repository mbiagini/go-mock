# GO Mock Service

Este servicio permite desplegar mocks fácilmente configurables definidos a través de un archivo .json.

## Instrucciones de uso

### Ejecutables

En el directorio base del repositorio están los ejecutables zerver_lin.exe y zerver_win.exe. Cada uno puede ser ejecutado desde el SO correspondiente (ambos fueron compilados para una arquitectura de 64-bit).

### Compilación

Si se realizó alguna modificación al código (a cargo de cada desarrollador), se puede volver a compilar para generar los ejecutables corriendo el script "server_compile.sh" que se encuentra en el directorio principal.

### Docker

En el directorio "docker" se encuentran todos los archivos necesarios para hacer el build de la imagen Docker (script "build_image.sh") y, luego, levantar un container que utilice la misma (script "run_container.sh"). Dependiendo de dónde vaya a utilizar estos archivos (localmente o en un servidor), se debe configurar el archivo ".env" acordemente.

## Próximos pasos

### v1.2.0:

1. Agregar documentación para el agregado de endpoints al mock. Considerar agregar un .docx o un swagger con un endpoint POST para simular el agregado de un nuevo endpoint y documentarlo allí.
2. Agregar validación de la configuración de forma automática y completa utilizando la librería gslog.
3. Permitir configurar una política de balanceo entre las respuestas, es decir, que no se utilice un "discriminator" sino que se elija qué respuesta devolver según una política, ignorando el input recibido.
   1. Permitir balancear por round-robin.
   2. Permitir balancear por porcentajes.
   3. Permitir balancear random (a definir).
4. Permitir configurar una lista de discriminators.

### v2.0.0:

1. Agregar interfaz REST y configuración almacenada en DB.

### v2.1.0:

1. Analizar la posibilidad de armar una documentación automática de todos los casos posibles configurados.
   1. Para esto, para las condiciones que tengan una expresión regular compleja, se podría agregar un campo "description" para describir human-readable lo que va a matchear la condición.
   2. Que la docu la genere al pegarle con un GET o POST pasándole el/los endpoints que se quieran documentar.
