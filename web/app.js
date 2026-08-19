const $ = s => document.querySelector(s);
const SIG = ["♈","♉","♊","♋","♌","♍","♎","♏","♐","♑","♒","♓"];
const NOM = ["Aries","Tauro","Géminis","Cáncer","Leo","Virgo","Libra","Escorpio","Sagitario","Capricornio","Acuario","Piscis"];

// ── navegación ──
document.querySelectorAll("nav button").forEach(b => b.onclick = () => {
  document.querySelectorAll("nav button").forEach(x => x.classList.toggle("on", x === b));
  document.querySelectorAll("section").forEach(s => s.classList.toggle("on", s.id === b.dataset.s));
});

const gms = g => { const d = Math.floor(g), m = Math.round((g - d) * 60);
  return m === 60 ? `${d+1}° 00′` : `${d}° ${String(m).padStart(2,"0")}′`; };

// ── rueda ──
function rueda(c) {
  const CX = 250, CY = 250, R = 196, RS = 152, RH = 112, asc = c.angulos[0];
  const P = (lon, r) => { const a = (180 - (lon - asc)) * Math.PI / 180;
    return [CX + r * Math.cos(a), CY - r * Math.sin(a)]; };
  let s = `<svg viewBox="0 0 500 500" role="img" aria-label="Rueda de la carta natal">`;
  // sectores de signo
  for (let i = 0; i < 12; i++) {
    const l0 = i * 30, [x0,y0] = P(l0, RS), [x1,y1] = P(l0, R);
    s += `<line x1="${x0.toFixed(1)}" y1="${y0.toFixed(1)}" x2="${x1.toFixed(1)}" y2="${y1.toFixed(1)}" stroke="currentColor" stroke-width="1" opacity=".3"/>`;
    const [gx,gy] = P(l0 + 15, (R + RS) / 2);
    s += `<text x="${gx.toFixed(1)}" y="${(gy+5).toFixed(1)}" text-anchor="middle" font-size="17" fill="currentColor" opacity=".8">${SIG[i]}</text>`;
  }
  s += `<circle cx="${CX}" cy="${CY}" r="${R}" fill="none" stroke="currentColor" stroke-width="1.4" opacity=".45"/>`;
  s += `<circle cx="${CX}" cy="${CY}" r="${RS}" fill="none" stroke="currentColor" stroke-width="1.1" opacity=".35"/>`;
  s += `<circle cx="${CX}" cy="${CY}" r="${RH}" fill="none" stroke="currentColor" stroke-width="1" opacity=".2"/>`;
  // cúspides
  c.cuspP.forEach((cu, i) => {
    const [x0,y0] = P(cu, RH), [x1,y1] = P(cu, RS), ang = i % 3 === 0;
    s += `<line x1="${x0.toFixed(1)}" y1="${y0.toFixed(1)}" x2="${x1.toFixed(1)}" y2="${y1.toFixed(1)}" stroke="${ang?"var(--horizonte)":"currentColor"}" stroke-width="${ang?2.2:1}" opacity="${ang?.95:.28}"/>`;
    const [nx,ny] = P((cu + c.cuspP[(i+1)%12] + (((c.cuspP[(i+1)%12]-cu+360)%360)<180?0:360))/2, RH + 20);
    s += `<text x="${nx.toFixed(1)}" y="${(ny+4).toFixed(1)}" text-anchor="middle" font-size="12" fill="currentColor" opacity=".55">${i+1}</text>`;
  });
  // planetas: desapilado angular. Se separan los glifos hasta la distancia
  // mínima legible y se traza una línea desde el glifo a su grado real, que es
  // como lo resuelven los programas profesionales.
  const RG = RS - 26, MINSEP = 8.5;
  const cuerpos = c.cuerpos.map(p => ({ ...p, rel: ((p.lon - asc) % 360 + 360) % 360 }))
                           .sort((a, b) => a.rel - b.rel);
  cuerpos.forEach(p => p.dib = p.rel);
  for (let pasada = 0; pasada < 250; pasada++) {
    let movido = false;
    for (let i = 0; i < cuerpos.length; i++) {
      const a = cuerpos[i], b = cuerpos[(i + 1) % cuerpos.length];
      let d = b.dib - a.dib; if (d < 0) d += 360;
      if (d < MINSEP) {
        const emp = (MINSEP - d) / 2;
        a.dib -= emp; b.dib += emp; movido = true;
      }
    }
    if (!movido) break;
  }
  cuerpos.forEach(p => {
    const real = asc + p.rel, dib = asc + p.dib;
    const [tx,ty] = P(real, RS), [t2x,t2y] = P(real, RS - 7);
    s += `<line x1="${tx.toFixed(1)}" y1="${ty.toFixed(1)}" x2="${t2x.toFixed(1)}" y2="${t2y.toFixed(1)}" stroke="var(--horizonte)" stroke-width="1.4" opacity=".85"/>`;
    if (Math.abs(p.dib - p.rel) > 0.4) {
      const [cx1,cy1] = P(real, RS - 7), [cx2,cy2] = P(dib, RG + 9);
      s += `<line x1="${cx1.toFixed(1)}" y1="${cy1.toFixed(1)}" x2="${cx2.toFixed(1)}" y2="${cy2.toFixed(1)}" stroke="currentColor" stroke-width=".8" opacity=".3"/>`;
    }
    const [x,y] = P(dib, RG);
    s += `<text x="${x.toFixed(1)}" y="${(y+6).toFixed(1)}" text-anchor="middle" font-size="17" fill="var(--horizonte)">${p.glifo}</text>`;
    const [gx,gy] = P(dib, RG - 17);
    s += `<text x="${gx.toFixed(1)}" y="${(gy+4).toFixed(1)}" text-anchor="middle" font-size="9" font-family="ui-monospace,Menlo,monospace" fill="currentColor" opacity=".55">${Math.floor(p.grado)}°${p.retro?"℞":""}</text>`;
  });
  // ángulos
  [["Asc",c.angulos[0]],["Dsc",c.angulos[1]],["MC",c.angulos[2]],["IC",c.angulos[3]]].forEach(([n,l]) => {
    const [x,y] = P(l, R + 16);
    s += `<text x="${x.toFixed(1)}" y="${(y+4).toFixed(1)}" text-anchor="middle" font-size="10.5" font-family="ui-monospace,Menlo,monospace" font-weight="700" fill="${n==="Asc"||n==="Dsc"?"var(--horizonte)":"var(--meridiano)"}">${n}</text>`;
  });
  return s + "</svg>";
}

