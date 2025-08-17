import { ModalWrapper, Button } from './utils.js';
const React = window.React;

export function ShareModal({ gameId, isRequester, playerCount, onFinalize }){
  const url = `${window.location.origin}${window.location.pathname}?gameId=${gameId}`;
  const [copied, setCopied] = React.useState(false);
  
  function copy(){
    navigator.clipboard.writeText(url);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }
  
  return React.createElement(ModalWrapper, { onClose: null }, [
    React.createElement('h2', { className: 'text-lg font-bold text-center' }, 
      isRequester ? 'Game Created!' : 'Game Link'
    ),
    React.createElement('p', { className: 'text-sm text-center text-gray-600 dark:text-gray-400' }, 
      isRequester 
        ? 'Share this URL with up to 7 friends to join your game:'
        : 'Share this URL to invite more players:'
    ),
    React.createElement('div', { className: 'bg-gray-100 dark:bg-gray-800 p-2 rounded text-xs break-all' }, url),
    React.createElement(Button, {
      variant: copied ? 'secondary' : 'warning',
      onClick: copied ? null : copy,
      disabled: copied,
      className: 'w-full'
    }, copied ? 'URL Copied to Clipboard!' : 'Copy URL'),
    
    isRequester && React.createElement(Button, {
      variant: 'success',
      onClick: playerCount > 1 ? onFinalize : null,
      disabled: playerCount <= 1,
      className: 'w-full'
    }, playerCount > 1 ? 'Finalize and Deal Cards' : `Need ${2 - playerCount} more player${2 - playerCount !== 1 ? 's' : ''}`),
    
    !isRequester && React.createElement('div', { 
      className: 'text-center text-sm text-gray-600 dark:text-gray-400 italic' 
    }, 'Waiting for the game creator to start...')
  ]);
}
