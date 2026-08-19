package efem

import "fmt"

// Sistema de Palabras Clave (Margaret Hone). La app COMPONE la frase a partir de
// los componentes; no guarda interpretaciones prefabricadas. Es el paso 1 de los
// cuatro que enseña el módulo 10: traducción literal y deliberadamente torpe.

// Sujetos siempre en singular para que la frase concuerde al componerla.
var FuncionPlaneta = map[string]string{
	"Sol": "tu vitalidad", "Luna": "tu manera de reaccionar",
	"Mercurio": "tu forma de pensar y de hablar", "Venus": "tu manera de valorar y de unirte",
	"Marte": "tu empuje", "Júpiter": "tu manera de expandirte y confiar",
	"Saturno": "tu manera de poner límites", "Urano": "tu manera de romper y despertar",
	"Neptuno": "tu manera de disolver e imaginar", "Plutón": "tu manera de vaciar y regenerar",
	"Nodo Norte": "hacia dónde tiras", "Nodo Sur": "de dónde vienes",
}

var ModoSigno = [12]string{"con arranque y sin rodeos", "con paciencia y apego a lo tangible",
	"con curiosidad y muchas versiones", "con sensibilidad y ganas de proteger",
	"con brillo y ganas de crear", "con detalle y sentido crítico",
	"buscando acuerdo y equilibrio", "con intensidad y sin soltar",
	"con amplitud y ganas de sentido", "con prudencia y mirada larga",
	"con distancia y cabeza fría", "con porosidad y sin bordes claros"}

var TerrenoCasa = [12]string{"tu cuerpo y tu manera de asomarte al mundo",
	"lo que posees y lo que valoras", "tu entorno cercano y cómo te comunicas",
	"tu casa y tus raíces", "lo que creas, disfrutas y amas",
	"tu trabajo diario y tu salud", "tus vínculos íntimos y quien te hace de espejo",
	"lo que compartes y las crisis que transforman", "lo que estudias y aquello en lo que crees",
	"tu vocación y lo que se ve de ti desde fuera", "tus amistades y tus proyectos",
	"tu retiro, lo que sueltas y lo que no se ve"}

var CualidadAspecto = map[string]string{
	"conjunción": "se funde con", "oposición": "tira en sentido contrario a",
	"trígono": "fluye con facilidad hacia", "cuadratura": "roza y cuesta con",
	"sextil": "colabora, si le haces caso, con", "semicuadratura": "roza levemente con",
	"sesquicuadratura": "roza levemente con", "quincuncio": "no acaba de encajar con",
	"semisextil": "se relaciona con cierta incomodidad con",
}

var MatizDignidad = map[string]string{
	"domicilio": "y aquí funciona a pleno rendimiento",
	"exaltación": "y aquí va sobrado, quizá de más",
	"exilio": "y aquí trabaja a contrapelo, lo que con los años da oficio",
	"caída": "y aquí se siente incómodo y poco valorado",
}

// Categorías de recogida (Hone): las frases se agrupan antes de sintetizar.
var Categorias = []string{"Carácter", "Mente", "Trabajo y dinero", "Vínculos",
	"Familia y raíces", "Salud", "Retiro e interior"}

var casaCategoria = [12]string{"Carácter", "Trabajo y dinero", "Mente", "Familia y raíces",
	"Vínculos", "Salud", "Vínculos", "Vínculos", "Mente", "Trabajo y dinero",
	"Vínculos", "Retiro e interior"}

type Frase struct {
	Categoria string `json:"categoria"`
	Texto     string `json:"texto"`
	Fuente    string `json:"fuente"`
	Peso      float64 `json:"peso"`
}

type Lectura struct {
	Frases          []Frase  `json:"frases"`
	Contradicciones []string `json:"contradicciones"`
	Dominante       string   `json:"dominante"`
	Nota            string   `json:"nota"`
}

// duro indica si un aspecto es de tensión.
func duro(n string) bool {
	return n == "cuadratura" || n == "oposición" || n == "semicuadratura" || n == "sesquicuadratura"
}
func blando(n string) bool { return n == "trígono" || n == "sextil" }

