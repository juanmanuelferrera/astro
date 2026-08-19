package occidental

import (
	"fmt"

	"astro/internal/efem"
)

// Sistema de Palabras Clave (Margaret Hone). La app COMPONE la frase a partir de
// los componentes; no guarda interpretaciones prefabricadas. Es el paso 1 de los
// cuatro que enseña el módulo 10: traducción literal y deliberadamente torpe.
//
// Todo sale en el idioma que se pida. El orden de las palabras cambia entre el
// español y el inglés, así que cada idioma trae sus propias plantillas y no se
// arma la frase pegando trozos traducidos sueltos.

type tabla struct {
	planeta map[string]string
	signo   [12]string
	casa    [12]string
	aspecto map[string]string
	dign    map[string]string
	cat     [12]string
	nombres map[string]string
	f       map[string]string
}

var es = tabla{
	planeta: map[string]string{
		"Sol": "tu vitalidad", "Luna": "tu manera de reaccionar",
		"Mercurio": "tu forma de pensar y de hablar", "Venus": "tu manera de valorar y de unirte",
		"Marte": "tu empuje", "Júpiter": "tu manera de expandirte y confiar",
		"Saturno": "tu manera de poner límites", "Urano": "tu manera de romper y despertar",
		"Neptuno": "tu manera de disolver e imaginar", "Plutón": "tu manera de vaciar y regenerar",
		"Nodo Norte": "hacia dónde tiras", "Nodo Sur": "de dónde vienes",
	},
	signo: [12]string{"con arranque y sin rodeos", "con paciencia y apego a lo tangible",
		"con curiosidad y muchas versiones", "con sensibilidad y ganas de proteger",
		"con brillo y ganas de crear", "con detalle y sentido crítico",
		"buscando acuerdo y equilibrio", "con intensidad y sin soltar",
		"con amplitud y ganas de sentido", "con prudencia y mirada larga",
		"con distancia y cabeza fría", "con porosidad y sin bordes claros"},
	casa: [12]string{"tu cuerpo y tu manera de asomarte al mundo",
		"lo que posees y lo que valoras", "tu entorno cercano y cómo te comunicas",
		"tu casa y tus raíces", "lo que creas, disfrutas y amas",
		"tu trabajo diario y tu salud", "tus vínculos íntimos y quien te hace de espejo",
		"lo que compartes y las crisis que transforman", "lo que estudias y aquello en lo que crees",
		"tu vocación y lo que se ve de ti desde fuera", "tus amistades y tus proyectos",
		"tu retiro, lo que sueltas y lo que no se ve"},
	aspecto: map[string]string{
		"conjunción": "se funde con", "oposición": "tira en sentido contrario a",
		"trígono": "fluye con facilidad hacia", "cuadratura": "roza y cuesta con",
		"sextil": "colabora, si le haces caso, con", "semicuadratura": "roza levemente con",
		"sesquicuadratura": "roza levemente con", "quincuncio": "no acaba de encajar con",
		"semisextil": "se relaciona con cierta incomodidad con",
	},
	dign: map[string]string{
		"domicilio":  "y aquí funciona a pleno rendimiento",
		"exaltación": "y aquí va sobrado, quizá de más",
		"exilio":     "y aquí trabaja a contrapelo, lo que con los años da oficio",
		"caída":      "y aquí se siente incómodo y poco valorado",
	},
	cat: [12]string{"Carácter", "Trabajo y dinero", "Mente", "Familia y raíces",
		"Vínculos", "Salud", "Vínculos", "Vínculos", "Mente", "Trabajo y dinero",
		"Vínculos", "Retiro e interior"},
	f: map[string]string{
		"pos":    "%s se expresa %s, en el terreno de %s",
		"asp":    "%s %s %s.",
		"reg":    "Lo relativo a %s pasa por %s: el regente de la casa %d está alojado en la %d.",
		"contra": "%s recibe a la vez %s de %s (orbe %.2f°) y %s de %s (orbe %.2f°). Por exactitud manda %s, pero las dos cosas están ahí: descríbelas juntas con «y», no con «pero».",
		"tension": "la tensión", "facilidad": "la facilidad",
		"dom":    "%s (regente del Ascendente: %s)",
		"nota":   "Esto es el paso 1 de 4: traducción literal, deliberadamente torpe. Agrupar por categorías es el paso 2 y ya está hecho. Los pasos 3 y 4 —resolver las contradicciones y escribirlo como texto seguido— los haces tú. La carta no se cuadra sola.",
	},
}

