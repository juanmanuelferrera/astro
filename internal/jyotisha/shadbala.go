package jyotisha

import "math"

// Ṣaḍbala — «seis fuerzas». Es el intento clásico de poner un número a lo que
// hasta ahora se venía juzgando a ojo: cuánto puede de verdad un graha.
//
// Se mide en virūpas; 60 virūpas son un rūpa. Cada graha tiene un mínimo
// exigido, y lo que cuenta no es la cifra bruta sino la RAZÓN entre lo que
// saca y lo que se le pide. Un Saturno con 5 rūpas está fuerte; un Sol con
// esos mismos 5 rūpas se queda corto.
//
// Las seis fuentes:
//   1. Sthāna   — de dónde está: exaltación, dignidad en siete vargas, par o
//                 impar, kendra, y en qué tercio del signo
//   2. Dig      — de la dirección: cada graha rinde en un punto del cielo
//   3. Kāla     — del tiempo: día o noche, quincena, año, mes, día y hora
//   4. Cheṣṭā   — del movimiento: retrógrado, rápido, lento o parado
//   5. Naisargika — la que trae de fábrica, fija y siempre la misma
//   6. Dṛk      — de quién le mira, con signo: los benéficos suman
//
// Lo que aquí NO está, y conviene decirlo: el yuddha bala (la guerra entre dos
// grahas a menos de un grado) y el bhāva bala (la fuerza de las casas, que es
// otro cálculo aparte). Todo lo demás está entero.

const virupasPorRupa = 60

// mínimos exigidos a cada graha, en rūpas
var minimoRupas = map[string]float64{
	"Sol": 5, "Luna": 6, "Marte": 5, "Mercurio": 7,
	"Júpiter": 6.5, "Venus": 5.5, "Saturno": 5,
}

// naisargika: la fuerza natural, por brillo. Fija.
var naisargika = map[string]float64{
	"Sol": 60, "Luna": 51.43, "Venus": 42.86, "Júpiter": 34.29,
	"Mercurio": 25.71, "Marte": 17.14, "Saturno": 8.57,
}

// amistades naturales. 1 amigo, 0 neutro, -1 enemigo.
var amistadNat = map[string]map[string]int{
	"Sol":      {"Luna": 1, "Marte": 1, "Júpiter": 1, "Mercurio": 0, "Venus": -1, "Saturno": -1},
	"Luna":     {"Sol": 1, "Mercurio": 1, "Marte": 0, "Júpiter": 0, "Venus": 0, "Saturno": 0},
	"Marte":    {"Sol": 1, "Luna": 1, "Júpiter": 1, "Venus": 0, "Saturno": 0, "Mercurio": -1},
	"Mercurio": {"Sol": 1, "Venus": 1, "Marte": 0, "Júpiter": 0, "Saturno": 0, "Luna": -1},
	"Júpiter":  {"Sol": 1, "Luna": 1, "Marte": 1, "Saturno": 0, "Mercurio": -1, "Venus": -1},
	"Venus":    {"Mercurio": 1, "Saturno": 1, "Marte": 0, "Júpiter": 0, "Sol": -1, "Luna": -1},
	"Saturno":  {"Mercurio": 1, "Venus": 1, "Júpiter": 0, "Sol": -1, "Luna": -1, "Marte": -1},
}

// El punto de máxima debilidad: el grado exacto opuesto a la exaltación.
func gradoDebil(g string) (float64, bool) {
	e, ok := exalta[g]
	if !ok {
		return 0, false
	}
	return math.Mod(float64(e.signo)*30+e.grado+180, 360), true
}

// vargas que entran en el saptavargaja bala
var sieteVargas = []int{1, 2, 3, 7, 9, 12, 30}

type Bala struct {
	Graha      string  `json:"graha"`
	Yuddha     float64 `json:"yuddha"` // ajuste por guerra; 0 si no pelea
	Rival      string  `json:"rival"`  // con quién, si pelea
	Sthana     float64 `json:"sthana"`
	Dig        float64 `json:"dig"`
	Kala       float64 `json:"kala"`
	Chesta     float64 `json:"chesta"`
	Naisargika float64 `json:"naisargika"`
	Drik       float64 `json:"drik"`
	Total      float64 `json:"total"`  // en virūpas
	Rupas      float64 `json:"rupas"`  // total / 60
	Minimo     float64 `json:"minimo"` // lo que se le exige, en rūpas
	Razon      float64 `json:"razon"`  // rūpas / mínimo. Por debajo de 1 va corto
	Rango      int     `json:"rango"`  // 1 es el más fuerte de la carta
}

