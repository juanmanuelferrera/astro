package jyotisha

import (
	"fmt"
	"math"
	"sort"
	"time"

	"astro/internal/efem"
)

type Graha struct {
	Nombre    string  `json:"nombre"`
	Glifo     string  `json:"glifo"`
	Lon       float64 `json:"lon"`
	Rasi      string  `json:"rasi"`
	RasiIdx   int     `json:"rasiIdx"`
	Grado     float64 `json:"grado"`
	Nak       string  `json:"nak"`
	Pada      int     `json:"pada"`
	SenorNak  string  `json:"senorNak"`
	Bhava     int     `json:"bhava"`
	Dignidad  string  `json:"dignidad"`
	Retro     bool    `json:"retro"`
	Gandanta  bool    `json:"gandanta"`
	Posicion  string  `json:"posicion"`
	Mula      bool    `json:"mula"`
	Combusto  bool    `json:"combusto"`
	DelSol    float64 `json:"delSol"`
	DigBala   bool    `json:"digBala"`
}

type Bhava struct {
	Numero  int      `json:"numero"`
	Rasi    string   `json:"rasi"`
	RasiIdx int      `json:"rasiIdx"`
	Senor   string   `json:"senor"`
	SenorEn int      `json:"senorEn"`
	Ocupan  []string `json:"ocupan"`
	Aspectan []string `json:"aspectan"`
}

type Periodo struct {
	Senor  string `json:"senor"`
	Desde  string `json:"desde"`
	Hasta  string `json:"hasta"`
	Anios  float64 `json:"anios"`
	Sub    []Periodo `json:"sub,omitempty"`
	Actual bool   `json:"actual"`
}

type Carta struct {
	JD         float64            `json:"jd"`
	UT         string             `json:"ut"`
	Ayanamsa   float64            `json:"ayanamsa"`
	Lagna      float64            `json:"lagna"`
	LagnaRasi  int                `json:"lagnaRasi"`
	LagnaNak   string             `json:"lagnaNak"`
	LagnaPada  int                `json:"lagnaPada"`
	LagnaPos   string             `json:"lagnaPos"`
	SenorLagna string             `json:"senorLagna"`
	MC         float64            `json:"mc"`
	Grahas     []Graha            `json:"grahas"`
	Bhavas     []Bhava            `json:"bhavas"`
	Vargas     map[string][]Graha `json:"vargas"`
	Karakas    map[string]string  `json:"karakas"`
	Dasas      []Periodo          `json:"dasas"`
	Yogas      []string           `json:"yogas"`
	Karakamsa  string             `json:"karakamsa"`
	Gocara     Gocara             `json:"gocara"`
	Pancanga   Pancanga           `json:"pancanga"`
	Ashtaka    Ashtakavarga       `json:"ashtaka"`
	Shadbala   Shadbala           `json:"shadbala"`
	Arudhas    Arudhas            `json:"arudhas"`
	LagnasEsp  LagnasEsp          `json:"lagnasEsp"`
	NodoVerdad bool               `json:"nodoVerdad"`
}

func norm(x float64) float64 {
	x = math.Mod(x, 360)
	if x < 0 {
		x += 360
	}
	return x
}

// ─────────────── vargas ───────────────

