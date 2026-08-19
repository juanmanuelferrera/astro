const $ = s => document.querySelector(s);
let TRAD = "occidental", LANG = "es", ESTILO = "norte", DATOS = null;
const t = () => T[LANG];

const SIGW = ["♈","♉","♊","♋","♌","♍","♎","♏","♐","♑","♒","♓"];
const NOMW = ["Aries","Tauro","Géminis","Cáncer","Leo","Virgo","Libra","Escorpio","Sagitario","Capricornio","Acuario","Piscis"];
const SUR = [[1,0],[2,0],[3,0],[3,1],[3,2],[3,3],[2,3],[1,3],[0,3],[0,2],[0,1],[0,0]];
const NORTE = [[.50,.26],[.25,.11],[.11,.25],[.26,.50],[.11,.75],[.25,.89],[.50,.74],[.75,.89],[.89,.75],[.74,.50],[.89,.25],[.75,.11]];
const gms = g => { const d=Math.floor(g), m=Math.round((g-d)*60); return m===60?`${d+1}° 00′`:`${d}° ${String(m).padStart(2,"0")}′`; };

// ═══ selector de tradición e idioma ═══
document.querySelectorAll(".selector button").forEach(b => b.onclick = () => {
  document.querySelectorAll(".selector button").forEach(x => x.classList.toggle("on", x===b));
  TRAD = b.dataset.t; document.documentElement.dataset.trad = TRAD;
  pintarNav(); pintarCurso(); $("#estilosBox").hidden = TRAD !== "jyotisha";
  aplicarIdioma(); levantar();
});
document.querySelectorAll(".idioma button").forEach(b => b.onclick = () => {
  document.querySelectorAll(".idioma button").forEach(x => x.classList.toggle("on", x===b));
  LANG = b.dataset.i; document.documentElement.lang = LANG; aplicarIdioma(); levantar();
});
document.querySelectorAll(".estilos button").forEach(b => b.onclick = () => {
  document.querySelectorAll(".estilos button").forEach(x => x.classList.toggle("on", x===b));
  ESTILO = b.dataset.e; if (DATOS) render();
});

function pintarNav() {
  const n = $("#nav");
  n.innerHTML = NAV[TRAD].map((k,i) =>
    `<button data-s="${k}" class="${i===0?"on":""}">${t().nav[k]}</button>`).join("") +
    `<button class="imp" id="imprimir">${t().nav.imprimir}</button>`;
  n.querySelectorAll("button[data-s]").forEach(b => b.onclick = () => {
    n.querySelectorAll("button[data-s]").forEach(x => x.classList.toggle("on", x===b));
    document.querySelectorAll("section").forEach(s => s.classList.toggle("on", s.id===b.dataset.s));
    if (b.dataset.s === "comparar") comparar();
    if (b.dataset.s === "ejercicio") pintarEjercicio();
  });
  $("#imprimir").onclick = () => window.print();
  document.querySelectorAll("section").forEach(s => s.classList.toggle("on", s.id === NAV[TRAD][0]));
}

function aplicarIdioma() {
  const x = t();
  $("#titulo").textContent = x.titulo;
  $("#lead").textContent = TRAD === "jyotisha" ? x.lead_jyotisha : x.lead_occidental;
  const et = {lCiudad:"ciudad",lFecha:"fecha",lHora:"hora",lTz:"tz",lLat:"lat",lLon:"lon"};
  for (const id in et) { const l=$("#"+id), inp=l.querySelector("input,div");
    l.childNodes[0].nodeValue = x[et[id]]; }
  $("#bLevantar").textContent = x.levantar; $("#guardar").textContent = x.guardar;
  $("#porque").textContent = x.porque; $("#pie").textContent = x.pie;
  if (!$("#husoTxt").textContent) $("#husoTxt").textContent = x.elige;
  pintarNav(); pintarCurso();
}

