package guardadas

import (
	"os"
	"testing"
)

// Las cartas del usuario no se tocan: ASTRO_DIR manda el almacén a un
// directorio de usar y tirar.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "astro-test")
	if err != nil {
		panic(err)
	}
	os.Setenv("ASTRO_DIR", dir)
	cargar()
	cod := m.Run()
	os.RemoveAll(dir)
	os.Exit(cod)
}

func ejemplo(nombre string) Carta {
	return Carta{Nombre: nombre, Ciudad: "Barcelona", Zona: "Europe/Madrid",
		Fecha: "1961-12-19", Hora: "16:30", TZ: 1, Lat: 41.58, Lon: 2.55}
}

func TestGuardarListarBorrar(t *testing.T) {
	if n := len(Listar()); n != 0 {
		t.Fatalf("el almacén empieza con %d cartas y debería estar vacío", n)
	}

	a, err := Guardar(ejemplo("una"))
	if err != nil {
		t.Fatalf("no guarda: %v", err)
	}
	if a.ID == "" {
		t.Error("la carta guardada no recibe identificador")
	}
	if a.Creada == "" {
		t.Error("la carta guardada no recibe fecha")
	}

	b, _ := Guardar(ejemplo("otra"))
	if b.ID == a.ID {
		t.Error("dos cartas con el mismo identificador")
	}
	if n := len(Listar()); n != 2 {
		t.Fatalf("hay %d cartas y deberían ser 2", n)
	}

	// Los datos vuelven tal cual se metieron.
	for _, c := range Listar() {
		if c.Ciudad != "Barcelona" || c.Lat != 41.58 || c.Zona != "Europe/Madrid" {
			t.Errorf("la carta vuelve cambiada: %+v", c)
		}
	}

	if err := Borrar(a.ID); err != nil {
		t.Fatalf("no borra: %v", err)
	}
	l := Listar()
	if len(l) != 1 || l[0].ID != b.ID {
		t.Errorf("tras borrar quedan %d cartas y no es la que tocaba", len(l))
	}

	// Borrar algo que no existe no puede llevarse por delante lo que hay.
	_ = Borrar("no-existe")
	if len(Listar()) != 1 {
		t.Error("borrar un identificador inexistente ha cambiado el almacén")
	}
}

// Lo guardado tiene que sobrevivir a un reinicio: es un fichero, no memoria.
func TestSobreviveAlReinicio(t *testing.T) {
	for _, c := range Listar() {
		_ = Borrar(c.ID)
	}
	g, _ := Guardar(ejemplo("persistente"))

	todas = nil // se olvida lo que hay en memoria
	cargar()    // y se relee del disco, como al arrancar

	l := Listar()
	if len(l) != 1 || l[0].ID != g.ID || l[0].Nombre != "persistente" {
		t.Errorf("tras releer del disco quedan %d cartas: %+v", len(l), l)
	}
}