// Varga devuelve el signo (0-11) que ocupa una longitud en la divisional n.
func Varga(lon float64, n int) int {
	s := int(lon / 30)
	g := math.Mod(lon, 30)
	par := s%2 == 1 // 0-indexado: impar del índice = signo par (Vṛṣabha, Karka…)
	switch n {
	case 1:
		return s
	case 2: // hora
		if !par {
			if g < 15 {
				return 4
			}
			return 3
		}
		if g < 15 {
			return 3
		}
		return 4
	case 3: // drekkāṇa
		return (s + 4*int(g/10)) % 12
	case 7: // saptāṁśa
		base := s
		if par {
			base = (s + 6) % 12
		}
		return (base + int(g/(30.0/7))) % 12
	case 9: // navāṁśa
		return ([4]int{0, 9, 6, 3}[s%4] + int(g/(30.0/9))) % 12
	case 10: // daśāṁśa
		base := s
		if par {
			base = (s + 8) % 12
		}
		return (base + int(g/3)) % 12
	case 12: // dvādaśāṁśa
		return (s + int(g/2.5)) % 12
	case 16:
		base := 0
		if s%3 == 1 {
			base = 4
		} else if s%3 == 2 {
			base = 8
		}
		return (base + int(g/(30.0/16))) % 12
	case 30: // triṁśāṁśa
		lim := []float64{5, 10, 18, 25, 30}
		impar := []int{0, 10, 8, 2, 6}
		parS := []int{1, 5, 11, 9, 7}
		for i, l := range lim {
			if g < l {
				if !par {
					return impar[i]
				}
				return parS[i]
			}
		}
		return s
	case 60:
		return (s + int(g/0.5)) % 12
	}
	return s
}

var VargasUsadas = []struct {
	Clave  string
	N      int
	Asunto string
}{
	{"D1", 1, "la vida tal como se vive"},
	{"D2", 2, "riqueza y sustento"},
	{"D3", 3, "hermanos, coraje, esfuerzo"},
	{"D7", 7, "hijos y descendencia"},
	{"D9", 9, "el alma, el cónyuge, la fuerza real"},
	{"D10", 10, "profesión y karma público"},
	{"D12", 12, "los padres"},
	{"D16", 16, "vehículos y confort"},
	{"D30", 30, "males y desgracias"},
	{"D60", 60, "karma acumulado"},
}

// ─────────────── vimśottarī ───────────────

var ciclo = []struct {
	senor string
	anios float64
}{{"Ketu", 7}, {"Venus", 20}, {"Sol", 6}, {"Luna", 10}, {"Marte", 7},
	{"Rāhu", 18}, {"Júpiter", 16}, {"Saturno", 19}, {"Mercurio", 17}}

const anioSideral = 365.2425

// Vimsottari calcula las mahādaśās desde la Luna, con sus bhuktis.
func Vimsottari(lonLuna float64, nacimiento time.Time, hasta int) []Periodo {
	const span = 360.0 / 27.0
	ni := int(lonLuna / span)
	frac := math.Mod(lonLuna, span) / span
	i := ni % 9
	transcurrido := ciclo[i].anios * frac
	inicio := nacimiento.Add(-time.Duration(transcurrido * anioSideral * 24 * float64(time.Hour)))

	var out []Periodo
	ahora := time.Now()
	for k := 0; k < hasta; k++ {
		c := ciclo[(i+k)%9]
		fin := inicio.Add(time.Duration(c.anios * anioSideral * 24 * float64(time.Hour)))
		p := Periodo{Senor: c.senor, Desde: inicio.Format("2006-01-02"),
			Hasta: fin.Format("2006-01-02"), Anios: c.anios,
			Actual: ahora.After(inicio) && ahora.Before(fin)}
		// bhuktis
		j := 0
		for x, y := range ciclo {
			if y.senor == c.senor {
				j = x
			}
		}
		bi := inicio
		for b := 0; b < 9; b++ {
			sc := ciclo[(j+b)%9]
			dur := c.anios * sc.anios / 120
			bf := bi.Add(time.Duration(dur * anioSideral * 24 * float64(time.Hour)))
			if bf.After(nacimiento) {
				p.Sub = append(p.Sub, Periodo{Senor: sc.senor, Desde: bi.Format("2006-01-02"),
					Hasta: bf.Format("2006-01-02"), Anios: math.Round(dur*100) / 100,
					Actual: ahora.After(bi) && ahora.Before(bf)})
			}
			bi = bf
		}
		out = append(out, p)
		inicio = fin
	}
	return out
}

// ─────────────── carta completa ───────────────

// Opciones son las decisiones que no cambian la astronomía pero sí la salida.
type Opciones struct {
	NodoVerdadero bool   // nodo verdadero en vez del medio
	Lang          string // "es" o "en"; en qué idioma salen los yogas
}

