package occidental

import (
	"strings"
	"testing"

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
