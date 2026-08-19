package jyotisha

import "fmt"

// Lectura védica. Mismo criterio que el motor occidental: la app COMPONE la
// frase a partir de sus piezas y no guarda interpretaciones prefabricadas. Es
// traducción literal, deliberadamente torpe; cuadrar la carta es trabajo del
// que lee, no del programa.
//
// Lo que distingue a este motor del occidental es la CADENA DE SEÑORES. En
// jyotiṣa un bhāva no se lee solo por quien lo ocupa: se lee por dónde está su
// señor, y por dónde está el señor de aquel. Ahí es donde una causa aparece
// metida dentro de otra, que es la manera védica de razonar.
//
// Todo el texto sale en el idioma que se pida. Nada de cadenas fijas.

type tabla struct {
	graha  map[string]string
	rasi   [12]string
	bhava  [12]string
	nak    [27]string
	dign   map[string]string
	cat    [12]string
	cats   []string
	frases map[string]string
}

var es = tabla{
	graha: map[string]string{
		// Sujetos en singular: la plantilla lleva el verbo en singular detrás.
		"Sol": "tu alma", "Luna": "tu mente",
		"Marte": "tu empuje", "Mercurio": "tu inteligencia",
		"Júpiter": "tu criterio", "Venus": "tu deseo",
		"Saturno": "tu resistencia",
		"Rāhu": "aquello que te obsesiona", "Ketu": "aquello que ya traes sabido",
	},
	rasi: [12]string{"con arranque y sin cálculo", "con paciencia y apego a lo tangible",
		"con curiosidad y en varias direcciones", "con memoria y ganas de amparar",
		"con orgullo y ganas de mandar", "con detalle y afán de corregir",
		"buscando trato y equilibrio", "con hondura y sin soltar",
		"con amplitud y sentido de la norma", "con frialdad y mirada larga",
		"con distancia y criterio propio", "con porosidad y sin bordes"},
	bhava: [12]string{"tu cuerpo y tu manera de presentarte", "lo que ganas, guardas y dices",
		"tu iniciativa, tus hermanos y tu esfuerzo corto", "tu madre, tu casa y tu paz interior",
		"tu inteligencia, tus hijos y lo que traes de otras vidas",
		"la enfermedad, la deuda y el enemigo", "tu pareja y todo pacto entre dos",
		"lo heredado, lo oculto y lo que te transforma de golpe",
		"tu padre, tu suerte y aquello en lo que crees",
		"tu obra, tu oficio y lo que se ve de ti", "lo que ganas sin buscarlo y tus mayores",
		"lo que gastas, sueltas y pierdes de vista"},
	nak: [27]string{"con prisa por curar", "con ganas de arrancar de raíz", "con filo que corta",
		"con hambre de crecer", "buscando sin parar", "con la herida abierta",
		"volviendo siempre a empezar", "alimentando a los suyos", "enroscándose en lo que teme",
		"con orgullo de linaje", "buscando gozo", "cumpliendo el pacto", "con la mano diestra",
		"puliendo la forma", "doblándose sin romperse", "queriendo dos cosas a la vez",
		"guardando la amistad", "queriendo mandar", "arrancando de cuajo", "sin cansarse nunca",
		"venciendo tarde y de veras", "escuchando", "marcando el ritmo", "curando en secreto",
		"llevando lo que quema", "cumpliendo lo prometido", "llegando a puerto"},
	dign: map[string]string{
		"exaltado":     "y aquí va sobrado, quizá de más",
		"debilitado":   "y aquí trabaja a contrapelo, lo que con los años da oficio",
		"signo propio": "y aquí está en su casa y rinde entero",
		"mūlatrikoṇa":  "y aquí está en su mejor asiento",
	},
	cat: [12]string{"Cuerpo y carácter", "Dinero y palabra", "Empuje y hermanos", "Madre y raíz",
		"Mente e hijos", "Salud y enemigos", "Pareja", "Lo oculto y las crisis",
		"Padre y creencia", "Obra y oficio", "Ganancias", "Gasto y soltar"},
	frases: map[string]string{
		"pos":     "%s se expresa %s, en el terreno de %s",
		"nakdet":  ", %s (nakṣatra %s, pada %d, de %s)",
		"retro":   " Va retrógrado: lo suyo se vuelve hacia dentro antes de salir.",
		"comb":    " Está combusto, a %.1f° del Sol: su asunto queda tapado por el del padre o la autoridad.",
		"gand":    " Cae en gaṇḍānta: nudo kármico, y ahí se pide cuidado, no alarma.",
		"dig":     " Tiene dig-bala: está en la dirección donde más rinde.",
		"cadena":  "Lo relativo a %s pasa por %s: el señor del bhāva %d (%s) está alojado en el %d.",
		"cadena2": "Y de ahí sigue: %s, señor del bhāva %d, está en el %d — o sea que %s depende a su vez de %s. Una causa dentro de otra.",
		"propia":  "El señor del bhāva %d (%s) está en su propio bhāva: el asunto se sostiene solo, sin depender de otro.",
		"mira":    "El bhāva %d (%s) recibe la mirada de %s. %s",
		"benef":   "Es mirada benéfica: protege y ablanda.",
		"malef":   "Es mirada maléfica: aprieta, y lo que aprieta también forma.",
		"mixta":   "Le miran a la vez benéficos y maléficos: la casa está disputada.",
		"dasa":    "Ahora corre la daśā de %s hasta %s. %s es señor del bhāva %d y está en el %d, así que el periodo tiende a traer los asuntos de %s por la vía de %s.",
		"antar":   " Dentro de ella corre la antardaśā de %s hasta %s, que estrecha el asunto.",
		"atma":    "El ātmakāraka es %s: el hilo del alma va por %s. Mira dónde cae en el D-9 antes que ninguna otra cosa.",
		"amatya":  "El amātyakāraka es %s: la carrera y el sustento se explican por ahí.",
		"f_senor": "señor de %d en %d",
		"f_senorEn": "señor de %d (%s) en %d",
		"f_drishti": "dṛṣṭi sobre %d",
		"f_dasa": "vimśottarī %s",
		"f_pos": "%s %s, bhāva %d",
		"contra":  "El bhāva %d (%s) está disputado: lo ocupan o lo miran %s por un lado y %s por otro. No elijas un lado. Las dos cosas están ahí y hay que describirlas juntas, con «y», no con «pero».",
		"dom":     "%s. El señor del lagna es %s, en el bhāva %d.",
		"nota":    "Esto es traducción literal, no una lectura. Agrupar por temas ya está hecho; resolver lo que se contradice y escribirlo seguido lo haces tú. Y cuando algo salga borroso, no le pongas adjetivos: haz más cálculos. Ahí está el curso, módulo 16.",
	},
}