// ═══ dibujo ═══
function ruedaOcc(c) {
  const CX=250,CY=250,R=196,RS=152,RH=112,asc=c.angulos[0];
  const P=(l,r)=>{const a=(180-(l-asc))*Math.PI/180;return [CX+r*Math.cos(a),CY-r*Math.sin(a)];};
  let s=`<svg viewBox="0 0 500 500" role="img" aria-label="Rueda natal">`;
  for(let i=0;i<12;i++){const [x0,y0]=P(i*30,RS),[x1,y1]=P(i*30,R);
    s+=`<line x1="${x0.toFixed(1)}" y1="${y0.toFixed(1)}" x2="${x1.toFixed(1)}" y2="${y1.toFixed(1)}" stroke="currentColor" stroke-width="1" opacity=".3"/>`;
    const [gx,gy]=P(i*30+15,(R+RS)/2);
    s+=`<text x="${gx.toFixed(1)}" y="${(gy+5).toFixed(1)}" text-anchor="middle" font-size="17" fill="currentColor" opacity=".8">${SIGW[i]}</text>`;}
  [R,RS,RH].forEach((r,i)=>s+=`<circle cx="${CX}" cy="${CY}" r="${r}" fill="none" stroke="currentColor" stroke-width="${i?1:1.4}" opacity="${i?.3:.45}"/>`);
  c.cuspP.forEach((cu,i)=>{const [x0,y0]=P(cu,RH),[x1,y1]=P(cu,RS),ang=i%3===0;
    s+=`<line x1="${x0.toFixed(1)}" y1="${y0.toFixed(1)}" x2="${x1.toFixed(1)}" y2="${y1.toFixed(1)}" stroke="${ang?"var(--ac)":"currentColor"}" stroke-width="${ang?2.2:1}" opacity="${ang?.95:.28}"/>`;
    let d=c.cuspP[(i+1)%12]-cu; if(d<0)d+=360;
    const [nx,ny]=P(cu+d/2,RH+20);
    s+=`<text x="${nx.toFixed(1)}" y="${(ny+4).toFixed(1)}" text-anchor="middle" font-size="12" fill="currentColor" opacity=".55">${i+1}</text>`;});
  const RG=RS-26,MIN=8.5;
  const cu=c.cuerpos.map(p=>({...p,rel:((p.lon-asc)%360+360)%360})).sort((a,b)=>a.rel-b.rel);
  cu.forEach(p=>p.dib=p.rel);
  for(let k=0;k<250;k++){let mv=false;
    for(let i=0;i<cu.length;i++){const a=cu[i],b=cu[(i+1)%cu.length];let d=b.dib-a.dib;if(d<0)d+=360;
      if(d<MIN){const e=(MIN-d)/2;a.dib-=e;b.dib+=e;mv=true;}}
    if(!mv)break;}
  cu.forEach(p=>{const [tx,ty]=P(asc+p.rel,RS),[t2x,t2y]=P(asc+p.rel,RS-7);
    s+=`<line x1="${tx.toFixed(1)}" y1="${ty.toFixed(1)}" x2="${t2x.toFixed(1)}" y2="${t2y.toFixed(1)}" stroke="var(--ac)" stroke-width="1.4" opacity=".85"/>`;
    if(Math.abs(p.dib-p.rel)>.4){const [a1,b1]=P(asc+p.rel,RS-7),[a2,b2]=P(asc+p.dib,RG+9);
      s+=`<line x1="${a1.toFixed(1)}" y1="${b1.toFixed(1)}" x2="${a2.toFixed(1)}" y2="${b2.toFixed(1)}" stroke="currentColor" stroke-width=".8" opacity=".3"/>`;}
    const [x,y]=P(asc+p.dib,RG);
    s+=`<text x="${x.toFixed(1)}" y="${(y+6).toFixed(1)}" text-anchor="middle" font-size="17" fill="var(--ac)">${p.glifo}<title>${p.nombre}</title></text>`;
    const [gx,gy]=P(asc+p.dib,RG-17);
    s+=`<text x="${gx.toFixed(1)}" y="${(gy+4).toFixed(1)}" text-anchor="middle" font-size="9" font-family="ui-monospace,Menlo,monospace" fill="currentColor" opacity=".55">${Math.floor(p.grado)}°${p.retro?"℞":""}</text>`;});
  [["Asc",c.angulos[0]],["Dsc",c.angulos[1]],["MC",c.angulos[2]],["IC",c.angulos[3]]].forEach(([n,l])=>{
    const [x,y]=P(l,R+16);
    s+=`<text x="${x.toFixed(1)}" y="${(y+4).toFixed(1)}" text-anchor="middle" font-size="10.5" font-family="ui-monospace,Menlo,monospace" font-weight="700" fill="${n[0]==="A"||n[0]==="D"?"var(--ac)":"var(--fr)"}">${n}</text>`;});
  return s+"</svg>";
}

