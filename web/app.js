const $ = s => document.querySelector(s);
let TRAD = "occidental", LANG = "es", ESTILO = "norte", DATOS = null;
let SECCION = "carta";   // la pestaña abierta; sobrevive a los cambios de idioma
const t = () => T[LANG];

const SIGW = ["♈","♉","♊","♋","♌","♍","♎","♏","♐","♑","♒","♓"];
const nom = i => SIGNOS[LANG][i];          // el signo, en el idioma activo
const cue = n => (CUERPOS[LANG] || {})[n] || n;   // el cuerpo, idem
const SUR = [[1,0],[2,0],[3,0],[3,1],[3,2],[3,3],[2,3],[1,3],[0,3],[0,2],[0,1],[0,0]];
const NORTE = [[.50,.26],[.25,.11],[.11,.25],[.26,.50],[.11,.75],[.25,.89],[.50,.74],[.75,.89],[.89,.75],[.74,.50],[.89,.25],[.75,.11]];
const gms = g => { const d=Math.floor(g), m=Math.round((g-d)*60); return m===60?`${d+1}° 00′`:`${d}° ${String(m).padStart(2,"0")}′`; };

// ═══ selector de tradición e idioma ═══
document.querySelectorAll(".selector button").forEach(b => b.onclick = () => {
  document.querySelectorAll(".selector button").forEach(x => x.classList.toggle("on", x===b));
  TRAD = b.dataset.t; document.documentElement.dataset.trad = TRAD;
  $("#estilosBox").hidden = TRAD !== "jyotisha";
  // Al cambiar de tradición el curso es otro: se cierra el módulo abierto.
  delete $("#texto").dataset.f; $("#texto").hidden = true; $("#texto").innerHTML = "";
  repintarTodo();
});
document.querySelectorAll(".idioma button").forEach(b => b.onclick = () => {
  document.querySelectorAll(".idioma button").forEach(x => x.classList.toggle("on", x===b));
  LANG = b.dataset.i; document.documentElement.lang = LANG;
  repintarTodo();
});

// repintarTodo rehace TODO lo que hay en pantalla en el idioma activo, no solo
// lo que se ve en ese momento: las pestañas ocultas conservan su contenido y
// asomarían en el idioma anterior al abrirlas.
//
// Parte del texto lo compone el servidor —los yogas, la lectura, la historia
// del huso— asi que no basta con repintar: hay que volver a pedirlo.
async function repintarTodo(){
  aplicarIdioma();                                   // rótulos, pestañas, índice del curso
  await pintarGuardadas();                           // el rótulo «Guardadas»
  if ($("#zona").value) await resolverHuso();        // «horario de verano» / «estándar»
  if (!$("#hist").hidden) await pintarHistoria();
  if ($("#texto").dataset.f) await abrirModulo($("#texto").dataset.f);
  if ($("#cmp").innerHTML) await comparar();
  if ($("#prd").innerHTML) await predecir($("#prFecha")?$("#prFecha").value:"");
  if ($("#ejBox").innerHTML) pintarEjercicio();
  if (DATOS) await levantar();                       // vuelve a pedir la carta y la lectura
}
document.querySelectorAll(".estilos button").forEach(b => b.onclick = () => {
  document.querySelectorAll(".estilos button").forEach(x => x.classList.toggle("on", x===b));
  ESTILO = b.dataset.e; if (DATOS) render();
});

function pintarNav() {
  const n = $("#nav");
  // Se conserva la pestaña abierta. Cambiar de idioma no debe echarte de donde
  // estabas: si estás en el curso, sigues en el curso. Solo se vuelve a la
  // primera cuando la pestaña no existe en la tradición nueva —el pañcāṅga no
  // está en occidental— porque entonces no hay adónde volver.
  const abierta = NAV[TRAD].includes(SECCION) ? SECCION : NAV[TRAD][0];
  SECCION = abierta;
  n.innerHTML = NAV[TRAD].map(k =>
    `<button data-s="${k}" class="${k===abierta?"on":""}">${t().nav[k]}</button>`).join("") +
    `<button class="imp" id="imprimir">${t().nav.imprimir}</button>`;
  n.querySelectorAll("button[data-s]").forEach(b => b.onclick = () => {
    SECCION = b.dataset.s;
    n.querySelectorAll("button[data-s]").forEach(x => x.classList.toggle("on", x===b));
    document.querySelectorAll("section").forEach(s => s.classList.toggle("on", s.id===b.dataset.s));
    if (b.dataset.s === "comparar") comparar();
    if (b.dataset.s === "prediccion") predecir();
    if (b.dataset.s === "ejercicio") pintarEjercicio();
  });
  $("#imprimir").onclick = () => window.print();
  document.querySelectorAll("section").forEach(s => s.classList.toggle("on", s.id === abierta));
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
  $("#btOcc").textContent = x.trad_occidental; $("#btJyo").textContent = x.trad_jyotisha;
  $("#btNorte").textContent = x.norte; $("#btSur").textContent = x.sur;
  $("#nodoTxt").textContent = x.nodo + ": " + x.nodo_verdadero;
  $("#nodoBox").title = x.nodo_ayuda;
  // Mientras no se haya resuelto un huso, este hueco lleva el aviso, y el
  // aviso tiene que cambiar de idioma como todo lo demás. Con la marca puesta
  // se sabe que lo que hay es el aviso y no un huso ya calculado.
  if ($("#husoTxt").dataset.aviso !== "no") $("#husoTxt").textContent = x.elige;
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
    s+=`<text x="${S/2}" y="${S+11}" text-anchor="middle" font-size="10" font-family="ui-monospace,Menlo,monospace" fill="currentColor" opacity=".55">${titulo} · ${t().sur} · ${t().surNota}</text>`;
  }else{const x0=P0,y0=P0,x1=P0+L,y1=P0+L,mx=(x0+x1)/2,my=(y0+y1)/2;
    s+=`<rect x="${x0}" y="${y0}" width="${L}" height="${L}" fill="var(--cas)" stroke="currentColor" stroke-width="1.4"/>`;
    s+=`<line x1="${x0}" y1="${y0}" x2="${x1}" y2="${y1}" stroke="currentColor" stroke-width="1.1"/><line x1="${x1}" y1="${y0}" x2="${x0}" y2="${y1}" stroke="currentColor" stroke-width="1.1"/>`;
    s+=`<path d="M ${mx},${y0} L ${x1},${my} L ${mx},${y1} L ${x0},${my} Z" fill="none" stroke="currentColor" stroke-width="1.1"/>`;
    for(let h=0;h<12;h++){const [fx,fy]=NORTE[h],cx=x0+fx*L,cy=y0+fy*L,r=(lagnaRasi+h)%12;
      s+=`<text x="${cx}" y="${cy-13}" text-anchor="middle" font-size="9.5" font-family="ui-monospace,Menlo,monospace" fill="currentColor" opacity=".5">${r+1}</text>`;
      (por[r]||[]).forEach((g,i)=>{const n=(por[r]||[]).length,dx=n>1?((i%2)-.5)*20:0;
        s+=`<text x="${cx+dx}" y="${cy+4+Math.floor(i/2)*15}" text-anchor="middle" font-size="13.5" fill="var(--ac)">${g.glifo}<title>${g.nombre} ${g.posicion||""}</title></text>`;});}
    s+=`<text x="${S/2}" y="${S+11}" text-anchor="middle" font-size="10" font-family="ui-monospace,Menlo,monospace" fill="currentColor" opacity=".55">${titulo} · ${t().norte} · ${t().norteNota}</text>`;}
  return s+"</svg>";
}

