<template>
  <div 
    class="discard-pile-container"
    data-discard-pile
  >
    <!-- Discard Pile Stack -->
    <div 
      v-if="discardPile.length > 0"
      class="discard-pile-stack"
    >
      <!-- Show last 5 cards with stacking effect -->
      <GameCard
        v-for="(card, index) in visibleDiscardCards"
        :key="`discard-${card.id}-${index}`"
        :rank="card.rank"
        :suit="card.suit"
        :back-color="backColor"
        :flipped="false"
        :interactive="false"
        size="sm"
        :style="getDiscardCardStyle(index)"
        class="discard-card"
        :data-card-id="card.id"
      />
      
      <!-- Card count indicator if more than 5 cards -->
      <div 
        v-if="discardPile.length > 5"
        class="discard-count-indicator"
        :style="getCountIndicatorStyle()"
      >
        +{{ discardPile.length - 5 }}
      </div>
    </div>
    
    <!-- Empty discard pile placeholder -->
    <div 
      v-else
      class="discard-pile-empty"
    >
      <div class="empty-pile-outline">
        <span class="empty-pile-text">Discard Pile</span>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import GameCard from './GameCard.vue'

export default {
  name: 'DiscardPile',
  components: {
    GameCard
  },
  props: {
    discardPile: {
      type: Array,
      required: true,
      default: () => []
    },
    backColor: {
      type: String,
      default: 'red'
    },
    animated: {
      type: Boolean,
      default: true
    }
  },
  setup(props) {
    // Show last 5 cards for visual stacking
    const visibleDiscardCards = computed(() => {
      return props.discardPile.slice(-5)
    })
    
    // Calculate card position in stack
    const getDiscardCardStyle = (index) => {
      const baseZIndex = 10
      const offsetX = index * 2 // 2px horizontal offset per card
      const offsetY = index * -1 // 1px vertical offset per card (upward)
      const rotation = (index - 2) * 2 // Slight rotation for natural look
      
      return {
        position: 'absolute',
        transform: `translate(${offsetX}px, ${offsetY}px) rotate(${rotation}deg)`,
        zIndex: baseZIndex + index,
        transition: props.animated ? 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)' : 'none'
      }
    }
    
    // Position for count indicator
    const getCountIndicatorStyle = () => {
      const topCardIndex = Math.min(4, props.discardPile.length - 1)
      const offsetX = topCardIndex * 2 + 10
      const offsetY = topCardIndex * -1 - 5
      
      return {
        position: 'absolute',
        transform: `translate(${offsetX}px, ${offsetY}px)`,
        zIndex: 20
      }
    }
    
    return {
      visibleDiscardCards,
      getDiscardCardStyle,
      getCountIndicatorStyle
    }
  }
}
</script>

<style scoped>
.discard-pile-container {
  position: fixed;
  top: 5rem; /* Below header */
  right: 1rem;
  z-index: 15;
  min-width: 4rem;
}

.discard-pile-stack {
  position: relative;
  width: 2rem; /* w-8 */
  height: 3rem; /* h-12 */
  min-height: 3rem;
}

.discard-card {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  border: 1px solid var(--color-border);
}

.discard-count-indicator {
  background: var(--color-primary);
  color: white;
  font-size: 0.625rem;
  font-weight: 600;
  padding: 0.125rem 0.25rem;
  border-radius: 0.5rem;
  min-width: 1.25rem;
  text-align: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.discard-pile-empty {
  width: 2rem; /* w-8 */
  height: 3rem; /* h-12 */
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-pile-outline {
  width: 100%;
  height: 100%;
  border: 2px dashed var(--color-border);
  border-radius: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-background-secondary);
  opacity: 0.6;
}

.empty-pile-text {
  font-size: 0.5rem;
  color: var(--color-text);
  text-align: center;
  line-height: 1.1;
  padding: 0.125rem;
  writing-mode: vertical-rl;
  text-orientation: mixed;
}

/* Mobile optimizations */
@media (max-width: 768px) {
  .discard-pile-container {
    top: 4rem;
    right: 0.5rem;
  }
  
  .discard-pile-stack {
    width: 1.5rem; /* Slightly smaller on mobile */
    height: 2.25rem;
    min-height: 2.25rem;
  }
  
  .discard-pile-empty .empty-pile-outline {
    width: 1.5rem;
    height: 2.25rem;
  }
  
  .empty-pile-text {
    font-size: 0.4rem;
  }
}

/* Animation for cards being added */
@keyframes discard-card-enter {
  0% {
    opacity: 0;
    transform: translate(var(--enter-x, 0), var(--enter-y, 0)) scale(1.2) rotate(var(--enter-rotation, 0deg));
  }
  60% {
    opacity: 0.8;
    transform: translate(var(--final-x), var(--final-y)) scale(1.1) rotate(var(--final-rotation));
  }
  100% {
    opacity: 1;
    transform: translate(var(--final-x), var(--final-y)) scale(1) rotate(var(--final-rotation));
  }
}

.discard-card.entering {
  animation: discard-card-enter 0.8s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Reduced motion support */
@media (prefers-reduced-motion: reduce) {
  .discard-card {
    transition: none !important;
  }
  
  .discard-card.entering {
    animation: none !important;
  }
  
  .discard-count-indicator {
    transition: none !important;
  }
}
</style>