package occidental

import (
	"fmt"
	"math"
	"sort"
	"time"

	"astro/internal/efem"
)

// Predicción occidental: tránsitos, progresiones secundarias y revolución
// solar. Es lo que enseña el módulo 12 del curso, y hasta ahora el programa lo
// explicaba sin hacerlo.
//
// Las tres técnicas NO se leen por separado. Un tránsito solo, sin apoyo, casi
// nunca produce nada visible; lo que marca un periodo es que las tres apunten a
// lo mismo. Eso es la regla de convergencia del módulo 9 aplicada al tiempo, y
// aquí se calcula en lugar de dejarla al ojo: ver Convergencias más abajo.
//
// Y lo que no hace, que importa tanto como lo que hace: no fecha sucesos. Da
// periodos y qué se activa en ellos. El módulo 13 explica por qué.

// Los lentos son los que marcan años. Los rápidos —Luna, Mercurio, Venus—
// marcan días y se dejan fuera a propósito: producen ruido, no periodos.
var lentos = []string{"Marte", "Júpiter", "Saturno", "Urano", "Neptuno", "Plutón"}

// Los orbes de tránsito son estrechos a propósito: en natal un orbe ancho
// matiza, pero en tránsito convierte cualquier semana en un acontecimiento.
//
// Y los menores van más apretados todavía. Con el mismo orbe que los mayores,
// la lista se llena de sesquicuadraturas y semisextiles —que son muchos más— y
// el tránsito que importa queda enterrado debajo.
const (
	orbeMayor = 1.5
	orbeMenor = 0.4
)

func esMayor(a string) bool {
	switch a {
	case "conjunción", "oposición", "trígono", "cuadratura", "sextil":
		return true
	}
	return false
}

func orbeDe(a string) float64 {
	if esMayor(a) {
		return orbeMayor
	}
	return orbeMenor
}

type Transito struct {
	Planeta  string  `json:"planeta"`  // el que transita
	Aspecto  string  `json:"aspecto"`
	Glifo    string  `json:"glifo"`
	Natal    string  `json:"natal"`    // a qué punto natal, ya en el idioma pedido
	interno  string  // el mismo punto sin traducir, para agrupar convergencias
	Orbe     float64 `json:"orbe"`
	Aplica   bool    `json:"aplica"`   // se está formando, no soltando
	Retro    bool    `json:"retro"`
	Casa     int     `json:"casa"`     // casa natal por la que pasa
	Pasadas  int     `json:"pasadas"`  // 1 o 3, según vaya a retrogradar encima
}

type Progresion struct {
	Planeta string  `json:"planeta"`
	SignoIdx int    `json:"signoIdx"` // el nombre del signo lo pone el navegador
	Grado   float64 `json:"grado"`
	Casa    int     `json:"casa"`
	Aspecto string  `json:"aspecto"` // a un punto natal, si lo hay
	Natal   string  `json:"natal"`
	Orbe    float64 `json:"orbe"`
	interno string  // sin traducir, para agrupar convergencias
}

type Revolucion struct {
	Cuando string  `json:"cuando"` // instante exacto del retorno, en UT
	Asc    float64 `json:"asc"`
	AscSig int     `json:"ascSig"`
	MC     float64 `json:"mc"`
	MCSig  int     `json:"mcSig"`
	Casa   int     `json:"casa"` // casa natal donde cae el Ascendente del año
}

type Prediccion struct {
	puntos map[string]float64 // los puntos natales; no sale en el JSON
	Transitos     []Transito   `json:"transitos"`
	Progresiones  []Progresion `json:"progresiones"`
	Revolucion    Revolucion   `json:"revolucion"`
	Convergencias []string     `json:"convergencias"`
	Edad          int          `json:"edad"`
	Nota          string       `json:"nota"`
}

// puntosNatales son los que reciben tránsitos: los diez planetas y los cuatro
// ángulos. Un tránsito al Ascendente cuenta tanto como uno a un planeta.
func puntosNatales(c efem.Carta) map[string]float64 {
	m := map[string]float64{}
	for _, b := range c.Cuerpos {
		if b.Nombre == "Nodo Norte" || b.Nombre == "Nodo Sur" {
			continue
		}
		m[b.Nombre] = b.Lon
	}
	m["Ascendente"] = c.Asc
	m["Medio Cielo"] = c.MC
	m["Descendente"] = math.Mod(c.Asc+180, 360)
	m["Fondo del Cielo"] = math.Mod(c.MC+180, 360)
	return m
}