// ═══ tablas ═══
function tablasOcc(c){const x=t();
  const ang=[x.ang_asc,x.ang_dsc,x.ang_mc,x.ang_ic];
  let h=`<div class="caja"><h3>${x.angulos}</h3><table><tbody>`;
  c.angulos.forEach((l,i)=>{const s=Math.floor(l/30);
    h+=`<tr><td>${ang[i]}</td><td class="num">${gms(l-s*30)}</td><td class="gl">${SIGW[s]}</td><td>${nom(s)}</td></tr>`;});
  h+=`</tbody></table></div><div class="caja"><h3>${x.planetas}</h3><table><thead><tr><th colspan="2">${x.planetas}</th><th class="num">${x.grados}</th><th colspan="2">${x.signo}</th><th class="num">${x.casa}</th></tr></thead><tbody>`;
  c.cuerpos.forEach(p=>h+=`<tr><td class="gl">${p.glifo}</td><td>${cue(p.nombre)}${p.retro?' ℞':''}</td><td class="num">${gms(p.grado)}</td><td class="gl">${p.glifoSig}</td><td>${nom(p.signoIdx)}</td><td class="num">${p.casaP}</td></tr>`);
  h+=`</tbody></table></div><div class="caja"><h3>${x.aspectos}</h3><table><tbody>`;
  c.aspectos.slice(0,14).forEach(a=>h+=`<tr><td>${cue(a.a)}</td><td class="gl">${a.glifo}</td><td>${a.nombre}</td><td>${cue(a.b)}</td><td class="num">${a.orbe.toFixed(2)}°</td></tr>`);
  h+=`</tbody></table></div><div class="caja"><h3>${x.casas}</h3><table><thead><tr><th>${x.casa}</th><th>${x.signo}</th><th>${x.senor}</th><th>${x.alojado}</th></tr></thead><tbody>`;
  c.regentes.forEach((r,i)=>{const s=Math.floor(c.cuspP[i]/30);
    h+=`<tr><td class="num">${i+1}</td><td class="gl">${SIGW[s]}</td><td>${cue(r)}</td><td>${x.casa} ${c.regenteEn[i]}</td></tr>`;});
  return h+`</tbody></table></div>`;}

