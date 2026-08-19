#!/usr/bin/env bash
# Arma el paquete que se le pasa a alguien con Linux: los dos binarios, un
# lanzador y las instrucciones. No lleva nada de macOS ni basura de __MACOSX.
set -e
cd "$(dirname "${BASH_SOURCE[0]}")/.."
DEST="${1:-dist}"
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT
mkdir -p "$DEST" "$TMP/astro"

echo "→ compilando para Linux"
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$TMP/astro/astro-amd64" .
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o "$TMP/astro/astro-arm64" .

cat > "$TMP/astro/astro" <<'LANZA'
#!/usr/bin/env bash
# Elige el binario que toca y abre el navegador.
cd "$(dirname "$0")"
case "$(uname -m)" in
  x86_64)          BIN=./astro-amd64 ;;
  aarch64|arm64)   BIN=./astro-arm64 ;;
  *) echo "Arquitectura no contemplada: $(uname -m)"; exit 1 ;;
esac
chmod +x "$BIN"
exec "$BIN" "$@"
LANZA
chmod +x "$TMP/astro/astro"

cat > "$TMP/astro/LEEME.txt" <<'DOC'
Astro — cartas astrales, occidental y jyotiṣa

Para arrancarlo:

    ./astro

Abre el navegador solo. Si no lo abre, la dirección sale en el terminal.

Opciones:

    ./astro -puerto 9000    otro puerto
    ./astro -red            accesible desde otros equipos de la red
    ./astro -abrir=false    no abrir el navegador

No necesita instalar nada: no tiene dependencias y lleva dentro las
efemérides, el curso y las ciudades.

Las cartas guardadas quedan en ~/.astro/
DOC

( cd "$TMP" && zip -qr "astro-linux.zip" astro -x '.*' -x '__MACOSX/*' )
mv "$TMP/astro-linux.zip" "$DEST/astro-linux.zip"
rm -f "$DEST/astro_linux.zip"   # el resto viejo, mal nombrado y con binarios de Mac
echo "→ listo: $DEST/astro-linux.zip"
unzip -l "$DEST/astro-linux.zip" | sed -n '4,9p'
