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

// Shared UI component utilities
const React = window.React;

// Base modal wrapper component
export function ModalWrapper({ children, onClose }) {
  return React.createElement('div', {
    className: 'fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50'
  }, 
    React.createElement('div', {
      className: 'bg-white dark:bg-slate-700 text-black dark:text-white p-4 rounded space-y-2'
    }, children)
  );
}

// Standard button component with consistent styling
export function Button({ 
  children, 
  onClick, 
  variant = 'primary', 
  disabled = false, 
  className = '',
  ...props 
}) {
  const baseClasses = 'px-4 py-1 rounded transition-colors';
  const variants = {
    primary: 'bg-pink-200 dark:bg-pink-700 text-black dark:text-white hover:bg-pink-300 dark:hover:bg-pink-600',
    success: 'bg-green-200 dark:bg-green-700 text-black dark:text-white hover:bg-green-300 dark:hover:bg-green-600',
    warning: 'bg-amber-200 dark:bg-amber-700 text-black dark:text-white hover:bg-amber-300 dark:hover:bg-amber-600',
    secondary: 'bg-gray-200 dark:bg-gray-700 text-black dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600'
  };
  
  const disabledClasses = disabled ? 'opacity-50 cursor-not-allowed' : '';
  const variantClasses = variants[variant] || variants.primary;
  
  return React.createElement('button', {
    className: `${baseClasses} ${variantClasses} ${disabledClasses} ${className}`.trim(),
    onClick: disabled ? null : onClick,
    disabled,
    ...props
  }, children);
}

// Standard input component with consistent styling
export function Input({ 
  className = '', 
  ...props 
}) {
  const baseClasses = 'border p-2 bg-white text-black dark:bg-gray-800 dark:text-white dark:border-gray-600';
  
  return React.createElement('input', {
    className: `${baseClasses} ${className}`.trim(),
    ...props
  });
}

// Standard select component with consistent styling
export function Select({ 
  children, 
  className = '', 
  ...props 
}) {
  const baseClasses = 'border p-1 bg-white text-black dark:bg-slate-800 dark:text-white dark:border-gray-600';
  
  return React.createElement('select', {
    className: `${baseClasses} ${className}`.trim(),
    ...props
  }, children);
}
