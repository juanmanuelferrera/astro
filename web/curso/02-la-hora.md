# Módulo 2 — La hora

**Objetivo:** que entienda que el dato más frágil de una carta es la hora, y que sepa
convertir una hora de reloj en hora universal sin equivocarse.

## La idea

El módulo 1 dejó una consecuencia incómoda: si el horizonte se mueve un grado cada cuatro
minutos, **un error de veinte minutos mueve el Ascendente cinco grados**. Un error de dos
horas lo cambia de signo. Casi todo lo que sale mal en una carta sale mal aquí.

Hay tres horas distintas y conviene no mezclarlas.

**Hora del reloj (hora civil).** Lo que decía el reloj de la pared. Es un acuerdo político:
depende del huso que ese país tuviera ese día.

**Hora universal (UT).** La de Greenwich. Es la que usan las efemérides. A ella hay que llegar.

**Hora local verdadera.** La que marcaría el Sol en ese punto exacto. Casi nunca coincide con
el reloj, porque los husos son franjas anchas.

## Las cuatro trampas

**1. El horario de verano.** La peor de todas. Cambia según el país y según el año, y las reglas
se han modificado muchas veces. España estuvo en horario de Berlín desde 1940. Muchos países
lo aplicaron en años sueltos durante guerras y crisis del petróleo. **Nunca lo supongas:
compruébalo para ese país y ese año.**

**2. El huso no es la longitud.** España está casi entera al oeste de Greenwich, pero usa la
hora de Europa central. Un nacimiento en Galicia lleva más de dos horas de desfase entre el
reloj y el Sol.

**3. La hora redondeada.** "Nació sobre las dos" no es un dato: es un intervalo de media hora,
y media hora son unos siete grados de Ascendente. Hay que preguntar de dónde sale la hora:
del parte médico, del registro civil, o de la memoria de alguien.

**4. Medianoche y mediodía.** Las doce de la noche del día 3 son las 00:00 del 3, no del 4.
Es un error tonto que se comete constantemente.

## Cómo se valida una hora

Preguntar «¿a qué hora naciste?» no basta. Casi todo el mundo contesta con una hora que nunca
ha comprobado, heredada de una conversación familiar. Validarla es un trabajo de archivo, y se
hace antes de calcular nada.

**Las fuentes, de la mejor a la peor:**

1. **Historia clínica del hospital.** Registra el minuto. Es lo más fiable que existe. En España
   cualquiera puede pedir su propia historia clínica al hospital donde nació, aunque los
   archivos antiguos a veces ya no están.
2. **Certificación literal de nacimiento** del Registro Civil. Recoge la hora que declaró el
   hospital. En España se pide gratis por la sede electrónica del Ministerio de Justicia. Es la
   fuente práctica: documental, accesible y con hora.
3. **Libro de familia.** A veces la trae, a veces no.
4. **Algo escrito de la época** — una agenda, una carta, un telegrama. Vale más que un recuerdo
   porque se escribió en caliente.
5. **La memoria de alguien.** La peor fuente, y la que se usa el noventa por ciento de las veces.

**Señales de que una hora está inventada:**

- **Termina en 00 o en 30.** Si al preguntar por varias personas salen muchas horas redondas,
  esas horas son estimaciones. Los nacimientos no se agrupan en punto y media.
- **Es la hora del registro, no la del parto.** Se anota cuando alguien rellena el papel, que
  puede ser horas después. Ojo con las horas de oficina.
- **Dos familiares dan horas distintas.** Pregunta a cada uno por separado y sin sugerirle
  nada — decir «¿fue sobre las tres?» contamina la respuesta.

Al revés, una **cesárea programada** suele traer hora muy fiable: estaba en el parte quirúrgico.

**Y si no se puede validar**, la respuesta correcta no es elegir una hora bonita. Es trabajar
con un **intervalo**: «entre las 9 y las 11». Después se mira qué cambia dentro de ese intervalo
— si no cambia nada relevante, se sigue adelante; si cambia el Ascendente de signo, hay que
rectificar antes de leer nada de casas. Eso es el módulo 4.

Y hay que decirlo en voz alta al consultante: **con hora dudosa, los planetas en signos siguen
valiendo y las casas no.**

## El procedimiento

1. Hora del reloj y fecha, tal como se registraron
2. ¿Había horario de verano ese día en ese país? Si sí, resta una hora
3. Aplica el huso estándar del lugar → obtienes **UT**
4. Comprueba: ¿el resultado cae en la fecha correcta? A veces se cambia de día

## Preguntas

1. Un niño nace en Barcelona el 15 de julio de 1985 a las 03:20. ¿Qué UT es?
   *(Verano → UTC+2. UT = 01:20 del mismo día.)*

2. Uno nace en Madrid a las 00:30 del 1 de enero. ¿Qué fecha se usa para la efeméride?
   *(UTC+1 en invierno → UT = 23:30 del 31 de diciembre. Cambia el día y el año.)*

3. La madre dice "fue a media mañana". ¿Qué margen de Ascendente estás manejando?
   *(De 09:00 a 12:00 son tres horas: unos 45°, es decir, hasta signo y medio. La carta no
   se puede levantar con eso. Hay que rectificar o conseguir el dato.)*

4. ¿Por qué la longitud del lugar no basta para saber el huso?
   *(Porque el huso es una decisión administrativa, no astronómica.)*

## Ejercicio

Que reconstruya la conversión de **su propia** hora, paso a paso, y compare con lo que imprime
`carta.py --pasos` en las líneas 1 y 2. Si no coincide, que encuentre en cuál de las cuatro
trampas cayó.

Y el trabajo de archivo, que se hace fuera de la sesión: **rastrear su hora hasta la fuente más
alta que pueda alcanzar** de la lista de arriba. Que vuelva diciendo de qué documento sale, no
de quién se lo contó. Si acaba en «me lo dijo mi madre», que se lo pregunte también a otra
persona por separado y compare.

## Pregunta de cierre

**¿Por qué una carta con hora aproximada sigue sirviendo para algunas cosas y no para otras?**

Debe llegar a: los planetas en signos siguen siendo válidos; los ángulos y las casas, no. O sea
que se puede hablar de temperamento y no de las áreas de la vida.

## Errores frecuentes

- Aplicar el horario de verano de hoy a un nacimiento de hace cincuenta años
- Sumar cuando había que restar. Al oeste de Greenwich se suma para llegar a UT; al este se resta
- Dar por buena una hora "de familia" sin rastrearla hasta un documento
- Sugerirle la hora al familiar al preguntarle, y quedarse con lo que contesta
- Confundir la hora del parto con la hora en que se inscribió en el registro
