const React = window.React;
import { backs } from './utils.js';

export function SettingsModal({theme,setTheme,backColor,setBackColor,onClose,user,onRename}){
  const [name,setName]=React.useState(user.name);
  function rename(){
    const cleaned=name.replace(/[^\w\s]/g,'').trim().slice(0,20);
    if(cleaned && cleaned!==user.name) onRename(cleaned);
  }
  return React.createElement('div',{className:'fixed inset-0 flex items-center justify-center bg-black bg-opacity-50'},
    React.createElement('div',{className:'bg-white dark:bg-slate-700 text-black dark:text-white p-4 rounded space-y-4'},[
      React.createElement('h2',{className:'text-lg font-bold'},'Settings'),
      React.createElement('div',{className:'space-y-2'},[
        React.createElement('div',null,[
          React.createElement('label',{className:'mr-2'},'Name:'),
          React.createElement('input',{className:'border p-1 bg-white text-black dark:bg-slate-800 dark:text-white',value:name,maxLength:20,onChange:e=>setName(e.target.value.replace(/[^\w\s]/g,'')),onKeyDown:e=>{if(e.key==='Enter')rename();}}),
          React.createElement('button',{className:'ml-2 px-2 py-1 bg-blue-200 dark:bg-blue-700 text-black dark:text-white rounded',onClick:rename},'Update')
        ]),
        React.createElement('div',null,[
          React.createElement('label',{className:'mr-2'},'Theme:'),
          React.createElement('select',{value:theme,onChange:e=>setTheme(e.target.value),className:'border p-1 bg-white text-black dark:bg-slate-800 dark:text-white'},[
            React.createElement('option',{value:'system'},'System'),
            React.createElement('option',{value:'light'},'Light'),
            React.createElement('option',{value:'dark'},'Dark')
          ])
        ]),
        React.createElement('div',null,[
          React.createElement('label',{className:'mr-2'},'Card Back:'),
          React.createElement('select',{value:backColor,onChange:e=>setBackColor(e.target.value),className:'border p-1 bg-white text-black dark:bg-slate-800 dark:text-white'},
            backs.map(b=>React.createElement('option',{key:b,value:b},b)))
        ])
      ]),
      React.createElement('button',{className:'px-4 py-1 bg-pink-200 dark:bg-pink-700 text-black dark:text-white rounded',onClick:onClose},'Close')
    ])
  );
}
