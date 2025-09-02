# Donkey Card Game - Complete Game Design Specification

Create a comprehensive multiplayer card game application called "Donkey" that brings families together through an engaging digital card game experience with modern web technologies and mobile-first design.

## 1. Application Overview

**Purpose Statement**: Donkey is a multiplayer card game application that allows 2-8 players to play a fast-paced, family-friendly card game together in real-time through web browsers on any device.

**Core Value Proposition**: Transforms the traditional "Donkey" card game into a digital experience that maintains the social interaction and excitement of the physical game while adding modern conveniences like real-time synchronization, chat, and cross-device compatibility.

**Success Metrics**: 
- Players can complete a full game session (registration to final scoring) in under 15 minutes
- 95% uptime for real-time game state synchronization
- Intuitive mobile interface requiring no tutorial for basic gameplay
- Support for simultaneous games with up to 8 players each

## 2. Technical Foundation

**Technology Stack**:
- **Frontend**: Vue.js 3.x with Composition API, Vite 4.x for build tooling, Tailwind CSS 3.x for styling
- **Backend**: Go 1.24.3+ with Gin web framework, GORM for database ORM
- **Database**: SQLite for development/testing, PostgreSQL for production
- **Real-time Communication**: Server-Sent Events (SSE) for game state updates
- **Development Tools**: Node.js v16.14.1 (NVM managed), Go modules for dependency management
- **API Documentation**: Swagger/OpenAPI 3.0

**Architecture Pattern**: 
- Single Page Application (SPA) with RESTful API backend
- Real-time event-driven architecture using SSE streams
- Component-based frontend with Vue 3 composition patterns
- Modular Go backend with clean separation of concerns

**Data Structure**:
```
User: { id, username, created_at }
Game: { id, status, requester_id, created_at, finalized_at }
GamePlayer: { user_id, game_id, position, joined_at }
GameCard: { id, game_id, user_id, card_value, card_suit }
GameState: { game_id, current_phase, data }
GameSessionLog: { game_id, user_id, message_type, content, timestamp }
```

**State Management**: 
- Frontend: Vue 3 Composition API with reactive state management
- Backend: GORM-managed database state with in-memory caching for active games
- Real-time synchronization via SSE streams per game/user combination

## 3. Feature Decomposition

### 3.1 User Registration
**Purpose**: Allow players to join the game system with a unique identifier

**User Journey**:
1. Player visits the application URL
2. Sees registration form with username input field
3. Enters desired username (3-20 characters, alphanumeric)
4. Clicks "Join Game" button
5. System generates unique hex ID and stores user
6. Player is redirected to game lobby

