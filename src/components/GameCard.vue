<template>
  <div
    :class="cardClasses"
    @click="handleClick"
    @touchstart="handleTouchStart"
    @touchend="handleTouchEnd"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
    :tabindex="interactive ? 0 : -1"
    @keydown="handleKeydown"
    :aria-label="ariaLabel"
    role="button"
    ref="cardRef"
    :data-card-id="cardId || `${rank}${suit}`"
  >
    <!-- Card face (front) - using static image -->
    <div
      v-if="!flipped"
      class="card-face card-front"
    >
      <img
        :src="cardImageSrc"
        :alt="cardImageAlt"
        class="card-image"
        @error="handleImageError"
      />
      
      <!-- Custom card content slot -->
      <div
        v-if="$slots.default"
        class="card-content"
      >
        <slot />
      </div>
    </div>
    
    <!-- Card back - using static image -->
    <div
      v-else
      class="card-face card-back"
    >
      <img
        :src="cardBackImageSrc"
        :alt="`${backColor} card back`"
        class="card-back-image"
        @error="handleBackImageError"
      />
      
      <!-- Custom back content slot -->
      <div
        v-if="$slots.back"
        class="card-back-content"
      >
        <slot name="back" />
      </div>
    </div>
    
    <!-- Card overlay for states -->
    <div
      v-if="overlay"
      class="card-overlay"
    >
      <slot name="overlay">
        <div class="card-overlay-content">
          {{ overlay }}
        </div>
      </slot>
    </div>
    
    <!-- Selection indicator -->
    <div
      v-if="selected"
      class="card-selection-indicator"
      aria-hidden="true"
    />
  </div>
</template>

<script>
import { computed, ref } from 'vue'