function tablasVed(c){const x=t();
  let h=`<div class="caja"><h3>${x.lagnaT}</h3><table><tbody>
    <tr><td>${x.lagnaT}</td><td>${c.lagnaPos}</td><td>${c.lagnaNak} · pada ${c.lagnaPada}</td></tr>
    <tr><td>${x.senor}</td><td colspan="2">${cue(c.senorLagna)}</td></tr>
    <tr><td>ayanāṁśa</td><td colspan="2">${c.ayanamsa.toFixed(4)}°</td></tr></tbody></table></div>`;
  h+=`<div class="caja"><h3>${x.grahas}</h3><table><thead><tr><th colspan="2">${x.grahas}</th><th>${x.posicion}</th><th>${x.nak}</th><th class="num">${x.pada}</th><th class="num">${x.casa}</th><th>${x.estado}</th></tr></thead><tbody>`;
  c.grahas.forEach(g=>{const cl=g.dignidad==="exaltado"?"exalt":g.dignidad==="debilitado"?"debil":"";
    const marcas=[g.combusto?`<span class="debil" title="combusto: a ${g.delSol}° del Sol">☌sol</span>`:"",
      g.digBala?`<span class="exalt" title="fuerza direccional plena">dig</span>`:""].filter(Boolean).join(" ");
    h+=`<tr><td class="gl">${g.glifo}</td><td>${cue(g.nombre)}${g.retro?' ℞':''}</td><td>${g.posicion}${g.gandanta?' <span class="gan" title="gaṇḍānta">⚠</span>':''}</td><td>${g.nak}</td><td class="num">${g.pada}</td><td class="num">${g.bhava}</td><td class="${cl}">${g.dignidad} ${marcas}</td></tr>`;});
  h+=`</tbody></table></div><div class="caja"><h3>${x.bhavas}</h3><table><thead><tr><th class="num">${x.casa}</th><th>${x.rasi}</th><th>${x.senor}</th><th>${x.alojado}</th><th>${x.ocupan}</th><th>${x.aspectan}</th></tr></thead><tbody>`;
  c.bhavas.forEach(b=>h+=`<tr><td class="num">${b.numero}</td><td>${b.rasi}</td><td>${cue(b.senor)}</td><td>${x.casa} ${b.senorEn}</td><td>${(b.ocupan||[]).map(cue).join(" ")||"—"}</td><td style="color:var(--muted)">${(b.aspectan||[]).map(cue).join(" ")||"—"}</td></tr>`);
  h+=`</tbody></table></div><div class="caja"><h3>${x.karakas}</h3>
    ${c.karakamsa?`<p style="margin:0 0 8px;color:var(--muted);font-size:.85rem">${x.karakamsa_txt.replace("{r}",`<b style="color:var(--ink)">${c.karakamsa}</b>`)}</p>`:""}
    <table><tbody>`;
  const kn={AK:"Ātmakāraka",AmK:"Amātyakāraka",BK:"Bhrātṛkāraka",MK:"Mātṛkāraka",PiK:"Pitṛkāraka",PK:"Putrakāraka",GK:"Ñātikāraka"};
  Object.keys(kn).forEach(k=>h+=`<tr><td><b>${k}</b></td><td>${cue(c.karakas[k]||"—")}</td><td style="color:var(--muted)">${kn[k]}</td></tr>`);
  h+=`</tbody></table></div>`;
  if(c.yogas&&c.yogas.length){h+=`<div class="caja"><h3>${x.yogas}</h3>`;c.yogas.forEach(y=>h+=`<p class="yoga">${y}</p>`);h+=`</div>`;}
  return h;}

// ═══ render ═══
function render(){
  if(TRAD==="occidental"){$("#dibujo").innerHTML=ruedaOcc(DATOS);$("#tablas").innerHTML=tablasOcc(DATOS);}
  else{$("#dibujo").innerHTML=cuadroVed(DATOS.grahas,DATOS.lagnaRasi,"D1");$("#tablas").innerHTML=tablasVed(DATOS);
    pintarVargas(DATOS);pintarDasas(DATOS);pintarPancanga(DATOS);pintarFuerza(DATOS);}
}
function params(){const [a,m,d]=$("#fecha").value.split("-"),[hh,mm]=$("#hora").value.split(":");
  return `anio=${+a}&mes=${+m}&dia=${+d}&hh=${+hh}&mm=${+mm}&tz=${$("#tz").value}&lat=${$("#lat").value}&lon=${$("#lon").value}`;}
async function levantar(e){if(e)e.preventDefault();
  const nodo=TRAD==="jyotisha"&&$("#nodoV")&&$("#nodoV").checked?"&nodo=verdadero":"";
  DATOS=await (await fetch(`/api/${TRAD==="jyotisha"?"vedica":"carta"}?`+params()+nodo+"&lang="+LANG)).json();
  const [Y,M,D]=$("#fecha").value.split("-"),[H,Mi]=$("#hora").value.split(":");
  const extra=TRAD==="jyotisha"?` · ${t().lagnaT} ${DATOS.lagnaPos}`:"";
  $("#ficha").innerHTML=`<h2>${$("#ciudad").value||"—"}</h2><p>${D}/${M}/${Y} · ${H}:${Mi} · UTC${+$("#tz").value>=0?"+":""}${$("#tz").value} · ${DATOS.ut}${extra}</p>`;
  render(); await leer();}
$("#f").onsubmit=levantar;

async function leer(){const url=(TRAD==="jyotisha"?"/api/lecturaved?":"/api/lectura?")+params()+"&lang="+LANG;
  const L=await (await fetch(url)).json();
  const cats={};L.frases.forEach(f=>(cats[f.categoria]=cats[f.categoria]||[]).push(f));
  let h=`<div class="caja"><h3>${t().dominante}</h3><p style="margin:0">${L.dominante}</p></div>`;
  for(const c in cats){h+=`<div class="caja"><h3>${c}</h3>`;
    cats[c].sort((a,b)=>b.peso-a.peso).forEach(f=>h+=`<p style="margin:0 0 9px">${f.texto}<br><span style="font-family:ui-monospace,Menlo,monospace;font-size:.71rem;color:var(--muted)">← ${f.fuente}</span></p>`);h+=`</div>`;}
  if(L.contradicciones.length){h+=`<div class="caja"><h3>${t().contradicciones}</h3>`;
    L.contradicciones.forEach(x=>h+=`<p style="margin:0 0 9px">${x}</p>`);h+=`</div>`;}
  $("#lec").innerHTML=h+`<div class="aviso">${L.nota}</div>`;}

