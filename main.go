// astro — cartas natales y curso, en las dos tradiciones y en dos idiomas.
// Un binario único: sin Python, sin red, sin instalar nada.
package main

import (
	"astro/internal/efem"
	"astro/internal/guardadas"
	"astro/internal/jyotisha"
	"astro/internal/lugares"
	"astro/internal/occidental"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"math"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

//go:embed web
var contenido embed.FS

// Version es el único sitio donde se escribe. De aquí sale el pie de la
// página, la respuesta de /api/version y la opción -version; empaquetar/todo.sh
// comprueba que coincide con la etiqueta que se va a publicar, para que no se
// publique una versión que dentro dice otra cosa.
const Version = "1.8.0"

func main() {
	puerto := flag.Int("puerto", 8733, "puerto del servidor")
	abrir := flag.Bool("abrir", true, "abrir el navegador al arrancar")
	red := flag.Bool("red", false, "aceptar conexiones de otros equipos de la red local")
	verVersion := flag.Bool("version", false, "decir la versión y salir")
	flag.Parse()
	if *verVersion {
		fmt.Println("astro " + Version)
		return
	}

	sub, _ := fs.Sub(contenido, "web")
	http.Handle("/", http.FileServer(http.FS(sub)))
	http.HandleFunc("/api/carta", apiCarta)             // occidental
	http.HandleFunc("/api/vedica", apiVedica)           // jyotiṣa
	http.HandleFunc("/api/comparar", apiComparar)       // las dos, lado a lado
	http.HandleFunc("/api/lectura", apiLectura)         // occidental
	http.HandleFunc("/api/lecturaved", apiLecturaVed)   // jyotiṣa
	http.HandleFunc("/api/prediccion", apiPrediccion)   // occidental: tiempo
	http.HandleFunc("/api/verificar", apiVerificar)
	http.HandleFunc("/api/verificarved", apiVerificarVed)
	http.HandleFunc("/api/prasna", apiPrasna)          // la carta de la pregunta
	http.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		jsonOut(w, map[string]string{"version": Version})
	})
	http.HandleFunc("/api/lugares", apiLugares)
	http.HandleFunc("/api/huso", apiHuso)
	http.HandleFunc("/api/husohistoria", apiHistoria)
	http.HandleFunc("/api/guardadas", apiGuardadas)

	var ln net.Listener
	var err error
	usado := *puerto
	for i := 0; i < 40; i++ {
		host := "127.0.0.1"
		if *red {
			host = "0.0.0.0"
		}
		if ln, err = net.Listen("tcp", fmt.Sprintf("%s:%d", host, usado)); err == nil {
			break
		}
		usado++
	}
	if ln == nil {
		fmt.Printf("\n  No he podido abrir ningún puerto entre %d y %d.\n"+
			"  Prueba con:  %s -puerto=9500\n\n", *puerto, usado, os.Args[0])
		os.Exit(1)
	}
	dir := fmt.Sprintf("http://localhost:%d", usado)
	aviso := ""
	if usado != *puerto {
		aviso = fmt.Sprintf("  (el puerto %d estaba ocupado)\n", *puerto)
	}
	fmt.Printf("\n  Astro — occidental y jyotiṣa\n  %s\n%s", dir, aviso)
	if *red {
		fmt.Println("\n  Abierto a la red local:")
		for _, ip := range ipsLocales() {
			fmt.Printf("    http://%s:%d\n", ip, usado)
		}
		fmt.Println("\n  AVISO: cualquiera en esta red puede ver o borrar tus cartas.")
	}
	fmt.Print("\n  Ctrl-C para parar.\n\n")
	if *abrir {
		go func() { time.Sleep(400 * time.Millisecond); abrirNavegador(dir) }()
	}
	if err := http.Serve(ln, nil); err != nil {
		fmt.Println("  error:", err)
	}
}

func ipsLocales() []string {
	var out []string
	dirs, _ := net.InterfaceAddrs()
	for _, a := range dirs {
		if n, ok := a.(*net.IPNet); ok && !n.IP.IsLoopback() {
			if v4 := n.IP.To4(); v4 != nil {
				out = append(out, v4.String())
			}
		}
	}
	if len(out) == 0 {
		out = append(out, "la-ip-de-esta-maquina")
	}
	return out
}

func abrirNavegador(url string) {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("open", url).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
}

