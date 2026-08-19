package efem

import "math"

// decl devuelve la declinación (grados) del punto de la eclíptica de longitud lon.
func decl(lon, eps float64) float64 {
	return math.Asin(math.Sin(eps*Grados)*math.Sin(lon*Grados)) * Rad
}

// ascRecta devuelve la ascensión recta (grados) del punto eclíptico de longitud lon.
func ascRecta(lon, eps float64) float64 {
	return norm360(math.Atan2(
		math.Sin(lon*Grados)*math.Cos(eps*Grados),
		math.Cos(lon*Grados)) * Rad)
}

// MedioCielo devuelve la longitud eclíptica del MC. tsl = tiempo sidéreo local en grados.
func MedioCielo(tsl, eps float64) float64 {
	return norm360(math.Atan2(
		math.Sin(tsl*Grados),
		math.Cos(tsl*Grados)*math.Cos(eps*Grados)) * Rad)
}

// Ascendente devuelve la longitud eclíptica del Ascendente.
func Ascendente(tsl, lat, eps float64) float64 {
	t := tsl * Grados
	e := eps * Grados
	f := lat * Grados
	a := math.Atan2(math.Cos(t), -(math.Sin(t)*math.Cos(e) + math.Tan(f)*math.Sin(e)))
	return norm360(a * Rad)
}

// semiarco devuelve el semiarco diurno (grados) del punto eclíptico de longitud lon.
// Devuelve false si el punto es circumpolar (no sale o no se pone).
func semiarco(lon, lat, eps float64) (float64, bool) {
	d := decl(lon, eps)
	x := -math.Tan(lat*Grados) * math.Tan(d*Grados)
	if x <= -1 || x >= 1 {
		return 0, false
	}
	return math.Acos(x) * Rad, true
}

// anguloHorario devuelve el ángulo horario (−180..180) del punto de longitud lon.
func anguloHorario(lon, tsl, eps float64) float64 {
	h := math.Mod(tsl-ascRecta(lon, eps)+540, 360) - 180
	return h
}

// resolverCuspide halla la longitud eclíptica cuyo ángulo horario coincide con
// el objetivo, que depende de su propio semiarco. Bisección: robusta y sin
// depender de convenios de signo.
//
// Ángulos horarios de referencia: 0 en el MC, −SD al salir (Asc), −180 en el IC.
//   cúspide 11: H = −SD/3      cúspide 12: H = −2·SD/3
//   cúspide 2:  H = −SD − SN/3 cúspide 3:  H = −SD − 2·SN/3   (SN = 180 − SD)
func resolverCuspide(tsl, lat, eps float64, objetivo func(sd float64) float64, aprox float64) (float64, bool) {
	g := func(lon float64) (float64, bool) {
		sd, ok := semiarco(lon, lat, eps)
		if !ok {
			return 0, false
		}
		return math.Mod(anguloHorario(lon, tsl, eps)-objetivo(sd)+540, 360) - 180, true
	}
	lo, hi := aprox-50.0, aprox+50.0
	flo, ok1 := g(lo)
	fhi, ok2 := g(hi)
	if !ok1 || !ok2 || flo*fhi > 0 {
		return 0, false
	}
	for k := 0; k < 90; k++ {
		mid := (lo + hi) / 2
		fm, ok := g(mid)
		if !ok {
			return 0, false
		}
		if flo*fm <= 0 {
			hi = mid
		} else {
			lo, flo = mid, fm
		}
	}
	return norm360((lo + hi) / 2), true
}

// Cuspides devuelve las 12 cúspides por Plácido. Si el sistema no puede
// resolverse (latitudes polares) devuelve casas iguales y ok=false.
func Cuspides(tsl, lat, eps float64) ([12]float64, bool) {
	var c [12]float64
	asc := Ascendente(tsl, lat, eps)
	mc := MedioCielo(tsl, eps)

	iguales := func() [12]float64 {
		var e [12]float64
		for i := range e {
			e[i] = norm360(asc + float64(i)*30)
		}
		return e
	}
	if math.Abs(lat) > 66.0 {
		return iguales(), false
	}

	arcoMC := math.Mod(asc-mc+360, 360)  // del MC al Asc, por encima del horizonte
	arcoIC := math.Mod(norm360(mc+180)-asc+360, 360) // del Asc al IC, por debajo

	tipos := []struct {
		idx   int
		obj   func(float64) float64
		aprox float64
	}{
		{10, func(sd float64) float64 { return -sd / 3 }, norm360(mc + arcoMC/3)},
		{11, func(sd float64) float64 { return -2 * sd / 3 }, norm360(mc + 2*arcoMC/3)},
		{1, func(sd float64) float64 { return -sd - (180-sd)/3 }, norm360(asc + arcoIC/3)},
		{2, func(sd float64) float64 { return -sd - 2*(180-sd)/3 }, norm360(asc + 2*arcoIC/3)},
	}
	for _, t := range tipos {
		v, ok := resolverCuspide(tsl, lat, eps, t.obj, t.aprox)
		if !ok {
			return iguales(), false
		}
		c[t.idx] = v
	}
	c[0], c[9] = asc, mc
	c[6], c[3] = norm360(asc+180), norm360(mc+180)
	c[4] = norm360(c[10] + 180)
	c[5] = norm360(c[11] + 180)
	c[7] = norm360(c[1] + 180)
	c[8] = norm360(c[2] + 180)
	return c, true
}

// CuspidesIguales devuelve las doce cúspides del sistema de casas iguales.
func CuspidesIguales(asc float64) [12]float64 {
	var c [12]float64
	for i := range c {
		c[i] = norm360(asc + float64(i)*30)
	}
	return c
}