var en = tabla{
	graha: map[string]string{
		"Sol": "your soul", "Luna": "your mind",
		"Marte": "your drive", "Mercurio": "your intelligence",
		"Júpiter": "your judgement", "Venus": "your desire",
		"Saturno": "your endurance",
		"Rāhu": "what you crave", "Ketu": "what you already know",
	},
	rasi: [12]string{"headlong, without reckoning", "patiently, attached to what is solid",
		"curiously, in several directions at once", "with memory and a wish to shelter",
		"with pride and a wish to command", "in detail, eager to correct",
		"seeking terms and balance", "deeply, and without letting go",
		"broadly, with a sense of the rule", "coldly, taking the long view",
		"at a distance, on your own terms", "porously, without clear edges"},
	bhava: [12]string{"your body and the way you present yourself", "what you earn, keep and say",
		"your initiative, your siblings and your short efforts", "your mother, your home and your inner peace",
		"your intelligence, your children and what you bring from other lives",
		"illness, debt and the enemy", "your partner and every pact between two",
		"what is inherited, what is hidden and what transforms you at a stroke",
		"your father, your fortune and what you believe",
		"your work, your calling and what shows of you", "what comes to you unsought, and your elders",
		"what you spend, release and lose sight of"},
	nak: [27]string{"in a hurry to heal", "wanting to pull it up by the root", "with an edge that cuts",
		"hungry to grow", "searching without stopping", "with the wound still open",
		"always coming back to begin again", "feeding its own", "coiling around what it fears",
		"proud of its line", "seeking delight", "keeping the bargain", "with a deft hand",
		"polishing the form", "bending without breaking", "wanting two things at once",
		"keeping the friendship", "wanting to command", "tearing it out whole", "never tiring",
		"winning late and for real", "listening", "keeping the beat", "healing in secret",
		"carrying what burns", "keeping its word", "coming into harbour"},
	dign: map[string]string{
		"exaltado":     "and here it runs high, perhaps too high",
		"debilitado":   "and here it works against the grain, which over the years makes a craftsman",
		"signo propio": "and here it is at home and gives its full measure",
		"mūlatrikoṇa":  "and here it sits in its best seat",
	},
	cat: [12]string{"Body and character", "Money and speech", "Drive and siblings", "Mother and root",
		"Mind and children", "Health and enemies", "Partner", "The hidden and the crises",
		"Father and belief", "Work and calling", "Gains", "Spending and letting go"},
	frases: map[string]string{
		"pos":     "%s expresses itself %s, in the field of %s",
		"nakdet":  ", %s (nakṣatra %s, pada %d, of %s)",
		"retro":   " It is retrograde: its business turns inward before it comes out.",
		"comb":    " It is combust, %.1f° from the Sun: its matter is covered over by that of the father or the authority.",
		"gand":    " It falls in gaṇḍānta: a karmic knot, and that calls for care, not alarm.",
		"dig":     " It has dig-bala: it stands in the direction where it gives most.",
		"cadena":  "What concerns %s passes through %s: the lord of bhāva %d (%s) is lodged in the %dth.",
		"cadena2": "And from there it goes on: %s, lord of bhāva %d, is in the %dth — so %s depends in turn on %s. One cause inside another.",
		"propia":  "The lord of bhāva %d (%s) is in its own bhāva: the matter stands on its own, depending on nothing else.",
		"mira":    "Bhāva %d (%s) receives the gaze of %s. %s",
		"benef":   "It is a benefic gaze: it protects and softens.",
		"malef":   "It is a malefic gaze: it presses, and what presses also forms.",
		"mixta":   "Benefics and malefics look at it at once: the house is contested.",
		"dasa":    "The daśā of %s is running now, until %s. %s is lord of bhāva %d and sits in the %dth, so the period tends to bring the matters of %s by way of %s.",
		"antar":   " Within it the antardaśā of %s runs until %s, which narrows the matter.",
		"atma":    "The ātmakāraka is %s: the soul's thread runs through %s. Look at where it falls in the D-9 before anything else.",
		"amatya":  "The amātyakāraka is %s: career and livelihood are explained from there.",
		"f_senor": "lord of %d in %d",
		"f_senorEn": "lord of %d (%s) in %d",
		"f_drishti": "dṛṣṭi on %d",
		"f_dasa": "vimśottarī %s",
		"f_pos": "%s %s, bhāva %d",
		"contra":  "Bhāva %d (%s) is contested: %s occupy or look at it on one side, %s on the other. Do not pick a side. Both are there, and they have to be described together, with \"and\", not with \"but\".",
		"dom":     "%s. The lord of the lagna is %s, in bhāva %d.",
		"nota":    "This is a literal translation, not a reading. Grouping by theme is already done; resolving what contradicts and writing it as running prose is your job. And when something comes out blurred, do not add adjectives to it: do more calculation. That is module 16 of the course.",
	},
}

