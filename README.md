# Astro

Cartas natales y curso en **dos tradiciones** —occidental y jyotiṣa— y **dos idiomas**, en un
único binario de Go. Sin Python, sin red, sin instalar nada: las efemérides, los dos cursos y
41.451 poblaciones van dentro del ejecutable.

```
./astro
```

Abre el navegador en `http://localhost:8733`. Si ese puerto está ocupado, busca el siguiente.

Versión actual: **1.8.0** (`./astro -version`).

---

## Descargar

En [releases](https://github.com/juanmanuelferrera/astro/releases) hay un paquete por sistema:

| | |
|---|---|
| `Astro-mac-app.zip` | aplicación de macOS con icono, sin terminal. Trae un LEEME con lo de la cuarentena de Gatekeeper |
| `astro-linux.zip` | los dos binarios y un lanzador que elige amd64 o arm64 mirando `uname` |
| `astro-mac` | binario universal suelto (x86_64 + arm64) |
| `astro-linux-amd64`, `astro-linux-arm64`, `astro-windows-amd64.exe` | binarios sueltos |

La aplicación de macOS **no está firmada** con cuenta de Apple, así que al bajarla el sistema se
niega a abrirla. Se quita así, y el `-r` hace falta porque es un paquete y no un fichero:

```
xattr -dr com.apple.quarantine /Applications/Astro.app
```

O sin terminal: Control-clic sobre la aplicación, «Abrir», confirmar. Solo la primera vez.

---

## Las dos tradiciones

Un selector arriba cambia **todo a la vez**: el cálculo, el vocabulario, el dibujo y el curso.
En modo védico la palabra «Ascendente» no aparece en ninguna parte; en occidental no aparece
«Lagna». No se pueden mezclar por accidente.

La excepción es deliberada: la pestaña **Comparar** pone los dos sistemas lado a lado. El mayor
riesgo de confusión convertido en la función que mejor enseña.

### Occidental

Zodíaco tropical · casas de Plácido y casas iguales en paralelo · aspectos por grados con orbes ·
rueda circular con desapilado angular de glifos y línea al grado real · regentes de casa.

Pestañas: **Carta · Lectura · Predicción · Comparar · Curso · Cálculo a mano**

### Jyotiṣa

Zodíaco sidéreo con ayanāṁśa Lahiri · bhāvas de signo entero · 27 nakṣatras con pada y señor ·
dṛṣṭi por signo con los aspectos especiales de Marte, Júpiter y Saturno · dignidades védicas con
grados exactos, mūlatrikoṇa, combustión, dig-bala y gaṇḍānta · diez vargas dibujadas · yogas
(pañca-mahāpuruṣa, rāja-yoga, nīca-bhaṅga, gaja-kesari, kemadruma) · tránsitos y Sade Sati con la
fase actual y su fecha de salida.

Carta cuadrada en los **dos estilos**, norte y sur de la India, con las vargas dibujadas por el
mismo componente para poder compararlas.

Pestañas: **Carta · Lectura · Vargas · Daśās · Pañcāṅga · Fuerza · Praśna · Comparar · Curso ·
Cálculo a mano**

### Comunes a las dos

41.451 lugares (España completa, resto del mundo desde 15.000 habitantes) · husos horarios
históricos con el botón «¿por qué este desfase?» · guardar y recuperar cartas · imprimir a PDF ·
interfaz en español e inglés · **curso de 35 módulos** (18 occidentales, 17 védicos, 50.000
palabras) traducido entero.

---

## Qué hace cada pestaña

### Lectura

Las dos tradiciones tienen motor de interpretación, y los dos **componen la frase a partir de sus
piezas** en lugar de guardar textos prefabricados.

**Occidental** — `internal/occidental/claves.go`. Sistema de Palabras Clave de Margaret Hone:
función del planeta + modo del signo + terreno de la casa + dignidad, más los regentes, que
convierten un aspecto suelto en un argumento. Detecta cuándo un planeta recibe a la vez aspectos
duros y blandos y lo declara como contradicción en lugar de elegir un lado.

**Jyotiṣa** — `internal/jyotisha/lectura.go`. Función del graha + modo del rāśi + terreno del
bhāva + nakṣatra + dignidad, y sobre todo la **cadena de señores**, que es como razona jyotiṣa: de
qué depende un asunto, y de qué depende aquello. Añade dṛṣṭi, la daśā que corre ahora y los
kārakas de Jaimini. Marca como disputado el bhāva que reciben benéficos y maléficos a la vez.

Ninguno de los dos sintetiza. Los dos dicen por qué.

### Predicción (occidental)

Tránsitos, progresiones secundarias y revolución solar — las tres que enseña el módulo 12.

Los tránsitos, solo de los lentos: los rápidos marcan días y no periodos. Orbes separados, 1,5°
los aspectos mayores y 0,4° los menores; con el mismo orbe la lista se llenaba de
sesquicuadraturas y el tránsito que importaba quedaba enterrado. Dice si aplica o separa, y si el
planeta va a retrogradar encima del punto y **pasar tres veces**.

Y lo que las une: las **convergencias**. Un periodo importa cuando dos técnicas **distintas**
señalan el mismo punto natal. Dos progresiones sobre el mismo sitio son una voz repetida y no
cuentan. Cuando no coincide nada, lo dice en lugar de rellenar.

### Pañcāṅga (jyotiṣa)

Los cinco miembros del calendario hindú: tithi con su quincena, vāra, nakṣatra con su pada, yoga y
karaṇa, cada uno con lo que lleva recorrido. El tithi y el karaṇa salen de la **diferencia** entre
la Luna y el Sol, así que el ayanāṁśa se cancela; el yoga sale de la **suma**, y no se cancela.

Los tres lagnas especiales —Bhāva, Horā y Ghaṭī— se cuentan desde el amanecer, así que hacen falta
salida y puesta del Sol: están en `internal/efem/orto.go`. Y los arudha padas de los doce bhāvas,
con el Arudha Lagna y el Upapada.

### Fuerza (jyotiṣa)

**Aṣṭakavarga** — las ocho tablas de Parāśara, con el bhinnāṣṭakavarga por graha y el
sarvāṣṭakavarga sumado. Las tablas se comprueban solas: los siete BAV suman exactamente 337, y hay
un test que falla si al teclearlas se coló un número de más.

**Ṣaḍbala** — las seis fuerzas en virūpas: sthāna (ucca, saptavargaja sobre siete vargas con
amistad compuesta, ojayugma, kendrādi, drekkāṇa), dig, kāla (nathonnatha, pakṣa, tribhāga, señores
del día y de la hora, ayana), cheṣṭā por los ocho estados del movimiento, naisargika, dṛk y
**yuddha** — la guerra entre dos planetas a menos de un grado, que sale en 8 de cada 100 cartas.

Lo que se muestra no es la cifra bruta sino la **razón** entre lo que cada graha saca y lo que se
le exige, porque el listón es distinto para cada uno.

Aparte va el **bhāva bala**, que contesta otra pregunta: no cuánto puede un graha sino cuánto
puede un asunto.

### Daśās (jyotiṣa)

Vimśottarī con sus bhuktis, más tres sistemas para contrastarla — cuando dos señalan el mismo
periodo la lectura se sostiene, y cuando solo lo dice uno es un murmullo.

- **Aṣṭottarī**, 108 años entre ocho grahas. Sin Ketu: eso es lo que la distingue.
- **Yoginī**, 36 años entre ocho yoginīs. Es corta, así que una vida la recorre tres veces.
- **Cara**, la de Jaimini — no cuelga de la Luna sino de los **rāśis**, y la duración de cada
  periodo sale de contar del signo al sitio donde está su señor.

### Praśna (jyotiṣa)

La carta del instante en que se hace la pregunta, no la del nacimiento. Es el recurso clásico
cuando no hay hora fiable — que, como dice el módulo 2, es casi siempre.

Lo primero que hace **no es contestar: es decidir si la pregunta se puede contestar**. Lagna en
los tres primeros grados de su rāśi, en los tres últimos, en gaṇḍānta, o la Luna en gaṇḍānta —
cualquiera de ésas y la respuesta es «espera y vuelve a preguntar», que es una respuesta completa.

Después juzga el bhāva del asunto (diecisiete a elegir): su señor y dónde está, quién lo ocupa, si
Júpiter lo mira, si el señor del lagna y el del asunto se relacionan, y los bindus del
aṣṭakavarga de ese rāśi.

### Comparar

La única pantalla donde las dos tradiciones aparecen juntas, y enseña las **dos** diferencias:

- El **zodíaco**: tropical contra sidéreo, con el ayanāṁśa que los separa.
- Las **casas**: Plácido, desiguales y dependientes de la latitud, contra el signo entero. Ésta
  suele pesar más y casi nunca se cuenta. En la carta de prueba cambian de casa 6 de 9 cuerpos y
  de signo solo 4.

### Cálculo a mano

Escribes lo que te salió en cada paso y dice en cuál te desviaste, **sin darte el resultado
bueno**. Los tres primeros pasos son los mismos en las dos tradiciones, porque la astronomía no
cambia; jyotiṣa añade restar el ayanāṁśa y sacar el Lagna sidéreo.

El corrector reconoce el error de **restar el ayanāṁśa dos veces**, que es el que más se comete al
empezar, y lo dice con esas palabras.

---

## Precisión

Motor astronómico propio, **cero dependencias externas**. Meeus para Sol, Luna, tiempo sidéreo y
orto; elementos de Standish con corrección de precesión al equinoccio de la fecha para los
planetas; casas de Plácido por bisección numérica; ayanāṁśa Lahiri por polinomio ajustado.

Los tests lo miden contra **Swiss Ephemeris**, no contra sí mismo — eso solo detectaría cambios,
no errores. Sobre 690 posiciones entre 1800 y 2100:

| | 1800-2050 | hasta 2100 |
|---|---|---|
| Sol | 0,37′ | 0,28′ |
| Luna | 0,74′ | 0,88′ |
| Júpiter | 7,1′ | 3,4′ |
| Saturno | 11,0′ | 18,6′ |
| resto de planetas | < 3,2′ | < 3,8′ |
| casas de Plácido | 3,7″ | |
| ayanāṁśa Lahiri | 0,0002″ | |
| salida y puesta del Sol | 35 s | |

Saturno es el peor porque la gran desigualdad con Júpiter no cabe en unos elementos keplerianos.
La tabla de Standish está dada para 1800-2050 y fuera se degrada; por eso el test mide las dos
ventanas por separado, en lugar de aflojar el margen de la buena para tapar la mala.

Y por encima de todo eso, lo que de verdad importa en astrología: **ninguna de las 690 posiciones
cambia de signo**. Diez minutos de arco no mueven a nadie de casa; salir en otro signo, sí.

En latitudes polares detecta que Plácido no puede resolverse y cae a casas iguales avisándolo, y
que el Sol no sale, en lugar de inventarse una hora.

---

## Opciones

| | |
|---|---|
| `-puerto=9000` | puerto fijo; por defecto 8733, y salta al siguiente libre si está ocupado |
| `-abrir=false` | no abrir el navegador al arrancar |
| `-red` | aceptar conexiones de otros equipos de la red local |
| `-version` | decir la versión y salir |

Por defecto se ata solo a `127.0.0.1`: nadie de la red puede entrar.

Las cartas guardadas van al directorio de configuración del sistema, dentro de `astro/`. Con la
variable **`ASTRO_DIR`** se mandan a otro sitio — sirve para llevar el programa en un pendrive con
sus cartas dentro, y para que los tests no escriban en las del usuario.

---

## Cómo está hecho

```
main.go                 servidor y endpoints; aquí se declara la versión
internal/efem/          astronomía compartida: sol, luna, planetas, casas, tiempo, orto
internal/occidental/    palabras clave de Hone y predicción
internal/jyotisha/      sidéreo, nakṣatras, vargas, daśās, aṣṭakavarga, ṣaḍbala,
                        pañcāṅga, arudhas, praśna, yogas, tránsitos
internal/lugares/       41.451 poblaciones y husos históricos
internal/guardadas/     persistencia de cartas
web/                    interfaz y los dos cursos
pruebas/                comprobaciones de la interfaz, que Go no puede hacer
empaquetar/             los paquetes de macOS y Linux
```

Compilar:

```
go build -o astro .
```

**Datos de terceros:** lugares de [GeoNames](https://www.geonames.org) (CC BY 4.0); husos de la
base IANA, embebida vía `time/tzdata`.

---

## Comprobaciones

```
./verificar.sh
```

Compila, pasa `vet`, corre los 47 tests con cobertura, y luego lo que Go no alcanza:

| | |
|---|---|
| sintaxis del JavaScript | `node --check` |
| `pruebas/estructura.mjs` | las costuras de la página: elementos que el JS pide contra los que el HTML tiene, identificadores repetidos, pestañas sin sección, variables de color sin definir, etiquetas sin cerrar, y que la versión esté declarada una sola vez |
| `pruebas/castellano.mjs` | frases en castellano escritas a mano dentro del código |
| paridad de idiomas | que las dos tablas tengan las mismas claves y todas las pestañas nombre |
| el curso | que los 35 módulos existan en los dos idiomas |
| `pruebas/interfaz.mjs` | **ejecuta la interfaz** contra un DOM de mentira con datos reales del servidor, en las dos tradiciones y los dos idiomas |

Cobertura: occidental 97,7% · lugares 92,9% · efem 93,1% · jyotiṣa 91,5% · guardadas 84,1%.

`.github/workflows/comprobar.yml` lo ejecuta todo en cada push, más la compilación cruzada para
las cinco combinaciones de sistema y arquitectura.

### Por qué existe tanta comprobación de la interfaz

El JavaScript viaja embebido con `go:embed`, así que **el compilador de Go no lo mira**. Estos tres
fallos compilaban, arrancaban, respondían a todos los endpoints y pasaban los tests:

- Un literal de plantilla mal cerrado. Un error de sintaxis tumba el fichero entero: la interfaz
  no pintaba nada. Salió en dos versiones publicadas.
- `var(--acento)`, `var(--linea)` y `var(--texto)`, que no existen. Una variable CSS sin definir
  no da error: la declaración se vuelve inválida y la propiedad cae a su valor inicial, que para
  un `background` es transparente. Las barras se dibujaban con su ancho exacto y en transparente.
- Los botones de estilo védico salían en occidental. El JavaScript los escondía con `hidden`, y
  `hidden` no hacía nada: la regla `[hidden]{display:none}` del navegador tiene la especificidad
  de una clase, y `.estilos{display:flex}` la anulaba. La página ya declara
  `[hidden]{display:none!important}`, que cubre a todos los elementos y a los que vengan.

Ninguno lo habría cazado un test de datos. Cada uno dejó una comprobación detrás.

---

## Empaquetar

```
./empaquetar/todo.sh
```

Comprueba primero y, si algo falla, **no construye nada**. Después arma el paquete de macOS, el de
Linux y los binarios sueltos, y dice el comando exacto para publicar la versión que acaba de
construir.

Existe porque los venía compilando a mano cada vez, y lo que se hace a mano se olvida: el paquete
de Linux se quedó tres versiones atrás sin que nada avisara, y encima llevaba dentro los binarios
de macOS.

---

## Pendiente

Nada declarado.

Lo único que el proyecto no puede comprobarse a sí mismo es **el aspecto**. Todo se verifica sin
navegador, y eso caza lógica, traducciones y estructura, pero no si algo se ve torcido.
