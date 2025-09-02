# DONKEY CARD GAME - VISUAL SPECIFICATION

This document defines the complete visual specification for the Donkey card game interface.

## LAYOUT OVERVIEW

### Main Container
- **Background**: Radial gradient from center (background-secondary) to edges (background)
- **Viewport**: Full screen responsive layout
- **Theme Support**: Light/dark mode with CSS variables

### Player Positioning
- **Viewport**: The area of the screen below the title and menu bar. Most of the remaining dimensions will be in terms of the viewport.
- **Layout**: Bottom 40% of the viewport dedicated to the current player's cards. Top 60% of the viewport dedicated to the opponents' name, cards, opponents' In-Game cards, Discard pile.
- **Current Player**: Fixed at bottom position
- **Opponents**: Distributed equally in the top 60% vertically-separated. 
- **Max Players**: 8 players maximum

### Turn Indicators
- **Active Player**: Pulsing blue background, white text, scale-110
- **Arrow Indicator**: Horizontally bouncing right arrow (▶) to the left of the active player's name
- **Inactive Players**: Gray background with border

### Opponent Grid Layout
Total positioning: Top-left
Total horizontal size: 100% of the width of the viewport (full width)
Total vertical size: 60% of the height of the viewport
Dedicated height for each opponent:
1. When fewer than 3 opponents: (Grid Height / 3)
2. When 3 or more opponents (upto N <= 8): (Grid Height / N)

Layout within each opponent cell:
Bottom 20% of the cell: Name in a pill, aligned below the fan
Left 75% area breakdown:
  - 16px left margin (to prevent clipping from left edge)
  - 65% fan area (cards display area) with left-aligned positioning  
  - Fan spacing calculated to keep all cards within container bounds with 16px margins
  - Cards positioned sequentially from left to right starting from left margin (no negative offsets)
Right area for In-Game card:
  - Card positioned with right edge at 5% margin from cell right edge
  - Card extends leftward as needed to show full card without squishing
  - Height auto-adjusted to maintain proper card aspect ratio and readability
Vertically align all Hands and In-Game card areas for all opponents

### Opponent Arrangement
The opponents and current player are still logically in a circle, even though physically in a grid. The sequence is such that the cyclic order will go from the top of the grid down one-by-down ending with the current player. So, in terms of turns, after the current player, the turn goes to the top player, and the next player, and so on.

## CARD SPECIFICATIONS

### Card Sizes
- **Current Player Cards**: `w-24 h-32` (mobile) → `w-28 h-40` (sm) → `w-32 h-44` (md+)
- **Opponent Cards**: `w-8 h-12` (small for visibility count)
- **In-Game Cards for Opponents**: 80% of the Opponent Grid Cell height, adjust width accordingly.
- **In-Game Cards for Current Player**: use same size established for the opponents above, and positioned top right, away from the fan of cards
- **Discard Pile Cards**: Same as opponent cards

