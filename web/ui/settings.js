import { backs, ModalWrapper, Button, Select } from './utils.js';
const React = window.React;

export function SettingsModal({
  theme = 'system',
  setTheme = () => {},
  backColor = 'red',
  setBackColor = () => {},
  onClose = () => {}
} = {}){
  return React.createElement(ModalWrapper, { onClose }, [
    React.createElement('h2', { className: 'text-lg font-bold' }, 'Settings'),
    React.createElement('div', { className: 'space-y-3' }, [
      React.createElement('div', { className: 'flex items-center gap-2' }, [
        React.createElement('label', { className: 'min-w-0 flex-shrink-0' }, 'Theme:'),
        React.createElement(Select, {
          value: theme,
          onChange: e => setTheme(e.target.value),
          className: 'flex-1'
        }, [
          React.createElement('option', { value: 'system' }, 'System'),
          React.createElement('option', { value: 'light' }, 'Light'),
          React.createElement('option', { value: 'dark' }, 'Dark')
        ])
      ]),
      React.createElement('div', { className: 'flex items-center gap-2' }, [
        React.createElement('label', { className: 'min-w-0 flex-shrink-0' }, 'Card Back:'),
        React.createElement(Select, {
          value: backColor,
          onChange: e => setBackColor(e.target.value),
          className: 'flex-1 capitalize'
        }, backs.map(back => 
          React.createElement('option', { key: back, value: back, className: 'capitalize' }, back)
        ))
      ])
    ]),
    React.createElement(Button, { onClick: onClose, className: 'mt-4' }, 'Close')
  ]);
}
