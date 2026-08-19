// astro — curso de astrología con calculadora de cartas natales.
// Binario único: no necesita Python, ni red, ni instalar nada.
package main

import (
	"astro/internal/efem"
	"astro/internal/guardadas"
	"astro/internal/lugares"
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

type paso struct {
	Nombre    string  `json:"nombre"`
	Tuyo      float64 `json:"tuyo"`
	Correcto  float64 `json:"correcto"`
	Desvio    float64 `json:"desvio"`
	Unidad    string  `json:"unidad"`
	Bien      bool    `json:"bien"`
	Comentario string `json:"comentario"`
}

func main() {
	puerto := flag.Int("puerto", 8733, "puerto del servidor")
	abrir := flag.Bool("abrir", true, "abrir el navegador al arrancar")
	red := flag.Bool("red", false, "aceptar conexiones de otros equipos de la red local")
	flag.Parse()

	sub, _ := fs.Sub(contenido, "web")
	http.Handle("/", http.FileServer(http.FS(sub)))
	http.HandleFunc("/api/carta", apiCarta)
	http.HandleFunc("/api/verificar", apiVerificar)
	http.HandleFunc("/api/lectura", apiLectura)
	http.HandleFunc("/api/lugares", apiLugares)
	http.HandleFunc("/api/huso", apiHuso)
	http.HandleFunc("/api/husohistoria", apiHistoria)
	http.HandleFunc("/api/guardadas", apiGuardadas)

// Si el puerto está ocupado se prueban los siguientes: una app de escritorio no
	// debe morirse porque otro programa tenga tomado un puerto.
	var ln net.Listener
	var err error
	usado := *puerto
	for i := 0; i < 40; i++ {
		host := "127.0.0.1"
		if *red {
			host = "0.0.0.0"
		}
		ln, err = net.Listen("tcp", fmt.Sprintf("%s:%d", host, usado))
		if err == nil {
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
	fmt.Printf("\n  Curso de Astrología\n  %s\n%s", dir, aviso)
	if *red {
		fmt.Println("\n  Abierto a la red local. Desde otro equipo o desde el móvil:")
		for _, ip := range ipsLocales() {
			fmt.Printf("    http://%s:%d\n", ip, usado)
		}
		fmt.Println("\n  AVISO: cualquiera que esté en esta red puede entrar y ver o borrar")
		fmt.Println("  tus cartas guardadas. No lo dejes puesto en una red pública.")
	}
	fmt.Println("\n  Ctrl-C para parar.\n")
	if *abrir {
		go func() { time.Sleep(400 * time.Millisecond); abrirNavegador(dir) }()
	}
	if err := http.Serve(ln, nil); err != nil {
		fmt.Println("  error:", err)
	}
}

// ipsLocales devuelve las direcciones IPv4 de la máquina en la red local.
func ipsLocales() []string {
	var out []string
	dirs, err := net.InterfaceAddrs()
	if err != nil {
		return out
	}
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
func ent(r *http.Request, k string) int {
	v, _ := strconv.Atoi(r.URL.Query().Get(k))
	return v
}

func apiCarta(w http.ResponseWriter, r *http.Request) {
	c := efem.Calcular(ent(r, "anio"), ent(r, "mes"), ent(r, "dia"),
		ent(r, "hh"), ent(r, "mm"), num(r, "tz"), num(r, "lat"), num(r, "lon"))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(c)
}

// apiVerificar corrige el cálculo a mano del módulo 3, paso a paso.
// Esto es aritmética pura: no hace falta ninguna IA para decir dónde falló.
func apiLugares(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(lugares.Buscar(r.URL.Query().Get("q"), 8))
}

// apiHuso resuelve el desfase REAL de un lugar en un instante concreto, con las
// reglas históricas de la base IANA. Es el paso que el módulo 2 llama la trampa
// número uno: nadie debería teclear el huso a mano.
func apiHuso(w http.ResponseWriter, r *http.Request) {
	off, nom, verano, err := lugares.Huso(r.URL.Query().Get("zona"),
		ent(r, "anio"), ent(r, "mes"), ent(r, "dia"), ent(r, "hh"), ent(r, "mm"))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"offset": off, "zona": nom, "verano": verano})
}

// apiGuardadas gestiona las cartas guardadas: listar, guardar y borrar.
func apiGuardadas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch r.Method {
	case http.MethodPost:
		var c guardadas.Carta
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "datos inválidos", http.StatusBadRequest)
			return
		}
		if c.Nombre == "" {
			http.Error(w, "hace falta un nombre", http.StatusBadRequest)
			return
		}
		if _, err := guardadas.Guardar(c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		if err := guardadas.Borrar(r.URL.Query().Get("id")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(map[string]any{
		"cartas": guardadas.Listar(), "fichero": guardadas.Ruta()})
}

func apiHistoria(w http.ResponseWriter, r *http.Request) {
	h, err := lugares.HistoriaHuso(r.URL.Query().Get("zona"), ent(r, "anio"), ent(r, "mes"),
		ent(r, "dia"), ent(r, "hh"), ent(r, "mm"), num(r, "lon"))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(h)
}

func apiLectura(w http.ResponseWriter, r *http.Request) {
	c := efem.Calcular(ent(r, "anio"), ent(r, "mes"), ent(r, "dia"),
		ent(r, "hh"), ent(r, "mm"), num(r, "tz"), num(r, "lat"), num(r, "lon"))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(c.Interpretar())
}

func apiVerificar(w http.ResponseWriter, r *http.Request) {
	c := efem.Calcular(ent(r, "anio"), ent(r, "mes"), ent(r, "dia"),
		ent(r, "hh"), ent(r, "mm"), num(r, "tz"), num(r, "lat"), num(r, "lon"))

	comprueba := func(nombre string, tuyo, correcto, tol float64, unidad string, angular bool) paso {
		d := tuyo - correcto
		if angular {
			d = math.Mod(d+540, 360) - 180
		}
		p := paso{Nombre: nombre, Tuyo: tuyo, Correcto: correcto,
			Desvio: math.Round(math.Abs(d)*10000) / 10000, Unidad: unidad,
			Bien: math.Abs(d) <= tol}
		switch {
		case p.Bien:
			p.Comentario = "correcto"
		case nombre == "Tiempo sidéreo local" && math.Abs(math.Abs(d)-math.Abs(num(r, "lon")/15)) < 0.05:
			p.Comentario = "te has dejado la corrección por longitud, o la has sumado al revés"
		case nombre == "Ascendente" && math.Abs(d) > 100:
			p.Comentario = "estás medio zodíaco desviado: revisa el cuadrante del arcotangente"
		default:
			p.Comentario = "revisa este paso antes de seguir"
		}
		return p
	}

	pasos := []paso{
		comprueba("Día juliano", num(r, "jd"), c.JD, 0.0007, "días", false),
		comprueba("T. sidéreo Greenwich", num(r, "tsg"), c.TSG, 0.05, "grados", true),
		comprueba("T. sidéreo local", num(r, "tsl"), c.TSL, 0.05, "grados", true),
		comprueba("Ascendente", num(r, "asc"), c.Asc, 0.05, "grados", true),
		comprueba("Medio Cielo", num(r, "mc"), c.MC, 0.05, "grados", true),
	}
	primero := -1
	for i, p := range pasos {
		if !p.Bien {
			primero = i
			break
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]any{"pasos": pasos, "primerFallo": primero})
}