function cuadroVed(grahas, lagnaRasi, titulo) {
  const S=400,P0=10,L=S-2*P0,por={};
  grahas.forEach(g=>(por[g.rasiIdx]=por[g.rasiIdx]||[]).push(g));
  let s=`<svg viewBox="0 0 ${S} ${S+16}" role="img" aria-label="${titulo}">`;
  if(ESTILO==="sur"){const c=L/4;
    for(let r=0;r<12;r++){const [cx,cy]=SUR[r],x=P0+cx*c,y=P0+cy*c;
      s+=`<rect x="${x}" y="${y}" width="${c}" height="${c}" fill="var(--cas)" stroke="currentColor" stroke-width="1.1" opacity=".95"/>`;
      s+=`<text x="${x+5}" y="${y+13}" font-size="9.5" font-family="ui-monospace,Menlo,monospace" fill="currentColor" opacity=".5">${["Me","Vṛ","Mi","Ka","Si","Kn","Tu","Vś","Dh","Mk","Ku","Mn"][r]}</text>`;
      if(r===lagnaRasi){s+=`<line x1="${x}" y1="${y}" x2="${x+c*.34}" y2="${y+c*.34}" stroke="var(--ac)" stroke-width="2"/>`;
        s+=`<text x="${x+c-5}" y="${y+13}" text-anchor="end" font-size="9" font-family="ui-monospace,Menlo,monospace" fill="var(--ac)" font-weight="700">La</text>`;}
      (por[r]||[]).forEach((g,i)=>s+=`<text x="${x+c/2}" y="${y+30+i*17}" text-anchor="middle" font-size="14" fill="var(--ac)">${g.glifo}<title>${g.nombre} ${g.posicion||""}</title></text>`);}
    s+=`<text x="${S/2}" y="${S+11}" text-anchor="middle" font-size="10" font-family="ui-monospace,Menlo,monospace" fill="currentColor" opacity=".55">${titulo} · sur · signos fijos</text>`;
  }else{const x0=P0,y0=P0,x1=P0+L,y1=P0+L,mx=(x0+x1)/2,my=(y0+y1)/2;
    s+=`<rect x="${x0}" y="${y0}" width="${L}" height="${L}" fill="var(--cas)" stroke="currentColor" stroke-width="1.4"/>`;
    s+=`<line x1="${x0}" y1="${y0}" x2="${x1}" y2="${y1}" stroke="currentColor" stroke-width="1.1"/><line x1="${x1}" y1="${y0}" x2="${x0}" y2="${y1}" stroke="currentColor" stroke-width="1.1"/>`;
    s+=`<path d="M ${mx},${y0} L ${x1},${my} L ${mx},${y1} L ${x0},${my} Z" fill="none" stroke="currentColor" stroke-width="1.1"/>`;
    for(let h=0;h<12;h++){const [fx,fy]=NORTE[h],cx=x0+fx*L,cy=y0+fy*L,r=(lagnaRasi+h)%12;
      s+=`<text x="${cx}" y="${cy-13}" text-anchor="middle" font-size="9.5" font-family="ui-monospace,Menlo,monospace" fill="currentColor" opacity=".5">${r+1}</text>`;
      (por[r]||[]).forEach((g,i)=>{const n=(por[r]||[]).length,dx=n>1?((i%2)-.5)*20:0;
        s+=`<text x="${cx+dx}" y="${cy+4+Math.floor(i/2)*15}" text-anchor="middle" font-size="13.5" fill="var(--ac)">${g.glifo}<title>${g.nombre} ${g.posicion||""}</title></text>`;});}
    s+=`<text x="${S/2}" y="${S+11}" text-anchor="middle" font-size="10" font-family="ui-monospace,Menlo,monospace" fill="currentColor" opacity=".55">${titulo} · norte · casa 1 arriba</text>`;}
  return s+"</svg>";
}

