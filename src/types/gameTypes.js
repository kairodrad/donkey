// Enhanced game state types to match backend models and support all prompt requirements
// Based on internal/model/models.go and prompt specifications

/**
 * User/Player Types
 */

// Player representation with full backend state support
export const createPlayer = (data = {}) => ({
  id: data.id || '',
  name: data.name || '',
  isBot: data.isBot || false,
  botDifficulty: data.botDifficulty || null, // 'easy', 'medium', 'difficult'
  isConnected: data.isConnected !== undefined ? data.isConnected : true,
  donkeyLetters: data.donkeyLetters || '', // '', 'D', 'DO', 'DON', 'DONK', 'DONKE', 'DONKEY'
  joinOrder: data.joinOrder || 0,
  joinedAt: data.joinedAt || new Date().toISOString(),
  lastSeenAt: data.lastSeenAt || new Date().toISOString()
})

// Player position in oval layout (from gameUtils.js)
export const createPlayerSeat = (data = {}) => ({
  ...createPlayer(data),
  seat: {
    top: data.seat?.top || '50%',
    left: data.seat?.left || '50%',
    transform: data.seat?.transform || 'translate(-50%, -50%)'
  },
  seatAngle: data.seatAngle || 0 // Angle on oval for In-Game card positioning
})

/**
 * Card Types
 */

// Enhanced card with backend state support
export const createCard = (data = {}) => ({
  id: data.id || '',
  rank: data.rank || '2', // '2'-'10', 'J', 'Q', 'K', 'A'
  suit: data.suit || 'hearts', // 'hearts', 'diamonds', 'clubs', 'spades'
  value: data.value || getCardValue(data.rank || '2'), // Numeric value for comparison
  location: data.location || 'hand', // 'deck', 'hand', 'in_play', 'discard'
  ownerId: data.ownerId || null,
  sortOrder: data.sortOrder || 0,
  
  // Frontend display state
  flipped: data.flipped !== undefined ? data.flipped : false,
  selected: data.selected !== undefined ? data.selected : false,
  disabled: data.disabled !== undefined ? data.disabled : false,
  disabledFeedback: data.disabledFeedback !== undefined ? data.disabledFeedback : false,
  
  // Animation state
  animating: data.animating !== undefined ? data.animating : false,
  animationType: data.animationType || null // 'dealing', 'playing', 'collecting', 'discarding'
})

// Played card in a turn
export const createPlayedCard = (data = {}) => ({
  id: data.id || '',
  turnId: data.turnId || '',
  cardId: data.cardId || '',
  playerId: data.playerId || '',
  playOrder: data.playOrder || 1,
  playedAt: data.playedAt || new Date().toISOString(),
  card: createCard(data.card || {}),
  player: createPlayer(data.player || {})
})

/**
 * Game Structure Types
 */

// Turn within a round
export const createTurn = (data = {}) => ({
  id: data.id || '',
  roundId: data.roundId || '',
  turnNumber: data.turnNumber || 1,
  startPlayerId: data.startPlayerId || '',
  leadSuit: data.leadSuit || null, // 'hearts', 'diamonds', 'clubs', 'spades'
  status: data.status || 'active', // 'active', 'cut', 'completed'
  winnerId: data.winnerId || null,
  cutPlayerId: data.cutPlayerId || null,
  startedAt: data.startedAt || new Date().toISOString(),
  completedAt: data.completedAt || null,
  playedCards: (data.playedCards || []).map(createPlayedCard),
  
  // Frontend computed properties
  expectedPlayerId: data.expectedPlayerId || null, // Computed from playedCards and round players
  isComplete: data.isComplete !== undefined ? data.isComplete : false
})

// Round player status
export const createRoundPlayer = (data = {}) => ({
  roundId: data.roundId || '',
  userId: data.userId || '',
  user: createPlayer(data.user || {}),
  position: data.position || 0,
  isFinished: data.isFinished !== undefined ? data.isFinished : false,
  finishedAt: data.finishedAt || null,
  cardsInHand: data.cardsInHand || 0
})

