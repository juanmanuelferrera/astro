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

## Pañcāṅga y lagnas especiales

Los cinco miembros del calendario hindú: tithi con su quincena, vāra, nakṣatra
con su pada, yoga y karaṇa, cada uno con lo que lleva recorrido. El tithi y el
karaṇa salen de la diferencia entre la Luna y el Sol, así que el ayanāṁśa se
cancela; el yoga sale de la suma y no se cancela.

Los tres lagnas especiales —Bhāva, Horā y Ghaṭī— se cuentan como tiempo
transcurrido desde el amanecer, así que hacen falta salida y puesta del Sol.
Están en `internal/efem/orto.go`, por el método de Meeus con dos iteraciones:
el error queda por debajo del medio minuto y en latitudes polares devuelve que
no hay orto en lugar de inventarse uno.

Y los arudha padas de los doce bhāvas, con el Arudha Lagna y el Upapada.

## Aṣṭakavarga

Las ocho tablas de Parāśara, con su bhinnāṣṭakavarga por graha y el
sarvāṣṭakavarga sumado. Las tablas se comprueban solas: los siete BAV tienen
que sumar exactamente 337, y hay un test que falla si al teclearlas se coló un
número de más o de menos.

## Ṣaḍbala

Las seis fuerzas en virūpas: sthāna (con ucca, saptavargaja sobre siete vargas
y amistad compuesta, ojayugma, kendrādi y drekkāṇa), dig, kāla (nathonnatha,
pakṣa, tribhāga, señores del día y de la hora, ayana), cheṣṭā por los ocho
estados del movimiento, naisargika y dṛk.

Y el **yuddha bala**: cuando dos de los cinco planetas quedan a menos de un
grado están en guerra, y la diferencia de fuerza se reparte según sus diámetros
aparentes. Gana el que más ṣaḍbala lleva acumulado — la tradición usa la
latitud, y ese cambio está declarado en el código. Sale en 8 de cada 100 cartas.

Aparte va el **bhāva bala**, que contesta otra pregunta: no cuánto puede un
graha sino cuánto puede un asunto. Tres partes: la fuerza del señor de la casa
—que es la que pesa con diferencia—, si la ocupan los grahas que esa casa pide,
y quién la mira.

Lo que se muestra no es la cifra bruta sino la razón entre lo que cada graha
saca y lo que se le exige, porque el listón es distinto para cada uno.

## Nodo medio o verdadero

Por defecto el medio, que es lo más extendido. La casilla de la pestaña de
carta cambia al verdadero, que oscila alrededor del medio hasta grado y medio
y mueve a Rāhu de pada. Las daśās no cambian: cuelgan del nakṣatra de la Luna.

## Predicción occidental

El módulo 12 enseña tránsitos, progresiones secundarias y revolución solar, y
hasta ahora el programa lo explicaba sin hacerlo. Ya están las tres, en
`internal/occidental/prediccion.go`.

- **Tránsitos** — solo los lentos, de Marte a Plutón: los rápidos marcan días y
  no periodos. Orbe estrecho a propósito (1,5° los aspectos mayores, 0,4° los
  menores; con el mismo orbe la lista se llenaba de sesquicuadraturas y el
  tránsito que importaba quedaba enterrado). Dice si aplica o separa, y si el
  planeta va a retrogradar encima del punto y pasar tres veces.
- **Progresiones secundarias** — un día después del nacimiento por año de vida.
- **Revolución solar** — el instante en que el Sol vuelve a su grado natal, por
  bisección.

Y lo que las une, que es lo que el módulo 12 pide de verdad: **las
convergencias**. Un periodo importa cuando dos técnicas DISTINTAS señalan el
mismo punto natal. Dos progresiones sobre el mismo sitio son una voz repetida y
no cuentan — es la regla del módulo 9 aplicada al tiempo. Cuando no coincide
nada, lo dice en lugar de rellenar.

## La pestaña de comparar

Es la única pantalla donde las dos tradiciones aparecen juntas, y enseña las
**dos** diferencias, no una:

- El **zodíaco**: tropical contra sidéreo, con el ayanāṁśa que los separa y qué
  cuerpos cambian de signo por él.
- Las **casas**: Plácido, que las hace desiguales y depende de la latitud,
  contra el signo entero. Esta suele pesar más y casi nunca se cuenta. En la
  carta de prueba cambian de casa 6 de 9 cuerpos y de signo solo 4.

Respeta el selector de nodo, para que Rāhu no salga en un sitio en una pestaña
y en otro en la de al lado.

## El switch de idioma

Cambiar de idioma rehace **todo** lo que hay montado, no solo lo que se está
mirando: las pestañas ocultas guardan su contenido y asomarían en el idioma
anterior al abrirlas. `repintarTodo()` vuelve a pintar los rótulos, el índice
del curso, el módulo que estuviera abierto, la lista de guardadas, el huso
resuelto, su panel de historia, la comparación, el ejercicio, la carta y la
lectura.