type Shadbala struct {
	Balas  []Bala      `json:"balas"`
	Bhavas []BalaBhava `json:"bhavas"`
	Nota   string      `json:"nota"`
}

// EntradaBala es lo que el cálculo necesita saber de cada graha.
type EntradaBala struct {
	Lon   float64 // sidérea
	Bhava int
	Vel   float64 // grados por día
}

// CalcShadbala. lagna es la longitud sidérea del ascendente, mc la del medio
// cielo, y esDeDia dice si el nacimiento fue entre el orto y el ocaso.
func CalcShadbala(gr map[string]EntradaBala, lagna, mc float64,
	esDeDia bool, tithiNum int, senorVara string, senorHora string,
	bhavas []Bhava) Shadbala {

	var s Shadbala
	rasiDe := map[string]int{}
	for g, e := range gr {
		rasiDe[g] = int(e.Lon / 30)
	}

	for _, g := range GrahasAV {
		e, ok := gr[g]
		if !ok {
			continue
		}
		b := Bala{Graha: g, Naisargika: naisargika[g], Minimo: minimoRupas[g]}
		b.Sthana = sthanaBala(g, e, gr, rasiDe)
		b.Dig = digBalaVirupas(g, e.Lon, lagna, mc)
		b.Kala = kalaBala(g, esDeDia, tithiNum, senorVara, senorHora, e.Lon)
		b.Chesta = chestaBala(g, e.Vel)
		b.Drik = drikBala(g, gr)
		b.Total = b.Sthana + b.Dig + b.Kala + b.Chesta + b.Naisargika + b.Drik
		b.Rupas = b.Total / virupasPorRupa
		if b.Minimo > 0 {
			b.Razon = b.Rupas / b.Minimo
		}
		s.Balas = append(s.Balas, b)
	}

	// La guerra se resuelve al final, porque necesita las otras cinco fuerzas
	// ya sumadas para saber quién gana.
	previo := map[string]float64{}
	for _, b := range s.Balas {
		previo[b.Graha] = b.Total
	}
	aj, riv := yuddhaBala(gr, previo)
	for i := range s.Balas {
		g := s.Balas[i].Graha
		if v, ok := aj[g]; ok {
			s.Balas[i].Yuddha, s.Balas[i].Rival = v, riv[g]
			s.Balas[i].Total += v
			s.Balas[i].Rupas = s.Balas[i].Total / virupasPorRupa
			if s.Balas[i].Minimo > 0 {
				s.Balas[i].Razon = s.Balas[i].Rupas / s.Balas[i].Minimo
			}
		}
	}

	// El bhāva bala cuelga del ṣaḍbala del señor de cada casa, así que va
	// después de todo lo anterior.
	balaDe := map[string]float64{}
	for _, b := range s.Balas {
		balaDe[b.Graha] = b.Total
	}
	s.Bhavas = calcBhavaBala(bhavas, balaDe, rasiDe)

	// rango por razón, no por total: cada graha se compara con su propio listón
	for i := range s.Balas {
		r := 1
		for j := range s.Balas {
			if s.Balas[j].Razon > s.Balas[i].Razon {
				r++
			}
		}
		s.Balas[i].Rango = r
	}
	s.Nota = "Las seis fuerzas completas, con yuddha bala. El bhāva bala no mide " +
		"grahas sino casas: cuánto puede el asunto, no cuánto puede el planeta."
	return s
}

// ── 1. Sthāna bala ──

