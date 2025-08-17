import { ModalWrapper, Button } from './utils.js';
const React = window.React;

export function AboutModal({ version, onClose }){
  return React.createElement(ModalWrapper, { onClose }, [
    React.createElement('h2', { className: 'text-lg font-bold text-center' }, 'About'),
    React.createElement('img', {
      src: '/apple-touch-icon.png',
      className: 'mx-auto w-24 h-24',
      alt: 'Donkey icon'
    }),
    React.createElement('div', { className: 'text-right space-y-1' }, [
      React.createElement('p', null, `Version: ${version}`),
      React.createElement('p', null, 'Created by Deepak Amin.')
    ]),
    React.createElement(Button, { onClick: onClose, className: 'mt-4 w-full' }, 'Close')
  ]);
}

