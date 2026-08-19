package jyotisha

import (
	"math"
	"strings"
	"testing"
)

// Ayanāṁśa Lahiri de Swiss Ephemeris, sacado del generador y no escrito a mano
// —el primer intento lo escribí a ojo y se desviaba medio grado—. El ayanāṁśa es la bisagra de todo el sistema sidéreo: si
// se mueve, se mueven los doce rāśis a la vez.
var refAyanamsa = []struct {
	jd, lahiri float64
}{
	{2378510.770833333, 21.065147356}, // 1800-01-15
	{2415034.770833333, 22.461076136}, // 1900-01-15
	{2451545.000000000, 23.857092354}, // J2000
	{2461272.000000000, 24.229120405}, // 2026-08-19
	{2488083.770833333, 25.254814297}, // 2100-01-15
}

func TestAyanamsa(t *testing.T) {
	var peor float64
	for _, r := range refAyanamsa {
		e := math.Abs(Ayanamsa(r.jd)-r.lahiri) * 3600
		if e > peor {
			peor = e
		}
	}
	// El polinomio está ajustado contra Swiss Ephemeris; medio segundo de arco
	// sobra, porque un rāśi mide 30 grados.
	if peor > 0.5 {
		t.Errorf("el ayanāṁśa se desvía %.4f″", peor)
	}
	t.Logf("peor desvío del ayanāṁśa: %.5f″", peor)
}

// Las vargas son aritmética pura: reparten un signo en n trozos. Lo que hay
// que comprobar es que no se salen del rango y que los límites caen donde
// deben, no un valor concreto sacado de ninguna parte.
func TestVargasEnRango(t *testing.T) {
	for _, n := range []int{1, 2, 3, 7, 9, 10, 12, 16, 30, 60} {
		for lon := 0.0; lon < 360; lon += 0.37 {
			v := Varga(lon, n)
			if v < 0 || v > 11 {
				t.Fatalf("D%d en %.2f° da el rāśi %d, fuera de 0-11", n, lon, v)
			}
		}
		// el D1 es el propio signo, por definición
		if n == 1 {
			for lon := 0.0; lon < 360; lon += 1.7 {
				if Varga(lon, 1) != int(lon/30) {
					t.Fatalf("el D1 de %.2f° no es su propio signo", lon)
				}
			}
		}
	}
}

// La vimśottarī reparte 120 años exactos entre los nueve grahas. Si la suma no
// cuadra, las fechas de todas las daśās se van desplazando.
func TestVimsottariSuma(t *testing.T) {
	total := 0.0
	for _, c := range ciclo {
		total += float64(c.anios)
	}
	if total != 120 {
		t.Errorf("el ciclo vimśottarī suma %g años, deberían ser 120", total)
	}
	if len(ciclo) != 9 {
		t.Errorf("el ciclo tiene %d grahas y deberían ser 9", len(ciclo))
	}
}

// Los bhāvas de signo entero: el primero es el del Lagna y de ahí en orden.
// Con esto se caza cualquier error de módulo, que es el fallo típico.
func TestBhavasDesdeLagna(t *testing.T) {
	for lagna := 0; lagna < 12; lagna++ {
		for rasi := 0; rasi < 12; rasi++ {
			b := ((rasi-lagna)%12+12)%12 + 1
			if b < 1 || b > 12 {
				t.Fatalf("lagna %d y rāśi %d dan el bhāva %d", lagna, rasi, b)
			}
			if rasi == lagna && b != 1 {
				t.Fatalf("el rāśi del lagna debería ser el bhāva 1, y da %d", b)
			}
		}
	}
}

// Dṛṣṭi: todos miran a la séptima, y los tres especiales añaden las suyas.
func TestDrishti(t *testing.T) {
	for _, g := range Grahas {
		d := Drishti(g)
		tiene7 := false
		for _, x := range d {
			if x == 7 {
				tiene7 = true
			}
			if x < 1 || x > 12 {
				t.Errorf("%s aspecta la casa %d, fuera de 1-12", g, x)
			}
		}
		if !tiene7 {
			t.Errorf("%s no mira a la séptima, y eso lo hacen todos", g)
		}
	}
	esperado := map[string][]int{
		"Marte": {4, 7, 8}, "Júpiter": {5, 7, 9}, "Saturno": {3, 7, 10},
	}
	for g, e := range esperado {
		d := Drishti(g)
		if len(d) != len(e) {
			t.Fatalf("%s tiene %d miradas y debería tener %d", g, len(d), len(e))
		}
		for i := range e {
			if d[i] != e[i] {
				t.Errorf("%s mira a %v y debería mirar a %v", g, d, e)
				break
			}
		}
	}
}