func sthanaBala(g string, e EntradaBala, gr map[string]EntradaBala, rasiDe map[string]int) float64 {
	var total float64

	// a) Ucca bala: lo lejos que está de su punto de máxima debilidad.
	if deb, ok := gradoDebil(g); ok {
		d := math.Abs(e.Lon - deb)
		if d > 180 {
			d = 360 - d
		}
		total += d / 3 // hasta 60
	}

	// b) Saptavargaja bala: la dignidad, contada en siete cartas distintas.
	for _, n := range sieteVargas {
		total += dignidadVirupas(g, Varga(e.Lon, n), gr, rasiDe)
	}

	// c) Ojayugma bala: si el signo par o impar le conviene. En el rāśi y en
	// el navāṁśa. La Luna y Venus quieren pares; los demás, impares.
	quierePar := g == "Luna" || g == "Venus"
	for _, n := range []int{1, 9} {
		par := Varga(e.Lon, n)%2 == 1 // índice 1 = signo 2º = par
		if par == quierePar {
			total += 15
		}
	}

	// d) Kendrādi bala: los pilares valen más que las esquinas.
	switch e.Bhava {
	case 1, 4, 7, 10:
		total += 60
	case 2, 5, 8, 11:
		total += 30
	default:
		total += 15
	}

	// e) Drekkāṇa bala: el tercio del signo. Los masculinos rinden en el
	// primero, los neutros en el segundo, los femeninos en el tercero.
	tercio := int(math.Mod(e.Lon, 30) / 10)
	switch g {
	case "Sol", "Marte", "Júpiter":
		if tercio == 0 {
			total += 15
		}
	case "Mercurio", "Saturno":
		if tercio == 1 {
			total += 15
		}
	case "Luna", "Venus":
		if tercio == 2 {
			total += 15
		}
	}
	return total
}

// dignidadVirupas puntúa en qué casa ajena está el graha, con la amistad
// compuesta: la natural más la de ese momento.
func dignidadVirupas(g string, rasi int, gr map[string]EntradaBala, rasiDe map[string]int) float64 {
	if m, ok := mulatrikona[g]; ok && m.signo == rasi {
		return 45
	}
	duenio := SenorRasi[rasi]
	if duenio == g {
		return 30
	}
	nat := amistadNat[g][duenio]
	// amistad temporal: los que caen en la 2, 3, 4, 10, 11 o 12 desde el graha
	// son amigos de circunstancia; los demás, enemigos de circunstancia.
	tmp := -1
	if rd, ok := rasiDe[duenio]; ok {
		d := ((rd-rasiDe[g])%12 + 12) % 12
		switch d {
		case 1, 2, 3, 9, 10, 11:
			tmp = 1
		}
	}
	switch {
	case nat == 1 && tmp == 1:
		return 22.5 // gran amigo
	case nat == 1 && tmp == -1:
		return 7.5 // neutro
	case nat == 0 && tmp == 1:
		return 15 // amigo
	case nat == 0 && tmp == -1:
		return 3.75 // enemigo
	case nat == -1 && tmp == 1:
		return 7.5 // neutro
	default:
		return 1.875 // gran enemigo
	}
}

// ── 2. Dig bala ──

// Cada graha tiene un punto del cielo donde rinde y el opuesto donde no puede.
// La fuerza es la distancia al punto muerto, dividida entre tres.
func digBalaVirupas(g string, lon, lagna, mc float64) float64 {
	var muerto float64
	switch g {
	case "Júpiter", "Mercurio":
		muerto = math.Mod(lagna+180, 360) // rinden en el Lagna
	case "Sol", "Marte":
		muerto = math.Mod(mc+180, 360) // rinden en el mediocielo
	case "Saturno":
		muerto = lagna // rinde en el descendente
	case "Luna", "Venus":
		muerto = mc // rinden en el fondo del cielo
	}
	d := math.Abs(lon - muerto)
	if d > 180 {
		d = 360 - d
	}
	return d / 3
}

// ── 3. Kāla bala ──