export default {
  name: 'GameCard',
  emits: ['click', 'select', 'hover'],
  props: {
    cardId: {
      type: String,
      default: null
    },
    rank: {
      type: String,
      required: true,
      validator: (value) => [
        'A', '2', '3', '4', '5', '6', '7', '8', '9', '10', 'J', 'Q', 'K'
      ].includes(value)
    },
    suit: {
      type: String,
      required: true,
      validator: (value) => ['hearts', 'diamonds', 'clubs', 'spades'].includes(value)
    },
    backColor: {
      type: String,
      default: 'red',
      validator: (value) => ['blue', 'gray', 'green', 'purple', 'red', 'yellow'].includes(value)
    },
    flipped: {
      type: Boolean,
      default: false
    },
    selected: {
      type: Boolean,
      default: false
    },
    disabledFeedback: {
      type: Boolean,
      default: false
    },
    disabled: {
      type: Boolean,
      default: false
    },
    interactive: {
      type: Boolean,
      default: true
    },
    size: {
      type: String,
      default: 'md',
      validator: (value) => ['xs', 'sm', 'md', 'lg', 'xl', 'user', 'current-player', 'opponent', 'ingame-opponent', 'ingame-current', 'discard'].includes(value)
    },
    overlay: {
      type: String,
      default: null
    },
    animation: {
      type: String,
      default: null,
      validator: (value) => !value || [
        'deal', 'flip', 'bounce', 'slide', 'fade'
      ].includes(value)
    },
    hoverEffect: {
      type: Boolean,
      default: true
    },
    fanRotation: {
      type: Number,
      default: 0
    },
    highestHighlight: {
      type: Boolean,
      default: false
    }
  },
  
  setup(props, { emit }) {
    const cardRef = ref(null)
    const isHovered = ref(false)
    
    // Suit abbreviation mapping for file names
    const suitAbbreviations = {
      hearts: 'H',
      diamonds: 'D',
      clubs: 'C',
      spades: 'S'
    }
    
    // Generate card image source path
    const cardImageSrc = computed(() => {
      const suitAbbr = suitAbbreviations[props.suit]
      const rankValue = props.rank
      return `/web/assets/${rankValue}${suitAbbr}.png`
    })
    
    // Generate card back image source path
    const cardBackImageSrc = computed(() => {
      return `/web/assets/${props.backColor}_back.png`
    })
    
    // Alt text for card image
    const cardImageAlt = computed(() => {
      return `${props.rank} of ${props.suit}`
    })
    
    // Size classes mapping - VISUAL.md specification
    const sizeClasses = {
      // Legacy sizes
      xs: 'w-6 h-8 text-xs',
      sm: 'w-8 h-12 text-xs',
      md: 'w-12 h-16 text-sm',
      lg: 'w-16 h-20 text-base',
      xl: 'w-20 h-28 text-lg',
      user: 'w-24 h-32 sm:w-28 sm:h-40 md:w-32 md:h-44 text-base',
      
      // New grid layout sizes (VISUAL.md specification)
      'current-player': 'w-24 h-32 sm:w-28 sm:h-40 md:w-32 md:h-44 text-base', // Current player cards
      'opponent': 'w-8 h-12 text-xs',                                            // Opponent cards (small for count visibility)
      // In-game cards reduced to 80% of previous Tailwind sizes for tighter layout
      'ingame-opponent': 'ingame-card-opponent text-sm',
      'ingame-current': 'ingame-card-current text-sm',
      'discard': 'w-8 h-12 text-xs'                                             // Discard pile cards (same as opponent cards)
    }
    
    // Animation classes mapping
    const animationClasses = {
      deal: 'animate-card-deal',
      flip: 'animate-card-flip',
      bounce: 'animate-bounce-in',
      slide: 'animate-slide-up',
      fade: 'animate-fade-in'
    }
    
    // Computed card classes
    const cardClasses = computed(() => {
      const classes = [
        'game-card',
        'relative',
        'select-none',
        'transition-all',
        'duration-300',
        sizeClasses[props.size]
      ]
      
      if (props.interactive && !props.disabled) {
        classes.push('cursor-pointer')
        
        if (props.hoverEffect) {
          classes.push('hover:card-hovered')
        }
      }
      
      if (props.selected && !props.disabled) {
        classes.push('card-selected', 'transform', '-translate-y-6', 'z-50')
      }
      
      if (props.disabled) {
        classes.push('card-disabled')
      }
      
      if (props.disabledFeedback) {
        classes.push('card-disabled-feedback')
      }
      
      if (props.highestHighlight) {
        classes.push('card-highest-highlight')
      }
      
      if (props.animation) {
        classes.push(animationClasses[props.animation])
      }
      
      if (isHovered.value && !props.disabled && (props.size === 'user' || props.size === 'current-player')) {
        classes.push('transform', '-translate-y-4', 'z-10')
      } else if (isHovered.value && !props.disabled) {
        classes.push('transform', '-translate-y-2')
      }
      
      return classes.join(' ')
    })
    
    // Fan rotation style
    const cardStyle = computed(() => {
      if (props.fanRotation !== 0) {
        return {
          transform: `rotate(${props.fanRotation}deg)`,
          transformOrigin: 'bottom center'
        }
      }
      return {}
    })
    
    // Accessibility label
    const ariaLabel = computed(() => {
      if (props.flipped) {
        return 'Face down card'
      }
      return `${props.rank} of ${props.suit}`
    })
    
    // Event handlers
    const handleClick = (event) => {
      if (!props.interactive || props.disabled) return
      
      emit('click', {
        rank: props.rank,
        suit: props.suit,
        selected: props.selected,
        id: props.cardId || `${props.rank}${props.suit}`,
        event
      })
      
      // Only emit select if parent doesn't handle it manually
      // This prevents conflicts with manual selection logic in App.vue
      if (!event.defaultPrevented) {
        emit('select', !props.selected)
      }
    }
    
    const handleMouseEnter = () => {
      if (!props.interactive || props.disabled) return
      
      isHovered.value = true
      emit('hover', {
        rank: props.rank,
        suit: props.suit,
        hovered: true
      })
    }
    
    const handleMouseLeave = () => {
      // Always clear hover state on mouse leave, even for disabled cards
      isHovered.value = false
      
      if (!props.disabled) {
        emit('hover', {
          rank: props.rank,
          suit: props.suit,
          hovered: false
        })
      }
    }
    
    const handleKeydown = (event) => {
      if (!props.interactive || props.disabled) return
      
      if (event.key === 'Enter' || event.key === ' ') {
        event.preventDefault()
        handleClick(event)
      }
    }
    
    // Touch event handlers for mobile support
    let touchStartTime = 0
    
    const handleTouchStart = (event) => {
      if (!props.interactive || props.disabled) return
      touchStartTime = Date.now()
      
      // Add visual feedback for touch
      if (cardRef.value) {
        cardRef.value.style.transform = 'scale(0.95)'
      }
    }
    
    const handleTouchEnd = (event) => {
      if (!props.interactive || props.disabled) return
      
      // Remove visual feedback
      if (cardRef.value) {
        cardRef.value.style.transform = ''
      }
      
      // Only trigger click if it was a quick tap (not a scroll/swipe)
      const touchDuration = Date.now() - touchStartTime
      if (touchDuration < 500) {
        event.preventDefault() // Prevent double-tap zoom and mouse click
        handleClick(event)
      }
    }
    
    // Error handling for missing images
    const handleImageError = (event) => {
      event.target.style.display = 'none'
    }
    
    const handleBackImageError = (event) => {
      event.target.style.display = 'none'
    }
    
    return {
      cardRef,
      cardClasses,
      cardStyle,
      cardImageSrc,
      cardBackImageSrc,
      cardImageAlt,
      ariaLabel,
      handleClick,
      handleTouchStart,
      handleTouchEnd,
      handleMouseEnter,
      handleMouseLeave,
      handleKeydown,
      handleImageError,
      handleBackImageError
    }
  }
}
</script>

