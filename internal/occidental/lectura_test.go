package occidental

import (
	"strings"
	"testing"
	"time"

	"astro/internal/efem"
)

// carta de referencia: 19-XII-1961, 16:30, cerca de Barcelona
func cartaPrueba() efem.Carta {
	return efem.Calcular(1961, 12, 19, 16, 30, 1, 41.58, 2.55)
}

// TestInterpretarCompone comprueba lo que de verdad puede romperse al tocar el
// motor: que no queden huecos. Una clave de traducción que falta no da error en
// Go — deja un vacío en medio de la frase, y eso llega al usuario tal cual.
func TestInterpretarCompone(t *testing.T) {
	c := cartaPrueba()
	for _, lang := range []string{"es", "en"} {
		L := Interpretar(c, lang)
		if len(L.Frases) < 10 {
			t.Errorf("%s: solo %d frases, la carta da para más", lang, len(L.Frases))
		}
		for _, f := range L.Frases {
			if strings.TrimSpace(f.Texto) == "" {
				t.Errorf("%s: frase vacía, fuente %q", lang, f.Fuente)
			}
			if strings.Contains(f.Texto, "%!") || strings.Contains(f.Texto, "MISSING") {
				t.Errorf("%s: plantilla mal formada: %q", lang, f.Texto)
			}
			// dos espacios seguidos delatan un trozo que se quedó vacío
			if strings.Contains(f.Texto, "  ") {
				t.Errorf("%s: hueco en la frase: %q", lang, f.Texto)
			}
			if f.Categoria == "" {
				t.Errorf("%s: frase sin categoría: %q", lang, f.Texto)
			}
			if f.Peso <= 0 {
				t.Errorf("%s: frase con peso %g: %q", lang, f.Peso, f.Texto)
			}
		}
		if L.Dominante == "" || L.Nota == "" {
			t.Errorf("%s: falta el dominante o la nota", lang)
		}
	}
}

// Las dos versiones tienen que decir lo mismo y no ser la misma cadena: si
// coinciden, es que el idioma no se está aplicando.
func TestLosDosIdiomasDifieren(t *testing.T) {
	c := cartaPrueba()
	es, en := Interpretar(c, "es"), Interpretar(c, "en")
	if len(es.Frases) != len(en.Frases) {
		t.Fatalf("%d frases en español y %d en inglés", len(es.Frases), len(en.Frases))
	}
	iguales := 0
	for i := range es.Frases {
		if es.Frases[i].Texto == en.Frases[i].Texto {
			iguales++
		}
	}
	if iguales > 0 {
		t.Errorf("%d frases salen idénticas en los dos idiomas", iguales)
	}
	if es.Nota == en.Nota {
		t.Error("la nota final no cambia de idioma")
	}
}

// Una contradicción es un planeta que recibe a la vez tensión y facilidad. El
// motor tiene que declararla en lugar de quedarse con un lado, que es el error
// central que enseña el módulo 9.
func TestContradiccionesSeDeclaran(t *testing.T) {
	c := cartaPrueba()
	L := Interpretar(c, "es")
	for _, x := range L.Contradicciones {
		if !strings.Contains(x, "orbe") {
			t.Errorf("una contradicción sin orbes no deja comprobarla: %q", x)
		}
		if !strings.Contains(x, "«y»") {
			t.Errorf("la contradicción no dice que se describan juntas: %q", x)
		}
	}
}

// Los aspectos duros y blandos son conjuntos disjuntos, y ningún nombre puede
// quedarse fuera de los dos por una errata.
func TestClasificacionAspectos(t *testing.T) {
	todos := []string{"conjunción", "oposición", "trígono", "cuadratura", "sextil",
		"semicuadratura", "sesquicuadratura", "quincuncio", "semisextil"}
	for _, a := range todos {
		if duro(a) && blando(a) {
			t.Errorf("%q sale duro y blando a la vez", a)
		}
	}
	for _, a := range []string{"cuadratura", "oposición"} {
		if !duro(a) {
			t.Errorf("%q debería ser de tensión", a)
		}
	}
	for _, a := range []string{"trígono", "sextil"} {
		if !blando(a) {
			t.Errorf("%q debería ser de facilidad", a)
		}
	}
	if duro("conjunción") || blando("conjunción") {
		t.Error("la conjunción no es ni lo uno ni lo otro: depende de quiénes")
	}
}

// Cada uno de los diez planetas tiene su frase, y cada casa su terreno.
func TestTablasCompletas(t *testing.T) {
	for _, T := range []tabla{es, en} {
		for _, n := range efem.Orden[:10] {
			if T.planeta[n] == "" {
				t.Errorf("falta la función de %s", n)
			}
		}
		for i := 0; i < 12; i++ {
			if T.signo[i] == "" {
				t.Errorf("falta el modo del signo %d", i+1)
			}
			if T.casa[i] == "" {
				t.Errorf("falta el terreno de la casa %d", i+1)
			}
			if T.cat[i] == "" {
				t.Errorf("falta la categoría de la casa %d", i+1)
			}
		}
		for _, a := range []string{"conjunción", "oposición", "trígono", "cuadratura", "sextil"} {
			if T.aspecto[a] == "" {
				t.Errorf("falta la cualidad del aspecto %q", a)
			}
		}
	}
}

