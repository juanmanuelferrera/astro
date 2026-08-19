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
