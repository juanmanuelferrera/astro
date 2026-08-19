#!/usr/bin/env python3
"""Genera la tabla de referencia contra la que se mide el motor.

Los números salen de Swiss Ephemeris, que es la implementación con la que se
compara todo el mundo. Se guardan clavados en un fichero Go y se meten en el
repositorio: así el test no compara el motor consigo mismo —eso solo detecta
cambios, no errores— sino contra una referencia externa que no se mueve.

Solo hay que volver a ejecutarlo si se amplía la cobertura de fechas. No para
"arreglar" un test que falla: si falla, el que está mal es el motor.

    ~/.claude/skills/astrologia/.venv/bin/python pruebas/generar_referencia.py
"""
import swisseph as swe

swe.set_sid_mode(swe.SIDM_LAHIRI, 0, 0)
FLAGS = swe.FLG_SWIEPH
try:
    swe.calc_ut(2451545.0, swe.SUN, FLAGS)
except Exception:
    FLAGS = swe.FLG_MOSEPH   # sin ficheros de efemérides, Moshier basta

CUERPOS = [("Sol", swe.SUN), ("Luna", swe.MOON), ("Mercurio", swe.MERCURY),
           ("Venus", swe.VENUS), ("Marte", swe.MARS), ("Júpiter", swe.JUPITER),
           ("Saturno", swe.SATURN), ("Urano", swe.URANUS), ("Neptuno", swe.NEPTUNE),
           ("Plutón", swe.PLUTO)]

# Fechas repartidas por tres siglos, más algunos casos que suelen romper cosas.
FECHAS = []
for anio in range(1800, 2101, 15):
    for mes, dia, hh in ((1, 15, 6), (5, 3, 14), (9, 22, 21)):
        FECHAS.append((anio, mes, dia, hh, 30))
FECHAS += [
    (2000, 1, 1, 12, 0),    # J2000, la época de referencia
    (1900, 1, 1, 0, 0),
    (2026, 8, 19, 12, 0),
    (1961, 12, 19, 15, 30),
    (2024, 2, 29, 3, 15),   # bisiesto
    (1999, 12, 31, 23, 59), # cambio de año, mes y siglo
]

# Sitios: ecuador, latitudes medias norte y sur, y una alta donde Plácido sufre.
LUGARES = [("Quito", -0.18, -78.47), ("Barcelona", 41.39, 2.17),
           ("Sídney", -33.87, 151.21), ("Oslo", 59.91, 10.75)]

def jd_de(a, m, d, hh, mi):
    return swe.julday(a, m, d, hh + mi / 60.0)

sal = []
w = sal.append
w("package efem")
w("")
w("// GENERADO por pruebas/generar_referencia.py — no editar a mano.")
w("//")
w("// Valores de Swiss Ephemeris %s. Son la referencia externa contra la que" % swe.version)
w("// se mide el motor: comparar el motor consigo mismo solo detectaría cambios,")
w("// no errores. Si un test falla, el que está mal es el motor.")
w("")
w("type refCuerpo struct {")
w("\tNombre string")
w("\tLon    float64 // longitud eclíptica aparente de la fecha, en grados")
w("}")
w("")
w("type refFecha struct {")
w("\tAnio, Mes, Dia, HH, MM int")
w("\tJD                     float64")
w("\tAyanamsa               float64 // Lahiri")
w("\tNodoMedio, NodoVerdad  float64")
w("\tCuerpos                []refCuerpo")
w("}")
w("")
w("type refCasas struct {")
w("\tJD        float64")
w("\tLat, Lon  float64")
w("\tAsc, MC   float64")
w("\tCuspides  [12]float64 // Plácido")
w("}")
w("")
w("var refFechas = []refFecha{")
for (a, m, d, hh, mi) in FECHAS:
    jd = jd_de(a, m, d, hh, mi)
    ayan = swe.get_ayanamsa_ut(jd)
    nm = swe.calc_ut(jd, swe.MEAN_NODE, FLAGS)[0][0]
    nv = swe.calc_ut(jd, swe.TRUE_NODE, FLAGS)[0][0]
    w("\t{%d, %d, %d, %d, %d, %.9f, %.9f, %.9f, %.9f, []refCuerpo{" % (a, m, d, hh, mi, jd, ayan, nm, nv))
    for nombre, cuerpo in CUERPOS:
        lon = swe.calc_ut(jd, cuerpo, FLAGS)[0][0]
        w('\t\t{"%s", %.9f},' % (nombre, lon))
    w("\t}},")
w("}")
w("")
w("var refCasasTabla = []refCasas{")
for (a, m, d, hh, mi) in FECHAS[::4]:
    jd = jd_de(a, m, d, hh, mi)
    for nombre, lat, lon in LUGARES:
        cusps, ascmc = swe.houses_ex(jd, lat, lon, b'P')
        w("\t{%.9f, %.6f, %.6f, %.9f, %.9f, [12]float64{%s}}," % (
            jd, lat, lon, ascmc[0], ascmc[1],
            ", ".join("%.9f" % c for c in cusps[:12])))
w("}")
w("")

# Orto y ocaso, que se comprueban aparte porque usan otra rutina.
w("type refOrto struct {")
w("\tJD0      float64 // día juliano a 0h UT")
w("\tLat, Lon float64")
w("\tSalida   float64 // horas UT")
w("\tPuesta   float64")
w("}")
w("")
w("var refOrtos = []refOrto{")
for (a, m, d, _, _) in FECHAS[::6]:
    jd0 = swe.julday(a, m, d, 0.0)
    for nombre, lat, lon in LUGARES:
        try:
            # Sin BIT_DISC_CENTER: se quiere el criterio corriente —borde
            # superior del disco más refracción—, que es el mismo -0°50' que
            # usa el motor. Con DISC_CENTER swisseph mide el centro del disco y
            # los dos números dejan de ser comparables.
            r1 = swe.rise_trans(jd0, swe.SUN, swe.CALC_RISE,
                                (lon, lat, 0), 0, 0, FLAGS)
            r2 = swe.rise_trans(jd0, swe.SUN, swe.CALC_SET,
                                (lon, lat, 0), 0, 0, FLAGS)
            if r1[0] != 0 or r2[0] != 0:
                continue
            s = (r1[1][0] - jd0) * 24
            p = (r2[1][0] - jd0) * 24
            if not (0 <= s < 48 and 0 <= p < 48):
                continue
            w("\t{%.9f, %.6f, %.6f, %.6f, %.6f}," % (jd0, lat, lon, s, p))
        except Exception:
            continue
w("}")

open("internal/efem/referencia_test.go", "w", encoding="utf-8").write("\n".join(sal) + "\n")
print("internal/efem/referencia_test.go — %d fechas, %d juegos de casas" % (len(FECHAS), len(FECHAS[::4]) * len(LUGARES)))