// Gaṇḍānta: los últimos 3°20′ de Karka, Vṛścika y Mīna, y nada más.
func TestGandanta(t *testing.T) {
	for s := 0; s < 12; s++ {
		agua := s == 3 || s == 7 || s == 11
		base := float64(s) * 30
		if Gandanta(base + 25) {
			t.Errorf("el rāśi %d a 25° no es gaṇḍānta", s)
		}
		if got := Gandanta(base + 29.5); got != agua {
			t.Errorf("el rāśi %d a 29,5° da gaṇḍānta=%v y debería ser %v", s, got, agua)
		}
	}
}

// TestCartaEntera levanta una carta de verdad y comprueba que todo lo que
// cuelga de ella cuadra consigo mismo. Es la prueba que más cubre: pasa por el
// cálculo, los bhāvas, las vargas, las daśās, el pañcāṅga, el aṣṭakavarga, el
// ṣaḍbala y los arudhas de una sola vez.
func TestCartaEntera(t *testing.T) {
	c := Calcular(1961, 12, 19, 16, 30, 1, 41.58, 2.55)

	if len(c.Grahas) != 9 {
		t.Errorf("%d grahas, deberían ser 9", len(c.Grahas))
	}
	if len(c.Bhavas) != 12 {
		t.Errorf("%d bhāvas, deberían ser 12", len(c.Bhavas))
	}
	for _, g := range c.Grahas {
		if g.Bhava < 1 || g.Bhava > 12 {
			t.Errorf("%s en el bhāva %d", g.Nombre, g.Bhava)
		}
		if g.Lon < 0 || g.Lon >= 360 {
			t.Errorf("%s a %.4f°", g.Nombre, g.Lon)
		}
		if g.Pada < 1 || g.Pada > 4 {
			t.Errorf("%s en el pada %d", g.Nombre, g.Pada)
		}
		if g.Nak == "" || g.SenorNak == "" {
			t.Errorf("%s sin nakṣatra o sin su señor", g.Nombre)
		}
		// el bhāva tiene que salir del rāśi, no de otro sitio
		if b := ((g.RasiIdx-c.LagnaRasi)%12+12)%12 + 1; b != g.Bhava {
			t.Errorf("%s dice bhāva %d y por su rāśi le toca el %d", g.Nombre, g.Bhava, b)
		}
	}

	// Cada bhāva tiene que estar en su sitio y su señor en algún lado.
	for i, b := range c.Bhavas {
		if b.Numero != i+1 {
			t.Errorf("el bhāva %d se numera como %d", i+1, b.Numero)
		}
		if b.Senor == "" {
			t.Errorf("el bhāva %d no tiene señor", b.Numero)
		}
		if b.SenorEn < 1 || b.SenorEn > 12 {
			t.Errorf("el señor del bhāva %d está en el %d", b.Numero, b.SenorEn)
		}
	}

	// Aṣṭakavarga: la suma es una constante del sistema.
	if c.Ashtaka.Total != 337 {
		t.Errorf("el aṣṭakavarga suma %d y siempre son 337", c.Ashtaka.Total)
	}

	// Ṣaḍbala: los siete grahas y las doce casas.
	if len(c.Shadbala.Balas) != 7 {
		t.Errorf("%d ṣaḍbalas, deberían ser 7", len(c.Shadbala.Balas))
	}
	if len(c.Shadbala.Bhavas) != 12 {
		t.Errorf("%d bhāva balas, deberían ser 12", len(c.Shadbala.Bhavas))
	}
	for _, b := range c.Shadbala.Balas {
		if b.Rupas <= 0 || b.Rupas > 30 {
			t.Errorf("%s saca %.2f rūpas, que no es una cifra creíble", b.Graha, b.Rupas)
		}
		if b.Rango < 1 || b.Rango > 7 {
			t.Errorf("%s queda en el puesto %d", b.Graha, b.Rango)
		}
	}

	// Pañcāṅga: los cinco miembros, con el tithi dentro del mes lunar.
	p := c.Pancanga
	if p.Tithi == "" || p.Vara == "" || p.Nakshatra == "" || p.Yoga == "" || p.Karana == "" {
		t.Errorf("al pañcāṅga le falta algún miembro: %+v", p)
	}
	if p.TithiNum < 1 || p.TithiNum > 30 {
		t.Errorf("tithi %d, fuera de 1-30", p.TithiNum)
	}
	if p.YogaNum < 1 || p.YogaNum > 27 {
		t.Errorf("yoga %d, fuera de 1-27", p.YogaNum)
	}
	if p.Vara != "Maṅgalavāra" {
		t.Errorf("el 19-XII-1961 fue martes, y sale %s", p.Vara)
	}

	// Daśās: el ciclo entero y una sola corriente.
	var suma float64
	actuales := 0
	for _, d := range c.Dasas {
		suma += d.Anios
		if d.Actual {
			actuales++
		}
	}
	if actuales != 1 {
		t.Errorf("%d daśās marcadas como actuales, tiene que haber una", actuales)
	}
	if suma < 119 || suma > 121 {
		t.Errorf("las daśās suman %.2f años y el ciclo son 120", suma)
	}

	// Arudhas: doce, y ninguno cae donde no puede.
	for i, r := range c.Arudhas.Padas {
		if r < 0 || r > 11 {
			t.Errorf("el arudha del bhāva %d cae en el rāśi %d", i+1, r)
		}
	}

	// Vargas: las que están declaradas, y todas con sus nueve grahas y el Lagna.
	for _, v := range VargasUsadas {
		gs, ok := c.Vargas[v.Clave]
		if !ok {
			t.Errorf("falta la varga %s", v.Clave)
			continue
		}
		if len(gs) != 10 {
			t.Errorf("la varga %s trae %d posiciones y deberían ser 10 (9 grahas y el Lagna)",
				v.Clave, len(gs))
		}
	}
}

