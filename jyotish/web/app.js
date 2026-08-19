const $ = s => document.querySelector(s);
let CARTA = null, ESTILO = "norte";

document.querySelectorAll("nav button[data-s]").forEach(b => b.onclick = () => {
  document.querySelectorAll("nav button[data-s]").forEach(x => x.classList.toggle("on", x === b));
  document.querySelectorAll("section").forEach(s => s.classList.toggle("on", s.id === b.dataset.s));
});
$("#imprimir").onclick = () => window.print();
document.querySelectorAll(".estilos button").forEach(b => b.onclick = () => {
  document.querySelectorAll(".estilos button").forEach(x => x.classList.toggle("on", x === b));
  ESTILO = b.dataset.e; if (CARTA) $("#rasi").innerHTML = dibujar(CARTA.grahas, CARTA.lagnaRasi, "D1");
});

// ── Sur de la India: los signos están fijos, rota la casa ──
const SUR = [[1,0],[2,0],[3,0],[3,1],[3,2],[3,3],[2,3],[1,3],[0,3],[0,2],[0,1],[0,0]];
// ── Norte de la India: las casas están fijas, rota el signo ──
const NORTE = [[.50,.26],[.25,.11],[.11,.25],[.26,.50],[.11,.75],[.25,.89],
               [.50,.74],[.75,.89],[.89,.75],[.74,.50],[.89,.25],[.75,.11]];

function dibujar(grahas, lagnaRasi, titulo) {
  const S = 400, P = 10, L = S - 2*P;
  const porRasi = {}; grahas.forEach(g => (porRasi[g.rasiIdx] = porRasi[g.rasiIdx] || []).push(g));
  let s = `<svg viewBox="0 0 ${S} ${S+16}" role="img" aria-label="Carta ${titulo} estilo ${ESTILO}">`;
  const cel = (x,y,w,h) => `<rect x="${x}" y="${y}" width="${w}" height="${h}" fill="var(--casilla)" stroke="currentColor" stroke-width="1.1" opacity=".9"/>`;

  if (ESTILO === "sur") {
    const c = L/4;
    for (let r = 0; r < 12; r++) {
      const [cx, cy] = SUR[r], x = P + cx*c, y = P + cy*c;
      s += cel(x, y, c, c);
      s += `<text x="${x+5}" y="${y+13}" font-size="9.5" font-family="ui-monospace,Menlo,monospace" fill="currentColor" opacity=".5">${["Me","Vṛ","Mi","Ka","Si","Kn","Tu","Vś","Dh","Mk","Ku","Mn"][r]}</text>`;
      if (r === lagnaRasi) {
        s += `<line x1="${x}" y1="${y}" x2="${x+c*.34}" y2="${y+c*.34}" stroke="var(--acento)" stroke-width="2"/>`;
        s += `<text x="${x+c-5}" y="${y+13}" text-anchor="end" font-size="9" font-family="ui-monospace,Menlo,monospace" fill="var(--acento)" font-weight="700">La</text>`;
      }
      (porRasi[r]||[]).forEach((g,i) =>
        s += `<text x="${x+c/2}" y="${y+30+i*17}" text-anchor="middle" font-size="14" fill="var(--acento)">${g.glifo}<title>${g.nombre} ${g.posicion||""}</title></text>`);
    }
    s += `<text x="${S/2}" y="${S+11}" text-anchor="middle" font-size="10" font-family="ui-monospace,Menlo,monospace" fill="currentColor" opacity=".55">${titulo} · sur de la India · signos fijos</text>`;
  } else {
    const x0=P, y0=P, x1=P+L, y1=P+L, mx=(x0+x1)/2, my=(y0+y1)/2;
    s += `<rect x="${x0}" y="${y0}" width="${L}" height="${L}" fill="var(--casilla)" stroke="currentColor" stroke-width="1.4"/>`;
    s += `<line x1="${x0}" y1="${y0}" x2="${x1}" y2="${y1}" stroke="currentColor" stroke-width="1.1"/>`;
    s += `<line x1="${x1}" y1="${y0}" x2="${x0}" y2="${y1}" stroke="currentColor" stroke-width="1.1"/>`;
    s += `<path d="M ${mx},${y0} L ${x1},${my} L ${mx},${y1} L ${x0},${my} Z" fill="none" stroke="currentColor" stroke-width="1.1"/>`;
    for (let h = 0; h < 12; h++) {
      const [fx, fy] = NORTE[h], cx = x0 + fx*L, cy = y0 + fy*L;
      const rasi = (lagnaRasi + h) % 12;
      s += `<text x="${cx}" y="${cy-13}" text-anchor="middle" font-size="9.5" font-family="ui-monospace,Menlo,monospace" fill="currentColor" opacity=".5">${rasi+1}</text>`;
      (porRasi[rasi]||[]).forEach((g,i) => {
        const fila = Math.floor(i/2), col = i%2, n = (porRasi[rasi]||[]).length;
        const dx = n > 1 ? (col-0.5)*20 : 0;
        s += `<text x="${cx+dx}" y="${cy+4+fila*15}" text-anchor="middle" font-size="13.5" fill="var(--acento)">${g.glifo}<title>${g.nombre} ${g.posicion||""}</title></text>`;
      });
    }
    s += `<text x="${S/2}" y="${S+11}" text-anchor="middle" font-size="10" font-family="ui-monospace,Menlo,monospace" fill="currentColor" opacity=".55">${titulo} · norte de la India · casas fijas, casa 1 arriba</text>`;
  }
  return s + "</svg>";
}

