export const backs = ['red','blue','green','gray','purple','yellow'];

export function getCookie(name){
  const m=document.cookie.match('(?:^|; )'+name+'=([^;]*)');
  return m?decodeURIComponent(m[1]):null;
}

export function setCookie(name,value){
  document.cookie=name+'='+encodeURIComponent(value)+'; path=/';
}

export function seats(players){
  const pos=[{top:0,left:200},{top:50,left:350},{top:200,left:400},{top:350,left:350},{top:400,left:200},{top:350,left:50},{top:200,left:0},{top:50,left:50}];
  const n=players.length;
  return players.map((p,i)=>({...p,seat:pos[Math.floor(i*8/n)]}));
}
