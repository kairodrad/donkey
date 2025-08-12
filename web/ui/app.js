import { getCookie, setCookie, seats, sortCards } from './utils.js';
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
  const gameIdRef=React.useRef(null);
  if(gameIdRef.current===null){
    const params=new URLSearchParams(window.location.search);
    gameIdRef.current=params.get('gameId');
    if(gameIdRef.current) window.history.replaceState(null,'',window.location.pathname);
  }
  const [user,setUser]=React.useState({id:null,name:null});
  const [gameId,setGameId]=React.useState(null);
  const [state,setState]=React.useState(null);
  const [backColor,setBackColor]=React.useState(getCookie('cardBack')||'red');
  const [theme,setTheme]=React.useState(initialTheme);
  const [showReg,setShowReg]=React.useState(false);
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
  const [selectedCard,setSelectedCard]=React.useState(null);

  React.useEffect(()=>{setCookie('cardBack',backColor);},[backColor]);
  React.useEffect(()=>{setCookie('theme',theme);applyTheme(theme);},[theme]);
  React.useEffect(()=>{
    const id=getCookie('userId');
    if(!id){setShowReg(true);return;}
    fetch(`/api/user/${id}`)
      .then(r=>{if(!r.ok) throw new Error('not found'); return r.json();})
      .then(d=>{setUser({id:d.id,name:d.name});setCookie('userId',d.id);setCookie('userName',d.name);})
      .catch(()=>{setShowReg(true);});
  },[]);
  React.useEffect(()=>{if(gameId && user.id){fetchState();fetchLogs();}},[gameId,user.id]);
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
    const body={name};
    fetch('/api/register',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)})
      .then(r=>r.json()).then(d=>{
        setUser(d);
        setCookie('userId',d.id);
        setCookie('userName',d.name);
        setShowReg(false);
        if(gameIdRef.current) joinGame(gameIdRef.current,d.id);
      });
  }
  function startGame(){
    if(!user.id){setShowReg(true);return;}
    fetch('/api/game/start',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({requesterId:user.id})})
      .then(r=>{if(!r.ok) throw new Error('start'); return r.json();})
      .then(d=>{setGameId(d.gameId);})
      .catch(()=>{});
  }
  function joinGame(gid,uid){
    const id=uid||user.id;
    if(!id) return;
    fetch('/api/game/join',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({gameId:gid,userId:id})})
      .then(r=>{if(!r.ok) throw new Error('join');})
      .then(()=>{setGameId(gid);window.history.replaceState(null,'','/');})
      .catch(()=>{});
  }
  React.useEffect(()=>{
    if(user.id && gameIdRef.current && !gameId){
      joinGame(gameIdRef.current,user.id);
    }
  },[user.id,gameId]);
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
  const actualTheme = theme==='system' ? (window.matchMedia('(prefers-color-scheme: dark)').matches?'dark':'light') : theme;

  return React.createElement('div',{className:'h-full flex flex-col items-center relative bg-white dark:bg-black'},[
    React.createElement('img',{src:`/assets/donkey-background-${actualTheme}.png`,className:'absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 max-w-full max-h-full opacity-20 pointer-events-none select-none'}),
    React.createElement('div',{className:'fixed top-0 left-0 w-full h-20 flex items-center justify-center bg-white dark:bg-black z-30'},[
      React.createElement('div',{className:'absolute left-2 h-full flex items-center relative'},[
        React.createElement('button',{className:'h-full aspect-square bg-amber-200 dark:bg-amber-700 text-black dark:text-white rounded',onClick:()=>setMenuOpen(!menuOpen)},'☰'),
        menuOpen && React.createElement('div',{className:'absolute top-full left-0 mt-2 bg-amber-100 dark:bg-amber-800 text-black dark:text-white rounded shadow z-40'},[
          (!gameId && React.createElement('button',{className:'block w-full text-left px-4 py-2 whitespace-nowrap hover:bg-amber-200 dark:hover:bg-amber-700',onClick:()=>{setMenuOpen(false);startGame();}},'New Game')),
          (gameId && isRequester && React.createElement('button',{className:'block w-full text-left px-4 py-2 whitespace-nowrap hover:bg-amber-200 dark:hover:bg-amber-700',onClick:()=>{setMenuOpen(false);abandon();}},'Abandon Game')),
          (gameId && !isRequester && React.createElement('button',{className:'block w-full text-left px-4 py-2 whitespace-nowrap opacity-50 cursor-not-allowed',disabled:true},'New Game')),
          React.createElement('button',{className:'block w-full text-left px-4 py-2 whitespace-nowrap hover:bg-amber-200 dark:hover:bg-amber-700',onClick:()=>{setMenuOpen(false);setShowRename(true);}},'Rename Player'),
          React.createElement('button',{className:'block w-full text-left px-4 py-2 whitespace-nowrap hover:bg-amber-200 dark:hover:bg-amber-700',onClick:()=>{setMenuOpen(false);setShowSettings(true);}},'Settings'),
          React.createElement('button',{className:'block w-full text-left px-4 py-2 whitespace-nowrap hover:bg-amber-200 dark:hover:bg-amber-700',onClick:()=>{setMenuOpen(false);setShowHelp(true);}},'Help'),
          React.createElement('button',{className:'block w-full text-left px-4 py-2 whitespace-nowrap hover:bg-amber-200 dark:hover:bg-amber-700',onClick:()=>{setMenuOpen(false);openAbout();}},'About')
        ])
      ]),
      React.createElement('img',{src:`/assets/donkey-title-${actualTheme}.png`,className:'h-full w-auto',alt:'Donkey title'})
    ]),
    React.createElement('div',{className:'p-4 mt-20 space-y-4 flex flex-col items-center z-10'},[
      !showShare && state && !state.hasStarted && isRequester &&
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
    !showShare && React.createElement('div',{className:`fixed bottom-0 right-0 m-2 text-xs px-2 py-1 rounded ${connected?'bg-green-500':'bg-red-500'} text-white z-10`},
      user.name ? `${user.name}: ${connected?'Connected':'Disconnected'}` : (connected?'Connected':'Disconnected')
    ),
    gameId && React.createElement('div',{className:'fixed bottom-0 left-0 m-2 w-64 z-10'},[
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
      return React.createElement(
        'div',
        {className:'flex justify-center items-end h-32'},
        cards.map((c,i)=>
          React.createElement(
            'div',
            {
              key:i,
              onMouseEnter:()=>setSelectedCard(i),
              className:`relative w-16 h-32 -ml-12 first:ml-0 transition-all ${selectedCard===i?'z-10':''}`
            },
            React.createElement('img',{
              src:`/assets/${c}.png`,
              className:`absolute bottom-0 w-16 h-24 transition-transform ${selectedCard===i?'-translate-y-2':''}`,
              style:{pointerEvents:'none'}
            })
          )
        )
      );
    }
    return React.createElement('div',{className:'flex justify-center items-end'},
      Array(p.cardCount).fill(0).map((_,i)=>React.createElement('img',{key:i,src:`/assets/${backColor}_back.png`,className:'w-8 h-12 -ml-6 first:ml-0'}))
    );
  }

  function closeAbandoned(){
    setShowAbandoned(false);setGameId(null);setState(null);setLogs([]);window.history.replaceState(null,'','/');
  }
}

ReactDOM.render(React.createElement(App), document.getElementById('root'));
