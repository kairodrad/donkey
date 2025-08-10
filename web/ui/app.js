import { backs, getCookie, setCookie, seats } from './utils.js';
import { RegistrationModal } from './registration.js';
import { SettingsModal } from './settings.js';
import { HelpModal } from './help.js';
import { AboutModal } from './about.js';
const React = window.React;
const ReactDOM = window.ReactDOM;

function applyTheme(t){
  const actual=t==='system'? (window.matchMedia('(prefers-color-scheme: dark)').matches?'dark':'light') : t;
  document.documentElement.classList.toggle('dark', actual==='dark');
}
const initialTheme=getCookie('theme')||'system';
applyTheme(initialTheme);

function App(){
  const params=new URLSearchParams(window.location.search);
  const gameIdParam=params.get('gameId');
  const [user,setUser]=React.useState({id:getCookie('userId'),name:getCookie('userName')});
  const [gameId,setGameId]=React.useState(gameIdParam);
  const [state,setState]=React.useState(null);
  const [backColor,setBackColor]=React.useState(getCookie('cardBack')||'red');
  const [theme,setTheme]=React.useState(initialTheme);
  const [showReg,setShowReg]=React.useState(!(user.id && user.name));
  const [showSettings,setShowSettings]=React.useState(false);
  const [showHelp,setShowHelp]=React.useState(false);
  const [showAbout,setShowAbout]=React.useState(false);
  const [version,setVersion]=React.useState('');

  React.useEffect(()=>{setCookie('cardBack',backColor);},[backColor]);
  React.useEffect(()=>{setCookie('theme',theme);applyTheme(theme);},[theme]);
  React.useEffect(()=>{if(!(user.id && user.name)) setShowReg(true);},[user]);
  React.useEffect(()=>{
    if(!gameId||!user.id) return;
    fetchState();
    const t=setInterval(fetchState,2000);
    return ()=>clearInterval(t);
  },[gameId,user]);
  React.useEffect(()=>{if(gameIdParam && user.id){joinGame(gameIdParam);}},[user.id]);

  function fetchState(){
    fetch(`/api/game/state?gameId=${gameId}&userId=${user.id}`).then(r=>r.json()).then(setState);
  }
  function register(name){
    fetch('/api/register',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name})})
      .then(r=>r.json()).then(d=>{setUser(d);setCookie('userId',d.id);setCookie('userName',d.name);setShowReg(false);if(gameIdParam) joinGame(gameIdParam);});
  }
  function startGame(){
    fetch('/api/game/start',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({requesterId:user.id})})
      .then(r=>r.json()).then(d=>{setGameId(d.gameId);window.history.replaceState({},'',`?gameId=${d.gameId}`);fetchState();});
  }
  function joinGame(gid){
    fetch('/api/game/join',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({gameId:gid,userId:user.id})})
      .then(()=>{setGameId(gid);fetchState();});
  }
  function finalize(){
    fetch('/api/game/finalize',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({gameId,userId:user.id})})
      .then(()=>fetchState());
  }
  function copyLink(){navigator.clipboard.writeText(window.location.href);}
  function openAbout(){
    fetch('/api/version').then(r=>r.json()).then(d=>setVersion(d.version));
    setShowAbout(true);
  }

  return React.createElement('div',{className:'h-full'},[
    React.createElement('nav',{className:'fixed top-0 left-0 p-2 space-x-2'},[
      React.createElement('button',{className:'px-2 py-1 bg-blue-200 text-black rounded',onClick:startGame},'New Game'),
      React.createElement('button',{className:'px-2 py-1 bg-blue-200 text-black rounded',onClick:()=>setShowSettings(true)},'Settings'),
      React.createElement('button',{className:'px-2 py-1 bg-blue-200 text-black rounded',onClick:()=>setShowHelp(true)},'Help'),
      React.createElement('button',{className:'px-2 py-1 bg-blue-200 text-black rounded',onClick:openAbout},'About')
    ]),
    React.createElement('h1',{className:'text-3xl font-bold text-center mt-4'},'DONKEY'),
    React.createElement('div',{className:'p-4 mt-8 space-y-4'},[
      user.name && React.createElement('div',null,`Welcome, ${user.name}`),
      gameId && state && !state.hasStarted && React.createElement('div',{className:'space-x-2'},[
        React.createElement('button',{className:'px-3 py-1 bg-blue-200 text-black rounded',onClick:copyLink},'Copy Game URL'),
        state.requesterId==user.id && React.createElement('button',{className:'px-3 py-1 bg-green-200 text-black rounded',onClick:finalize},'Finalize Players and Deal')
      ]),
      state && renderPlayers()
    ]),
    showSettings && React.createElement(SettingsModal,{theme,setTheme,backColor,setBackColor,onClose:()=>setShowSettings(false)}),
    showHelp && React.createElement(HelpModal,{onClose:()=>setShowHelp(false)}),
    showAbout && React.createElement(AboutModal,{version,onClose:()=>setShowAbout(false)}),
    showReg && React.createElement(RegistrationModal,{onSubmit:register})
  ]);

  function renderPlayers(){
    const playersWithSeats=seats(state.players);
    return React.createElement('div',{className:'relative w-[400px] h-[400px] mx-auto mt-4'},
      playersWithSeats.map(p=>React.createElement('div',{key:p.id,style:{position:'absolute',top:p.seat.top,left:p.seat.left,textAlign:'center'}},[
        React.createElement('div',{className:'font-semibold'},p.name),
        state.hasStarted?renderCards(p):null
      ]))
    );
  }

  function renderCards(p){
    const cards=p.id==user.id?p.cards:Array(p.cardCount).fill(`${backColor}_back`);
    return React.createElement('div',{className:'flex'},
      cards.map((c,i)=>React.createElement('img',{key:i,src:`/assets/${c}.png`,className:'w-8 h-12 -ml-4 first:ml-0'}))
    );
  }
}

ReactDOM.render(React.createElement(App), document.getElementById('root'));

