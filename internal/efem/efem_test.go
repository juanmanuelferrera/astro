package efem

import (
	"math"
	"testing"
)

// difAng devuelve la diferencia entre dos ángulos en grados, en [-180, 180].
func difAng(a, b float64) float64 {
	d := math.Mod(a-b+540, 360) - 180
	return d
}

// arcmin convierte grados a minutos de arco.
func arcmin(g float64) float64 { return math.Abs(g) * 60 }

// TestPosiciones mide cada cuerpo contra Swiss Ephemeris. El Sol y la Luna
// llevan series completas de Meeus; los planetas salen de los elementos
// keplerianos de Standish, que son más burdos a propósito — no necesitan
// ficheros de efemérides y valen de sobra para astrología.
//
// Se mide en dos ventanas porque la tabla de Standish que usa el motor está
// dada para 1800-2050 y fuera de ahí se degrada. Dentro, que es donde caen
// todas las cartas natales de gente viva, el margen es estrecho. Fuera se
// afloja, y queda escrito cuánto: Saturno llega a 19′ hacia 2070, por la gran
// desigualdad con Júpiter, que unos elementos keplerianos no pueden recoger.
//
// Nada de esto mueve un planeta de signo en tres siglos — eso lo comprueba
// TestNoCambianDeSigno, y es la prueba que de verdad importa.
func TestPosiciones(t *testing.T) {
	t.Run("dentro de 1800-2050", func(t *testing.T) {
		comprobar(t, map[string]float64{
			"Sol": 1, "Luna": 3, "Mercurio": 12, "Venus": 6, "Marte": 12,
			"Júpiter": 8, "Saturno": 12, "Urano": 12, "Neptuno": 5, "Plutón": 5,
		}, func(a int) bool { return a >= 1800 && a <= 2050 })
	})
	t.Run("fuera, hasta 2100", func(t *testing.T) {
		comprobar(t, map[string]float64{
			"Sol": 1, "Luna": 3, "Mercurio": 12, "Venus": 6, "Marte": 12,
			"Júpiter": 8, "Saturno": 20, "Urano": 25, "Neptuno": 25, "Plutón": 25,
		}, func(a int) bool { return a > 2050 })
	})
}

func comprobar(t *testing.T, margen map[string]float64, dentro func(int) bool) {
	t.Helper()
	peor := map[string]float64{}
	peorFecha := map[string]refFecha{}

	for _, f := range refFechas {
		if !dentro(f.Anio) {
			continue
		}
		for _, c := range f.Cuerpos {
			var mio float64
			switch c.Nombre {
			case "Sol":
				mio = Sol(f.JD)
			case "Luna":
				mio = Luna(f.JD)
			default:
				mio = Planeta(c.Nombre, f.JD)
			}
			e := arcmin(difAng(mio, c.Lon))
			if e > peor[c.Nombre] {
				peor[c.Nombre], peorFecha[c.Nombre] = e, f
			}
		}
	}
	for _, n := range Orden[:10] {
		e, f := peor[n], peorFecha[n]
		if e > margen[n] {
			t.Errorf("%s se desvía %.2f′ (margen %.0f′) el %d-%02d-%02d",
				n, e, margen[n], f.Anio, f.Mes, f.Dia)
		} else {
			t.Logf("%-9s peor %6.2f′  (margen %2.0f′)", n, e, margen[n])
		}
	}
}

// TestNoCambianDeSigno es la prueba que de verdad importa en astrología. Da
// igual desviarse diez minutos de arco si el planeta sigue en el mismo signo;
// lo que no puede pasar nunca es que salga en otro.
func TestNoCambianDeSigno(t *testing.T) {
	total, distintos := 0, 0
	for _, f := range refFechas {
		for _, c := range f.Cuerpos {
			var mio float64
			switch c.Nombre {
			case "Sol":
				mio = Sol(f.JD)
			case "Luna":
				mio = Luna(f.JD)
			default:
				mio = Planeta(c.Nombre, f.JD)
			}
			total++
			if int(norm360(mio)/30) != int(norm360(c.Lon)/30) {
				distintos++
				t.Errorf("%s cae en otro signo el %d-%02d-%02d: %.4f° contra %.4f°",
					c.Nombre, f.Anio, f.Mes, f.Dia, mio, c.Lon)
			}
		}
	}
	t.Logf("%d posiciones comprobadas, %d cambian de signo", total, distintos)
}