// Calcular usa el nodo medio y español, que son los valores por defecto.
func Calcular(anio, mes, dia, hh, mm int, tz, lat, lonGeo float64) Carta {
	return CalcularOpts(anio, mes, dia, hh, mm, tz, lat, lonGeo, Opciones{})
}

// CalcularCon se conserva para quien solo quiera cambiar el nodo.
func CalcularCon(anio, mes, dia, hh, mm int, tz, lat, lonGeo float64, nodoVerdadero bool) Carta {
	return CalcularOpts(anio, mes, dia, hh, mm, tz, lat, lonGeo,
		Opciones{NodoVerdadero: nodoVerdadero})
}

func CalcularOpts(anio, mes, dia, hh, mm int, tz, lat, lonGeo float64, o Opciones) Carta {
	nodoVerdadero := o.NodoVerdadero
	base := efem.CalcularCon(anio, mes, dia, hh, mm, tz, lat, lonGeo, nodoVerdadero)
	jd := base.JD
	ayan := Ayanamsa(jd)

	c := Carta{JD: jd, UT: base.UT, Ayanamsa: ayan, NodoVerdad: nodoVerdadero,
		Lagna: Sidereo(base.Asc, jd), MC: Sidereo(base.MC, jd),
		Vargas: map[string][]Graha{}, Karakas: map[string]string{}}
	c.LagnaRasi = int(c.Lagna / 30)
	c.LagnaNak, c.LagnaPada, _ = Nakshatra(c.Lagna)
	c.LagnaPos = Formato(c.Lagna)
	c.SenorLagna = SenorRasi[c.LagnaRasi]

	pos := map[string]float64{}
	for _, b := range base.Cuerpos {
		var n string
		switch b.Nombre {
		case "Nodo Norte":
			n = "Rāhu"
		case "Nodo Sur":
			n = "Ketu"
		case "Urano", "Neptuno", "Plutón":
			continue
		default:
			n = b.Nombre
		}
		l := Sidereo(b.Lon, jd)
		pos[n] = l
		nk, pd, sn := Nakshatra(l)
		c.Grahas = append(c.Grahas, Graha{Nombre: n, Glifo: Glifo[n], Lon: l,
			Rasi: Rasis[int(l/30)], RasiIdx: int(l / 30), Grado: math.Mod(l, 30),
			Nak: nk, Pada: pd, SenorNak: sn,
			Bhava: ((int(l/30)-c.LagnaRasi)%12+12)%12 + 1,
			Dignidad: Dignidad(n, l), Retro: b.Retro || n == "Rāhu" || n == "Ketu",
			Gandanta: Gandanta(l), Posicion: Formato(l)})
	}
	// estado fino: mūlatrikoṇa, combustión y dig-bala
	for i := range c.Grahas {
		g := &c.Grahas[i]
		g.Mula = EnMulatrikona(g.Nombre, g.Lon)
		if g.Nombre != "Sol" {
			g.Combusto, g.DelSol = Combusto(g.Nombre, g.Lon, pos["Sol"], g.Retro)
		}
		g.DigBala = TieneDigBala(g.Nombre, g.Bhava)
		if g.Mula && g.Dignidad == "signo propio" {
			g.Dignidad = "mūlatrikoṇa"
		}
	}
	sort.Slice(c.Grahas, func(i, j int) bool {
		oi, oj := 99, 99
		for k, g := range Grahas {
			if g == c.Grahas[i].Nombre {
				oi = k
			}
			if g == c.Grahas[j].Nombre {
				oj = k
			}
		}
		return oi < oj
	})

	// bhāvas de signo entero
	for h := 1; h <= 12; h++ {
		r := (c.LagnaRasi + h - 1) % 12
		bh := Bhava{Numero: h, Rasi: Rasis[r], RasiIdx: r, Senor: SenorRasi[r]}
		bh.SenorEn = ((int(pos[SenorRasi[r]]/30)-c.LagnaRasi)%12+12)%12 + 1
		for _, g := range c.Grahas {
			if g.RasiIdx == r {
				bh.Ocupan = append(bh.Ocupan, g.Nombre)
			}
			for _, d := range Drishti(g.Nombre) {
				if (g.RasiIdx+d-1)%12 == r {
					bh.Aspectan = append(bh.Aspectan, g.Nombre)
				}
			}
		}
		c.Bhavas = append(c.Bhavas, bh)
	}

	// vargas
	for _, v := range VargasUsadas {
		var lista []Graha
		lagV := Varga(c.Lagna, v.N)
		for _, g := range c.Grahas {
			s := Varga(g.Lon, v.N)
			lista = append(lista, Graha{Nombre: g.Nombre, Glifo: g.Glifo, Rasi: Rasis[s],
				RasiIdx: s, Bhava: ((s-lagV)%12+12)%12 + 1})
		}
		lista = append(lista, Graha{Nombre: "Lagna", Rasi: Rasis[lagV], RasiIdx: lagV, Bhava: 1})
		c.Vargas[v.Clave] = lista
	}

	// kārakas de Jaimini: los siete grahas ordenados por grado dentro de su signo
	nombres := []string{"AK", "AmK", "BK", "MK", "PiK", "PK", "GK"}
	tipo := []struct {
		n string
		g float64
	}{}
	for _, g := range c.Grahas {
		if g.Nombre == "Rāhu" || g.Nombre == "Ketu" {
			continue
		}
		tipo = append(tipo, struct {
			n string
			g float64
		}{g.Nombre, g.Grado})
	}
	sort.Slice(tipo, func(i, j int) bool { return tipo[i].g > tipo[j].g })
	for i, t := range tipo {
		if i < len(nombres) {
			c.Karakas[nombres[i]] = t.n
		}
	}

	// daśās
	nac := time.Date(anio, time.Month(mes), dia, hh, mm, 0, 0,
		time.FixedZone("local", int(tz*3600)))
	c.Dasas = Vimsottari(pos["Luna"], nac, 9)
	if ak := c.Karakas["AK"]; ak != "" {
		for _, g := range c.Grahas {
			if g.Nombre == ak {
				c.Karakamsa = Rasis[Varga(g.Lon, 9)]
			}
		}
	}
	c.Yogas = detectarYogas(c, pos, o.Lang)
	c.Gocara = Transitos(c.LagnaRasi, int(pos["Luna"]/30), time.Now())

	// ── pañcāṅga, el calendario del día ──
	c.Pancanga = CalcPancanga(jd, pos["Sol"], pos["Luna"])

	// ── aṣṭakavarga: los bindus por rāśi ──
	rasiDe := map[string]int{"Lagna": c.LagnaRasi}
	for g, l := range pos {
		rasiDe[g] = int(l / 30)
	}
	c.Ashtaka = CalcAshtakavarga(rasiDe)

	// ── arudha padas: la imagen de cada bhāva ──
	c.Arudhas = CalcArudhas(c.Bhavas, rasiDe)

	// ── lagnas especiales: se cuentan desde el amanecer ──
	// jd0 es el día juliano a las 0h UT del día del nacimiento.
	jd0 := math.Floor(jd-0.5) + 0.5
	orto, _, ocaso, hay := efem.Orto(jd0, lat, lonGeo)
	horasUT := (jd - jd0) * 24
	if hay {
		desde := horasUT - orto
		if desde < 0 {
			// nació antes del amanecer: cuenta desde el orto del día anterior
			o2, _, _, h2 := efem.Orto(jd0-1, lat, lonGeo)
			if h2 {
				orto, desde = o2, horasUT+24-o2
			}
		}
		solOrto := Sidereo(efem.Sol(jd0+orto/24), jd)
		c.LagnasEsp = CalcLagnasEsp(solOrto, desde, orto, true)
	}

	// ── ṣaḍbala ──
	entradas := map[string]EntradaBala{}
	for _, b := range base.Cuerpos {
		n := b.Nombre
		if n == "Urano" || n == "Neptuno" || n == "Plutón" ||
			n == "Nodo Norte" || n == "Nodo Sur" {
			continue
		}
		l := Sidereo(b.Lon, jd)
		entradas[n] = EntradaBala{Lon: l, Vel: b.Vel,
			Bhava: ((int(l/30)-c.LagnaRasi)%12+12)%12 + 1}
	}
	esDeDia := hay && horasUT >= orto && horasUT < ocaso
	// el señor de la hora gira en el mismo orden que los días de la semana
	senorHora := "Sol"
	if hay {
		// Horas planetarias, contadas de reloj desde el amanecer. La primera
		// es del señor del día y las siguientes bajan por el orden caldeo.
		h := int(horasUT-orto+24) % 24
		senorHora = ordenHoras[(indiceCaldeo(c.Pancanga.SenorVara)+h)%7]
	}
	c.Shadbala = CalcShadbala(entradas, c.Lagna, c.MC, esDeDia,
		c.Pancanga.TithiNum, c.Pancanga.SenorVara, senorHora)

	return c
}