function tablas(c) {
  let h = `<div class="caja"><h3>Lagna</h3><table><tbody>
    <tr><td>Lagna</td><td>${c.lagnaPos}</td><td>${c.lagnaNak} pada ${c.lagnaPada}</td></tr>
    <tr><td>señor del Lagna</td><td colspan="2">${c.senorLagna}</td></tr>
    <tr><td>ayanāṁśa</td><td colspan="2">${c.ayanamsa.toFixed(4)}°</td></tr></tbody></table></div>`;
  h += `<div class="caja"><h3>Grahas</h3><table><thead><tr><th colspan="2">Graha</th>
    <th>Posición</th><th>Nakṣatra</th><th class="num">Pada</th><th class="num">Casa</th><th>Estado</th></tr></thead><tbody>`;
  c.grahas.forEach(g => {
    const cl = g.dignidad === "exaltado" ? "exalt" : g.dignidad === "debilitado" ? "debil" : "";
    h += `<tr><td class="gl">${g.glifo}</td><td>${g.nombre}${g.retro?' ℞':''}</td>
      <td>${g.posicion}${g.gandanta?' <span class="gan">⚠</span>':''}</td><td>${g.nak}</td>
      <td class="num">${g.pada}</td><td class="num">${g.bhava}</td><td class="${cl}">${g.dignidad}</td></tr>`;
  });
  h += `</tbody></table></div><div class="caja"><h3>Bhāvas y sus señores</h3><table>
    <thead><tr><th class="num">Casa</th><th>Rāśi</th><th>Señor</th><th>alojado en</th><th>Ocupan</th><th>Aspectan</th></tr></thead><tbody>`;
  c.bhavas.forEach(b => h += `<tr><td class="num">${b.numero}</td><td>${b.rasi}</td><td>${b.senor}</td>
    <td>casa ${b.senorEn}</td><td>${(b.ocupan||[]).join(" ")||"—"}</td>
    <td style="color:var(--muted)">${(b.aspectan||[]).join(" ")||"—"}</td></tr>`);
  h += `</tbody></table></div><div class="caja"><h3>Kārakas de Jaimini</h3><table><tbody>`;
  const kn = {AK:"Ātmakāraka · el alma",AmK:"Amātyakāraka · mente y carrera",BK:"Bhrātṛkāraka · hermanos",
    MK:"Mātṛkāraka · madre",PiK:"Pitṛkāraka · padre",PK:"Putrakāraka · hijos",GK:"Ñātikāraka · obstáculos"};
  Object.keys(kn).forEach(k => h += `<tr><td><b>${k}</b></td><td>${c.karakas[k]||"—"}</td><td style="color:var(--muted)">${kn[k]}</td></tr>`);
  h += `</tbody></table></div>`;
  if (c.yogas && c.yogas.length) {
    h += `<div class="caja"><h3>Yogas detectados</h3>`;
    c.yogas.forEach(y => h += `<p class="yoga">${y}</p>`);
    h += `</div>`;
  }
  return h;
}

function pintarVargas(c) {
  const desc = {D1:"la vida tal como se vive",D2:"riqueza y sustento",D3:"hermanos, coraje, esfuerzo",
    D7:"hijos y descendencia",D9:"el alma, el cónyuge, la fuerza real",D10:"profesión y karma público",
    D12:"los padres",D16:"vehículos y confort",D30:"males y desgracias",D60:"karma acumulado"};
  $("#vg").innerHTML = Object.keys(c.vargas).map(k => {
    const gs = c.vargas[k], lag = gs.find(g => g.nombre === "Lagna");
    return `<div class="vg"><h4>${k}</h4><p>${desc[k]||""}</p>
      ${dibujar(gs.filter(g=>g.nombre!=="Lagna"), lag ? lag.rasiIdx : 0, k)}</div>`;
  }).join("");
}

