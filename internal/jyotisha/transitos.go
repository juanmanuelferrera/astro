package jyotisha

import (
	"fmt"
	"math"
	"time"

	"astro/internal/efem"
)

type Transito struct {
	Graha    string  `json:"graha"`
	Glifo    string  `json:"glifo"`
	Posicion string  `json:"posicion"`
	RasiIdx  int     `json:"rasiIdx"`
	DesdeLag int     `json:"desdeLagna"`
	DesdeLun int     `json:"desdeLuna"`
	Retro    bool    `json:"retro"`
}

type SadeSati struct {
	Activo bool   `json:"activo"`
	Fase   string `json:"fase"`
	Nota   string `json:"nota"`
	Desde  string `json:"desde"`
	Hasta  string `json:"hasta"`
}

type Gocara struct {
	Fecha     string     `json:"fecha"`
	Transitos []Transito `json:"transitos"`
	Sade      SadeSati   `json:"sade"`
}

// posicionesSidereas devuelve las longitudes sidéreas de los grahas en un instante.
func posicionesSidereas(t time.Time) (map[string]float64, map[string]bool) {
	u := t.UTC()
	h := float64(u.Hour()) + float64(u.Minute())/60
	jd := efem.DiaJuliano(u.Year(), int(u.Month()), float64(u.Day()), h)
	pos := map[string]float64{}
	retro := map[string]bool{}
	set := func(n string, l float64) { pos[n] = Sidereo(l, jd) }
	set("Sol", efem.Sol(jd))
	set("Luna", efem.Luna(jd))
	set("Rāhu", efem.NodoLunarMedio(jd))
	set("Ketu", efem.NodoLunarMedio(jd)+180)
	retro["Rāhu"], retro["Ketu"] = true, true
	for _, n := range []string{"Marte", "Mercurio", "Júpiter", "Venus", "Saturno"} {
		l := efem.Planeta(n, jd)
		set(n, l)
		d := math.Mod(efem.Planeta(n, jd+1)-l+540, 360) - 180
		retro[n] = d < 0
	}
	return pos, retro
}

// Transitos devuelve el cielo de hoy respecto a una carta natal, con el estado
// del Sade Sati — el paso de Śani por la 12ª, la 1ª y la 2ª desde la Luna natal.
func Transitos(lagnaRasi, lunaRasi int, cuando time.Time) Gocara {
	pos, retro := posicionesSidereas(cuando)
	g := Gocara{Fecha: cuando.Format("2006-01-02")}
	for _, n := range Grahas {
		l := pos[n]
		r := int(l / 30)
		g.Transitos = append(g.Transitos, Transito{Graha: n, Glifo: Glifo[n],
			Posicion: Formato(l), RasiIdx: r, Retro: retro[n],
			DesdeLag: ((r-lagnaRasi)%12+12)%12 + 1,
			DesdeLun: ((r-lunaRasi)%12+12)%12 + 1})
	}
	// Sade Sati
	casa := ((int(pos["Saturno"]/30)-lunaRasi)%12 + 12) % 12 + 1
	switch casa {
	case 12:
		g.Sade = SadeSati{Activo: true, Fase: "primera",
			Nota: "Śani en la 12ª desde la Luna: empieza el ciclo. Suele notarse como desgaste y gasto."}
	case 1:
		g.Sade = SadeSati{Activo: true, Fase: "central",
			Nota: "Śani sobre el signo de la Luna: la fase más exigente. Emocional y de responsabilidad."}
	case 2:
		g.Sade = SadeSati{Activo: true, Fase: "última",
			Nota: "Śani en la 2ª desde la Luna: la salida. Suele tocar recursos y familia."}
	default:
		g.Sade = SadeSati{Activo: false,
			Nota: fmt.Sprintf("Śani transita la casa %d desde la Luna natal: no hay Sade Sati.", casa)}
	}
	// buscar el próximo comienzo o el final, sondeando mes a mes
	paso := cuando
	for i := 0; i < 400; i++ {
		paso = paso.AddDate(0, 1, 0)
		p, _ := posicionesSidereas(paso)
		c := ((int(p["Saturno"]/30)-lunaRasi)%12 + 12) % 12 + 1
		dentro := c == 12 || c == 1 || c == 2
		if dentro != g.Sade.Activo {
			if g.Sade.Activo {
				g.Sade.Hasta = paso.Format("2006-01")
			} else {
				g.Sade.Desde = paso.Format("2006-01")
			}
			break
		}
	}
	return g
}