var en = tabla{
	planeta: map[string]string{
		"Sol": "your vitality", "Luna": "the way you react",
		"Mercurio": "the way you think and speak", "Venus": "the way you value and join",
		"Marte": "your drive", "Júpiter": "the way you expand and trust",
		"Saturno": "the way you set limits", "Urano": "the way you break and wake up",
		"Neptuno": "the way you dissolve and imagine", "Plutón": "the way you empty out and regenerate",
		"Nodo Norte": "where you are pulling", "Nodo Sur": "where you come from",
	},
	signo: [12]string{"headlong and without detours", "patiently, attached to what is solid",
		"curiously, in many versions", "sensitively, wanting to protect",
		"brightly, wanting to create", "in detail, with a critical eye",
		"seeking agreement and balance", "intensely, without letting go",
		"broadly, wanting it to mean something", "cautiously, taking the long view",
		"at a distance, with a cool head", "porously, without clear edges"},
	casa: [12]string{"your body and the way you show up in the world",
		"what you own and what you value", "your immediate surroundings and how you communicate",
		"your home and your roots", "what you create, enjoy and love",
		"your daily work and your health", "your close ties and whoever mirrors you",
		"what you share and the crises that transform you", "what you study and what you believe",
		"your calling and what others see of you", "your friendships and your projects",
		"your retreat, what you let go of and what stays out of sight"},
	aspecto: map[string]string{
		"conjunción": "merges with", "oposición": "pulls the opposite way from",
		"trígono": "flows easily towards", "cuadratura": "grates and comes hard with",
		"sextil": "works with, if you let it,", "semicuadratura": "grates slightly with",
		"sesquicuadratura": "grates slightly with", "quincuncio": "never quite fits with",
		"semisextil": "relates, somewhat uneasily, to",
	},
	dign: map[string]string{
		"domicilio":  "and here it works at full strength",
		"exaltación": "and here it runs high, perhaps too high",
		"exilio":     "and here it works against the grain, which over the years makes a craftsman",
		"caída":      "and here it feels awkward and little valued",
	},
	cat: [12]string{"Character", "Work and money", "Mind", "Family and roots",
		"Ties", "Health", "Ties", "Ties", "Mind", "Work and money",
		"Ties", "Retreat and inner life"},
	nombres: map[string]string{"Sol": "Sun", "Luna": "Moon", "Mercurio": "Mercury",
		"Venus": "Venus", "Marte": "Mars", "Júpiter": "Jupiter", "Saturno": "Saturn",
		"Urano": "Uranus", "Neptuno": "Neptune", "Plutón": "Pluto",
		"conjunción": "conjunction", "oposición": "opposition", "trígono": "trine",
		"cuadratura": "square", "sextil": "sextile", "semicuadratura": "semisquare",
		"sesquicuadratura": "sesquisquare", "quincuncio": "quincunx", "semisextil": "semisextile"},
	f: map[string]string{
		"pos":    "%s expresses itself %s, in the field of %s",
		"asp":    "%s %s %s.",
		"reg":    "What concerns %s passes through %s: the ruler of house %d is lodged in the %dth.",
		"contra": "%s receives at once a %s from %s (orb %.2f°) and a %s from %s (orb %.2f°). By exactness %s rules, but both are there: describe them together, with \"and\", not with \"but\".",
		"tension": "the tension", "facilidad": "the ease",
		"dom":    "%s (ruler of the Ascendant: %s)",
		"nota":   "This is step 1 of 4: a literal translation, deliberately clumsy. Grouping into categories is step 2 and is already done. Steps 3 and 4 — resolving the contradictions and writing it as running prose — are yours. A chart does not square itself.",
	},
}

type Frase struct {
	Categoria string  `json:"categoria"`
	Texto     string  `json:"texto"`
	Fuente    string  `json:"fuente"`
	Peso      float64 `json:"peso"`
}

type Lectura struct {
	Frases          []Frase  `json:"frases"`
	Contradicciones []string `json:"contradicciones"`
	Dominante       string   `json:"dominante"`
	Nota            string   `json:"nota"`
}

// duro indica si un aspecto es de tensión.
func duro(n string) bool {
	return n == "cuadratura" || n == "oposición" || n == "semicuadratura" || n == "sesquicuadratura"
}
func blando(n string) bool { return n == "trígono" || n == "sextil" }

