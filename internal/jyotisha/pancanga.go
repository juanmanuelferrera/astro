package jyotisha

import (
	"math"

	"astro/internal/efem"
)

// El pañcāṅga — «cinco miembros» — es el calendario hindú: tithi, vāra,
// nakṣatra, yoga y karaṇa. No describe a la persona, describe el día. Se usa
// para elegir momentos (muhūrta) y para fechar, y es lo que cuelga en la pared
// de cualquier casa india.
//
// Nota sobre el ayanāṁśa: el tithi y el karaṇa salen de la DIFERENCIA entre la
// Luna y el Sol, así que el ayanāṁśa se cancela y da igual el zodíaco. El yoga
// sale de la SUMA, así que no se cancela: se calcula sobre longitudes sidéreas.

var tithis = [15]string{"Pratipat", "Dvitīyā", "Tṛtīyā", "Caturthī", "Pañcamī",
	"Ṣaṣṭhī", "Saptamī", "Aṣṭamī", "Navamī", "Daśamī", "Ekādaśī", "Dvādaśī",
	"Trayodaśī", "Caturdaśī", "Pūrṇimā"}

var yogasPanca = [27]string{"Viṣkambha", "Prīti", "Āyuṣmān", "Saubhāgya", "Śobhana",
	"Atigaṇḍa", "Sukarma", "Dhṛti", "Śūla", "Gaṇḍa", "Vṛddhi", "Dhruva",
	"Vyāghāta", "Harṣaṇa", "Vajra", "Siddhi", "Vyatīpāta", "Varīyān", "Parigha",
	"Śiva", "Siddha", "Sādhya", "Śubha", "Śukla", "Brahma", "Indra", "Vaidhṛti"}

// Los siete karaṇas móviles giran en rueda; los cuatro fijos aparecen una vez
// por mes lunar, al principio y al final.
var karanasCara = [7]string{"Bava", "Bālava", "Kaulava", "Taitila", "Gara", "Vaṇija", "Viṣṭi"}

var varas = [7]string{"Ravivāra", "Somavāra", "Maṅgalavāra", "Budhavāra",
	"Guruvāra", "Śukravāra", "Śanivāra"}
var senorVara = [7]string{"Sol", "Luna", "Marte", "Mercurio", "Júpiter", "Venus", "Saturno"}

type Pancanga struct {
	Tithi     string  `json:"tithi"`
	TithiNum  int     `json:"tithiNum"`  // 1-30 dentro del mes lunar
	Paksha    string  `json:"paksha"`    // śukla (creciente) o kṛṣṇa (menguante)
	TithiPct  float64 `json:"tithiPct"`  // cuánto lleva recorrido, 0-100
	Vara      string  `json:"vara"`
	SenorVara string  `json:"senorVara"`
	Nakshatra string  `json:"nakshatra"`
	Pada      int     `json:"pada"`
	SenorNak  string  `json:"senorNak"`
	NakPct    float64 `json:"nakPct"`
	Yoga      string  `json:"yoga"`
	YogaNum   int     `json:"yogaNum"`
	Karana    string  `json:"karana"`
	Visti     bool    `json:"visti"`  // Viṣṭi/Bhadrā: el karaṇa que se evita
	Luna      float64 `json:"luna"`   // fase, 0-100
}

// CalcPancanga arma los cinco miembros. lonSol y lonLuna son sidéreas.
func CalcPancanga(jd, lonSol, lonLuna float64) Pancanga {
	var p Pancanga

	// Tithi: cada uno son 12° que la Luna le saca al Sol.
	dif := math.Mod(lonLuna-lonSol+360, 360)
	n := int(dif / 12)
	p.TithiNum = n + 1
	p.TithiPct = math.Mod(dif, 12) / 12 * 100
	if n < 15 {
		p.Paksha = "śukla"
		p.Tithi = tithis[n]
	} else {
		p.Paksha = "kṛṣṇa"
		if n == 29 {
			p.Tithi = "Amāvāsyā"
		} else {
			p.Tithi = tithis[n-15]
		}
	}
	p.Luna = dif / 360 * 100

	// Vāra: el día de la semana. Cuenta desde el mediodía juliano; en la
	// tradición el día empieza al amanecer, no a medianoche.
	d := int(math.Floor(jd+1.5)) % 7
	if d < 0 {
		d += 7
	}
	p.Vara, p.SenorVara = varas[d], senorVara[d]

	// Nakṣatra de la Luna.
	p.Nakshatra, p.Pada, p.SenorNak = Nakshatra(lonLuna)
	const span = 360.0 / 27.0
	p.NakPct = math.Mod(lonLuna, span) / span * 100

	// Yoga: la suma de los dos, repartida en 27 tramos de 13°20′.
	suma := math.Mod(lonSol+lonLuna, 360)
	y := int(suma / span)
	if y > 26 {
		y = 26
	}
	p.YogaNum, p.Yoga = y+1, yogasPanca[y]

	// Karaṇa: medio tithi. Sesenta por mes lunar, y de ellos cuatro son fijos.
	k := int(dif / 6)
	switch {
	case k == 0:
		p.Karana = "Kiṁstughna"
	case k == 57:
		p.Karana = "Śakuni"
	case k == 58:
		p.Karana = "Catuṣpāda"
	case k == 59:
		p.Karana = "Nāga"
	default:
		p.Karana = karanasCara[(k-1)%7]
	}
	p.Visti = p.Karana == "Viṣṭi"

	return p
}

var _ = efem.Signos
