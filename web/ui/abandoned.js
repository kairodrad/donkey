import { ModalWrapper, Button } from './utils.js';
const React = window.React;

export function AbandonedModal({ onClose }){
  return React.createElement(ModalWrapper, { onClose }, [
    React.createElement('div', { className: 'text-center space-y-3' }, [
      React.createElement('h3', { className: 'text-lg font-semibold text-red-600 dark:text-red-400' }, 
        'Game Ended'
      ),
      React.createElement('p', { className: 'text-gray-700 dark:text-gray-300' }, 
        'The creator disconnected and the game has ended.'
      )
    ]),
    React.createElement(Button, { onClick: onClose, className: 'mt-4 w-full' }, 'OK')
  ]);
}