// Los nombres de graha se guardan en español dentro del motor porque así los
// produce el cálculo. Aquí se traducen al salir.
var grahaEn = map[string]string{"Sol": "Sun", "Luna": "Moon", "Marte": "Mars",
	"Mercurio": "Mercury", "Júpiter": "Jupiter", "Venus": "Venus", "Saturno": "Saturn",
	"Rāhu": "Rāhu", "Ketu": "Ketu"}

type FraseVed struct {
	Categoria string  `json:"categoria"`
	Texto     string  `json:"texto"`
	Fuente    string  `json:"fuente"`
	Peso      float64 `json:"peso"`
}

type LecturaVed struct {
	Frases          []FraseVed `json:"frases"`
	Contradicciones []string   `json:"contradicciones"`
	Dominante       string     `json:"dominante"`
	Nota            string     `json:"nota"`
}

// benefico clasifica los grahas al modo clásico. Mercurio y la Luna dependen
// de la compañía; se toman como neutros y no se fuerzan a un bando.
func benefico(g string) int {
	switch g {
	case "Júpiter", "Venus":
		return 1
	case "Saturno", "Marte", "Sol", "Rāhu", "Ketu":
		return -1
	}
	return 0
}

// Interpretar compone la lectura literal. lang es "es" o "en".
func Interpretar(c Carta, lang string) LecturaVed {
	T := es
	if lang == "en" {
		T = en
	}
	nom := func(g string) string {
		if lang == "en" {
			if v, ok := grahaEn[g]; ok {
				return v
			}
		}
		return g
	}
	nombres := func(l []string) string {
		out := ""
		for i, g := range l {
			if i > 0 {
				out += ", "
			}
			out += nom(g)
		}
		return out
	}
	f := T.frases
	var L LecturaVed
	add := func(cat, txt, fte string, peso float64) {
		L.Frases = append(L.Frases, FraseVed{Categoria: cat, Texto: txt, Fuente: fte, Peso: peso})
	}

	// 1. Cada graha: función, modo, terreno, nakṣatra y estado.
	for _, g := range c.Grahas {
		if g.Bhava < 1 || g.Bhava > 12 {
			continue
		}
		nakIdx := int(g.Lon / (360.0 / 27.0))
		if nakIdx > 26 {
			nakIdx = 26
		}
		t := fmt.Sprintf(f["pos"], T.graha[g.Nombre], T.rasi[g.RasiIdx], T.bhava[g.Bhava-1])
		t += fmt.Sprintf(f["nakdet"], T.nak[nakIdx], g.Nak, g.Pada, nom(g.SenorNak))
		if m, ok := T.dign[g.Dignidad]; ok {
			t += " — " + m
		}
		if g.Mula {
			t += " — " + T.dign["mūlatrikoṇa"]
		}
		t += "."
		if g.Retro {
			t += f["retro"]
		}
		if g.Combusto {
			t += fmt.Sprintf(f["comb"], g.DelSol)
		}
		if g.Gandanta {
			t += f["gand"]
		}
		if g.DigBala {
			t += f["dig"]
		}
		add(T.cat[g.Bhava-1], t, fmt.Sprintf(f["f_pos"], nom(g.Nombre), g.Posicion, g.Bhava), 1)
	}

	// 2. La cadena de señores. Es el corazón del método: dice de qué depende
	// cada asunto, y luego de qué depende aquello.
	senorEn := map[int]int{}
	for _, b := range c.Bhavas {
		senorEn[b.Numero] = b.SenorEn
	}
	for _, b := range c.Bhavas {
		if b.SenorEn < 1 || b.SenorEn > 12 {
			continue
		}
		if b.SenorEn == b.Numero {
			add(T.cat[b.Numero-1], fmt.Sprintf(f["propia"], b.Numero, nom(b.Senor)),
				fmt.Sprintf(f["f_senor"], b.Numero, b.SenorEn), 1.3)
			continue
		}
		t := fmt.Sprintf(f["cadena"], T.bhava[b.Numero-1], T.bhava[b.SenorEn-1],
			b.Numero, nom(b.Senor), b.SenorEn)
		// el segundo escalón: dónde está el señor de aquel bhāva
		seg := c.Bhavas[b.SenorEn-1]
		if seg.SenorEn >= 1 && seg.SenorEn <= 12 && seg.SenorEn != b.SenorEn && seg.SenorEn != b.Numero {
			t += " " + fmt.Sprintf(f["cadena2"], nom(seg.Senor), seg.Numero, seg.SenorEn,
				T.bhava[b.Numero-1], T.bhava[seg.SenorEn-1])
		}
		add(T.cat[b.Numero-1], t, fmt.Sprintf(f["f_senorEn"], b.Numero, nom(b.Senor), b.SenorEn), 1.4)
	}

	// 3. Dṛṣṭi: quién mira cada bhāva. En jyotiṣa el aspecto es por signo entero.
	for _, b := range c.Bhavas {
		if len(b.Aspectan) == 0 {
			continue
		}
		pos, neg := 0, 0
		for _, g := range b.Aspectan {
			switch benefico(g) {
			case 1:
				pos++
			case -1:
				neg++
			}
		}
		q := f["mixta"]
		switch {
		case pos > 0 && neg == 0:
			q = f["benef"]
		case neg > 0 && pos == 0:
			q = f["malef"]
		}
		add(T.cat[b.Numero-1], fmt.Sprintf(f["mira"], b.Numero, T.bhava[b.Numero-1],
			nombres(b.Aspectan), q),
			fmt.Sprintf(f["f_drishti"], b.Numero), 1.1)

		// 4. Contradicción: la misma casa tirada por los dos lados.
		if pos > 0 && neg > 0 {
			var bs, ms []string
			for _, g := range append(append([]string{}, b.Ocupan...), b.Aspectan...) {
				if benefico(g) == 1 {
					bs = append(bs, g)
				} else if benefico(g) == -1 {
					ms = append(ms, g)
				}
			}
			L.Contradicciones = append(L.Contradicciones, fmt.Sprintf(f["contra"],
				b.Numero, T.bhava[b.Numero-1], nombres(bs), nombres(ms)))
		}
	}

	// 5. La daśā que corre ahora: es lo que fecha todo lo anterior.
	for _, d := range c.Dasas {
		if !d.Actual {
			continue
		}
		bh := 0
		for _, b := range c.Bhavas {
			if b.Senor == d.Senor {
				bh = b.Numero
				break
			}
		}
		en := 0
		for _, g := range c.Grahas {
			if g.Nombre == d.Senor {
				en = g.Bhava
			}
		}
		if bh == 0 || en == 0 {
			break
		}
		t := fmt.Sprintf(f["dasa"], nom(d.Senor), d.Hasta, nom(d.Senor), bh, en,
			T.bhava[bh-1], T.bhava[en-1])
		for _, s := range d.Sub {
			if s.Actual {
				t += fmt.Sprintf(f["antar"], nom(s.Senor), s.Hasta)
				break
			}
		}
		add(T.cat[en-1], t, fmt.Sprintf(f["f_dasa"], nom(d.Senor)), 2)
		break
	}

	// 6. Los kārakas de Jaimini: por dónde va el hilo del alma.
	if ak := c.Karakas["Ātmakāraka"]; ak != "" {
		add(T.cat[0], fmt.Sprintf(f["atma"], nom(ak), T.graha[ak]), "Ātmakāraka", 1.8)
	}
	if am := c.Karakas["Amātyakāraka"]; am != "" {
		add(T.cat[9], fmt.Sprintf(f["amatya"], nom(am), T.graha[am]), "Amātyakāraka", 1.5)
	}

	L.Dominante = dominanteVed(c, T, nom, lang)
	L.Nota = f["nota"]
	return L
}

