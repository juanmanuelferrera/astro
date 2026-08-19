package main

import (
	"astro/internal/efem"
	"fmt"
	"os"
	"strconv"
)

func main() {
	jd, _ := strconv.ParseFloat(os.Args[1], 64)
	nombres := []string{"Mercurio", "Venus", "Marte", "Júpiter", "Saturno", "Urano", "Neptuno", "Plutón"}
	fmt.Printf("%.6f %.6f %.6f", efem.Sol(jd), efem.Luna(jd), efem.NodoLunarMedio(jd))
	for _, n := range nombres {
		fmt.Printf(" %.6f", efem.Planeta(n, jd))
	}
	fmt.Println()
}