**Visual Design**:
- Clean, centered form on pastel background (#f8fafc)
- Username input with rounded borders (8px), focus states with primary color ring
- "Join Game" button with primary gradient (bg-gradient-to-r from-blue-500 to-purple-600)
- Donkey game logo/title at top with playful typography
- Mobile-first responsive design, minimum 320px width support

**Functionality**:
- Real-time username validation (no duplicates, length checks)
- Auto-focus on username field
- Enter key submission
- Loading state during registration
- Error handling for network failures and validation errors

**Data Flow**: Form submission → API validation → Database insert → Session creation → Redirect to lobby

**Error States**: 
- Username taken: "Username already exists, please choose another"
- Invalid characters: "Username can only contain letters and numbers"
- Too short/long: "Username must be 3-20 characters"
- Network error: "Connection failed, please try again"

**Success States**: Smooth transition to game lobby with welcome animation

### 3.2 Game Creation and Joining
**Purpose**: Enable players to create new games or join existing ones

**User Journey**:
1. From lobby, player can either:
   - Create new game (becomes requester)
   - Enter game ID to join existing game
2. Game creator sees game ID prominently displayed for sharing
3. Other players enter this game ID to join
4. All players see real-time updates as others join
5. Game supports 2-8 players maximum

**Visual Design**:
- Split interface: "Create Game" and "Join Game" sections
- Game ID displayed in large, copyable text with one-click copy button
- Player list showing joined participants with avatars/initials
- Real-time join animations as players enter
- Clear "Start Game" button (only visible to requester when 2+ players)

**Functionality**:
- Unique 6-character alphanumeric game ID generation
- Real-time player list updates via SSE
- Copy-to-clipboard functionality for game ID sharing
- Auto-start option at 8 players or manual start by requester
- Leave game functionality with proper cleanup

### 3.3 Card Dealing and Game Setup
**Purpose**: Distribute cards and initialize game state for all players

**User Journey**:
1. Requester clicks "Start Game" or game auto-starts at 8 players
2. All players see "Dealing cards..." loading animation
3. System shuffles standard 52-card deck
4. Cards are dealt evenly starting after the requester (clockwise)
5. Remaining cards (if any) are set aside
6. Each player sees their hand in fan arrangement
7. Game phase transitions to "playing"

**Visual Design**:
- Card dealing animation showing cards flying to each player position
- Fan-style hand layout optimized for mobile touch interaction
- Card backs during dealing, revealing player's cards only to them
- Circular seating arrangement visualization showing all player positions
- Smooth transition animations between game phases

**Functionality**:
- Standard 52-card deck (13 ranks × 4 suits)
- Even distribution algorithm with remainder handling
- Secure card visibility (players only see their own cards)
- Touch-optimized card selection with proper spacing for fingers
- Real-time synchronization of dealing process

### 3.4 Real-Time Game State Synchronization
**Purpose**: Keep all players synchronized with current game state

**User Journey**:
1. Any game state change triggers updates to all connected players
2. Players see immediate visual feedback for their actions
3. Other players see updates within 100ms via SSE
4. Connection issues are handled gracefully with reconnection
5. Game state persists if players temporarily disconnect

**Technical Implementation**:
- SSE endpoint: `/api/game/:gameId/stream/:userId`
- Event types: player_joined, card_dealt, game_started, game_ended, chat_message
- Automatic reconnection with exponential backoff
- Server-side connection timeout handling (5 minutes idle)
- State reconciliation on reconnection

### 3.5 In-Game Chat System
**Purpose**: Enable social interaction during gameplay

**User Journey**:
1. Chat panel accessible via expandable interface
2. Players type messages and press Enter or click Send
3. Messages appear in real-time for all players
4. Message history persists for the game session
5. System messages for game events (player joined, cards dealt, etc.)

**Visual Design**:
- Collapsible chat panel that doesn't interfere with game area
- Message bubbles with player name and timestamp
- Different styling for system messages vs player messages
- Auto-scroll to latest messages
- Character limit (280 characters) with counter

## 4. Interaction Specifications

**Input Methods**:
- Touch-optimized for mobile devices (minimum 44px touch targets)
- Keyboard shortcuts: Enter for form submission, Escape for modal closing
- Mouse and touch support for card selection and dragging
- Copy-to-clipboard button functionality

**Validation Rules**:
- Real-time form validation with visual feedback
- Username: 3-20 characters, alphanumeric only
- Game ID: 6 characters, case-insensitive
- Chat messages: 1-280 characters, no empty messages

**Animation and Transitions**:
- Page transitions: 300ms ease-in-out
- Card dealing animation: 150ms per card with staggered timing
- Hover states: 200ms color/shadow transitions
- Loading states: Smooth spinner animations
- Success feedback: 500ms confirmation animations

**Responsive Behavior**:
- Mobile-first design starting at 320px width
- Tablet optimization at 768px+ with enhanced card layouts
- Desktop optimization at 1024px+ with sidebar layouts
- Card fan arrangements that adapt to screen orientation

**Accessibility**:
- ARIA labels for all interactive elements
- Keyboard navigation support
- Screen reader announcements for game state changes
- High contrast mode support
- Focus management for modal dialogs

## 5. Visual Design System

**Color Palette**:
- Primary: #3b82f6 (Blue 500)
- Secondary: #8b5cf6 (Purple 500)
- Accent: #10b981 (Emerald 500)
- Background: #f8fafc (Slate 50)
- Surface: #ffffff (White)
- Error: #ef4444 (Red 500)
- Warning: #f59e0b (Amber 500)
- Success: #10b981 (Emerald 500)
- Text Primary: #1e293b (Slate 800)
- Text Secondary: #64748b (Slate 500)

**Dark Mode Palette**:
- Background: #0f172a (Slate 900)
- Surface: #1e293b (Slate 800)
- Text Primary: #f1f5f9 (Slate 100)
- Text Secondary: #94a3b8 (Slate 400)

**Typography**:
- Font Family: 'Inter', system-ui, sans-serif
- Headings: font-weight: 600, line-height: 1.2
- Body: font-weight: 400, line-height: 1.6
- Size Scale: 12px, 14px, 16px, 18px, 24px, 32px, 48px

**Spacing System**:
- Base unit: 4px (0.25rem)
- Scale: 4px, 8px, 12px, 16px, 24px, 32px, 48px, 64px
- Component padding: 16px default, 24px for cards
- Section margins: 32px vertical, 16px horizontal on mobile

**Component Library**:
- Cards: 8px border radius, subtle shadow, hover elevation
- Buttons: 6px border radius, 300ms transitions, focus rings
- Inputs: 6px border radius, border on focus, error states
- Modals: 12px border radius, backdrop blur effect
- Notifications: Toast-style with auto-dismiss

## 6. Performance and Optimization

**Loading Strategy**:
- Critical CSS inlined, non-critical CSS lazy loaded
- Vue components lazy loaded for different game phases
- Image assets optimized and served via CDN
- Service worker for offline capability

**Performance Targets**:
- Initial page load: <2 seconds on 3G connection
- Real-time updates: <100ms latency for SSE events
- Smooth animations: 60fps on mobile devices
- Memory usage: <50MB for typical game session

**SEO Requirements**:
- Meta tags for social sharing
- Open Graph protocol implementation
- Schema.org structured data for gaming application
- Semantic HTML structure

**Analytics Integration**:
- Game completion rates
- Player engagement metrics
- Performance monitoring
- Error tracking and reporting

## Recent UI Updates

**Removed Components (Backend functionality maintained)**:
- **Discard Pile UI**: Removed from visual interface to simplify gameplay focus. Backend still tracks discarded cards for game logic.

**Enhanced Features**:
- **Highest Card Highlighting**: In-Game cards with highest value show yellow pulsing border
- **Improved Chat System**: "Session Updates" renamed to "Chat Logs" with full-width bottom panel when expanded
- **Optimized Card Animation**: Cards animate directly from fan position to In-Game area without intermediate center positioning

## 7. Security and Validation

**Input Sanitization**:
- All user inputs escaped and validated server-side
- SQL injection prevention via GORM parameterized queries
- XSS prevention via Vue's automatic escaping
- CSRF protection for state-changing operations

**Rate Limiting**:
- API endpoints: 100 requests per minute per IP
- Chat messages: 10 messages per minute per user
- Game creation: 5 games per hour per user

**Data Privacy**:
- No personal information stored beyond username
- Game data automatically deleted after 24 hours of inactivity
- No tracking cookies, minimal session data

## 8. Testing Scenarios

**Happy Path**:
1. User registers with valid username
2. Creates game and shares ID with friends
3. Friends join using game ID
4. Game starts automatically or manually
5. Cards are dealt successfully
6. Players interact and complete game
7. Final scores displayed

**Edge Cases**:
- Maximum 8 players joining simultaneously
- Game creator leaves before starting
- Network interruption during card dealing
- Multiple rapid chat messages
- Browser refresh during active game

**Error Scenarios**:
- Database connection failure
- SSE stream disconnection
- Invalid game ID entry
- Duplicate username registration
- API timeout during critical operations

**Stress Tests**:
- 100 simultaneous games
- Rapid SSE event generation
- Memory leaks during long sessions
- Database performance under load
- Mobile device battery impact

## 9. Game Rules Implementation (To Be Detailed)

**Note**: The specific rules of the "Donkey" card game need to be defined and added to this specification. This should include:
- Win/lose conditions
- Card play mechanics
- Turn order and timing
- Scoring system
- Special card effects or rules
- Game ending conditions

**Current Implementation**: The existing codebase has basic card dealing but lacks the specific game logic. This section should be expanded with the complete rules once they are defined.

## 10. Development Phases

**Phase 1**: Core Infrastructure
- User registration and game creation
- Basic real-time synchronization
- Card dealing mechanics

**Phase 2**: Game Logic Implementation
- Complete game rules implementation
- Turn management and validation
- Win/lose condition handling

**Phase 3**: Enhanced Features
- Chat system with rich features
- Game statistics and history
- Advanced UI polish and animations

**Phase 4**: Production Optimization
- Performance optimization
- Security hardening
- Deployment automation
- Monitoring and analytics

This specification provides a comprehensive foundation for developing the Donkey card game. Each section should be refined and expanded as development progresses and game rules are finalized.