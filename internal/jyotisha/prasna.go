package jyotisha

import "fmt"

// Praśna: la carta del momento en que se hace la pregunta.
//
// No es la carta de la persona sino la del instante, y se lee para ESA
// pregunta. Es el recurso clásico cuando no hay hora de nacimiento fiable —
// que, como dice el módulo 2, es casi siempre.
//
// Lo primero que hace un praśna bien hecho no es contestar: es decidir si la
// pregunta se puede contestar. Las clásicas dan varias señales de que no, y
// están todas aquí. Decir «esta pregunta no está madura, vuelve más tarde» es
// una respuesta completa, igual que «hasta aquí llega la carta» en el módulo 11.

type Prasna struct {
	Tema       string   `json:"tema"`
	Bhava      int      `json:"bhava"`      // el bhāva del asunto preguntado
	Apta       bool     `json:"apta"`       // ¿se puede juzgar la pregunta?
	Reparos    []string `json:"reparos"`    // por qué no, si no
	Lagna      string   `json:"lagna"`
	SenorLagna string   `json:"senorLagna"`
	Frases     []string `json:"frases"`
	Nota       string   `json:"nota"`
}

// TemasPrasna son los asuntos que se pueden preguntar, con el bhāva que los
// gobierna. No es una lista cerrada por capricho: cada asunto tiene su casa, y
// preguntar sin saber cuál es equivale a no preguntar.
var TemasPrasna = []struct {
	Clave string
	Bhava int
}{
	{"salud", 1}, {"dinero", 2}, {"hermanos", 3}, {"madre", 4}, {"casa", 4},
	{"hijos", 5}, {"estudios", 5}, {"enemigos", 6}, {"enfermedad", 6},
	{"pareja", 7}, {"sociedad", 7}, {"herencia", 8}, {"padre", 9}, {"viaje", 9},
	{"trabajo", 10}, {"ganancia", 11}, {"perdida", 12},
}

func bhavaDelTema(tema string) int {
	for _, t := range TemasPrasna {
		if t.Clave == tema {
			return t.Bhava
		}
	}
	return 1
}

// HacerPrasna juzga la carta del momento para el tema pedido.
func HacerPrasna(c Carta, tema, lang string) Prasna {
	T := es
	if lang == "en" {
		T = en
	}
	f := T.frases
	nom := func(g string) string {
		if lang == "en" {
			if v, ok := grahaEn[g]; ok {
				return v
			}
		}
		return g
	}
	p := Prasna{Tema: tema, Bhava: bhavaDelTema(tema), Apta: true,
		Lagna: c.LagnaPos, SenorLagna: nom(c.SenorLagna)}

	// ── ¿está la pregunta madura? ──
	//
	// Los primeros y últimos grados de un signo: el asunto acaba de empezar o
	// ya se ha acabado, y en ninguno de los dos casos hay nada que juzgar.
	grado := c.Lagna - float64(c.LagnaRasi)*30
	if grado < 3 {
		p.Apta = false
		p.Reparos = append(p.Reparos, f["pr_lagnaJoven"])
	}
	if grado > 27 {
		p.Apta = false
		p.Reparos = append(p.Reparos, f["pr_lagnaViejo"])
	}
	if Gandanta(c.Lagna) {
		p.Apta = false
		p.Reparos = append(p.Reparos, f["pr_lagnaGandanta"])
	}
	// La Luna es la mente del que pregunta. Si está en gaṇḍānta, la pregunta
	// sale de un nudo y no de una duda.
	for _, g := range c.Grahas {
		if g.Nombre != "Luna" {
			continue
		}
		if g.Gandanta {
			p.Apta = false
			p.Reparos = append(p.Reparos, f["pr_lunaGandanta"])
		}
		if g.Bhava == 6 || g.Bhava == 8 || g.Bhava == 12 {
			p.Reparos = append(p.Reparos, fmt.Sprintf(f["pr_lunaDuhsthana"], g.Bhava))
		}
	}

	// ── el juicio ──
	b := c.Bhavas[p.Bhava-1]
	p.Frases = append(p.Frases, fmt.Sprintf(f["pr_bhava"], p.Bhava, T.bhava[p.Bhava-1]))
	p.Frases = append(p.Frases, fmt.Sprintf(f["pr_senor"], p.Bhava, nom(b.Senor), b.SenorEn,
		T.bhava[b.SenorEn-1]))

	if len(b.Ocupan) > 0 {
		var buenos, malos []string
		for _, g := range b.Ocupan {
			if benefico(g) == 1 {
				buenos = append(buenos, nom(g))
			} else if benefico(g) == -1 {
				malos = append(malos, nom(g))
			}
		}
		if len(buenos) > 0 {
			p.Frases = append(p.Frases, fmt.Sprintf(f["pr_ocupanBien"], lista(buenos)))
		}
		if len(malos) > 0 {
			p.Frases = append(p.Frases, fmt.Sprintf(f["pr_ocupanMal"], lista(malos)))
		}
	} else {
		p.Frases = append(p.Frases, f["pr_vacio"])
	}

	// La mirada de Júpiter sobre la casa preguntada es la mejor señal que hay.
	for _, g := range b.Aspectan {
		if g == "Júpiter" {
			p.Frases = append(p.Frases, f["pr_guruMira"])
		}
	}

	// El señor del lagna y el del asunto: si están juntos o se miran, lo que se
	// pregunta y el que pregunta están conectados, y el asunto avanza.
	sl, sa := c.SenorLagna, b.Senor
	if sl == sa {
		p.Frases = append(p.Frases, f["pr_mismoSenor"])
	} else {
		var rl, ra int = -1, -1
		for _, g := range c.Grahas {
			if g.Nombre == sl {
				rl = g.RasiIdx
			}
			if g.Nombre == sa {
				ra = g.RasiIdx
			}
		}
		if rl >= 0 && ra >= 0 {
			switch {
			case rl == ra:
				p.Frases = append(p.Frases, fmt.Sprintf(f["pr_juntos"], nom(sl), nom(sa)))
			case seMiran(sl, sa, rl, ra):
				p.Frases = append(p.Frases, fmt.Sprintf(f["pr_seMiran"], nom(sl), nom(sa)))
			default:
				p.Frases = append(p.Frases, fmt.Sprintf(f["pr_sinRelacion"], nom(sl), nom(sa)))
			}
		}
	}

	// El aṣṭakavarga del bhāva preguntado: es el único número que da el
	// sistema, y aquí sirve justo para eso, para no decidir a ojo.
	sav := c.Ashtaka.SAV[b.RasiIdx]
	switch {
	case sav >= 30:
		p.Frases = append(p.Frases, fmt.Sprintf(f["pr_savAlto"], sav))
	case sav <= 25:
		p.Frases = append(p.Frases, fmt.Sprintf(f["pr_savBajo"], sav))
	default:
		p.Frases = append(p.Frases, fmt.Sprintf(f["pr_savMedio"], sav))
	}

	p.Nota = f["pr_nota"]
	return p
}

func lista(x []string) string {
	out := ""
	for i, s := range x {
		if i > 0 {
			out += ", "
		}
		out += s
	}
	return out
}
