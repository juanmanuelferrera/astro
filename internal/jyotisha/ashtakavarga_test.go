package jyotisha

import "testing"

// Las ocho tablas de cada graha suman una cifra fija que trae la tradición.
// Si al teclearlas se coló un número de más o de menos, esto lo caza.
func TestSumaTablas(t *testing.T) {
	esperado := map[string]int{"Sol": 48, "Luna": 49, "Marte": 39,
		"Mercurio": 54, "Júpiter": 56, "Venus": 52, "Saturno": 39}
	m, total := SumaTablas()
	for g, n := range esperado {
		if m[g] != n {
			t.Errorf("%s suma %d bindus, deberían ser %d", g, m[g], n)
		}
	}
	if total != 337 {
		t.Errorf("el total es %d, debería ser 337", total)
	}
}

// Ninguna casa puede pasar de 8: hay ocho referencias y cada una aporta como
// mucho un bindu por casa.
func TestMaximoOcho(t *testing.T) {
	rasi := map[string]int{"Sol": 8, "Luna": 1, "Marte": 7, "Mercurio": 8,
		"Júpiter": 10, "Venus": 9, "Saturno": 8, "Lagna": 1}
	a := CalcAshtakavarga(rasi)
	for g, fila := range a.BAV {
		for i, n := range fila {
			if n > 8 || n < 0 {
				t.Errorf("%s en el rāśi %d tiene %d bindus, fuera de 0-8", g, i+1, n)
			}
		}
	}
	if a.Total != 337 {
		t.Errorf("el SAV suma %d, debería ser 337", a.Total)
	}
}

// El arudha nunca puede quedarse en su propio bhāva ni en el séptimo: esa es
// justamente la excepción de la regla.
func TestArudhaNoCaeEnSiMismo(t *testing.T) {
	for b := 0; b < 12; b++ {
		for s := 0; s < 12; s++ {
			a := ArudhaDe(b, s)
			d := ((a-b)%12 + 12) % 12
			if d == 0 || d == 6 {
				t.Errorf("bhāva en rāśi %d con señor en %d da arudha en %d, a %d casas: prohibido",
					b, s, a, d)
			}
		}
	}
}