// TestCasas comprueba el Ascendente, el Medio Cielo y las doce cúspides de
// Plácido. Aquí el margen es mucho más estrecho: son geometría pura, sin
// series aproximadas de por medio.
func TestCasas(t *testing.T) {
	var peorAsc, peorMC, peorCusp float64
	for _, r := range refCasasTabla {
		eps := Oblicuidad(r.JD)
		tsl := TiempoSidereoGreenwich(r.JD) + r.Lon
		if e := arcmin(difAng(Ascendente(tsl, r.Lat, eps), r.Asc)); e > peorAsc {
			peorAsc = e
		}
		if e := arcmin(difAng(MedioCielo(tsl, eps), r.MC)); e > peorMC {
			peorMC = e
		}
		c, ok := Cuspides(tsl, r.Lat, eps)
		if !ok {
			continue // latitud donde Plácido no puede: se cae a casas iguales
		}
		for i := 0; i < 12; i++ {
			if e := arcmin(difAng(c[i], r.Cuspides[i])); e > peorCusp {
				peorCusp = e
			}
		}
	}
	// medio minuto de arco es de sobra: son segundos de arco en la práctica
	for n, e := range map[string]float64{"Ascendente": peorAsc, "Medio Cielo": peorMC, "cúspides": peorCusp} {
		if e > 0.5 {
			t.Errorf("%s se desvía %.4f′ (%.2f″)", n, e, e*60)
		}
	}
	t.Logf("Asc %.4f″  MC %.4f″  cúspides %.4f″", peorAsc*60, peorMC*60, peorCusp*60)
}

// TestNodos: el medio es una fórmula cerrada y tiene que salir casi exacto.
// El verdadero lleva cinco términos periódicos, así que admite algo más.
func TestNodos(t *testing.T) {
	var peorM, peorV float64
	for _, f := range refFechas {
		if e := arcmin(difAng(NodoLunarMedio(f.JD), f.NodoMedio)); e > peorM {
			peorM = e
		}
		if e := arcmin(difAng(NodoLunarVerdadero(f.JD), f.NodoVerdad)); e > peorV {
			peorV = e
		}
	}
	if peorM > 0.5 {
		t.Errorf("el nodo medio se desvía %.3f′", peorM)
	}
	// El nodo verdadero admite más porque no se está comparando lo mismo: el de
	// swisseph es el nodo osculador de la órbita instantánea, y el del motor es
	// la fórmula de Meeus, que lo aproxima con cinco términos periódicos. La
	// diferencia oscila, no crece, así que es de definición y no un error que
	// se vaya acumulando.
	if peorV > 15 {
		t.Errorf("el nodo verdadero se desvía %.3f′", peorV)
	}
	t.Logf("nodo medio %.4f′  ·  nodo verdadero %.4f′", peorM, peorV)
}

// TestOrto compara salida y puesta del Sol. Un minuto de reloj es un cuarto de
// grado de Ghaṭī Lagna, así que el margen tiene que ser estrecho.
func TestOrto(t *testing.T) {
	var peorS, peorP float64
	n := 0
	for _, r := range refOrtos {
		s, _, p, hay := Orto(r.JD0, r.Lat, r.Lon)
		if !hay {
			continue
		}
		n++
		if e := math.Abs(s-r.Salida) * 60; e > peorS {
			peorS = e
		}
		if e := math.Abs(p-r.Puesta) * 60; e > peorP {
			peorP = e
		}
	}
	if n == 0 {
		t.Skip("sin datos de referencia de orto")
	}
	if peorS > 1.5 || peorP > 1.5 {
		t.Errorf("orto %.2f min, ocaso %.2f min de desvío (margen 1,5)", peorS, peorP)
	}
	t.Logf("%d días: orto %.2f min  ·  ocaso %.2f min", n, peorS, peorP)
}

// TestDiaJuliano contra los valores tabulados de Meeus, capítulo 7.
func TestDiaJuliano(t *testing.T) {
	casos := []struct {
		a, m int
		d    float64
		jd   float64
	}{
		{2000, 1, 1.5, 2451545.0},
		{1999, 1, 1.0, 2451179.5},
		{1987, 1, 27.0, 2446822.5},
		{1900, 1, 1.0, 2415020.5},
		{1600, 1, 1.0, 2305447.5},
		{837, 4, 10.3, 2026871.8},
	}
	for _, c := range casos {
		d := math.Floor(c.d)
		jd := DiaJuliano(c.a, c.m, d, (c.d-d)*24)
		if math.Abs(jd-c.jd) > 1e-6 {
			t.Errorf("%d-%02d-%.1f da %.6f, debería ser %.6f", c.a, c.m, c.d, jd, c.jd)
		}
	}
}

