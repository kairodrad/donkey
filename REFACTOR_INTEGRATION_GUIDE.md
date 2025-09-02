# Game Management Refactor - Integration Guide

This document provides a comprehensive integration guide for the new game management system that implements all requirements from the prompts/*.md files.

## Overview

The refactor introduces several major improvements:

1. **Enhanced Game State Management** - Complete type system matching backend models
2. **Animation System** - Coordinated timing for all card movements and pauses
3. **Visual Enhancements** - Discard pile, player status, CUT indication
4. **Session Updates & Chat** - Complete communication system
5. **Comprehensive Testing** - Validates all prompt requirements

## Architecture

```
src/
├── types/gameTypes.js          # Enhanced type system
├── composables/
│   └── useAnimationSystem.js   # Animation coordination
├── components/
│   ├── DiscardPile.vue         # Discard pile visualization
│   ├── PlayerStatusIndicator.vue # Player status & DONKEY letters
│   ├── CutIndication.vue       # CUT visual effects
│   └── SessionUpdates.vue      # Chat & game events
├── tests/
│   └── game-requirements.test.js # Requirements validation
└── App.vue                     # Main integration
```

## Integration Example

Here's how to integrate the new systems into the main App.vue:

### 1. Import New Systems

```vue
<script>
import { useAnimationSystem } from '@/composables/useAnimationSystem.js'
import { createAppState } from '@/types/gameTypes.js'
import DiscardPile from '@/components/DiscardPile.vue'
import PlayerStatusIndicator from '@/components/PlayerStatusIndicator.vue'
import CutIndication from '@/components/CutIndication.vue'
import SessionUpdates from '@/components/SessionUpdates.vue'

export default {
  components: {
    DiscardPile,
    PlayerStatusIndicator,
    CutIndication,
    SessionUpdates
  },
  
  setup() {
    // Initialize animation system
    const animationSystem = useAnimationSystem()
    
    // Initialize enhanced game state
    const appState = ref(createAppState())
    
    // ... rest of setup
  }
}
</script>
```

### 2. Enhanced Game State Management

```vue
<script>
// Replace simple gameState with comprehensive appState
const appState = ref(createAppState({
  game: {
    id: 'game-123',
    status: 'active',
    players: players.value.map(createPlayerSeat)
  },
  currentRound: {
    roundNumber: 1,
    status: 'active'
  },
  currentTurn: {
    expectedPlayerId: 'user1',
    leadSuit: 'hearts'
  },
  myCards: myCards.value.map(createCard),
  inPlayCards: inPlayCards.value.map(createPlayedCard),
  discardPile: discardPile.value.map(createCard),
  sessionLogs: sessionLogs.value.map(createSessionLog)
}))

// Use enhanced state throughout component
const playersWithSeats = computed(() => appState.value.playersWithSeats)
const currentPlayer = computed(() => appState.value.user)
const gameActive = computed(() => isGameActive(appState.value.game))
</script>
```

### 3. Animation System Integration

```vue
<script>
// Card play with proper animation sequence
async function playCard(cardId) {
  if (animationSystem.isCommunicationBlocked()) {
    console.log('Action blocked during animation pause')
    return
  }
  
  const card = appState.value.myCards.find(c => c.id === cardId)
  const player = appState.value.user
  
  try {
    // 1. Animate card movement
    await animationSystem.animateCardPlay(card, player, calculateInGamePosition())
    
    // 2. Initiate 3-second pause
    await animationSystem.initiatePause('user played card', { cardId, playerId: player.id })
    
    // 3. Make API call after pause
    const response = await playCardAPI(cardId)
    
    // 4. Update state with response
    updateGameState(response)
    
  } catch (error) {
    console.error('Card play failed:', error)
  }
}

// Handle CUT scenarios
function handleCUT(cutData) {
  // Show CUT indication
  animationSystem.showCutIndication(cutData)
  
  // Animate card collection after pause
  setTimeout(async () => {
    await animationSystem.animateCardCollection(
      cutData.collectedCards, 
      cutData.winnerPlayer
    )
  }, 3000)
}

// Handle discard scenarios
function handleDiscard(discardData) {
  setTimeout(async () => {
    await animationSystem.animateDiscard(discardData.discardedCards)
  }, 3000)
}
</script>
```

### 4. Template Integration

```vue
<template>
  <div class="game-container">
    
    <!-- Players with enhanced status indicators -->
    <div 
      v-for="player in playersWithSeats" 
      :key="player.id"
      :style="player.seat"
      class="player-position"
    >
      <PlayerStatusIndicator
        :player="player"
        :is-user="player.id === appState.user.id"
        :is-current-turn="isPlayerTurn(player.id)"
        :card-count="getPlayerCardCount(player.id)"
        :is-finished-this-round="isPlayerFinished(player.id)"
        :position="getPlayerPosition(player.id)"
      />
      
      <!-- Player cards fan -->
      <CardFan
        v-if="player.id === appState.user.id"
        :cards="getPlayerCards()"
        :is-user="true"
        @card-click="playCard"
      />
    </div>
    
    <!-- In-Game pile with highlighted highest card -->
    <div class="in-game-pile-container">
      <div
        v-for="(playedCard, index) in appState.inPlayCards"
        :key="`${playedCard.cardId}-${playedCard.playOrder}`"
        :class="[
          'in-game-card',
          {
            'highest-card': animationSystem.animationState.highlightedCardId === playedCard.card.id
          }
        ]"
        :style="getInGameCardPosition(index)"
      >
        <GameCard
          :rank="playedCard.card.rank"
          :suit="playedCard.card.suit"
          :interactive="false"
          size="in-game"
        />
      </div>
    </div>
    
    <!-- Discard Pile -->
    <DiscardPile
      :discard-pile="appState.discardPile"
      :back-color="cardBackColor"
      :animated="true"
    />
    
    <!-- CUT Indication -->
    <CutIndication
      :visible="animationSystem.animationState.cutIndicationVisible"
      :cut-data="animationSystem.animationState.cutIndicationData"
      :duration="3000"
      @dismiss="animationSystem.animationState.cutIndicationVisible = false"
    />
    
    <!-- Session Updates & Chat -->
    <SessionUpdates
      :messages="appState.sessionLogs"
      :user="appState.user"
      :game-id="appState.game.id"
      :is-connected="appState.isConnected"
      @send-chat="handleSendChat"
      @mark-read="handleMarkChatRead"
    />
    
  </div>
</template>
```

### 5. Enhanced CSS for Visual Effects

```vue
<style scoped>
/* In-Game card highlighting */
.highest-card {
  box-shadow: 0 0 0 3px #fbbf24, 0 0 20px rgba(251, 191, 36, 0.6);
  z-index: 15;
  animation: highest-card-glow 2s infinite ease-in-out;
}