func kalaBala(g string, esDeDia bool, tithiNum int, senorVara, senorHora string, lon float64) float64 {
	var total float64

	// a) Nathonnatha: los de la noche mandan de noche y al revés. Mercurio
	// siempre saca los 60: no tiene preferencia.
	deNoche := map[string]bool{"Luna": true, "Marte": true, "Saturno": true}
	deDia := map[string]bool{"Sol": true, "Júpiter": true, "Venus": true}
	switch {
	case g == "Mercurio":
		total += 60
	case deDia[g] && esDeDia, deNoche[g] && !esDeDia:
		total += 60
	}

	// b) Pakṣa bala: los benéficos crecen con la Luna y los maléficos menguan.
	// La Luna llena da 60 al benéfico y 0 al maléfico.
	// 0 en luna nueva, 60 en llena, y de vuelta a 0.
	pakBenef := float64(tithiNum) / 15 * 60
	if tithiNum > 15 {
		pakBenef = 60 - float64(tithiNum-15)/15*60
	}
	benef := map[string]bool{"Júpiter": true, "Venus": true, "Mercurio": true, "Luna": true}
	if benef[g] {
		total += pakBenef
	} else {
		total += 60 - pakBenef
	}
	if g == "Luna" {
		total += pakBenef // la Luna cobra el pakṣa dos veces
	}

	// c) Tribhāga bala: el día y la noche en tres tercios, uno por graha.
	// Júpiter cobra siempre.
	if g == "Júpiter" {
		total += 60
	}

	// d) Los señores del año, el mes, el día y la hora.
	if senorVara == g {
		total += 45
	}
	if senorHora == g {
		total += 60
	}

	// e) Ayana bala: la declinación. Casi todos quieren declinación norte;
	// la Luna y Saturno, sur. El Sol cobra el doble.
	dec := declinacionDe(lon)
	norte := !(g == "Luna" || g == "Saturno")
	if !norte {
		dec = -dec
	}
	ay := (23.45 + dec) / 46.9 * 60
	if g == "Mercurio" {
		ay = 60 // Mercurio rinde en las dos
	}
	if g == "Sol" {
		ay *= 2
	}
	total += ay
	return total
}

// declinacionDe da la declinación de una longitud eclíptica, tomando la
// oblicuidad como constante. Basta para el ayana bala.
func declinacionDe(lon float64) float64 {
	const eps = 23.4392911 * math.Pi / 180
	l := lon * math.Pi / 180
	return math.Asin(math.Sin(eps)*math.Sin(l)) * 180 / math.Pi
}

// ── 4. Cheṣṭā bala ──

// velocidad media diaria de cada graha, en grados
var velMedia = map[string]float64{
	"Sol": 0.9856, "Luna": 13.176, "Marte": 0.524, "Mercurio": 1.383,
	"Júpiter": 0.083, "Venus": 1.602, "Saturno": 0.033,
}

// Los ocho estados del movimiento. El Sol y la Luna no retrogradan nunca, así
// que se les mide solo por lo rápido que van.
func chestaBala(g string, vel float64) float64 {
	m := velMedia[g]
	if m == 0 {
		return 30
	}
	switch {
	case vel < 0 && math.Abs(vel) > m*0.5:
		return 60 // vakra: retrógrado franco
	case vel < 0:
		return 30 // anuvakra: retrocediendo poco
	case math.Abs(vel) < m*0.05:
		return 15 // vikala: casi parado
	case vel < m*0.5:
		return 7.5 // mandatara: muy lento
	case vel < m*0.9:
		return 15 // manda: lento
	case vel > m*1.5:
		return 45 // aticāra: disparado
	case vel > m*1.1:
		return 30 // cāra: rápido
	default:
		return 30 // sama: a su paso
	}
}

// ── 6. Dṛk bala ──

// La mirada de un benéfico suma y la de un maléfico resta. Se cuenta por
// signo entero, que es como aspecta jyotiṣa.
func drikBala(g string, gr map[string]EntradaBala) float64 {
	miRasi := int(gr[g].Lon / 30)
	var total float64
	for otro, e := range gr {
		if otro == g {
			continue
		}
		r := int(e.Lon / 30)
		mira := false
		for _, d := range Drishti(otro) {
			if (r+d-1)%12 == miRasi {
				mira = true
			}
		}
		if !mira {
			continue
		}
		switch benefico(otro) {
		case 1:
			total += 15
		case -1:
			total -= 15
		default:
			total += 7.5
		}
	}
	return total
}

// ── Yuddha bala: la guerra entre grahas ──
//
// Cuando dos grahas —de Marte a Saturno, nunca las luminarias ni los nodos—
// quedan a menos de un grado, la tradición dice que están en guerra. Gana el
// que esté más al norte en latitud; aquí se toma un criterio más simple y
// declarado: gana el que tenga más ṣaḍbala acumulado.
//
// La diferencia de fuerza dividida por la diferencia de diámetros aparentes se
// suma al vencedor y se resta al vencido. Los diámetros son los clásicos, en
// minutos de arco.
var diametro = map[string]float64{
	"Marte": 9.4, "Mercurio": 6.6, "Júpiter": 190.4, "Venus": 16.6, "Saturno": 158.0,
}

