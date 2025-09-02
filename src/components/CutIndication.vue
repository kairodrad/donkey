<template>
  <Teleport to="body">
    <Transition name="cut-indication" appear>
      <div 
        v-if="visible && cutData"
        class="cut-indication-overlay"
        @click="handleClick"
      >
        <!-- Main CUT Banner -->
        <div class="cut-banner">
          <div class="cut-banner-content">
            <!-- CUT Text with Swords -->
            <div class="cut-title">
              ⚔️ CUT! ⚔️
            </div>
            
            <!-- Cut Action Description -->
            <div class="cut-description">
              <div class="cut-player">
                {{ getCutPlayerName() }} played out of suit
              </div>
              <div class="cut-arrow">↓</div>
              <div class="cut-winner">
                {{ getWinnerName() }} gets {{ cutData.cardsCollected }} cards
              </div>
            </div>
            
            <!-- Visual Card Collection Animation -->
            <div class="card-collection-visual">
              <div 
                v-for="index in Math.min(cutData.cardsCollected || 0, 6)"
                :key="index"
                class="flying-card"
                :style="getFlyingCardStyle(index)"
              >
                🂠
              </div>
            </div>
            
            <!-- Timer Bar -->
            <div class="cut-timer-bar">
              <div 
                class="cut-timer-progress"
                :style="timerProgressStyle"
              ></div>
            </div>
          </div>
          
          <!-- Pulsing Effect Background -->
          <div class="cut-pulse-bg"></div>
        </div>
        
        <!-- Background Particles/Effects -->
        <div class="cut-particles">
          <div 
            v-for="index in 12"
            :key="index"
            class="particle"
            :style="getParticleStyle(index)"
          ></div>
        </div>
        
        <!-- Dismissal Hint -->
        <div class="cut-dismiss-hint">
          Click to dismiss or wait {{ Math.ceil(remainingTime / 1000) }}s
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'

export default {
  name: 'CutIndication',
  props: {
    visible: {
      type: Boolean,
      required: true
    },
    cutData: {
      type: Object,
      default: null
    },
    duration: {
      type: Number,
      default: 3000 // 3 seconds
    },
    dismissible: {
      type: Boolean,
      default: true
    }
  },
  emits: ['dismiss'],
  setup(props, { emit }) {
    const remainingTime = ref(props.duration)
    let timer = null
    let countdown = null
    
    // Timer progress for visual countdown
    const timerProgressStyle = computed(() => {
      const progress = (props.duration - remainingTime.value) / props.duration * 100
      return {
        width: `${progress}%`,
        transition: 'width 0.1s linear'
      }
    })
    
    // Get player names from cut data
    const getCutPlayerName = () => {
      return props.cutData?.cutPlayerName || 'Someone'
    }
    
    const getWinnerName = () => {
      return props.cutData?.winnerPlayerName || 'Someone'
    }
    
    // Flying card animation styles
    const getFlyingCardStyle = (index) => {
      const delay = index * 0.1 // Staggered animation
      const angle = (index * 60) - 30 // Spread cards in arc
      const distance = 50 + (index * 10) // Varying distances
      
      return {
        animationDelay: `${delay}s`,
        '--fly-angle': `${angle}deg`,
        '--fly-distance': `${distance}px`
      }
    }
    
    // Particle effect styles
    const getParticleStyle = (index) => {
      const angle = (index * 30) % 360
      const distance = 100 + Math.random() * 100
      const size = 4 + Math.random() * 8
      const duration = 1 + Math.random() * 2
      
      return {
        left: '50%',
        top: '50%',
        width: `${size}px`,
        height: `${size}px`,
        animationDuration: `${duration}s`,
        animationDelay: `${Math.random() * 0.5}s`,
        '--particle-angle': `${angle}deg`,
        '--particle-distance': `${distance}px`
      }
    }
    
    // Handle click to dismiss
    const handleClick = () => {
      if (props.dismissible) {
        dismiss()
      }
    }
    
    // Dismiss the indication
    const dismiss = () => {
      if (timer) {
        clearTimeout(timer)
        timer = null
      }
      if (countdown) {
        clearInterval(countdown)
        countdown = null
      }
      emit('dismiss')
    }
    
    // Start timers when component becomes visible
    const startTimers = () => {
      if (!props.visible) return
      
      remainingTime.value = props.duration
      
      // Auto-dismiss timer
      timer = setTimeout(() => {
        dismiss()
      }, props.duration)
      
      // Countdown timer for visual progress
      countdown = setInterval(() => {
        remainingTime.value = Math.max(0, remainingTime.value - 100)
      }, 100)
    }
    
    // Cleanup timers
    const cleanupTimers = () => {
      if (timer) {
        clearTimeout(timer)
        timer = null
      }
      if (countdown) {
        clearInterval(countdown)
        countdown = null
      }
    }
    
    // Watch for visibility changes
    const { watch } = Vue
    watch(() => props.visible, (visible) => {
      if (visible) {
        startTimers()
      } else {
        cleanupTimers()
      }
    }, { immediate: true })
    
    onMounted(() => {
      if (props.visible) {
        startTimers()
      }
    })
    
    onUnmounted(() => {
      cleanupTimers()
    })
    
    return {
      remainingTime,
      timerProgressStyle,
      getCutPlayerName,
      getWinnerName,
      getFlyingCardStyle,
      getParticleStyle,
      handleClick,
      dismiss
    }
  }
}
</script>

