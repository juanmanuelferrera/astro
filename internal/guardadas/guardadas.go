// Package guardadas — persistencia de cartas en el disco del usuario.
// Se guarda en la carpeta de configuración del sistema, no en el navegador:
// así sobrevive a limpiar datos y el fichero se puede copiar o respaldar.
package guardadas

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Carta struct {
	ID     string  `json:"id"`
	Nombre string  `json:"nombre"`
	Ciudad string  `json:"ciudad"`
	Zona   string  `json:"zona"`
	Fecha  string  `json:"fecha"`
	Hora   string  `json:"hora"`
	TZ     float64 `json:"tz"`
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Creada string  `json:"creada"`
}

var (
	mu    sync.Mutex
	ruta  string
	todas []Carta
)

func Ruta() string { return ruta }

func init() {
	base, err := os.UserConfigDir()
	if err != nil {
		base, _ = os.UserHomeDir()
	}
	dir := filepath.Join(base, "astro")
	os.MkdirAll(dir, 0o755)
	ruta = filepath.Join(dir, "cartas.json")
	if b, err := os.ReadFile(ruta); err == nil {
		json.Unmarshal(b, &todas)
	}
}

func grabar() error {
	b, err := json.MarshalIndent(todas, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ruta, b, 0o644)
}

// Listar devuelve las cartas, la más reciente primero.
func Listar() []Carta {
	mu.Lock()
	defer mu.Unlock()
	r := make([]Carta, len(todas))
	copy(r, todas)
	sort.Slice(r, func(i, j int) bool { return r[i].Creada > r[j].Creada })
	if r == nil {
		r = []Carta{}
	}
	return r
}

// Guardar añade una carta nueva, o sustituye la que tenga el mismo nombre.
func Guardar(c Carta) (Carta, error) {
	mu.Lock()
	defer mu.Unlock()
	c.Creada = time.Now().Format(time.RFC3339)
	for i, v := range todas {
		if v.Nombre == c.Nombre {
			c.ID = v.ID
			todas[i] = c
			return c, grabar()
		}
	}
	c.ID = strconv.FormatInt(time.Now().UnixNano(), 36)
	todas = append(todas, c)
	return c, grabar()
}

// Borrar elimina una carta por su id.
func Borrar(id string) error {
	mu.Lock()
	defer mu.Unlock()
	for i, v := range todas {
		if v.ID == id {
			todas = append(todas[:i], todas[i+1:]...)
			return grabar()
		}
	}
	return nil
}
