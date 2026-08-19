# Curso de Astrología

Curso de astrología occidental con calculadora de cartas natales. Un único binario:
sin Python, sin red, sin instalar nada.

```
./astro
```

Abre el navegador en `http://localhost:8733`.

## Qué hace

- **Carta natal** — rueda dibujada, ángulos, planetas, aspectos por exactitud y regentes de casa.
- **Lectura** — traducción literal por el Sistema de Palabras Clave de Margaret Hone, agrupada
  por categorías, con las contradicciones señaladas. No sintetiza: eso es del que lee.
- **Curso** — los catorce módulos más la vía corta y el plan de una semana.
- **Cálculo a mano** — corrige paso a paso el levantamiento manual de una carta y dice en cuál
  te desviaste, sin dar el resultado bueno.

## Opciones

| | |
|---|---|
| `-puerto=9000` | puerto fijo (por defecto 8733, y salta al siguiente libre si está ocupado) |
| `-abrir=false` | no abrir el navegador al arrancar |
| `-red` | aceptar conexiones de otros equipos de la red local |

## Precisión

Motor astronómico propio (Meeus para Sol, Luna y tiempo sidéreo; elementos de Standish para los
planetas, con corrección de precesión al equinoccio de la fecha). Verificado contra Swiss
Ephemeris en 60 fechas aleatorias entre 1900 y 2040:

| | error medio | máximo |
|---|---|---|
| Sol | 0,12′ | 0,39′ |
| Luna | 0,39′ | 0,78′ |
| Mercurio, Venus, Marte | < 0,5′ | 1,7′ |
| Urano, Neptuno, Plutón | < 1,6′ | 4,4′ |
| Júpiter y Saturno | 3-5′ | 11,5′ |

En 660 posiciones comprobadas, **el error no movió ningún planeta de signo**.

Casas por Plácido resueltas por bisección numérica: exactas a 0,6″ en ocho ciudades de Quito a
Reikiavik. En latitudes polares detecta que el sistema no puede resolverse y cae a casas iguales
avisando.

## Compilar

```
go build -o astro .
```

Sin dependencias externas. Para las demás plataformas:

```
GOOS=linux  GOARCH=amd64 go build -ldflags="-s -w" -o dist/astro-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/astro-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/astro-windows-amd64.exe .
```

## Datos de terceros

- Lugares: 41.451 poblaciones de [GeoNames](https://www.geonames.org) (CC BY 4.0) — España
  completa y el resto del mundo desde 15.000 habitantes.
- Husos horarios: base IANA, embebida vía `time/tzdata`.
