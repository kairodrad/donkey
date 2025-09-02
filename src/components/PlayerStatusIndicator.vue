<template>
  <div 
    :class="playerStatusClasses"
    :style="playerStatusStyle"
  >
    <!-- Player Name -->
    <div class="player-name">
      {{ player.name }}
      <span v-if="isUser" class="user-indicator">(You)</span>
      <span v-if="player.isBot" class="bot-indicator">🤖</span>
    </div>
    
    <!-- Connection Status Indicator -->
    <div 
      :class="connectionStatusClasses"
      :title="connectionStatusTitle"
    >
      <div class="status-dot"></div>
    </div>
    
    <!-- DONKEY Letters Progress -->
    <div 
      v-if="player.donkeyLetters"
      class="donkey-letters"
      :title="donkeyProgressTitle"
    >
      <span 
        v-for="(letter, index) in donkeyLettersArray"
        :key="index"
        :class="donkeyLetterClasses(index)"
      >
        {{ letter }}
      </span>
    </div>
    
    <!-- Turn Indicator -->
    <div 
      v-if="isCurrentTurn"
      class="turn-indicator"
      :class="turnIndicatorClasses"
    >
      <div class="turn-arrow">▼</div>
      <div class="turn-pulse"></div>
    </div>
    
    <!-- Card Count (for opponents) -->
    <div 
      v-if="!isUser && cardCount !== null"
      class="card-count"
      :title="`${player.name} has ${cardCount} cards`"
    >
      {{ cardCount }}
    </div>
    
    <!-- Finished Round Indicator -->
    <div 
      v-if="isFinishedThisRound"
      class="finished-indicator"
      title="Finished this round"
    >
      ✓
    </div>
    
    <!-- Disconnection Warning -->
    <div 
      v-if="!player.isConnected && !player.isBot"
      class="disconnection-warning"
      title="Player disconnected"
    >
      ⚠️
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'PlayerStatusIndicator',
  props: {
    player: {
      type: Object,
      required: true
    },
    isUser: {
      type: Boolean,
      default: false
    },
    isCurrentTurn: {
      type: Boolean,
      default: false
    },
    cardCount: {
      type: Number,
      default: null
    },
    isFinishedThisRound: {
      type: Boolean,
      default: false
    },
    position: {
      type: String,
      default: 'bottom', // 'top', 'bottom', 'left', 'right'
      validator: (value) => ['top', 'bottom', 'left', 'right'].includes(value)
    }
  },
  setup(props) {
    // Player status container classes
    const playerStatusClasses = computed(() => {
      const classes = ['player-status-indicator']
      
      if (props.isUser) {
        classes.push('player-status-user')
      } else {
        classes.push('player-status-opponent')
      }
      
      if (props.isCurrentTurn) {
        classes.push('player-status-current-turn')
      }
      
      if (props.isFinishedThisRound) {
        classes.push('player-status-finished')
      }
      
      if (!props.player.isConnected && !props.player.isBot) {
        classes.push('player-status-disconnected')
      }
      
      classes.push(`player-status-${props.position}`)
      
      return classes.join(' ')
    })
    
    // Dynamic positioning style
    const playerStatusStyle = computed(() => {
      // This would be overridden by parent component for exact positioning
      return {}
    })
    
    // Connection status classes
    const connectionStatusClasses = computed(() => {
      const classes = ['connection-status']
      
      if (props.player.isBot) {
        classes.push('connection-bot')
      } else if (props.player.isConnected) {
        classes.push('connection-online')
      } else {
        classes.push('connection-offline')
      }
      
      return classes.join(' ')
    })
    
    // Connection status tooltip
    const connectionStatusTitle = computed(() => {
      if (props.player.isBot) {
        return `Bot (${props.player.botDifficulty || 'easy'})`
      }
      return props.player.isConnected ? 'Connected' : 'Disconnected'
    })
    
    // Turn indicator classes with animation
    const turnIndicatorClasses = computed(() => {
      const classes = ['turn-indicator-active']
      
      if (props.isUser) {
        classes.push('turn-indicator-user')
      } else {
        classes.push('turn-indicator-opponent')
      }
      
      return classes.join(' ')
    })
    
    // DONKEY letters array for display
    const donkeyLettersArray = computed(() => {
      return props.player.donkeyLetters.split('')
    })
    
    // DONKEY progress tooltip
    const donkeyProgressTitle = computed(() => {
      const remaining = 6 - props.player.donkeyLetters.length
      if (remaining === 0) {
        return 'DONKEY - Game Loser!'
      }
      return `${props.player.donkeyLetters.length}/6 letters - ${remaining} away from DONKEY`
    })
    
    // Individual DONKEY letter classes
    const donkeyLetterClasses = (index) => {
      const classes = ['donkey-letter']
      
      // Add color based on progress
      if (index === props.player.donkeyLetters.length - 1) {
        classes.push('donkey-letter-latest') // Most recent letter
      }
      
      // Different colors based on danger level
      const progress = props.player.donkeyLetters.length
      if (progress >= 5) {
        classes.push('donkey-letter-danger') // Red - very close
      } else if (progress >= 3) {
        classes.push('donkey-letter-warning') // Orange - warning
      } else {
        classes.push('donkey-letter-safe') // Blue - safe
      }
      
      return classes.join(' ')
    }
    
    return {
      playerStatusClasses,
      playerStatusStyle,
      connectionStatusClasses,
      connectionStatusTitle,
      turnIndicatorClasses,
      donkeyLettersArray,
      donkeyProgressTitle,
      donkeyLetterClasses
    }
  }
}
</script>