function pintarDasas(c) {
  let h = `<div class="caja"><h3>Vimśottarī — mahādaśās y bhuktis</h3><div class="dasa">`;
  c.dasas.forEach(p => {
    h += `<div class="p ${p.actual?'act':''}"><span>${p.senor}</span>
      <span>${p.desde} → ${p.hasta}</span><span>${p.anios} años</span></div>`;
    if (p.actual && p.sub) p.sub.forEach(b =>
      h += `<div class="p b ${b.actual?'act':''}"><span>${b.senor}</span>
        <span>${b.desde} → ${b.hasta}</span><span>${b.anios} a</span></div>`);
  });
  $("#ds").innerHTML = h + `</div></div>`;
}

function params(){const [a,m,d]=$("#fecha").value.split("-"),[hh,mm]=$("#hora").value.split(":");
  return `anio=${+a}&mes=${+m}&dia=${+d}&hh=${+hh}&mm=${+mm}&tz=${$("#tz").value}&lat=${$("#lat").value}&lon=${$("#lon").value}`;}

async function levantar(e){
  if(e) e.preventDefault();
  const c = await (await fetch("/api/carta?"+params())).json();
  CARTA = c;
  const [Y,M,D]=$("#fecha").value.split("-"),[H,Mi]=$("#hora").value.split(":");
  $("#ficha").innerHTML = `<h2>${$("#ciudad").value||"coordenadas manuales"}</h2>
    <p>${D}/${M}/${Y} · ${H}:${Mi} · UTC${+$("#tz").value>=0?"+":""}${$("#tz").value} ·
    ${c.ut} · Lagna ${c.lagnaPos} · ayanāṁśa ${c.ayanamsa.toFixed(4)}°</p>`;
  $("#rasi").innerHTML = dibujar(c.grahas, c.lagnaRasi, "D1");
  $("#tablas").innerHTML = tablas(c);
  pintarVargas(c); pintarDasas(c);
}
$("#f").onsubmit = levantar;

// ── ciudades, husos, guardadas ──
async function buscarCiudad(){const q=$("#ciudad").value,c=$("#sug");
  if(q.length<2){c.innerHTML="";return;}
  const ls=await (await fetch("/api/lugares?q="+encodeURIComponent(q))).json();
  c.innerHTML=(ls||[]).map((l,i)=>`<button type="button" data-i="${i}">${l.nombre}<span>${l.region?l.region+" · ":""}${l.pais}</span></button>`).join("");
  c.dataset.datos=JSON.stringify(ls||[]);}
async function resolverHuso(){const z=$("#zona").value;if(!z)return;
  const [a,m,d]=$("#fecha").value.split("-"),[hh,mm]=$("#hora").value.split(":");
  const r=await (await fetch(`/api/huso?zona=${encodeURIComponent(z)}&anio=${+a}&mes=${+m}&dia=${+d}&hh=${+hh}&mm=${+mm}`)).json();
  if(r.error)return; $("#tz").value=r.offset;
  $("#huso").innerHTML=`UTC${r.offset>=0?"+":""}${r.offset} · ${r.zona}`+(r.verano?` · <b style="color:var(--acento)">horario de verano activo</b>`:` · horario estándar`);}
$("#ciudad").oninput=buscarCiudad;
$("#sug").onclick=async ev=>{const b=ev.target.closest("button");if(!b)return;
  const l=JSON.parse($("#sug").dataset.datos)[+b.dataset.i];
  $("#lat").value=l.lat;$("#lon").value=l.lon;$("#ciudad").value=l.nombre;$("#zona").value=l.zona;$("#sug").innerHTML="";
  await resolverHuso();};
$("#fecha").onchange=resolverHuso; $("#hora").onchange=resolverHuso;

async function pintarGuardadas(d){d=d||await (await fetch("/api/guardadas")).json();
  const c=d.cartas||[];
  $("#guardadas").innerHTML=c.length?`<span class="et">Guardadas</span>`+c.map(x=>
    `<span class="chip"><button class="abrir" data-id="${x.id}">${x.nombre}<small>${x.ciudad}</small></button>
     <button class="x" data-del="${x.id}">×</button></span>`).join(""):"";
  $("#guardadas").dataset.datos=JSON.stringify(c);}
