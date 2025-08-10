const React = window.React;

export function ShareModal({gameId,isRequester,playerCount,onFinalize}){
  const url=`${window.location.origin}${window.location.pathname}?gameId=${gameId}`;
  const [copied,setCopied]=React.useState(false);
  function copy(){
    navigator.clipboard.writeText(url);
    setCopied(true);
    setTimeout(()=>setCopied(false),2000);
  }
  const copyClasses=copied?'bg-gray-200 dark:bg-gray-700 text-gray-500 cursor-not-allowed':'bg-blue-200 dark:bg-blue-700 text-black dark:text-white';
  return React.createElement('div',{className:'fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center'},[
    React.createElement('div',{className:'bg-white dark:bg-gray-800 text-black dark:text-white p-4 rounded space-y-2'},[
      React.createElement('div',null,isRequester?'Game created! Please copy the URL to share.':'Please copy the URL to share.'),
      React.createElement('button',{className:`px-3 py-1 rounded w-full ${copyClasses}`,onClick:copied?null:copy,disabled:copied},copied?'URL copied to clipboard':'Copy URL'),
      isRequester && React.createElement('button',{
        className:`px-3 py-1 bg-green-200 dark:bg-green-700 text-black dark:text-white rounded w-full ${playerCount>1?'':'opacity-50 cursor-not-allowed'}`,
        onClick:playerCount>1?onFinalize:null
      },'Finalize and Deal'),
      !isRequester && React.createElement('div',null,'Awaiting Creator to start the game...')
    ])
  ]);
}