// El orden caldeo, de más lento a más rápido, es el que rige las horas
// planetarias. Cada hora desde el amanecer avanza un puesto.
var ordenHoras = [7]string{"Saturno", "Júpiter", "Marte", "Sol", "Venus", "Mercurio", "Luna"}

func indiceCaldeo(g string) int {
	for i, n := range ordenHoras {
		if n == g {
			return i
		}
	}
	return 0
}

// detectarYogas busca las combinaciones que se pueden afirmar con seguridad
// desde las posiciones. No es un catálogo: solo lo comprobable.
func detectarYogas(c Carta, pos map[string]float64, lang string) []string {
	f := yogaTextos(lang)
	gn := nombreGraha(lang)
	var y []string
	rLuna := int(pos["Luna"] / 30)

	// Gaja-kesari: Guru en kendra DESDE LA LUNA, no desde el Lagna
	if Kendra(int(pos["Júpiter"]/30), rLuna) {
		y = append(y, f["gaja"])
	}

	// Pañca-mahāpuruṣa: exaltado o en signo propio Y en kendra desde el Lagna
	for _, g := range c.Grahas {
		nom, ok := mahapurusha[g.Nombre]
		if !ok {
			continue
		}
		fuerte := g.Dignidad == "exaltado" || g.Dignidad == "signo propio" || g.Dignidad == "mūlatrikoṇa"
		if fuerte && (g.Bhava == 1 || g.Bhava == 4 || g.Bhava == 7 || g.Bhava == 10) {
			y = append(y, fmt.Sprintf(f["mahapurusha"], nom, gn(g.Nombre),
				f["dig_"+g.Dignidad], g.Bhava))
		}
	}

	// Rāja-yoga: relación entre señor de kendra y señor de trikona
	senorDe := func(h int) string { return SenorRasi[(c.LagnaRasi+h-1)%12] }
	kendras := []int{1, 4, 7, 10}
	trikonas := []int{1, 5, 9}
	vistos := map[string]bool{}
	for _, k := range kendras {
		for _, tr := range trikonas {
			if k == tr {
				continue
			}
			a, b := senorDe(k), senorDe(tr)
			if a == b {
				clave := "solo:" + a
				if !vistos[clave] {
					vistos[clave] = true
					y = append(y, fmt.Sprintf(f["raja_solo"], gn(a), k, tr))
				}
				continue
			}
			ra, rb := int(pos[a]/30), int(pos[b]/30)
			clave := a + "-" + b
			if vistos[clave] || vistos[b+"-"+a] {
				continue
			}
			switch {
			case ra == rb:
				vistos[clave] = true
				y = append(y, fmt.Sprintf(f["raja_juntos"], k, tr, gn(a), gn(b), Rasis[ra]))
			case seMiran(a, b, ra, rb):
				vistos[clave] = true
				y = append(y, fmt.Sprintf(f["raja_miran"], k, tr, gn(a), gn(b)))
			}
		}
	}

	// Nīca-bhaṅga: debilidad cancelada
	for _, g := range c.Grahas {
		if g.Dignidad != "debilitado" {
			continue
		}
		sd := SenorRasi[g.RasiIdx]
		if Kendra(int(pos[sd]/30), c.LagnaRasi) || Kendra(int(pos[sd]/30), rLuna) {
			y = append(y, fmt.Sprintf(f["nica"], gn(g.Nombre), gn(sd)))
		}
	}

	// Kemadruma: Luna sin grahas en las casas contiguas
	solo := true
	for _, g := range c.Grahas {
		if g.Nombre == "Luna" || g.Nombre == "Rāhu" || g.Nombre == "Ketu" {
			continue
		}
		d := ((g.RasiIdx-rLuna)%12 + 12) % 12
		if d == 0 || d == 1 || d == 11 {
			solo = false
		}
	}
	if solo {
		y = append(y, f["kemadruma"])
	}
	return y
}

