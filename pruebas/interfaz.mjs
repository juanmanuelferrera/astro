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
globalThis.fetch = (u, o) => {
  const s = String(u);
  return fetchReal(/^https?:/.test(s) ? s : base + (s.startsWith("/") ? s : "/" + s), o);
};
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
  set DATOS(v){ DATOS = v }, get DATOS(){ return DATOS },
  repintarTodo, comparar, pintarEjercicio, abrirModulo, leer, predecir,
  set SECCION(v){ SECCION = v }, get SECCION(){ return SECCION } };
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
// El texto que el usuario ve de verdad: se quitan las etiquetas para que no
// cuenten los identificadores ni las clases —«Fecha» vive dentro del id
// «prFecha» y daba un falso positivo—, pero se rescatan antes los atributos
// que sí se leen: los title, los placeholder y las etiquetas de accesibilidad.
function visible(html) {
  if (!html) return "";
  const atributos = [...html.matchAll(/(?:title|placeholder|aria-label|alt)="([^"]*)"/g)]
    .map(m => m[1]).join(" ");
  return atributos + " " + html.replace(/<[^>]*>/g, " ");
}

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
const huecos = ["pc","fz","vg","ds","tablas","lec","cmp","prd","lista","ficha","husoTxt","pie"];
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
  return visible(n.innerHTML) + " " + (n.textContent || ""); }).join(" ");
const coladas = soloEs.filter(([, v]) => paja.includes(v));
if (coladas.length) coladas.forEach(([k, v]) =>
  mal(`en inglés se cuela el texto español de "${k}"`, v.slice(0, 60)));
else ok(`ningún texto español en la página en inglés (${soloEs.length} comprobados)`);

// ── el switch de idioma tiene que arrastrarlo todo ────────────────────────
//
// Esta es la prueba de verdad. Se monta la página entera en español —incluidas
// las pestañas que solo se pintan al abrirlas: comparar, el ejercicio, un
// módulo del curso abierto y el huso resuelto— y luego se toca el switch. Si
// después queda una sola frase en español, es que algo no se repinta.
{
  const A3 = globalThis.__api;
  for (const id of huecos.concat(["texto", "hist", "guardadas", "lista"])) {
    const n = doc.querySelector("#" + id); n.innerHTML = ""; n.textContent = "";
  }
  A3.LANG = "es"; A3.TRAD = "jyotisha"; A3.DATOS = datos.es.jyotisha;
  doc.querySelector("#zona").value = "Europe/Madrid";
  doc.querySelector("#fecha").value = "1961-12-19";
  doc.querySelector("#hora").value = "16:30";
  doc.querySelector("#lat").value = "41.58";
  doc.querySelector("#lon").value = "2.55";
  doc.querySelector("#tz").value = "1";
  doc.querySelector("#hist").hidden = false;

  A3.aplicarIdioma(); A3.render();
  await A3.comparar();
  A3.TRAD = "occidental"; A3.DATOS = datos.es.occidental; await A3.predecir("2015-06-01");
  A3.TRAD = "jyotisha"; A3.DATOS = datos.es.jyotisha;
  A3.pintarEjercicio();
  await A3.abrirModulo("01-el-cielo");
  await A3.leer();

  const antes = huecos.concat(["texto", "hist", "cmp", "prd", "ejBox"])
    .map(id => visible(doc.querySelector("#" + id).innerHTML)).join(" ");
  const habiaEs = soloEs.filter(([, v]) => antes.includes(v)).length;
  if (habiaEs < 5) mal("la página en español no se llegó a montar", `solo ${habiaEs} textos`);
  else ok(`página montada en español (${habiaEs} textos localizados)`);

  // …y ahora se toca el switch, que es lo único que hace el usuario.
  A3.LANG = "en";
  await A3.repintarTodo();

  const despues = huecos.concat(["texto", "hist", "cmp", "prd", "ejBox"])
    .map(id => { const n = doc.querySelector("#" + id);
                 return visible(n.innerHTML) + " " + (n.textContent || ""); }).join(" ");
  const quedan = soloEs.filter(([, v]) => despues.includes(v));
  if (quedan.length) {
    const dónde = id => visible(doc.querySelector("#" + id).innerHTML);
    quedan.slice(0, 6).forEach(([k, v]) => {
      const sitio = huecos.concat(["texto","hist","cmp","prd","ejBox"]).find(id => dónde(id).includes(v)) || "?";
      mal(`tras el switch queda español en #${sitio}`, `${k}: ${v.slice(0, 50)}`);
    });
  } else ok("tras tocar el switch no queda una sola frase en español");
}

// ── el switch no debe echarte de la pestaña donde estás ───────────────────
{
  const A5 = globalThis.__api;
  A5.TRAD = "jyotisha"; A5.DATOS = datos.es.jyotisha; A5.LANG = "es";
  for (const p of ["curso", "fuerza", "dasas", "comparar"]) {
    A5.SECCION = p;
    A5.LANG = "en"; await A5.repintarTodo();
    if (A5.SECCION !== p) mal(`el switch te saca de la pestaña "${p}"`, `acabas en "${A5.SECCION}"`);
    else ok(`sigues en "${p}" tras cambiar de idioma`);
    A5.LANG = "es"; await A5.repintarTodo();
  }
  // Y al cambiar de tradición, si la pestaña no existe allí, vuelve a la primera.
  A5.SECCION = "pancanga"; A5.TRAD = "occidental"; A5.DATOS = datos.es.occidental;
  await A5.repintarTodo();
  if (A5.SECCION === "carta") ok('"pancanga" no existe en occidental: vuelve a "carta"');
  else mal("al cambiar de tradición se queda en una pestaña que no existe", A5.SECCION);
}

console.log(fallos ? `\n${fallos} fallo(s)` : "\ninterfaz correcta.");
process.exit(fallos ? 1 : 0);
