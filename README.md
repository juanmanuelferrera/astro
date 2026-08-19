# Astro

Cartas natales y curso en **dos tradiciones** —occidental y jyotiṣa— y **dos idiomas**, en un
único binario. Sin Python, sin red, sin instalar nada.

```
./astro
```

Abre el navegador en `http://localhost:8733`. Si ese puerto está ocupado, busca el siguiente libre.

## Las dos tradiciones

Un selector arriba cambia **todo a la vez**: el cálculo, el vocabulario, el dibujo y el curso.
En modo védico la palabra «Ascendente» no aparece en ninguna parte; en occidental no aparece
«Lagna». No se pueden mezclar por accidente.

La excepción es deliberada: la pestaña **Comparar** pone los dos zodíacos lado a lado y marca qué
cuerpos cambian de signo. El mayor riesgo de confusión convertido en la función que mejor enseña.

### Occidental

Zodíaco tropical · casas de Plácido y casas iguales en paralelo · aspectos por grados con orbes ·
rueda circular con desapilado angular de glifos y línea al grado real · regentes de casa.

**Lectura** por el Sistema de Palabras Clave de Margaret Hone: traduce la carta a frases literales
agrupadas por categorías y **señala las contradicciones**. No sintetiza — eso es del que lee, y el
módulo 10 explica por qué.

**Corrector del cálculo a mano**: escribes lo que te salió en cada paso y dice en cuál te
desviaste, sin darte el resultado bueno.

**Curso de 14 módulos**, más la vía corta de un cuarto de hora y el plan de una semana.

### Jyotiṣa

Zodíaco sidéreo con ayanāṁśa Lahiri · bhāvas de signo entero · 27 nakṣatras con pada y señor ·
dṛṣṭi por signo con los aspectos especiales de Marte, Júpiter y Saturno · dignidades védicas con
grados exactos, mūlatrikoṇa, combustión, dig-bala y gaṇḍānta · **diez vargas dibujadas** ·
daśās vimśottarī con bhuktis · kārakas de Jaimini con karakāṁśa · detección de yogas
(pañca-mahāpuruṣa, rāja-yoga por señores de kendra y trikona, nīca-bhaṅga, gaja-kesari,
kemadruma) · **tránsitos y Sade Sati** con la fase actual y su fecha de salida.

Carta cuadrada en los **dos estilos**, norte y sur de la India, con las vargas dibujadas por el
mismo componente para poder compararlas.

**Curso de 16 módulos.**

### Comunes a las dos

41.451 lugares (España completa, resto del mundo desde 15.000 habitantes) · husos horarios
históricos con el botón «¿por qué este desfase?» que muestra los cambios del país ·
guardar y recuperar cartas · imprimir a PDF · interfaz en español e inglés.

## Opciones

| | |
|---|---|
| `-puerto=9000` | puerto fijo; por defecto 8733, y salta al siguiente libre si está ocupado |
| `-abrir=false` | no abrir el navegador al arrancar |
| `-red` | aceptar conexiones de otros equipos de la red local |

Por defecto se ata solo a `127.0.0.1`: nadie de la red puede entrar.

## Precisión

Motor astronómico propio, **cero dependencias externas**. Meeus para Sol, Luna y tiempo sidéreo;
elementos de Standish con corrección de precesión al equinoccio de la fecha para los planetas;
casas de Plácido por bisección numérica; ayanāṁśa Lahiri por polinomio ajustado.

Verificado contra Swiss Ephemeris en 60 fechas aleatorias entre 1900 y 2040:

| | error medio | máximo |
|---|---|---|
| Sol | 0,12′ | 0,39′ |
| Luna | 0,39′ | 0,78′ |
| Mercurio, Venus, Marte | < 0,5′ | 1,7′ |
| Urano, Neptuno, Plutón | < 1,6′ | 4,4′ |
| Júpiter y Saturno | 3-5′ | 11,5′ |
| Casas de Plácido | — | 0,6″ |
| Ayanāṁśa Lahiri | — | 0,0002″ |

En 660 posiciones comprobadas, **el error no movió ningún planeta de signo**.

En latitudes polares detecta que Plácido no puede resolverse y cae a casas iguales avisándolo.

## Compilar

```
go build -o astro .
```

Para las demás plataformas:

```
GOOS=linux   GOARCH=amd64 go build -ldflags="-s -w" -o dist/astro-linux-amd64 .
GOOS=darwin  GOARCH=arm64 go build -ldflags="-s -w" -o dist/astro-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/astro-windows-amd64.exe .
```

`cmd/probar` y `cmd/jprobar` son las herramientas con las que se midió el motor contra
pyswisseph. Si tocas la parte astronómica, repite la verificación con ellas.

## Estructura

```
internal/efem/       astronomía compartida por las dos tradiciones
internal/occidental/ sistema de palabras clave
internal/jyotisha/   sidéreo, nakṣatras, vargas, daśās, yogas, tránsitos
internal/lugares/    41.451 poblaciones y husos históricos
internal/guardadas/  persistencia de cartas
web/                 interfaz y los dos cursos
```

## Datos de terceros

- Lugares: [GeoNames](https://www.geonames.org) (CC BY 4.0)
- Husos horarios: base IANA, embebida vía `time/tzdata`

## Lectura

Las dos tradiciones tienen motor de interpretación, y los dos componen la frase a partir de sus
piezas en lugar de guardar textos prefabricados.

**Occidental** — `internal/occidental/claves.go`. Sistema de Palabras Clave de Margaret Hone:
función del planeta + modo del signo + terreno de la casa + dignidad, más los regentes, que
convierten un aspecto suelto en un argumento. Detecta cuándo un planeta recibe a la vez aspectos
duros y blandos y lo declara como contradicción en lugar de elegir un lado.

**Jyotiṣa** — `internal/jyotisha/lectura.go`. Función del graha + modo del rāśi + terreno del
bhāva + nakṣatra + dignidad, y sobre todo la **cadena de señores**, que es como razona jyotiṣa:
de qué depende un asunto, y de qué depende aquello. Añade dṛṣṭi por signo entero, la daśā que
corre ahora y los kārakas de Jaimini. Marca como disputado el bhāva que reciben benéficos y
maléficos a la vez.

Ninguno de los dos sintetiza. Los dos dicen por qué.

## Pendiente

- Ashtakavarga y shadbala
- Pañcāṅga: tithi, vara, yoga, karaṇa
- Arudha padas y lagnas especiales de Jaimini
- Nodos verdaderos como alternativa a los medios