// Interpretar compone la traducción literal de la carta y agrupa por categorías.
// NO sintetiza: eso es trabajo del que lee, y el módulo 10 explica por qué.
func Interpretar(c efem.Carta, lang string) Lectura {
	T := es
	if lang == "en" {
		T = en
	}
	// nom traduce nombres propios de planeta y de aspecto al salir.
	nom := func(s string) string {
		if T.nombres != nil {
			if v, ok := T.nombres[s]; ok {
				return v
			}
		}
		return s
	}
	f := T.f
	var L Lectura
	for _, p := range c.Cuerpos {
		if p.Nombre == "Nodo Norte" || p.Nombre == "Nodo Sur" {
			continue
		}
		t := fmt.Sprintf(f["pos"], T.planeta[p.Nombre], T.signo[p.SignoIdx], T.casa[p.CasaP-1])
		if m, ok := T.dign[p.Dignidad]; ok {
			t += " — " + m
		}
		t += "."
		L.Frases = append(L.Frases, Frase{Categoria: T.cat[p.CasaP-1], Texto: t,
			Fuente: fmt.Sprintf("%s %d°, casa %d", nom(p.Nombre), int(p.Grado), p.CasaP), Peso: 1})
	}
	// aspectos: solo los que pesan de verdad
	for _, a := range c.Aspectos {
		if a.Orbe > 4 {
			continue
		}
		t := fmt.Sprintf(f["asp"], T.planeta[a.A], T.aspecto[a.Nombre], T.planeta[a.B])
		L.Frases = append(L.Frases, Frase{Categoria: T.cat[0], Texto: t,
			Fuente: fmt.Sprintf("%s %s %s, %.2f°", nom(a.A), nom(a.Nombre), nom(a.B), a.Orbe),
			Peso: 1 + (4-a.Orbe)/2})
	}
	// regentes: lo que convierte casas sueltas en argumentos
	for i, r := range c.Regentes {
		if c.RegenteEn[i] == i+1 {
			continue
		}
		L.Frases = append(L.Frases, Frase{Categoria: T.cat[i],
			Texto: fmt.Sprintf(f["reg"], T.casa[i], T.casa[c.RegenteEn[i]-1], i+1, c.RegenteEn[i]),
			Fuente: fmt.Sprintf("%s → %d", nom(r), c.RegenteEn[i]), Peso: 1.2})
	}
	// contradicciones: mismo planeta con aspectos duros y blandos a la vez
	porPlaneta := map[string][]efem.Aspecto{}
	for _, a := range c.Aspectos {
		if a.Orbe <= 6 {
			porPlaneta[a.A] = append(porPlaneta[a.A], a)
			porPlaneta[a.B] = append(porPlaneta[a.B], a)
		}
	}
	for _, n := range efem.Orden[:10] {
		var d, b *efem.Aspecto
		for i := range porPlaneta[n] {
			a := porPlaneta[n][i]
			if duro(a.Nombre) && (d == nil || a.Orbe < d.Orbe) {
				d = &porPlaneta[n][i]
			}
			if blando(a.Nombre) && (b == nil || a.Orbe < b.Orbe) {
				b = &porPlaneta[n][i]
			}
		}
		if d != nil && b != nil {
			otroD, otroB := d.A, b.A
			if otroD == n {
				otroD = d.B
			}
			if otroB == n {
				otroB = b.B
			}
			manda := f["tension"]
			if b.Orbe < d.Orbe {
				manda = f["facilidad"]
			}
			L.Contradicciones = append(L.Contradicciones, fmt.Sprintf(f["contra"],
				nom(n), nom(d.Nombre), nom(otroD), d.Orbe, nom(b.Nombre), nom(otroB), b.Orbe, manda))
		}
	}
	L.Dominante = dominante(c, nom, f)
	L.Nota = f["nota"]
	return L
}

// dominante estima qué planeta manda: regencia del Ascendente, dignidad,
// angularidad y aspectos exactos.
func dominante(c efem.Carta, nom func(string) string, f map[string]string) string {
	p := map[string]float64{}
	regAsc := efem.RegenteSigno[int(c.Asc/30)]
	p[regAsc] += 3
	for _, b := range c.Cuerpos {
		switch b.Dignidad {
		case "domicilio":
			p[b.Nombre] += 2
		case "exaltación":
			p[b.Nombre] += 1.5
		}
		switch b.CasaP {
		case 1, 4, 7, 10:
			p[b.Nombre] += 1.5
		case 2, 5, 8, 11:
			p[b.Nombre] += 0.5
		}
	}
	for _, a := range c.Aspectos {
		if a.Orbe <= 2 {
			p[a.A] += 0.7
			p[a.B] += 0.7
		}
	}
	mejor, max := regAsc, -1.0
	for _, n := range efem.Orden[:10] {
		if p[n] > max {
			mejor, max = n, p[n]
		}
	}
	return fmt.Sprintf(f["dom"], nom(mejor), nom(regAsc))
}
