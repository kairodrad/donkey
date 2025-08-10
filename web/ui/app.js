import { backs, getCookie, setCookie, seats } from './utils.js';
import { RegistrationModal } from './registration.js';

// theme handling
(function(){
  let theme=getCookie('theme');
  if(!theme){
    theme=window.matchMedia('(prefers-color-scheme: dark)').matches?'dark':'light';
    setCookie('theme',theme);
  }
  if(theme==='dark'){
    document.documentElement.classList.add('dark');
  }
})();

function App(){
  const params=new URLSearchParams(window.location.search);
  const gameIdParam=params.get('gameId');
  const [user,setUser]=React.useState({id:getCookie('userId'),name:getCookie('userName')});
  const [gameId,setGameId]=React.useState(gameIdParam);
  const [state,setState]=React.useState(null);
  const [backColor,setBackColor]=React.useState(getCookie('cardBack')||'red');
  const [showReg,setShowReg]=React.useState(!user.id);

  React.useEffect(()=>{setCookie('cardBack',backColor);},[backColor]);

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

  return React.createElement('div',{className:'p-4 space-y-4'},[
    React.createElement('div',{className:'flex justify-between'},[
      React.createElement('div',null,user.name?`Welcome, ${user.name}`:''),
      React.createElement('select',{value:backColor,onChange:e=>setBackColor(e.target.value),className:'border rounded p-1 bg-white dark:bg-slate-700'},backs.map(b=>React.createElement('option',{key:b,value:b},b)))
    ]),
    !gameId && React.createElement('button',{className:'px-4 py-2 bg-pink-200 text-black rounded',onClick:startGame},'New Game'),
    gameId && state && !state.hasStarted && React.createElement('div',{className:'space-x-2'},[
      React.createElement('button',{className:'px-3 py-1 bg-blue-200 text-black rounded',onClick:copyLink},'Copy Game URL'),
      state.requesterId==user.id && React.createElement('button',{className:'px-3 py-1 bg-green-200 text-black rounded',onClick:finalize},'Finalize Players and Deal')
    ]),
    state && renderPlayers(),
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
