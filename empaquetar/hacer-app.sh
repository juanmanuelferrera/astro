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

echo "→ listo: $APP"
du -sh "$APP" | awk '{print "   tamaño: "$1}'
