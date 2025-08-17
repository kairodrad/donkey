import { ModalWrapper, Button, Input } from './utils.js';
const React = window.React;

export function RegistrationModal({ onSubmit }){
  const [name, setName] = React.useState('');
  
  function submit(){
    const cleaned = name.replace(/[^\w\s]/g,'').trim().slice(0,20);
    if(cleaned) onSubmit(cleaned);
  }
  
  function handleNameChange(e) {
    setName(e.target.value.replace(/[^\w\s]/g,''));
  }
  
  function handleKeyDown(e) {
    if(e.key === 'Enter') submit();
  }
  
  return React.createElement(ModalWrapper, { onClose: null }, [
    React.createElement('h2', { className: 'text-lg font-bold' }, 'Register'),
    React.createElement('p', { className: 'text-sm text-gray-600 dark:text-gray-400' }, 
      'Choose a name to join the game (alphanumeric and spaces only)'
    ),
    React.createElement(Input, {
      className: 'w-full',
      value: name,
      maxLength: 20,
      placeholder: 'Enter your name...',
      onChange: handleNameChange,
      onKeyDown: handleKeyDown,
      autoFocus: true
    }),
    React.createElement(Button, { 
      onClick: submit, 
      disabled: !name.trim(),
      className: 'mt-4 w-full' 
    }, 'Join Game')
  ]);
}