// seMiran indica si dos grahas se aspectan mutuamente por dṛṣṭi de signo.
func seMiran(a, b string, ra, rb int) bool {
	mira := func(g string, desde, hasta int) bool {
		for _, d := range Drishti(g) {
			if (desde+d-1)%12 == hasta {
				return true
			}
		}
		return false
	}
	return mira(a, ra, rb) && mira(b, rb, ra)
}

func itoa(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}


// yogaTextos son las plantillas de cada yoga en el idioma que se pida. Igual
// que en el motor de lectura, cada idioma trae la frase entera: no se arma
// pegando trozos, porque el orden de las palabras no coincide.
func yogaTextos(lang string) map[string]string {
	if lang == "en" {
		return map[string]string{
			"gaja": "Gaja-kesari — Jupiter in a kendra from the Moon: intelligence, good name, protection",
			"mahapurusha": "%s yoga — %s %s in a kendra (house %d): one of the five mahāpuruṣa",
			"raja_solo":   "Rāja-yoga — %s rules a kendra (%d) and a trikona (%d) at once",
			"raja_juntos": "Rāja-yoga — the lords of the %dth and the %dth (%s and %s) sit together in %s",
			"raja_miran":  "Rāja-yoga — the lords of the %dth and the %dth (%s and %s) look at each other",
			"nica": "Nīca-bhaṅga — %s is debilitated, but %s, the lord of its sign, occupies a kendra: the debilitation is cancelled",
			"kemadruma": "Kemadruma — the Moon has no grahas beside it or with it. It reads as mental solitude, but it is easily cancelled: do not alarm anyone with this",
			"dig_exaltado": "exalted", "dig_signo propio": "in its own sign",
			"dig_mūlatrikoṇa": "in mūlatrikoṇa", "dig_debilitado": "debilitated", "dig_—": "",
		}
	}
	return map[string]string{
		"gaja": "Gaja-kesari — Júpiter en kendra desde la Luna: inteligencia, buen nombre, protección",
		"mahapurusha": "%s yoga — %s %s en kendra (casa %d): uno de los cinco mahāpuruṣa",
		"raja_solo":   "Rāja-yoga — %s rige a la vez un kendra (%d) y un trikona (%d)",
		"raja_juntos": "Rāja-yoga — los señores de la %d y la %d (%s y %s) están juntos en %s",
		"raja_miran":  "Rāja-yoga — los señores de la %d y la %d (%s y %s) se miran mutuamente",
		"nica": "Nīca-bhaṅga — %s está debilitado pero %s, señor de su signo, ocupa un kendra: la debilidad se cancela",
		"kemadruma": "Kemadruma — la Luna no tiene grahas a los lados ni consigo. Se lee como soledad mental, pero se cancela con facilidad: no alarmes con esto",
		"dig_exaltado": "exaltado", "dig_signo propio": "en signo propio",
		"dig_mūlatrikoṇa": "en mūlatrikoṇa", "dig_debilitado": "debilitado", "dig_—": "",
	}
}

// nombreGraha traduce el nombre del graha al salir. Dentro del motor siempre
// van en español, que es como los produce el cálculo.
func nombreGraha(lang string) func(string) string {
	if lang != "en" {
		return func(g string) string { return g }
	}
	return func(g string) string {
		if v, ok := grahaEn[g]; ok {
			return v
		}
		return g
	}
}
