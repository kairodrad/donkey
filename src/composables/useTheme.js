import { ref, computed, watch, onMounted } from 'vue'

// Theme state - reactive and persistent
const theme = ref('system')
const systemPrefersDark = ref(false)

// Theme utilities
export function useTheme() {
  // Computed current theme
  const currentTheme = computed(() => {
    if (theme.value === 'system') {
      return systemPrefersDark.value ? 'dark' : 'light'
    }
    return theme.value
  })

  // Check if current theme is dark
  const isDark = computed(() => currentTheme.value === 'dark')

  // Media query for system preference
  const mediaQuery = computed(() => {
    if (typeof window !== 'undefined') {
      return window.matchMedia('(prefers-color-scheme: dark)')
    }
    return null
  })

  // Update system preference
  const updateSystemPreference = () => {
    if (mediaQuery.value) {
      systemPrefersDark.value = mediaQuery.value.matches
    }
  }

  // Apply theme to document
  const applyTheme = (newTheme) => {
    if (typeof document !== 'undefined') {
      const root = document.documentElement
      
      // Remove existing theme classes
      root.classList.remove('light', 'dark')
      
      // Add new theme class
      root.classList.add(newTheme)
      
      // Update meta theme-color for mobile browsers
      const metaThemeColor = document.querySelector('meta[name="theme-color"]')
      if (metaThemeColor) {
        metaThemeColor.setAttribute(
          'content', 
          newTheme === 'dark' ? '#111827' : '#ffffff'
        )
      }
    }
  }

  // Set theme
  const setTheme = (newTheme) => {
    if (['light', 'dark', 'system'].includes(newTheme)) {
      theme.value = newTheme
      
      // Persist to localStorage
      if (typeof localStorage !== 'undefined') {
        localStorage.setItem('donkey-theme', newTheme)
      }
    }
  }

  // Toggle between light and dark (skips system)
  const toggleTheme = () => {
    const newTheme = currentTheme.value === 'dark' ? 'light' : 'dark'
    setTheme(newTheme)
  }

  // Cycle through all themes: light -> dark -> system
  const cycleTheme = () => {
    const themes = ['light', 'dark', 'system']
    const currentIndex = themes.indexOf(theme.value)
    const nextIndex = (currentIndex + 1) % themes.length
    setTheme(themes[nextIndex])
  }

  // Initialize theme
  const initializeTheme = () => {
    // Check for saved theme preference
    if (typeof localStorage !== 'undefined') {
      const savedTheme = localStorage.getItem('donkey-theme')
      if (savedTheme && ['light', 'dark', 'system'].includes(savedTheme)) {
        theme.value = savedTheme
      }
    }

    // Set up system preference detection
    updateSystemPreference()
    if (mediaQuery.value) {
      mediaQuery.value.addEventListener('change', updateSystemPreference)
    }

    // Apply initial theme
    applyTheme(currentTheme.value)
  }

  // Watch for theme changes
  watch(
    currentTheme,
    (newTheme) => {
      applyTheme(newTheme)
    },
    { immediate: false }
  )

  // Cleanup function
  const cleanup = () => {
    if (mediaQuery.value) {
      mediaQuery.value.removeEventListener('change', updateSystemPreference)
    }
  }

  // Initialize on mount
  onMounted(() => {
    initializeTheme()
  })

  // Return public API
  return {
    // State
    theme: computed(() => theme.value),
    currentTheme,
    isDark,
    systemPrefersDark: computed(() => systemPrefersDark.value),
    
    // Actions
    setTheme,
    toggleTheme,
    cycleTheme,
    initializeTheme,
    cleanup,
    
    // Utilities
    isLight: computed(() => !isDark.value),
    themeLabel: computed(() => {
      const labels = {
        light: 'Light',
        dark: 'Dark',
        system: 'System'
      }
      return labels[theme.value] || 'System'
    })
  }
}

// Utility for getting theme colors in JavaScript
export function useThemeColors() {
  const { isDark } = useTheme()
  
  const colors = computed(() => {
    const baseColors = {
      background: isDark.value ? '#111827' : '#ffffff',
      backgroundSecondary: isDark.value ? '#1f2937' : '#f9fafb',
      surface: isDark.value ? '#1f2937' : '#ffffff',
      border: isDark.value ? '#4b5563' : '#e5e7eb',
      text: isDark.value ? '#f3f4f6' : '#111827',
      textSecondary: isDark.value ? '#d1d5db' : '#6b7280',
      primary: isDark.value ? '#7dd3fc' : '#38bdf8',
      secondary: isDark.value ? '#d8b4fe' : '#c084fc',
      accent: isDark.value ? '#f8b4b4' : '#f38686',
      success: isDark.value ? '#86efac' : '#4ade80',
      warning: isDark.value ? '#fdba74' : '#fb923c',
      error: isDark.value ? '#fca5a5' : '#f87171'
    }
    
    return baseColors
  })
  
  return {
    colors,
    isDark
  }
}