func num(r *http.Request, k string) float64 {
	v, _ := strconv.ParseFloat(r.URL.Query().Get(k), 64)
	return v
}
func ent(r *http.Request, k string) int { v, _ := strconv.Atoi(r.URL.Query().Get(k)); return v }
func datos(r *http.Request) (int, int, int, int, int, float64, float64, float64) {
	return ent(r, "anio"), ent(r, "mes"), ent(r, "dia"), ent(r, "hh"), ent(r, "mm"),
		num(r, "tz"), num(r, "lat"), num(r, "lon")
}
func jsonOut(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

func apiCarta(w http.ResponseWriter, r *http.Request) {
	a, m, d, h, mi, tz, lat, lo := datos(r)
	jsonOut(w, efem.Calcular(a, m, d, h, mi, tz, lat, lo))
}

func apiVedica(w http.ResponseWriter, r *http.Request) {
	a, m, d, h, mi, tz, lat, lo := datos(r)
	jsonOut(w, jyotisha.CalcularOpts(a, m, d, h, mi, tz, lat, lo,
		jyotisha.Opciones{NodoVerdadero: nodoVerdadero(r), Lang: idioma(r)}))
}

// nodoVerdadero lee ?nodo=verdadero. Por defecto se usa el medio.
func nodoVerdadero(r *http.Request) bool {
	return r.URL.Query().Get("nodo") == "verdadero"
}

func apiLectura(w http.ResponseWriter, r *http.Request) {
	a, m, d, h, mi, tz, lat, lo := datos(r)
	jsonOut(w, occidental.Interpretar(efem.Calcular(a, m, d, h, mi, tz, lat, lo), idioma(r)))
}

// apiLecturaVed es la lectura védica. Razona por cadena de señores, que es lo
// que en jyotiṣa convierte una posición suelta en una causa.
func apiLecturaVed(w http.ResponseWriter, r *http.Request) {
	a, m, d, h, mi, tz, lat, lo := datos(r)
	jsonOut(w, jyotisha.Interpretar(jyotisha.Calcular(a, m, d, h, mi, tz, lat, lo), idioma(r)))
}

// apiPrediccion da tránsitos, progresiones y revolución solar para una fecha.
// Sin ?cuando= se toma hoy, que es lo que se quiere el noventa por ciento de
// las veces.
func apiPrediccion(w http.ResponseWriter, r *http.Request) {
	a, m, d, h, mi, tz, lat, lo := datos(r)
	natal := efem.Calcular(a, m, d, h, mi, tz, lat, lo)
	cuando := time.Now()
	if s := r.URL.Query().Get("cuando"); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			cuando = t
		}
	}
	jsonOut(w, occidental.Predecir(natal, cuando, idioma(r)))
}

// idioma lee ?lang= y cae en español si no viene o no se reconoce.
func idioma(r *http.Request) string {
	if r.URL.Query().Get("lang") == "en" {
		return "en"
	}
	return "es"
}

// apiComparar es la única pantalla donde las dos tradiciones aparecen juntas.
// Se muestra a propósito: enseña por qué el mismo nacimiento da signos distintos.
func apiComparar(w http.ResponseWriter, r *http.Request) {
	a, m, d, h, mi, tz, lat, lo := datos(r)
	nodoV := nodoVerdadero(r)
	occ := efem.CalcularCon(a, m, d, h, mi, tz, lat, lo, nodoV)
	ved := jyotisha.CalcularCon(a, m, d, h, mi, tz, lat, lo, nodoV)

	// Se mandan indices y grados, no cadenas ya compuestas: el nombre del signo
	// depende del idioma y lo pone el navegador.
	type fila struct {
		Cuerpo   string  `json:"cuerpo"`
		Glifo    string  `json:"glifo"`
		TropIdx  int     `json:"tropIdx"`
		TropGr   float64 `json:"tropGr"`
		SidIdx   int     `json:"sidIdx"`
		SidGr    float64 `json:"sidGr"`
		Cambia   bool    `json:"cambia"`
		CasaOcc  int     `json:"casaOcc"`  // Plácido
		CasaVed  int     `json:"casaVed"`  // signo entero
		CambiaC  bool    `json:"cambiaC"`  // la casa no es la misma
	}
	var filas []fila
	ved2 := map[string]jyotisha.Graha{}
	for _, g := range ved.Grahas {
		ved2[g.Nombre] = g
	}
	equiv := map[string]string{"Nodo Norte": "Rāhu", "Nodo Sur": "Ketu"}
	for _, c := range occ.Cuerpos {
		n := c.Nombre
		if e, ok := equiv[n]; ok {
			n = e
		}
		v, ok := ved2[n]
		if !ok {
			continue
		}
		filas = append(filas, fila{Cuerpo: n, Glifo: v.Glifo,
			TropIdx: c.SignoIdx, TropGr: c.Grado,
			SidIdx: v.RasiIdx, SidGr: v.Grado,
			Cambia:  c.Signo != jyotisha.RasisEs[v.RasiIdx],
			CasaOcc: c.CasaP, CasaVed: v.Bhava,
			CambiaC: c.CasaP != v.Bhava})
	}
	jsonOut(w, map[string]any{
		"ayanamsa":     ved.Ayanamsa,
		"ascIdx":       int(occ.Asc / 30),
		"ascGr":        math.Mod(occ.Asc, 30),
		"lagIdx":       ved.LagnaRasi,
		"lagGr":        math.Mod(ved.Lagna, 30),
		"cambiaLagna":  efem.Signos[int(occ.Asc/30)] != jyotisha.RasisEs[ved.LagnaRasi],
		"filas":        filas,
	})
}