function tablas(c) {
  const ang = ["Ascendente","Descendente","Medio Cielo","Fondo del Cielo"];
  let h = `<div class="caja"><h3>Los cuatro ángulos</h3><table><tbody>`;
  c.angulos.forEach((l,i) => { const si = Math.floor(l/30);
    h += `<tr><td>${ang[i]}</td><td class="num">${gms(l-si*30)}</td><td class="gl">${SIG[si]}</td><td>${NOM[si]}</td></tr>`; });
  h += `</tbody></table></div><div class="caja"><h3>Planetas</h3><table>
    <thead><tr><th colspan="2">Planeta</th><th class="num">Grados</th><th colspan="2">Signo</th><th class="num">Casa</th></tr></thead><tbody>`;
  c.cuerpos.forEach(p => h += `<tr><td class="gl">${p.glifo}</td><td>${p.nombre}${p.retro?' <span style="color:var(--horizonte)">℞</span>':''}</td>
    <td class="num">${gms(p.grado)}</td><td class="gl">${p.glifoSig}</td><td>${p.signo}</td><td class="num">${p.casaP}</td></tr>`);
  h += `</tbody></table></div><div class="caja"><h3>Aspectos, por exactitud</h3><table><tbody>`;
  c.aspectos.slice(0,14).forEach(a => h += `<tr><td>${a.a}</td><td class="gl">${a.glifo}</td><td>${a.nombre}</td><td>${a.b}</td><td class="num">${a.orbe.toFixed(2)}°</td></tr>`);
  h += `</tbody></table></div><div class="caja"><h3>Regentes de casa</h3><table>
    <thead><tr><th>Casa</th><th>Signo</th><th>Regente</th><th>alojado en</th></tr></thead><tbody>`;
  c.regentes.forEach((r,i) => { const si = Math.floor(c.cuspP[i]/30);
    h += `<tr><td class="num">${i+1}</td><td class="gl">${SIG[si]}</td><td>${r}</td><td>casa ${c.regenteEn[i]}</td></tr>`; });
  return h + `</tbody></table></div>`;
}

function params(pre="") {
  const [a,m,d] = $("#"+pre+"fecha").value.split("-"), [hh,mm] = $("#"+pre+"hora").value.split(":");
  return `anio=${+a}&mes=${+m}&dia=${+d}&hh=${+hh}&mm=${+mm}&tz=${$("#"+pre+"tz").value}&lat=${$("#"+pre+"lat").value}&lon=${$("#"+pre+"lon").value}`;
}