// TestTiempoSidereo contra el ejemplo 12.a de Meeus.
func TestTiempoSidereo(t *testing.T) {
	// 1987 abril 10, 0h UT → 13h 10m 46,3668s de tiempo sidéreo medio
	got := TiempoSidereoGreenwich(2446895.5)
	esperado := (13 + 10.0/60 + 46.3668/3600) * 15
	if d := math.Abs(difAng(got, esperado)); d > 0.01 {
		t.Errorf("tiempo sidéreo %.6f°, esperado %.6f° (aparente contra medio: %.2f″)",
			got, esperado, d*3600)
	}
}

// TestCasasIguales: doce tramos de 30° desde el Ascendente, sin más. Es el
// sistema que Margaret Hone prefería, y el que salva las latitudes donde
// Plácido no puede calcular.
func TestCasasIguales(t *testing.T) {
	for _, asc := range []float64{0, 29.99, 137.5, 359.9} {
		c := CuspidesIguales(asc)
		if math.Abs(c[0]-asc) > 1e-9 {
			t.Errorf("la casa 1 no arranca en el Ascendente: %.4f contra %.4f", c[0], asc)
		}
		for i := 0; i < 12; i++ {
			sig := c[(i+1)%12]
			d := math.Mod(sig-c[i]+360, 360)
			if math.Abs(d-30) > 1e-9 {
				t.Errorf("de la casa %d a la %d hay %.6f° y deberían ser 30", i+1, (i+1)%12+1, d)
			}
			if c[i] < 0 || c[i] >= 360 {
				t.Errorf("la cúspide %d vale %.4f, fuera de rango", i+1, c[i])
			}
		}
	}
}

// En latitudes polares Plácido no puede: hay grados que no salen ni se ponen,
// así que no tienen arcos que trisecar. Lo que no puede pasar es que devuelva
// números sin avisar.
func TestPlacidoAvisaEnLatitudPolar(t *testing.T) {
	jd := DiaJuliano(2000, 6, 21, 12)
	eps := Oblicuidad(jd)
	tsl := norm360(TiempoSidereoGreenwich(jd))
	for _, lat := range []float64{78, -78, 85} {
		if _, ok := Cuspides(tsl, lat, eps); ok {
			t.Errorf("a %.0f° de latitud Plácido dice que sí puede, y no puede", lat)
		}
	}
	for _, lat := range []float64{0, 41.4, -33.9, 59.9} {
		if _, ok := Cuspides(tsl, lat, eps); !ok {
			t.Errorf("a %.1f° de latitud Plácido debería poder", lat)
		}
	}
}

// Las dignidades esenciales: domicilio, exaltación, exilio y caída. La tabla
// tiene que ser coherente consigo misma — el exilio es siempre el opuesto del
// domicilio, y la caída el opuesto de la exaltación.
func TestDignidades(t *testing.T) {
	domicilios := map[string][]int{
		"Sol": {4}, "Luna": {3}, "Mercurio": {2, 5}, "Venus": {1, 6},
		"Marte": {0, 7}, "Júpiter": {8, 11}, "Saturno": {9, 10},
	}
	for planeta, signos := range domicilios {
		for _, s := range signos {
			if d := dignidad(planeta, s); d != "domicilio" {
				t.Errorf("%s en el signo %d da %q y debería ser domicilio", planeta, s+1, d)
			}
			if d := dignidad(planeta, (s+6)%12); d != "exilio" {
				t.Errorf("%s enfrente de su domicilio da %q y debería ser exilio", planeta, d)
			}
		}
	}
	exaltaciones := map[string]int{"Sol": 0, "Luna": 1, "Mercurio": 5,
		"Venus": 11, "Marte": 9, "Júpiter": 3, "Saturno": 6}
	for planeta, s := range exaltaciones {
		// Mercurio es el caso raro y va aparte: se exalta en Virgo, que además
		// es su propio domicilio. Ahí se declara el domicilio, y enfrente el
		// exilio, que es lo que hace todo el mundo cuando las dos cosas caen
		// en el mismo signo.
		if planeta == "Mercurio" {
			continue
		}
		if d := dignidad(planeta, s); d != "exaltación" {
			t.Errorf("%s en el signo %d da %q y debería ser exaltación", planeta, s+1, d)
		}
		if d := dignidad(planeta, (s+6)%12); d != "caída" {
			t.Errorf("%s enfrente de su exaltación da %q y debería ser caída", planeta, d)
		}
	}
	if d := dignidad("Mercurio", 5); d != "domicilio" {
		t.Errorf("Mercurio en Virgo da %q: es domicilio y exaltación a la vez", d)
	}
	if d := dignidad("Mercurio", 11); d != "exilio" && d != "caída" {
		t.Errorf("Mercurio en Piscis da %q, y es exilio y caída a la vez", d)
	}
}