func sep(a, b float64) float64 {
	d := math.Abs(math.Mod(a-b+360, 360))
	if d > 180 {
		d = 360 - d
	}
	return d
}

// Predecir arma las tres técnicas para una fecha dada.
func Predecir(natal efem.Carta, cuando time.Time, lang string) Prediccion {
	T := es
	if lang == "en" {
		T = en
	}
	var p Prediccion
	puntos := puntosNatales(natal)
	p.puntos = puntos
	// La edad manda en las progresiones: un día por año cumplido.
	edad := cuando.Year() - natal.Anio
	if int(cuando.Month()) < natal.Mes ||
		(int(cuando.Month()) == natal.Mes && cuando.Day() < natal.Dia) {
		edad--
	}
	p.Edad = edad

	// ── tránsitos ──
	hoy := efem.Calcular(cuando.Year(), int(cuando.Month()), cuando.Day(),
		cuando.Hour(), cuando.Minute(), 0, natal.Lat, natal.Lon)
	for _, b := range hoy.Cuerpos {
		if !esLento(b.Nombre) {
			continue
		}
		for nombre, lonNatal := range puntos {
			s := sep(b.Lon, lonNatal)
			for _, a := range efem.TablaAspectos() {
				dv := math.Abs(s - a.Angulo)
				if dv > orbeDe(a.Nombre) {
					continue
				}
				// aplicativo: mañana el orbe será menor
				manana := math.Abs(sep(b.Lon+b.Vel, lonNatal) - a.Angulo)
				t := Transito{Planeta: nombreLlano(b.Nombre, T), Aspecto: a.Nombre, Glifo: a.Glifo,
					Natal: nombreLlano(nombre, T), interno: nombre,
					Orbe: math.Round(dv*100) / 100,
					Aplica: manana < dv, Retro: b.Retro,
					Casa: casaNatalDe(b.Lon, natal), Pasadas: 1}
				// Los lentos retrogradan encima del punto y pasan tres veces:
				// directo, retrógrado y directo. Son tres momentos del mismo
				// proceso, no tres cosas distintas.
				if b.Nombre != "Marte" && math.Abs(b.Vel) < 0.25 {
					t.Pasadas = 3
				}
				p.Transitos = append(p.Transitos, t)
				break
			}
		}
	}
	ordenarPorOrbe(p.Transitos)

	// ── progresiones secundarias: un día después del nacimiento por año ──
	prog := efem.CalcularJD(natal.JD+float64(edad), natal.Lat, natal.Lon)
	for _, b := range prog.Cuerpos {
		if b.Nombre != "Sol" && b.Nombre != "Luna" && b.Nombre != "Venus" && b.Nombre != "Marte" {
			continue
		}
		pr := Progresion{Planeta: nombreLlano(b.Nombre, T), SignoIdx: b.SignoIdx, Grado: b.Grado,
			Casa: casaNatalDe(b.Lon, natal)}
		mejor := 99.0
		for nombre, lonNatal := range puntos {
			s := sep(b.Lon, lonNatal)
			for _, a := range efem.TablaAspectos() {
				if dv := math.Abs(s - a.Angulo); dv <= 1 && dv < mejor {
					mejor = dv
					pr.Aspecto, pr.Orbe = a.Nombre, math.Round(dv*100)/100
					pr.Natal, pr.interno = nombreLlano(nombre, T), nombre
				}
			}
		}
		p.Progresiones = append(p.Progresiones, pr)
	}

	// ── revolución solar ──
	p.Revolucion = revolucionSolar(natal, cuando.Year())

	// ── convergencias ──
	p.Convergencias = convergencias(p, T)
	p.Nota = T.f["predNota"]
	return p
}

func esLento(n string) bool {
	for _, x := range lentos {
		if x == n {
			return true
		}
	}
	return false
}

func casaNatalDe(lon float64, c efem.Carta) int {
	for i := 0; i < 12; i++ {
		a, b := c.CuspP[i], c.CuspP[(i+1)%12]
		if a < b {
			if lon >= a && lon < b {
				return i + 1
			}
		} else if lon >= a || lon < b { // la casa cruza el 0° de Aries
			return i + 1
		}
	}
	return 0
}

func ordenarPorOrbe(t []Transito) {
	for i := range t {
		for j := i + 1; j < len(t); j++ {
			if t[j].Orbe < t[i].Orbe {
				t[i], t[j] = t[j], t[i]
			}
		}
	}
}

