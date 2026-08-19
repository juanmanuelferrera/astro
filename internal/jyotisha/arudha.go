package jyotisha

import "math"

// Arudha padas y lagnas especiales, los dos de Jaimini.
//
// El arudha es la IMAGEN de un bhāva: no lo que la casa es, sino lo que de
// ella se ve desde fuera. El Arudha Lagna (AL) es la reputación frente al
// Lagna, que es la persona. Que los dos coincidan o se opongan dice mucho.
//
// La regla: se cuenta del bhāva a su señor, y desde el señor se vuelve a
// contar lo mismo. Con dos excepciones — si el resultado cae en el propio
// bhāva o en el séptimo desde él, la imagen no se sostiene ahí y se toma el
// décimo desde ese punto.

// ArudhaDe devuelve el rāśi (0-11) del arudha de un bhāva.
// rasiBhava es el rāśi del bhāva y rasiSenor el rāśi donde está su señor.
func ArudhaDe(rasiBhava, rasiSenor int) int {
	n := ((rasiSenor-rasiBhava)%12 + 12) % 12 // distancia del bhāva a su señor
	a := (rasiSenor + n) % 12
	// Excepciones: el arudha no puede quedarse en su propia casa ni enfrente.
	d := ((a-rasiBhava)%12 + 12) % 12
	if d == 0 || d == 6 {
		a = (a + 9) % 12 // el décimo desde ahí, contando inclusive
	}
	return a
}

type Arudhas struct {
	Padas [12]int `json:"padas"` // rāśi del arudha de cada bhāva, A1(=AL) a A12
	AL    int     `json:"al"`    // Arudha Lagna
	UL    int     `json:"ul"`    // Upapada, el arudha de la 12: la pareja
}

// CalcArudhas necesita el rāśi de cada bhāva y dónde está su señor.
func CalcArudhas(bhavas []Bhava, rasiDe map[string]int) Arudhas {
	var a Arudhas
	for i, b := range bhavas {
		if i > 11 {
			break
		}
		rs, ok := rasiDe[b.Senor]
		if !ok {
			a.Padas[i] = b.RasiIdx
			continue
		}
		a.Padas[i] = ArudhaDe(b.RasiIdx, rs)
	}
	a.AL, a.UL = a.Padas[0], a.Padas[11]
	return a
}

// ── Lagnas especiales ──
//
// Los tres se cuentan como tiempo transcurrido desde el amanecer, a partir de
// la posición del Sol en ese instante. Cambian de rāśi a velocidades
// distintas, y por eso miden cosas distintas: el Ghaṭī Lagna corre setenta y
// cinco veces más rápido que el zodíaco y sirve para el poder y el mando; el
// Horā Lagna, para la riqueza; el Bhāva Lagna, para el cuerpo y la vida.

type LagnasEsp struct {
	Bhava   float64 `json:"bhava"`   // Bhāva Lagna: un rāśi cada 2 h
	Hora    float64 `json:"hora"`    // Horā Lagna: un rāśi cada hora
	Ghati   float64 `json:"ghati"`   // Ghaṭī Lagna: un rāśi cada 24 min
	Amanece float64 `json:"amanece"` // hora UT del orto, para poder comprobarlo
	Hay     bool    `json:"hay"`     // falso en latitudes donde no sale el Sol
}

// CalcLagnasEsp toma el Sol sidéreo del amanecer y las horas transcurridas
// desde entonces hasta el nacimiento.
func CalcLagnasEsp(solAlAmanecer, horasDesdeOrto, orto float64, hay bool) LagnasEsp {
	if !hay {
		return LagnasEsp{Hay: false}
	}
	n := func(x float64) float64 { return math.Mod(math.Mod(x, 360)+360, 360) }
	return LagnasEsp{
		Bhava:   n(solAlAmanecer + horasDesdeOrto*15), // 30° / 2 h
		Hora:    n(solAlAmanecer + horasDesdeOrto*30), // 30° / 1 h
		Ghati:   n(solAlAmanecer + horasDesdeOrto*75), // 30° / 24 min
		Amanece: orto,
		Hay:     true,
	}
}