// Interpretar compone la traducción literal de la carta y agrupa por categorías.
// NO sintetiza: eso es trabajo del que lee, y el módulo 10 explica por qué.
func (c Carta) Interpretar() Lectura {
	var L Lectura
	for _, p := range c.Cuerpos {
		if p.Nombre == "Nodo Norte" || p.Nombre == "Nodo Sur" {
			continue
		}
		t := fmt.Sprintf("%s se expresa %s, en el terreno de %s",
			FuncionPlaneta[p.Nombre], ModoSigno[p.SignoIdx], TerrenoCasa[p.CasaP-1])
		if m, ok := MatizDignidad[p.Dignidad]; ok {
			t += " — " + m
		}
		t += "."
		L.Frases = append(L.Frases, Frase{Categoria: casaCategoria[p.CasaP-1], Texto: t,
			Fuente: fmt.Sprintf("%s en %s, casa %d", p.Nombre, p.Signo, p.CasaP), Peso: 1})
	}
	// aspectos: solo los que pesan de verdad
	for _, a := range c.Aspectos {
		if a.Orbe > 4 {
			continue
		}
		t := fmt.Sprintf("%s %s %s.", FuncionPlaneta[a.A], CualidadAspecto[a.Nombre], FuncionPlaneta[a.B])
		L.Frases = append(L.Frases, Frase{Categoria: "Carácter", Texto: t,
			Fuente: fmt.Sprintf("%s %s %s, orbe %.2f°", a.A, a.Nombre, a.B, a.Orbe),
			Peso: 1 + (4-a.Orbe)/2})
	}
	// regentes: lo que convierte casas sueltas en argumentos
	for i, r := range c.Regentes {
		if c.RegenteEn[i] == i+1 {
			continue
		}
		L.Frases = append(L.Frases, Frase{Categoria: casaCategoria[i],
			Texto: fmt.Sprintf("Lo relativo a %s pasa por %s: el regente de la casa %d está alojado en la %d.",
				TerrenoCasa[i], TerrenoCasa[c.RegenteEn[i]-1], i+1, c.RegenteEn[i]),
			Fuente: fmt.Sprintf("regente de %d (%s) en %d", i+1, r, c.RegenteEn[i]), Peso: 1.2})
	}
	// contradicciones: mismo planeta con aspectos duros y blandos a la vez
	porPlaneta := map[string][]Aspecto{}
	for _, a := range c.Aspectos {
		if a.Orbe <= 6 {
			porPlaneta[a.A] = append(porPlaneta[a.A], a)
			porPlaneta[a.B] = append(porPlaneta[a.B], a)
		}
	}
	for _, n := range Orden[:10] {
		var d, b *Aspecto
		for i := range porPlaneta[n] {
			a := porPlaneta[n][i]
			if duro(a.Nombre) && (d == nil || a.Orbe < d.Orbe) {
				d = &porPlaneta[n][i]
			}
			if blando(a.Nombre) && (b == nil || a.Orbe < b.Orbe) {
				b = &porPlaneta[n][i]
			}
		}
		if d != nil && b != nil {
			otroD, otroB := d.A, b.A
			if otroD == n {
				otroD = d.B
			}
			if otroB == n {
				otroB = b.B
			}
			manda := "la tensión"
			if b.Orbe < d.Orbe {
				manda = "la facilidad"
			}
			L.Contradicciones = append(L.Contradicciones, fmt.Sprintf(
				"%s recibe a la vez %s de %s (orbe %.2f°) y %s de %s (orbe %.2f°). "+
					"Por exactitud manda %s, pero las dos cosas están ahí: descríbelas juntas con «y», no con «pero».",
				n, d.Nombre, otroD, d.Orbe, b.Nombre, otroB, b.Orbe, manda))
		}
	}
	L.Dominante = c.dominante()
	L.Nota = "Esto es el paso 1 de 4: traducción literal, deliberadamente torpe. " +
		"Agrupar por categorías es el paso 2 y ya está hecho. Los pasos 3 y 4 —resolver las " +
		"contradicciones y escribirlo como texto seguido— los haces tú. La carta no se cuadra sola."
	return L
}

// dominante estima qué planeta manda: regencia del Ascendente, dignidad,
// angularidad y aspectos exactos.
func (c Carta) dominante() string {
	p := map[string]float64{}
	regAsc := RegenteSigno[int(c.Asc/30)]
	p[regAsc] += 3
	for _, b := range c.Cuerpos {
		switch b.Dignidad {
		case "domicilio":
			p[b.Nombre] += 2
		case "exaltación":
			p[b.Nombre] += 1.5
		}
		switch b.CasaP {
		case 1, 4, 7, 10:
			p[b.Nombre] += 1.5
		case 2, 5, 8, 11:
			p[b.Nombre] += 0.5
		}
	}
	for _, a := range c.Aspectos {
		if a.Orbe <= 2 {
			p[a.A] += 0.7
			p[a.B] += 0.7
		}
	}
	mejor, max := regAsc, -1.0
	for _, n := range Orden[:10] {
		if p[n] > max {
			mejor, max = n, p[n]
		}
	}
	return fmt.Sprintf("%s (regente del Ascendente: %s)", mejor, regAsc)
}
