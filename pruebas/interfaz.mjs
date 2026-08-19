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

// Los datos se piden en cada idioma: parte del texto —los yogas, por ejemplo—
// lo compone el servidor y viaja ya traducido dentro del JSON.
const datos = {};
for (const lang of ["es", "en"]) {
  datos[lang] = {
    occidental: await (await fetch(`${base}/api/carta?${Q}&lang=${lang}`)).json(),
    jyotisha:   await (await fetch(`${base}/api/vedica?${Q}&lang=${lang}`)).json(),
  };
}
const occ = datos.es.occidental, ved = datos.es.jyotisha;

for (const lang of ["es", "en"]) {
  A.LANG = lang;
  for (const trad of ["occidental", "jyotisha"]) {
    A.TRAD = trad; A.DATOS = datos[lang][trad];
    const etiqueta = `${trad}/${lang}`;
    try { A.aplicarIdioma(); ok(`${etiqueta}: aplicarIdioma`); } catch (e) { mal(`${etiqueta}: aplicarIdioma`, e); }
    try { A.pintarNav();     ok(`${etiqueta}: pintarNav`); }     catch (e) { mal(`${etiqueta}: pintarNav`, e); }
    try { A.pintarCurso();   ok(`${etiqueta}: pintarCurso`); }   catch (e) { mal(`${etiqueta}: pintarCurso`, e); }
    try { A.render();        ok(`${etiqueta}: render`); }        catch (e) { mal(`${etiqueta}: render`, e); }
  }
}

// Ninguna traducción puede salir como "undefined" en la página.
A.LANG = "en"; A.TRAD = "jyotisha"; A.DATOS = datos.en.jyotisha;
A.render();
for (const id of ["pc", "fz", "vg", "ds", "tablas"]) {
  const html = doc.querySelector("#" + id).innerHTML;
  if (/undefined/.test(html)) mal(`#${id} contiene "undefined"`, html.match(/.{0,50}undefined.{0,30}/)[0]);
  else ok(`#${id} sin "undefined"`);
}

// ── que en inglés no se cuele ni una frase en castellano ─────────────────
//
// El chequeo no adivina qué es castellano: coge el propio diccionario. Si un
// texto sale distinto en los dos idiomas, entonces con la página en inglés la
// versión española no puede aparecer por ningún lado. Cero falsos positivos, y
// caza el caso que se escapa siempre: la cadena que se escribe una vez al
// arrancar y luego no se vuelve a traducir al cambiar de idioma.
function textosDe(obj, salida = []) {
  for (const v of Object.values(obj)) {
    if (typeof v === "string" && v.length > 12) salida.push(v);
    else if (v && typeof v === "object") textosDe(v, salida);
  }
  return salida;
}
const A2 = globalThis.__api;
const soloEs = [];
{
  const es = A2.t.call(null), enDic = null;
  // se accede a los diccionarios por el idioma activo
  A2.LANG = "es"; const dEs = A2.t();
  A2.LANG = "en"; const dEn = A2.t();
  const planos = (o, pre = "") => Object.entries(o).flatMap(([k, v]) =>
    typeof v === "string" ? [[pre + k, v]] : (v && typeof v === "object" ? planos(v, pre + k + ".") : []));
  const mEs = new Map(planos(dEs)), mEn = new Map(planos(dEn));
  for (const [k, v] of mEs) if (v.length > 4 && mEn.get(k) !== v) soloEs.push([k, v]);
}

// Se pinta todo en inglés, desde cero, y se mira el HTML resultante.
A2.LANG = "en";
const huecos = ["pc","fz","vg","ds","tablas","lec","cmp","lista","ficha","husoTxt","pie"];
for (const id of huecos) { const n = doc.querySelector("#" + id); n.innerHTML = ""; n.textContent = ""; }
// La página arranca en español y luego el usuario cambia de idioma. Se simula
// esa primera pasada, porque es justo ahí donde se esconde el fallo: la cadena
// que se escribe una vez y ya nadie vuelve a tocar.
A2.LANG = "es"; A2.TRAD = "occidental"; A2.DATOS = datos.es.occidental; A2.aplicarIdioma();
A2.LANG = "en";
for (const trad of ["occidental", "jyotisha"]) {
  A2.TRAD = trad; A2.DATOS = datos.en[trad];
  A2.aplicarIdioma(); A2.pintarNav(); A2.pintarCurso(); A2.render();
}
const paja = huecos.map(id => { const n = doc.querySelector("#" + id);
  return (n.innerHTML || "") + " " + (n.textContent || ""); }).join(" ");
const coladas = soloEs.filter(([, v]) => paja.includes(v));
if (coladas.length) coladas.forEach(([k, v]) =>
  mal(`en inglés se cuela el texto español de "${k}"`, v.slice(0, 60)));
else ok(`ningún texto español en la página en inglés (${soloEs.length} comprobados)`);

console.log(fallos ? `\n${fallos} fallo(s)` : "\ninterfaz correcta.");
process.exit(fallos ? 1 : 0);
