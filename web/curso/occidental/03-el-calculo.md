# Módulo 3 — El cálculo

**Objetivo:** que levante una carta **a mano**, una vez en la vida. Después usará software,
pero sabiendo qué hace el software.

Este módulo es lento. No lo aceleres. Es el que separa a quien sabe astrología de quien sabe
manejar un programa.

## La idea

Calcular una carta es responder a una pregunta geométrica: **¿por qué grado de la eclíptica
pasaban el horizonte y el meridiano en ese instante y en ese punto de la Tierra?**

Todo lo demás son herramientas para contestarla.

## Los pasos

**1. De hora civil a UT.** Ya está en el módulo 2.

**2. Posiciones planetarias.** La efeméride da los planetas a mediodía (o a medianoche) de
cada día. El nacimiento casi nunca cae ahí, así que hay que **interpolar**: si el Sol avanza
57′ ese día y el nacimiento fue a un tercio del día, avanzó un tercio de 57′.

Aquí entran los **logaritmos de proporción** que traen las efemérides antiguas: convierten la
regla de tres en una suma. Explícale el mecanismo aunque hoy no haga falta.

**3. Tiempo sidéreo.** Este es el corazón y el punto donde todo el mundo se atasca.

El día solar dura 24 h; el sidéreo, unos 4 minutos menos, porque mientras la Tierra gira sobre
sí misma también avanza en su órbita. El **tiempo sidéreo** dice qué grado del zodíaco está
culminando ahora mismo en Greenwich.

Se toma de la efeméride para ese día y se corrige:
- por las horas transcurridas desde el momento tabulado
- por la aceleración (unos 10 segundos por hora transcurrida)
- por la **longitud** del lugar: cada 15° de longitud es una hora

El resultado es el **tiempo sidéreo local**.

**4. Ascendente y Medio Cielo.** Con el tiempo sidéreo local y la **latitud**, se entra en las
Tablas de Casas y se leen los dos ángulos. La latitud hace falta porque la inclinación del
horizonte respecto a la eclíptica cambia con ella: por eso en latitudes altas unos signos
tardan muchísimo en salir y otros salen en minutos.

**5. Cúspides intermedias.** Según el sistema. Es el módulo 6.

**6. Colocar los planetas** en las casas resultantes.

## Cómo lo corriges

Que haga **cada** cuenta él. Cuando dé un resultado, lanza:

```
python3 carta.py --fecha ... --hora ... --tz ... --lat ... --lon ... --pasos
```

Eso imprime UT, día juliano, tiempo sidéreo de Greenwich, corrección por longitud, tiempo
sidéreo local, Ascendente y MC.

**Dile en qué paso se desvió y cuánto. No le des el número bueno.** Que rehaga desde ahí.

## Preguntas

1. ¿Por qué el día sidéreo es más corto que el solar?
2. Dos personas nacen a la misma hora UT, una en Quito y otra en Reikiavik. ¿Por qué el
   Ascendente es distinto si el tiempo sidéreo de Greenwich es el mismo?
   *(Por la latitud y por la longitud: una entra en la tabla por otra fila y con otro TSL.)*
3. Si te equivocas en 4 minutos de tiempo sidéreo, ¿cuánto te equivocas en el Ascendente?
   *(Aproximadamente un grado, según latitud.)*
4. ¿Para qué sirve la latitud y para qué la longitud?
   *(Longitud: corregir el tiempo sidéreo. Latitud: entrar en las tablas de casas.)*

## Ejercicio

Levantar su carta entera a mano. Con lápiz. Tiempo estimado: entre una y dos horas la primera
vez. Que apunte cada paso para poder comparar.

## Pregunta de cierre

**¿Qué dato necesitas para el tiempo sidéreo local, y qué dato para el Ascendente?**

## Errores frecuentes

- Olvidar la corrección de aceleración
- Sumar la longitud cuando había que restarla
- Entrar en las tablas con la latitud equivocada de signo
- Interpolar la Luna como si fuera lineal en un día entero: se mueve 12-15° diarios y merece cuidado
