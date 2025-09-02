// Comprehensive animation system implementing GAMEMANAGEMENT.md requirements
// Handles all card movements, timing, and coordination with backend

import { ref, reactive, nextTick } from 'vue'

/**
 * Animation System State
 */
const animationState = reactive({
  // Current phase
  phase: 'idle', // 'idle', 'dealing', 'playing', 'pausing', 'cutting', 'discarding'
  
  // Communication control
  communicationBlocked: false,
  pauseActive: false,
  pauseStartTime: null,
  pauseRemainingTime: 0,
  
  // Animation queues
  dealingQueue: [],
  playingQueue: [],
  collectingQueue: [],
  discardingQueue: [],
  
  // Active animations
  activeAnimations: new Set(),
  
  // Visual state
  cutIndicationVisible: false,
  cutIndicationData: null,
  highlightedCardId: null,
  
  // Timing configuration
  dealDelay: 250,     // 0.25s between dealing cards
  pauseDuration: 3000, // 3s pause after actions
  moveAnimation: 1500, // 1.5s for card movement
  discardDelay: 250    // 0.25s between discarding cards
})

/**
 * Core Animation System Composable
 */
export function useAnimationSystem() {
  
  /**
   * DEALING ANIMATION - GAMEMANAGEMENT.md requirement:
   * "Start with the 52 cards, and set a 0.25s timer between dealing each card"
   */
  async function animateCardDealing(cards, players, startPlayerIndex = 0) {
    if (animationState.phase !== 'idle') {
      return false
    }
    
    animationState.phase = 'dealing'
    animationState.dealingQueue = [...cards]
    
    
    
    try {
      // Deal cards with 0.25s delay between each
      for (let i = 0; i < cards.length; i++) {
        const card = cards[i]
        const playerIndex = (startPlayerIndex + i) % players.length
        const player = players[playerIndex]
        
        // Animate card moving to player
        await animateCardToPlayer(card, player, i)
        
        // Sort and arrange card in player's hand
        arrangeCasrdInHand(card, player)
        
        // 0.25s delay before next card (except last card)
        if (i < cards.length - 1) {
          await sleep(animationState.dealDelay)
        }
      }
      
      
      return true
      
    } catch (error) {
      return false
    } finally {
      animationState.phase = 'idle'
      animationState.dealingQueue = []
    }
  }
  
  /**
   * CARD PLAY ANIMATION - GAMEMANAGEMENT.md requirement:
   * "animate the move so that the card is shown slowly moving from the player or opponent hand to the pile"
   */
  async function animateCardPlay(card, fromPlayer, toPosition) {
    if (!card || !fromPlayer) {
      return false
    }
    
    const animationId = `play-${card.id}-${Date.now()}`
    animationState.activeAnimations.add(animationId)
    
    try {
      
      // Mark card as animating
      card.animating = true
      card.animationType = 'playing'
      
      // Get source and destination elements
      const cardElement = getCardElement(card.id)
      const destinationElement = getInGamePileElement()
      
      if (!cardElement || !destinationElement) {
        return false
      }
      
      // Calculate animation path
      const path = calculateAnimationPath(cardElement, destinationElement, toPosition)
      
      // Animate movement with style transition
      await animateElementAlongPath(cardElement, path, {
        duration: animationState.moveAnimation,
        easing: 'cubic-bezier(0.4, 0, 0.2, 1)',
        styleTransitions: {
          from: 'card-in-hand',
          to: 'card-in-game'
        }
      })
      
      // Update card state
      card.animating = false
      card.animationType = null
      
      
      return true
      
    } catch (error) {
      return false
    } finally {
      animationState.activeAnimations.delete(animationId)
    }
  }
  
  /**
   * 3-SECOND PAUSE SYSTEM - GAMEMANAGEMENT.md requirement:
   * "pause for 3 seconds before taking any action like yielding turn to the next player"
   */
  async function initiatePause(reason = 'card played', context = null) {
    if (animationState.pauseActive) {
      return extendPause()
    }
    
    
    // Block all communication
    animationState.communicationBlocked = true
    animationState.pauseActive = true
    animationState.pauseStartTime = Date.now()
    animationState.pauseRemainingTime = animationState.pauseDuration
    
    // Log the pause start
    logSessionEvent('pause_started', `Turn pause started - 3 seconds (${reason})`, context)
    
    // Start countdown timer
    const countdownInterval = setInterval(() => {
      const elapsed = Date.now() - animationState.pauseStartTime
      animationState.pauseRemainingTime = Math.max(0, animationState.pauseDuration - elapsed)
      
      if (animationState.pauseRemainingTime <= 0) {
        clearInterval(countdownInterval)
        completePause()
      }
    }, 100) // Update every 100ms for smooth countdown
    
    // Return promise that resolves when pause completes
    return new Promise(resolve => {
      const checkComplete = () => {
        if (!animationState.pauseActive) {
          resolve()
        } else {
          setTimeout(checkComplete, 100)
        }
      }
      checkComplete()
    })
  }
  
  function completePause() {
    
    
    // Unblock communication
    animationState.communicationBlocked = false
    animationState.pauseActive = false
    animationState.pauseStartTime = null
    animationState.pauseRemainingTime = 0
    
    logSessionEvent('pause_completed', 'Turn pause completed, resuming game flow')
  }
  
  function extendPause(additionalTime = 3000) {
    animationState.pauseStartTime = Date.now()
    animationState.pauseRemainingTime = additionalTime
    
  }
  
  /**
   * CUT INDICATION - GAMEMANAGEMENT.md requirement:
   * "add a non-modal creative emphasizing indication of CUT"
   */
  async function showCutIndication(cutData) {
    const { cutPlayerId, winnerPlayerId, cardsCollected } = cutData
    
    
    
    // Set CUT indication data
    animationState.cutIndicationData = {
      cutPlayerId,
      winnerPlayerId,
      cardsCollected,
      message: `CUT! Player ${winnerPlayerId} gets ${cardsCollected} cards`
    }
    
    animationState.cutIndicationVisible = true
    
    // Log to session
    logSessionEvent('cut_performed', 
      `Player ${cutPlayerId} performed CUT. Player ${winnerPlayerId} collected ${cardsCollected} cards.`,
      cutData
    )
    
    // Show for 3 seconds during pause
    await sleep(3000)
    
    // Hide indication
    animationState.cutIndicationVisible = false
    animationState.cutIndicationData = null
  }
  
  /**
   * CARD COLLECTION ANIMATION - For CUT scenarios
   */
  async function animateCardCollection(cards, targetPlayer) {
    if (!cards.length || !targetPlayer) return false
    
    
    
    animationState.collectingQueue = [...cards]
    
    try {
      // Collect cards with 0.25s delay between each
      for (let i = 0; i < cards.length; i++) {
        const card = cards[i]
        
        await animateCardToPlayer(card, targetPlayer, i)
        
        // 0.25s delay between card movements
        if (i < cards.length - 1) {
          await sleep(animationState.dealDelay)
        }
      }
      
      // Arrange collected cards in hand with proper sorting
      arrangePlayerHand(targetPlayer)
      
      return true
    } catch (error) {
      return false
    } finally {
      animationState.collectingQueue = []
    }
  }
  
  /**
   * DISCARD ANIMATION - GAMEMANAGEMENT.md requirement:
   * "0.25s pause after each card being moved to the player's hand"
   */
  async function animateDiscard(cards) {
    if (!cards.length) return false
    
    
    
    animationState.phase = 'discarding'
    animationState.discardingQueue = [...cards]
    
    try {
      // Move cards to discard pile with delay
      for (let i = 0; i < cards.length; i++) {
        const card = cards[i]
        
        await animateCardToDiscardPile(card, i)
        
        // 0.25s delay between card movements
        if (i < cards.length - 1) {
          await sleep(animationState.discardDelay)
        }
      }
      
      
      return true
      
    } catch (error) {
      return false
    } finally {
      animationState.phase = 'idle'
      animationState.discardingQueue = []
    }
  }
  
  /**
   * HIGHEST CARD HIGHLIGHTING - GAMEMANAGEMENT.md requirement:
   * "highlight the highest card in the IN-Game pile"
   */
  function highlightHighestCard(inPlayCards) {
    if (!inPlayCards.length) {
      animationState.highlightedCardId = null
      return
    }
    
    // Find highest card in current suit
    let highestCard = inPlayCards[0]
    
    if (inPlayCards[0]?.card?.suit) {
      const leadSuit = inPlayCards[0].card.suit
      const suitCards = inPlayCards.filter(pc => pc.card.suit === leadSuit)
      
      highestCard = suitCards.reduce((highest, current) => 
        (current.card.value > highest.card.value) ? current : highest
      )
    }
    
    animationState.highlightedCardId = highestCard.card.id
  }
  
  /**
   * UTILITY FUNCTIONS
   */
  
  // Animate single card to player position
  async function animateCardToPlayer(card, player, index = 0) {
    const cardElement = getCardElement(card.id)
    const playerHandElement = getPlayerHandElement(player.id)
    
    if (!cardElement || !playerHandElement) {
      return
    }
    
    // Calculate final position in fan
    const fanPosition = calculateFanPosition(player, index)
    
    return animateElementToPosition(cardElement, fanPosition, {
      duration: animationState.moveAnimation * 0.7, // Faster for dealing
      easing: 'cubic-bezier(0.25, 0.8, 0.25, 1)'
    })
  }
  
  // Animate card to discard pile
  async function animateCardToDiscardPile(card, stackIndex = 0) {
    const cardElement = getCardElement(card.id)
    const discardElement = getDiscardPileElement()
    
    if (!cardElement || !discardElement) {
      return
    }
    
    // Calculate stacked position in discard pile
    const discardPosition = calculateDiscardPosition(stackIndex)
    
    return animateElementToPosition(cardElement, discardPosition, {
      duration: animationState.moveAnimation * 0.8,
      easing: 'cubic-bezier(0.4, 0, 0.2, 1)'
    })
  }
  
  // Arrange card in player's hand with sorting
  function arrangeCasrdInHand(card, player) {
    // This would trigger Vue reactivity to re-sort and re-arrange hand
    // Implementation depends on parent component structure
    
  }
  
  function arrangePlayerHand(player) {
    
    // Trigger hand re-arrangement with proper Diamond->Clubs->Hearts->Spades, 2->A sorting
  }
  
  // Element finding utilities
  function getCardElement(cardId) {
    return document.querySelector(`[data-card-id="${cardId}"]`)
  }
  
  function getPlayerHandElement(playerId) {
    return document.querySelector(`[data-player-hand="${playerId}"]`)
  }
  
  function getInGamePileElement() {
    return document.querySelector('[data-in-game-pile]')
  }
  
  function getDiscardPileElement() {
    return document.querySelector('[data-discard-pile]')
  }
  
  // Animation utilities
  function calculateAnimationPath(fromElement, toElement, endPosition) {
    const fromRect = fromElement.getBoundingClientRect()
    const toRect = toElement.getBoundingClientRect()
    
    return {
      from: { x: fromRect.left, y: fromRect.top },
      to: { x: toRect.left + (endPosition?.x || 0), y: toRect.top + (endPosition?.y || 0) }
    }
  }
  
  function calculateFanPosition(player, cardIndex) {
    // Implementation would calculate proper fan position based on player seating and card index
    return { x: 0, y: 0, rotation: 0 }
  }
  
  function calculateDiscardPosition(stackIndex) {
    // Calculate stacked position in discard pile
    return { 
      x: stackIndex * 2, // 2px horizontal offset
      y: stackIndex * -1, // 1px vertical offset
      zIndex: stackIndex + 1
    }
  }
  
  async function animateElementAlongPath(element, path, options = {}) {
    return new Promise(resolve => {
      const { duration = 1000, easing = 'ease', styleTransitions = {} } = options
      
      // Apply starting styles
      if (styleTransitions.from) {
        element.classList.remove(styleTransitions.to)
        element.classList.add(styleTransitions.from)
      }
      
      // Set up transition
      element.style.transition = `transform ${duration}ms ${easing}`
      
      // Start animation
      requestAnimationFrame(() => {
        const deltaX = path.to.x - path.from.x
        const deltaY = path.to.y - path.from.y
        
        element.style.transform = `translate(${deltaX}px, ${deltaY}px)`
        
        // Apply ending styles
        if (styleTransitions.to) {
          element.classList.remove(styleTransitions.from)
          element.classList.add(styleTransitions.to)
        }
      })
      
      // Complete animation
      setTimeout(() => {
        element.style.transition = ''
        resolve()
      }, duration)
    })
  }
  
  async function animateElementToPosition(element, position, options = {}) {
    return new Promise(resolve => {
      const { duration = 1000, easing = 'ease' } = options
      
      element.style.transition = `transform ${duration}ms ${easing}`
      
      requestAnimationFrame(() => {
        let transform = `translate(${position.x}px, ${position.y}px)`
        if (position.rotation !== undefined) {
          transform += ` rotate(${position.rotation}deg)`
        }
        element.style.transform = transform
        
        if (position.zIndex !== undefined) {
          element.style.zIndex = position.zIndex
        }
      })
      
      setTimeout(() => {
        element.style.transition = ''
        resolve()
      }, duration)
    })
  }
  
  // Utility function for delays
  function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms))
  }
  
  // Session logging
  function logSessionEvent(type, message, data = null) {
    // This would emit to parent component for actual logging
    
  }
  
  /**
   * COMMUNICATION CONTROL
   */
  
  function isCommunicationBlocked() {
    return animationState.communicationBlocked
  }
  
  function blockCommunication() {
    animationState.communicationBlocked = true
    
  }
  
  function unblockCommunication() {
    animationState.communicationBlocked = false
    
  }
  
  /**
   * STATE MANAGEMENT
   */
  
  function resetAnimationSystem() {
    animationState.phase = 'idle'
    animationState.communicationBlocked = false
    animationState.pauseActive = false
    animationState.pauseStartTime = null
    animationState.pauseRemainingTime = 0
    animationState.dealingQueue = []
    animationState.playingQueue = []
    animationState.collectingQueue = []
    animationState.discardingQueue = []
    animationState.activeAnimations.clear()
    animationState.cutIndicationVisible = false
    animationState.cutIndicationData = null
    animationState.highlightedCardId = null
    
    
  }
  
  /**
   * EXPOSED API
   */
  
  return {
    // State
    animationState: animationState,
    
    // Core animations
    animateCardDealing,
    animateCardPlay,
    animateCardCollection,
    animateDiscard,
    
    // Pause system
    initiatePause,
    completePause,
    extendPause,
    
    // Visual effects
    showCutIndication,
    highlightHighestCard,
    
    // Communication control
    isCommunicationBlocked,
    blockCommunication,
    unblockCommunication,
    
    // Utilities
    resetAnimationSystem,
    sleep
  }
}
