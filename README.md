# Analizador Léxico, Sintáctico y Semántico Multi-Lenguaje

**Curso:** [Nombre de tu Materia]
**Profesor:** [Nombre de tu Profesor]
**Alumno:** [Tu Nombre Completo]
**Fecha:** Junio 2025

## Descripción del Proyecto

Este proyecto es un analizador full-stack desarrollado con React para el frontend y Go (Golang) para el backend.

El sistema es capaz de realizar análisis léxico, sintáctico y semántico para dos lenguajes diferentes:
1.  Un **lenguaje simplificado tipo C** con bucles `do-while`.
2.  Un **subconjunto básico de Swift**, con soporte para variables (`var`), constantes (`let`) y anotaciones de tipo.

El usuario puede seleccionar el lenguaje a analizar a través de una interfaz interactiva.

## Tecnologías Utilizadas
* **Frontend:** React
* **Backend:** Go (Golang)

## Cómo Ejecutar el Proyecto

Es necesario tener **Go** y **Node.js** (con npm) instalados.

1.  Clonar el repositorio:
    ```bash
    git clone [https://github.com/dolthub/dolt](https://github.com/dolthub/dolt)
    cd [Nombre de tu repositorio]
    ```

2.  Iniciar el **Backend** (Servidor en Go):
    ```bash
    # Navega a la carpeta del backend
    cd analyzerBack

    # Ejecuta el servidor
    go run main.go
    ```
    El servidor se iniciará en `http://localhost:8080`.

3.  Iniciar el **Frontend** (App en React):
    *Abrir una nueva terminal para este paso.*
    ```bash
    # Navega a la carpeta del frontend
    cd analyzer_front

    # Instala las dependencias (solo la primera vez)
    npm install

    # Inicia la aplicación de React
    npm start
    ```
    La aplicación se abrirá automáticamente en `http://localhost:3000`.