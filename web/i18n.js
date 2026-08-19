// Cadenas de interfaz. El vocabulario depende de la TRADICIÓN, no solo del idioma:
// en modo védico no aparece la palabra "Ascendente" en ninguna parte, y viceversa.
const T = {
  es: {
    titulo: "Astro", lead_occidental: "Carta natal occidental, zodíaco tropical con casas de Plácido.",
    lead_jyotisha: "Carta védica, zodíaco sidéreo con ayanāṁśa Lahiri y casas de signo entero.",
    ciudad: "Ciudad", fecha: "Fecha", hora: "Hora", tz: "Desfase UTC", lat: "Latitud", lon: "Longitud",
    levantar: "Levantar", guardar: "Guardar", porque: "¿por qué este desfase?",
    elige: "Elige una ciudad y el desfase se resuelve solo.",
    nav: { carta: "Carta", lectura: "Lectura", vargas: "Vargas", dasas: "Daśās",
           comparar: "Comparar", curso: "Curso", ejercicio: "Cálculo a mano", imprimir: "Imprimir / PDF" },
    angulos: "Los cuatro ángulos", lagnaT: "Lagna",
    planetas: "Planetas", grahas: "Grahas", casas: "Casas y regentes", bhavas: "Bhāvas y sus señores",
    aspectos: "Aspectos, por exactitud", karakas: "Kārakas de Jaimini", yogas: "Yogas detectados",
    posicion: "Posición", signo: "Signo", grados: "Grados", casa: "Casa", estado: "Estado",
    nak: "Nakṣatra", pada: "Pada", senor: "Señor", alojado: "alojado en", ocupan: "Ocupan", aspectan: "Aspectan",
    dominante: "Planeta dominante", contradicciones: "Contradicciones — aquí empieza el oficio",
    compTit: "El mismo nacimiento en los dos zodíacos",
    compTxt: "Ésta es la única pantalla donde las dos tradiciones aparecen juntas, y está aquí a propósito: es la manera más rápida de entender por qué el mismo cielo da signos distintos.",
    tropical: "Tropical (occidental)", sidereo: "Sidéreo (jyotiṣa)", cambia: "cambia de signo",
    ejTit: "Corrección del cálculo a mano",
    ejTxt: "Levanta la carta con lápiz y efemérides. Escribe aquí lo que te haya salido y te digo en qué paso te desviaste — no el resultado bueno.",
    comprobar: "Comprobar", modulo: "Módulo", indice: "Índice",
    pie: "Motor astronómico propio en Go, sin dependencias. Verificado contra Swiss Ephemeris. " +
         "Lugares de GeoNames (CC BY 4.0). Husos de la base IANA.",
  },
  en: {
    titulo: "Astro", lead_occidental: "Western natal chart, tropical zodiac with Placidus houses.",
    lead_jyotisha: "Vedic chart, sidereal zodiac with Lahiri ayanāṁśa and whole-sign houses.",
    ciudad: "City", fecha: "Date", hora: "Time", tz: "UTC offset", lat: "Latitude", lon: "Longitude",
    levantar: "Cast chart", guardar: "Save", porque: "why this offset?",
    elige: "Pick a city and the offset resolves itself.",
    nav: { carta: "Chart", lectura: "Reading", vargas: "Vargas", dasas: "Daśās",
           comparar: "Compare", curso: "Course", ejercicio: "Manual calculation", imprimir: "Print / PDF" },
    angulos: "The four angles", lagnaT: "Lagna",
    planetas: "Planets", grahas: "Grahas", casas: "Houses and rulers", bhavas: "Bhāvas and their lords",
    aspectos: "Aspects, by exactness", karakas: "Jaimini kārakas", yogas: "Yogas found",
    posicion: "Position", signo: "Sign", grados: "Degrees", casa: "House", estado: "State",
    nak: "Nakṣatra", pada: "Pada", senor: "Lord", alojado: "placed in", ocupan: "Occupy", aspectan: "Aspect",
    dominante: "Dominant planet", contradicciones: "Contradictions — where the craft begins",
    compTit: "The same birth in both zodiacs",
    compTxt: "This is the only screen where both traditions appear together, and it is here on purpose: it is the fastest way to understand why the same sky yields different signs.",
    tropical: "Tropical (western)", sidereo: "Sidereal (jyotiṣa)", cambia: "changes sign",
    ejTit: "Checking your hand calculation",
    ejTxt: "Cast the chart with pencil and ephemeris. Type what you got and I will tell you which step went wrong — not the right answer.",
    comprobar: "Check", modulo: "Module", indice: "Index",
    pie: "Own astronomical engine in Go, no dependencies. Verified against Swiss Ephemeris. " +
         "Places from GeoNames (CC BY 4.0). Time zones from the IANA database.",
  }
};

// Pestañas por tradición: cada modo enseña solo lo suyo.
const NAV = {
  occidental: ["carta", "lectura", "comparar", "curso", "ejercicio"],
  jyotisha:   ["carta", "vargas", "dasas", "comparar", "curso"],
};

const CURSOS = {
  occidental: [["00-mapa","Mapa del curso",""],["01-el-cielo","El cielo desde un punto","1"],
    ["02-la-hora","La hora","2"],["03-el-calculo","El cálculo","3"],["04-angulos","Los cuatro ángulos","4"],
    ["05-posiciones","Posiciones","5"],["06-casas","Casas y regentes","6"],["07-aspectos","Aspectos","7"],
    ["08-dignidades","Dignidades","8"],["09-combinar","Combinar y contrarrestar","9"],
    ["10-sintesis","Síntesis","10"],["11-profundizar","Profundizar","11"],["12-prediccion","Predicción","12"],
    ["13-oficio","Oficio y límites","13"],["14-jyotisha","El zodíaco sidéreo","14"],
    ["corta","La vía corta en cinco bloques","·"],["plan-semana","Plan de una semana","·"],
    ["motor-profundizacion","Motor de profundización","·"]],
  jyotisha: [["00-mapa","Mapa del curso",""],["01-el-cielo","El cielo desde un punto","1"],
    ["02-la-hora","La hora","2"],["03-el-calculo","El cálculo","3"],["04-angulos","Los ángulos","4"],
    ["05-rasis","Los rāśis","5"],["06-nakshatras","Los nakṣatras","6"],["07-bhavas","Los bhāvas","7"],
    ["08-grahas","Los grahas","8"],["09-drishti","Dṛṣṭi: la mirada","9"],["10-dignidades","Dignidades y fuerza","10"],
    ["11-vargas","Las vargas","11"],["12-dasas","Las daśās","12"],["13-oficio","Oficio y límites","13"],
    ["14-karakas","Los kārakas","14"],["15-yogas","Los yogas","15"],["16-profundizar","Profundizar","16"]],
};
