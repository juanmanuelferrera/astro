package jyotisha

import (
	"math"
	"time"
)

// Daśās distintas de la vimśottarī.
//
// La vimśottarī es la que se usa siempre, pero no es la única, y las clásicas
// recomiendan contrastarla. Cuando dos sistemas señalan el mismo periodo, la
// lectura se sostiene; cuando solo lo dice uno, es un murmullo — la misma regla
// de convergencia de todo lo demás, aplicada al reloj.
//
// Aquí van tres:
//
//   Aṣṭottarī — 108 años en ocho grahas, sin Ketu. Arranca del nakṣatra de la
//               Luna igual que la vimśottarī, pero reparte distinto.
//   Yoginī    — 36 años en ocho yoginīs, cada una con su graha. Es corta, así
//               que en una vida se recorre tres veces; se usa para afinar.
//   Cara      — la de Jaimini, y funciona con otra lógica: no cuelga de la
//               Luna sino de los RĀŚIS, y la duración de cada periodo sale de
//               contar del signo a su señor.

// ── Aṣṭottarī: 108 años ──
var cicloAstottari = []struct {
	senor string
	anios float64
}{{"Sol", 6}, {"Luna", 15}, {"Marte", 8}, {"Mercurio", 17},
	{"Saturno", 10}, {"Júpiter", 19}, {"Rāhu", 12}, {"Venus", 21}}

// El reparto de nakṣatras por señor no es regular como en la vimśottarī: cada
// señor gobierna un número distinto de nakṣatras. Esta es la tabla clásica,
// contando desde Ārdrā (el nakṣatra 6, índice 5).
var nakAstottari = [27]int{
	// índice de señor en cicloAstottari para cada nakṣatra, de Aśvinī a Revatī
	6, 6, 7, 7, 7, 0, 0, 0, 0, 1, 1, 1, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7,
}

// ── Yoginī: 36 años ──
var cicloYogini = []struct {
	nombre, senor string
	anios         float64
}{{"Maṅgalā", "Luna", 1}, {"Piṅgalā", "Sol", 2}, {"Dhānyā", "Júpiter", 3},
	{"Bhrāmarī", "Marte", 4}, {"Bhadrikā", "Mercurio", 5}, {"Ulkā", "Saturno", 6},
	{"Siddhā", "Venus", 7}, {"Saṅkaṭā", "Rāhu", 8}}

type Dasa struct {
	Sistema string    `json:"sistema"`
	Total   float64   `json:"total"` // años del ciclo entero
	Nota    string    `json:"nota"`
	Ciclos  []Periodo `json:"ciclos"`
}

// periodos arma la lista a partir de un ciclo cualquiera, con la parte ya
// consumida al nacer descontada del primero.
func periodos(nombres []string, anios []float64, i int, frac float64,
	nacimiento time.Time, cuantos int) []Periodo {
	n := len(nombres)
	transcurrido := anios[i] * frac
	inicio := nacimiento.Add(-time.Duration(transcurrido * anioSideral * 24 * float64(time.Hour)))
	ahora := time.Now()
	var out []Periodo
	for k := 0; k < cuantos; k++ {
		j := (i + k) % n
		fin := inicio.Add(time.Duration(anios[j] * anioSideral * 24 * float64(time.Hour)))
		out = append(out, Periodo{Senor: nombres[j],
			Desde: inicio.Format("2006-01-02"), Hasta: fin.Format("2006-01-02"),
			Anios: anios[j], Actual: ahora.After(inicio) && ahora.Before(fin)})
		inicio = fin
	}
	return out
}

// Astottari reparte 108 años entre ocho grahas. No entra Ketu.
func Astottari(lonLuna float64, nacimiento time.Time, cuantos int) Dasa {
	const span = 360.0 / 27.0
	ni := int(lonLuna / span)
	if ni > 26 {
		ni = 26
	}
	i := nakAstottari[ni]
	frac := math.Mod(lonLuna, span) / span

	nombres := make([]string, len(cicloAstottari))
	anios := make([]float64, len(cicloAstottari))
	total := 0.0
	for k, c := range cicloAstottari {
		nombres[k], anios[k] = c.senor, c.anios
		total += c.anios
	}
	return Dasa{Sistema: "Aṣṭottarī", Total: total,
		Ciclos: periodos(nombres, anios, i, frac, nacimiento, cuantos)}
}

// Yogini reparte 36 años entre ocho yoginīs. El punto de partida sale del
// nakṣatra de la Luna por una cuenta distinta: (nakṣatra + 3) módulo 8.
func Yogini(lonLuna float64, nacimiento time.Time, cuantos int) Dasa {
	const span = 360.0 / 27.0
	ni := int(lonLuna / span)
	if ni > 26 {
		ni = 26
	}
	i := (ni + 3) % 8
	frac := math.Mod(lonLuna, span) / span

	nombres := make([]string, len(cicloYogini))
	anios := make([]float64, len(cicloYogini))
	total := 0.0
	for k, c := range cicloYogini {
		nombres[k], anios[k] = c.nombre+" ("+c.senor+")", c.anios
		total += c.anios
	}
	return Dasa{Sistema: "Yoginī", Total: total,
		Ciclos: periodos(nombres, anios, i, frac, nacimiento, cuantos)}
}

// ── Cara daśā de Jaimini ──
//
// No cuelga de la Luna sino del Lagna, y no reparte grahas sino RĀŚIS. La
// duración de cada uno sale de contar del rāśi al sitio donde está su señor,
// menos uno. El sentido de la cuenta depende de si el rāśi es impar o par.

// duracionCara da los años que dura el periodo de un rāśi.
func duracionCara(rasi int, rasiDe map[string]int) float64 {
	senor := SenorRasi[rasi]
	// Escorpio y Acuario tienen dos señores en Jaimini; se toma el que esté
	// mejor colocado, y a falta de criterio fino, el tradicional.
	pos, ok := rasiDe[senor]
	if !ok {
		return 1
	}
	impar := rasi%2 == 0 // índice 0 es Meṣa, que es impar en la cuenta clásica
	var n int
	if impar {
		n = ((pos-rasi)%12 + 12) % 12 // hacia delante
	} else {
		n = ((rasi-pos)%12 + 12) % 12 // hacia atrás
	}
	if n == 0 {
		return 12 // el señor en su propio signo da el periodo entero
	}
	return float64(n)
}

// Cara arma la secuencia de rāśis desde el Lagna. El sentido también depende
// de si el Lagna cae en signo impar o par.
func Cara(lagnaRasi int, rasiDe map[string]int, nacimiento time.Time, cuantos int) Dasa {
	adelante := lagnaRasi%2 == 0
	var nombres []string
	var anios []float64
	total := 0.0
	for k := 0; k < 12; k++ {
		var r int
		if adelante {
			r = (lagnaRasi + k) % 12
		} else {
			r = ((lagnaRasi-k)%12 + 12) % 12
		}
		a := duracionCara(r, rasiDe)
		nombres = append(nombres, Rasis[r])
		anios = append(anios, a)
		total += a
	}
	if cuantos > 12 {
		cuantos = 12
	}
	return Dasa{Sistema: "Cara (Jaimini)", Total: total,
		Ciclos: periodos(nombres, anios, 0, 0, nacimiento, cuantos)}
}
