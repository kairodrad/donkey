# DONKEY CARD GAME - WIREFRAMES

This document provides ASCII wireframes for the new grid-based layout replacing the oval table design.

## MAIN GAME LAYOUT (8 Players)

```
┌─────────────────────────────────────────────────────────────────────┐
│ DONKEY GAME - [Settings] [Chat] [Help]                    [Theme]   │
├─────────────────────────────────────────────────────────────────────┤
│                          VIEWPORT (100vh - header)                 │
│ ┌───────────────────────────────────────────┐ ┌──────────────────┐ │
│ │           OPPONENT GRID (80% width)        │ │  DISCARD PILE    │ │
│ │              (60% height)                  │ │   (20% width)    │ │
│ │ ┌─────────────────────────────────────────┐ │ │ ┌──────────────┐ │ │
│ │ │ P2: █ █ █ █ █ █        [In-Game Card]  │ │ │ │  ████████    │ │ │
│ │ │         Player 2                        │ │ │ │  ████████    │ │ │
│ │ ├─────────────────────────────────────────┤ │ │ │  "Discarded" │ │ │
│ │ │ P3: █ █ █ █ █          [In-Game Card]  │ │ │ └──────────────┘ │ │
│ │ │         Player 3                        │ │ └──────────────────┘ │
│ │ ├─────────────────────────────────────────┤ │                    │
│ │ │ P4: █ █ █ █ █ █ █      [In-Game Card]  │ │                    │
│ │ │         Player 4                        │ │                    │
│ │ ├─────────────────────────────────────────┤ │                    │
│ │ │ P5: █ █ █ █ █          [In-Game Card]  │ │                    │
│ │ │         Player 5                        │ │                    │
│ │ ├─────────────────────────────────────────┤ │                    │
│ │ │ P6: █ █ █ █ █ █        [In-Game Card]  │ │                    │
│ │ │         Player 6                        │ │                    │
│ │ ├─────────────────────────────────────────┤ │                    │
│ │ │ P7: █ █ █ █ █ █ █ █    [In-Game Card]  │ │                    │
│ │ │         Player 7                        │ │                    │
│ │ ├─────────────────────────────────────────┤ │                    │
│ │ │ P8: █ █ █ █ █          [In-Game Card]  │ │                    │
│ │ │         Player 8                        │ │                    │
│ │ └─────────────────────────────────────────┘ │                    │
│ └───────────────────────────────────────────────┘                    │
│                                                                     │
│ ┌─────────────────────────────────────────────────────────────────┐ │
│ │                    CURRENT PLAYER (40% height)                 │ │
│ │                                                  [In-Game Card] │ │
│ │      ► Player 1 (You) - Your Turn                              │ │
│ │                                                                 │ │
│ │              🂠  🂡  🂢  🂣  🂤  🂥  🂦                        │ │
│ │             ╱ ╱  ╱   ╱   ╱   ╲   ╲  ╲                        │ │
│ │            █    █    █    █    █    █   █                       │ │
│ │               FAN LAYOUT (Rotated Cards)                       │ │
│ └─────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

## OPPONENT CELL DETAILED LAYOUT

```
┌─────────────────────────────────────────┐
│ OPPONENT CELL (Height: GridHeight/N)    │
│                                         │
│ ┌───────────────────┐ ┌───────────────┐ │
│ │ LEFT 75%: CARDS   │ │ RIGHT 25%:    │ │
│ │                   │ │  IN-GAME CARD │ │
│ │ █ █ █ █ █ █       │ │               │ │
│ │ (Fan Layout)      │ │    ████████   │ │
│ │                   │ │    ████████   │ │
│ │                   │ │    (80% cell  │ │
│ │                   │ │     height)   │ │
│ └───────────────────┘ └───────────────┘ │
│                                         │
│ ┌─────────────────────────────────────┐ │
│ │ BOTTOM 20%: PLAYER NAME PILL        │ │
│ │        [Player Name]                │ │
│ └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

## TURN INDICATOR STATES

### Active Player (Current Turn)
```
┌─────────────────────────────────────────┐
│ ► P2: █ █ █ █ █ █      [In-Game Card]  │ ← Bouncing Arrow
│     🟦Player 2 (Active)🟦              │ ← Pulsing Blue BG
└─────────────────────────────────────────┘
```

### Inactive Player
```
┌─────────────────────────────────────────┐
│   P3: █ █ █ █ █        [In-Game Card]  │
│      Player 3                           │ ← Gray Background
└─────────────────────────────────────────┘
```

## RESPONSIVE LAYOUTS

