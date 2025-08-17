import { ModalWrapper, Button } from './utils.js';
const React = window.React;

export function HelpModal({ onClose }){
  const steps = [
    'Register a name to begin.',
    'Start a new game and share the URL with up to seven friends.',
    'Players join via the shared link; the requester can finalize players.',
    'Once finalized, the deck is shuffled and dealt clockwise.',
    'You see your own cards while opponents show only their card backs.'
  ];

  return React.createElement(ModalWrapper, { onClose }, [
    React.createElement('h2', { className: 'text-lg font-bold' }, 'How to Play'),
    React.createElement('div', { className: 'max-w-md' }, [
      React.createElement('ul', { className: 'list-disc pl-5 space-y-2 text-sm' },
        steps.map((step, index) => 
          React.createElement('li', { key: index }, step)
        )
      )
    ]),
    React.createElement(Button, { onClick: onClose, className: 'mt-4 w-full' }, 'Close')
  ]);
}

