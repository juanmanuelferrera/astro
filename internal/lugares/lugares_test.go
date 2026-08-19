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