### 3 Players Layout
```
┌─────────────────────────────────────────────────────────────────────┐
│ ┌───────────────────────────────────────────┐ ┌──────────────────┐ │
│ │           OPPONENT GRID                    │ │  DISCARD PILE    │ │
│ │ ┌─────────────────────────────────────────┐ │ │                  │ │
│ │ │ P2: █ █ █ █ █ █        [In-Game Card]  │ │ │                  │ │
│ │ │         Player 2                        │ │ │                  │ │
│ │ ├─────────────────────────────────────────┤ │ │                  │ │
│ │ │ P3: █ █ █ █ █          [In-Game Card]  │ │ │                  │ │
│ │ │         Player 3                        │ │ │                  │ │
│ │ ├─────────────────────────────────────────┤ │ │                  │ │
│ │ │            (Empty Space)                │ │ │                  │ │
│ │ │      (GridHeight/3 each)                │ │ │                  │ │
│ │ └─────────────────────────────────────────┘ │ │                  │ │
│ └───────────────────────────────────────────────┘ └──────────────────┘ │
│                                                                     │
│ ┌─────────────────────────────────────────────────────────────────┐ │
│ │                    CURRENT PLAYER                               │ │
│ │                              [In-Game Card]                     │ │
│ │                   Player 1 (You)                               │ │
│ │                                                                 │ │
│ │              🂠  🂡  🂢  🂣  🂤  🂥  🂦                        │ │
│ └─────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

### Mobile Layout (≤768px)
```
┌─────────────────────────────────┐
│ DONKEY [⚙] [💬] [❓] [🌙]    │
├─────────────────────────────────┤
│ ┌─────────────────┐ ┌─────────┐ │
│ │  OPPONENT GRID  │ │DISCARD  │ │
│ │ P2: ██ ██    [█]│ │   ██    │ │
│ │     Player 2    │ │   ██    │ │
│ │ P3: ██ ██    [█]│ │  "Disc" │ │
│ │     Player 3    │ └─────────┘ │
│ │ P4: ██ ██    [█]│             │
│ │     Player 4    │             │
│ └─────────────────┘             │
│                                 │
│ ┌─────────────────────────────┐ │
│ │    CURRENT PLAYER           │ │
│ │          [In-Game]          │ │
│ │     ► Player 1 (You)        │ │
│ │                             │ │
│ │    █  █  █  █  █  █         │ │
│ │   (Smaller Cards)           │ │
│ └─────────────────────────────┘ │
└─────────────────────────────────┘
```

## CARD ANIMATION FLOWS

### Card Play Animation
```
FROM: Current Player Fan Position
  │
  ├─ Scale down to In-Game size
  │
  ├─ Move to top-right In-Game area
  │
  └─ 3-second pause for all players

FROM: Bot Player Cards
  │
  ├─ Card appears in In-Game area
  │
  ├─ Fan updates (card removed)
  │
  └─ 3-second pause for all players
```

### Turn Completion Animation
```
CUT Scenario:
  In-Game Cards → Winner's Fan (animate movement)
  Display "CUT!" banner
  3-second pause
  Clear In-Game area
  Update turn indicator

DISCARD Scenario:
  In-Game Cards → Discard Pile (animate movement)
  Update discard stack visual
  3-second pause  
  Clear In-Game area
  Update turn indicator
```

## TECHNICAL SPECIFICATIONS

### CSS Grid Structure
```css
.game-viewport {
  height: calc(100vh - header-height);
  display: grid;
  grid-template-rows: 60% 40%;
  grid-template-columns: 80% 20%;
}

.opponent-area {
  grid-row: 1;
  grid-column: 1;
  display: flex;
  flex-direction: column;
}

.discard-area {
  grid-row: 1;
  grid-column: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
}

.current-player-area {
  grid-row: 2;
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  position: relative;
}
```

### Opponent Cell Layout
```css
.opponent-cell {
  height: calc(var(--grid-height) / var(--opponent-count));
  min-height: calc(var(--grid-height) / 3);
  display: grid;
  grid-template-columns: 75% 25%;
  grid-template-rows: 80% 20%;
}

.opponent-cards {
  grid-column: 1;
  grid-row: 1;
}

.opponent-ingame {
  grid-column: 2;
  grid-row: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.opponent-name {
  grid-column: 1 / -1;
  grid-row: 2;
  display: flex;
  align-items: center;
  justify-content: center;
}
```

This wireframe specification provides the complete visual blueprint for implementing the new grid-based layout while maintaining intuitive gameplay flow.