$("#guardar").onclick=async()=>{const n=prompt("¿Con qué nombre la guardo?",$("#ciudad").value||"Carta");if(!n)return;
  pintarGuardadas(await (await fetch("/api/guardadas",{method:"POST",headers:{"Content-Type":"application/json"},
    body:JSON.stringify({nombre:n,ciudad:$("#ciudad").value,zona:$("#zona").value,fecha:$("#fecha").value,
      hora:$("#hora").value,tz:+$("#tz").value,lat:+$("#lat").value,lon:+$("#lon").value})})).json());};
$("#guardadas").onclick=async ev=>{
  const del=ev.target.closest("[data-del]");
  if(del){const c=JSON.parse($("#guardadas").dataset.datos).find(x=>x.id===del.dataset.del);
    if(!confirm(`¿Borrar «${c.nombre}»?`))return;
    pintarGuardadas(await (await fetch("/api/guardadas?id="+del.dataset.del,{method:"DELETE"})).json());return;}
  const ab=ev.target.closest("[data-id]");if(!ab)return;
  const c=JSON.parse($("#guardadas").dataset.datos).find(x=>x.id===ab.dataset.id);
  $("#ciudad").value=c.ciudad;$("#zona").value=c.zona;$("#fecha").value=c.fecha;$("#hora").value=c.hora;
  $("#tz").value=c.tz;$("#lat").value=c.lat;$("#lon").value=c.lon;await resolverHuso();levantar();};
// ── curso ──
const MODS=[["00-mapa","Mapa del curso",""],["01-el-cielo","El cielo desde un punto","1"],
["02-la-hora","La hora","2"],["03-el-calculo","El cálculo","3"],["04-angulos","Los ángulos","4"],
["05-rasis","Los rāśis","5"],["06-nakshatras","Los nakṣatras","6"],["07-bhavas","Los bhāvas","7"],
["08-grahas","Los grahas","8"],["09-drishti","Dṛṣṭi: la mirada","9"],["10-dignidades","Dignidades y fuerza","10"],
["11-vargas","Las vargas","11"],["12-dasas","Las daśās","12"],["13-oficio","Oficio, ética y límites","13"],
["14-karakas","Los kārakas","14"],["15-yogas","Los yogas","15"],["16-profundizar","Profundizar","16"]];
$("#lista").innerHTML=MODS.map(([f,t,n])=>`<a href="#" data-f="${f}"><b>${n?"Módulo "+n:"Índice"}</b>${t}</a>`).join("");
$("#lista").onclick=async ev=>{const a=ev.target.closest("a");if(!a)return;ev.preventDefault();
  const md=await (await fetch("curso/"+a.dataset.f+".md")).text();
  $("#texto").hidden=false;$("#texto").innerHTML=markdown(md);
  $("#texto").scrollIntoView({behavior:"smooth",block:"start"});};
function markdown(t){const esc=s=>s.replace(/&/g,"&amp;").replace(/</g,"&lt;");
  const inl=s=>esc(s).replace(/`([^`]+)`/g,"<code>$1</code>").replace(/\*\*([^*]+)\*\*/g,"<strong>$1</strong>").replace(/\*([^*]+)\*/g,"<em>$1</em>");
  const o=[];let tb=null;
  for(const ln of t.split("\n")){
    const fila=ln.trim().startsWith("|")&&ln.trim().endsWith("|");
    if(fila){const cl=ln.trim().slice(1,-1).split("|").map(x=>x.trim());
      if(cl.every(c=>/^:?-+:?$/.test(c)))continue;
      if(!tb){tb=1;o.push("<table><tbody>");}
      o.push("<tr>"+cl.map(c=>`<td>${inl(c)}</td>`).join("")+"</tr>");continue;}
    if(tb){o.push("</tbody></table>");tb=null;}
    if(/^#{1,4} /.test(ln)){const n=ln.match(/^#+/)[0].length;o.push(`<h${n}>${inl(ln.slice(n+1))}</h${n}>`);}
    else if(/^> /.test(ln))o.push(`<blockquote>${inl(ln.slice(2))}</blockquote>`);
    else if(/^[-*] /.test(ln))o.push(`<li>${inl(ln.slice(2))}</li>`);
    else if(ln.trim()==="---")o.push("<hr>");
    else if(ln.trim())o.push(`<p>${inl(ln)}</p>`);}
  if(tb)o.push("</tbody></table>");
  return o.join("\n").replace(/(<li>[\s\S]*?<\/li>)(?!\s*<li>)/g,"<ul>$1</ul>");}

pintarGuardadas(); levantar();
