const React = window.React;

export function AbandonedModal({onClose}){
  return React.createElement('div',{className:'fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50'},
    React.createElement('div',{className:'bg-white dark:bg-slate-700 text-black dark:text-white p-4 rounded space-y-2 text-center'},[
      React.createElement('p',null,'The creator disconnected; the game has ended.'),
      React.createElement('button',{className:'mt-2 px-4 py-1 bg-pink-200 dark:bg-pink-700 text-black dark:text-white rounded',onClick:onClose},'OK')
    ])
  );
}