// ═══ tablas ═══
function tablasOcc(c){const x=t();
  const ang=["Ascendente","Descendente","Medio Cielo","Fondo del Cielo"];
  let h=`<div class="caja"><h3>${x.angulos}</h3><table><tbody>`;
  c.angulos.forEach((l,i)=>{const s=Math.floor(l/30);
    h+=`<tr><td>${ang[i]}</td><td class="num">${gms(l-s*30)}</td><td class="gl">${SIGW[s]}</td><td>${NOMW[s]}</td></tr>`;});
  h+=`</tbody></table></div><div class="caja"><h3>${x.planetas}</h3><table><thead><tr><th colspan="2">${x.planetas}</th><th class="num">${x.grados}</th><th colspan="2">${x.signo}</th><th class="num">${x.casa}</th></tr></thead><tbody>`;
  c.cuerpos.forEach(p=>h+=`<tr><td class="gl">${p.glifo}</td><td>${p.nombre}${p.retro?' ℞':''}</td><td class="num">${gms(p.grado)}</td><td class="gl">${p.glifoSig}</td><td>${p.signo}</td><td class="num">${p.casaP}</td></tr>`);
  h+=`</tbody></table></div><div class="caja"><h3>${x.aspectos}</h3><table><tbody>`;
  c.aspectos.slice(0,14).forEach(a=>h+=`<tr><td>${a.a}</td><td class="gl">${a.glifo}</td><td>${a.nombre}</td><td>${a.b}</td><td class="num">${a.orbe.toFixed(2)}°</td></tr>`);
  h+=`</tbody></table></div><div class="caja"><h3>${x.casas}</h3><table><thead><tr><th>${x.casa}</th><th>${x.signo}</th><th>${x.senor}</th><th>${x.alojado}</th></tr></thead><tbody>`;
  c.regentes.forEach((r,i)=>{const s=Math.floor(c.cuspP[i]/30);
    h+=`<tr><td class="num">${i+1}</td><td class="gl">${SIGW[s]}</td><td>${r}</td><td>${x.casa} ${c.regenteEn[i]}</td></tr>`;});
  return h+`</tbody></table></div>`;}

function tablasVed(c){const x=t();
  let h=`<div class="caja"><h3>${x.lagnaT}</h3><table><tbody>
    <tr><td>${x.lagnaT}</td><td>${c.lagnaPos}</td><td>${c.lagnaNak} · pada ${c.lagnaPada}</td></tr>
    <tr><td>${x.senor}</td><td colspan="2">${c.senorLagna}</td></tr>
    <tr><td>ayanāṁśa</td><td colspan="2">${c.ayanamsa.toFixed(4)}°</td></tr></tbody></table></div>`;
  h+=`<div class="caja"><h3>${x.grahas}</h3><table><thead><tr><th colspan="2">${x.grahas}</th><th>${x.posicion}</th><th>${x.nak}</th><th class="num">${x.pada}</th><th class="num">${x.casa}</th><th>${x.estado}</th></tr></thead><tbody>`;
  c.grahas.forEach(g=>{const cl=g.dignidad==="exaltado"?"exalt":g.dignidad==="debilitado"?"debil":"";
    const marcas=[g.combusto?`<span class="debil" title="combusto: a ${g.delSol}° del Sol">☌sol</span>`:"",
      g.digBala?`<span class="exalt" title="fuerza direccional plena">dig</span>`:""].filter(Boolean).join(" ");
    h+=`<tr><td class="gl">${g.glifo}</td><td>${g.nombre}${g.retro?' ℞':''}</td><td>${g.posicion}${g.gandanta?' <span class="gan" title="gaṇḍānta">⚠</span>':''}</td><td>${g.nak}</td><td class="num">${g.pada}</td><td class="num">${g.bhava}</td><td class="${cl}">${g.dignidad} ${marcas}</td></tr>`;});
  h+=`</tbody></table></div><div class="caja"><h3>${x.bhavas}</h3><table><thead><tr><th class="num">${x.casa}</th><th>Rāśi</th><th>${x.senor}</th><th>${x.alojado}</th><th>${x.ocupan}</th><th>${x.aspectan}</th></tr></thead><tbody>`;
  c.bhavas.forEach(b=>h+=`<tr><td class="num">${b.numero}</td><td>${b.rasi}</td><td>${b.senor}</td><td>${x.casa} ${b.senorEn}</td><td>${(b.ocupan||[]).join(" ")||"—"}</td><td style="color:var(--muted)">${(b.aspectan||[]).join(" ")||"—"}</td></tr>`);
  h+=`</tbody></table></div><div class="caja"><h3>${x.karakas}</h3>
    ${c.karakamsa?`<p style="margin:0 0 8px;color:var(--muted);font-size:.85rem">Karakāṁśa — el Ātmakāraka cae en <b style="color:var(--ink)">${c.karakamsa}</b> en el navāṁśa. Se usa como tercer ascendente para el dharma.</p>`:""}
    <table><tbody>`;
  const kn={AK:"Ātmakāraka",AmK:"Amātyakāraka",BK:"Bhrātṛkāraka",MK:"Mātṛkāraka",PiK:"Pitṛkāraka",PK:"Putrakāraka",GK:"Ñātikāraka"};
  Object.keys(kn).forEach(k=>h+=`<tr><td><b>${k}</b></td><td>${c.karakas[k]||"—"}</td><td style="color:var(--muted)">${kn[k]}</td></tr>`);
  h+=`</tbody></table></div>`;
  if(c.yogas&&c.yogas.length){h+=`<div class="caja"><h3>${x.yogas}</h3>`;c.yogas.forEach(y=>h+=`<p class="yoga">${y}</p>`);h+=`</div>`;}
  return h;}

