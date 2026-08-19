#!/usr/bin/env python3
"""Genera el icono de Astro: una rueda de carta, dibujada sin dependencias."""
import math, os, subprocess, sys, tempfile, zlib, struct

def png(w, h, pixels):
    def chunk(t, d):
        c = t + d
        return struct.pack('>I', len(d)) + c + struct.pack('>I', zlib.crc32(c) & 0xffffffff)
    raw = b''.join(b'\x00' + bytes(row) for row in pixels)
    return (b'\x89PNG\r\n\x1a\n'
            + chunk(b'IHDR', struct.pack('>IIBBBBB', w, h, 8, 6, 0, 0, 0))
            + chunk(b'IDAT', zlib.compress(raw, 9)) + chunk(b'IEND', b''))

def dibuja(S):
    cx = cy = S / 2
    F = 4                                   # muestreo para bordes suaves
    fondo = (250, 247, 242); tinta = (26, 22, 17); ac = (154, 91, 30)
    filas = []
    R, RI = S * 0.415, S * 0.30
    for y in range(S):
        fila = []
        for x in range(S):
            sr = sg = sb = sa = 0
            for oy in range(F):
                for ox in range(F):
                    px, py = x + (ox + .5) / F, y + (oy + .5) / F
                    d = math.hypot(px - cx, py - cy)
                    col, al = fondo, 0
                    if d <= S * 0.47:
                        col, al = fondo, 255
                        ang = math.degrees(math.atan2(cy - py, px - cx)) % 360
                        # doce radios
                        for k in range(12):
                            a = k * 30
                            dif = abs(((ang - a + 180) % 360) - 180)
                            if RI < d < R and dif < math.degrees(math.atan2(S*0.006, max(d,1))):
                                col = tinta
                        if abs(d - R) < S * 0.008 or abs(d - RI) < S * 0.006:
                            col = tinta
                        # eje del horizonte, en color
                        if abs(py - cy) < S * 0.011 and d < R * 1.06:
                            col = ac
                        if abs(px - cx) < S * 0.009 and d < R * 1.06:
                            col = (45, 96, 114)
                    sr += col[0]*al//255; sg += col[1]*al//255; sb += col[2]*al//255; sa += al
            n = F * F
            fila += [sr//n, sg//n, sb//n, sa//n]
        filas.append(fila)
    return filas

def main(destino):
    d = tempfile.mkdtemp(); ic = os.path.join(d, 'astro.iconset'); os.makedirs(ic)
    for s, nom in [(16,'16x16'),(32,'16x16@2x'),(32,'32x32'),(64,'32x32@2x'),
                   (128,'128x128'),(256,'128x128@2x'),(256,'256x256'),
                   (512,'256x256@2x'),(512,'512x512'),(1024,'512x512@2x')]:
        open(os.path.join(ic, f'icon_{nom}.png'), 'wb').write(png(s, s, dibuja(s)))
    subprocess.run(['iconutil','-c','icns',ic,'-o',destino], check=True)

if __name__ == '__main__':
    main(sys.argv[1])
