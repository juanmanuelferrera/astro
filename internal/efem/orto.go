package efem

import "math"

// Salida y puesta del Sol. Hace falta para el pañcāṅga —en la tradición el día
// empieza al amanecer, no a medianoche— y para los lagnas especiales de
// jyotiṣa, que se cuentan como tiempo transcurrido desde el orto.
//
// Método de Meeus, cap. 15. Se itera dos veces sobre la posición del Sol
// porque cambia mientras se resuelve la ecuación, y con dos pasadas el error
// baja del medio minuto.

// h0 es la altura del centro del Sol en el instante del orto aparente: medio
// diámetro solar por debajo del horizonte, más la refracción atmosférica.
const h0 = -0.8333

// ecuatoriales convierte la longitud eclíptica del Sol en ascensión recta y
// declinación, las dos en grados.
func ecuatoriales(jd float64) (alfa, delta float64) {
	lam := Sol(jd) * Grados
	eps := Oblicuidad(jd) * Grados
	alfa = math.Atan2(math.Cos(eps)*math.Sin(lam), math.Cos(lam)) / Grados
	delta = math.Asin(math.Sin(eps)*math.Sin(lam)) / Grados
	return norm360(alfa), delta
}

// Orto devuelve salida, tránsito y puesta del Sol en horas UT del día de jd0,
// que debe ser el día juliano a las 0h UT. lon es positiva al este.
// El tercer valor dice si hubo salida: en latitudes altas puede no haberla.
func Orto(jd0, lat, lon float64) (salida, transito, puesta float64, hay bool) {
	theta0 := TiempoSidereoGreenwich(jd0) // ya viene en grados

	calc := func(m float64) (float64, float64, float64) {
		alfa, delta := ecuatoriales(jd0 + m)
		// ángulo horario en el orto
		cosH := (math.Sin(h0*Grados) - math.Sin(lat*Grados)*math.Sin(delta*Grados)) /
			(math.Cos(lat*Grados) * math.Cos(delta*Grados))
		return alfa, delta, cosH
	}

	alfa, _, cosH := calc(0.5)
	if cosH < -1 || cosH > 1 {
		// el Sol no cruza el horizonte ese día: circumpolar o noche polar
		return 0, 0, 0, false
	}
	H0 := math.Acos(cosH) / Grados

	// tránsito: cuando la ascensión recta del Sol iguala al tiempo sidéreo local
	m0 := math.Mod((alfa-lon-theta0)/360+1, 1)
	m1 := math.Mod(m0-H0/360+1, 1)
	m2 := math.Mod(m0+H0/360+1, 1)

	// segunda pasada, ya con el Sol en su sitio de cada momento
	for i := 0; i < 2; i++ {
		for _, p := range []*float64{&m1, &m0, &m2} {
			a, d, _ := calc(*p)
			// tiempo sidéreo local en ese instante
			tsl := norm360(theta0 + 360.985647**p + lon)
			Hloc := tsl - a
			for Hloc > 180 {
				Hloc -= 360
			}
			for Hloc < -180 {
				Hloc += 360
			}
			if p == &m0 {
				*p = math.Mod(*p-Hloc/360+1, 1)
				continue
			}
			alt := math.Asin(math.Sin(lat*Grados)*math.Sin(d*Grados)+
				math.Cos(lat*Grados)*math.Cos(d*Grados)*math.Cos(Hloc*Grados)) / Grados
			den := 360 * math.Cos(d*Grados) * math.Cos(lat*Grados) * math.Sin(Hloc*Grados)
			if math.Abs(den) < 1e-9 {
				continue
			}
			*p = math.Mod(*p+(alt-h0)/den+1, 1)
		}
	}
	return m1 * 24, m0 * 24, m2 * 24, true
}