// ═══ pañcāṅga, arudhas y lagnas especiales ═══
function pintarPancanga(c){const x=t(),p=c.pancanga;if(!p)return;
  const barra=(pct,txt)=>`<div style="margin:0 0 10px"><div style="display:flex;justify-content:space-between;font-size:.78rem"><span>${txt}</span><span style="color:var(--muted)">${pct.toFixed(0)}% ${x.pc_recorrido}</span></div>
    <div style="height:4px;background:var(--linea);border-radius:2px;margin-top:3px"><div style="height:4px;width:${pct}%;background:var(--acento);border-radius:2px"></div></div></div>`;
  let h=`<div class="caja"><h3>${x.nav.pancanga}</h3><p style="color:var(--muted);margin:0 0 12px">${x.pc_lead}</p>
    <table><tbody>
    <tr><td><b>${x.pc_tithi}</b></td><td>${p.tithi}</td><td>${p.paksha} ${x.pc_paksha} · ${p.tithiNum}/30</td></tr>
    <tr><td><b>${x.pc_vara}</b></td><td>${p.vara}</td><td>${x.pc_senor}: ${cue(p.senorVara)}</td></tr>
    <tr><td><b>${x.pc_nak}</b></td><td>${p.nakshatra}</td><td>pada ${p.pada} · ${x.pc_senor}: ${cue(p.senorNak)}</td></tr>
    <tr><td><b>${x.pc_yoga}</b></td><td>${p.yoga}</td><td>${p.yogaNum}/27</td></tr>
    <tr><td><b>${x.pc_karana}</b></td><td>${p.karana}</td><td>${p.visti?"⚠":""}</td></tr>
    </tbody></table>`;
  h+=barra(p.tithiPct,x.pc_tithi)+barra(p.nakPct,x.pc_nak)+barra(p.luna,x.pc_luna);
  if(p.visti)h+=`<div class="aviso">${x.pc_visti}</div>`;
  h+=`</div>`;
  // lagnas especiales
  const l=c.lagnasEsp;
  if(l&&l.hay){const hh=Math.floor(l.amanece),mm=Math.round((l.amanece-hh)*60);
    h+=`<div class="caja"><h3>${x.pc_lagnas}</h3><p style="color:var(--muted);margin:0 0 4px">${x.pc_lagnasLead}</p>
      <p style="color:var(--muted);margin:0 0 12px">${x.pc_amanece}: ${String(hh).padStart(2,"0")}:${String(mm).padStart(2,"0")} UT</p><table><tbody>`;
    [["bhava",x.pc_bl],["hora",x.pc_hl],["ghati",x.pc_gl]].forEach(([k,d])=>{
      h+=`<tr><td>${gms(l[k]%30)}</td><td><b>${RASIS[Math.floor(l[k]/30)]}</b></td><td style="color:var(--muted)">${d}</td></tr>`;});
    h+=`</tbody></table></div>`;}
  // arudhas
  const a=c.arudhas;
  if(a){h+=`<div class="caja"><h3>${x.ar_titulo}</h3><p style="color:var(--muted);margin:0 0 12px">${x.ar_lead}</p>
    <table><tbody><tr><td colspan="2"><b>${x.ar_al}</b></td><td><b>${RASIS[a.al]}</b></td></tr>
    <tr><td colspan="2"><b>${x.ar_ul}</b></td><td><b>${RASIS[a.ul]}</b></td></tr>
    <tr><td colspan="3" style="height:6px"></td></tr>`;
    a.padas.forEach((r,i)=>{h+=`<tr><td>A${i+1}</td><td style="color:var(--muted)">${x.ar_bhava} ${i+1}</td><td>${RASIS[r]}</td></tr>`;});
    h+=`</tbody></table></div>`;}
  $("#pc").innerHTML=h;}

