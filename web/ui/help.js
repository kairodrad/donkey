const React = window.React;

export function HelpModal({onClose}){
  return React.createElement('div',{className:'fixed inset-0 flex items-center justify-center bg-black bg-opacity-50'},
    React.createElement('div',{className:'bg-white dark:bg-slate-700 p-4 rounded space-y-2 max-w-md'},[
      React.createElement('h2',{className:'text-lg font-bold'},'How to Play'),
      React.createElement('ul',{className:'list-disc pl-5 space-y-1 text-sm'},[
        React.createElement('li',null,'Register a name to begin.'),
        React.createElement('li',null,'Start a new game and share the URL with up to seven friends.'),
        React.createElement('li',null,'Players join via the shared link; the requester can finalize players.'),
        React.createElement('li',null,'Once finalized, the deck is shuffled and dealt clockwise.'),
        React.createElement('li',null,'You see your own cards while opponents show only their card backs.')
      ]),
      React.createElement('button',{className:'mt-2 px-4 py-1 bg-pink-200 text-black rounded',onClick:onClose},'Close')
    ])
  );
}