<style scoped>
.player-status-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
  padding: 0.5rem;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 0.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: relative;
  min-width: 4rem;
  text-align: center;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* User vs Opponent styling */
.player-status-user {
  border-color: var(--color-primary);
  background: var(--color-background-secondary);
}

.player-status-opponent {
  background: var(--color-surface);
}

/* Current turn highlighting */
.player-status-current-turn {
  background: rgba(59, 130, 246, 0.1);
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2), 0 4px 12px rgba(59, 130, 246, 0.3);
  transform: scale(1.05);
}

/* Finished player styling */
.player-status-finished {
  opacity: 0.7;
  background: rgba(16, 185, 129, 0.1);
  border-color: #10b981;
}

/* Disconnected player styling */
.player-status-disconnected {
  opacity: 0.6;
  background: rgba(239, 68, 68, 0.1);
  border-color: #ef4444;
}

/* Player name */
.player-name {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text);
  line-height: 1.2;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-indicator {
  font-size: 0.625rem;
  color: var(--color-primary);
  font-weight: 500;
}

.bot-indicator {
  font-size: 0.75rem;
  opacity: 0.8;
}

/* Connection status dot */
.connection-status {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  transition: all 0.3s ease;
}

.connection-online .status-dot {
  background: #10b981; /* Green */
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.3);
}

.connection-offline .status-dot {
  background: #ef4444; /* Red */
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.3);
  animation: connection-warning 2s infinite;
}

.connection-bot .status-dot {
  background: #8b5cf6; /* Purple */
  box-shadow: 0 0 0 2px rgba(139, 92, 246, 0.3);
}

