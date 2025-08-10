import { backs, getCookie, setCookie, seats } from './utils.js';
import { RegistrationModal } from './registration.js';
import { SettingsModal } from './settings.js';
import { HelpModal } from './help.js';
import { AboutModal } from './about.js';
import { ShareModal } from './share.js';
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
  const [showShare,setShowShare]=React.useState(false);
  const [menuOpen,setMenuOpen]=React.useState(false);
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
      .then(r=>r.json()).then(d=>{setGameId(d.gameId);setShowShare(true);fetchState();});
  }
  function joinGame(gid){
    fetch('/api/game/join',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({gameId:gid,userId:user.id})})
      .then(()=>{setGameId(gid);setShowShare(true);fetchState();});
  }
  function finalize(){
    fetch('/api/game/finalize',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({gameId,userId:user.id})})
      .then(()=>{setShowShare(false);fetchState();});
  }
  function openAbout(){
    fetch('/api/version').then(r=>r.json()).then(d=>setVersion(d.version));
    setShowAbout(true);
  }

  return React.createElement('div',{className:'h-full'},[
    React.createElement('nav',{className:'fixed top-0 left-0 p-2'},[
      React.createElement('div',{className:'relative'},[
        React.createElement('button',{className:'px-2 py-1 bg-blue-200 text-black rounded',onClick:()=>setMenuOpen(!menuOpen)},'☰'),
        menuOpen && React.createElement('div',{className:'absolute mt-2 bg-white dark:bg-gray-800 rounded shadow'},[
          React.createElement('button',{className:'block w-full text-left px-4 py-2',onClick:()=>{setMenuOpen(false);startGame();}},'New Game'),
          React.createElement('button',{className:'block w-full text-left px-4 py-2',onClick:()=>{setMenuOpen(false);setShowSettings(true);}},'Settings'),
          React.createElement('button',{className:'block w-full text-left px-4 py-2',onClick:()=>{setMenuOpen(false);setShowHelp(true);}},'Help'),
          React.createElement('button',{className:'block w-full text-left px-4 py-2',onClick:()=>{setMenuOpen(false);openAbout();}},'About')
        ])
      ])
    ]),
    React.createElement('h1',{className:'text-3xl font-bold text-center mt-4'},'DONKEY'),
    React.createElement('div',{className:'p-4 mt-8 space-y-4'},[
      user.name && React.createElement('div',null,`Welcome, ${user.name}`),
      state && !state.hasStarted && state.requesterId==user.id &&
        React.createElement('button',{
          className:`px-3 py-1 bg-green-200 text-black rounded ${state.players.length>1?'':'opacity-50 cursor-not-allowed'}`,
          onClick:state.players.length>1?finalize:null
        },'Finalize Players and Deal'),
      state && !state.hasStarted && state.requesterId!=user.id && React.createElement('div',null,'Awaiting Creator to start the game...'),
      state && renderPlayers()
    ]),
    showSettings && React.createElement(SettingsModal,{theme,setTheme,backColor,setBackColor,onClose:()=>setShowSettings(false)}),
    showHelp && React.createElement(HelpModal,{onClose:()=>setShowHelp(false)}),
    showAbout && React.createElement(AboutModal,{version,onClose:()=>setShowAbout(false)}),
    showShare && React.createElement(ShareModal,{gameId,isRequester:state&&state.requesterId==user.id,playerCount:state?state.players.length:1,onFinalize:finalize,onClose:()=>setShowShare(false)}),
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

