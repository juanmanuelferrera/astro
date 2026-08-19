package efem

import (
	"math"
	"time"
)

var Signos = [12]string{"Aries", "Tauro", "Géminis", "Cáncer", "Leo", "Virgo",
	"Libra", "Escorpio", "Sagitario", "Capricornio", "Acuario", "Piscis"}
var GlifoSigno = [12]string{"♈", "♉", "♊", "♋", "♌", "♍", "♎", "♏", "♐", "♑", "♒", "♓"}
var RegenteSigno = [12]string{"Marte", "Venus", "Mercurio", "Luna", "Sol", "Mercurio",
	"Venus", "Marte", "Júpiter", "Saturno", "Saturno", "Júpiter"}

var Glifo = map[string]string{"Sol": "☉", "Luna": "☽", "Mercurio": "☿", "Venus": "♀",
	"Marte": "♂", "Júpiter": "♃", "Saturno": "♄", "Urano": "♅", "Neptuno": "♆",
	"Plutón": "♇", "Nodo Norte": "☊", "Nodo Sur": "☋"}

var Orden = []string{"Sol", "Luna", "Mercurio", "Venus", "Marte", "Júpiter",
	"Saturno", "Urano", "Neptuno", "Plutón", "Nodo Norte", "Nodo Sur"}

var domicilio = map[string][]int{"Sol": {4}, "Luna": {3}, "Mercurio": {2, 5},
	"Venus": {1, 6}, "Marte": {0, 7}, "Júpiter": {8, 11}, "Saturno": {9, 10}}
var exaltacion = map[string]int{"Sol": 0, "Luna": 1, "Mercurio": 5, "Venus": 11,
	"Marte": 9, "Júpiter": 3, "Saturno": 6}

type Aspecto struct {
	A      string `json:"a"`
	B      string `json:"b"`
	Nombre string `json:"nombre"`
	Glifo  string `json:"glifo"`
	Exacto              float64 `json:"exacto"`
	Orbe                float64 `json:"orbe"`
}

type Cuerpo struct {
	Nombre   string  `json:"nombre"`
	Glifo    string  `json:"glifo"`
	Lon      float64 `json:"lon"`
	Signo    string  `json:"signo"`
	SignoIdx int     `json:"signoIdx"`
	Grado    float64 `json:"grado"`
	GlifoSig string  `json:"glifoSig"`
	Retro    bool    `json:"retro"`
	Vel      float64 `json:"vel"` // grados por día; negativa si retrograda
	Dignidad string  `json:"dignidad"`
	CasaP    int     `json:"casaP"`
	CasaI    int     `json:"casaI"`
}

type Carta struct {
	Fecha, Hora, Lugar string     `json:"-"`
	Etiqueta           string     `json:"etiqueta"`
	JD                 float64    `json:"jd"`
	UT                 string     `json:"ut"`
	TSG                float64    `json:"tsg"`
	TSL                float64    `json:"tsl"`
	Oblicuidad         float64    `json:"oblicuidad"`
	Asc, MC, Dsc, IC   float64    `json:"-"`
	Angulos            [4]float64 `json:"angulos"`
	CuspP              [12]float64 `json:"cuspP"`
	CuspI              [12]float64 `json:"cuspI"`
	PlacidusOK         bool        `json:"placidusOK"`
	Cuerpos            []Cuerpo    `json:"cuerpos"`
	Aspectos           []Aspecto   `json:"aspectos"`
	Regentes           [12]string  `json:"regentes"`
	RegenteEn          [12]int     `json:"regenteEn"`
}

var tablaAspectos = []struct {
	nombre, glifo string
	angulo, orbe  float64
}{
	{"conjunción", "☌", 0, 8}, {"oposición", "☍", 180, 8}, {"trígono", "△", 120, 7},
	{"cuadratura", "□", 90, 7}, {"sextil", "✶", 60, 5}, {"semicuadratura", "∠", 45, 2},
	{"sesquicuadratura", "⚼", 135, 2}, {"quincuncio", "⚻", 150, 3}, {"semisextil", "⚺", 30, 2},
}

func dignidad(nombre string, signo int) string {
	if d, ok := domicilio[nombre]; ok {
		for _, s := range d {
			if s == signo {
				return "domicilio"
			}
			if (s+6)%12 == signo {
				return "exilio"
			}
		}
	}
	if e, ok := exaltacion[nombre]; ok {
		if e == signo {
			return "exaltación"
		}
		if (e+6)%12 == signo {
			return "caída"
		}
	}
	return "—"
}

func casaDe(lon float64, c [12]float64) int {
	for i := 0; i < 12; i++ {
		a, b := c[i], c[(i+1)%12]
		if a <= b {
			if lon >= a && lon < b {
				return i + 1
			}
		} else if lon >= a || lon < b {
			return i + 1
		}
	}
	return 1
}

// Calcular levanta la carta completa. tz en horas (este positivo), lonGeo este positivo.
// Calcular usa el nodo lunar medio, que es lo tradicional en occidental.
func Calcular(anio, mes, dia, hh, mm int, tz, lat, lonGeo float64) Carta {
	return CalcularCon(anio, mes, dia, hh, mm, tz, lat, lonGeo, false)
}

