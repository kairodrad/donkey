import test from 'node:test';
import assert from 'node:assert/strict';

// Mock DOM and window environment for theme testing
let appliedTheme = null;
let cookieStore = {};
let currentSystemTheme = 'light';

global.window = {
  matchMedia: (query) => ({
    matches: query.includes('dark') ? currentSystemTheme === 'dark' : currentSystemTheme === 'light'
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

// Mock utils functions for testing
function getCookie(name) {
  return cookieStore[name] || null;
}

function setCookie(name, value) {
  cookieStore[name] = value;
}

function applyTheme(theme) {
  const actualTheme = theme === 'system' 
    ? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light') 
    : theme;
  document.documentElement.classList.toggle('dark', actualTheme === 'dark');
  return actualTheme;
}

test('applyTheme correctly applies light theme', () => {
  currentSystemTheme = 'light';
  const result = applyTheme('light');
  assert.strictEqual(result, 'light');
  assert.strictEqual(appliedTheme, 'light');
});

test('applyTheme correctly applies dark theme', () => {
  const result = applyTheme('dark');
  assert.strictEqual(result, 'dark');
  assert.strictEqual(appliedTheme, 'dark');
});

test('applyTheme handles system theme - light', () => {
  currentSystemTheme = 'light';
  const result = applyTheme('system');
  assert.strictEqual(result, 'light');
  assert.strictEqual(appliedTheme, 'light');
});

test('applyTheme handles system theme - dark', () => {
  currentSystemTheme = 'dark';
  const result = applyTheme('system');
  assert.strictEqual(result, 'dark');
  assert.strictEqual(appliedTheme, 'dark');
});

test('cookie persistence for theme setting', () => {
  // Clear cookie store
  cookieStore = {};
  
  // Set theme cookie
  setCookie('theme', 'dark');
  assert.strictEqual(getCookie('theme'), 'dark');
  
  // Test initial theme loading
  const initialTheme = getCookie('theme') || 'system';
  assert.strictEqual(initialTheme, 'dark');
});

test('cookie persistence for card back color', () => {
  // Clear cookie store
  cookieStore = {};
  
  // Set card back cookie
  setCookie('cardBack', 'blue');
  assert.strictEqual(getCookie('cardBack'), 'blue');
  
  // Test default fallback
  cookieStore = {};
  const defaultColor = getCookie('cardBack') || 'red';
  assert.strictEqual(defaultColor, 'red');
});

test('theme asset URL generation', () => {
  currentSystemTheme = 'light';
  
  // Test light theme asset paths
  const lightTheme = applyTheme('light');
  assert.strictEqual(lightTheme, 'light');
  const lightBgPath = `/assets/donkey-background-${lightTheme}.png`;
  const lightTitlePath = `/assets/donkey-title-${lightTheme}.png`;
  assert.strictEqual(lightBgPath, '/assets/donkey-background-light.png');
  assert.strictEqual(lightTitlePath, '/assets/donkey-title-light.png');
  
  // Test dark theme asset paths
  const darkTheme = applyTheme('dark');
  assert.strictEqual(darkTheme, 'dark');
  const darkBgPath = `/assets/donkey-background-${darkTheme}.png`;
  const darkTitlePath = `/assets/donkey-title-${darkTheme}.png`;
  assert.strictEqual(darkBgPath, '/assets/donkey-background-dark.png');
  assert.strictEqual(darkTitlePath, '/assets/donkey-title-dark.png');
});

test('system theme switching behavior', () => {
  // Test system theme changes with light preference
  currentSystemTheme = 'light';
  let result = applyTheme('system');
  assert.strictEqual(result, 'light');
  assert.strictEqual(appliedTheme, 'light');
  
  // Simulate system theme change to dark
  currentSystemTheme = 'dark';
  result = applyTheme('system');
  assert.strictEqual(result, 'dark');
  assert.strictEqual(appliedTheme, 'dark');
});

test('theme persistence workflow', () => {
  // Clear cookies
  cookieStore = {};
  
  // Simulate user selecting dark theme
  setCookie('theme', 'dark');
  const savedTheme = getCookie('theme');
  applyTheme(savedTheme);
  
  assert.strictEqual(savedTheme, 'dark');
  assert.strictEqual(appliedTheme, 'dark');
  
  // Simulate page reload - theme should persist
  const reloadedTheme = getCookie('theme') || 'system';
  applyTheme(reloadedTheme);
  
  assert.strictEqual(reloadedTheme, 'dark');
  assert.strictEqual(appliedTheme, 'dark');
});

test('invalid theme values fallback gracefully', () => {
  // Test with invalid theme value - should not crash
  assert.doesNotThrow(() => {
    applyTheme('invalid-theme');
  });
  
  // Invalid theme should be treated as light (not 'system' or 'dark')
  applyTheme('invalid-theme');
  assert.strictEqual(appliedTheme, 'light');
});