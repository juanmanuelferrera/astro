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
	d := dia + horaUT/24

	// La corrección gregoriana solo se aplica a partir del 15 de octubre de
	// 1582; antes de esa fecha las cuentas van por el calendario juliano, que
	// no tiene la regla de los siglos. Sin esto, cualquier fecha anterior sale
	// desplazada — para el año 837 son cuatro días enteros.
	gregoriano := anio > 1582 ||
		(anio == 1582 && (mes > 10 || (mes == 10 && d >= 15)))
	b := 0.0
	if gregoriano {
		a := math.Floor(float64(y) / 100)
		b = 2 - a + math.Floor(a/4)
	}
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

// DeDiaJuliano es la vuelta de DiaJuliano: del día juliano a fecha civil y
// hora UT. Meeus, capítulo 7. Respeta el corte gregoriano igual que la ida.
func DeDiaJuliano(jd float64) (anio, mes int, dia, horaUT float64) {
	jd += 0.5
	z := math.Floor(jd)
	f := jd - z
	a := z
	if z >= 2299161 { // a partir del 15 de octubre de 1582
		alfa := math.Floor((z - 1867216.25) / 36524.25)
		a = z + 1 + alfa - math.Floor(alfa/4)
	}
	b := a + 1524
	c := math.Floor((b - 122.1) / 365.25)
	d := math.Floor(365.25 * c)
	e := math.Floor((b - d) / 30.6001)

	diaEnt := b - d - math.Floor(30.6001*e) + f
	if e < 14 {
		mes = int(e - 1)
	} else {
		mes = int(e - 13)
	}
	if mes > 2 {
		anio = int(c - 4716)
	} else {
		anio = int(c - 4715)
	}
	dia = math.Floor(diaEnt)
	horaUT = (diaEnt - dia) * 24
	return
}
