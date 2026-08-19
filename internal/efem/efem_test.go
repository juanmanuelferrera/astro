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
