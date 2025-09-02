// Utility functions ported from the original React app

export function getCookie(name) {
  const match = document.cookie.match('(?:^|; )' + name + '=([^;]*)')
  return match ? decodeURIComponent(match[1]) : null
}

export function setCookie(name, value) {
  document.cookie = name + '=' + encodeURIComponent(value) + '; path=/'
}

// Card back color options
export const cardBackColors = ['red', 'blue', 'green', 'gray', 'purple', 'yellow']

// Player grid layout for new design - opponents in vertical grid, current player at bottom
export function getPlayerSeats(players, currentUserId) {
  const currentUser = players.find(p => p.id === currentUserId)
  const opponents = players.filter(p => p.id !== currentUserId)
  const opponentCount = opponents.length
  
  // Determine opponent grid dimensions per VISUAL.md specification
  const gridHeight = 60 // 60% of viewport height for opponent grid
  const cellHeight = opponentCount < 3 
    ? gridHeight / 3 // VISUAL.md: When fewer than 3 opponents: (Grid Height / 3)
    : gridHeight / opponentCount // VISUAL.md: When 3 or more opponents: (Grid Height / N)
  
  const result = []
  
  // Position opponents in vertical grid (top 60% of viewport)
  opponents.forEach((opponent, index) => {
    const top = index * cellHeight // Position at exact cell boundaries per VISUAL.md
    
    result.push({
      ...opponent,
      seat: {
        position: 'opponent-grid',
        gridIndex: index,
        cellHeight: `${cellHeight}vh`,
        top: `${top}vh`,
        left: '0',
        width: '100vw', // 100% width for opponent grid (discard pile removed)
        height: `${cellHeight}vh`
      },
      isOpponent: true
    })
  })
  
  // Position current user at bottom (40% of viewport)
  if (currentUser) {
    result.push({
      ...currentUser,
      seat: {
        position: 'current-player',
        top: '60vh', // Start at 60vh (after opponent grid)
        left: '0',
        width: '100vw', // Full width for current player area
        height: '40vh' // 40% height for current player
      },
      isCurrentUser: true
    })
  }
  
  // Maintain logical circular order for turn sequence
  // Opponents are indexed 0 to N-1 from top to bottom
  // Current user comes after last opponent in turn order
  return result.map((player) => ({
    ...player,
    turnOrder: player.isCurrentUser ? opponents.length : player.seat.gridIndex
  }))
}

// Card sorting function
const suitOrder = { D: 0, C: 1, H: 2, S: 3 }
const rankOrder = { '2': 0, '3': 1, '4': 2, '5': 3, '6': 4, '7': 5, '8': 6, '9': 7, '10': 8, 'J': 9, 'Q': 10, 'K': 11, 'A': 12 }

export function sortCards(cards) {
  return [...cards].sort((a, b) => {
    const suitA = a.slice(-1)
    const suitB = b.slice(-1)
    const rankA = a.slice(0, -1)
    const rankB = b.slice(0, -1)
    
    if (suitOrder[suitA] !== suitOrder[suitB]) {
      return suitOrder[suitA] - suitOrder[suitB]
    }
    return rankOrder[rankA] - rankOrder[rankB]
  })
}

// Convert card string to object format
export function parseCard(cardString) {
  const suit = cardString.slice(-1)
  const rank = cardString.slice(0, -1)
  
  const suitMap = {
    'H': 'hearts',
    'D': 'diamonds',
    'C': 'clubs',
    'S': 'spades'
  }
  
  return {
    rank,
    suit: suitMap[suit] || 'hearts',
    flipped: false,
    selected: false,
    disabled: false
  }
}

// Convert array of card strings to card objects
export function parseCards(cardStrings) {
  return cardStrings.map(parseCard)
}

// API helper functions
export async function apiCall(endpoint, options = {}) {
  const response = await fetch(endpoint, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers
    },
    ...options
  })
  
  if (!response.ok) {
    throw new Error(`API call failed: ${response.status}`)
  }
  
  // Check if response has content before parsing JSON
  const text = await response.text()
  if (!text.trim()) {
    // Return empty object for empty responses (common for POST operations)
    return {}
  }
  
  try {
    return JSON.parse(text)
  } catch (error) {
    throw new Error(`Invalid JSON response: ${error.message}`)
  }
}

export async function registerUser(name) {
  return apiCall('/api/register', {
    method: 'POST',
    body: JSON.stringify({ name })
  })
}

export async function createGame(requesterId) {
  return apiCall('/api/game/create', {
    method: 'POST',
    body: JSON.stringify({ requesterId })
  })
}

export async function startGame(gameId, userId) {
  return apiCall('/api/game/start', {
    method: 'POST',
    body: JSON.stringify({ gameId, userId })
  })
}

export async function joinGame(gameId, userId) {
  return apiCall('/api/game/join', {
    method: 'POST',
    body: JSON.stringify({ gameId, userId })
  })
}

export async function addBot(gameId, userId, difficulty = 'easy') {
  return apiCall('/api/game/add-bot', {
    method: 'POST',
    body: JSON.stringify({ gameId, userId, difficulty })
  })
}

export async function playCard(gameId, userId, cardId) {
  return apiCall('/api/game/play-card', {
    method: 'POST',
    body: JSON.stringify({ gameId, userId, cardId })
  })
}

export async function getGameList(userId, status = null) {
  const params = new URLSearchParams({ userId })
  if (status) params.append('status', status)
  return apiCall(`/api/games?${params}`)
}

// Legacy finalize function (deprecated - use startGame instead)
export async function finalizeGame(gameId, userId) {
  return startGame(gameId, userId)
}

export async function abandonGame(gameId, userId) {
  return apiCall('/api/game/abandon', {
    method: 'POST',
    body: JSON.stringify({ gameId, userId })
  })
}

export async function sendChatMessage(gameId, userId, message) {
  return apiCall('/api/game/chat', {
    method: 'POST',
    body: JSON.stringify({ gameId, userId, message })
  })
}

export async function getGameState(gameId, userId) {
  return apiCall(`/api/game/${gameId}/state/${userId}`)
}

export async function getGameLogs(gameId) {
  return apiCall(`/api/game/${gameId}/logs`)
}

export async function getUser(userId) {
  return apiCall(`/api/user/${userId}`)
}

export async function getVersion() {
  return apiCall('/api/version')
}
