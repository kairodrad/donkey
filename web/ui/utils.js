export const backs = ['red','blue','green','gray','purple','yellow'];

export function getCookie(name){
  const m=document.cookie.match('(?:^|; )'+name+'=([^;]*)');
  return m?decodeURIComponent(m[1]):null;
}

export function setCookie(name,value){
  document.cookie=name+'='+encodeURIComponent(value)+'; path=/';
}

export function seats(players,currentUserId){
  const pos=[{top:0,left:200},{top:50,left:350},{top:200,left:400},{top:350,left:350},{top:400,left:200},{top:350,left:50},{top:200,left:0},{top:50,left:50}];
  const n=players.length;
  const start=players.findIndex(p=>p.id===currentUserId);
  const ordered=start>=0?players.slice(start).concat(players.slice(0,start)):players;
  return ordered.map((p,i)=>({...p,seat:pos[(4+Math.floor(i*8/n))%8]}));
}

const suitOrder={D:0,C:1,H:2,S:3};
const rankOrder={'2':0,'3':1,'4':2,'5':3,'6':4,'7':5,'8':6,'9':7,'10':8,'J':9,'Q':10,'K':11,'A':12};

export function sortCards(cards){
  return [...cards].sort((a,b)=>{
    const sA=a.slice(-1), sB=b.slice(-1);
    const rA=a.slice(0,-1), rB=b.slice(0,-1);
    if(suitOrder[sA]!==suitOrder[sB]) return suitOrder[sA]-suitOrder[sB];
    return rankOrder[rA]-rankOrder[rB];
  });
}
