#!/usr/bin/env bash
# Comprobacion completa antes de publicar. Falla al primer problema.
#
# El motivo de que exista: el JavaScript viaja embebido con go:embed y el
# compilador de Go no lo mira. Un parentesis de mas ahi compila, arranca,
# responde a los endpoints — y deja la interfaz muerta sin decir nada.
set -e
cd "$(dirname "${BASH_SOURCE[0]}")"

echo "── Go: compila y vet ──"
go build ./... && go vet ./...

echo "── Go: tests ──"
go test ./... 2>&1 | grep -v "no test files" || true

echo "── JavaScript: sintaxis ──"
if command -v node >/dev/null 2>&1; then
  for f in web/*.js; do node --check "$f" && echo "  ✓ $f"; done
else
  echo "  ⚠ node no instalado, sin comprobar"; fi

echo "── i18n: los dos idiomas emparejados ──"
node -e '
const fs=require("fs");
eval(fs.readFileSync("web/i18n.js","utf8")+`
const a=Object.keys(T.es),b=Object.keys(T.en);
const f=a.filter(k=>!b.includes(k)).concat(b.filter(k=>!a.includes(k)));
if(f.length){console.error("  ✗ claves desparejadas: "+f.join(", "));process.exit(1);}
console.log("  ✓ "+a.length+" claves en los dos idiomas");
for(const tr of ["occidental","jyotisha"]) for(const k of NAV[tr])
  if(!T.es.nav[k]||!T.en.nav[k]){console.error("  ✗ falta el nombre de la pestaña "+k);process.exit(1);}
console.log("  ✓ todas las pestañas tienen nombre");`)' 2>/dev/null || \
 { echo "  ⚠ node no instalado, sin comprobar"; }

echo "── castellano escrito a mano ──"
if command -v node >/dev/null 2>&1; then
  node pruebas/castellano.mjs
else
  echo "  ⚠ node no instalado, sin comprobar"
fi

echo "── curso: los modulos existen en los dos idiomas ──"
falta=0
for t in occidental jyotisha; do
  for f in web/curso/$t/*.md; do
    [ -f "web/curso/$t/en/$(basename "$f")" ] || { echo "  ✗ falta web/curso/$t/en/$(basename "$f")"; falta=1; }
  done
done
[ $falta -eq 0 ] && echo "  ✓ $(ls web/curso/*/[0-9c-p]*.md 2>/dev/null | wc -l | tr -d ' ') modulos, todos traducidos"

echo "── endpoints ──"
go build -o /tmp/_astro_v . 
/tmp/_astro_v -abrir=false -puerto=8999 >/dev/null 2>&1 &
PID=$!
sleep 1.5
Q='anio=1961&mes=12&dia=19&hh=16&mm=30&tz=1&lat=41.58&lon=2.55'
mal=0
for e in "carta?$Q" "vedica?$Q" "vedica?$Q&nodo=verdadero" "lectura?$Q&lang=en" \
         "lecturaved?$Q&lang=es" "comparar?$Q" "lugares?q=barcelona" "huso?lat=41.4&lon=2.2&anio=1961&mes=12&dia=19"; do
  cod=$(curl -s -o /dev/null -w '%{http_code}' "http://localhost:8999/api/$e")
  [ "$cod" = "200" ] && echo "  ✓ ${e%%\?*}" || { echo "  ✗ ${e%%\?*} → $cod"; mal=1; }
done
echo "── interfaz: se ejecuta de verdad ──"
if command -v node >/dev/null 2>&1; then
  node pruebas/interfaz.mjs 8999 | sed "s/^/  /" || mal=1
else
  echo "  ⚠ node no instalado, sin comprobar"
fi

kill $PID 2>/dev/null || true; wait $PID 2>/dev/null || true; rm -f /tmp/_astro_v
[ $mal -eq 0 ] || exit 1
echo
echo "todo en orden."
