import test from 'node:test';
import assert from 'node:assert/strict';

// Enhanced React stub for server-side testing with theme functionality
let appliedTheme = null;
let cookieStore = {};
global.window = { 
  React: { 
    createElement: (...args) => ({ tag: args[0], props: args[1], children: args.slice(2) }),
    useState: (initial) => [initial, () => {}],
    useEffect: () => {}
  },
  matchMedia: (query) => ({
    matches: query.includes('dark') ? false : true
  })
};
global.document = {
  documentElement: {
    classList: {
      toggle: (className, condition) => {
        appliedTheme = condition ? 'dark' : 'light';
      }
    }
  },
  cookie: ''
};

// Mock utils
const mockUtils = {
  getCookie: (name) => cookieStore[name] || null,
  setCookie: (name, value) => { cookieStore[name] = value; }
};

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

test('SettingsModal handles all theme options', () => {
  const themes = ['light', 'dark', 'system'];
  
  themes.forEach(theme => {
    assert.doesNotThrow(() => {
      SettingsModal({
        theme,
        setTheme: () => {},
        backColor: 'red',
        setBackColor: () => {},
        onClose: () => {}
      });
    }, `Should handle ${theme} theme`);
  });
});

test('SettingsModal handles all card back colors', () => {
  const colors = ['red', 'blue', 'green', 'purple', 'yellow', 'gray'];
  
  colors.forEach(color => {
    assert.doesNotThrow(() => {
      SettingsModal({
        theme: 'light',
        setTheme: () => {},
        backColor: color,
        setBackColor: () => {},
        onClose: () => {}
      });
    }, `Should handle ${color} card back color`);
  });
});

test('SettingsModal handles edge cases', () => {
  // Test with undefined/null values
  assert.doesNotThrow(() => {
    SettingsModal({
      theme: undefined,
      setTheme: () => {},
      backColor: null,
      setBackColor: () => {},
      onClose: () => {}
    });
  });
});