async function levantar(e) {
  if (e) e.preventDefault();
  const c = await (await fetch("/api/carta?" + params())).json();
  const [Y,M,D] = $("#fecha").value.split("-"), [H,Mi] = $("#hora").value.split(":");
  const ciudad = $("#ciudad").value || "coordenadas manuales";
  const tz = +$("#tz").value;
  $("#ficha").innerHTML = `<h2>${ciudad}</h2><p>${D}/${M}/${Y} · ${H}:${Mi} ·
    UTC${tz>=0?"+":""}${tz} · lat ${$("#lat").value} · lon ${$("#lon").value} ·
    ${c.ut} · ${c.placidusOK ? "casas de Plácido" : "casas iguales (latitud extrema)"}</p>`;
  $("#rueda").innerHTML = rueda(c);
  $("#tablas").innerHTML = tablas(c);
  if (!c.placidusOK) $("#tablas").insertAdjacentHTML("afterbegin",
    `<div class="aviso"><b>Latitud extrema:</b> Plácido no puede resolverse aquí, se muestran casas iguales.</div>`);
}
async function leer() {
  const L = await (await fetch("/api/lectura?" + params())).json();
  const cats = {};
  L.frases.forEach(f => (cats[f.categoria] = cats[f.categoria] || []).push(f));
  let h = `<div class="caja"><h3>Planeta dominante</h3><p style="margin:0">${L.dominante}</p></div>`;
  for (const c in cats) {
    h += `<div class="caja"><h3>${c}</h3>`;
    cats[c].sort((a,b) => b.peso - a.peso).forEach(f =>
      h += `<p style="margin:0 0 10px"> ${f.texto}<br><span style="font-family:ui-monospace,Menlo,monospace;font-size:.72rem;color:var(--muted)">← ${f.fuente}</span></p>`);
    h += `</div>`;
  }
  if (L.contradicciones.length) {
    h += `<div class="caja"><h3>Contradicciones — aquí empieza el oficio</h3>`;
    L.contradicciones.forEach(t => h += `<p style="margin:0 0 10px">${t}</p>`);
    h += `</div>`;
  }
  h += `<div class="aviso">${L.nota}</div>`;
  $("#lec").innerHTML = h;
}
// ── ciudades y huso ──
async function buscarCiudad(pre) {
  const q = $("#" + pre + "ciudad").value;
  const cont = $("#" + pre + "sug");
  if (q.length < 2) { cont.innerHTML = ""; return; }
  const ls = await (await fetch("/api/lugares?q=" + encodeURIComponent(q))).json();
  cont.innerHTML = (ls || []).map((l,i) =>
    `<button type="button" data-i="${i}">${l.nombre}<span>${l.region?l.region+" · ":""}${l.pais}</span></button>`).join("");
  cont.dataset.datos = JSON.stringify(ls || []);
}
async function elegirCiudad(pre, l) {
  $("#" + pre + "lat").value = l.lat; $("#" + pre + "lon").value = l.lon;
  $("#" + pre + "ciudad").value = l.nombre; $("#" + pre + "sug").innerHTML = "";
  $("#" + pre + "zona").value = l.zona;
  await resolverHuso(pre);
}
async function resolverHuso(pre) {
  const zona = $("#" + pre + "zona").value; if (!zona) return;
  const [a,m,d] = $("#" + pre + "fecha").value.split("-");
  const [hh,mm] = $("#" + pre + "hora").value.split(":");
  const r = await (await fetch(`/api/huso?zona=${encodeURIComponent(zona)}&anio=${+a}&mes=${+m}&dia=${+d}&hh=${+hh}&mm=${+mm}`)).json();
  if (r.error) return;
  $("#" + pre + "tz").value = r.offset;
  $("#" + pre + "huso").innerHTML = `UTC${r.offset>=0?"+":""}${r.offset} · ${r.zona}` +
    (r.verano ? ` · <b style="color:var(--horizonte)">horario de verano activo</b>` : ` · horario estándar`);
}
["", "e"].forEach(pre => {
  const ci = $("#" + pre + "ciudad"); if (!ci) return;
  ci.oninput = () => buscarCiudad(pre);
  $("#" + pre + "sug").onclick = ev => {
    const b = ev.target.closest("button"); if (!b) return;
    elegirCiudad(pre, JSON.parse($("#" + pre + "sug").dataset.datos)[+b.dataset.i]);
  };
  $("#" + pre + "fecha").onchange = () => resolverHuso(pre);
  $("#" + pre + "hora").onchange = () => resolverHuso(pre);
});

$("#imprimir").onclick = () => window.print();

