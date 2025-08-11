const React = window.React;

export function AboutModal({version,onClose}){
  return React.createElement('div',{className:'fixed inset-0 flex items-center justify-center bg-black bg-opacity-50'},
    React.createElement('div',{className:'bg-white dark:bg-slate-700 text-black dark:text-white p-4 rounded space-y-2 text-center'},[
      React.createElement('h2',{className:'text-lg font-bold'},'DONKEY'),
      React.createElement('p',null,`Version: ${version}`),
      React.createElement('p',null,'Created by Deepak Amin.'),
      React.createElement('button',{className:'mt-2 px-4 py-1 bg-pink-200 dark:bg-pink-700 text-black dark:text-white rounded',onClick:onClose},'Close')
    ])
  );
}

