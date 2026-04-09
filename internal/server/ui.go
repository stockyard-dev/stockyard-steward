package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Steward</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:960px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center}
.st-v{font-size:1.3rem;font-weight:700;color:var(--gold)}
.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.count-label{font-size:.6rem;color:var(--cm);margin-bottom:.5rem}
.item{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}
.item:hover{border-color:var(--leather)}
.item-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.item-title{font-size:.85rem;font-weight:700;flex:1}
.item-amount{font-size:1rem;font-weight:700;color:var(--gold);white-space:nowrap}
.item-meta{font-size:.55rem;color:var(--cm);margin-top:.4rem;display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
.item-meta-sep{color:var(--bg3)}
.item-actions{display:flex;gap:.3rem;flex-shrink:0;margin-left:.5rem}
.item-extra{font-size:.58rem;color:var(--cd);margin-top:.4rem;padding-top:.35rem;border-top:1px dashed var(--bg3);display:flex;flex-direction:column;gap:.15rem}
.item-extra-row{display:flex;gap:.4rem}
.item-extra-label{color:var(--cm);text-transform:uppercase;letter-spacing:.5px;min-width:90px}
.item-extra-val{color:var(--cream)}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm)}
.badge.paid{border-color:var(--green);color:var(--green)}
.badge.pending{border-color:var(--gold);color:var(--gold)}
.badge.overdue{border-color:var(--red);color:var(--red)}
.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s;font-family:var(--mono)}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:480px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.fr-section{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.fr-section-label{font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.row2{grid-template-columns:1fr}.toolbar{flex-direction:column;align-items:stretch}.search{min-width:100%}.filter-sel{width:100%}.trial-bar{flex-direction:column;align-items:stretch}.trial-bar input.key-input{width:100%}}
.trial-bar{display:none;background:linear-gradient(90deg,#3a2419,#2e1c14);border-bottom:2px solid var(--rust);padding:.7rem 1.5rem;font-family:var(--mono);font-size:.68rem;color:var(--cream);align-items:center;gap:1rem;flex-wrap:wrap}
.trial-bar.show{display:flex}
.trial-bar-msg{flex:1;min-width:240px;line-height:1.5}
.trial-bar-msg strong{color:var(--rust);text-transform:uppercase;letter-spacing:1px;font-size:.6rem;display:block;margin-bottom:.15rem}
.trial-bar-actions{display:flex;gap:.5rem;align-items:center;flex-wrap:wrap}
.trial-bar a.btn-trial{background:var(--rust);color:#fff;padding:.4rem .8rem;text-decoration:none;font-size:.65rem;text-transform:uppercase;letter-spacing:1px;font-weight:700;border:1px solid var(--rust);transition:all .2s}
.trial-bar a.btn-trial:hover{background:#f08545;border-color:#f08545}
.trial-bar-divider{color:var(--cm);font-size:.6rem}
.trial-bar input.key-input{padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.6rem;width:200px}
.trial-bar input.key-input:focus{outline:none;border-color:var(--rust)}
.trial-bar button.btn-activate{padding:.4rem .7rem;background:var(--bg2);color:var(--cream);border:1px solid var(--leather);font-family:var(--mono);font-size:.6rem;cursor:pointer;text-transform:uppercase;letter-spacing:1px}
.trial-bar button.btn-activate:hover{background:var(--bg3)}
.trial-bar button.btn-activate:disabled{opacity:.5;cursor:wait}
.trial-msg{font-size:.6rem;color:var(--cm);margin-left:.5rem}
.trial-msg.error{color:#e74c3c}
.trial-msg.success{color:#4ade80}
.btn-disabled-trial{opacity:.45;cursor:not-allowed!important}
</style>
</head>
<body>

<div class="trial-bar" id="trial-bar">
<div class="trial-bar-msg">
<strong>Trial Required</strong>
You can view your existing expenses, but adding, editing, or deleting is locked until you start a 14-day free trial.
</div>
<div class="trial-bar-actions">
<a class="btn-trial" href="https://stockyard.dev/" target="_blank" rel="noopener">Start 14-Day Trial</a>
<span class="trial-bar-divider">or</span>
<input type="text" class="key-input" id="trial-key-input" placeholder="SY-..." autocomplete="off" spellcheck="false">
<button class="btn-activate" id="trial-activate-btn" onclick="activateLicense()">Activate</button>
<span class="trial-msg" id="trial-msg"></span>
</div>
</div>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> STEWARD</h1>
<button class="btn btn-p" onclick="openForm()">+ Add Expense</button>
</div>

<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search description, vendor, category..." oninput="render()">
<select class="filter-sel" id="category-filter" onchange="render()">
<option value="">All Categories</option>
</select>
<select class="filter-sel" id="status-filter" onchange="render()">
<option value="">All Status</option>
<option value="pending">Pending</option>
<option value="paid">Paid</option>
<option value="overdue">Overdue</option>
</select>
</div>
<div class="count-label" id="count"></div>
<div id="list"></div>
</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var RESOURCE='expenses';

// Field defs drive the form, the items, and the submit body.
// Categories starts empty and is populated from /api/config.categories
// (or falls back to a generic default list if no config is present).
var defaultCategories=['Office','Travel','Meals','Software','Equipment','Marketing','Utilities','Other'];

var fields=[
{name:'description',label:'Description',type:'textarea',required:true,placeholder:'What was this for?'},
{name:'amount',label:'Amount',type:'money',required:true,placeholder:'0.00'},
{name:'category',label:'Category',type:'select',options:defaultCategories},
{name:'vendor',label:'Vendor',type:'text',placeholder:'Who you paid'},
{name:'date',label:'Date',type:'date'},
{name:'payment_method',label:'Payment Method',type:'select',options:['Cash','Credit Card','Debit Card','Check','Bank Transfer','Other']},
{name:'receipt',label:'Receipt #',type:'text'},
{name:'status',label:'Status',type:'select',options:['pending','paid','overdue']}
];

var items=[],editId=null;

// ─── Money helpers ────────────────────────────────────────────────
// Amounts are stored as integer cents to avoid float math errors.

function fmtMoney(cents){
if(cents===null||cents===undefined||cents==='')return'$0.00';
var n=parseInt(cents,10);
if(isNaN(n))return'$0.00';
var neg=n<0;
n=Math.abs(n);
var dollars=Math.floor(n/100);
var rem=n%100;
var s='$'+dollars.toLocaleString()+'.'+(rem<10?'0':'')+rem;
return neg?'-'+s:s;
}

function parseMoney(str){
if(!str)return 0;
// Strip everything except digits, dot, and leading minus
var s=String(str).replace(/[^0-9.\-]/g,'');
var f=parseFloat(s);
if(isNaN(f))return 0;
return Math.round(f*100);
}

function fmtDate(s){
if(!s)return'';
try{return new Date(s).toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'})}catch(e){return s}
}

// ─── Loading and rendering ────────────────────────────────────────

async function load(){
try{
var r=await fetch(A+'/'+RESOURCE).then(function(r){return r.json()});
var list=r[RESOURCE]||[];
try{
var extras=await fetch(A+'/extras/'+RESOURCE).then(function(r){return r.json()});
list.forEach(function(it){
var ex=extras[it.id];
if(!ex)return;
Object.keys(ex).forEach(function(k){if(it[k]===undefined)it[k]=ex[k]});
});
}catch(e){}
items=list;
}catch(e){
console.error('load failed',e);
items=[];
}
populateCategoryFilter();
renderStats();
render();
}

function populateCategoryFilter(){
var sel=document.getElementById('category-filter');
if(!sel)return;
var current=sel.value;
// Collect categories from items + the configured defaults
var seen={};
var cats=[];
var catField=fieldByName('category');
if(catField&&catField.options){catField.options.forEach(function(c){if(!seen[c]){seen[c]=true;cats.push(c)}})}
items.forEach(function(i){if(i.category&&!seen[i.category]){seen[i.category]=true;cats.push(i.category)}});
sel.innerHTML='<option value="">All Categories</option>'+cats.map(function(c){return'<option value="'+esc(c)+'"'+(c===current?' selected':'')+'>'+esc(c)+'</option>'}).join('');
}

function renderStats(){
var total=items.length;
var totalCents=0;
items.forEach(function(i){totalCents+=parseInt(i.amount||0,10)});
// "This month" = current calendar month based on the date field (not created_at)
var nowD=new Date();
var ym=nowD.getFullYear()+'-'+String(nowD.getMonth()+1).padStart(2,'0');
var thisMonthCents=0;
items.forEach(function(i){
if(i.date&&String(i.date).startsWith(ym)){thisMonthCents+=parseInt(i.amount||0,10)}
});
document.getElementById('stats').innerHTML=
'<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Expenses</div></div>'+
'<div class="st"><div class="st-v">'+fmtMoney(totalCents)+'</div><div class="st-l">Total</div></div>'+
'<div class="st"><div class="st-v">'+fmtMoney(thisMonthCents)+'</div><div class="st-l">This Month</div></div>';
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var cf=document.getElementById('category-filter').value;
var sf=document.getElementById('status-filter').value;
var f=items;
if(cf)f=f.filter(function(i){return i.category===cf});
if(sf)f=f.filter(function(i){return i.status===sf});
if(q)f=f.filter(function(i){
return(i.description||'').toLowerCase().includes(q)||
       (i.vendor||'').toLowerCase().includes(q)||
       (i.category||'').toLowerCase().includes(q);
});
document.getElementById('count').textContent=f.length+' expense'+(f.length!==1?'s':'');
if(!f.length){
var msg=window._emptyMsg||'No expenses found.';
document.getElementById('list').innerHTML='<div class="empty">'+esc(msg)+'</div>';
return;
}
var h='';
f.forEach(function(i){h+=itemHTML(i)});
document.getElementById('list').innerHTML=h;
}

function itemHTML(i){
var h='<div class="item"><div class="item-top">';
h+='<div class="item-title">'+esc(i.description)+'</div>';
h+='<div class="item-amount">'+fmtMoney(i.amount)+'</div>';
if(!window._trialRequired){
h+='<div class="item-actions">';
h+='<button class="btn btn-sm" onclick="openEdit(\''+i.id+'\')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(\''+i.id+'\')" style="color:var(--red)">&#10005;</button>';
h+='</div>';
}
h+='</div>';

h+='<div class="item-meta">';
var parts=[];
if(i.category)parts.push('<span>'+esc(i.category)+'</span>');
if(i.vendor)parts.push('<span>'+esc(i.vendor)+'</span>');
if(i.date)parts.push('<span>'+esc(fmtDate(i.date))+'</span>');
if(i.payment_method)parts.push('<span>'+esc(i.payment_method)+'</span>');
if(i.receipt)parts.push('<span>#'+esc(i.receipt)+'</span>');
h+=parts.join('<span class="item-meta-sep">·</span>');
if(i.status)h+=' <span class="badge '+esc(i.status)+'">'+esc(i.status)+'</span>';
h+='</div>';

// Custom fields from personalization in their own block
var customRows='';
fields.forEach(function(f){
if(!f.isCustom)return;
var v=i[f.name];
if(v===undefined||v===null||v==='')return;
customRows+='<div class="item-extra-row">';
customRows+='<span class="item-extra-label">'+esc(f.label)+'</span>';
customRows+='<span class="item-extra-val">'+esc(String(v))+'</span>';
customRows+='</div>';
});
if(customRows)h+='<div class="item-extra">'+customRows+'</div>';

h+='</div>';
return h;
}

// ─── Form ─────────────────────────────────────────────────────────

function fieldByName(n){
for(var i=0;i<fields.length;i++)if(fields[i].name===n)return fields[i];
return null;
}

function fieldHTML(f,value){
var v=value;
if(v===undefined||v===null)v='';
var req=f.required?' *':'';
var ph='';
if(f.placeholder)ph=' placeholder="'+esc(f.placeholder)+'"';
else if(f.name==='description'&&window._placeholderName)ph=' placeholder="'+esc(window._placeholderName)+'"';

var h='<div class="fr"><label>'+esc(f.label)+req+'</label>';

if(f.type==='select'){
h+='<select id="f-'+f.name+'">';
if(!f.required)h+='<option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=(String(v)===String(o))?' selected':'';
var disp=(typeof o==='string')?(o.charAt(0).toUpperCase()+o.slice(1)):String(o);
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(disp)+'</option>';
});
h+='</select>';
}else if(f.type==='textarea'){
h+='<textarea id="f-'+f.name+'" rows="2"'+ph+'>'+esc(String(v))+'</textarea>';
}else if(f.type==='checkbox'){
h+='<input type="checkbox" id="f-'+f.name+'"'+(v?' checked':'')+' style="width:auto">';
}else if(f.type==='money'){
// Display dollars/cents in the input even though we store cents.
var displayVal=v?(parseInt(v,10)/100).toFixed(2):'';
h+='<input type="text" id="f-'+f.name+'" value="'+esc(displayVal)+'"'+ph+' inputmode="decimal">';
}else{
var inputType=f.type||'text';
h+='<input type="'+esc(inputType)+'" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}

h+='</div>';
return h;
}

function formHTML(item){
var i=item||{};
var isEdit=!!item;
var h='<h2>'+(isEdit?'EDIT EXPENSE':'NEW EXPENSE')+'</h2>';

// Render description first (full width), then natural pairs
h+=fieldHTML(fieldByName('description'),i.description);

// Pair amount + date, category + vendor, payment_method + status
var pairs=[['amount','date'],['category','vendor'],['payment_method','status']];
pairs.forEach(function(p){
var a=fieldByName(p[0]);
var b=fieldByName(p[1]);
if(a&&b){
h+='<div class="row2">'+fieldHTML(a,i[p[0]])+fieldHTML(b,i[p[1]])+'</div>';
}
});
h+=fieldHTML(fieldByName('receipt'),i.receipt);

// Custom fields injected by personalization get their own labeled section
var customFields=fields.filter(function(f){return f.isCustom});
if(customFields.length){
var sectionLabel=window._customSectionLabel||'Additional Details';
h+='<div class="fr-section"><div class="fr-section-label">'+esc(sectionLabel)+'</div>';
customFields.forEach(function(f){h+=fieldHTML(f,i[f.name])});
h+='</div>';
}

h+='<div class="acts">';
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add Expense')+'</button>';
h+='</div>';
return h;
}

function openForm(){
editId=null;
document.getElementById('mdl').innerHTML=formHTML();
document.getElementById('mbg').classList.add('open');
var d=document.getElementById('f-description');
if(d)d.focus();
}

function openEdit(id){
var x=null;
for(var j=0;j<items.length;j++){if(items[j].id===id){x=items[j];break}}
if(!x)return;
editId=id;
document.getElementById('mdl').innerHTML=formHTML(x);
document.getElementById('mbg').classList.add('open');
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
editId=null;
}

// ─── Submit ───────────────────────────────────────────────────────

async function submit(){
var descEl=document.getElementById('f-description');
if(!descEl||!descEl.value.trim()){alert('Description is required');return}

var body={};
var extras={};
fields.forEach(function(f){
var el=document.getElementById('f-'+f.name);
if(!el)return;
var val;
if(f.type==='checkbox')val=el.checked;
else if(f.type==='money')val=parseMoney(el.value);
else if(f.type==='number')val=parseFloat(el.value)||0;
else val=el.value.trim();
if(f.isCustom)extras[f.name]=val;
else body[f.name]=val;
});

var savedId=editId;
try{
if(editId){
var r1=await fetch(A+'/'+RESOURCE+'/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/'+RESOURCE,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Save failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/'+RESOURCE+'/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){
alert('Network error: '+e.message);
return;
}

closeModal();
load();
}

async function del(id){
if(!confirm('Delete this expense?'))return;
await fetch(A+'/'+RESOURCE+'/'+id,{method:'DELETE'});
load();
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;
if(cfg.placeholder_name)window._placeholderName=cfg.placeholder_name;
if(cfg.primary_label)window._customSectionLabel=cfg.primary_label+' Details';

// Categories from config replace the default category options.
// This is the killer feature for steward — a brewery gets brewery
// categories, a therapy practice gets therapy categories, etc.
if(Array.isArray(cfg.categories)&&cfg.categories.length){
var catField=fieldByName('category');
if(catField)catField.options=cfg.categories.slice();
}

if(Array.isArray(cfg.custom_fields)){
cfg.custom_fields.forEach(function(cf){
if(!cf||!cf.name||!cf.label)return;
if(fieldByName(cf.name))return;
fields.push({
name:cf.name,
label:cf.label,
type:cf.type||'text',
options:cf.options||[],
isCustom:true
});
});
}
}).catch(function(){
// No config — use defaults
}).finally(function(){
checkTrialState();
load();
});
})();

// ─── trial-required license gating ───
window._trialRequired=false;

async function checkTrialState(){
try{
var resp=await fetch(A+'/tier');
if(!resp.ok)return;
var data=await resp.json();
window._trialRequired=!!data.trial_required;
if(window._trialRequired){
document.getElementById('trial-bar').classList.add('show');
disableWriteControls();
if(typeof render==='function')render();
}else{
document.getElementById('trial-bar').classList.remove('show');
}
}catch(e){}
}

function disableWriteControls(){
var buttons=document.querySelectorAll('.hdr .btn, .hdr .btn-p');
buttons.forEach(function(b){
var t=b.textContent||'';
if(t.indexOf('Add')!==-1||t.indexOf('New')!==-1||t.indexOf('Log')!==-1){
b.classList.add('btn-disabled-trial');
b.title='Locked: trial required';
b.onclick=function(e){
e.preventDefault();
showTrialNudge();
return false;
};
}
});
}

function showTrialNudge(){
var input=document.getElementById('trial-key-input');
if(input){
input.focus();
input.style.borderColor='var(--rust)';
setTimeout(function(){if(input)input.style.borderColor=''},1500);
}
}

async function activateLicense(){
var input=document.getElementById('trial-key-input');
var btn=document.getElementById('trial-activate-btn');
var msg=document.getElementById('trial-msg');
if(!input||!btn||!msg)return;
var key=(input.value||'').trim();
if(!key){
msg.className='trial-msg error';
msg.textContent='Paste your license key first';
input.focus();
return;
}
btn.disabled=true;
msg.className='trial-msg';
msg.textContent='Activating...';
try{
var resp=await fetch(A+'/license/activate',{
method:'POST',
headers:{'Content-Type':'application/json'},
body:JSON.stringify({license_key:key})
});
var data=await resp.json();
if(!resp.ok){
msg.className='trial-msg error';
msg.textContent=data.error||'Activation failed';
btn.disabled=false;
return;
}
msg.className='trial-msg success';
msg.textContent='Activated. Reloading...';
setTimeout(function(){location.reload()},800);
}catch(e){
msg.className='trial-msg error';
msg.textContent='Network error: '+e.message;
btn.disabled=false;
}
}

document.addEventListener('DOMContentLoaded',function(){
var input=document.getElementById('trial-key-input');
if(input){
input.addEventListener('keydown',function(e){
if(e.key==='Enter')activateLicense();
});
}
});
</script>
</body>
</html>`