// ── cartas guardadas ──
async function pintarGuardadas(d) {
  d = d || await (await fetch("/api/guardadas")).json();
  const c = d.cartas || [];
  $("#guardadas").innerHTML = c.length
    ? `<span class="et">Guardadas</span>` + c.map(x =>
        `<span class="chip"><button class="abrir" data-id="${x.id}">${x.nombre}<small>${x.ciudad}</small></button>
         <button class="x" data-del="${x.id}" title="Borrar">×</button></span>`).join("")
    : "";
  $("#guardadas").dataset.datos = JSON.stringify(c);
}
$("#guardar").onclick = async () => {
  const sug = $("#ciudad").value || "Carta";
  const nombre = prompt("¿Con qué nombre la guardo?", sug);
  if (!nombre) return;
  const d = await (await fetch("/api/guardadas", { method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ nombre, ciudad: $("#ciudad").value, zona: $("#zona").value,
      fecha: $("#fecha").value, hora: $("#hora").value, tz: +$("#tz").value,
      lat: +$("#lat").value, lon: +$("#lon").value }) })).json();
  pintarGuardadas(d);
};
$("#guardadas").onclick = async ev => {
  const del = ev.target.closest("[data-del]");
  if (del) {
    const c = JSON.parse($("#guardadas").dataset.datos).find(x => x.id === del.dataset.del);
    if (!confirm(`¿Borrar «${c.nombre}»?`)) return;
    pintarGuardadas(await (await fetch("/api/guardadas?id=" + del.dataset.del, { method: "DELETE" })).json());
    return;
  }
  const ab = ev.target.closest("[data-id]");
  if (!ab) return;
  const c = JSON.parse($("#guardadas").dataset.datos).find(x => x.id === ab.dataset.id);
  $("#ciudad").value = c.ciudad; $("#zona").value = c.zona; $("#fecha").value = c.fecha;
  $("#hora").value = c.hora; $("#tz").value = c.tz; $("#lat").value = c.lat; $("#lon").value = c.lon;
  await resolverHuso(""); levantar(); leer();
};
pintarGuardadas();

$("#porque").onclick = async () => {
  const caja = $("#hist");
  if (!caja.hidden) { caja.hidden = true; return; }
  const [a,m,d] = $("#fecha").value.split("-"), [hh,mm] = $("#hora").value.split(":");
  const h = await (await fetch(`/api/husohistoria?zona=${encodeURIComponent($("#zona").value)}&anio=${+a}&mes=${+m}&dia=${+d}&hh=${+hh}&mm=${+mm}&lon=${$("#lon").value}`)).json();
  if (h.error) { caja.innerHTML = h.error; caja.hidden = false; return; }
  const sg = n => (n >= 0 ? "+" : "") + (Math.round(n*100)/100);
  let t = `<h4>De dónde sale UTC${sg(h.offset)}</h4>`;
  t += `<div class="fila">zona horaria · <b>${h.zona}</b> (${h.abrev})</div>`;
  t += `<div class="fila">desfase estándar de ese año · <b>UTC${sg(h.estandar)}</b></div>`;
  t += `<div class="fila">horario de verano ese día · <b>${h.verano ? "SÍ, +1 h" : "no"}</b></div>`;

  t += `<div class="sec"><h4>El reloj frente al Sol</h4>`;
  t += `<div class="fila">por longitud le correspondería · UTC${sg(h.solar)}</div>`;
  t += `<div class="fila">el reloj va <b>${sg(h.desfase)} h</b> respecto al Sol de ese lugar</div>`;
  t += `<p style="margin:6px 0 0">Por eso el desfase no se puede deducir de la longitud: es una
    decisión administrativa, no astronómica.</p></div>`;

  t += `<div class="sec"><h4>Horario de verano en ${a}</h4>`;
  t += h.delAnio.length
    ? h.delAnio.map(c => `<div class="fila">${c.fecha} · UTC${sg(c.de)} → <b>UTC${sg(c.a)}</b> — ${c.motivo}</div>`).join("")
    : `<div class="fila">ese año no hubo cambios de hora en este país</div>`;
  t += `</div>`;

  if (h.historicos.length) {
    t += `<div class="sec"><h4>Cambios de huso del país</h4>`;
    t += h.historicos.map(c => `<div class="fila">${c.fecha} · UTC${sg(c.de)} → <b>UTC${sg(c.a)}</b></div>`).join("");
    t += `</div>`;
  }
  caja.innerHTML = t; caja.hidden = false;
};

$("#f").onsubmit = e => { levantar(e); leer(); };
levantar(); leer();

