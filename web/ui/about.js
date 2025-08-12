const React = window.React;

export function AboutModal({version,onClose}){
  return React.createElement('div',{className:'fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50'},
    React.createElement('div',{className:'bg-white dark:bg-slate-700 text-black dark:text-white p-4 rounded space-y-2'},[
      React.createElement('h2',{className:'text-lg font-bold text-center'},'About'),
      React.createElement('img',{src:'/apple-touch-icon.png',className:'mx-auto w-24 h-24',alt:'Donkey icon'}),
      React.createElement('p',{className:'text-right'},`Version: ${version}`),
      React.createElement('p',{className:'text-right'},'Created by Deepak Amin.'),
      React.createElement('button',{className:'mt-2 px-4 py-1 bg-pink-200 dark:bg-pink-700 text-black dark:text-white rounded',onClick:onClose},'Close')
    ])
  );
}