@keyframes connection-warning {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* DONKEY letters */
.donkey-letters {
  display: flex;
  gap: 0.125rem;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  max-width: 100%;
}

.donkey-letter {
  font-size: 0.625rem;
  font-weight: 700;
  padding: 0.125rem 0.25rem;
  border-radius: 0.25rem;
  line-height: 1;
  text-align: center;
  min-width: 0.875rem;
  transition: all 0.2s ease;
}

.donkey-letter-safe {
  background: rgba(59, 130, 246, 0.2);
  color: #1d4ed8;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.donkey-letter-warning {
  background: rgba(245, 158, 11, 0.2);
  color: #b45309;
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.donkey-letter-danger {
  background: rgba(239, 68, 68, 0.2);
  color: #dc2626;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.donkey-letter-latest {
  animation: letter-pulse 2s infinite;
  transform: scale(1.1);
}

@keyframes letter-pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

/* Turn indicator */
.turn-indicator {
  position: absolute;
  top: -1.5rem;
  left: 50%;
  transform: translateX(-50%);
  z-index: 10;
}

.turn-indicator-active .turn-arrow {
  color: var(--color-primary);
  font-size: 1rem;
  font-weight: bold;
  animation: turn-bounce 1s infinite;
  filter: drop-shadow(0 2px 4px rgba(59, 130, 246, 0.4));
}

.turn-indicator-user .turn-arrow {
  color: #10b981; /* Green for user */
  filter: drop-shadow(0 2px 4px rgba(16, 185, 129, 0.4));
}

@keyframes turn-bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}

.turn-pulse {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 20px;
  height: 20px;
  background: rgba(59, 130, 246, 0.3);
  border-radius: 50%;
  transform: translate(-50%, -50%);
  animation: turn-pulse-animation 1.5s infinite;
  pointer-events: none;
}

@keyframes turn-pulse-animation {
  0% {
    transform: translate(-50%, -50%) scale(0.8);
    opacity: 1;
  }
  100% {
    transform: translate(-50%, -50%) scale(2);
    opacity: 0;
  }
}

/* Card count */
.card-count {
  font-size: 0.625rem;
  font-weight: 600;
  color: var(--color-text);
  background: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: 0.375rem;
  padding: 0.125rem 0.375rem;
  min-width: 1.25rem;
  text-align: center;
}

/* Finished round indicator */
.finished-indicator {
  position: absolute;
  top: -0.25rem;
  right: -0.25rem;
  width: 1rem;
  height: 1rem;
  background: #10b981;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.625rem;
  font-weight: bold;
  border: 2px solid var(--color-surface);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

/* Disconnection warning */
.disconnection-warning {
  position: absolute;
  top: -0.25rem;
  left: -0.25rem;
  width: 1rem;
  height: 1rem;
  background: #ef4444;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.5rem;
  border: 2px solid var(--color-surface);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  animation: disconnection-pulse 2s infinite;
}

@keyframes disconnection-pulse {
  0%, 100% { 
    background: #ef4444;
    transform: scale(1);
  }
  50% { 
    background: #dc2626;
    transform: scale(1.1);
  }
}

/* Position-specific styles */
.player-status-bottom {
  margin-bottom: 1rem;
}

.player-status-top .turn-indicator {
  top: auto;
  bottom: -1.5rem;
}

.player-status-top .turn-indicator .turn-arrow {
  transform: rotate(180deg);
}

/* Mobile optimizations */
@media (max-width: 768px) {
  .player-status-indicator {
    padding: 0.375rem;
    gap: 0.125rem;
    min-width: 3rem;
  }
  
  .player-name {
    font-size: 0.625rem;
  }
  
  .status-dot {
    width: 6px;
    height: 6px;
  }
  
  .donkey-letter {
    font-size: 0.5rem;
    padding: 0.0625rem 0.125rem;
    min-width: 0.75rem;
  }
  
  .card-count {
    font-size: 0.5rem;
    padding: 0.0625rem 0.25rem;
  }
  
  .turn-indicator .turn-arrow {
    font-size: 0.875rem;
  }
}

/* Reduced motion support */
@media (prefers-reduced-motion: reduce) {
  .player-status-current-turn {
    transform: none;
  }
  
  .turn-indicator .turn-arrow {
    animation: none;
  }
  
  .turn-pulse {
    animation: none;
    opacity: 0.3;
  }
  
  .donkey-letter-latest {
    animation: none;
    transform: none;
  }
  
  .connection-offline .status-dot {
    animation: none;
  }
  
  .disconnection-warning {
    animation: none;
  }
}
</style>