// Round within a game
export const createRound = (data = {}) => ({
  id: data.id || '',
  gameId: data.gameId || '',
  roundNumber: data.roundNumber || 1,
  status: data.status || 'setup', // 'setup', 'dealing', 'active', 'completed'
  startedAt: data.startedAt || new Date().toISOString(),
  completedAt: data.completedAt || null,
  loserId: data.loserId || null,
  turns: (data.turns || []).map(createTurn),
  roundPlayers: (data.roundPlayers || []).map(createRoundPlayer),
  cards: (data.cards || []).map(createCard),
  
  // Frontend computed properties
  activePlayers: data.activePlayers || [], // Players still in this round
  finishedPlayers: data.finishedPlayers || [] // Players who finished early
})

// Main game state
export const createGame = (data = {}) => ({
  id: data.id || '',
  requesterId: data.requesterId || '',
  status: data.status || 'waiting', // 'waiting', 'active', 'completed', 'abandoned', 'paused'
  maxPlayers: data.maxPlayers || 8,
  minPlayers: data.minPlayers || 2,
  createdAt: data.createdAt || new Date().toISOString(),
  startedAt: data.startedAt || null,
  completedAt: data.completedAt || null,
  loserId: data.loserId || null, // Final DONKEY loser
  
  // Related data
  players: (data.players || []).map(createPlayerSeat), // Players with seating
  rounds: (data.rounds || []).map(createRound),
  sessionLogs: (data.sessionLogs || []).map(createSessionLog)
})

/**
 * Session and Communication Types
 */

// Session log entry
export const createSessionLog = (data = {}) => ({
  id: data.id || '',
  gameId: data.gameId || '',
  userId: data.userId || null,
  type: data.type || 'system', // 'chat', 'system', 'game_event', 'round_event', 'turn_event'
  message: data.message || '',
  eventData: data.eventData || null,
  createdAt: data.createdAt || new Date().toISOString(),
  
  // Frontend display properties
  isSystemMessage: data.type !== 'chat',
  displayTime: formatTimeForDisplay(data.createdAt || new Date().toISOString())
})

/**
 * Animation and Visual State Types
 */

// Animation state for coordinated timing
export const createAnimationState = (data = {}) => ({
  // Current animation phase
  phase: data.phase || 'idle', // 'idle', 'dealing', 'playing', 'pausing', 'cutting', 'discarding'
  
  // Timing control
  isPaused: data.isPaused !== undefined ? data.isPaused : false,
  pauseStartTime: data.pauseStartTime || null,
  pauseDuration: data.pauseDuration || 3000, // 3 seconds default
  
  // Animation queues
  dealingQueue: data.dealingQueue || [],
  playingQueue: data.playingQueue || [],
  collectingQueue: data.collectingQueue || [],
  discardingQueue: data.discardingQueue || [],
  
  // Communication control (for 3-second pause system)
  communicationBlocked: data.communicationBlocked !== undefined ? data.communicationBlocked : false,
  
  // Visual state
  highlightedCardId: data.highlightedCardId || null, // For highest card indication
  cutIndicationVisible: data.cutIndicationVisible !== undefined ? data.cutIndicationVisible : false,
  cutIndicationData: data.cutIndicationData || null
})

