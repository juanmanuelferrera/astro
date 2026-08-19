#!/usr/bin/env bash
# Construye Astro.app: un paquete de macOS que arranca el servidor en silencio,
# abre una ventana sin barra de direcciones y se apaga al cerrarla.
#
# Por dentro sigue siendo el mismo binario de Go. El .app no es un programa
# distinto: es un envoltorio para que el Finder, el Dock y Spotlight lo traten
# como una aplicacion.
set -e
cd "$(dirname "${BASH_SOURCE[0]}")/.."
DEST="${1:-dist}"
APP="$DEST/Astro.app"
rm -rf "$APP"
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"

# El JavaScript va embebido con go:embed, asi que el compilador de Go no lo
# mira: un error de sintaxis ahi pasa la compilacion y rompe la interfaz
# entera sin que nada avise. Se comprueba aqui antes de empaquetar nada.
if command -v node >/dev/null 2>&1; then
  echo "→ comprobando el JavaScript"
  for f in web/*.js; do
    node --check "$f" || { echo "ERROR de sintaxis en $f — no se empaqueta"; exit 1; }
  done
  # Y que ademas se ejecute: la sintaxis correcta no garantiza que pinte nada.
  # El binario se compila delante y solo despues se lanza al fondo; si se
  # encadenan con && y un solo &, el build tambien se va al fondo y la espera
  # se agota antes de que exista el binario.
  echo "→ comprobando que la interfaz se ejecuta"
  go build -o /tmp/_astro_p .
  /tmp/_astro_p -abrir=false -puerto=8998 >/dev/null 2>&1 &
  PRUEBA=$!
  for _ in $(seq 40); do
    curl -sf -o /dev/null http://localhost:8998/ && break
    sleep 0.2
  done
  if ! node pruebas/interfaz.mjs 8998 >/dev/null 2>&1; then
    echo "ERROR: la interfaz falla al ejecutarse — no se empaqueta"
    node pruebas/interfaz.mjs 8998 2>&1 | grep "✗" | head -5
    kill $PRUEBA 2>/dev/null; exit 1
  fi
  kill $PRUEBA 2>/dev/null || true
  wait $PRUEBA 2>/dev/null || true
  rm -f /tmp/_astro_p
else
  echo "⚠ node no esta instalado: el JavaScript no se ha comprobado"
fi

echo "→ compilando el binario universal"
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o /tmp/_a1 .
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o /tmp/_a2 .
lipo -create -output "$APP/Contents/MacOS/astro-bin" /tmp/_a1 /tmp/_a2
rm -f /tmp/_a1 /tmp/_a2

cat > "$APP/Contents/Info.plist" <<'PLIST'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
  <key>CFBundleName</key><string>Astro</string>
  <key>CFBundleDisplayName</key><string>Astro</string>
  <key>CFBundleIdentifier</key><string>cc.vedabase.astro</string>
  <key>CFBundleVersion</key><string>1.1.0</string>
  <key>CFBundleShortVersionString</key><string>1.1.0</string>
  <key>CFBundlePackageType</key><string>APPL</string>
  <key>CFBundleExecutable</key><string>Astro</string>
  <key>CFBundleIconFile</key><string>astro.icns</string>
  <key>LSMinimumSystemVersion</key><string>11.0</string>
  <!-- Sin icono en el Dock no se puede cerrar con Cmd-Q, asi que se deja. -->
  <key>LSUIElement</key><false/>
  <key>NSHighResolutionCapable</key><true/>
</dict></plist>
PLIST

# El ejecutable del bundle: arranca el servidor, abre la ventana y espera.
cat > "$APP/Contents/MacOS/Astro" <<'LAUNCH'
#!/bin/bash
AQUI="$(cd "$(dirname "$0")" && pwd)"
BIN="$AQUI/astro-bin"

# Puerto libre: si ya hay una instancia, se reutiliza en vez de levantar otra.
PUERTO=8733
for i in $(seq 0 40); do
  if ! nc -z 127.0.0.1 $((PUERTO+i)) 2>/dev/null; then PUERTO=$((PUERTO+i)); LIBRE=1; break; fi
  if curl -sf "http://127.0.0.1:$((PUERTO+i))/api/guardadas" >/dev/null 2>&1; then
    PUERTO=$((PUERTO+i)); YA=1; break
  fi
done
URL="http://localhost:$PUERTO"

if [ -z "$YA" ]; then
  "$BIN" -puerto=$PUERTO -abrir=false >/tmp/astro-app.log 2>&1 &
  SRV=$!
  for i in $(seq 1 60); do curl -sf "$URL/api/guardadas" >/dev/null 2>&1 && break; sleep 0.1; done
fi

# Ventana sin barra de direcciones. Se prueban los navegadores que la admiten;
# Safari no tiene modo aplicacion, asi que se abre en ventana propia.
ABIERTO=""
for NAV in "Google Chrome" "Brave Browser" "Microsoft Edge" "Chromium"; do
  if [ -d "/Applications/$NAV.app" ]; then
    open -na "$NAV" --args --app="$URL" --user-data-dir="$HOME/Library/Application Support/astro/ventana"
    ABIERTO=1; break
  fi
done
# Safari no tiene modo aplicacion, y `make new document` falla si Safari no
# esta ya arrancado (error -600). `open -a` lo lanza si hace falta.
if [ -z "$ABIERTO" ]; then
  open -a Safari "$URL" 2>/dev/null || open "$URL"
fi

# El proceso del bundle vive mientras viva el servidor: asi el icono se queda
# en el Dock y Cmd-Q lo cierra todo.
if [ -n "$SRV" ]; then
  trap 'kill $SRV 2>/dev/null' EXIT INT TERM
  wait $SRV
fi
LAUNCH
chmod +x "$APP/Contents/MacOS/Astro" "$APP/Contents/MacOS/astro-bin"

# Icono: una rueda sencilla, generada aqui para no depender de ficheros sueltos.
python3 empaquetar/icono.py "$APP/Contents/Resources/astro.icns" 2>/dev/null || \
  echo "  (sin icono: falta el generador)"

# El binario suelto, por si alguien lo prefiere sin envoltorio. El de dentro
# del paquete ya es universal, asi que basta copiarlo.
cp "$APP/Contents/MacOS/astro-bin" "$DEST/astro-mac"

# El zip lo hace el script, no la mano. Hacerlo aparte es como se queda viejo.
# ditto y no zip: un .app lleva enlaces simbolicos y atributos que zip pierde.
echo "→ armando el zip"
# El directorio de montaje se llama Astro y no tmp.loquesea: --keepParent
# conserva el nombre del padre, y con un mktemp pelado el zip acababa
# llevando dentro una carpeta llamada tmp.QUnAYL3oOo.
BASE=$(mktemp -d)
STAGE="$BASE/Astro"
mkdir -p "$STAGE"
cp -R "$APP" "$STAGE/"
cat > "$STAGE/LEEME.txt" <<'DOC'
Astro — cartas astrales, occidental y jyotiṣa

Arrastra Astro.app a Aplicaciones y ábrelo. Levanta el servidor solo y abre una
ventana sin barra de direcciones. Al cerrarla se apaga todo.


SI MACOS DICE QUE NO SE PUEDE ABRIR

La aplicación no está firmada con una cuenta de desarrollador de Apple, así que
al bajarla de internet macOS le pone una marca de cuarentena y se niega a
abrirla. Hay dos maneras de quitarla.

La rápida, en el Terminal — ojo al -r, que hace falta porque es un paquete y no
un fichero suelto:

    xattr -dr com.apple.quarantine /Applications/Astro.app

Si la tienes en otro sitio, pon esa ruta. El error «No such file» significa que
la ruta no es esa, no que el comando esté mal.

La otra, sin Terminal: Control-clic sobre Astro.app, «Abrir», y confirmar. Solo
la primera vez.


OPCIONES

Para las opciones hay que llamar al binario de dentro:

    /Applications/Astro.app/Contents/MacOS/astro-bin -puerto 9000
    /Applications/Astro.app/Contents/MacOS/astro-bin -red
    /Applications/Astro.app/Contents/MacOS/astro-bin -abrir=false

-red lo hace accesible desde otros equipos de la red local.

No necesita instalar nada más: lleva dentro las efemérides, el curso y las
ciudades. Las cartas guardadas quedan en ~/.astro/
DOC
rm -f "$DEST/Astro-mac-app.zip"
ditto -c -k --sequesterRsrc --keepParent "$STAGE" "$DEST/Astro-mac-app.zip"
rm -rf "$BASE"

echo "→ listo: $APP"
du -sh "$APP" | awk '{print "   app:    "$1}'
du -h "$DEST/Astro-mac-app.zip" | awk '{print "   zip:    "$1}'
du -h "$DEST/astro-mac" | awk '{print "   binario: "$1}'
