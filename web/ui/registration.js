const React = window.React;

export function RegistrationModal({onSubmit}){
  const [name,setName]=React.useState('');
  function submit(){
    const cleaned=name.replace(/[^\w\s]/g,'').trim().slice(0,20);
    if(cleaned) onSubmit(cleaned);
  }
  return React.createElement('div',{className:'fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50'},
    React.createElement('div',{className:'bg-white dark:bg-slate-700 text-black dark:text-white p-4 rounded space-y-2'},[
      React.createElement('h2',{className:'text-lg font-bold mb-2'},'Register'),
      React.createElement('input',{className:'border p-2 w-full bg-white text-black dark:bg-gray-800 dark:text-white',value:name,maxLength:20,onChange:e=>setName(e.target.value.replace(/[^\w\s]/g,'')),onKeyDown:e=>{if(e.key==='Enter')submit();}}),
      React.createElement('button',{className:'mt-2 px-4 py-1 bg-pink-200 dark:bg-pink-700 text-black dark:text-white rounded',onClick:submit},'Submit')
    ])
  );
}
