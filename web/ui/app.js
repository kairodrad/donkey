import { backs, getCookie, setCookie, seats, sortCards } from './utils.js';
import { RegistrationModal } from './registration.js';
import { SettingsModal } from './settings.js';
import { HelpModal } from './help.js';
import { AboutModal } from './about.js';
import { ShareModal } from './share.js';
import { AbandonedModal } from './abandoned.js';
import { RenameModal } from './rename.js';
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
  if(gameIdParam) window.history.replaceState(null,'',window.location.pathname);
  const [user,setUser]=React.useState({id:getCookie('userId'),name:getCookie('userName')});
  const [gameId,setGameId]=React.useState(gameIdParam);
  const [state,setState]=React.useState(null);
  const [backColor,setBackColor]=React.useState(getCookie('cardBack')||'red');
  const [theme,setTheme]=React.useState(initialTheme);
  const [showReg,setShowReg]=React.useState(!(user.id && user.name));
  const [showSettings,setShowSettings]=React.useState(false);
  const [showHelp,setShowHelp]=React.useState(false);
  const [showAbout,setShowAbout]=React.useState(false);
  const [menuOpen,setMenuOpen]=React.useState(false);
  const [version,setVersion]=React.useState('');
  const [logs,setLogs]=React.useState([]);
  const [connected,setConnected]=React.useState(false);
  const [showLog,setShowLog]=React.useState(true);
  const [chat,setChat]=React.useState('');
  const [connKey,setConnKey]=React.useState(0);
  const [showAbandoned,setShowAbandoned]=React.useState(false);
  const [showRename,setShowRename]=React.useState(false);

  React.useEffect(()=>{setCookie('cardBack',backColor);},[backColor]);
  React.useEffect(()=>{setCookie('theme',theme);applyTheme(theme);},[theme]);
  React.useEffect(()=>{if(!(user.id && user.name)) setShowReg(true);},[user]);
  React.useEffect(()=>{if(gameIdParam && user.id){joinGame(gameIdParam);}},[user.id]);
  React.useEffect(()=>{if(gameId){fetchState();fetchLogs();}},[gameId]);
  React.useEffect(()=>{if(state && state.isAbandoned){setShowAbandoned(true);}},[state]);
  React.useEffect(()=>{
    if(!gameId||!user.id) return;
    const es=new EventSource(`/api/game/${gameId}/stream/${user.id}`);
    es.onopen=()=>setConnected(true);
    es.onerror=()=>{setConnected(false);es.close();setTimeout(()=>setConnKey(k=>k+1),2000);};
    es.onmessage=e=>{
      const data=JSON.parse(e.data);
      if(data.type==='state') fetchState();
      if(data.type==='log') setLogs(l=>[data.log,...l]);
    };
    return ()=>{es.close();setConnected(false);};
  },[gameId,user.id,connKey]);

  function fetchState(){
    fetch(`/api/game/${gameId}/state/${user.id}`).then(r=>r.json()).then(setState);
  }
  function fetchLogs(){
    fetch(`/api/game/${gameId}/logs`).then(r=>r.json()).then(setLogs);
  }
  function register(name){
    fetch('/api/register',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name})})
      .then(r=>r.json()).then(d=>{setUser(d);setCookie('userId',d.id);setCookie('userName',d.name);setShowReg(false);if(gameIdParam) joinGame(gameIdParam);});
  }
  function startGame(){
    fetch('/api/game/start',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({requesterId:user.id})})
      .then(r=>r.json()).then(d=>{setGameId(d.gameId);});
  }
  function joinGame(gid){
    fetch('/api/game/join',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({gameId:gid,userId:user.id})})
      .then(()=>{setGameId(gid);window.history.replaceState(null,'','/');});
  }
  function finalize(){
    fetch('/api/game/finalize',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({gameId,userId:user.id})});
  }
  function abandon(){
    fetch('/api/game/abandon',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({gameId,userId:user.id})})
      .then(()=>{setGameId(null);setState(null);setLogs([]);window.history.replaceState(null,'','/');});
  }
  function rename(newName){
    fetch(`/api/user/${user.id}/rename`,{
      method:'POST',
      headers:{'Content-Type':'application/json'},
      body:JSON.stringify({gameId,name:newName})
    })
      .then(r=>{if(!r.ok) throw new Error('rename'); return r.json();})
      .then(d=>{if(d && d.name){setUser(u=>({...u,name:d.name}));setCookie('userName',d.name);}})
      .catch(()=>{});
  }
  function openAbout(){
    fetch('/api/version').then(r=>r.json()).then(d=>setVersion(d.version));
    setShowAbout(true);
  }
  function sendChat(){
    const msg=chat.trim();
    if(!msg) return;
    fetch('/api/game/chat',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({gameId,userId:user.id,message:msg})});
    setChat('');
  }

  const showShare = gameId && state && !state.hasStarted && !state.isAbandoned;
  const isRequester = state && state.requesterId==user.id;

  return React.createElement('div',{className:'h-full flex flex-col items-center'},[
    React.createElement('nav',{className:'fixed top-0 left-0 p-2'},[
      React.createElement('div',{className:'relative'},[
        React.createElement('button',{className:'px-2 py-1 bg-blue-200 dark:bg-blue-700 text-black dark:text-white rounded',onClick:()=>setMenuOpen(!menuOpen)},'☰'),
        menuOpen && React.createElement('div',{className:'absolute mt-2 bg-white dark:bg-gray-800 text-black dark:text-white rounded shadow'},[
          (!gameId && React.createElement('button',{className:'block w-full text-left px-4 py-2 whitespace-nowrap',onClick:()=>{setMenuOpen(false);startGame();}},'New Game')),
          (gameId && isRequester && React.createElement('button',{className:'block w-full text-left px-4 py-2 whitespace-nowrap',onClick:()=>{setMenuOpen(false);abandon();}},'Abandon Game')),
          (gameId && !isRequester && React.createElement('button',{className:'block w-full text-left px-4 py-2 whitespace-nowrap opacity-50 cursor-not-allowed',disabled:true},'New Game')),
          React.createElement('button',{className:'block w-full text-left px-4 py-2 whitespace-nowrap',onClick:()=>{setMenuOpen(false);setShowRename(true);}},'Rename Player'),
          React.createElement('button',{className:'block w-full text-left px-4 py-2 whitespace-nowrap',onClick:()=>{setMenuOpen(false);setShowSettings(true);}},'Settings'),
          React.createElement('button',{className:'block w-full text-left px-4 py-2 whitespace-nowrap',onClick:()=>{setMenuOpen(false);setShowHelp(true);}},'Help'),
          React.createElement('button',{className:'block w-full text-left px-4 py-2 whitespace-nowrap',onClick:()=>{setMenuOpen(false);openAbout();}},'About')
        ])
      ])
    ]),
    React.createElement('h1',{className:'text-3xl font-bold text-center mt-4'},'DONKEY'),
    React.createElement('div',{className:'p-4 mt-8 space-y-4 flex flex-col items-center'},[
      state && !state.hasStarted && isRequester &&
        React.createElement('button',{
          className:`px-3 py-1 bg-green-200 dark:bg-green-700 text-black dark:text-white rounded ${state.players.length>1?'':'opacity-50 cursor-not-allowed'}`,
          onClick:state.players.length>1?finalize:null
        },'Finalize Players and Deal'),
      state && !state.hasStarted && !isRequester && React.createElement('div',null,'Awaiting Creator to start the game...'),
      state && renderPlayers()
    ]),
    showSettings && React.createElement(SettingsModal,{theme,setTheme,backColor,setBackColor,user,onRename:rename,onClose:()=>setShowSettings(false)}),
    showRename && React.createElement(RenameModal,{currentName:user.name,onSubmit:rename,onClose:()=>setShowRename(false)}),
    showHelp && React.createElement(HelpModal,{onClose:()=>setShowHelp(false)}),
    showAbout && React.createElement(AboutModal,{version,onClose:()=>setShowAbout(false)}),
    showShare && React.createElement(ShareModal,{gameId,isRequester,playerCount:state?state.players.length:1,onFinalize:finalize}),
    showReg && React.createElement(RegistrationModal,{onSubmit:register}),
    showAbandoned && React.createElement(AbandonedModal,{onClose:closeAbandoned}),
    React.createElement('div',{className:`fixed bottom-0 right-0 m-2 text-xs px-2 py-1 rounded ${connected?'bg-green-500':'bg-red-500'} text-white`},
      user.name ? (connected?`Connected to game as ${user.name}`:`Disconnected as ${user.name}`) : (connected?'Connected':'Disconnected')
    ),
    gameId && React.createElement('div',{className:'fixed bottom-0 left-0 m-2 w-64'},[
      React.createElement('div',{className:'bg-white dark:bg-gray-800 text-black dark:text-white border rounded'},[
        React.createElement('div',{className:'flex justify-between items-center px-2 py-1 border-b'},[
          React.createElement('span',null,'Session Updates'),
          React.createElement('button',{onClick:()=>setShowLog(!showLog)},showLog?'-':'+')
        ]),
        showLog?[
          React.createElement('div',{className:'h-24 overflow-y-auto px-2 text-sm'},logs.map(l=>React.createElement('div',{key:l.id},l.message))),
          React.createElement('div',{className:'flex border-t'},[
            React.createElement('input',{className:'flex-grow p-1 bg-white text-black dark:bg-gray-700 dark:text-white',value:chat,maxLength:128,onChange:e=>setChat(e.target.value),onKeyDown:e=>{if(e.key==='Enter')sendChat();}}),
            React.createElement('button',{className:'px-2',onClick:sendChat},'Send')
          ])
        ]:
          React.createElement('div',{className:'px-2 py-1 text-sm truncate'},logs.length?logs[0].message:'')
      ])
    ])
  ]);

  function renderPlayers(){
    const playersWithSeats=seats(state.players,user.id);
    return React.createElement('div',{className:'relative w-[400px] h-[400px] mx-auto mt-4'},
      playersWithSeats.map(p=>React.createElement('div',{key:p.id,className:'absolute flex flex-col items-center',style:{top:p.seat.top,left:p.seat.left,transform:'translate(-50%,-50%)'}},[
        state.hasStarted?renderCards(p):null,
        React.createElement('div',{className:'mt-1 font-semibold'},p.name)
      ]))
    );
  }

  function renderCards(p){
    if(p.id==user.id){
      const cards=sortCards(p.cards||[]);
      return React.createElement('div',{className:'flex justify-center items-end h-28'},
        cards.map((c,i)=>React.createElement('img',{key:i,src:`/assets/${c}.png`,className:'w-16 h-24 -ml-4 first:ml-0 relative hover:z-10 hover:-translate-y-2 hover:ml-0 transition-all'}))
      );
    }
    return React.createElement('div',{className:'flex justify-center items-end'},
      Array(p.cardCount).fill(0).map((_,i)=>React.createElement('img',{key:i,src:`/assets/${backColor}_back.png`,className:'w-8 h-12 -ml-2 first:ml-0'}))
    );
  }

  function closeAbandoned(){
    setShowAbandoned(false);setGameId(null);setState(null);setLogs([]);window.history.replaceState(null,'','/');
  }
}

ReactDOM.render(React.createElement(App), document.getElementById('root'));
