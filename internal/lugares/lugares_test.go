package lugares

import "testing"

// Los datos van embebidos con go:embed, así que un fichero mal formado no da
// error al compilar: deja la búsqueda vacía y nadie se entera hasta usarla.
func TestHayDatos(t *testing.T) {
	r := Buscar("barcelona", 8)
	if len(r) == 0 {
		t.Fatal("no se encuentra Barcelona: los datos no se han cargado")
	}
}

func TestBusquedasCorrientes(t *testing.T) {
	casos := []struct {
		q                string
		latMin, latMax   float64
		lonMin, lonMax   float64
	}{
		{"barcelona", 41.2, 41.6, 1.9, 2.4},
		{"madrid", 40.2, 40.6, -3.9, -3.5},
		{"london", 51.3, 51.7, -0.3, 0.1},
		{"tokyo", 35.5, 35.9, 139.5, 139.9},
	}
	for _, c := range casos {
		r := Buscar(c.q, 8)
		if len(r) == 0 {
			t.Errorf("%q no da ningún resultado", c.q)
			continue
		}
		l := r[0]
		if l.Lat < c.latMin || l.Lat > c.latMax || l.Lon < c.lonMin || l.Lon > c.lonMax {
			t.Errorf("%q da %.3f, %.3f, que no cae donde debería", c.q, l.Lat, l.Lon)
		}
		if l.Zona == "" {
			t.Errorf("%q sale sin huso horario", c.q)
		}
	}
}

// Sin acentos y con distinta caja tiene que encontrar lo mismo: nadie escribe
// «Córdoba» con tilde en un buscador.
func TestSinAcentosYSinCaja(t *testing.T) {
	for _, par := range [][2]string{{"córdoba", "cordoba"}, {"MÁLAGA", "malaga"}, {"Gijón", "gijon"}} {
		a, b := Buscar(par[0], 8), Buscar(par[1], 8)
		if len(a) == 0 || len(b) == 0 {
			t.Errorf("%q o %q no dan nada", par[0], par[1])
			continue
		}
		if a[0].Nombre != b[0].Nombre {
			t.Errorf("%q da %q y %q da %q", par[0], a[0].Nombre, par[1], b[0].Nombre)
		}
	}
}

func TestNoRevientaConBasura(t *testing.T) {
	for _, q := range []string{"", " ", "a", "zzzzzzzzzz", "'; DROP TABLE", "日本"} {
		_ = Buscar(q, 8) // basta con que no entre en pánico
	}
}

// El huso tiene que salir para cualquier zona conocida y en cualquier época.
func TestHuso(t *testing.T) {
	r := Buscar("barcelona", 8)
	if len(r) == 0 {
		t.Skip("sin datos")
	}
	z := r[0].Zona
	for _, c := range []struct {
		anio, mes, dia int
		esperado       float64
	}{
		{2026, 1, 15, 1}, // invierno: hora central europea
		{2026, 7, 15, 2}, // verano
		{1961, 12, 19, 1},
	} {
		off, _, _, err := Huso(z, c.anio, c.mes, c.dia, 12, 0)
		if err != nil {
			t.Errorf("%s no se reconoce como huso: %v", z, err)
			continue
		}
		if off != c.esperado {
			t.Errorf("%s el %d-%02d-%02d da UTC%+g y debería ser UTC%+g",
				z, c.anio, c.mes, c.dia, off, c.esperado)
		}
	}
}

// La historia del huso es lo que contesta «¿por qué este desfase?». Es la parte
// que más confunde a la gente y no la comprobaba nada.
func TestHistoriaHuso(t *testing.T) {
	// España lleva hora de Berlín desde 1940, así que en diciembre de 1961 el
	// reloj iba una hora por delante del Sol de Barcelona, y algo más.
	h, err := HistoriaHuso("Europe/Madrid", 1961, 12, 19, 16, 30, 2.55)
	if err != nil {
		t.Fatalf("no calcula la historia: %v", err)
	}
	if h.Offset != 1 {
		t.Errorf("en diciembre de 1961 España iba en UTC+1, y da UTC%+g", h.Offset)
	}
	if h.Verano {
		t.Error("en diciembre no hay horario de verano")
	}
	// el desfase con el Sol: Barcelona está a 2,55° este, o sea +10 minutos de
	// sol, mientras el reloj marca +1 hora
	if h.Solar < 0.1 || h.Solar > 0.3 {
		t.Errorf("la hora solar de 2,55° este debería rondar UTC+0,17 y da %+g", h.Solar)
	}
	if h.Desfase < 0.7 || h.Desfase > 0.9 {
		t.Errorf("el desfase reloj-Sol debería rondar los 50 minutos y da %+g h", h.Desfase)
	}
	if h.Zona != "Europe/Madrid" || h.Abrev == "" {
		t.Errorf("falta zona o abreviatura: %+v", h)
	}

	// En verano tiene que detectarlo, y ese año tiene que haber cambios de hora.
	v, err := HistoriaHuso("Europe/Madrid", 2020, 7, 15, 12, 0, 2.55)
	if err != nil {
		t.Fatal(err)
	}
	if !v.Verano || v.Offset != 2 {
		t.Errorf("julio de 2020 en Madrid es UTC+2 con horario de verano, y da UTC%+g verano=%v",
			v.Offset, v.Verano)
	}
	if len(v.DelAnio) != 2 {
		t.Errorf("2020 tuvo dos cambios de hora y encuentra %d", len(v.DelAnio))
	}

	// Un año sin cambios de hora no debe inventarse ninguno.
	q, err := HistoriaHuso("Europe/Madrid", 1950, 6, 15, 12, 0, 2.55)
	if err != nil {
		t.Fatal(err)
	}
	if len(q.DelAnio) != 0 {
		t.Errorf("en 1950 España no cambiaba la hora y encuentra %d cambios: %+v",
			len(q.DelAnio), q.DelAnio)
	}

	if _, err := HistoriaHuso("No/Existe", 2000, 1, 1, 12, 0, 0); err == nil {
		t.Error("una zona inventada debería dar error")
	}
}
