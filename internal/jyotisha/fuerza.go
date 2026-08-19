package jyotisha

import "math"

// Mūlatrikoṇa: tramo dentro del signo propio donde el graha va especialmente cómodo.
var mulatrikona = map[string]struct {
	signo      int
	desde, has float64
}{
	"Sol": {4, 0, 20}, "Luna": {1, 4, 30}, "Marte": {0, 0, 12},
	"Mercurio": {5, 16, 20}, "Júpiter": {8, 0, 10}, "Venus": {6, 0, 15}, "Saturno": {10, 0, 20},
}

// EnMulatrikona indica si la posición cae en el tramo mūlatrikoṇa del graha.
func EnMulatrikona(graha string, lon float64) bool {
	m, ok := mulatrikona[graha]
	if !ok {
		return false
	}
	s, g := int(lon/30), math.Mod(lon, 30)
	return s == m.signo && g >= m.desde && g < m.has
}

// Combustión (asta): márgenes en grados respecto al Sol.
var margenAsta = map[string]float64{
	"Luna": 12, "Marte": 17, "Mercurio": 14, "Júpiter": 11, "Venus": 10, "Saturno": 15,
}

// Combusto indica si el graha está quemado por cercanía al Sol, y a cuántos grados.
// Mercurio y Venus retrógrados tienen margen menor.
func Combusto(graha string, lon, lonSol float64, retro bool) (bool, float64) {
	m, ok := margenAsta[graha]
	if !ok {
		return false, 0
	}
	if retro && graha == "Mercurio" {
		m = 12
	}
	if retro && graha == "Venus" {
		m = 8
	}
	d := math.Abs(lon - lonSol)
	if d > 180 {
		d = 360 - d
	}
	return d <= m, math.Round(d*100) / 100
}

// Dig-bala: casa donde cada graha alcanza su fuerza direccional plena.
var digBala = map[string]int{
	"Júpiter": 1, "Mercurio": 1, "Sol": 10, "Marte": 10, "Saturno": 7, "Luna": 4, "Venus": 4,
}

// TieneDigBala indica si el graha está en su casa de fuerza direccional.
func TieneDigBala(graha string, bhava int) bool { return digBala[graha] == bhava }

// Kendra dice si dos signos guardan relación de ángulo (1, 4, 7, 10).
func Kendra(a, b int) bool {
	d := ((a-b)%12 + 12) % 12
	return d == 0 || d == 3 || d == 6 || d == 9
}

// Trikona dice si dos signos guardan relación de trígono (1, 5, 9).
func Trikona(a, b int) bool {
	d := ((a-b)%12 + 12) % 12
	return d == 0 || d == 4 || d == 8
}

var mahapurusha = map[string]string{
	"Marte": "Ruchaka", "Mercurio": "Bhadra", "Júpiter": "Haṁsa",
	"Venus": "Mālavya", "Saturno": "Śaśa",
}
