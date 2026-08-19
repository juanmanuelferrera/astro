// Package efem — cálculo astronómico en Go puro, sin cgo.
// Algoritmos de Jean Meeus, "Astronomical Algorithms" (2ª ed.).
package efem

import "math"

const (
	Grados = math.Pi / 180
	Rad    = 180 / math.Pi
)

// norm360 lleva un ángulo al rango [0,360).
func norm360(x float64) float64 {
	x = math.Mod(x, 360)
	if x < 0 {
		x += 360
	}
	return x
}

// DiaJuliano convierte fecha y hora UT en día juliano. Meeus cap. 7.
func DiaJuliano(anio, mes int, dia, horaUT float64) float64 {
	y, m := anio, mes
	if m <= 2 {
		y--
		m += 12
	}
	a := math.Floor(float64(y) / 100)
	b := 2 - a + math.Floor(a/4) // calendario gregoriano
	d := dia + horaUT/24
	return math.Floor(365.25*float64(y+4716)) +
		math.Floor(30.6001*float64(m+1)) + d + b - 1524.5
}

// T devuelve los siglos julianos desde J2000.0.
func T(jd float64) float64 { return (jd - 2451545.0) / 36525.0 }

// TiempoSidereoGreenwich devuelve el TSG aparente en grados. Meeus cap. 12.
func TiempoSidereoGreenwich(jd float64) float64 {
	t := T(jd)
	theta := 280.46061837 + 360.98564736629*(jd-2451545.0) +
		0.000387933*t*t - t*t*t/38710000.0
	// corrección de nutación → tiempo sidéreo aparente
	dpsi, eps := Nutacion(jd)
	return norm360(theta + dpsi*math.Cos(eps*Grados))
}

// Nutacion devuelve (Δψ en grados, ε verdadera en grados). Meeus cap. 22, términos principales.
func Nutacion(jd float64) (float64, float64) {
	t := T(jd)
	omega := norm360(125.04452 - 1934.136261*t)
	L := norm360(280.4665 + 36000.7698*t)
	Lp := norm360(218.3165 + 481267.8813*t)

	dpsi := (-17.20*math.Sin(omega*Grados) - 1.32*math.Sin(2*L*Grados) -
		0.23*math.Sin(2*Lp*Grados) + 0.21*math.Sin(2*omega*Grados)) / 3600
	deps := (9.20*math.Cos(omega*Grados) + 0.57*math.Cos(2*L*Grados) +
		0.10*math.Cos(2*Lp*Grados) - 0.09*math.Cos(2*omega*Grados)) / 3600

	// oblicuidad media, Meeus (22.3)
	u := t / 100
	eps0 := 23.43929111 - (46.8150*t+0.00059*t*t-0.001813*t*t*t)/3600
	_ = u
	return dpsi, eps0 + deps
}

// Oblicuidad devuelve la oblicuidad verdadera de la eclíptica en grados.
func Oblicuidad(jd float64) float64 {
	_, eps := Nutacion(jd)
	return eps
}