// CalcularCon permite pedir el nodo verdadero en vez del medio. En jyotiṣa hay
// escuelas para las dos cosas y la diferencia llega a grado y medio, que basta
// para mover a Rāhu de pada y a veces de nakṣatra. Las daśās no cambian: esas
// cuelgan del nakṣatra de la Luna, no del de Rāhu.
func CalcularCon(anio, mes, dia, hh, mm int, tz, lat, lonGeo float64, nodoVerdadero bool) Carta {
	horaLocal := float64(hh) + float64(mm)/60
	horaUT := horaLocal - tz
	d := float64(dia)
	if horaUT < 0 {
		horaUT += 24
		d--
	} else if horaUT >= 24 {
		horaUT -= 24
		d++
	}
	jd := DiaJuliano(anio, mes, d, horaUT)
	eps := Oblicuidad(jd)
	tsg := TiempoSidereoGreenwich(jd)
	tsl := norm360(tsg + lonGeo)

	asc := Ascendente(tsl, lat, eps)
	mc := MedioCielo(tsl, eps)
	cp, ok := Cuspides(tsl, lat, eps)
	ci := CuspidesIguales(asc)

	c := Carta{JD: jd, TSG: tsg, TSL: tsl, Oblicuidad: eps,
		Asc: asc, MC: mc, Dsc: norm360(asc + 180), IC: norm360(mc + 180),
		Angulos: [4]float64{asc, norm360(asc + 180), mc, norm360(mc + 180)},
		CuspP:   cp, CuspI: ci, PlacidusOK: ok}
	c.UT = time.Date(anio, time.Month(mes), int(d), int(horaUT),
		int((horaUT-math.Floor(horaUT))*60), 0, 0, time.UTC).Format("2006-01-02 15:04 UT")

	pos := map[string]float64{}
	for _, n := range Orden {
		var l float64
		switch n {
		case "Sol":
			l = Sol(jd)
		case "Luna":
			l = Luna(jd)
		case "Nodo Norte":
			l = nodo(jd, nodoVerdadero)
		case "Nodo Sur":
			l = norm360(nodo(jd, nodoVerdadero) + 180)
		default:
			l = Planeta(n, jd)
		}
		pos[n] = l
		si := int(l / 30)
		// La velocidad se saca por diferencia central, medio día a cada lado.
		// Sirve para saber si retrograda y, en jyotiṣa, para el cheṣṭā bala.
		var vel float64
		switch n {
		case "Sol":
			vel = math.Mod(Sol(jd+0.5)-Sol(jd-0.5)+540, 360) - 180
		case "Luna":
			vel = math.Mod(Luna(jd+0.5)-Luna(jd-0.5)+540, 360) - 180
		case "Nodo Norte", "Nodo Sur":
			vel = math.Mod(nodo(jd+0.5, nodoVerdadero)-nodo(jd-0.5, nodoVerdadero)+540, 360) - 180
		default:
			vel = math.Mod(Planeta(n, jd+0.5)-Planeta(n, jd-0.5)+540, 360) - 180
		}
		retro := vel < 0
		c.Cuerpos = append(c.Cuerpos, Cuerpo{Nombre: n, Glifo: Glifo[n], Lon: l,
			Signo: Signos[si], SignoIdx: si, Grado: l - float64(si)*30,
			GlifoSig: GlifoSigno[si], Retro: retro, Vel: vel, Dignidad: dignidad(n, si),
			CasaP: casaDe(l, cp), CasaI: casaDe(l, ci)})
	}

	// aspectos entre los diez planetas
	pl := Orden[:10]
	for i := 0; i < len(pl); i++ {
		for j := i + 1; j < len(pl); j++ {
			sep := math.Abs(pos[pl[i]] - pos[pl[j]])
			if sep > 180 {
				sep = 360 - sep
			}
			for _, a := range tablaAspectos {
				if dv := math.Abs(sep - a.angulo); dv <= a.orbe {
					c.Aspectos = append(c.Aspectos, Aspecto{A: pl[i], B: pl[j],
						Nombre: a.nombre, Glifo: a.glifo, Exacto: a.angulo,
						Orbe: math.Round(dv*100) / 100})
					break
				}
			}
		}
	}
	for i := 0; i < len(c.Aspectos); i++ {
		for j := i + 1; j < len(c.Aspectos); j++ {
			if c.Aspectos[j].Orbe < c.Aspectos[i].Orbe {
				c.Aspectos[i], c.Aspectos[j] = c.Aspectos[j], c.Aspectos[i]
			}
		}
	}
	// regentes de casa y dónde están alojados
	for i := 0; i < 12; i++ {
		r := RegenteSigno[int(cp[i]/30)]
		c.Regentes[i] = r
		c.RegenteEn[i] = casaDe(pos[r], cp)
	}
	return c
}

// nodo elige entre el nodo medio y el verdadero.
func nodo(jd float64, verdadero bool) float64 {
	if verdadero {
		return NodoLunarVerdadero(jd)
	}
	return NodoLunarMedio(jd)
}