// ═══ aṣṭakavarga y ṣaḍbala ═══
function pintarFuerza(c){const x=t(),a=c.ashtaka,s=c.shadbala;if(!a)return;
  // el SAV se pinta como barras: la media es 28 y la vista tiene que enseñar
  // de un vistazo qué casas van sobradas y cuáles no llegan.
  const max=Math.max(...a.sav);
  let h=`<div class="caja"><h3>${x.av_titulo} · ${x.av_sav}</h3>
    <p style="color:var(--muted);margin:0 0 12px">${x.av_lead}<br>${x.av_savLead}</p><table><tbody>`;
  a.sav.forEach((n,i)=>{const bh=((i-c.lagnaRasi)%12+12)%12+1;
    const col=n>=30?"var(--acento)":n<25?"var(--muted)":"var(--texto)";
    h+=`<tr><td style="width:2.2em;color:var(--muted)">${bh}</td><td style="width:6em">${RASIS[i]}</td>
      <td style="width:2.4em;text-align:right;font-variant-numeric:tabular-nums;color:${col}"><b>${n}</b></td>
      <td><div style="height:8px;width:${n/max*100}%;background:${col};border-radius:2px;min-width:2px"></div></td></tr>`;});
  h+=`</tbody></table><p style="color:var(--muted);margin:10px 0 0">${x.av_total}: ${a.total} · ${x.av_media} ${a.media.toFixed(1)}</p></div>`;
  // BAV, un graha por fila
  h+=`<div class="caja"><h3>${x.av_bav}</h3><p style="color:var(--muted);margin:0 0 12px">${x.av_bavLead}</p>
    <div style="overflow-x:auto"><table><tbody><tr><td></td>`;
  RASIS.forEach(r=>h+=`<td style="font-size:.7rem;color:var(--muted)">${r.slice(0,4)}</td>`);
  h+=`</tr>`;
  Object.keys(a.bav).forEach(g=>{h+=`<tr><td><b>${cue(g)}</b></td>`;
    a.bav[g].forEach(n=>h+=`<td style="text-align:center;font-variant-numeric:tabular-nums;color:${n>=5?"var(--acento)":n<=2?"var(--muted)":"inherit"}">${n}</td>`);
    h+=`</tr>`;});
  h+=`</tbody></table></div></div>`;
  // ṣaḍbala
  if(s&&s.balas){h+=`<div class="caja"><h3>${x.sb_titulo}</h3><p style="color:var(--muted);margin:0 0 12px">${x.sb_lead}</p>
    <div style="overflow-x:auto"><table><tbody><tr>
    <td></td><td>${x.sb_graha}</td><td>${x.sb_sthana}</td><td>${x.sb_dig}</td><td>${x.sb_kala}</td>
    <td>${x.sb_chesta}</td><td>${x.sb_nais}</td><td>${x.sb_drik}</td><td>${x.sb_yuddha}</td><td><b>${x.sb_rupas}</b></td>
    <td>${x.sb_min}</td><td><b>${x.sb_razon}</b></td></tr>`;
    const n1=v=>v.toFixed(0),n2=v=>v.toFixed(2);
    [...s.balas].sort((p,q)=>p.rango-q.rango).forEach(b=>{
      const corto=b.razon<1;
      h+=`<tr><td style="color:var(--muted)">${b.rango}</td><td><b>${cue(b.graha)}</b></td>
        <td>${n1(b.sthana)}</td><td>${n1(b.dig)}</td><td>${n1(b.kala)}</td><td>${n1(b.chesta)}</td>
        <td>${n1(b.naisargika)}</td><td>${n1(b.drik)}</td>
        <td${b.rival?` title="${x.sb_rival} ${cue(b.rival)}"`:""}>${b.rival?(b.yuddha>0?"+":"")+n1(b.yuddha)+" ⚔":"—"}</td>
        <td style="font-variant-numeric:tabular-nums"><b>${n2(b.rupas)}</b></td>
        <td style="color:var(--muted)">${b.minimo.toFixed(1)}</td>
        <td style="font-variant-numeric:tabular-nums;color:${corto?"var(--muted)":"var(--acento)"}"><b>${n2(b.razon)}</b>${corto?" ↓":""}</td></tr>`;});
    h+=`</tbody></table></div><div class="aviso">${s.nota}</div></div>`;}
  // bhāva bala: la fuerza del asunto, no la del planeta
  if(s&&s.bhavas&&s.bhavas.length){const mx=Math.max(...s.bhavas.map(b=>b.total));
    h+=`<div class="caja"><h3>${x.bb_titulo}</h3><p style="color:var(--muted);margin:0 0 12px">${x.bb_lead}</p>
      <div style="overflow-x:auto"><table><tbody><tr><td></td><td>${x.bb_bhava}</td><td>${x.bb_senor}</td>
      <td>${x.bb_dig}</td><td>${x.bb_drishti}</td><td><b>${x.sb_rupas}</b></td><td></td></tr>`;
    [...s.bhavas].sort((p,q)=>p.rango-q.rango).forEach(b=>{
      h+=`<tr><td style="color:var(--muted)">${b.rango}</td><td><b>${b.numero}</b> ${RASIS[(c.lagnaRasi+b.numero-1)%12]}</td>
        <td>${b.senor.toFixed(0)}</td><td>${b.dig>0?"+":""}${b.dig.toFixed(0)}</td>
        <td>${b.drishti>0?"+":""}${b.drishti.toFixed(1)}</td>
        <td style="font-variant-numeric:tabular-nums"><b>${b.rupas.toFixed(2)}</b></td>
        <td style="width:38%"><div style="height:8px;width:${Math.max(2,b.total/mx*100)}%;background:var(--acento);border-radius:2px"></div></td></tr>`;});
    h+=`</tbody></table></div></div>`;}
  $("#fz").innerHTML=h;}

function pintarVargas(c){const x=t(); const d=k=>x["v_"+k]||"";
  $("#vg").innerHTML=Object.keys(c.vargas).map(k=>{const gs=c.vargas[k],lg=gs.find(g=>g.nombre==="Lagna");
    return `<div class="vg"><h4>${k}</h4><p>${d(k)}</p>${cuadroVed(gs.filter(g=>g.nombre!=="Lagna"),lg?lg.rasiIdx:0,k)}</div>`;}).join("");}
function pintarGocara(c){const g=c.gocara;if(!g)return "";
  let h=`<div class="caja"><h3>${t().sadeTit}</h3>`;
  h+=g.sade.activo
    ? `<p class="yoga"><b>${t().sadeAct.replace("{f}",g.sade.fase)}</b> ${g.sade.nota}${g.sade.hasta?t().sadeSale.replace("{d}","<b>"+g.sade.hasta+"</b>"):""}</p>`
    : `<p style="margin:0;color:var(--muted)">${g.sade.nota}${g.sade.desde?t().sadeProx.replace("{d}","<b style=\"color:var(--ink)\">"+g.sade.desde+"</b>"):""}</p>`;
  h+=`</div><div class="caja"><h3>${t().transitos.replace("{f}",g.fecha)}</h3><table>
    <thead><tr><th colspan="2">${t().graha}</th><th>${t().posicion}</th><th class="num">${t().desdeLagna}</th><th class="num">${t().desdeLuna}</th></tr></thead><tbody>`;
  g.transitos.forEach(t=>h+=`<tr><td class="gl">${t.glifo}</td><td>${cue(t.graha)}${t.retro?' ℞':''}</td><td>${t.posicion}</td><td class="num">${t.desdeLagna}</td><td class="num">${t.desdeLuna}</td></tr>`);
  return h+`</tbody></table></div>`;}

