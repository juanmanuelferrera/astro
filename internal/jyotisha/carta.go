package jyotisha

import (
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

func Calcular(anio, mes, dia, hh, mm int, tz, lat, lonGeo float64) Carta {
	base := efem.Calcular(anio, mes, dia, hh, mm, tz, lat, lonGeo)
	jd := base.JD
	ayan := Ayanamsa(jd)

	c := Carta{JD: jd, UT: base.UT, Ayanamsa: ayan,
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
	c.Yogas = detectarYogas(c, pos)
	c.Gocara = Transitos(c.LagnaRasi, int(pos["Luna"]/30), time.Now())
	return c
}

// detectarYogas busca las combinaciones que se pueden afirmar con seguridad
// desde las posiciones. No es un catálogo: solo lo comprobable.
func detectarYogas(c Carta, pos map[string]float64) []string {
	var y []string
	rLuna := int(pos["Luna"] / 30)

	// Gaja-kesari: Guru en kendra DESDE LA LUNA, no desde el Lagna
	if Kendra(int(pos["Júpiter"]/30), rLuna) {
		y = append(y, "Gaja-kesari — Júpiter en kendra desde la Luna: inteligencia, buen nombre, protección")
	}

	// Pañca-mahāpuruṣa: exaltado o en signo propio Y en kendra desde el Lagna
	for _, g := range c.Grahas {
		nom, ok := mahapurusha[g.Nombre]
		if !ok {
			continue
		}
		fuerte := g.Dignidad == "exaltado" || g.Dignidad == "signo propio" || g.Dignidad == "mūlatrikoṇa"
		if fuerte && (g.Bhava == 1 || g.Bhava == 4 || g.Bhava == 7 || g.Bhava == 10) {
			y = append(y, nom+" yoga — "+g.Nombre+" "+g.Dignidad+" en kendra (casa "+
				itoa(g.Bhava)+"): uno de los cinco mahāpuruṣa")
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
					y = append(y, "Rāja-yoga — "+a+" rige a la vez un kendra ("+itoa(k)+
						") y un trikona ("+itoa(tr)+")")
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
				y = append(y, "Rāja-yoga — los señores de la "+itoa(k)+" y la "+itoa(tr)+
					" ("+a+" y "+b+") están juntos en "+Rasis[ra])
			case seMiran(a, b, ra, rb):
				vistos[clave] = true
				y = append(y, "Rāja-yoga — los señores de la "+itoa(k)+" y la "+itoa(tr)+
					" ("+a+" y "+b+") se miran mutuamente")
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
			y = append(y, "Nīca-bhaṅga — "+g.Nombre+" está debilitado pero "+sd+
				", señor de su signo, ocupa un kendra: la debilidad se cancela")
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
		y = append(y, "Kemadruma — la Luna no tiene grahas a los lados ni consigo. "+
			"Se lee como soledad mental, pero se cancela con facilidad: no alarmes con esto")
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