// revolucionSolar busca el instante en que el Sol vuelve a su grado natal.
// Se hace por bisección sobre la diferencia de longitud: el Sol avanza siempre
// hacia delante, así que la función es monótona dentro del año y no hay trampas
// de convención de signo.
func revolucionSolar(natal efem.Carta, anio int) Revolucion {
	objetivo := 0.0
	for _, b := range natal.Cuerpos {
		if b.Nombre == "Sol" {
			objetivo = b.Lon
		}
	}
	// se arranca del cumpleaños aproximado y se afina
	jd0 := efem.DiaJuliano(anio, natal.Mes, float64(natal.Dia), 0) - 2
	dif := func(jd float64) float64 {
		d := math.Mod(efem.Sol(jd)-objetivo+540, 360) - 180
		return d
	}
	a, b := jd0, jd0+5
	if dif(a)*dif(b) > 0 {
		a, b = jd0-3, jd0+8 // por si el arranque cayó al otro lado
	}
	for i := 0; i < 60; i++ {
		m := (a + b) / 2
		if dif(a)*dif(m) <= 0 {
			b = m
		} else {
			a = m
		}
	}
	jd := (a + b) / 2
	rs := efem.CalcularJD(jd, natal.Lat, natal.Lon)
	return Revolucion{
		Cuando: rs.UT, Asc: rs.Asc, AscSig: int(rs.Asc / 30),
		MC: rs.MC, MCSig: int(rs.MC / 30),
		Casa: casaNatalDe(rs.Asc, natal),
	}
}

// convergencias es lo que separa esto de una lista de tránsitos.
//
// El módulo 12 lo dice sin rodeos: un periodo importa cuando tránsito,
// progresión y revolución solar apuntan a lo mismo. Así que no basta con
// contar dos señales — hace falta que vengan de TÉCNICAS DISTINTAS. Dos
// progresiones sobre el mismo punto son una sola voz repetida, que es
// exactamente el error contra el que avisa la regla de convergencia.
func convergencias(p Prediccion, T tabla) []string {
	// por cada punto natal, qué técnicas lo señalan y con qué
	tecnicas := map[string]map[string][]string{}
	apunta := func(punto, tecnica, quien string) {
		if tecnicas[punto] == nil {
			tecnicas[punto] = map[string][]string{}
		}
		tecnicas[punto][tecnica] = append(tecnicas[punto][tecnica], quien)
	}

	for _, t := range p.Transitos {
		// solo los aspectos mayores y apretados: lo demás es ruido
		if esMayor(t.Aspecto) && t.Orbe <= 1 {
			apunta(t.interno, "transito", fmt.Sprintf(T.f["porTransito"], t.Planeta))
		}
	}
	for _, g := range p.Progresiones {
		if g.interno != "" && esMayor(g.Aspecto) {
			apunta(g.interno, "progresion", fmt.Sprintf(T.f["porProgresion"], g.Planeta))
		}
	}
	// La revolución solar entra por sus ángulos: dónde cae el Ascendente del
	// año respecto a la carta natal.
	for punto, lon := range p.puntos {
		s := sep(p.Revolucion.Asc, lon)
		for _, a := range efem.TablaAspectos() {
			if esMayor(a.Nombre) && math.Abs(s-a.Angulo) <= 2 {
				apunta(punto, "revolucion", T.f["porRevolucion"])
				break
			}
		}
	}

	// Se ordena por el nombre interno del punto y no por la frase ya traducida:
	// si no, el mismo periodo sale en distinto orden en cada idioma.
	var puntos []string
	for punto := range tecnicas {
		puntos = append(puntos, punto)
	}
	sort.Strings(puntos)

	var out []string
	for _, punto := range puntos {
		porTecnica := tecnicas[punto]
		if len(porTecnica) < 2 {
			continue // una sola técnica es un murmullo, por mucho que se repita
		}
		var quienes []string
		for _, lista := range porTecnica {
			quienes = append(quienes, lista...)
		}
		sort.Strings(quienes)
		lista := ""
		for i, q := range quienes {
			if i > 0 {
				lista += ", "
			}
			lista += q
		}
		out = append(out, fmt.Sprintf(T.f["convergen"], nombreLlano(punto, T),
			len(porTecnica), lista))
	}
	if len(out) == 0 {
		out = append(out, T.f["sinConvergencia"])
	}
	return out
}

// nombreLlano da el nombre propio del planeta o del ángulo en el idioma
// activo, sin la frase de función que usa la lectura.
func nombreLlano(n string, T tabla) string {
	if T.nombres != nil {
		if v, ok := T.nombres[n]; ok {
			return v
		}
	}
	return n
}