// ═══ render ═══
function render(){
  if(TRAD==="occidental"){$("#dibujo").innerHTML=ruedaOcc(DATOS);$("#tablas").innerHTML=tablasOcc(DATOS);}
  else{$("#dibujo").innerHTML=cuadroVed(DATOS.grahas,DATOS.lagnaRasi,"D1");$("#tablas").innerHTML=tablasVed(DATOS);
    pintarVargas(DATOS);pintarDasas(DATOS);}
}
function params(){const [a,m,d]=$("#fecha").value.split("-"),[hh,mm]=$("#hora").value.split(":");
  return `anio=${+a}&mes=${+m}&dia=${+d}&hh=${+hh}&mm=${+mm}&tz=${$("#tz").value}&lat=${$("#lat").value}&lon=${$("#lon").value}`;}
async function levantar(e){if(e)e.preventDefault();
  DATOS=await (await fetch(`/api/${TRAD==="jyotisha"?"vedica":"carta"}?`+params())).json();
  const [Y,M,D]=$("#fecha").value.split("-"),[H,Mi]=$("#hora").value.split(":");
  const extra=TRAD==="jyotisha"?` · ${t().lagnaT} ${DATOS.lagnaPos}`:"";
  $("#ficha").innerHTML=`<h2>${$("#ciudad").value||"—"}</h2><p>${D}/${M}/${Y} · ${H}:${Mi} · UTC${+$("#tz").value>=0?"+":""}${$("#tz").value} · ${DATOS.ut}${extra}</p>`;
  render(); if(TRAD==="occidental")leer();}
$("#f").onsubmit=levantar;

async function leer(){const L=await (await fetch("/api/lectura?"+params())).json();
  const cats={};L.frases.forEach(f=>(cats[f.categoria]=cats[f.categoria]||[]).push(f));
  let h=`<div class="caja"><h3>${t().dominante}</h3><p style="margin:0">${L.dominante}</p></div>`;
  for(const c in cats){h+=`<div class="caja"><h3>${c}</h3>`;
    cats[c].sort((a,b)=>b.peso-a.peso).forEach(f=>h+=`<p style="margin:0 0 9px">${f.texto}<br><span style="font-family:ui-monospace,Menlo,monospace;font-size:.71rem;color:var(--muted)">← ${f.fuente}</span></p>`);h+=`</div>`;}
  if(L.contradicciones.length){h+=`<div class="caja"><h3>${t().contradicciones}</h3>`;
    L.contradicciones.forEach(x=>h+=`<p style="margin:0 0 9px">${x}</p>`);h+=`</div>`;}
  $("#lec").innerHTML=h+`<div class="aviso">${L.nota}</div>`;}

function pintarVargas(c){const d={D1:"la vida tal como se vive",D2:"riqueza y sustento",D3:"hermanos y coraje",
  D7:"hijos",D9:"el alma y el cónyuge",D10:"profesión",D12:"los padres",D16:"vehículos y confort",
  D30:"males",D60:"karma acumulado"};
  $("#vg").innerHTML=Object.keys(c.vargas).map(k=>{const gs=c.vargas[k],lg=gs.find(g=>g.nombre==="Lagna");
    return `<div class="vg"><h4>${k}</h4><p>${d[k]||""}</p>${cuadroVed(gs.filter(g=>g.nombre!=="Lagna"),lg?lg.rasiIdx:0,k)}</div>`;}).join("");}
