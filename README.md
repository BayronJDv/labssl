# Lab SSL

Una aplicación de terminal interactiva para analizar certificados SSL de sitios web. Usa la API de SSLLabs para obtener información detallada sobre la seguridad SSL/TLS de un dominio.

## Qué es

Este es un proyecto hecho en Go que te permite verificar la configuración SSL de cualquier sitio web directamente desde tu terminal. En lugar de ir a la web de SSLLabs, tienes todo accesible desde la línea de comandos con una interfaz gráfica amigable.

## Características

- Analizar certificados SSL de dominios
- Ver detalles de los análisis anteriores
- Configurar parámetros como edad máxima del análisis, visibilidad pública, etc.
- Interfaz completamente en terminal con navegación por teclado

## Requisitos

- Go 1.18 o superior
- Conexión a internet para acceder a la API de SSLLabs

## Instalación

1. Clona el repositorio
2. Navega al directorio del proyecto
3. Ejecuta:

```bash
go build
```

## Uso

Para iniciar la aplicación:

```bash
./labssl
```

Una vez dentro, navega usando las flechas del teclado y presiona Enter para seleccionar opciones. Puedes analizar dominios, ver configuraciones o acceder a la ayuda.

## Estructura del Proyecto

- `main.go` - Punto de entrada de la aplicación
- `bubbletea/` - Lógica de la interfaz de usuario
  - `model.go` - Modelo de estado de la aplicación
  - `update.go` - Actualización del estado
  - `view.go` - Renderización de la interfaz
  - `analyze/` - Lógica para comunicarse con SSLLabs
    - `analyze.go` - Llamadas a la API
    - `utils.go` - Utilidades
- `style/` - Estilos y temas de la interfaz

## Dependencias

- BubbleTea - Framework para interfaces de terminal en Go
- Bubbles - Componentes reutilizables para BubbleTea

## Notas

La aplicación se comunica con la API de SSLLabs que puede tomar tiempo procesando análisis de certificados. Algunos escaneos pueden demorar varios segundos dependiendo de la carga del servicio.

## Licencia

Sin especificar por ahora.
