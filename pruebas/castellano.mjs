// Busca castellano escrito a mano en el JavaScript de la interfaz.
//
// La otra prueba compara contra el diccionario, así que solo ve los textos que
// ESTÁN en el diccionario. Una frase escrita directamente dentro de una
// plantilla no está en ninguna parte y se le escapa entera — que es justo como
// se colaron la línea del ayanāṁśa en la pestaña de comparar y la nota del
// karakāṁśa en la de carta.
//
// Aquí se mira el código fuente. No se intenta reconocer dónde empieza y acaba
// cada cadena: las plantillas ocupan varias líneas y eso no se resuelve línea a
// línea. Se quitan comentarios y expresiones ${...}, y en lo que queda se
// cuentan palabras que solo existen en castellano. Con dos en la misma línea ya
// no es código: es una frase.

import fs from "node:fs";

const palabras = ["el","la","los","las","del","que","con","para","una","por",
  "más","está","están","son","sus","esa","ésa","ese","esta","este","toda","todo",
  "entre","cada","sin","sobre","como","desde","hasta","pero","porque","cuando",
  "donde","aquí","según","también","tras","es","se","su","al","lo","un","y","o",
  "no","hay","ser","dos","tres"];

// Palabras que también son inglesas o sánscritas y no prueban nada por sí solas.
const ambiguas = new Set(["no","la","el","son","es","o","a"]);

let fallos = 0;
const fichero = "web/app.js";
const bruto = fs.readFileSync(fichero, "utf8")
  .replace(/\/\*[\s\S]*?\*\//g, "");          // comentarios de bloque

bruto.split("\n").forEach((linea, i) => {
  const l = linea
    .replace(/\/\/.*$/, "")                    // comentarios de línea
    .replace(/\$\{[^{}]*\}/g, " · ");          // expresiones interpoladas
  const halladas = new Set();
  for (const p of palabras) {
    if (new RegExp(`(^|[\\s>(,.;:¿¡"'\`—–])${p}([\\s<),.;:?!—–\`"']|$)`, "i").test(l)) halladas.add(p);
  }
  const fuertes = [...halladas].filter(p => !ambiguas.has(p));
  if (halladas.size < 3 || fuertes.length < 2) return;
  console.log(`  ✗ ${fichero}:${i + 1} — castellano a mano: «${fuertes.slice(0, 4).join("», «")}»`);
  console.log(`      ${linea.trim().slice(0, 100)}`);
  fallos++;
});

console.log(fallos ? `\n  ${fallos} línea(s) con castellano sin traducir` :
  "  ✓ sin castellano escrito a mano en app.js");
process.exit(fallos ? 1 : 0);