function pintarGocara(c){const g=c.gocara;if(!g)return "";
  let h=`<div class="caja"><h3>Sade Sati</h3>`;
  h+=g.sade.activo
    ? `<p class="yoga"><b>Activo — fase ${g.sade.fase}.</b> ${g.sade.nota}${g.sade.hasta?` Sale hacia <b>${g.sade.hasta}</b>.`:""}</p>`
    : `<p style="margin:0;color:var(--muted)">${g.sade.nota}${g.sade.desde?` El próximo empieza hacia <b style="color:var(--ink)">${g.sade.desde}</b>.`:""}</p>`;
  h+=`</div><div class="caja"><h3>Tránsitos de hoy · ${g.fecha}</h3><table>
    <thead><tr><th colspan="2">Graha</th><th>Posición</th><th class="num">desde Lagna</th><th class="num">desde Luna</th></tr></thead><tbody>`;
  g.transitos.forEach(t=>h+=`<tr><td class="gl">${t.glifo}</td><td>${t.graha}${t.retro?' ℞':''}</td><td>${t.posicion}</td><td class="num">${t.desdeLagna}</td><td class="num">${t.desdeLuna}</td></tr>`);
  return h+`</tbody></table></div>`;}

function pintarDasas(c){let h=`<div class="caja"><h3>Vimśottarī</h3><div class="dasa">`;
  c.dasas.forEach(p=>{h+=`<div class="p ${p.actual?'act':''}"><span>${p.senor}</span><span>${p.desde} → ${p.hasta}</span><span>${p.anios}</span></div>`;
    if(p.actual&&p.sub)p.sub.forEach(b=>h+=`<div class="p b ${b.actual?'act':''}"><span>${b.senor}</span><span>${b.desde} → ${b.hasta}</span><span>${b.anios}</span></div>`);});
  $("#ds").innerHTML=h+`</div></div>`+pintarGocara(c);}

// ═══ comparación: la única pantalla con las dos ═══
async function comparar(){const x=t();
  const d=await (await fetch("/api/comparar?"+params())).json();
  let h=`<div class="caja"><h3>${x.compTit}</h3><p style="margin:0 0 12px;color:var(--muted)">${x.compTxt}</p>
    <p style="margin:0 0 14px"><b>ayanāṁśa ${d.ayanamsa.toFixed(4)}°</b> — ésa es toda la diferencia entre las dos columnas.</p>
    <table><thead><tr><th colspan="2"></th><th>${x.tropical}</th><th>${x.sidereo}</th><th></th></tr></thead><tbody>
    <tr><td></td><td><b>Asc / Lagna</b></td><td>${d.ascendente}</td><td>${d.lagna}</td><td class="cambia">${d.cambiaLagna?x.cambia:""}</td></tr>`;
  d.filas.forEach(f=>h+=`<tr><td class="gl">${f.glifo}</td><td>${f.cuerpo}</td><td>${f.tropical}</td><td>${f.sidereo}</td><td class="cambia">${f.cambia?x.cambia:""}</td></tr>`);
  $("#cmp").innerHTML=h+`</tbody></table></div>`;}

// ═══ curso ═══
function pintarCurso(){const x=t();
  $("#lista").innerHTML=CURSOS[TRAD].map(([f,ti,n])=>
    `<a href="#" data-f="${f}"><b>${n?(n==="·"?"Extra":x.modulo+" "+n):x.indice}</b>${ti}</a>`).join("");
  $("#lista").onclick=async ev=>{const a=ev.target.closest("a");if(!a)return;ev.preventDefault();
    const md=await (await fetch(`curso/${TRAD}/${a.dataset.f}.md`)).text();
    $("#texto").hidden=false;$("#texto").innerHTML=markdown(md);
    $("#texto").scrollIntoView({behavior:"smooth",block:"start"});};}
