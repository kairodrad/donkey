const React = window.React;

export function ShareModal({gameId,isRequester,playerCount,onFinalize}){
  const url=`${window.location.origin}${window.location.pathname}?gameId=${gameId}`;
  function copy(){navigator.clipboard.writeText(url);}
  return React.createElement('div',{className:'fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center'},[
    React.createElement('div',{className:'bg-white dark:bg-gray-800 text-black dark:text-white p-4 rounded space-y-2'},[
      React.createElement('div',null,isRequester?'Game created! Please copy the URL to share.':'Please copy the URL to share.'),
      React.createElement('button',{className:'px-3 py-1 bg-blue-200 dark:bg-blue-700 text-black dark:text-white rounded w-full',onClick:copy},'Copy URL'),
      isRequester && React.createElement('button',{
        className:`px-3 py-1 bg-green-200 dark:bg-green-700 text-black dark:text-white rounded w-full ${playerCount>1?'':'opacity-50 cursor-not-allowed'}`,
        onClick:playerCount>1?onFinalize:null
      },'Finalize and Deal'),
      !isRequester && React.createElement('div',null,'Awaiting Creator to start the game...')
    ])
  ]);
}