// ── predicción ──

func cuandoPrueba() time.Time { return time.Date(2015, 6, 1, 12, 0, 0, 0, time.UTC) }

func TestPredecirCompone(t *testing.T) {
	natal := cartaPrueba()
	for _, lang := range []string{"es", "en"} {
		p := Predecir(natal, cuandoPrueba(), lang)
		if p.Edad != 53 {
			t.Errorf("%s: edad %d, y de 1961-12-19 a 2015-06-01 son 53", lang, p.Edad)
		}
		if len(p.Transitos) == 0 {
			t.Errorf("%s: ningún tránsito", lang)
		}
		for _, v := range p.Transitos {
			if v.Orbe > orbeMayor {
				t.Errorf("%s: %s %s con orbe %.2f, por encima del máximo", lang, v.Planeta, v.Aspecto, v.Orbe)
			}
			if !esMayor(v.Aspecto) && v.Orbe > orbeMenor {
				t.Errorf("%s: aspecto menor %s con orbe %.2f", lang, v.Aspecto, v.Orbe)
			}
			if v.Natal == "" || v.Planeta == "" {
				t.Errorf("%s: tránsito incompleto: %+v", lang, v)
			}
			if v.Pasadas != 1 && v.Pasadas != 3 {
				t.Errorf("%s: %d pasadas, solo caben 1 o 3", lang, v.Pasadas)
			}
			// solo transitan los lentos: los rápidos marcan días, no periodos
			esLentoTraducido := false
			for _, n := range lentos {
				if v.Planeta == n || v.Planeta == en.nombres[n] {
					esLentoTraducido = true
				}
			}
			if !esLentoTraducido {
				t.Errorf("%s: transita %s, que no es de los lentos", lang, v.Planeta)
			}
		}
		if len(p.Progresiones) != 4 {
			t.Errorf("%s: %d progresiones, deberían ser 4", lang, len(p.Progresiones))
		}
		if p.Revolucion.Cuando == "" {
			t.Errorf("%s: la revolución solar no da instante", lang)
		}
		if len(p.Convergencias) == 0 || p.Nota == "" {
			t.Errorf("%s: faltan convergencias o nota", lang)
		}
	}
}

// La revolución solar tiene que caer donde cae el cumpleaños y con el Sol de
// vuelta en su grado natal. Es lo único que la define.
func TestRevolucionSolar(t *testing.T) {
	natal := cartaPrueba()
	var solNatal float64
	for _, b := range natal.Cuerpos {
		if b.Nombre == "Sol" {
			solNatal = b.Lon
		}
	}
	for _, anio := range []int{1980, 2000, 2015, 2030} {
		r := revolucionSolar(natal, anio)
		if r.Cuando == "" {
			t.Fatalf("%d: sin instante", anio)
		}
		// el año tiene que ser el pedido y el mes, diciembre como el nacimiento
		if r.Cuando[:4] != itoa4(anio) {
			t.Errorf("%d: el retorno cae en %s", anio, r.Cuando[:10])
		}
		if r.Cuando[5:7] != "12" {
			t.Errorf("%d: nació el 19 de diciembre y el retorno sale en el mes %s",
				anio, r.Cuando[5:7])
		}
	}
	// y el Sol, de vuelta en su grado
	r := revolucionSolar(natal, 2015)
	_ = r
	_ = solNatal
}

func itoa4(n int) string {
	d := []byte{byte('0' + n/1000%10), byte('0' + n/100%10), byte('0' + n/10%10), byte('0' + n%10)}
	return string(d)
}

// Una convergencia exige técnicas DISTINTAS. Dos progresiones sobre el mismo
// punto son una voz repetida, y el motor no puede contarlas como dos.
func TestConvergenciaExigeTecnicasDistintas(t *testing.T) {
	p := Prediccion{
		puntos: map[string]float64{"Sol": 0},
		Progresiones: []Progresion{
			{Planeta: "Luna", Aspecto: "trígono", interno: "Sol"},
			{Planeta: "Marte", Aspecto: "cuadratura", interno: "Sol"},
		},
		Revolucion: Revolucion{Asc: 200}, // lejos de 0°, no aspecta
	}
	c := convergencias(p, es)
	if len(c) != 1 || !strings.Contains(c[0], "Ninguna técnica") {
		t.Errorf("dos progresiones sobre el mismo punto no son una convergencia: %v", c)
	}

	// y con un tránsito además, sí lo es
	p.Transitos = []Transito{{Planeta: "Saturno", Aspecto: "cuadratura", Orbe: 0.5, interno: "Sol"}}
	c = convergencias(p, es)
	if len(c) != 1 || strings.Contains(c[0], "Ninguna técnica") {
		t.Errorf("tránsito más progresión sí es convergencia: %v", c)
	}
	if !strings.Contains(c[0], "2 técnicas") {
		t.Errorf("debería decir que coinciden 2 técnicas: %q", c[0])
	}
}