### Card States
- **Normal**: Full color, standard border
- **Selected**: Blue border, elevated (-translate-y-6), pulsing glow animation in-place (only for valid/enabled cards)
- **Disabled**: Same appearance as normal but with red border feedback for 3 seconds when clicked, cursor: not-allowed, no elevation or selection movement
- **Highest In-Game Card**: Yellow border (4px solid #fbbf24) with pulsing glow animation to indicate the card with highest value in current round
- **Hover (not on mobile)**: Elevated (-translate-y-4) with smooth transition (to allow for easier selection)

### Card Arrangement
- **Current Player Area Layout** (Bottom 40% of viewport):
  - **Card fan area**: Top 60% of the 40vh area (positioned at top-0)
  - **Fan sizing**: 70vw width with 5vw left margin (spans from 5vw to 75vw)
  - **Cards arrangement**: Fan out horizontally with center-based positioning offset 20px to the right and rotation (max 45°)
  - **Spacing**: Adaptive up to 30px between cards, constrained to 70vw available width
  - **In-Game card**: Top-right corner of current player area at right-8, top-4
  - **Name area**: Bottom 40% of the 40vh area, positioned below card fan
  - **Name styling**: Same as opponent names with turn indicators
  - **Z-index**: Natural layering (left-bottom, right-top)
  - **Chat Logs**: Start minimized, positioned at bottom-right corner
  - **Minimized view**: Single line showing latest message with ellipsis truncation and + button
  - **No header when minimized**: Only message text and expand button visible
  - **Expanded view**: Full-width bottom panel (192px height) with "Chat Logs" title and − button to minimize
- **Opponent Fans**:
  - Small card backs showing count
  - Minimal spacing for compact display
  - Non-interactive

## USER INTERFACE ELEMENTS

### Current Player Name/Status
- **Current Player**: turn indicator + Player name + "(You)"

### Connection Status
- **Removed**: No longer displayed as separate pill widget
- **Integration**: Connection status reflected in player names within game area

## ANIMATIONS & TIMING

### 3-Second Pause System
- **Trigger**: After any player plays a card
- **Behavior**: Complete communication blackout
- **Visual State**: Cards remain visible in In-Game area
- **User Feedback**: All interactions disabled during pause

### Card Transitions
- **Play Animation**: Direct movement from card's fan position to In-Game position (no center positioning)
  - **Current Player**: Cards animate directly from their fan position to top-right corner of current player area (right-8, top-4)
  - **Opponents**: Cards animate to their respective In-Game areas within grid cells
  - **Selection Behavior**: Selection is cleared immediately when playing to prevent center positioning
  - **Resizing**: Smooth transition between fan size and In-Game size
  - **Duration**: 1.5 seconds with cubic-bezier easing
- **Selection**: Instant elevation with smooth transition
- **Hover Effects**: 300ms cubic-bezier(0.4, 0, 0.2, 1)
- **Fan Arrangement**: Smooth positioning updates

### Feedback Animations
- **Disabled Cards**: Red border feedback shown for 3 seconds when clicked, no card movement
- **Selection Glow**: Pulsing blue border (2s infinite) for valid selected cards only
- **Highest Value Highlight**: Yellow border with pulsing glow animation (2s infinite) for the In-Game card with highest value
- **Turn Indicator**: Bouncing arrow animation
- **CUT Indicator**: Red pulsing banner with sword emojis

## RESPONSIVE DESIGN

### Mobile Optimizations (≤768px)
- **Current Player Cards**: Smaller sizing with adequate touch targets (min 44×60px)
- **Card Spacing**: Reduced spacing (0.25rem padding)
- **Fan Position**: Bottom adjustment (13rem → 12rem)
- **Touch Feedback**: Scale transform on touch (0.95)

### Tablet/Desktop
- **Enhanced Spacing**: More generous card spacing
- **Hover Effects**: Full hover interactions enabled
- **Larger Cards**: Better visibility for current player

## COLOR SCHEME & THEMING

### Theme Variables
- `--color-primary`: Blue selection indicators
- `--color-surface`: Card backgrounds and UI elements
- `--color-background`: Main background
- `--color-background-secondary`: Gradient center
- `--color-text`: Primary text color
- `--color-border`: Standard borders

### Card Backs
- **Available Colors**: red, blue, green, gray, purple, yellow
- **User Preference**: Stored in cookies
- **Consistency**: Same back color for all cards in game

## ACCESSIBILITY FEATURES

### Keyboard Navigation
- **Tab Order**: Logical progression through interactive elements
- **Focus Indicators**: Clear visual focus states
- **Shortcuts**: Space/Enter for card selection

### Screen Reader Support
- **ARIA Labels**: Descriptive labels for cards and game state
- **Role Attributes**: Proper button/interactive element roles
- **State Announcements**: Selection and game state changes

### Reduced Motion Support
- **Preference Detection**: `@media (prefers-reduced-motion: reduce)`
- **Fallback**: Disable animations when requested
- **Essential Motion**: Keep critical game state transitions

## GAME STATE VISUALIZATION

### Waiting States
- **Pre-Game**: Central "Ready to play?" with start button
- **Host Waiting**: Player count display with bot controls
- **Guest Waiting**: "Waiting for host" message

### Active Game
- **Turn Visualization**: Clear current player indication
- **Card Count**: Visible card counts for all players
- **Game Progress**: Discard pile shows game progression

### End States
- **Abandoned Game**: Modal overlay with return options
- **Game Complete**: Victory/completion states

## VISUAL HIERARCHY

### Primary Elements
1. Opponent grid (top focus)
2. Current player's cards (bottom prominence)
3. Turn indicator (clear active state)

### Secondary Elements
1. Opponent cards (minimal distraction)
2. In-Game cards (clear but not dominant)
3. Discard pile (reference information)

### Tertiary Elements
1. Settings and menu options
2. Chat and log information
3. Background and decorative elements

## PERFORMANCE CONSIDERATIONS

### Optimization
- **CSS Transforms**: Hardware acceleration for animations
- **Image Loading**: Efficient card image assets
- **Render Cycles**: Minimal re-renders during pause states
- **Memory Management**: Proper cleanup of animation timers

### Loading States
- **Progressive Enhancement**: Base functionality without JavaScript
- **Skeleton Loading**: Placeholder states for async content
- **Error Boundaries**: Graceful degradation for missing assets

---

## IMPLEMENTATION NOTES

This specification should be used as the authoritative reference for all visual implementations. Any changes to the visual design should be documented here first, then applied to the codebase.

Key files implementing these specifications:
- `src/App.vue` - Main layout and game state, highest card highlighting logic
- `src/components/CardFan.vue` - Card arrangement and interactions  
- `src/components/GameCard.vue` - Individual card rendering and states, highest card highlighting
- `src/utils/gameUtils.js` - Player positioning calculations

**Removed Components:**
- `src/components/DiscardPile.vue` - Discard pile removed from UI (backend functionality maintained)