function pintarDasas(c){let h=`<div class="caja"><h3>${t().vimsottari}</h3><div class="dasa">`;
  c.dasas.forEach(p=>{h+=`<div class="p ${p.actual?'act':''}"><span>${cue(p.senor)}</span><span>${p.desde} → ${p.hasta}</span><span>${p.anios}</span></div>`;
    if(p.actual&&p.sub)p.sub.forEach(b=>h+=`<div class="p b ${b.actual?'act':''}"><span>${cue(b.senor)}</span><span>${b.desde} → ${b.hasta}</span><span>${b.anios}</span></div>`);});
  $("#ds").innerHTML=h+`</div></div>`+pintarGocara(c);}

// ═══ comparación: la única pantalla con las dos ═══
async function comparar(){const x=t();
    // El comparador tiene que enseñar el mismo nodo que la carta; si no, Rāhu
  // saldría en un sitio en una pestaña y en otro en la de al lado.
  const nodo=$("#nodoV")&&$("#nodoV").checked?"&nodo=verdadero":"";
  const d=await (await fetch("/api/comparar?"+params()+nodo)).json();
  let h=`<div class="caja"><h3>${x.compTit}</h3><p style="margin:0 0 12px;color:var(--muted)">${x.compTxt}</p>
    <p style="margin:0 0 14px">${x.comp_ayan.replace("{a}",`<b>${d.ayanamsa.toFixed(4)}</b>`)}</p>
    <table><thead><tr><th colspan="2"></th><th>${x.tropical}</th><th>${x.sidereo}</th><th></th></tr></thead><tbody>
    <tr><td></td><td><b>${x.ang_asc} / ${x.lagnaT}</b></td><td>${gms(d.ascGr)} ${nom(d.ascIdx)}</td><td>${gms(d.lagGr)} ${RASIS[d.lagIdx]}</td><td class="cambia">${d.cambiaLagna?x.cambia:""}</td></tr>`;
  d.filas.forEach(f=>h+=`<tr><td class="gl">${f.glifo}</td><td>${cue(f.cuerpo)}</td><td>${gms(f.tropGr)} ${nom(f.tropIdx)}</td><td>${gms(f.sidGr)} ${RASIS[f.sidIdx]}</td><td class="cambia">${f.cambia?x.cambia:""}</td></tr>`);
  h+=`</tbody></table>`;
  const nS=d.filas.filter(f=>f.cambia).length, nC=d.filas.filter(f=>f.cambiaC).length;
  h+=`<p style="margin:12px 0 0;color:var(--muted)">${x.comp_cuenta
        .replace("{n}",`<b>${nC}</b>`).replace("{t}",d.filas.length).replace("{s}",`<b>${nS}</b>`)}</p></div>`;
  // Las casas: la otra diferencia, y la que más cambia una lectura.
  h+=`<div class="caja"><h3>${x.comp_casas}</h3>
    <p style="margin:0 0 12px;color:var(--muted)">${x.comp_casasTxt}</p>
    <table><thead><tr><th colspan="2"></th><th>${x.comp_placido}</th><th>${x.comp_entero}</th><th></th></tr></thead><tbody>`;
  d.filas.forEach(f=>h+=`<tr><td class="gl">${f.glifo}</td><td>${cue(f.cuerpo)}</td>
    <td>${f.casaOcc||"—"}</td><td>${f.casaVed||"—"}</td>
    <td class="cambia">${f.cambiaC?x.comp_otraCasa:""}</td></tr>`);
  $("#cmp").innerHTML=h+`</tbody></table></div>`;}

// ═══ predicción occidental ═══
async function predecir(cuando){const x=t();
  const f=cuando||($("#prFecha")&&$("#prFecha").value)||"";
  const p=await (await fetch("/api/prediccion?"+params()+"&lang="+LANG+(f?"&cuando="+f:""))).json();
  const sig=l=>`${gms(l%30)} ${nom(Math.floor(l/30))}`;
  let h=`<div class="caja"><h3>${x.nav.prediccion}</h3>
    <p style="color:var(--muted);margin:0 0 12px">${x.pr_lead}</p>
    <label style="display:inline-flex;align-items:center;gap:8px;font-size:.85rem">
      ${x.pr_fecha} <input type="date" id="prFecha" value="${f}">
      <button class="go" id="prHoy" style="padding:4px 12px">${x.pr_hoy}</button>
    </label>
    <span style="color:var(--muted);margin-left:12px">${x.pr_edad} ${p.edad}</span></div>`;

  // convergencias primero: es lo único que marca un periodo
  h+=`<div class="caja"><h3>${x.pr_conv}</h3><p style="color:var(--muted);margin:0 0 12px">${x.pr_convTxt}</p>`;
  p.convergencias.forEach(c=>h+=`<p style="margin:0 0 9px">${c}</p>`);
  h+=`</div>`;

  h+=`<div class="caja"><h3>${x.pr_transitos}</h3><p style="color:var(--muted);margin:0 0 12px">${x.pr_transitosTxt}</p>
    <div style="overflow-x:auto"><table><thead><tr><th>${x.pr_planeta}</th><th></th><th>${x.pr_aspecto}</th>
    <th>${x.pr_natal}</th><th>${x.pr_orbe}</th><th></th><th>${x.pr_casa}</th><th>${x.pr_pasadas}</th></tr></thead><tbody>`;
  p.transitos.forEach(v=>h+=`<tr><td><b>${cue(v.planeta)}</b>${v.retro?" ℞":""}</td><td class="gl">${v.glifo}</td>
    <td>${v.aspecto}</td><td>${cue(v.natal)}</td>
    <td style="font-variant-numeric:tabular-nums">${v.orbe.toFixed(2)}°</td>
    <td style="color:var(--muted)">${v.aplica?x.pr_aplica:x.pr_separa}</td>
    <td>${v.casa||"—"}</td><td>${v.pasadas}</td></tr>`);
  h+=`</tbody></table></div><div class="aviso">${x.pr_pasadasTxt}</div></div>`;

  h+=`<div class="caja"><h3>${x.pr_prog}</h3><p style="color:var(--muted);margin:0 0 12px">${x.pr_progTxt}</p><table><tbody>`;
  p.progresiones.forEach(g=>h+=`<tr><td><b>${cue(g.planeta)}</b></td><td>${gms(g.grado)} ${nom(g.signoIdx)}</td>
    <td>${x.pr_casa} ${g.casa||"—"}</td><td style="color:var(--muted)">${g.aspecto?`${g.aspecto} → ${cue(g.natal)} (${g.orbe}°)`:""}</td></tr>`);
  h+=`</tbody></table></div>`;

  const r=p.revolucion;
  h+=`<div class="caja"><h3>${x.pr_rs}</h3><p style="color:var(--muted);margin:0 0 12px">${x.pr_rsTxt}</p><table><tbody>
    <tr><td>${x.pr_rsCuando}</td><td><b>${r.cuando}</b></td></tr>
    <tr><td>${x.pr_rsAsc}</td><td><b>${sig(r.asc)}</b></td></tr>
    <tr><td>${x.pr_rsMC}</td><td><b>${sig(r.mc)}</b></td></tr>
    <tr><td>${x.pr_rsCasa}</td><td><b>${r.casa||"—"}</b></td></tr>
    </tbody></table></div><div class="aviso">${p.nota}</div>`;
  $("#prd").innerHTML=h;
  $("#prFecha").onchange=()=>predecir($("#prFecha").value);
  $("#prHoy").onclick=e=>{e.preventDefault();predecir("");};}

