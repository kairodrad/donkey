import test from 'node:test';
import assert from 'node:assert/strict';

// Minimal React stub for server-side testing
global.window = { React: { createElement: (...args) => ({ tag: args[0], props: args[1], children: args.slice(2) }) } };

const { SettingsModal } = await import('./settings.js');

test('SettingsModal renders without crashing', () => {
  assert.doesNotThrow(() => {
    SettingsModal({
      theme: 'light',
      setTheme: () => {},
      backColor: 'red',
      setBackColor: () => {},
      onClose: () => {}
    });
  });
});
