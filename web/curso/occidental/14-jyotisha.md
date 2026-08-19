# Módulo 14 — El otro zodíaco: jyotiṣa

**Objetivo:** que entienda por qué el mismo nacimiento da signos distintos en India, y que
conozca el sistema que más lejos llega en profundizar. Es un módulo **informativo**: no se
aprende jyotiṣa aquí, se aprende que existe y qué ofrece.

## Por qué cambian los signos

El zodíaco occidental es **tropical**: empieza en el equinoccio de primavera, en el punto donde
el Sol cruza el ecuador celeste. Es una referencia estacional.

El zodíaco indio es **sidéreo**: empieza en las estrellas fijas.

Los dos coincidieron hace unos 1.700 años, pero el eje de la Tierra se bambolea —la precesión
de los equinoccios— y se han ido separando algo más de un grado cada setenta años. Hoy la
diferencia, el **ayanāṁśa**, ronda los **24 grados**.

Consecuencia práctica: casi todo el mundo tiene el Sol un signo antes en jyotiṣa. Quien es Libra
en occidental suele ser Virgo en jyotiṣa.

**No es que uno esté mal.** Miden cosas distintas: uno la relación con las estaciones, otro la
relación con las estrellas. Discutir cuál es "el verdadero" es discutir si el metro es más
verdadero que la milla.

Además, jyotiṣa usa casi siempre **casas de signo entero**: la casa 1 es todo el signo del
ascendente, la 2 el siguiente, y así. Sin cúspides intermedias.

`carta.py --sideral` calcula la versión sidérea de la misma carta, por si quiere verlo.

## Lo que jyotiṣa tiene y occidental no

Aquí está lo que importa de este módulo. Todo el módulo 11 —bajar de lo general a la causa— se
apoya en tener instrumentos para bajar. **Occidental tiene una escala de seis escalones.
Jyotiṣa tiene decenas**, y por eso llega mucho más hondo.

**Nakṣatras.** El zodíaco dividido no en 12 sino en **27 casas lunares** de 13°20′, cada una
con su regente planetario, su deidad y su carácter, y cada una subdividida en cuatro cuartos.
Un planeta no solo está en un signo: está en un nakṣatra, en un pada concreto, y bajo el
gobierno del regente de ese nakṣatra — cuya propia condición hay que mirar después. Es un
escalón entero de profundidad que occidental no tiene.

**Vargas.** Cartas divisionales: la carta se subdivide y produce **otra carta completa** para
cada asunto. La D9 para el matrimonio y el alma, la D10 para la profesión, la D12 para los
padres, la D7 para los hijos. Hay dieciséis en uso corriente. Cuando algo sale borroso en la
carta principal, se va a la varga de ese asunto y se vuelve a mirar todo desde allí.

**Daśās.** Sistemas de periodos que dicen **cuándo** se activa cada planeta, encajados unos
dentro de otros: un periodo mayor de años, dentro un sub-periodo, dentro un sub-sub-periodo.
Se puede bajar hasta niveles de días. Es la respuesta al "¿cuándo?" con una precisión que las
progresiones no alcanzan.

**Kārakas.** Significadores, y de dos clases: los fijos (el Sol es el padre siempre) y los
**variables**, que se asignan según el grado exacto que ocupe cada planeta en esa carta
concreta. Da un segundo juego de indicadores para leer los mismos temas.

**Cartas derivadas.** Tomar la casa de otra persona como su ascendente y leer su carta entera
desde dentro de la tuya.

## Las causas dentro de las causas

Con esos instrumentos, la lógica del módulo 11 se puede llevar mucho más lejos. Un ejemplo del
tipo de recorrido que permite:

> Hallazgo: *hay algo con la madre.*
> → El regente de la casa de la madre, y dónde está.
> → El nakṣatra de ese regente, y quién lo gobierna.
> → La condición de **ese** gobernante, que es otro planeta con su propia historia.
> → La varga de los padres, donde todo se recoloca.
> → La carta derivada de la madre, tomando su casa como ascendente suyo.
> → Y dentro de esa carta derivada, **su** casa de la madre — la abuela.

En ese punto aparece a menudo lo que se buscaba: que lo que parecía un asunto entre dos venía
de la generación anterior. **Una causa dentro de otra, provocada por una tercera.**

Se dice que jyotiṣa tiene cientos de miles de combinaciones calculables. Es literalmente cierto,
y es su fuerza y su peligro: con instrumentos infinitos, **las reglas de parada del módulo 11
dejan de ser una recomendación y pasan a ser lo único que separa el análisis de la invención.**

## Lo que se lleva de aquí

- Su Sol probablemente está un signo antes en jyotiṣa, y eso no invalida nada
- Occidental es más fuerte en psicología y en lectura de carácter
- Jyotiṣa es más fuerte en tiempo y en profundidad de causas
- **La disciplina es la misma en los dos**: contar convergencias, separar dato de inferencia, y
  saber parar

Si algún día quiere estudiarlo, que lo estudie como un sistema entero y no mezclando piezas.
Los sistemas mezclados no funcionan: cada uno tiene su lógica interna y sus reglas de peso.

## Preguntas

1. ¿Por qué el Sol cambia de signo entre los dos sistemas?
2. ¿Por qué no tiene sentido preguntar cuál es el verdadero?
3. ¿Qué es una varga y para qué sirve cuando una lectura sale borrosa?
4. Con cientos de miles de cálculos disponibles, ¿cuál es el riesgo mayor?

## Pregunta de cierre

**Si un sistema te permite bajar veinte capas, ¿qué es lo que impide que las últimas quince
sean inventadas?**