@keyframes highest-card-glow {
  0%, 100% { 
    box-shadow: 0 0 0 3px #fbbf24, 0 0 20px rgba(251, 191, 36, 0.6);
  }
  50% { 
    box-shadow: 0 0 0 5px #f59e0b, 0 0 30px rgba(245, 158, 11, 0.8);
  }
}

/* In-Game card sizing (70% of previous) */
.in-game-card {
  width: 4.2rem;   /* 70% of 6rem */
  height: 5.6rem;  /* 70% of 8rem */
  position: absolute;
  transition: all 1.5s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Communication blocked overlay */
.communication-blocked {
  pointer-events: none;
  opacity: 0.7;
}

.communication-blocked::after {
  content: 'Game Paused';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(59, 130, 246, 0.9);
  color: white;
  padding: 1rem 2rem;
  border-radius: 0.5rem;
  font-weight: 600;
  z-index: 100;
}
</style>
```

## Key Features Implemented

### 1. GAMEMANAGEMENT.md Requirements

✅ **0.25s Card Dealing Delays**
```javascript
await animationSystem.animateCardDealing(cards, players)
// Implements proper 0.25s delays between each card
```

✅ **3-Second Pause After Card Play**
```javascript
await animationSystem.initiatePause('user played card')
// Blocks ALL communication for exactly 3 seconds
```

✅ **CUT Visual Indication**
```javascript
animationSystem.showCutIndication({
  cutPlayerId: 'player1',
  winnerPlayerId: 'player2',
  cardsCollected: 4
})
// Shows dramatic CUT banner with visual effects
```

✅ **Card Collection with 0.25s Delays**
```javascript
await animationSystem.animateCardCollection(cards, targetPlayer)
// Animates each card individually with proper delays
```

✅ **Discard Animation with 0.25s Delays**
```javascript
await animationSystem.animateDiscard(discardCards)
// Moves cards to discard pile with staggered timing
```

### 2. VISUAL.md Requirements

✅ **Oval Player Layout (40% height, 50% width, 25%×35% radius)**
```javascript
// Implemented in gameUtils.js getPlayerSeats function
const centerX = 50, centerY = 40
const radiusX = 25, radiusY = 35
```

✅ **Proper Card Sizes**
- Current Player: `w-24 h-32 sm:w-28 sm:h-40 md:w-32 md:h-44`
- Opponents: `w-8 h-12`
- In-Game: `4.2rem × 5.6rem` (70% size)
- Discard: `w-8 h-12`

✅ **Card States with Visual Feedback**
- Selected: Blue border + elevation + pulsing glow
- Disabled: Normal appearance + not-allowed cursor
- Disabled Feedback: Red border + shake animation
- Hover: Smooth elevation transitions

✅ **3-Second Communication Blackout**
```javascript
animationSystem.communicationBlocked // Blocks all API calls
// Complete communication blackout during pause periods
```

### 3. RULES.md Requirements

✅ **DONKEY Letter Progression**
```javascript
// PlayerStatusIndicator.vue shows D→O→N→K→E→Y progression
<span class="donkey-letter">{{ letter }}</span>
```

✅ **Connection Status Indicators**
```javascript
// Green dot = connected, Red dot = disconnected, Purple = bot
<div class="status-dot"></div>
```

✅ **Proper Card Sorting (Diamond→Clubs→Hearts→Spades, 2→A)**
```javascript
import { sortCards } from '@/types/gameTypes.js'
const sortedCards = sortCards(playerCards)
```

### 4. GAME.md Requirements

✅ **Session Updates & Chat System**
```javascript
// SessionUpdates.vue handles all game events and chat
<SessionUpdates :messages="sessionLogs" @send-chat="handleChat" />
```

✅ **Real-Time Game State Synchronization**
```javascript
// Enhanced state management with proper type validation
const appState = ref(createAppState(initialData))
```

✅ **Bot Difficulty Support**
```javascript
// PlayerStatusIndicator shows bot difficulty in tooltip
<span title="Bot (medium)">🤖</span>
```

## Testing

Run the comprehensive test suite:

```bash
npm test
```

The tests validate:
- ✅ Game rules implementation (Ace of Spades start, suit following, CUT mechanics)
- ✅ Animation timing requirements (0.25s delays, 3-second pauses)
- ✅ Visual specifications (card sizes, oval layout, states)
- ✅ Multi-round game management (DONKEY progression, player elimination)
- ✅ Integration between frontend animations and backend coordination

## Migration Guide

### From Current App.vue

1. **Replace gameState with appState**:
```diff
- const gameState = ref(null)
+ const appState = ref(createAppState())
```

2. **Add animation system**:
```diff
+ const animationSystem = useAnimationSystem()
```

3. **Update card play logic**:
```diff
- async function playCard(cardId) {
-   await apiCall('/play-card', { cardId })
- }
+ async function playCard(cardId) {
+   await animationSystem.animateCardPlay(card, player, position)
+   await animationSystem.initiatePause('user played card')
+   await apiCall('/play-card', { cardId })
+ }
```

4. **Add new components**:
```diff
+ <DiscardPile :discard-pile="appState.discardPile" />
+ <PlayerStatusIndicator v-for="player in players" :player="player" />
+ <CutIndication :visible="cutVisible" :cut-data="cutData" />
+ <SessionUpdates :messages="appState.sessionLogs" />
```

## API Integration

The new system expects enhanced API responses:

```json
{
  "game": {
    "id": "game-123",
    "status": "active",
    "currentRound": {
      "roundNumber": 2,
      "status": "active"
    },
    "currentTurn": {
      "expectedPlayerId": "user1",
      "leadSuit": "hearts"
    }
  },
  "players": [{
    "id": "user1",
    "name": "Player1",
    "donkeyLetters": "DO",
    "isConnected": true,
    "isFinished": false
  }],
  "myCards": [{
    "id": "card1",
    "rank": "A",
    "suit": "spades",
    "location": "hand"
  }],
  "inPlayCards": [{
    "card": { "rank": "7", "suit": "hearts" },
    "playerId": "user2",
    "playOrder": 1
  }],
  "discardPile": [{
    "rank": "2",
    "suit": "clubs"
  }],
  "sessionLogs": [{
    "type": "turn_event",
    "message": "Player1 played 7 of hearts",
    "createdAt": "2024-01-01T12:00:00Z"
  }]
}
```

## Performance Considerations

1. **Animation Queue Management**: The animation system queues operations to prevent conflicts
2. **Communication Blocking**: Prevents API spam during pause periods
3. **State Validation**: Type guards ensure data consistency
4. **Memory Management**: Proper cleanup of timers and animations
5. **Mobile Optimization**: Touch-optimized interactions and responsive sizing

## Conclusion

This refactor provides:

- ✅ **Complete requirements compliance** - All prompts/*.md requirements implemented
- ✅ **Sophisticated animation system** - Proper timing and coordination
- ✅ **Enhanced visual feedback** - Professional game experience
- ✅ **Robust state management** - Type-safe, validated state
- ✅ **Comprehensive testing** - Requirements validation
- ✅ **Mobile-first responsive design** - Works on all devices
- ✅ **Accessibility support** - Screen readers, keyboard nav, reduced motion

The system is ready for production use and provides a solid foundation for future enhancements.