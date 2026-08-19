package efem

import "math"

// Sol devuelve la longitud eclíptica aparente del Sol en grados. Meeus cap. 25.
func Sol(jd float64) float64 {
	t := T(jd)
	L0 := 280.46646 + 36000.76983*t + 0.0003032*t*t
	M := norm360(357.52911 + 35999.05029*t - 0.0001537*t*t)
	mr := M * Grados
	C := (1.914602-0.004817*t-0.000014*t*t)*math.Sin(mr) +
		(0.019993-0.000101*t)*math.Sin(2*mr) +
		0.000289*math.Sin(3*mr)
	verdadera := L0 + C
	omega := 125.04 - 1934.136*t
	return norm360(verdadera - 0.00569 - 0.00478*math.Sin(omega*Grados))
}