func apiVerificar(w http.ResponseWriter, r *http.Request) {
	a, m, d, h, mi, tz, lat, lo := datos(r)
	c := efem.Calcular(a, m, d, h, mi, tz, lat, lo)
	comp := func(nombre string, tuyo, correcto, tol float64, unidad string, ang bool) map[string]any {
		dd := tuyo - correcto
		if ang {
			dd = math.Mod(dd+540, 360) - 180
		}
		bien := math.Abs(dd) <= tol
		com := "correcto"
		switch {
		case bien:
		case nombre == "Tiempo sidéreo local" && math.Abs(math.Abs(dd)-math.Abs(lo/15)) < 0.05:
			com = "te has dejado la corrección por longitud, o la has sumado al revés"
		case nombre == "Ascendente" && math.Abs(dd) > 100:
			com = "estás medio zodíaco desviado: revisa el cuadrante del arcotangente"
		default:
			com = "revisa este paso antes de seguir"
		}
		return map[string]any{"nombre": nombre, "tuyo": tuyo, "correcto": correcto,
			"desvio": math.Round(math.Abs(dd)*10000) / 10000, "unidad": unidad,
			"bien": bien, "comentario": com}
	}
	pasos := []map[string]any{
		comp("Día juliano", num(r, "jd"), c.JD, 0.0007, "días", false),
		comp("T. sidéreo Greenwich", num(r, "tsg"), c.TSG, 0.05, "grados", true),
		comp("Tiempo sidéreo local", num(r, "tsl"), c.TSL, 0.05, "grados", true),
		comp("Ascendente", num(r, "asc"), c.Asc, 0.05, "grados", true),
		comp("Medio Cielo", num(r, "mc"), c.MC, 0.05, "grados", true),
	}
	primero := -1
	for i, p := range pasos {
		if !p["bien"].(bool) {
			primero = i
			break
		}
	}
	jsonOut(w, map[string]any{"pasos": pasos, "primerFallo": primero})
}