// ═══ curso ═══
function pintarCurso(){const x=t();
  const nota = "";   // los módulos ya existen en los dos idiomas
  $("#lista").innerHTML=nota+(LANG==="en"?CURSOS_EN:CURSOS_ES)[TRAD].map(([f,ti,n])=>
    `<a href="#" data-f="${f}"><b>${n?(n==="·"?x.extra:x.modulo+" "+n):x.indice}</b>${ti}</a>`).join("");
  $("#lista").onclick=async ev=>{const a=ev.target.closest("a");if(!a)return;ev.preventDefault();
    await abrirModulo(a.dataset.f);
    $("#texto").scrollIntoView({behavior:"smooth",block:"start"});};}

// abrirModulo deja anotado cuál está abierto, para poder recargarlo en el otro
// idioma sin que el lector pierda el sitio.
async function abrirModulo(f){
  const sub=LANG==="en"?"en/":"";   // los módulos traducidos viven en en/
  const md=await (await fetch(`curso/${TRAD}/${sub}${f}.md`)).text();
  $("#texto").dataset.f=f;$("#texto").hidden=false;$("#texto").innerHTML=markdown(md);}
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
  // Se vuelve a pintar cada vez, porque si no nunca cambiaría de idioma. Lo
  // que el alumno lleve escrito se guarda antes y se devuelve después: perder
  // sus cuentas por tocar el switch sería una faena.
  const campos=["jd","tsg","tsl","asc","mc"];
  const escrito={};
  if($("#ejBox").innerHTML) campos.forEach(k=>{const n=$("#"+k); if(n) escrito[k]=n.value;});
  const resultado=$("#res")?$("#res").innerHTML:"";
  $("#ejBox").innerHTML=`<h2 style="font-size:1.2rem;margin:0 0 6px">${x.ejTit}</h2>
    <p class="lead">${x.ejTxt}</p><form id="fv">
    <label>${x.ej_jd}<input id="jd" class="w"></label>
    <label>${x.ej_tsg}<input id="tsg"></label>
    <label>${x.ej_tsl}<input id="tsl"></label>
    <label>${x.ej_asc}<input id="asc"></label>
    <label>${x.ej_mc}<input id="mc"></label>
    <button class="go" type="submit">${x.comprobar}</button></form><div id="res"></div>`;
  $("#fv").onsubmit=async e=>{e.preventDefault();
    const ex=["jd","tsg","tsl","asc","mc"].map(k=>`${k}=${$("#"+k).value||0}`).join("&");
    const r=await (await fetch("/api/verificar?"+params()+"&"+ex)).json();
    let h="";r.pasos.forEach(p=>h+=`<div class="paso"><span>${p.nombre}</span><span class="${p.bien?'bien':'malo'}">${p.bien?"✓":"✗"}</span><span style="color:var(--muted)">${p.tuyo}</span><span class="${p.bien?'bien':'malo'}">${p.bien?"":"±"+p.desvio+" "+p.unidad}</span></div>`);
    if(r.primerFallo>=0){const p=r.pasos[r.primerFallo];
      h+=`<div class="aviso"><b>${p.nombre}.</b> ±${p.desvio} ${p.unidad} — ${p.comentario}. ${t().ejRehaz}</div>`;}
    else h+=`<div class="aviso" style="border-color:var(--ok)"><b>${t().ejOk}</b></div>`;
    $("#res").innerHTML=h;};
  campos.forEach(k=>{const n=$("#"+k); if(n&&escrito[k]!==undefined) n.value=escrito[k];});
  // El resultado anterior se deja tal cual: lo compone el servidor y volver a
  // pedirlo sin que el alumno lo pida sería recalcularle la corrección.
  if(resultado&&$("#res")) $("#res").innerHTML=resultado;}

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
  $("#husoTxt").dataset.aviso="no";
  $("#husoTxt").innerHTML=`UTC${r.offset>=0?"+":""}${r.offset} · ${r.zona}`
    +(r.verano?` · <b style="color:var(--ac)">${t().verano}</b>`:` · ${t().estandar}`);}
