// jyotish — carta védica y curso, en un binario sin dependencias.
package main

import (
	"astro/internal/guardadas"
	"astro/internal/jyotisha"
	"astro/internal/lugares"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
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

func main() {
	puerto := flag.Int("puerto", 8744, "puerto del servidor")
	abrir := flag.Bool("abrir", true, "abrir el navegador al arrancar")
	red := flag.Bool("red", false, "aceptar conexiones de otros equipos de la red local")
	flag.Parse()

	sub, _ := fs.Sub(contenido, "web")
	http.Handle("/", http.FileServer(http.FS(sub)))
	http.HandleFunc("/api/carta", apiCarta)
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
		fmt.Printf("\n  No he podido abrir ningún puerto entre %d y %d.\n\n", *puerto, usado)
		os.Exit(1)
	}
	dir := fmt.Sprintf("http://localhost:%d", usado)
	aviso := ""
	if usado != *puerto {
		aviso = fmt.Sprintf("  (el puerto %d estaba ocupado)\n", *puerto)
	}
	fmt.Printf("\n  Jyotiṣa — carta védica\n  %s\n%s", dir, aviso)
	if *red {
		fmt.Println("\n  Abierto a la red local:")
		for _, ip := range ipsLocales() {
			fmt.Printf("    http://%s:%d\n", ip, usado)
		}
		fmt.Println("\n  AVISO: cualquiera en esta red puede ver o borrar tus cartas.")
	}
	fmt.Println("\n  Ctrl-C para parar.\n")
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

func num(r *http.Request, k string) float64 { v, _ := strconv.ParseFloat(r.URL.Query().Get(k), 64); return v }
func ent(r *http.Request, k string) int     { v, _ := strconv.Atoi(r.URL.Query().Get(k)); return v }

func apiCarta(w http.ResponseWriter, r *http.Request) {
	c := jyotisha.Calcular(ent(r, "anio"), ent(r, "mes"), ent(r, "dia"),
		ent(r, "hh"), ent(r, "mm"), num(r, "tz"), num(r, "lat"), num(r, "lon"))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(c)
}

func apiLugares(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(lugares.Buscar(r.URL.Query().Get("q"), 8))
}

func apiHuso(w http.ResponseWriter, r *http.Request) {
	off, nom, verano, err := lugares.Huso(r.URL.Query().Get("zona"),
		ent(r, "anio"), ent(r, "mes"), ent(r, "dia"), ent(r, "hh"), ent(r, "mm"))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"offset": off, "zona": nom, "verano": verano})
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

func apiGuardadas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
	json.NewEncoder(w).Encode(map[string]any{"cartas": guardadas.Listar(), "fichero": guardadas.Ruta()})
}