// dominanteVed pesa quién manda: señor del lagna, dignidad, kendra/trikona,
// dig-bala y ātmakāraka. No es Ṣaḍbala; es una estimación declarada como tal.
func dominanteVed(c Carta, T tabla, nom func(string) string, lang string) string {
	p := map[string]float64{}
	p[c.SenorLagna] += 3
	if ak := c.Karakas["Ātmakāraka"]; ak != "" {
		p[ak] += 2
	}
	for _, g := range c.Grahas {
		switch g.Dignidad {
		case "exaltado":
			p[g.Nombre] += 2
		case "signo propio":
			p[g.Nombre] += 1.5
		case "debilitado":
			p[g.Nombre] -= 1
		}
		if g.Mula {
			p[g.Nombre] += 1
		}
		if g.DigBala {
			p[g.Nombre] += 1
		}
		if g.Combusto {
			p[g.Nombre] -= 1
		}
		switch g.Bhava {
		case 1, 4, 7, 10:
			p[g.Nombre] += 1.5
		case 5, 9:
			p[g.Nombre] += 1
		case 6, 8, 12:
			p[g.Nombre] -= 0.5
		}
	}
	mejor, max := c.SenorLagna, -99.0
	for _, g := range Grahas {
		if p[g] > max {
			mejor, max = g, p[g]
		}
	}
	bh := 0
	for _, g := range c.Grahas {
		if g.Nombre == c.SenorLagna {
			bh = g.Bhava
		}
	}
	quien := mejor + " manda en esta carta"
	if lang == "en" {
		quien = nom(mejor) + " rules this chart"
	} else {
		quien = nom(mejor) + " manda en esta carta"
	}
	return fmt.Sprintf(T.frases["dom"], quien, nom(c.SenorLagna), bh)
}
