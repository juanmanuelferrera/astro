// Cadenas de interfaz. El vocabulario depende de la TRADICIÓN, no solo del idioma:
// en modo védico no aparece la palabra "Ascendente" en ninguna parte, y viceversa.
const SIGNOS = {
  es: ["Aries","Tauro","Géminis","Cáncer","Leo","Virgo","Libra","Escorpio","Sagitario","Capricornio","Acuario","Piscis"],
  en: ["Aries","Taurus","Gemini","Cancer","Leo","Virgo","Libra","Scorpio","Sagittarius","Capricorn","Aquarius","Pisces"],
};
// Los rāśis no se traducen: son los nombres que se usan en cualquier idioma.
const RASIS = ["Meṣa","Vṛṣabha","Mithuna","Karka","Siṁha","Kanyā","Tulā","Vṛścika","Dhanus","Makara","Kumbha","Mīna"];

// Nombres de los cuerpos. Rāhu y Ketu no se traducen: no tienen equivalente.
const CUERPOS = {
  es: {}, // el servidor ya los manda en español
  en: {"Sol":"Sun","Luna":"Moon","Mercurio":"Mercury","Venus":"Venus","Marte":"Mars",
       "Júpiter":"Jupiter","Saturno":"Saturn","Urano":"Uranus","Neptuno":"Neptune",
       "Plutón":"Pluto","Nodo Norte":"North Node","Nodo Sur":"South Node",
       "Rāhu":"Rāhu","Ketu":"Ketu","Lagna":"Lagna"},
};

const T = {
  es: {
    titulo: "Astro", lead_occidental: "Carta natal occidental, zodíaco tropical con casas de Plácido.",
    lead_jyotisha: "Carta védica, zodíaco sidéreo con ayanāṁśa Lahiri y casas de signo entero.",
    ciudad: "Ciudad", fecha: "Fecha", hora: "Hora", tz: "Desfase UTC", lat: "Latitud", lon: "Longitud",
    levantar: "Levantar", guardar: "Guardar", porque: "¿por qué este desfase?",
    elige: "Elige una ciudad y el desfase se resuelve solo.",
    nav: { carta: "Carta", lectura: "Lectura", vargas: "Vargas", dasas: "Daśās y tránsitos",
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
    trad_occidental: "Occidental", trad_jyotisha: "Jyotiṣa",
    norte: "Norte de la India", sur: "Sur de la India",
    norteNota: "casa 1 arriba", surNota: "signos fijos",
    guardadas: "Guardadas", nombrePrompt: "Nombre", borrar: "¿Borrar «{n}»?",
    ang_asc: "Ascendente", ang_dsc: "Descendente", ang_mc: "Medio Cielo", ang_ic: "Fondo del Cielo",
    rasi: "Rāśi", extra: "Extra", vimsottari: "Vimśottarī",
    sadeTit: "Sade Sati", sadeAct: "Activo — fase {f}.", sadeSale: " Sale hacia {d}.",
    sadeProx: " El próximo empieza hacia {d}.",
    transitos: "Tránsitos de hoy · {f}", desdeLagna: "desde Lagna", desdeLuna: "desde Luna",
    graha: "Graha", verano: "horario de verano", estandar: "horario estándar",
    hDe: "De dónde sale UTC{o}", hZona: "zona", hEstandar: "estándar del año",
    hVerano: "horario de verano", hSi: "sí, +1 h", hNo: "no",
    hRelojSol: "El reloj frente al Sol", hTocaria: "por longitud le tocaría",
    hVa: "el reloj va {d} h respecto al Sol", hDstAnio: "Horario de verano en {a}",
    hSinDst: "ese año no hubo cambios", hCambios: "Cambios de huso del país",
    ejOk: "Los cinco pasos correctos.", ejRehaz: "Rehaz desde ahí.",
    v_D1:"la vida tal como se vive", v_D2:"riqueza y sustento", v_D3:"hermanos y coraje",
    v_D7:"hijos", v_D9:"el alma y el cónyuge", v_D10:"profesión", v_D12:"los padres",
    v_D16:"vehículos y confort", v_D30:"males", v_D60:"karma acumulado",
    cursoSoloEs: "Los módulos del curso están por ahora solo en español.",
    pie: "Motor astronómico propio en Go, sin dependencias. Verificado contra Swiss Ephemeris. " +
         "Lugares de GeoNames (CC BY 4.0). Husos de la base IANA.",
  },
  en: {
    titulo: "Astro", lead_occidental: "Western natal chart, tropical zodiac with Placidus houses.",
    lead_jyotisha: "Vedic chart, sidereal zodiac with Lahiri ayanāṁśa and whole-sign houses.",
    ciudad: "City", fecha: "Date", hora: "Time", tz: "UTC offset", lat: "Latitude", lon: "Longitude",
    levantar: "Cast chart", guardar: "Save", porque: "why this offset?",
    elige: "Pick a city and the offset resolves itself.",
    nav: { carta: "Chart", lectura: "Reading", vargas: "Vargas", dasas: "Daśās and transits",
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
    trad_occidental: "Western", trad_jyotisha: "Jyotiṣa",
    norte: "North Indian", sur: "South Indian",
    norteNota: "house 1 on top", surNota: "fixed signs",
    guardadas: "Saved", nombrePrompt: "Name", borrar: "Delete “{n}”?",
    ang_asc: "Ascendant", ang_dsc: "Descendant", ang_mc: "Midheaven", ang_ic: "Imum Coeli",
    rasi: "Rāśi", extra: "Extra", vimsottari: "Vimśottarī",
    sadeTit: "Sade Sati", sadeAct: "Active — {f} phase.", sadeSale: " Ends around {d}.",
    sadeProx: " The next one begins around {d}.",
    transitos: "Transits today · {f}", desdeLagna: "from Lagna", desdeLuna: "from Moon",
    graha: "Graha", verano: "daylight saving", estandar: "standard time",
    hDe: "Where UTC{o} comes from", hZona: "zone", hEstandar: "standard for that year",
    hVerano: "daylight saving", hSi: "yes, +1 h", hNo: "no",
    hRelojSol: "The clock against the Sun", hTocaria: "longitude would give",
    hVa: "the clock runs {d} h against the local Sun", hDstAnio: "Daylight saving in {a}",
    hSinDst: "no clock changes that year", hCambios: "Country-wide offset changes",
    ejOk: "All five steps correct.", ejRehaz: "Redo from there.",
    v_D1:"life as it is lived", v_D2:"wealth and sustenance", v_D3:"siblings and courage",
    v_D7:"children", v_D9:"the soul and the spouse", v_D10:"profession", v_D12:"the parents",
    v_D16:"vehicles and comfort", v_D30:"misfortunes", v_D60:"accumulated karma",
    cursoSoloEs: "The course modules are in Spanish for now.",
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