function markdown(tx){const esc=s=>s.replace(/&/g,"&amp;").replace(/</g,"&lt;");
  const inl=s=>esc(s).replace(/`([^`]+)`/g,"<code>$1</code>").replace(/\*\*([^*]+)\*\*/g,"<strong>$1</strong>").replace(/\*([^*]+)\*/g,"<em>$1</em>");
  const o=[];let tb=null;
  for(const ln of tx.split("\n")){const fila=ln.trim().startsWith("|")&&ln.trim().endsWith("|");
    if(fila){const cl=ln.trim().slice(1,-1).split("|").map(s=>s.trim());
      if(cl.every(c=>/^:?-+:?$/.test(c)))continue;
      if(!tb){tb=1;o.push("<table><tbody>");}
      o.push("<tr>"+cl.map(c=>`<td>${inl(c)}</td>`).join("")+"</tr>");continue;}
    if(tb){o.push("</tbody></table>");tb=null;}
    if(/^#{1,4} /.test(ln)){const n=ln.match(/^#+/)[0].length;o.push(`<h${n}>${inl(ln.slice(n+1))}</h${n}>`);}
    else if(/^> /.test(ln))o.push(`<blockquote>${inl(ln.slice(2))}</blockquote>`);
    else if(/^[-*] /.test(ln))o.push(`<li>${inl(ln.slice(2))}</li>`);
    else if(/^\d+\. /.test(ln))o.push(`<li>${inl(ln.replace(/^\d+\. /,""))}</li>`);
    else if(ln.trim()==="---")o.push("<hr>");
    else if(ln.trim())o.push(`<p>${inl(ln)}</p>`);}
  if(tb)o.push("</tbody></table>");
  return o.join("\n").replace(/(<li>[\s\S]*?<\/li>)(?!\s*<li>)/g,"<ul>$1</ul>");}

// ═══ ejercicio ═══
function pintarEjercicio(){const x=t();
  if($("#ejBox").dataset.listo)return;
  $("#ejBox").dataset.listo="1";
  $("#ejBox").innerHTML=`<h2 style="font-size:1.2rem;margin:0 0 6px">${x.ejTit}</h2>
    <p class="lead">${x.ejTxt}</p><form id="fv">
    <label>Día juliano<input id="jd" class="w"></label>
    <label>T.S. Greenwich °<input id="tsg"></label>
    <label>T.S. local °<input id="tsl"></label>
    <label>Ascendente °<input id="asc"></label>
    <label>Medio Cielo °<input id="mc"></label>
    <button class="go" type="submit">${x.comprobar}</button></form><div id="res"></div>`;
  $("#fv").onsubmit=async e=>{e.preventDefault();
    const ex=["jd","tsg","tsl","asc","mc"].map(k=>`${k}=${$("#"+k).value||0}`).join("&");
    const r=await (await fetch("/api/verificar?"+params()+"&"+ex)).json();
    let h="";r.pasos.forEach(p=>h+=`<div class="paso"><span>${p.nombre}</span><span class="${p.bien?'bien':'malo'}">${p.bien?"✓":"✗"}</span><span style="color:var(--muted)">${p.tuyo}</span><span class="${p.bien?'bien':'malo'}">${p.bien?"":"±"+p.desvio+" "+p.unidad}</span></div>`);
    if(r.primerFallo>=0){const p=r.pasos[r.primerFallo];
      h+=`<div class="aviso"><b>${p.nombre}.</b> ±${p.desvio} ${p.unidad} — ${p.comentario}. Rehaz desde ahí.</div>`;}
    else h+=`<div class="aviso" style="border-color:var(--ok)"><b>Los cinco pasos correctos.</b></div>`;
    $("#res").innerHTML=h;};}

// ═══ ciudades, husos, guardadas ═══
async function buscarCiudad(){const q=$("#ciudad").value,c=$("#sug");
  if(q.length<2){c.innerHTML="";return;}
  const ls=await (await fetch("/api/lugares?q="+encodeURIComponent(q))).json();
  c.innerHTML=(ls||[]).map((l,i)=>`<button type="button" data-i="${i}">${l.nombre}<span>${l.region?l.region+" · ":""}${l.pais}</span></button>`).join("");
  c.dataset.datos=JSON.stringify(ls||[]);}
async function resolverHuso(){const z=$("#zona").value;if(!z)return;
  const [a,m,d]=$("#fecha").value.split("-"),[hh,mm]=$("#hora").value.split(":");
  const r=await (await fetch(`/api/huso?zona=${encodeURIComponent(z)}&anio=${+a}&mes=${+m}&dia=${+d}&hh=${+hh}&mm=${+mm}`)).json();
  if(r.error)return;$("#tz").value=r.offset;
  $("#husoTxt").innerHTML=`UTC${r.offset>=0?"+":""}${r.offset} · ${r.zona}`+(r.verano?` · <b style="color:var(--ac)">horario de verano</b>`:``);}
$("#ciudad").oninput=buscarCiudad;
$("#sug").onclick=async ev=>{const b=ev.target.closest("button");if(!b)return;
  const l=JSON.parse($("#sug").dataset.datos)[+b.dataset.i];
  $("#lat").value=l.lat;$("#lon").value=l.lon;$("#ciudad").value=l.nombre;$("#zona").value=l.zona;$("#sug").innerHTML="";
  await resolverHuso();};
$("#fecha").onchange=resolverHuso;$("#hora").onchange=resolverHuso;
$("#porque").onclick=async()=>{const c=$("#hist");
  if(!c.hidden){c.hidden=true;return;}
  const [a,m,d]=$("#fecha").value.split("-"),[hh,mm]=$("#hora").value.split(":");
  const h=await (await fetch(`/api/husohistoria?zona=${encodeURIComponent($("#zona").value)}&anio=${+a}&mes=${+m}&dia=${+d}&hh=${+hh}&mm=${+mm}&lon=${$("#lon").value}`)).json();
  if(h.error){c.innerHTML=h.error;c.hidden=false;return;}
  const sg=n=>(n>=0?"+":"")+(Math.round(n*100)/100);
  let x=`<h4>De dónde sale UTC${sg(h.offset)}</h4><div class="fila">zona · <b>${h.zona}</b> (${h.abrev})</div>
    <div class="fila">estándar del año · <b>UTC${sg(h.estandar)}</b></div>
    <div class="fila">horario de verano · <b>${h.verano?"sí, +1 h":"no"}</b></div>
    <div class="sec"><h4>El reloj frente al Sol</h4>
    <div class="fila">por longitud le tocaría · UTC${sg(h.solar)}</div>
    <div class="fila">el reloj va <b>${sg(h.desfase)} h</b> respecto al Sol</div></div>`;
  x+=`<div class="sec"><h4>Horario de verano en ${a}</h4>`+
    (h.delAnio.length?h.delAnio.map(v=>`<div class="fila">${v.fecha} · UTC${sg(v.de)} → <b>UTC${sg(v.a)}</b> — ${v.motivo}</div>`).join(""):`<div class="fila">ese año no hubo cambios</div>`)+`</div>`;
  if(h.historicos.length)x+=`<div class="sec"><h4>Cambios de huso del país</h4>`+h.historicos.map(v=>`<div class="fila">${v.fecha} · UTC${sg(v.de)} → <b>UTC${sg(v.a)}</b></div>`).join("")+`</div>`;
  c.innerHTML=x;c.hidden=false;};
async function pintarGuardadas(d){d=d||await (await fetch("/api/guardadas")).json();const c=d.cartas||[];
  $("#guardadas").innerHTML=c.length?`<span class="et">Guardadas</span>`+c.map(v=>`<span class="chip"><button class="abrir" data-id="${v.id}">${v.nombre}<small>${v.ciudad}</small></button><button class="x" data-del="${v.id}">×</button></span>`).join(""):"";
  $("#guardadas").dataset.datos=JSON.stringify(c);}
$("#guardar").onclick=async()=>{const n=prompt("Nombre",$("#ciudad").value||"Carta");if(!n)return;
  pintarGuardadas(await (await fetch("/api/guardadas",{method:"POST",headers:{"Content-Type":"application/json"},
    body:JSON.stringify({nombre:n,ciudad:$("#ciudad").value,zona:$("#zona").value,fecha:$("#fecha").value,
    hora:$("#hora").value,tz:+$("#tz").value,lat:+$("#lat").value,lon:+$("#lon").value})})).json());};
$("#guardadas").onclick=async ev=>{const del=ev.target.closest("[data-del]");
  if(del){const c=JSON.parse($("#guardadas").dataset.datos).find(x=>x.id===del.dataset.del);
    if(!confirm(`¿Borrar «${c.nombre}»?`))return;
    pintarGuardadas(await (await fetch("/api/guardadas?id="+del.dataset.del,{method:"DELETE"})).json());return;}
  const ab=ev.target.closest("[data-id]");if(!ab)return;
  const c=JSON.parse($("#guardadas").dataset.datos).find(x=>x.id===ab.dataset.id);
  $("#ciudad").value=c.ciudad;$("#zona").value=c.zona;$("#fecha").value=c.fecha;$("#hora").value=c.hora;
  $("#tz").value=c.tz;$("#lat").value=c.lat;$("#lon").value=c.lon;await resolverHuso();levantar();};

document.documentElement.dataset.trad=TRAD;
aplicarIdioma(); pintarGuardadas(); levantar();