// apiVerificarVed corrige el cálculo a mano de una carta védica. Comparte los
// tres primeros pasos con el occidental —la astronomía es la misma— y añade
// los dos que son propios: restar el ayanāṁśa y sacar el rāśi del Lagna, que
// es donde se equivoca todo el mundo al empezar.
func apiVerificarVed(w http.ResponseWriter, r *http.Request) {
	a, m, d, h, mi, tz, lat, lo := datos(r)
	c := jyotisha.Calcular(a, m, d, h, mi, tz, lat, lo)
	base := efem.Calcular(a, m, d, h, mi, tz, lat, lo)
	EN := idioma(r) == "en"

	tx := func(es, en string) string {
		if EN {
			return en
		}
		return es
	}
	comp := func(nombre string, tuyo, correcto, tol float64, unidad string, ang bool, pista string) map[string]any {
		dd := tuyo - correcto
		if ang {
			dd = math.Mod(dd+540, 360) - 180
		}
		bien := math.Abs(dd) <= tol
		com := tx("correcto", "correct")
		if !bien {
			com = pista
		}
		return map[string]any{"nombre": nombre, "tuyo": tuyo, "correcto": correcto,
			"desvio": math.Round(math.Abs(dd)*10000) / 10000, "unidad": unidad,
			"bien": bien, "comentario": com}
	}
	grados := tx("grados", "degrees")
	pasos := []map[string]any{
		comp(tx("Día juliano", "Julian day"), num(r, "jd"), c.JD, 0.0007,
			tx("días", "days"), false,
			tx("revisa la conversión a hora universal antes de seguir",
				"check the conversion to universal time before going on")),
		comp(tx("Tiempo sidéreo local", "Local sidereal time"), num(r, "tsl"), base.TSL, 0.05, grados, true,
			tx("te has dejado la corrección por longitud, o la has sumado al revés",
				"you have left out the longitude correction, or added it the wrong way")),
		comp(tx("Ascendente tropical", "Tropical ascendant"), num(r, "asctrop"), base.Asc, 0.05, grados, true,
			tx("este paso es el mismo que en occidental: revísalo antes de restar nada",
				"this step is the same as in western astrology: check it before subtracting anything")),
		comp(tx("Ayanāṁśa", "Ayanāṁśa"), num(r, "ayan"), c.Ayanamsa, 0.02, grados, false,
			tx("el ayanāṁśa cambia con los años: no vale restar 24° fijos, hay que calcularlo para la fecha",
				"the ayanāṁśa changes over the years: subtracting a flat 24° will not do, it has to be calculated for the date")),
		comp(tx("Lagna sidéreo", "Sidereal lagna"), num(r, "lagna"), c.Lagna, 0.05, grados, true,
			tx("el ayanāṁśa se RESTA una sola vez al ascendente tropical; restarlo dos veces es el error más común al empezar",
				"the ayanāṁśa is SUBTRACTED once from the tropical ascendant; subtracting it twice is the commonest beginner's mistake")),
	}
	primero := -1
	for i, p := range pasos {
		if !p["bien"].(bool) {
			primero = i
			break
		}
	}
	// el rāśi que sale, para que pueda comprobarlo de un vistazo
	jsonOut(w, map[string]any{"pasos": pasos, "primerFallo": primero,
		"lagnaRasi": c.LagnaRasi, "lagnaPos": c.LagnaPos})
}

// apiPrasna levanta la carta del momento en que se pregunta. Sin ?cuando= se
// toma el instante de ahora, que es lo propio de un praśna: la pregunta se hace
// cuando se hace.
func apiPrasna(w http.ResponseWriter, r *http.Request) {
	lat, lon := num(r, "lat"), num(r, "lon")
	t := time.Now()
	if s := r.URL.Query().Get("cuando"); s != "" {
		if p, err := time.Parse("2006-01-02T15:04", s); err == nil {
			t = p
		}
	}
	// El praśna se levanta en hora local del sitio desde el que se pregunta;
	// aquí se pasa la hora ya en UT y el huso a cero.
	u := t.UTC()
	c := jyotisha.Calcular(u.Year(), int(u.Month()), u.Day(), u.Hour(), u.Minute(), 0, lat, lon)
	tema := r.URL.Query().Get("tema")
	if tema == "" {
		tema = "salud"
	}
	jsonOut(w, map[string]any{
		"prasna": jyotisha.HacerPrasna(c, tema, idioma(r)),
		"carta":  c,
		"cuando": u.Format("2006-01-02 15:04 UT"),
	})
}

func apiLugares(w http.ResponseWriter, r *http.Request) {
	jsonOut(w, lugares.Buscar(r.URL.Query().Get("q"), 8))
}

func apiHuso(w http.ResponseWriter, r *http.Request) {
	off, nom, verano, err := lugares.Huso(r.URL.Query().Get("zona"),
		ent(r, "anio"), ent(r, "mes"), ent(r, "dia"), ent(r, "hh"), ent(r, "mm"))
	if err != nil {
		jsonOut(w, map[string]any{"error": err.Error()})
		return
	}
	jsonOut(w, map[string]any{"offset": off, "zona": nom, "verano": verano})
}

func apiHistoria(w http.ResponseWriter, r *http.Request) {
	h, err := lugares.HistoriaHuso(r.URL.Query().Get("zona"), ent(r, "anio"), ent(r, "mes"),
		ent(r, "dia"), ent(r, "hh"), ent(r, "mm"), num(r, "lon"))
	if err != nil {
		jsonOut(w, map[string]any{"error": err.Error()})
		return
	}
	jsonOut(w, h)
}

func apiGuardadas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var c guardadas.Carta
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil || c.Nombre == "" {
			http.Error(w, "datos inválidos", http.StatusBadRequest)
			return
		}
		guardadas.Guardar(c)
	case http.MethodDelete:
		guardadas.Borrar(r.URL.Query().Get("id"))
	}
	jsonOut(w, map[string]any{"cartas": guardadas.Listar(), "fichero": guardadas.Ruta()})
}
