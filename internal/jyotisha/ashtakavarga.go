package jyotisha

// Aṣṭakavarga — «ocho divisiones». Es el único sistema de jyotiṣa que da un
// número por casa en lugar de un juicio, y por eso se usa para decidir entre
// lecturas que se contradicen.
//
// Cómo funciona: para cada graha se mira, desde cada uno de ocho puntos de
// referencia (los siete grahas más el Lagna), en qué casas contadas desde ese
// punto aporta un bindu — un punto benéfico. Sumando se obtiene el
// bhinnāṣṭakavarga (BAV) de ese graha: 12 casillas, de 0 a 8 puntos.
//
// Sumando los siete BAV se obtiene el sarvāṣṭakavarga (SAV): cuánto apoyo
// total tiene cada signo. La media es 337/12 ≈ 28. Por encima de 30 el asunto
// de esa casa va rodado; por debajo de 25 cuesta.
//
// Las tablas son las de Parāśara. La suma de los siete BAV tiene que dar
// exactamente 337: eso lo comprueba TestSuma más abajo, y es lo que garantiza
// que no se ha colado una errata al teclearlas.

// benefico[graha][referencia] = casas, contadas desde la referencia, donde el
// graha recibe un bindu. La referencia incluye al Lagna en última posición.
var refs = []string{"Sol", "Luna", "Marte", "Mercurio", "Júpiter", "Venus", "Saturno", "Lagna"}

var bindus = map[string]map[string][]int{
	"Sol": {
		"Sol": {1, 2, 4, 7, 8, 9, 10, 11}, "Luna": {3, 6, 10, 11},
		"Marte": {1, 2, 4, 7, 8, 9, 10, 11}, "Mercurio": {3, 5, 6, 9, 10, 11, 12},
		"Júpiter": {5, 6, 9, 11}, "Venus": {6, 7, 12},
		"Saturno": {1, 2, 4, 7, 8, 9, 10, 11}, "Lagna": {3, 4, 6, 10, 11, 12},
	},
	"Luna": {
		"Sol": {3, 6, 7, 8, 10, 11}, "Luna": {1, 3, 6, 7, 10, 11},
		"Marte": {2, 3, 5, 6, 9, 10, 11}, "Mercurio": {1, 3, 4, 5, 7, 8, 10, 11},
		"Júpiter": {1, 4, 7, 8, 10, 11, 12}, "Venus": {3, 4, 5, 7, 9, 10, 11},
		"Saturno": {3, 5, 6, 11}, "Lagna": {3, 6, 10, 11},
	},
	"Marte": {
		"Sol": {3, 5, 6, 10, 11}, "Luna": {3, 6, 11},
		"Marte": {1, 2, 4, 7, 8, 10, 11}, "Mercurio": {3, 5, 6, 11},
		"Júpiter": {6, 10, 11, 12}, "Venus": {6, 8, 11, 12},
		"Saturno": {1, 4, 7, 8, 9, 10, 11}, "Lagna": {1, 3, 6, 10, 11},
	},
	"Mercurio": {
		"Sol": {5, 6, 9, 11, 12}, "Luna": {2, 4, 6, 8, 10, 11},
		"Marte": {1, 2, 4, 7, 8, 9, 10, 11}, "Mercurio": {1, 3, 5, 6, 9, 10, 11, 12},
		"Júpiter": {6, 8, 11, 12}, "Venus": {1, 2, 3, 4, 5, 8, 9, 11},
		"Saturno": {1, 2, 4, 7, 8, 9, 10, 11}, "Lagna": {1, 2, 4, 6, 8, 10, 11},
	},
	"Júpiter": {
		"Sol": {1, 2, 3, 4, 7, 8, 9, 10, 11}, "Luna": {2, 5, 7, 9, 11},
		"Marte": {1, 2, 4, 7, 8, 10, 11}, "Mercurio": {1, 2, 4, 5, 6, 9, 10, 11},
		"Júpiter": {1, 2, 3, 4, 7, 8, 10, 11}, "Venus": {2, 5, 6, 9, 10, 11},
		"Saturno": {3, 5, 6, 12}, "Lagna": {1, 2, 4, 5, 6, 7, 9, 10, 11},
	},
	"Venus": {
		"Sol": {8, 11, 12}, "Luna": {1, 2, 3, 4, 5, 8, 9, 11, 12},
		"Marte": {3, 4, 6, 9, 11, 12}, "Mercurio": {3, 5, 6, 9, 11},
		"Júpiter": {5, 8, 9, 10, 11}, "Venus": {1, 2, 3, 4, 5, 8, 9, 10, 11},
		"Saturno": {3, 4, 5, 8, 9, 10, 11}, "Lagna": {1, 2, 3, 4, 5, 8, 9, 11},
	},
	"Saturno": {
		"Sol": {1, 2, 4, 7, 8, 10, 11}, "Luna": {3, 6, 11},
		"Marte": {3, 5, 6, 10, 11, 12}, "Mercurio": {6, 8, 9, 10, 11, 12},
		"Júpiter": {5, 6, 11, 12}, "Venus": {6, 11, 12},
		"Saturno": {3, 5, 6, 11}, "Lagna": {1, 3, 4, 6, 10, 11},
	},
}

// GrahasAV son los siete que tienen aṣṭakavarga. Rāhu y Ketu no lo tienen:
// no son cuerpos y el sistema no los contempla.
var GrahasAV = []string{"Sol", "Luna", "Marte", "Mercurio", "Júpiter", "Venus", "Saturno"}

type Ashtakavarga struct {
	BAV   map[string][12]int `json:"bav"`   // por graha, bindus en cada rāśi
	SAV   [12]int            `json:"sav"`   // la suma de los siete
	Total int                `json:"total"` // tiene que ser 337
	Media float64            `json:"media"` // 337/12
}

// CalcAshtakavarga necesita el rāśi de cada referencia, Lagna incluido.
func CalcAshtakavarga(rasiDe map[string]int) Ashtakavarga {
	a := Ashtakavarga{BAV: map[string][12]int{}, Media: 337.0 / 12.0}
	for _, g := range GrahasAV {
		var fila [12]int
		for _, r := range refs {
			base, ok := rasiDe[r]
			if !ok {
				continue
			}
			for _, casa := range bindus[g][r] {
				// casa 1 es el propio signo de la referencia, así que se resta 1
				fila[(base+casa-1)%12]++
			}
		}
		a.BAV[g] = fila
		for i := 0; i < 12; i++ {
			a.SAV[i] += fila[i]
			a.Total += fila[i]
		}
	}
	return a
}

// SumaTablas devuelve cuántos bindus tiene cada graha en sus ocho tablas.
// Los valores clásicos son Sol 48, Luna 49, Marte 39, Mercurio 54,
// Júpiter 56, Venus 52 y Saturno 39, y suman 337.
func SumaTablas() (map[string]int, int) {
	m, total := map[string]int{}, 0
	for _, g := range GrahasAV {
		for _, r := range refs {
			m[g] += len(bindus[g][r])
		}
		total += m[g]
	}
	return m, total
}
