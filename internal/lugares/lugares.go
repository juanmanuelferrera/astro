// Package lugares — ciudades con coordenadas y zona horaria IANA.
// La resolución del huso la hace time/tzdata, embebida en el binario: incluye
// todas las reglas históricas (España a hora de Berlín en 1940, cambios de
// horario de verano por año y por país, etc.).
package lugares

import (
	"bufio"
	"bytes"
	"compress/gzip"
	_ "embed"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata" // embebe la base IANA: sin esto no hay husos históricos
)

type Lugar struct {
	Nombre   string  `json:"nombre"`
	Region   string  `json:"region"`
	Pais     string  `json:"pais"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Zona     string  `json:"zona"`
}

//go:embed datos/lugares.tsv.gz
var comprimido []byte

// Ciudades se carga una sola vez, ordenada por población descendente:
// ese orden es el ranking de resultados.
var Ciudades []Lugar

func init() {
	zr, err := gzip.NewReader(bytes.NewReader(comprimido))
	if err != nil {
		return
	}
	defer zr.Close()
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		c := strings.Split(sc.Text(), "\t")
		if len(c) < 6 {
			continue
		}
		lat, _ := strconv.ParseFloat(c[3], 64)
		lon, _ := strconv.ParseFloat(c[4], 64)
		Ciudades = append(Ciudades, Lugar{c[0], c[1], c[2], lat, lon, c[5]})
	}
}

func sinTildes(s string) string {
	r := strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u",
		"ü", "u", "ñ", "n", "à", "a", "è", "e", "ç", "c")
	return r.Replace(strings.ToLower(s))
}

// Buscar devuelve hasta n coincidencias. Tres niveles de prioridad —exacto,
// por prefijo, por subcadena— y dentro de cada nivel se conserva el orden del
// fichero, que está ordenado por población.
func Buscar(q string, n int) []Lugar {
	q = sinTildes(strings.TrimSpace(q))
	if len(q) < 2 {
		return []Lugar{}
	}
	var exacto, pref, sub []Lugar
	for _, l := range Ciudades {
		nom := sinTildes(l.Nombre)
		switch {
		case nom == q:
			exacto = append(exacto, l)
		case strings.HasPrefix(nom, q):
			pref = append(pref, l)
		case strings.Contains(nom, q):
			sub = append(sub, l)
		}
		if len(exacto)+len(pref) >= n*3 {
			break
		}
	}
	r := append(append(exacto, pref...), sub...)
	if len(r) > n {
		r = r[:n]
	}
	if r == nil {
		r = []Lugar{}
	}
	return r
}

type Cambio struct {
	Fecha  string  `json:"fecha"`
	De     float64 `json:"de"`
	A      float64 `json:"a"`
	Nombre string  `json:"nombre"`
	Motivo string  `json:"motivo"`
}

type Historia struct {
	Zona       string   `json:"zona"`
	Offset     float64  `json:"offset"`
	Abrev      string   `json:"abrev"`
	Verano     bool     `json:"verano"`
	Estandar   float64  `json:"estandar"`
	DelAnio    []Cambio `json:"delAnio"`
	Historicos []Cambio `json:"historicos"`
	Solar      float64  `json:"solar"`
	Desfase    float64  `json:"desfase"`
}

func offsetEn(loc *time.Location, t time.Time) (float64, string) {
	n, o := t.In(loc).Zone()
	return float64(o) / 3600, n
}

// HistoriaHuso explica de dónde sale el desfase: los saltos de horario de verano
// de ese año y los cambios estructurales del país a lo largo del siglo.
func HistoriaHuso(zona string, anio, mes, dia, hh, mm int, lonGeo float64) (Historia, error) {
	loc, err := time.LoadLocation(zona)
	if err != nil {
		return Historia{}, err
	}
	H := Historia{Zona: zona, DelAnio: []Cambio{}, Historicos: []Cambio{}}
	nacim := time.Date(anio, time.Month(mes), dia, hh, mm, 0, 0, loc)
	H.Offset, H.Abrev = offsetEn(loc, nacim)

	// transiciones dentro del año de nacimiento, día a día
	prev, prevN := offsetEn(loc, time.Date(anio, 1, 1, 12, 0, 0, 0, time.UTC))
	minOff := prev
	for d := 1; d <= 366; d++ {
		t := time.Date(anio, 1, 1, 12, 0, 0, 0, time.UTC).AddDate(0, 0, d)
		if t.Year() != anio {
			break
		}
		o, n := offsetEn(loc, t)
		if o < minOff {
			minOff = o
		}
		if o != prev {
			motivo := "empieza el horario de verano"
			if o < prev {
				motivo = "termina el horario de verano"
			}
			H.DelAnio = append(H.DelAnio, Cambio{Fecha: t.Format("2 Jan"), De: prev, A: o,
				Nombre: n, Motivo: motivo})
			prev, prevN = o, n
		}
	}
	_ = prevN
	H.Estandar = minOff
	H.Verano = H.Offset > minOff

	// cambios estructurales: se compara el desfase de enero año a año
	pe, _ := offsetEn(loc, time.Date(1890, 1, 15, 12, 0, 0, 0, time.UTC))
	for y := 1891; y <= 2050; y++ {
		o, n := offsetEn(loc, time.Date(y, 1, 15, 12, 0, 0, 0, time.UTC))
		if o != pe {
			H.Historicos = append(H.Historicos, Cambio{Fecha: strconvI(y), De: pe, A: o,
				Nombre: n, Motivo: "cambio de huso estándar del país"})
			pe = o
		}
	}
	// hora solar frente a hora del reloj
	H.Solar = lonGeo / 15
	H.Desfase = H.Offset - H.Solar
	return H, nil
}

func strconvI(y int) string {
	s := ""
	for y > 0 {
		s = string(rune('0'+y%10)) + s
		y /= 10
	}
	return s
}

// Huso resuelve el desfase real de un lugar en un instante concreto.
// Devuelve horas, abreviatura de la zona y si había horario de verano.
func Huso(zona string, anio, mes, dia, hh, mm int) (float64, string, bool, error) {
	loc, err := time.LoadLocation(zona)
	if err != nil {
		return 0, "", false, err
	}
	t := time.Date(anio, time.Month(mes), dia, hh, mm, 0, 0, loc)
	nom, off := t.Zone()
	_, oe := time.Date(anio, 1, 15, 12, 0, 0, 0, loc).Zone()
	_, oj := time.Date(anio, 7, 15, 12, 0, 0, 0, loc).Zone()
	estandar := oe
	if oj < oe {
		estandar = oj
	}
	return float64(off) / 3600, nom, off > estandar, nil
}