// Complete application state combining all subsystems
export const createAppState = (data = {}) => ({
  // Core game data
  game: createGame(data.game || {}),
  currentRound: data.currentRound ? createRound(data.currentRound) : null,
  currentTurn: data.currentTurn ? createTurn(data.currentTurn) : null,
  
  // User-specific data  
  user: createPlayer(data.user || {}),
  myCards: (data.myCards || []).map(createCard),
  
  // Game state
  playersWithSeats: (data.playersWithSeats || []).map(createPlayerSeat),
  inPlayCards: (data.inPlayCards || []).map(createPlayedCard),
  discardPile: (data.currentRound?.discardPile || data.discardPile || []).map(createCard),
  
  // Animation and visual state
  animationState: createAnimationState(data.animationState || {}),
  
  // Communication state
  isConnected: data.isConnected !== undefined ? data.isConnected : false,
  lastEventTime: data.lastEventTime || null,
  
  // UI state
  showMenu: data.showMenu !== undefined ? data.showMenu : false,
  currentTheme: data.currentTheme || 'system',
  selectedCardIds: data.selectedCardIds || [],
  
  // Session data
  sessionLogs: (data.sessionLogs || []).map(createSessionLog),
  chatMessages: (data.chatMessages || []).map(createSessionLog)
})

/**
 * Utility Functions
 */

// Get numeric value for card comparison
export function getCardValue(rank) {
  const values = {
    '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, '10': 10,
    'J': 11, 'Q': 12, 'K': 13, 'A': 14
  }
  return values[rank] || 2
}

// Sort cards according to RULES.md specification
export function sortCards(cards) {
  const suitOrder = { 'diamonds': 0, 'clubs': 1, 'hearts': 2, 'spades': 3 }
  const rankOrder = { '2': 0, '3': 1, '4': 2, '5': 3, '6': 4, '7': 5, '8': 6, '9': 7, '10': 8, 'J': 9, 'Q': 10, 'K': 11, 'A': 12 }
  
  return [...cards].sort((a, b) => {
    // First sort by suit (Diamond, Clubs, Hearts, Spades)
    const suitDiff = suitOrder[a.suit] - suitOrder[b.suit]
    if (suitDiff !== 0) return suitDiff
    
    // Then sort by rank within suit (2-A)
    return rankOrder[a.rank] - rankOrder[b.rank]
  })
}

// Format timestamp for display
export function formatTimeForDisplay(timestamp) {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('en-US', { 
    hour: '2-digit', 
    minute: '2-digit',
    hour12: false
  })
}

// Check if player has finished the game (DONKEY)
export function isDonkey(player) {
  return player.donkeyLetters === 'DONKEY'
}

// Get next letter in DONKEY progression
export function getNextDonkeyLetter(currentLetters) {
  const progression = 'DONKEY'
  return currentLetters.length < progression.length 
    ? progression[currentLetters.length]
    : null
}

// Validate card play according to rules
export function isValidCardPlay(card, gameState) {
  const { currentTurn, myCards } = gameState
  
  // Can always play if no lead suit
  if (!currentTurn?.leadSuit) return true
  
  // Must follow suit if possible
  if (card.suit === currentTurn.leadSuit) return true
  
  // Can cut if no cards of lead suit
  const hasLeadSuit = myCards.some(c => c.suit === currentTurn.leadSuit)
  return !hasLeadSuit
}

// Get legal cards for current turn
export function getLegalCards(cards, gameState) {
  return cards.filter(card => isValidCardPlay(card, gameState))
}

/**
 * State Validation Functions
 */

// Validate that game state is consistent
export function validateGameState(state) {
  const errors = []
  
  if (!state.game?.id) {
    errors.push('Game ID is required')
  }
  
  if (state.game?.status === 'active' && !state.currentRound) {
    errors.push('Active game must have current round')
  }
  
  if (state.currentRound?.status === 'active' && !state.currentTurn) {
    errors.push('Active round must have current turn')
  }
  
  // Validate player consistency
  const playerIds = new Set(state.playersWithSeats.map(p => p.id))
  if (state.user?.id && !playerIds.has(state.user.id)) {
    errors.push('User must be in players list')
  }
  
  return errors
}

/**
 * Type Guards
 */

export function isBot(player) {
  return player.isBot === true
}

export function isHuman(player) {
  return player.isBot === false
}

export function isGameActive(game) {
  return game.status === 'active'
}

export function isRoundActive(round) {
  return round?.status === 'active'
}

export function isTurnActive(turn) {
  return turn?.status === 'active'
}