// ── curso ──
const MODS = [["01-el-cielo","El cielo desde un punto"],["02-la-hora","La hora"],
["03-el-calculo","El cálculo"],["04-angulos","Los cuatro ángulos"],["05-posiciones","Posiciones"],
["06-casas","Casas y regentes"],["07-aspectos","Aspectos"],["08-dignidades","Dignidades"],
["09-combinar","Combinar y contrarrestar"],["10-sintesis","Síntesis"],["11-profundizar","Profundizar"],
["12-prediccion","Predicción"],["13-oficio","Oficio y límites"],["14-jyotisha","El zodíaco sidéreo"]];
$("#lista").innerHTML = MODS.map(([f,t],i) =>
  `<a href="#" data-f="${f}"><b>Módulo ${i+1}</b>${t}</a>`).join("") +
  `<a href="#" data-f="corta"><b>Extra</b>La vía corta en cinco bloques</a>` +
  `<a href="#" data-f="motor-profundizacion"><b>Extra</b>Motor de profundización</a>` +
  `<a href="#" data-f="plan-semana"><b>Extra</b>Plan de una semana</a>`;
$("#lista").onclick = async ev => {
  const a = ev.target.closest("a"); if (!a) return; ev.preventDefault();
  const md = await (await fetch("curso/" + a.dataset.f + ".md")).text();
  $("#texto").hidden = false; $("#texto").innerHTML = markdown(md);
  $("#texto").scrollIntoView({ behavior: "smooth", block: "start" });
};

function markdown(t) {
  const esc = s => s.replace(/&/g,"&amp;").replace(/</g,"&lt;");
  const inline = s => esc(s).replace(/`([^`]+)`/g,"<code>$1</code>")
    .replace(/\*\*([^*]+)\*\*/g,"<strong>$1</strong>").replace(/\*([^*]+)\*/g,"<em>$1</em>");
  const out = []; let tabla = null;
  for (const ln of t.split("\n")) {
    const fila = ln.trim().startsWith("|") && ln.trim().endsWith("|");
    if (fila) {
      const cel = ln.trim().slice(1,-1).split("|").map(x => x.trim());
      if (cel.every(c => /^:?-+:?$/.test(c))) continue;
      if (!tabla) { tabla = true; out.push("<table><tbody>"); }
      out.push("<tr>" + cel.map(c => `<td>${inline(c)}</td>`).join("") + "</tr>");
      continue;
    }
    if (tabla) { out.push("</tbody></table>"); tabla = null; }
    if (/^#{1,4} /.test(ln)) { const n = ln.match(/^#+/)[0].length;
      out.push(`<h${n}>${inline(ln.slice(n+1))}</h${n}>`); }
    else if (/^> /.test(ln)) out.push(`<blockquote>${inline(ln.slice(2))}</blockquote>`);
    else if (/^[-*] /.test(ln)) out.push(`<li>${inline(ln.slice(2))}</li>`);
    else if (/^\d+\. /.test(ln)) out.push(`<li>${inline(ln.replace(/^\d+\. /,""))}</li>`);
    else if (ln.trim() === "---") out.push("<hr>");
    else if (ln.trim()) out.push(`<p>${inline(ln)}</p>`);
  }
  if (tabla) out.push("</tbody></table>");
  return out.join("\n").replace(/(<li>[\s\S]*?<\/li>)(?!\s*<li>)/g, "<ul>$1</ul>");
}

// ── ejercicio ──
$("#fv").onsubmit = async e => {
  e.preventDefault();
  const extra = ["jd","tsg","tsl","asc","mc"].map(k => `${k}=${$("#"+k).value||0}`).join("&");
  const r = await (await fetch("/api/verificar?" + params("e") + "&" + extra)).json();
  let h = "";
  r.pasos.forEach((p,i) => {
    const marca = p.bien ? "✓" : "✗";
    h += `<div class="paso"><span>${p.nombre}</span>
      <span class="${p.bien?'bien':'mal'}">${marca}</span>
      <span style="color:var(--muted)">${p.tuyo}</span>
      <span class="${p.bien?'bien':'mal'}">${p.bien?'':'±'+p.desvio+' '+p.unidad}</span></div>`;
  });
  if (r.primerFallo >= 0) {
    const p = r.pasos[r.primerFallo];
    h += `<div class="aviso"><b>Primer paso que falla: ${p.nombre}.</b> Te desvías
      ${p.desvio} ${p.unidad} — ${p.comentario}. Rehaz desde ahí; lo de abajo arrastra este error.</div>`;
  } else {
    h += `<div class="aviso" style="border-color:var(--ok)"><b>Los cinco pasos correctos.</b>
      Has levantado la carta a mano. Eso es lo que separa saber astrología de saber usar un programa.</div>`;
  }
  $("#res").innerHTML = h;
};