<style scoped>
.game-card {
  perspective: 1000px;
}

.card-face {
  width: 100%;
  height: 100%;
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.card-front {
  border: 1px solid #e5e7eb;
}

.card-back {
  border: 1px solid #d1d5db;
}

.card-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  border-radius: 8px;
}

/* Precise in-game card sizing at ~80% of previous Tailwind sizes (w-16 h-24) */
.ingame-card-opponent,
.ingame-card-current {
  width: 3.2rem;  /* 80% of 4rem (w-16) */
  height: 4.8rem; /* 80% of 6rem (h-24) */
}

.card-back-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  border-radius: 8px;
}

.card-content {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

.card-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 20;
}

.card-overlay-content {
  color: white;
  font-weight: bold;
  text-align: center;
  padding: 1rem;
}

.card-selection-indicator {
  position: absolute;
  top: -3px;
  left: -3px;
  right: -3px;
  bottom: -3px;
  border: 4px solid var(--color-primary);
  border-radius: 12px;
  pointer-events: none;
  box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.3), 0 8px 24px rgba(56, 189, 248, 0.4);
  z-index: 30;
  animation: pulse-glow 2s infinite;
}

@keyframes pulse-glow {
  0%, 100% { 
    box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.3), 0 8px 24px rgba(56, 189, 248, 0.4);
  }
  50% { 
    box-shadow: 0 0 0 5px rgba(56, 189, 248, 0.5), 0 12px 32px rgba(56, 189, 248, 0.6);
  }
}

/* Mobile touch optimization */
@media (max-width: 768px) {
  .game-card {
    min-width: 32px;
    min-height: 48px;
  }
  
  /* Larger touch targets for user cards */
  .w-20.h-28 {
    min-width: 48px;
    min-height: 64px;
  }
}

/* Reduced motion support */
@media (prefers-reduced-motion: reduce) {
  .game-card {
    transition: none;
  }
  
  .card-face {
    animation: none !important;
  }
  
  .card-selection-indicator {
    animation: none !important;
  }
}

/* Disabled card styling - maintain normal appearance */
.card-disabled {
  /* Keep normal appearance, just not interactive */
  opacity: 1 !important;
  cursor: not-allowed;
  filter: none !important;
  -webkit-filter: none !important;
  color: inherit !important;
  background-color: inherit !important;
  border-color: inherit !important;
  /* Disable all transitions to prevent unwanted animations */
  transition: none !important;
}

.card-disabled .card-image,
.card-disabled .card-back-image {
  opacity: 1 !important;
  filter: none !important;
  -webkit-filter: none !important;
}

/* Red border feedback for disabled card clicks */
.card-disabled-feedback {
  border: 3px solid #ef4444 !important; /* Red border */
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.3), 0 4px 12px rgba(239, 68, 68, 0.2) !important;
  animation: shake 0.5s ease-in-out !important;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-2px); }
  75% { transform: translateX(2px); }
}

/* Yellow border highlight for highest value card */
.card-highest-highlight {
  border: 4px solid #fbbf24 !important; /* Yellow 400 */
  box-shadow: 0 0 0 2px rgba(251, 191, 36, 0.4), 0 4px 12px rgba(251, 191, 36, 0.3) !important;
  animation: highest-glow 2s infinite ease-in-out !important;
}

@keyframes highest-glow {
  0%, 100% { 
    box-shadow: 0 0 0 2px rgba(251, 191, 36, 0.4), 0 4px 12px rgba(251, 191, 36, 0.3);
    border-color: #fbbf24;
  }
  50% { 
    box-shadow: 0 0 0 4px rgba(251, 191, 36, 0.6), 0 8px 20px rgba(251, 191, 36, 0.5);
    border-color: #f59e0b; /* Slightly darker yellow for variation */
  }
}
</style>
