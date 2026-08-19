#!/usr/bin/env bash
# Construye TODO lo que se publica, en un solo paso.
#
# El motivo de que exista: los binarios de Windows y de Linux los venia
# compilando a mano en el terminal cada vez, y el zip de macOS igual. Lo que se
# hace a mano se olvida, y entonces la release lleva un fichero de tres
# versiones atras sin que nada avise. Ya paso con el paquete de Linux.
set -e
cd "$(dirname "${BASH_SOURCE[0]}")/.."
DEST="${1:-dist}"
mkdir -p "$DEST"

VERSION=$(go run . -version | awk '{print $2}')
echo "═══ astro $VERSION ═══"
echo
echo "═══ comprobaciones ═══"
./verificar.sh >/dev/null || { echo "verificar.sh falla — no se publica nada"; exit 1; }
echo "  todo en orden"

echo
echo "═══ macOS ═══"
./empaquetar/hacer-app.sh "$DEST" | sed 's/^/  /'

echo
echo "═══ Linux ═══"
./empaquetar/hacer-linux.sh "$DEST" | sed 's/^/  /'

echo
echo "═══ binarios sueltos ═══"
for par in linux/amd64 linux/arm64 windows/amd64; do
  os=${par%/*}; arq=${par#*/}
  sal="$DEST/astro-$os-$arq"; [ "$os" = "windows" ] && sal="$sal.exe"
  GOOS=$os GOARCH=$arq go build -ldflags="-s -w" -o "$sal" .
  echo "  $(basename "$sal")"
done

echo
echo "═══ listo — astro $VERSION ═══"
echo "  para publicar:  gh release create v$VERSION dist/astro-mac dist/Astro-mac-app.zip \\"
echo "                    dist/astro-linux.zip dist/astro-linux-* dist/astro-windows-amd64.exe"
ls -la "$DEST" | awk 'NR>1 && $5>10000 {printf "  %-26s %6.1f MB\n", $9, $5/1048576}'