Y como parte del texto lo compone el servidor —los yogas, la lectura entera con
sus fuentes— no basta con repintar: hay que volver a pedirlo con el idioma
nuevo.

El ejercicio de cálculo a mano se repinta guardando antes lo que el alumno
lleve escrito y devolviéndoselo después. Perder sus cuentas por tocar el switch
sería una faena.

## Empaquetar

```
./empaquetar/todo.sh
```

Comprueba primero y, si algo falla, no publica nada. Después construye lo de
macOS, lo de Linux y los binarios sueltos de Windows y Linux.

Existe porque los venía compilando a mano cada vez, y lo que se hace a mano se
olvida: el paquete de Linux se quedó tres versiones atrás sin que nada avisara,
y encima llevaba dentro los binarios de macOS.

Los dos paquetes llevan un LEEME. El de macOS explica lo de la cuarentena de
Gatekeeper —la aplicación no está firmada, así que al bajarla el sistema se
niega a abrirla— con el `xattr -dr` y la alternativa sin terminal. El de Linux
trae un lanzador que elige amd64 o arm64 mirando `uname`.

## Dónde se guardan las cartas

En el directorio de configuración del sistema, dentro de `astro/`. Con la
variable `ASTRO_DIR` se manda a otro sitio — sirve para llevar el programa en
un pendrive con sus cartas dentro, y para que los tests no escriban en las del
usuario.

## Los tests

```
go test ./...
```

Lo que mide la astronomía no es el motor contra sí mismo —eso solo detecta
cambios, no errores— sino contra **Swiss Ephemeris**, que es la referencia con
la que se compara todo el mundo. Los valores están clavados en
`internal/efem/referencia_test.go`, generados por
`pruebas/generar_referencia.py`. Si un test falla, el que está mal es el motor.

Sobre 690 posiciones repartidas entre 1800 y 2100:

| | dentro de 1800-2050 | fuera, hasta 2100 |
|---|---|---|
| Sol | 0,37′ | 0,28′ |
| Luna | 0,74′ | 0,88′ |
| Saturno | 11,0′ | 18,6′ |
| resto de planetas | < 7,1′ | < 3,8′ |

Saturno es el peor porque la gran desigualdad con Júpiter no cabe en unos
elementos keplerianos. La tabla de Standish que usa el motor está dada para
1800-2050, y fuera se degrada; por eso el test mide las dos ventanas por
separado y no afloja el margen de la buena para tapar la mala.

Y por encima de todo eso, lo que de verdad importa en astrología:
**ninguna de las 690 posiciones cambia de signo.** Diez minutos de arco no
mueven a nadie de casa; salir en otro signo, sí.

Casas de Plácido contra Swiss Ephemeris: Ascendente 2,9″, Medio Cielo 2,2″,
cúspides 3,7″. Ayanāṁśa Lahiri: 0,0002″. Salida y puesta del Sol: 35 segundos.

Escribiendo estos tests aparecieron dos cosas: el día juliano estaba mal para
fechas anteriores al 15 de octubre de 1582 —se aplicaba la regla gregoriana de
los siglos siempre, y para el año 837 eran cuatro días de desfase—, y la
referencia de orto la estaba pidiendo con otro criterio de altura, así que los
números no eran comparables.

## Integración continua

`.github/workflows/comprobar.yml` ejecuta en cada push lo mismo que
`verificar.sh`: compila, `vet`, los tests con cobertura, la sintaxis del
JavaScript, el detector de castellano a mano, la paridad de idiomas, que el
curso esté traducido entero, que la interfaz se ejecute de verdad, y que el
binario compile para las cinco combinaciones de sistema y arquitectura.

Existe porque los tests no servían de mucho dependiendo de que yo me acordara
de correrlos.

## Comprobar antes de publicar

```
./verificar.sh
```

Compila, pasa `vet` y los tests, comprueba la sintaxis del JavaScript, que los
dos idiomas tengan las mismas claves, que los 35 módulos del curso existan en
inglés, que los endpoints respondan, y **ejecuta la interfaz entera** contra un
DOM de mentira con datos reales.

Lo último es lo que más falta hacía. El JavaScript viaja embebido con
`go:embed` y el compilador de Go no lo mira: un paréntesis de más ahí compila,
arranca, responde a todos los endpoints — y deja la interfaz muerta sin decir
nada. Pasó, y en dos versiones publicadas.

## Pendiente

- Ejercicio de cálculo a mano para jyotiṣa (el occidental ya está)
- Daśās distintas de la vimśottarī (aṣṭottarī, yoginī, cara)
- Praśna: la carta del momento de la pregunta
