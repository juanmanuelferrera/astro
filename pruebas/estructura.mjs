// Comprueba que la página encaja consigo misma.
//
// Existe por un fallo que no cazaba nada: escribí var(--acento), var(--linea) y
// var(--texto), que no existen — las de verdad se llaman --ac, --rule y --ink.
// Una variable CSS que no está definida no da error: la declaración se vuelve
// inválida y la propiedad cae a su valor inicial. Para un `background` ese
// valor es transparente, así que las barras del pañcāṅga y del bhāva bala se
// dibujaban y no se veían. Compilaba, se ejecutaba, las pruebas pasaban.
//
// Aquí se miran las costuras: lo que el JavaScript pide y el HTML tiene, los
// identificadores repetidos, las pestañas sin sección, las variables de color
// y el equilibrio de las etiquetas.

import fs from "node:fs";

const html = fs.readFileSync("web/index.html", "utf8");
const app = fs.readFileSync("web/app.js", "utf8");
const i18n = fs.readFileSync("web/i18n.js", "utf8");

let fallos = 0;
const mal = (que, detalle) => { fallos++; console.log(`  ✗ ${que}\n      ${detalle}`); };
const bien = m => console.log(`  ✓ ${m}`);

// ── el JavaScript pide elementos que el HTML tiene ──
const idsHtml = new Set([...html.matchAll(/\bid="([^"]+)"/g)].map(m => m[1]));
const idsCreados = new Set([...app.matchAll(/\bid="([^"]+)"/g)].map(m => m[1]));
let huerfanos = 0;
for (const m of app.matchAll(/\$\("#([A-Za-z][\w-]*)"\)/g)) {
  const id = m[1];
  if (!idsHtml.has(id) && !idsCreados.has(id)) { mal(`app.js pide #${id} y no existe`, "ni en index.html ni creado por el propio JS"); huerfanos++; }
}
if (!huerfanos) bien("todos los elementos que pide el JavaScript existen");

// ── identificadores únicos ──
const cuenta = {};
for (const m of html.matchAll(/\bid="([^"]+)"/g)) cuenta[m[1]] = (cuenta[m[1]] || 0) + 1;
const repes = Object.entries(cuenta).filter(([, c]) => c > 1);
repes.forEach(([id, c]) => mal(`el id "${id}" aparece ${c} veces`, "querySelector solo vería el primero"));
if (!repes.length) bien("ningún identificador repetido");

// ── cada pestaña tiene su sección, y cada sección su pestaña ──
eval(i18n + ";globalThis.__NAV = NAV;");
const secciones = new Set([...html.matchAll(/<section id="([^"]+)"/g)].map(m => m[1]));
let sueltas = 0;
for (const [trad, lista] of Object.entries(globalThis.__NAV))
  for (const k of lista)
    if (!secciones.has(k)) { mal(`la pestaña "${k}" de ${trad} no tiene <section id="${k}">`, "al pulsarla no aparecería nada"); sueltas++; }
for (const s of secciones)
  if (!Object.values(globalThis.__NAV).some(l => l.includes(s))) { mal(`<section id="${s}"> no está en ningún NAV`, "queda inalcanzable"); sueltas++; }
if (!sueltas) bien(`las ${secciones.size} secciones encajan con las pestañas`);

// ── las variables de color existen ──
const definidas = new Set([...html.matchAll(/(--[\w-]+)\s*:/g)].map(m => m[1]));
const usadas = new Set([...(html + app).matchAll(/var\((--[\w-]+)\)/g)].map(m => m[1]));
const inventadas = [...usadas].filter(v => !definidas.has(v));
inventadas.forEach(v => mal(`se usa var(${v}) y no está definida`,
  "la declaración se vuelve inválida y la propiedad cae a su valor inicial: un fondo así queda transparente"));
if (!inventadas.length) bien(`las ${usadas.size} variables de color usadas están definidas`);

// ── el atributo hidden funciona de verdad ──
//
// hidden esconde porque la hoja del navegador trae [hidden]{display:none}. Esa
// regla tiene la misma especificidad que una clase, y el origen del autor gana
// al del navegador: basta con que la clase del elemento fije un display para
// que hidden deje de hacer nada, sin avisar. Le pasó a #estilosBox, que llevaba
// .estilos{display:flex} y enseñaba los estilos de carta védica en occidental.
//
// La cura es declarar [hidden] en la propia página con !important. Mientras
// esté, esto no puede repetirse con ningún elemento.
const declaraHidden = /\[hidden\][^{]*\{[^}]*display\s*:\s*none\s*!important/.test(html);
if (declaraHidden) {
  bien("la página declara [hidden]{display:none!important}: el atributo no lo puede tapar una clase");
} else {
  const conDisplay = {};
  for (const m of html.matchAll(/\.([\w-]+)\s*\{([^}]*)\}/g))
    if (/(^|;)\s*display\s*:/.test(m[2]))
      conDisplay[m[1]] = /display\s*:\s*([\w-]+)/.exec(m[2])[1];
  let choques = 0;
  for (const m of html.matchAll(/<\w+([^>]*\bhidden\b[^>]*)>/g)) {
    const attrs = m[1];
    if (attrs.includes('type="hidden"')) continue;
    const id = (/id="([^"]+)"/.exec(attrs) || [, "(sin id)"])[1];
    for (const c of ((/class="([^"]+)"/.exec(attrs) || [, ""])[1]).split(/\s+/).filter(Boolean))
      if (conDisplay[c]) {
        mal(`el atributo hidden de #${id} no hace nada`,
          `la clase .${c} pone display:${conDisplay[c]}, y el estilo del autor gana a [hidden] del navegador`);
        choques++;
      }
  }
  if (!choques) mal("la página no declara [hidden]{display:none!important}",
    "hoy no choca con ninguna clase, pero cualquier clase con display que se añada mañana romperá un hidden en silencio");
}

// ── la versión se escribe en un solo sitio ──
const go = fs.readFileSync("main.go", "utf8");
const ver = /const Version = "([^"]+)"/.exec(go);
if (!ver) mal("main.go no declara Version", "sin ella el pie no puede decir qué se está ejecutando");
else if (!/^\d+\.\d+\.\d+$/.test(ver[1]))
  mal(`la versión "${ver[1]}" no tiene la forma X.Y.Z`, "las etiquetas de publicación la usan tal cual");
else bien(`versión ${ver[1]}, declarada una sola vez en main.go`);

// ── las etiquetas cierran ──
const cuerpo = html.replace(/<script[\s\S]*?<\/script>/g, "").replace(/<style[\s\S]*?<\/style>/g, "");
const vacias = new Set(["br","hr","img","input","meta","link","source","path","circle","line","rect","use","col"]);
const pila = [];
let descuadre = 0;
for (const m of cuerpo.matchAll(/<(\/?)([a-zA-Z][\w-]*)([^>]*)>/g)) {
  const [, cierre, tag, resto] = m;
  if (vacias.has(tag.toLowerCase()) || resto.trimEnd().endsWith("/")) continue;
  if (cierre) {
    const ult = pila.pop();
    if (ult !== tag) { mal(`se cierra </${tag}> y estaba abierta <${ult}>`, "el navegador reordenaría el árbol a su manera"); descuadre++; }
  } else pila.push(tag);
}
if (pila.length) { mal(`quedan sin cerrar: ${pila.join(", ")}`, "en index.html"); descuadre++; }
if (!descuadre) bien("las etiquetas de index.html cierran todas");

console.log(fallos ? `\n  ${fallos} problema(s) de estructura` : "\n  la página encaja consigo misma");
process.exit(fallos ? 1 : 0);
