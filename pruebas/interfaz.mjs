// Prueba de la interfaz sin navegador.
//
// El JavaScript viaja embebido con go:embed y el compilador de Go no lo mira.
// `node --check` caza los errores de sintaxis, pero no los de ejecución: una
// clave de traducción que no existe, un elemento que se busca antes de estar
// en la página, un campo del JSON que se renombró en el servidor. Eso solo
// sale al ejecutar, y esto lo ejecuta.
//
// Se monta un DOM de mentira, se cargan i18n.js y app.js tal cual, y se llaman
// las funciones de pintado con datos de verdad traídos del servidor.

import fs from "node:fs";

const PUERTO = process.argv[2] || "8910";
const base = `http://localhost:${PUERTO}`;

// ── DOM de mentira ────────────────────────────────────────────────────────
const hechos = new Map();
function nuevoNodo(id = "") {
  const n = {
    id, hidden: false, title: "", value: "", textContent: "", innerHTML: "",
    checked: false, dataset: {}, style: {},
    classList: { add(){}, remove(){}, toggle(){}, contains: () => false },
    childNodes: [{ nodeValue: "" }],
    querySelectorAll: () => [], querySelector: () => nuevoNodo(),
    appendChild(){}, addEventListener(){}, scrollIntoView(){},
    getContext: () => null,
  };
  return n;
}
const doc = {
  querySelector(sel) {
    const id = sel.startsWith("#") ? sel.slice(1) : sel;
    if (!hechos.has(id)) hechos.set(id, nuevoNodo(id));
    return hechos.get(id);
  },
  querySelectorAll: () => [],
  getElementById(id) { return doc.querySelector("#" + id); },
  createElement: () => nuevoNodo(),
  addEventListener(){},
  documentElement: nuevoNodo(),
  body: nuevoNodo(),
};
globalThis.document = doc;
globalThis.window = { print(){}, matchMedia: () => ({ matches:false, addEventListener(){} }),
                      addEventListener(){}, location:{ href: base } };
globalThis.localStorage = { getItem: () => null, setItem(){}, removeItem(){} };
globalThis.alert = () => {};
// El código pide rutas relativas; aquí hay que anteponerle el servidor.
const fetchReal = globalThis.fetch;
globalThis.fetch = (u, o) => fetchReal(String(u).startsWith("/") ? base + u : u, o);
globalThis.prompt = () => null;
globalThis.confirm = () => true;

// ── cargar el código tal cual lo sirve el binario ─────────────────────────
const src = fs.readFileSync("web/i18n.js", "utf8") + "\n" + fs.readFileSync("web/app.js", "utf8");

let fallos = 0;
const ok = m => console.log("  ✓ " + m);
const mal = (m, e) => { fallos++; console.log("  ✗ " + m + "\n      " + (e && e.stack ? e.stack.split("\n")[0] : e)); };

// El código termina con llamadas de arranque; se ejecuta entero y se exponen
// las funciones que hay que probar.
const exportar = `
;globalThis.__api = { t, pintarNav, pintarCurso, aplicarIdioma, render,
  pintarPancanga, pintarFuerza, pintarVargas, pintarDasas,
  set TRAD(v){ TRAD = v }, get TRAD(){ return TRAD },
  set LANG(v){ LANG = v }, get LANG(){ return LANG },
  set DATOS(v){ DATOS = v }, get DATOS(){ return DATOS } };
`;
try { (0, eval)(src + exportar); ok("app.js se ejecuta entero"); }
catch (e) { console.log("  ✗ app.js revienta al cargar");
  console.log("    "+e.stack.split("\n").slice(0,3).join("\n    ")); process.exit(1); }

const A = globalThis.__api;
const Q = "anio=1961&mes=12&dia=19&hh=16&mm=30&tz=1&lat=41.58&lon=2.55";

const occ = await (await fetch(`${base}/api/carta?${Q}`)).json();
const ved = await (await fetch(`${base}/api/vedica?${Q}`)).json();

for (const lang of ["es", "en"]) {
  A.LANG = lang;
  for (const [trad, datos] of [["occidental", occ], ["jyotisha", ved]]) {
    A.TRAD = trad; A.DATOS = datos;
    const etiqueta = `${trad}/${lang}`;
    try { A.aplicarIdioma(); ok(`${etiqueta}: aplicarIdioma`); } catch (e) { mal(`${etiqueta}: aplicarIdioma`, e); }
    try { A.pintarNav();     ok(`${etiqueta}: pintarNav`); }     catch (e) { mal(`${etiqueta}: pintarNav`, e); }
    try { A.pintarCurso();   ok(`${etiqueta}: pintarCurso`); }   catch (e) { mal(`${etiqueta}: pintarCurso`, e); }
    try { A.render();        ok(`${etiqueta}: render`); }        catch (e) { mal(`${etiqueta}: render`, e); }
  }
}

// Ninguna traducción puede salir como "undefined" en la página.
A.LANG = "en"; A.TRAD = "jyotisha"; A.DATOS = ved;
A.render();
for (const id of ["pc", "fz", "vg", "ds", "tablas"]) {
  const html = doc.querySelector("#" + id).innerHTML;
  if (/undefined/.test(html)) mal(`#${id} contiene "undefined"`, html.match(/.{0,50}undefined.{0,30}/)[0]);
  else ok(`#${id} sin "undefined"`);
}

console.log(fallos ? `\n${fallos} fallo(s)` : "\ninterfaz correcta.");
process.exit(fallos ? 1 : 0);