// La cabecera fija solo marca su raya y su sombra cuando de verdad hay algo
// por encima; si no, se vería una línea suelta con la página sin desplazar.
window.addEventListener("scroll", () => {
  const h = document.querySelector(".fija");
  if (h) h.classList.toggle("pegada", window.scrollY > 4);
}, { passive: true });

$("#ciudad").oninput=buscarCiudad;
$("#sug").onclick=async ev=>{const b=ev.target.closest("button");if(!b)return;
  const l=JSON.parse($("#sug").dataset.datos)[+b.dataset.i];
  $("#lat").value=l.lat;$("#lon").value=l.lon;$("#ciudad").value=l.nombre;$("#zona").value=l.zona;$("#sug").innerHTML="";
  await resolverHuso();};
$("#fecha").onchange=resolverHuso;$("#hora").onchange=resolverHuso;
// Cambiar de nodo obliga a recalcular: mueve a Rāhu y a Ketu de pada.
$("#nodoV").onchange=()=>{if(DATOS)levantar();};
$("#porque").onclick=()=>{const c=$("#hist");
  if(!c.hidden){c.hidden=true;return;}
  pintarHistoria();};

async function pintarHistoria(){const c=$("#hist");
  const [a,m,d]=$("#fecha").value.split("-"),[hh,mm]=$("#hora").value.split(":");
  const h=await (await fetch(`/api/husohistoria?zona=${encodeURIComponent($("#zona").value)}&anio=${+a}&mes=${+m}&dia=${+d}&hh=${+hh}&mm=${+mm}&lon=${$("#lon").value}`)).json();
  if(h.error){c.innerHTML=h.error;c.hidden=false;return;}
  const sg=n=>(n>=0?"+":"")+(Math.round(n*100)/100);
  const L=t();
  let x=`<h4>${L.hDe.replace("{o}",sg(h.offset))}</h4><div class="fila">${L.hZona} · <b>${h.zona}</b> (${h.abrev})</div>
    <div class="fila">${L.hEstandar} · <b>UTC${sg(h.estandar)}</b></div>
    <div class="fila">${L.hVerano} · <b>${h.verano?L.hSi:L.hNo}</b></div>
    <div class="sec"><h4>${L.hRelojSol}</h4>
    <div class="fila">${L.hTocaria} · UTC${sg(h.solar)}</div>
    <div class="fila">${L.hVa.replace("{d}","<b>"+sg(h.desfase)+"</b>")}</div></div>`;
  x+=`<div class="sec"><h4>${L.hDstAnio.replace("{a}",a)}</h4>`+
    (h.delAnio.length?h.delAnio.map(v=>`<div class="fila">${v.fecha} · UTC${sg(v.de)} → <b>UTC${sg(v.a)}</b> — ${v.motivo}</div>`).join(""):`<div class="fila">${L.hSinDst}</div>`)+`</div>`;
  if(h.historicos.length)x+=`<div class="sec"><h4>${L.hCambios}</h4>`+h.historicos.map(v=>`<div class="fila">${v.fecha} · UTC${sg(v.de)} → <b>UTC${sg(v.a)}</b></div>`).join("")+`</div>`;
  c.innerHTML=x;c.hidden=false;}
async function pintarGuardadas(d){d=d||await (await fetch("/api/guardadas")).json();const c=d.cartas||[];
  $("#guardadas").innerHTML=c.length?`<span class="et">${t().guardadas}</span>`+c.map(v=>`<span class="chip"><button class="abrir" data-id="${v.id}">${v.nombre}<small>${v.ciudad}</small></button><button class="x" data-del="${v.id}">×</button></span>`).join(""):"";
  $("#guardadas").dataset.datos=JSON.stringify(c);}
$("#guardar").onclick=async()=>{const n=prompt(t().nombrePrompt,$("#ciudad").value||"Carta");if(!n)return;
  pintarGuardadas(await (await fetch("/api/guardadas",{method:"POST",headers:{"Content-Type":"application/json"},
    body:JSON.stringify({nombre:n,ciudad:$("#ciudad").value,zona:$("#zona").value,fecha:$("#fecha").value,
    hora:$("#hora").value,tz:+$("#tz").value,lat:+$("#lat").value,lon:+$("#lon").value})})).json());};
$("#guardadas").onclick=async ev=>{const del=ev.target.closest("[data-del]");
  if(del){const c=JSON.parse($("#guardadas").dataset.datos).find(x=>x.id===del.dataset.del);
    if(!confirm(t().borrar.replace("{n}",c.nombre)))return;
    pintarGuardadas(await (await fetch("/api/guardadas?id="+del.dataset.del,{method:"DELETE"})).json());return;}
  const ab=ev.target.closest("[data-id]");if(!ab)return;
  const c=JSON.parse($("#guardadas").dataset.datos).find(x=>x.id===ab.dataset.id);
  $("#ciudad").value=c.ciudad;$("#zona").value=c.zona;$("#fecha").value=c.fecha;$("#hora").value=c.hora;
  $("#tz").value=c.tz;$("#lat").value=c.lat;$("#lon").value=c.lon;await resolverHuso();levantar();};

document.documentElement.dataset.trad=TRAD;
aplicarIdioma(); pintarGuardadas(); levantar();