// enGuerra son los únicos que pueden pelear: los cinco planetas de verdad.
var enGuerra = []string{"Marte", "Mercurio", "Júpiter", "Venus", "Saturno"}

// yuddhaBala devuelve el ajuste que le toca a cada graha por las guerras en
// las que anda metido. Devuelve también con quién pelea, para poder decirlo.
func yuddhaBala(gr map[string]EntradaBala, previo map[string]float64) (map[string]float64, map[string]string) {
	ajuste, rival := map[string]float64{}, map[string]string{}
	for i, a := range enGuerra {
		for _, b := range enGuerra[i+1:] {
			ea, oka := gr[a]
			eb, okb := gr[b]
			if !oka || !okb {
				continue
			}
			d := math.Abs(ea.Lon - eb.Lon)
			if d > 180 {
				d = 360 - d
			}
			if d >= 1 {
				continue
			}
			gana, pierde := a, b
			if previo[b] > previo[a] {
				gana, pierde = b, a
			}
			// diferencia de fuerza repartida según los diámetros aparentes
			dif := math.Abs(previo[gana]-previo[pierde]) /
				math.Abs(diametro[gana]-diametro[pierde]+0.001)
			if dif > 60 {
				dif = 60 // la tradición no da un tope; se pone uno para no disparar
			}
			ajuste[gana] += dif
			ajuste[pierde] -= dif
			rival[gana], rival[pierde] = pierde, gana
		}
	}
	return ajuste, rival
}

// ── Bhāva bala: la fuerza de las casas, no de los grahas ──
//
// Es un cálculo aparte del ṣaḍbala y contesta otra pregunta: no cuánto puede un
// graha, sino cuánto puede un asunto. Tiene tres partes.
//
//   Bhāvādhipati — la fuerza del señor de la casa. Es la mayor con diferencia,
//                  y por eso el ṣaḍbala tiene que estar hecho antes.
//   Bhāva dig    — la casa rinde según qué clase de graha la ocupa: las de
//                  bienes quieren benéficos, las de esfuerzo quieren maléficos.
//   Bhāva dṛṣṭi  — quién la mira, con signo.

type BalaBhava struct {
	Numero   int     `json:"numero"`
	Senor    float64 `json:"senor"`  // ṣaḍbala de su señor, en virūpas
	Dig      float64 `json:"dig"`
	Drishti  float64 `json:"drishti"`
	Total    float64 `json:"total"`
	Rupas    float64 `json:"rupas"`
	Rango    int     `json:"rango"`
}

// Qué clase de graha le conviene a cada bhāva. La 1, 4, 7 y 10 son de bienes y
// prefieren benéficos; la 3, 6 y 11 son de esfuerzo y prefieren maléficos.
var quiereBenefico = [12]bool{true, true, false, true, true, false,
	true, true, true, true, false, true}

func calcBhavaBala(bhavas []Bhava, balaDe map[string]float64, rasiDe map[string]int) []BalaBhava {
	var out []BalaBhava
	for _, b := range bhavas {
		if b.Numero < 1 || b.Numero > 12 {
			continue
		}
		bb := BalaBhava{Numero: b.Numero, Senor: balaDe[b.Senor]}

		// dig: suma quien encaja con lo que la casa pide, resta quien no
		for _, g := range b.Ocupan {
			s := benefico(g)
			if s == 0 {
				continue
			}
			if (s == 1) == quiereBenefico[b.Numero-1] {
				bb.Dig += 30
			} else {
				bb.Dig -= 30
			}
		}
		// dṛṣṭi: la mirada de un benéfico suma, la de un maléfico resta
		for _, g := range b.Aspectan {
			switch benefico(g) {
			case 1:
				bb.Drishti += 15
			case -1:
				bb.Drishti -= 15
			default:
				bb.Drishti += 7.5
			}
		}
		bb.Total = bb.Senor + bb.Dig + bb.Drishti
		bb.Rupas = bb.Total / virupasPorRupa
		out = append(out, bb)
	}
	for i := range out {
		r := 1
		for j := range out {
			if out[j].Total > out[i].Total {
				r++
			}
		}
		out[i].Rango = r
	}
	return out
}
