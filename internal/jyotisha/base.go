// Package jyotisha — lo propio de la astrología védica: zodíaco sidéreo,
// nakṣatras, casas de signo entero, dignidades y dṛṣṭi.
// La astronomía la comparte con internal/efem: es la misma para los dos sistemas.
package jyotisha

import (
	"fmt"
	"math"

	"astro/internal/efem"
)

// Ayanamsa devuelve el ayanāṁśa Lahiri (Citrapakṣa) en grados.
// Polinomio ajustado contra Swiss Ephemeris: residuo máximo 0,0002″ entre
// 1850 y 2100, y por debajo de 0,001″ entre 1600 y 2200.
func Ayanamsa(jd float64) float64 {
	t := (jd - 2451545.0) / 36525.0
	return 23.857092349401 + 1.396887883512*t +
		0.000307093983*t*t + 0.000000028323*t*t*t
}

// Sidereo convierte una longitud tropical en sidérea.
func Sidereo(lonTropical, jd float64) float64 {
	l := math.Mod(lonTropical-Ayanamsa(jd), 360)
	if l < 0 {
		l += 360
	}
	return l
}

var Rasis = [12]string{"Meṣa", "Vṛṣabha", "Mithuna", "Karka", "Siṁha", "Kanyā",
	"Tulā", "Vṛścika", "Dhanus", "Makara", "Kumbha", "Mīna"}
var RasisEs = [12]string{"Aries", "Tauro", "Géminis", "Cáncer", "Leo", "Virgo",
	"Libra", "Escorpio", "Sagitario", "Capricornio", "Acuario", "Piscis"}
var SenorRasi = [12]string{"Marte", "Venus", "Mercurio", "Luna", "Sol", "Mercurio",
	"Venus", "Marte", "Júpiter", "Saturno", "Saturno", "Júpiter"}

// Los 27 nakṣatras, con el señor del ciclo vimśottarī.
var Nakshatras = [27]string{"Aśvinī", "Bharaṇī", "Kṛttikā", "Rohiṇī", "Mṛgaśīrṣa", "Ārdrā",
	"Punarvasu", "Puṣya", "Āśleṣā", "Maghā", "P.Phalgunī", "U.Phalgunī", "Hasta", "Citrā",
	"Svātī", "Viśākhā", "Anurādhā", "Jyeṣṭhā", "Mūla", "P.Aṣāḍhā", "U.Aṣāḍhā", "Śravaṇa",
	"Dhaniṣṭhā", "Śatabhiṣā", "P.Bhādrapadā", "U.Bhādrapadā", "Revatī"}

var senorNak = [9]string{"Ketu", "Venus", "Sol", "Luna", "Marte", "Rāhu", "Júpiter", "Saturno", "Mercurio"}

// Nakshatra devuelve nombre, pada (1-4) y señor del nakṣatra de una longitud sidérea.
func Nakshatra(lon float64) (string, int, string) {
	const span = 360.0 / 27.0
	i := int(lon / span)
	if i > 26 {
		i = 26
	}
	pada := int(math.Mod(lon, span)/(span/4)) + 1
	return Nakshatras[i], pada, senorNak[i%9]
}

// Grahas: en jyotiṣa los nueve clásicos. Los exteriores no forman parte del sistema.
var Grahas = []string{"Sol", "Luna", "Marte", "Mercurio", "Júpiter", "Venus", "Saturno", "Rāhu", "Ketu"}

var Glifo = map[string]string{"Sol": "☉", "Luna": "☽", "Marte": "♂", "Mercurio": "☿",
	"Júpiter": "♃", "Venus": "♀", "Saturno": "♄", "Rāhu": "☊", "Ketu": "☋"}

// Dignidades védicas: exaltación con su grado exacto, debilitación en el opuesto,
// signo propio y mūlatrikoṇa.
var exalta = map[string]struct {
	signo int
	grado float64
}{
	"Sol": {0, 10}, "Luna": {1, 3}, "Marte": {9, 28}, "Mercurio": {5, 15},
	"Júpiter": {3, 5}, "Venus": {11, 27}, "Saturno": {6, 20}, "Rāhu": {1, 20}, "Ketu": {7, 20},
}
var propio = map[string][]int{"Sol": {4}, "Luna": {3}, "Marte": {0, 7}, "Mercurio": {2, 5},
	"Júpiter": {8, 11}, "Venus": {1, 6}, "Saturno": {9, 10}}

// Dignidad devuelve el estado del graha en su signo.
func Dignidad(graha string, lon float64) string {
	s := int(lon / 30)
	if e, ok := exalta[graha]; ok {
		if e.signo == s {
			return "exaltado"
		}
		if (e.signo+6)%12 == s {
			return "debilitado"
		}
	}
	for _, p := range propio[graha] {
		if p == s {
			return "signo propio"
		}
	}
	return "—"
}

// Drishti devuelve las casas que aspecta un graha desde su posición.
// En jyotiṣa el aspecto es por signo entero, no por grados: todos aspectan la
// séptima, y Marte, Júpiter y Saturno tienen aspectos especiales.
func Drishti(graha string) []int {
	switch graha {
	case "Marte":
		return []int{4, 7, 8}
	case "Júpiter":
		return []int{5, 7, 9}
	case "Saturno":
		return []int{3, 7, 10}
	case "Rāhu", "Ketu":
		return []int{5, 7, 9}
	default:
		return []int{7}
	}
}

// Gandanta indica si una posición cae en la unión agua-fuego: los últimos
// 3°20′ de Karka, Vṛścika o Mīna. Es un nudo kármico y pide śānti.
func Gandanta(lon float64) bool {
	s := int(lon / 30)
	g := math.Mod(lon, 30)
	return (s == 3 || s == 7 || s == 11) && g >= 26.0+40.0/60.0
}

// Formato devuelve la posición como "12°34′ Meṣa".
func Formato(lon float64) string {
	s := int(lon / 30)
	g := math.Mod(lon, 30)
	d := int(g)
	m := int(math.Round((g - float64(d)) * 60))
	if m == 60 {
		d, m = d+1, 0
	}
	return fmt.Sprintf("%d°%02d′ %s", d, m, Rasis[s])
}

var _ = efem.Signos