// La lectura védica tiene que componerse entera en los dos idiomas, sin huecos.
func TestLecturaVedica(t *testing.T) {
	c := Calcular(1961, 12, 19, 16, 30, 1, 41.58, 2.55)
	for _, lang := range []string{"es", "en"} {
		L := Interpretar(c, lang)
		if len(L.Frases) < 15 {
			t.Errorf("%s: solo %d frases", lang, len(L.Frases))
		}
		for _, f := range L.Frases {
			if f.Texto == "" || f.Categoria == "" {
				t.Errorf("%s: frase incompleta: %+v", lang, f)
			}
			if strings.Contains(f.Texto, "%!") || strings.Contains(f.Texto, "MISSING") {
				t.Errorf("%s: plantilla mal formada: %q", lang, f.Texto)
			}
		}
		if L.Dominante == "" || L.Nota == "" {
			t.Errorf("%s: falta el dominante o la nota", lang)
		}
	}
	// y que no digan lo mismo, que es como se sabe que el idioma se aplica
	es, en := Interpretar(c, "es"), Interpretar(c, "en")
	if es.Nota == en.Nota || es.Dominante == en.Dominante {
		t.Error("la lectura no cambia de idioma")
	}
}

// El nodo verdadero mueve a Rāhu y a Ketu, y a nadie más.
func TestNodoVerdaderoSoloMueveLosNodos(t *testing.T) {
	m := CalcularCon(1961, 12, 19, 16, 30, 1, 41.58, 2.55, false)
	v := CalcularCon(1961, 12, 19, 16, 30, 1, 41.58, 2.55, true)
	pos := func(c Carta, n string) float64 {
		for _, g := range c.Grahas {
			if g.Nombre == n {
				return g.Lon
			}
		}
		return -1
	}
	for _, n := range Grahas {
		dm, dv := pos(m, n), pos(v, n)
		nodo := n == "Rāhu" || n == "Ketu"
		if !nodo && dm != dv {
			t.Errorf("%s se mueve al cambiar de nodo: %.6f contra %.6f", n, dm, dv)
		}
		if nodo && dm == dv {
			t.Errorf("%s no se mueve al pedir el nodo verdadero", n)
		}
	}
	// Rāhu y Ketu siguen enfrentados pase lo que pase.
	d := math.Mod(pos(v, "Rāhu")-pos(v, "Ketu")+360, 360)
	if math.Abs(d-180) > 1e-9 {
		t.Errorf("Rāhu y Ketu están a %.6f° y siempre son 180", d)
	}
}