// casaDe reparte una longitud entre doce cúspides. El caso que rompe siempre es
// la casa que cruza el 0° de Aries.
func TestCasaDe(t *testing.T) {
	// casas iguales desde 350°: la casa 1 va de 350° a 20°, cruzando el cero
	c := CuspidesIguales(350)
	// casa 1: 350-20 · casa 2: 20-50 · casa 3: 50-80 · … · casa 12: 320-350
	casos := map[float64]int{
		350: 1, 355: 1, 0: 1, 10: 1, 19.9: 1,
		20: 2, 45: 2, 50: 3, 180: 7, 320: 12, 349.9: 12,
	}
	for lon, esperada := range casos {
		if got := casaDe(lon, c); got != esperada {
			t.Errorf("%.1f° cae en la casa %d y debería caer en la %d", lon, got, esperada)
		}
	}
	// y toda longitud tiene que caer en alguna casa
	for lon := 0.0; lon < 360; lon += 0.3 {
		if h := casaDe(lon, c); h < 1 || h > 12 {
			t.Fatalf("%.1f° cae en la casa %d", lon, h)
		}
	}
}

// La carta occidental completa: que los planetas caigan en la casa que les toca
// según las cúspides, y que los regentes apunten a donde está su planeta.
func TestCartaOccidentalCoherente(t *testing.T) {
	c := Calcular(1961, 12, 19, 16, 30, 1, 41.58, 2.55)
	if len(c.Cuerpos) != 12 {
		t.Errorf("%d cuerpos, deberían ser 12 (diez planetas y los dos nodos)", len(c.Cuerpos))
	}
	pos := map[string]float64{}
	for _, b := range c.Cuerpos {
		pos[b.Nombre] = b.Lon
		if b.CasaP != casaDe(b.Lon, c.CuspP) {
			t.Errorf("%s dice casa %d y por su longitud le toca la %d",
				b.Nombre, b.CasaP, casaDe(b.Lon, c.CuspP))
		}
		if b.SignoIdx != int(b.Lon/30) {
			t.Errorf("%s dice el signo %d y está a %.2f°", b.Nombre, b.SignoIdx, b.Lon)
		}
	}
	// los regentes: el de cada casa está donde dice que está
	for i := 0; i < 12; i++ {
		r := c.Regentes[i]
		if r == "" {
			t.Errorf("la casa %d no tiene regente", i+1)
			continue
		}
		if c.RegenteEn[i] != casaDe(pos[r], c.CuspP) {
			t.Errorf("el regente de la casa %d (%s) dice estar en la %d y está en la %d",
				i+1, r, c.RegenteEn[i], casaDe(pos[r], c.CuspP))
		}
	}
	// los aspectos vienen ordenados por exactitud, que es como se leen
	for i := 1; i < len(c.Aspectos); i++ {
		if c.Aspectos[i].Orbe < c.Aspectos[i-1].Orbe {
			t.Errorf("los aspectos no están ordenados por orbe: %.2f después de %.2f",
				c.Aspectos[i].Orbe, c.Aspectos[i-1].Orbe)
			break
		}
	}
}

// El día juliano y su inverso tienen que cerrar el círculo, también antes de
// la reforma gregoriana.
func TestJulianoIdaYVuelta(t *testing.T) {
	casos := []struct {
		a, m, d int
		h       float64
	}{
		{2026, 8, 19, 15.5}, {2000, 1, 1, 12}, {1961, 12, 19, 15.5},
		{1582, 10, 15, 0}, {1500, 3, 1, 6}, {837, 4, 10, 7.2},
	}
	for _, c := range casos {
		jd := DiaJuliano(c.a, c.m, float64(c.d), c.h)
		a, m, d, h := DeDiaJuliano(jd)
		if a != c.a || m != c.m || int(d) != c.d || math.Abs(h-c.h) > 1e-4 {
			t.Errorf("%d-%02d-%02d %.2fh → jd %.6f → %d-%02d-%02.0f %.4fh",
				c.a, c.m, c.d, c.h, jd, a, m, d, h)
		}
	}
}
