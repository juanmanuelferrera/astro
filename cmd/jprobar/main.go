package main

import (
	"astro/internal/jyotisha"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func main() {
	a, _ := strconv.Atoi(os.Args[1]); m, _ := strconv.Atoi(os.Args[2]); d, _ := strconv.Atoi(os.Args[3])
	h, _ := strconv.Atoi(os.Args[4]); mi, _ := strconv.Atoi(os.Args[5])
	tz, _ := strconv.ParseFloat(os.Args[6], 64)
	lat, _ := strconv.ParseFloat(os.Args[7], 64); lo, _ := strconv.ParseFloat(os.Args[8], 64)
	c := jyotisha.Calcular(a, m, d, h, mi, tz, lat, lo)
	if len(os.Args) > 9 {
		b, _ := json.Marshal(c); fmt.Println(string(b)); return
	}
	fmt.Printf("Lagna %s · %s pada %d · señor %s\n", c.LagnaPos, c.LagnaNak, c.LagnaPada, c.SenorLagna)
	fmt.Printf("ayanāṁśa %.6f°\n\n", c.Ayanamsa)
	for _, g := range c.Grahas {
		r := ""; if g.Retro { r = "℞" }
		gd := ""; if g.Gandanta { gd = " ⚠gaṇḍānta" }
		fmt.Printf("  %s %-9s %-16s %-13s pada %d  casa %-3d %-13s %s%s\n",
			g.Glifo, g.Nombre, g.Posicion, g.Nak, g.Pada, g.Bhava, g.Dignidad, r, gd)
	}
	fmt.Printf("\nkārakas: ")
	for _, k := range []string{"AK","AmK","BK","MK","PiK","PK","GK"} { fmt.Printf("%s=%s  ", k, c.Karakas[k]) }
	fmt.Println("\n\ndaśās:")
	for _, p := range c.Dasas[:4] {
		mk := ""; if p.Actual { mk = "  ← ACTUAL" }
		fmt.Printf("  %-9s %s → %s%s\n", p.Senor, p.Desde, p.Hasta, mk)
		if p.Actual { for _, b := range p.Sub { if b.Actual { fmt.Printf("      bhukti %s: %s → %s\n", b.Senor, b.Desde, b.Hasta) } } }
	}
	fmt.Println("\nyogas:"); for _, y := range c.Yogas { fmt.Println("  ·", y) }
	fmt.Println("\nD9 navāṁśa:")
	for _, g := range c.Vargas["D9"] { fmt.Printf("  %-8s %-11s casa %d\n", g.Nombre, g.Rasi, g.Bhava) }
}
