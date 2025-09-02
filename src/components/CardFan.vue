<template>
  <div 
    :class="fanClasses"
    :style="fanStyle"
  >
    <GameCard
      v-for="(card, index) in cards"
      :key="`${card.rank}${card.suit}`"
      :card-id="card.id"
      :rank="card.rank"
      :suit="card.suit"
      :back-color="backColor"
      :flipped="card.flipped"
      :selected="card.selected"
      :disabled="card.disabled"
      :disabled-feedback="card.disabledFeedback"
      :interactive="interactive"
      :size="cardSize"
      :hover-effect="hoverEffect"
      :fan-rotation="calculateRotation(index)"
      :style="calculateCardPosition(index)"
      @click="handleCardClick(card, index)"
      @select="handleCardSelect(card, index, $event)"
      @hover="handleCardHover(card, index, $event)"
      class="card-fan-item"
      :class="cardFanItemClass"
    />
  </div>
</template>

<script>
import { computed } from 'vue'
import GameCard from './GameCard.vue'

export default {
  name: 'CardFan',
  components: {
    GameCard
  },
  emits: ['card-click', 'card-select', 'card-hover'],
  props: {
    cards: {
      type: Array,
      required: true,
      default: () => []
    },
    backColor: {
      type: String,
      default: 'red'
    },
    interactive: {
      type: Boolean,
      default: true
    },
    isUser: {
      type: Boolean,
      default: false
    },
    maxRotation: {
      type: Number,
      default: 45 // Max rotation in degrees for the fan spread
    },
    spacing: {
      type: Number,
      default: 0.6 // Multiplier for card overlap
    },
    hoverEffect: {
      type: Boolean,
      default: true
    },
    size: {
      type: String,
      default: null // Will use 'current-player' for user cards, 'opponent' for opponents
    }
  },
  
  setup(props, { emit }) {
    // Determine card size based on user status and layout
    const cardSize = computed(() => {
      if (props.size) return props.size
      return props.isUser ? 'current-player' : 'opponent'
    })
    
    // Fan container classes
    const fanClasses = computed(() => {
      const classes = ['card-fan', 'relative']
      
      if (props.isUser) {
        classes.push('card-fan-user')
      } else {
        classes.push('card-fan-opponent')
      }
      
      return classes.join(' ')
    })
    
    // Fan container style
    const fanStyle = computed(() => {
      return {
        perspective: '1000px'
      }
    })
    
    // Card fan item classes
    const cardFanItemClass = computed(() => {
      return props.isUser ? 'card-fan-item-user' : 'card-fan-item-opponent'
    })
    
    // Calculate rotation for each card in the fan
    const calculateRotation = (index) => {
      if (props.cards.length <= 1) return 0
      
      const totalCards = props.cards.length
      const centerIndex = (totalCards - 1) / 2
      const rotationStep = props.maxRotation / Math.max(1, centerIndex)
      
      return (index - centerIndex) * rotationStep
    }
    
    // Calculate position for each card based on new grid layout
    const calculateCardPosition = (index) => {
      if (props.cards.length <= 1) {
        return {
          zIndex: 1
        }
      }
      
      const totalCards = props.cards.length
      const centerIndex = (totalCards - 1) / 2
      
      // Different spacing for current player vs opponents
      let spacing
      if (props.isUser) {
        // Current player: adaptive spacing for bottom fan (70vw available)
        const maxCards = 51
        const viewportWidth = typeof window !== 'undefined' ? window.innerWidth : 375
        const availableWidth = viewportWidth * 0.70 // 70vw available
        const maxSpacing = Math.max(8, Math.floor((availableWidth - 100) / Math.min(totalCards, maxCards)))
        spacing = Math.min(maxSpacing, 30) // Max 30px spacing
      } else {
        // Opponents: center fan around the exact midpoint of the container
        // Use conservative spacing to avoid excessive spread
        spacing = Math.min(12, 24)
      }
      
      let offset
      if (props.isUser) {
        // Current player: center-based positioning moved left by 5% width for better accommodation
        const viewportWidth = typeof window !== 'undefined' ? window.innerWidth : 375
        const leftShift = viewportWidth * 0.05 // Move left by 5% of viewport width
        const centerOffset = 40 - leftShift // Reduced base offset and move left by 5% width
        offset = (index - centerIndex) * spacing * props.spacing + centerOffset
      } else {
        // Opponents: center cards relative to 50% container width
        const centerOffset = 0 // center baseline
        offset = (index - centerIndex) * spacing * props.spacing + centerOffset
      }
      
      const rotation = calculateRotation(index)
      const zIndex = index + 1
      
      const style = {
        transform: `translateX(${offset}px) rotate(${rotation}deg)`,
        transformOrigin: props.isUser ? 'bottom center' : 'center center',
        zIndex: zIndex,
        // Set CSS custom properties for hover state maintenance
        '--translate-x': `${offset}px`,
        '--translate-y': '0px',
        '--rotation': `${rotation}deg`,
        '--z-index': zIndex
      }

      // Anchor opponent cards to the horizontal center of the container
      if (!props.isUser) {
        style.left = '50%'
      }

      return style
    }
    
    // Event handlers
    const handleCardClick = (card, index) => {
      emit('card-click', { card, index })
    }
    
    const handleCardSelect = (card, index, selected) => {
      emit('card-select', { card, index, selected })
    }
    
    const handleCardHover = (card, index, hoverData) => {
      emit('card-hover', { card, index, ...hoverData })
    }
    
    return {
      cardSize,
      fanClasses,
      fanStyle,
      cardFanItemClass,
      calculateRotation,
      calculateCardPosition,
      handleCardClick,
      handleCardSelect,
      handleCardHover
    }
  }
}
</script>

<style scoped>
.card-fan {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

/* Current player fan - positioned within parent container */
.card-fan-user {
  position: relative;
  padding: 0 0.5rem;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

/* Opponent fan - positioned within grid cell */
.card-fan-opponent {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding-left: 0;
}

.card-fan-item {
  position: absolute;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Disable transitions for disabled cards to prevent unwanted animations */
.card-fan-item.card-disabled {
  transition: none !important;
}

.card-fan-item-user {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Disable transitions for disabled user cards */
.card-fan-item-user.card-disabled {
  transition: none !important;
}

.card-fan-item-user:not(.card-disabled):hover {
  transform: translateX(var(--translate-x, 0)) translateY(-1rem) rotate(var(--rotation, 0deg)) !important;
  z-index: 100 !important;
}

/* Ensure disabled cards have no hover effects */
.card-fan-item-user.card-disabled:hover {
  transform: translateX(var(--translate-x, 0)) translateY(var(--translate-y, 0)) rotate(var(--rotation, 0deg)) !important;
  z-index: var(--z-index, auto) !important;
}

.card-fan-item-opponent {
  pointer-events: none;
}

/* Mobile optimizations */
@media (max-width: 768px) {
  .card-fan-user {
    padding: 0 0.25rem;
    max-width: 99vw;
  }
  
  .card-fan-opponent {
    padding-left: 0.25rem;
  }
  
  .card-fan-item-user {
    /* Ensure adequate touch targets on mobile */
    min-width: 44px;
    min-height: 60px;
  }
}

@media (max-width: 480px) {
  .card-fan-user {
    padding: 0 0.125rem;
  }
  
  .card-fan-opponent {
    padding-left: 0.125rem;
  }
}

/* Ensure smooth animations */
@media (prefers-reduced-motion: reduce) {
  .card-fan-item {
    transition: none;
  }
}
</style>