<style scoped>
.cut-indication-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(4px);
  cursor: pointer;
  user-select: none;
}

.cut-banner {
  position: relative;
  background: linear-gradient(135deg, #dc2626, #ef4444, #f87171);
  border: 3px solid #ffffff;
  border-radius: 1rem;
  padding: 2rem;
  box-shadow: 
    0 0 0 4px rgba(220, 38, 38, 0.5),
    0 20px 40px rgba(0, 0, 0, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
  transform-origin: center;
  animation: cut-banner-entrance 0.6s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  max-width: 90vw;
  min-width: 320px;
  overflow: hidden;
}

.cut-banner-content {
  position: relative;
  z-index: 2;
  text-align: center;
  color: white;
}

.cut-pulse-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
  animation: cut-pulse 1.5s infinite ease-in-out;
  pointer-events: none;
}

.cut-title {
  font-size: 3rem;
  font-weight: 900;
  text-shadow: 
    2px 2px 0 #991b1b,
    4px 4px 8px rgba(0, 0, 0, 0.5);
  margin-bottom: 1rem;
  animation: cut-title-glow 2s infinite ease-in-out;
}

.cut-description {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.cut-player, .cut-winner {
  font-size: 1.25rem;
  font-weight: 700;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5);
  padding: 0.5rem 1rem;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 0.5rem;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.cut-winner {
  background: rgba(255, 255, 255, 0.2);
  animation: cut-winner-highlight 1s infinite alternate ease-in-out;
}

.cut-arrow {
  font-size: 1.5rem;
  animation: cut-arrow-bounce 1s infinite ease-in-out;
}

.card-collection-visual {
  position: relative;
  height: 60px;
  margin: 1rem 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.flying-card {
  position: absolute;
  font-size: 1.5rem;
  color: white;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5);
  animation: flying-card-animation 1s ease-out forwards;
  opacity: 0;
}

.cut-timer-bar {
  width: 100%;
  height: 4px;
  background: rgba(255, 255, 255, 0.3);
  border-radius: 2px;
  overflow: hidden;
  margin-top: 1rem;
}

.cut-timer-progress {
  height: 100%;
  background: linear-gradient(90deg, #fbbf24, #f59e0b);
  border-radius: 2px;
  box-shadow: 0 0 8px rgba(245, 158, 11, 0.6);
}

.cut-particles {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  overflow: hidden;
}

.particle {
  position: absolute;
  background: #fbbf24;
  border-radius: 50%;
  animation: particle-explosion 2s ease-out forwards;
  opacity: 0;
}

.cut-dismiss-hint {
  position: absolute;
  bottom: 2rem;
  left: 50%;
  transform: translateX(-50%);
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.875rem;
  text-align: center;
  animation: cut-hint-fade 2s infinite ease-in-out;
}

/* Animations */
@keyframes cut-banner-entrance {
  0% {
    opacity: 0;
    transform: scale(0.3) rotate(-10deg);
  }
  60% {
    transform: scale(1.1) rotate(2deg);
  }
  100% {
    opacity: 1;
    transform: scale(1) rotate(0deg);
  }
}

@keyframes cut-pulse {
  0%, 100% {
    opacity: 0.1;
    transform: scale(1);
  }
  50% {
    opacity: 0.3;
    transform: scale(1.05);
  }
}

@keyframes cut-title-glow {
  0%, 100% {
    text-shadow: 
      2px 2px 0 #991b1b,
      4px 4px 8px rgba(0, 0, 0, 0.5),
      0 0 20px rgba(255, 255, 255, 0.3);
  }
  50% {
    text-shadow: 
      2px 2px 0 #991b1b,
      4px 4px 8px rgba(0, 0, 0, 0.5),
      0 0 30px rgba(255, 255, 255, 0.6);
  }
}

@keyframes cut-winner-highlight {
  0% {
    background: rgba(255, 255, 255, 0.2);
    transform: scale(1);
  }
  100% {
    background: rgba(255, 255, 255, 0.4);
    transform: scale(1.02);
  }
}

@keyframes cut-arrow-bounce {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-8px);
  }
}

@keyframes flying-card-animation {
  0% {
    opacity: 1;
    transform: translate(0, 0) rotate(0deg) scale(1);
  }
  50% {
    opacity: 1;
    transform: 
      translate(
        calc(cos(var(--fly-angle)) * var(--fly-distance)), 
        calc(sin(var(--fly-angle)) * var(--fly-distance))
      ) 
      rotate(720deg) 
      scale(1.2);
  }
  100% {
    opacity: 0;
    transform: 
      translate(
        calc(cos(var(--fly-angle)) * var(--fly-distance) * 0.8), 
        calc(sin(var(--fly-angle)) * var(--fly-distance) * 0.8)
      ) 
      rotate(1080deg) 
      scale(0.5);
  }
}

@keyframes particle-explosion {
  0% {
    opacity: 1;
    transform: translate(0, 0) scale(1);
  }
  100% {
    opacity: 0;
    transform: 
      translate(
        calc(cos(var(--particle-angle)) * var(--particle-distance)), 
        calc(sin(var(--particle-angle)) * var(--particle-distance))
      ) 
      scale(0);
  }
}

@keyframes cut-hint-fade {
  0%, 100% {
    opacity: 0.6;
  }
  50% {
    opacity: 1;
  }
}

/* Transition animations */
.cut-indication-enter-active {
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.cut-indication-leave-active {
  transition: all 0.3s ease-in;
}

.cut-indication-enter-from {
  opacity: 0;
  transform: scale(0.8);
}

.cut-indication-leave-to {
  opacity: 0;
  transform: scale(1.1);
}

/* Mobile optimizations */
@media (max-width: 768px) {
  .cut-banner {
    padding: 1.5rem 1rem;
    margin: 1rem;
    min-width: auto;
  }
  
  .cut-title {
    font-size: 2rem;
    margin-bottom: 0.75rem;
  }
  
  .cut-player, .cut-winner {
    font-size: 1rem;
    padding: 0.375rem 0.75rem;
  }
  
  .card-collection-visual {
    height: 40px;
  }
  
  .flying-card {
    font-size: 1.25rem;
  }
  
  .cut-dismiss-hint {
    bottom: 1rem;
    font-size: 0.75rem;
  }
}

/* Reduced motion support */
@media (prefers-reduced-motion: reduce) {
  .cut-banner {
    animation: none;
    transform: none;
  }
  
  .cut-pulse-bg {
    animation: none;
    opacity: 0.1;
  }
  
  .cut-title {
    animation: none;
  }
  
  .cut-winner {
    animation: none;
  }
  
  .cut-arrow {
    animation: none;
  }
  
  .flying-card {
    animation: none;
    opacity: 0;
  }
  
  .particle {
    animation: none;
    opacity: 0;
  }
  
  .cut-hint-fade {
    animation: none;
    opacity: 0.8;
  }
}

/* High contrast mode support */
@media (prefers-contrast: high) {
  .cut-banner {
    background: #dc2626;
    border-color: #ffffff;
    box-shadow: 0 0 0 4px #000000;
  }
  
  .cut-player, .cut-winner {
    background: rgba(0, 0, 0, 0.8);
    border-color: #ffffff;
  }
  
  .cut-timer-progress {
    background: #fbbf24;
    box-shadow: none;
  }
}
</style>