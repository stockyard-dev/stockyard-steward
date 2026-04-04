package server
import "net/http"
func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) { w.Header().Set("Content-Type", "text/html"); w.Write([]byte(dashHTML)) }
const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Steward</title><link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet"><style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--mono:'JetBrains Mono',monospace}*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}.main{padding:1.5rem;max-width:960px;margin:0 auto}.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap}.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}.search:focus{outline:none;border-color:var(--leather)}.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}.exp{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem 1rem;margin-bottom:.4rem;display:flex;justify-content:space-between;align-items:center;transition:border-color .2s}.exp:hover{border-color:var(--leather)}.exp-left{flex:1}.exp-desc{font-size:.8rem}.exp-meta{font-size:.55rem;color:var(--cm);margin-top:.2rem;display:flex;gap:.5rem;flex-wrap:wrap}.exp-right{display:flex;align-items:center;gap:.5rem}.exp-amount{font-size:.9rem;font-weight:700}.cat-badge{font-size:.45rem;padding:.1rem .3rem;background:var(--bg3);color:var(--cd)}.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}.btn-sm{font-size:.55rem;padding:.2rem .4rem}.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:420px;max-width:92vw}.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust)}.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}.fr input,.fr select{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}.fr input:focus,.fr select:focus{outline:none;border-color:var(--leather)}.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> STEWARD</h1><button class="btn btn-p" onclick="openForm()">+ Add Expense</button></div>
<div class="main"><div class="stats" id="stats"></div><div class="toolbar"><input class="search" id="search" placeholder="Search expenses..." oninput="render()"><select class="filter-sel" id="cat-filter" onchange="render()"><option value="">All Categories</option></select></div><div id="list"></div></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',items=[],editId=null;
async function load(){var r=await fetch(A+'/expenses').then(function(r){return r.json()});items=r.expenses||[];renderStats();buildCatFilter();render();}
function renderStats(){var t=items.length,total=items.reduce(function(s,e){return s+(e.amount||0)},0),cats={};items.forEach(function(e){if(e.category)cats[e.category]=true});
document.getElementById('stats').innerHTML='<div class="st"><div class="st-v">'+t+'</div><div class="st-l">Expenses</div></div><div class="st"><div class="st-v">$'+(total/100).toLocaleString(undefined,{minimumFractionDigits:2})+'</div><div class="st-l">Total</div></div><div class="st"><div class="st-v">'+Object.keys(cats).length+'</div><div class="st-l">Categories</div></div>';}
function buildCatFilter(){var cats={};items.forEach(function(e){if(e.category)cats[e.category]=true});var sel=document.getElementById('cat-filter');var cur=sel.value;sel.innerHTML='<option value="">All Categories</option>';Object.keys(cats).sort().forEach(function(c){sel.innerHTML+='<option value="'+esc(c)+'"'+(cur===c?' selected':'')+'>'+esc(c)+'</option>';});}
function render(){var q=(document.getElementById('search').value||'').toLowerCase();var cf=document.getElementById('cat-filter').value;var f=items;
if(cf)f=f.filter(function(e){return e.category===cf});
if(q)f=f.filter(function(e){return(e.description||'').toLowerCase().includes(q)||(e.vendor||'').toLowerCase().includes(q)||(e.category||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('list').innerHTML='<div class="empty">No expenses recorded.</div>';return;}
var h='';f.forEach(function(e){
h+='<div class="exp"><div class="exp-left"><div class="exp-desc">'+esc(e.description)+'</div>';
h+='<div class="exp-meta">';
if(e.category)h+='<span class="cat-badge">'+esc(e.category)+'</span>';
if(e.vendor)h+='<span>'+esc(e.vendor)+'</span>';
if(e.date)h+='<span>'+e.date+'</span>';
if(e.status)h+='<span>'+esc(e.status)+'</span>';
h+='</div></div><div class="exp-right"><span class="exp-amount">$'+(e.amount/100).toFixed(2)+'</span>';
h+='<button class="btn btn-sm" onclick="openEdit(''+e.id+'')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(''+e.id+'')" style="color:var(--red)">&#10005;</button></div></div>';});
document.getElementById('list').innerHTML=h;}
async function del(id){if(!confirm('Delete?'))return;await fetch(A+'/expenses/'+id,{method:'DELETE'});load();}
function formHTML(exp){var i=exp||{description:'',amount:0,category:'',vendor:'',date:'',status:'pending'};var isEdit=!!exp;
var h='<h2>'+(isEdit?'EDIT':'ADD')+' EXPENSE</h2>';
h+='<div class="fr"><label>Description *</label><input id="f-desc" value="'+esc(i.description)+'"></div>';
h+='<div class="row2"><div class="fr"><label>Amount ($)</label><input id="f-amt" type="number" step="0.01" value="'+((i.amount||0)/100).toFixed(2)+'"></div>';
h+='<div class="fr"><label>Date</label><input id="f-date" type="date" value="'+esc(i.date)+'"></div></div>';
h+='<div class="row2"><div class="fr"><label>Category</label><input id="f-cat" value="'+esc(i.category)+'" placeholder="e.g. travel, office"></div>';
h+='<div class="fr"><label>Vendor</label><input id="f-vendor" value="'+esc(i.vendor)+'"></div></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add')+'</button></div>';return h;}
function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');}
function openEdit(id){var e=null;for(var j=0;j<items.length;j++){if(items[j].id===id){e=items[j];break;}}if(!e)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(e);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}
async function submit(){var desc=document.getElementById('f-desc').value.trim();if(!desc){alert('Description required');return;}
var body={description:desc,amount:Math.round(parseFloat(document.getElementById('f-amt').value||0)*100),date:document.getElementById('f-date').value,category:document.getElementById('f-cat').value.trim(),vendor:document.getElementById('f-vendor').value.trim()};
if(editId){await fetch(A+'/expenses/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/expenses',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}closeModal();load();}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});load();
</script></body></html>`
