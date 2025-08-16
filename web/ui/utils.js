export const backs = ['red','blue','green','gray','purple','yellow'];

export function getCookie(name){
  const m=document.cookie.match('(?:^|; )'+name+'=([^;]*)');
  return m?decodeURIComponent(m[1]):null;
}

export function setCookie(name,value){
  document.cookie=name+'='+encodeURIComponent(value)+'; path=/';
}

// Theme utility functions
export function getActualTheme(theme) {
  return theme === 'system' 
    ? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
    : theme;
}

export function applyTheme(theme) {
  const actualTheme = getActualTheme(theme);
  document.documentElement.classList.toggle('dark', actualTheme === 'dark');
  return actualTheme;
}

export function getThemeAssetUrl(assetName, theme) {
  const actualTheme = getActualTheme(theme);
  return `/assets/${assetName}-${actualTheme}.png`;
}

// Image preloading functionality
export function preloadThemeAssets(themes = ['light', 'dark']) {
  const assetNames = ['donkey-background', 'donkey-title'];
  
  themes.forEach(theme => {
    assetNames.forEach(assetName => {
      const img = new Image();
      img.src = getThemeAssetUrl(assetName, theme);
      // Store references to prevent garbage collection
      if (!window.preloadedAssets) window.preloadedAssets = [];
      window.preloadedAssets.push(img);
    });
  });
}

// Enhanced image creation with fallback
export function createThemedImage(assetName, theme, fallbackTheme = 'light') {
  const img = new Image();
  const primarySrc = getThemeAssetUrl(assetName, theme);
  const fallbackSrc = getThemeAssetUrl(assetName, fallbackTheme);
  
  img.src = primarySrc;
  img.onerror = () => {
    if (img.src !== fallbackSrc) {
      img.src = fallbackSrc;
    }
  };
  
  return img;
}

export function seats(players,currentUserId){
  const pos=[{top:20,left:200},{top:60,left:340},{top:200,left:380},{top:340,left:340},{top:380,left:200},{top:340,left:60},{top:200,left:20},{top:60,left:60}];
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
