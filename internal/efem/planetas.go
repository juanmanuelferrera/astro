package efem

import "math"

// Elementos keplerianos aproximados (Standish, JPL), época J2000, válidos ~1800-2050.
// a(UA) e I(°) L(°) varpi(°) Omega(°) y sus tasas por siglo juliano.
type elem struct{ a, e, i, l, w, o, da, de, di, dl, dw, do_ float64 }

var elems = map[string]elem{
	"Mercurio": {0.38709927, 0.20563593, 7.00497902, 252.25032350, 77.45779628, 48.33076593,
		0.00000037, 0.00001906, -0.00594749, 149472.67411175, 0.16047689, -0.12534081},
	"Venus": {0.72333566, 0.00677672, 3.39467605, 181.97909950, 131.60246718, 76.67984255,
		0.00000390, -0.00004107, -0.00078890, 58517.81538729, 0.00268329, -0.27769418},
	"Tierra": {1.00000261, 0.01671123, -0.00001531, 100.46457166, 102.93768193, 0.0,
		0.00000562, -0.00004392, -0.01294668, 35999.37244981, 0.32327364, 0.0},
	"Marte": {1.52371034, 0.09339410, 1.84969142, -4.55343205, -23.94362959, 49.55953891,
		0.00001847, 0.00007882, -0.00813131, 19140.30268499, 0.44441088, -0.29257343},
	"Júpiter": {5.20288700, 0.04838624, 1.30439695, 34.39644051, 14.72847983, 100.47390909,
		-0.00011607, -0.00013253, -0.00183714, 3034.74612775, 0.21252668, 0.20469106},
	"Saturno": {9.53667594, 0.05386179, 2.48599187, 49.95424423, 92.59887831, 113.66242448,
		-0.00125060, -0.00050991, 0.00193609, 1222.49362201, -0.41897216, -0.28867794},
	"Urano": {19.18916464, 0.04725744, 0.77263783, 313.23810451, 170.95427630, 74.01692503,
		-0.00196176, -0.00004397, -0.00242939, 428.48202785, 0.40805281, 0.04240589},
	"Neptuno": {30.06992276, 0.00859048, 1.77004347, -55.12002969, 44.96476227, 131.78422574,
		0.00026291, 0.00005105, 0.00035372, 218.45945325, -0.32241464, -0.00508664},
	"Plutón": {39.48211675, 0.24882730, 17.14001206, 238.92903833, 224.06891629, 110.30393684,
		-0.00031596, 0.00005170, 0.00004818, 145.20780515, -0.04062942, -0.01183482},
}

type vec struct{ x, y, z float64 }

// heliocentrico devuelve la posición heliocéntrica eclíptica (UA) de un planeta.
func heliocentrico(nombre string, jd float64) vec {
	el := elems[nombre]
	t := T(jd)
	a := el.a + el.da*t
	e := el.e + el.de*t
	inc := el.i + el.di*t
	L := el.l + el.dl*t
	w := el.w + el.dw*t
	om := el.o + el.do_*t

	M := norm360(L - w)
	if M > 180 {
		M -= 360
	}
	arg := w - om // argumento del perihelio

	// Kepler por Newton-Raphson
	Mr := M * Grados
	E := Mr + e*math.Sin(Mr)
	for k := 0; k < 60; k++ {
		dE := (E - e*math.Sin(E) - Mr) / (1 - e*math.Cos(E))
		E -= dE
		if math.Abs(dE) < 1e-12 {
			break
		}
	}
	// coordenadas en el plano orbital
	xp := a * (math.Cos(E) - e)
	yp := a * math.Sqrt(1-e*e) * math.Sin(E)

	cw, sw := math.Cos(arg*Grados), math.Sin(arg*Grados)
	co, so := math.Cos(om*Grados), math.Sin(om*Grados)
	ci, si := math.Cos(inc*Grados), math.Sin(inc*Grados)

	return vec{
		x: (cw*co-sw*so*ci)*xp + (-sw*co-cw*so*ci)*yp,
		y: (cw*so+sw*co*ci)*xp + (-sw*so+cw*co*ci)*yp,
		z: (sw*si)*xp + (cw*si)*yp,
	}
}

// Planeta devuelve la longitud eclíptica geocéntrica aparente en grados,
// corregida por aberración de la luz (iteración por tiempo-luz).
func Planeta(nombre string, jd float64) float64 {
	tierra := heliocentrico("Tierra", jd)
	p := heliocentrico(nombre, jd)
	for k := 0; k < 3; k++ { // corrección por tiempo-luz
		dx, dy, dz := p.x-tierra.x, p.y-tierra.y, p.z-tierra.z
		d := math.Sqrt(dx*dx + dy*dy + dz*dz)
		p = heliocentrico(nombre, jd-d*0.0057755183)
	}
	dx, dy := p.x-tierra.x, p.y-tierra.y
	lon2000 := math.Atan2(dy, dx) * Rad
	// Los elementos de Standish están referidos al equinoccio fijo J2000; las
	// posiciones astrológicas se dan en el equinoccio de la fecha. Hay que sumar
	// la precesión general en longitud (Lieske et al.).
	dpsi, _ := Nutacion(jd)
	return norm360(lon2000 + PrecesionDesdeJ2000(jd) + dpsi)
}

// PrecesionDesdeJ2000 devuelve la precesión general en longitud, en grados,
// acumulada entre J2000 y la fecha dada.
func PrecesionDesdeJ2000(jd float64) float64 {
	t := T(jd)
	return (5028.796195*t + 1.1054348*t*t + 0.00007964*t*t*t) / 3600
}

// LongitudHelio expone la longitud heliocéntrica para depuración.
func LongitudHelio(nombre string, jd float64) float64 {
	p := heliocentrico(nombre, jd)
	return norm360(math.Atan2(p.y, p.x) * Rad)
}
