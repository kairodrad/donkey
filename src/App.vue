<template>
  <div class="h-full flex flex-col bg-[var(--color-background)] text-[var(--color-text)] overflow-hidden">
    <!-- Header -->
    <header class="fixed top-0 left-0 w-full h-16 flex items-center justify-between px-4 bg-[var(--color-background)] border-b border-[var(--color-border)] z-30">
      <!-- Menu Button -->
      <div class="relative">
        <BaseButton
          variant="ghost"
          size="md"
          icon-only
          @click="toggleMenu"
          class="text-lg"
        >
          <template #icon>☰</template>
        </BaseButton>
        
        <!-- Dropdown Menu -->
        <div
          v-if="showMenu"
          class="absolute top-full left-0 mt-2 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg shadow-game z-40 min-w-max"
        >
          <button
            v-for="item in menuItems"
            :key="item.label"
            :disabled="item.disabled"
            @click="handleMenuAction(item.action)"
            class="block w-full text-left px-4 py-2 text-sm hover:bg-[var(--color-background-secondary)] disabled:opacity-50 disabled:cursor-not-allowed first:rounded-t-lg last:rounded-b-lg"
          >
            {{ item.label }}
          </button>
        </div>
      </div>
      
      <!-- Game Title -->
      <h1 class="text-2xl font-extrabold tracking-wide bg-clip-text text-transparent"
          style="background-image: linear-gradient(90deg, var(--color-primary), var(--color-text));">
        🫏 Donkey 🫏
      </h1>
      
      <!-- Theme Toggle -->
      <BaseButton
        variant="ghost"
        size="md"
        icon-only
        @click="cycleTheme"
        :title="`Current theme: ${themeLabel}`"
      >
        <template #icon>
          <span v-if="currentTheme === 'light'">☀️</span>
          <span v-else-if="currentTheme === 'dark'">🌙</span>
          <span v-else>🌓</span>
        </template>
      </BaseButton>
    </header>
    
    <!-- Turn Result Notification Popups -->
    <div v-if="turnResultNotification" class="fixed bottom-20 left-4 right-4 z-50 flex justify-center">
      <div class="bg-[var(--color-surface)] border-2 border-[var(--color-border)] rounded-lg shadow-xl px-6 py-4 max-w-md w-full">
        <div class="flex items-center justify-between">
          <div class="flex-1">
            <div class="text-sm font-medium text-[var(--color-text)]" v-html="turnResultNotification.message"></div>
          </div>
          <BaseButton
            variant="ghost"
            size="xs"
            icon-only
            @click="dismissTurnResultNotification"
            class="ml-3 flex-shrink-0"
          >
            <template #icon>×</template>
          </BaseButton>
        </div>
        <!-- Progress bar -->
        <div class="mt-3 w-full bg-[var(--color-background)] rounded-full h-1">
          <div 
            class="bg-[var(--color-primary)] h-1 rounded-full transition-all duration-100 ease-linear"
            :style="{ width: `${turnResultProgress}%` }"
          ></div>
        </div>
      </div>
    </div>
    
    <!-- Game End Popup -->
    <div v-if="gameEndPopup" class="fixed inset-0 z-60 flex items-center justify-center bg-black bg-opacity-50">
      <div class="bg-[var(--color-surface)] border-2 border-[var(--color-border)] rounded-lg shadow-2xl p-6 max-w-lg w-full mx-4">
        <div class="text-center">
          <h2 class="text-2xl font-bold text-[var(--color-text)] mb-2">🎉 Game Over!</h2>
          <p class="text-lg text-[var(--color-text-secondary)] mb-6">
            <strong>{{ gameEndPopup.loserName }}</strong> is the DONKEY!
          </p>
          
          <!-- Scoreboard -->
          <div class="bg-[var(--color-background-secondary)] rounded-lg p-4 mb-4">
            <h3 class="text-sm font-semibold text-[var(--color-text)] mb-3">Final Scores</h3>
            <div class="space-y-2">
              <div 
                v-for="player in gameEndPopup.scoreboard" 
                :key="player.playerId"
                class="flex items-center justify-between px-3 py-2 rounded"
                :class="{
                  'bg-red-100 border border-red-300': player.isLoser,
                  'bg-green-100 border border-green-300': !player.isLoser && player.donkeyLetters === 0,
                  'bg-[var(--color-surface)]': !player.isLoser && player.donkeyLetters > 0
                }"
              >
                <span class="text-sm font-medium" :class="{
                  'text-red-800': player.isLoser,
                  'text-green-800': !player.isLoser && player.donkeyLetters === 0,
                  'text-[var(--color-text)]': !player.isLoser && player.donkeyLetters > 0
                }">
                  {{ player.playerName }}
                  <span v-if="player.isBot" class="text-xs opacity-75">(Bot)</span>
                  <span v-if="player.isLoser" class="ml-1">🫏</span>
                  <span v-else-if="player.donkeyLetters === 0" class="ml-1">👑</span>
                </span>
                <span class="text-sm font-mono" :class="{
                  'text-red-800': player.isLoser,
                  'text-green-800': !player.isLoser && player.donkeyLetters === 0,
                  'text-[var(--color-text-secondary)]': !player.isLoser && player.donkeyLetters > 0
                }">
                  {{ player.donkeyLetters }}/6 letters
                </span>
              </div>
            </div>
          </div>
          
          <!-- Dismiss Button -->
          <div class="flex justify-center">
            <BaseButton
              variant="primary"
              size="md"
              @click="dismissGameEndPopup"
              class="min-w-24"
            >
              OK
            </BaseButton>
          </div>
          
          <!-- Auto-dismiss progress bar -->
          <div class="mt-4 w-full bg-[var(--color-background)] rounded-full h-2">
            <div 
              class="bg-[var(--color-primary)] h-2 rounded-full transition-all duration-100 ease-linear"
              :style="{ width: `${gameEndProgress}%` }"
            ></div>
          </div>
          <p class="text-xs text-[var(--color-text-tertiary)] mt-2">
            Auto-dismissing in {{ Math.ceil(gameEndProgress / 10) }}s
          </p>
        </div>
      </div>
    </div>
    
    <!-- Main Game Area -->
    <main class="flex-1 relative pt-16 overflow-hidden" @click="handleMainAreaClick">
      <!-- Game Viewport -->
      <div class="relative w-full game-viewport transition-all duration-300" 
           :style="{ 
             background: 'radial-gradient(ellipse at center, var(--color-background-secondary) 0%, var(--color-background) 70%)',
             height: 'calc(100vh - 4rem)'
           }">
        
        <!-- Opponent Grid (Top 60% of viewport, Full width) -->
        <div 
          v-if="displayedGameState && displayedGameState.game && displayedGameState.game.status === 'active'"
          class="opponent-area absolute left-0 overflow-hidden transition-transform duration-300"
          style="top: 2vh; width: 100vw; height: 45vh;"
        >
          <!-- Opponent Cells -->
          <div
            v-for="opponent in opponentPlayers"
            :key="opponent.id"
            class="opponent-cell absolute left-0 items-center"
            :style="{
              top: `calc(${opponent.seat.top} * 0.75)`,
              width: opponent.seat.width,
              height: `calc(${opponent.seat.cellHeight} * 0.75)`,
              display: 'grid',
              gridTemplateColumns: 'auto 1fr auto',
              alignItems: 'center',
              paddingLeft: '5vw',
              paddingRight: '5vw'
            }"
          >
            <!-- Opponent Name (Left-justified) -->
            <div class="opponent-name-area flex items-center justify-start" 
                 :class="{ 'opacity-60': isPlayerFinished(opponent.id) }"
                 :style="{ height: '100%', gridColumn: '1', gridRow: '1', zIndex: 2 }">
              <div class="flex flex-col items-start justify-center h-full">
                <div 
                  class="px-3 py-1 rounded-full text-sm font-medium transition-all duration-500 relative"
                  :class="{
                    'bg-[var(--color-primary)] text-white border-2 border-[var(--color-primary)] shadow-lg scale-110 animate-pulse': isCurrentTurn(opponent.id),
                    'bg-green-100 border-2 border-green-500 text-green-800 shadow-md': isPlayerFinished(opponent.id) && !isCurrentTurn(opponent.id),
                    'bg-[var(--color-surface)] border border-[var(--color-border)]': !isCurrentTurn(opponent.id) && !isPlayerFinished(opponent.id)
                  }"
                >
                  <!-- Turn indicator bouncing arrow -->
                  <div v-if="isCurrentTurn(opponent.id)" 
                       class="absolute -left-8 top-1/2 transform -translate-y-1/2 text-[var(--color-primary)] text-lg animate-horizontal-bounce">
                    ▶
                  </div>
                  
                  {{ opponent.name }}
                  <span v-if="isPlayerFinished(opponent.id)" class="ml-1">🎉</span>
                  <span v-if="hasAceOfSpades(opponent.id)" class="text-yellow-400 ml-1">♠A</span>
                </div>
                
                <!-- Finished player hint -->
                <div v-if="isPlayerFinished(opponent.id)" class="text-xs text-green-600 font-medium mt-1 px-2">
                  ✓ Finished - Waiting
                </div>
              </div>
            </div>
            
            <!-- Opponent Card Fan (Center-justified, spans middle column) -->
            <div class="opponent-cards-area flex items-center justify-center" 
                 :style="{ height: '100%', gridColumn: '1 / -1', gridRow: '1', justifySelf: 'center', zIndex: 1 }">
              <CardFan
                :cards="getPlayerCards(opponent)"
                :back-color="cardBackColor"
                :interactive="false"
                :is-user="false"
                size="opponent"
                @click.stop
                class="h-full w-full flex items-center justify-center"
              />
            </div>
            
            <!-- Opponent In-Game Card (Right-justified) -->
            <div class="opponent-ingame-area flex items-center justify-end" 
                 :style="{ height: '100%', gridColumn: '3', gridRow: '1', justifySelf: 'end', zIndex: 2 }">
              <div v-if="getOpponentInGameCard(opponent.id)" class="opponent-ingame-card border-2 border-gray-300 border-opacity-50 rounded-lg p-1">
                <GameCard
                  :rank="getOpponentInGameCard(opponent.id).card.rank"
                  :suit="getOpponentInGameCard(opponent.id).card.suit"
                  size="ingame-opponent"
                  :flipped="false"
                  :interactive="false"
                  :disabled="false"
                  :highest-highlight="isHighestInGameCard(getOpponentInGameCard(opponent.id))"
                />
              </div>
            </div>
          </div>
        </div>
        
        
        <!-- Current Player In-Game Card (horizontally with player area, vertically aligned with opponents) -->
        <div v-if="getCurrentPlayerInGameCard() && displayedGameState && displayedGameState.game && displayedGameState.game.status === 'active'" 
             class="current-player-ingame-area absolute flex items-center justify-end" 
             style="bottom: 35vh; right: 5vw; height: 10vh; z-index: 20;">
          <div class="current-player-ingame-card border-2 border-gray-300 border-opacity-50 rounded-lg p-1">
            <GameCard
              :rank="getCurrentPlayerInGameCard().card.rank"
              :suit="getCurrentPlayerInGameCard().card.suit"
              size="ingame-opponent"
              :flipped="false"
              :interactive="false"
              :disabled="false"
              :highest-highlight="isHighestInGameCard(getCurrentPlayerInGameCard())"
            />
          </div>
        </div>
        
        <!-- Current Player Area (Bottom 40% of viewport, Full width) -->
        <div 
          v-if="displayedGameState && displayedGameState.game && displayedGameState.game.status === 'active'"
          class="current-player-area absolute bottom-0 left-0 w-full"
          style="height: 55vh;"
        >
          
          <!-- Current Player Card Fan (Top section of player area, 70% width, left-positioned to prevent clipping) -->
          <div class="current-player-cards absolute top-0 flex justify-start" 
               style="left: 15vw; width: 70vw; height: 60%;">
            <CardFan
              v-if="displayedGameState.game && displayedGameState.game.status === 'active'"
              :cards="getPlayerCards(currentPlayerData)"
              :back-color="cardBackColor"
              :interactive="true"
              :is-user="true"
              @card-click="handleCardClick"
              @card-select="handleCardSelect"
              @click.stop
              class="w-full h-full"
            />
          </div>
          
          <!-- Current Player Name (Below card fan, moved up 5% height to be closer to cards) -->
          <div class="current-player-name absolute flex items-center justify-center flex-col" 
               :class="{ 'opacity-70': isPlayerFinished(user.id) }"
               style="top: calc(55% - 5vh); left: 25vw; width: 50vw; height: 40%;">
            <div 
              class="px-3 py-1 rounded-full text-sm font-medium transition-all duration-500 relative"
              :class="{
                'bg-[var(--color-primary)] text-white border-2 border-[var(--color-primary)] shadow-lg scale-110 animate-pulse': isCurrentTurn(user.id),
                'bg-green-100 border-2 border-green-500 text-green-800 shadow-md': isPlayerFinished(user.id) && !isCurrentTurn(user.id),
                'bg-[var(--color-surface)] border border-[var(--color-border)]': !isCurrentTurn(user.id) && !isPlayerFinished(user.id)
              }"
            >
              <!-- Turn indicator bouncing arrow -->
              <div v-if="isCurrentTurn(user.id)" 
                   class="absolute -left-8 top-1/2 transform -translate-y-1/2 text-[var(--color-primary)] text-lg animate-horizontal-bounce">
                ▶
              </div>
              
              {{ user.name }}
              <span class="opacity-75"> (You)</span>
              <span v-if="isPlayerFinished(user.id)" class="ml-1">🎉</span>
              <span v-if="hasAceOfSpades(user.id)" class="text-yellow-400 ml-1">♠A</span>
            </div>
            
            <!-- Finished player hint for current user -->
            <div v-if="isPlayerFinished(user.id)" class="text-xs text-green-600 font-medium mt-2 px-3 py-1 bg-green-50 rounded-full border border-green-200">
              ✓ You've finished! Waiting for others...
            </div>
          </div>
        </div>
        
        <!-- Center Game Controls -->
        <div class="absolute text-center z-10"
             style="top: 40%; left: 50%; transform: translate(-50%, -50%);"
             :class="{ 'mt-32': displayedGameState && displayedGameState.inPlayCards && displayedGameState.inPlayCards.length > 0 }">
          <!-- Waiting for Game Start -->
          <div v-if="!displayedGameState" class="space-y-4">
            <p class="text-lg text-[var(--color-text-secondary)]">Ready to play?</p>
            <BaseButton
              variant="primary"
              size="lg"
              @click="startNewGame"
              :loading="loading"
            >
              Start New Game
            </BaseButton>
          </div>
          
          <!-- Game Setup Phase -->
          <div v-else-if="displayedGameState.game && displayedGameState.game.status === 'waiting' && isRequester" class="space-y-4">
            <p class="text-lg text-[var(--color-text-secondary)]">
              Players: {{ displayedGameState.players.length }}/8
            </p>
            <BaseButton
              variant="success"
              size="lg"
              :disabled="displayedGameState.players.length < 2"
              @click="handleFinalizeGame"
              :loading="loading"
            >
              <span class="text-center">
                Start Game<br>
                ({{ displayedGameState.players.length }} players)
              </span>
            </BaseButton>
            <BaseButton
              variant="secondary"
              size="md"
              @click="showShare = true"
            >
              Share Game Link
            </BaseButton>
            
            <!-- Bot Management -->
            <div class="border-t border-[var(--color-border)] pt-3">
              <div class="flex items-center justify-between mb-2">
                <h3 class="text-sm font-medium text-[var(--color-text)]">Bots ({{ botPlayers.length }})</h3>
                <BaseButton
                  variant="outline"
                  size="sm"
                  @click="addBotPlayer('easy')"
                  :disabled="botPlayers.length >= maxBotCap || displayedGameState.players.length >= (displayedGameState.game?.maxPlayers || 8) || loading"
                >
                  Add Bot
                </BaseButton>
              </div>
              
              <!-- Bot list without fixed height container -->
              <div class="space-y-1">
                <div 
                  v-for="bot in botPlayers" 
                  :key="bot.id"
                  class="flex items-center justify-between text-xs bg-[var(--color-bg-secondary)] px-2 py-1.5 rounded relative"
                >
                  <span class="text-[var(--color-text-secondary)] truncate flex-1 mr-2">
                    {{ bot.name }}
                  </span>
                  <div class="flex items-center gap-1">
                    <!-- Custom dropdown to control positioning -->
                    <div class="relative">
                      <button
                        @click="toggleBotDifficultyDropdown(bot.id)"
                        class="text-xs bg-transparent border border-[var(--color-border)] rounded px-1 py-0.5 text-[var(--color-text-secondary)] min-w-16 text-left"
                        :disabled="loading"
                      >
                        {{ bot.botDifficulty === 'difficult' ? 'Hard' : bot.botDifficulty.charAt(0).toUpperCase() + bot.botDifficulty.slice(1) }}
                      </button>
                      <div 
                        v-if="dropdownOpen === bot.id"
                        class="absolute right-0 top-full mt-1 bg-[var(--color-surface)] border border-[var(--color-border)] rounded shadow-lg z-50 min-w-20"
                      >
                        <button
                          v-for="difficulty in ['easy', 'medium', 'difficult']"
                          :key="difficulty"
                          @click="changeBotDifficulty(bot.id, difficulty)"
                          class="block w-full text-left text-xs px-2 py-1 hover:bg-[var(--color-bg-secondary)] text-[var(--color-text-secondary)]"
                        >
                          {{ difficulty === 'difficult' ? 'Hard' : difficulty.charAt(0).toUpperCase() + difficulty.slice(1) }}
                        </button>
                      </div>
                    </div>
                    <button
                      @click="removeBotPlayer(bot.id)"
                      class="text-[var(--color-danger)] hover:text-[var(--color-danger-hover)] ml-1"
                      :disabled="loading"
                      title="Remove bot"
                    >
                      ×
                    </button>
                  </div>
                </div>
                <div v-if="botPlayers.length === 0" class="text-xs text-[var(--color-text-tertiary)] text-center py-2">
                  No bots added
                </div>
              </div>
            </div>
          </div>
          
          <!-- Waiting for Host -->
          <div v-else-if="displayedGameState.game && displayedGameState.game.status === 'waiting' && !isRequester" class="space-y-4">
            <p class="text-lg text-[var(--color-text-secondary)]">
              Waiting for game creator to start...
            </p>
            <p class="text-sm text-[var(--color-text-tertiary)]">
              Players: {{ displayedGameState.players.length }}/8
            </p>
          </div>
          
          <!-- Active Game Controls -->
          <div v-else-if="displayedGameState.game && displayedGameState.game.status === 'active' && allRemainingPlayersAreBots" class="space-y-2">
            <p class="text-sm text-[var(--color-text-secondary)] mb-2">
              All remaining players are bots
            </p>
            <BaseButton
              variant="outline"
              size="sm"
              @click="toggleSpeedUp"
              :class="{
                'bg-[var(--color-primary)] text-white border-[var(--color-primary)]': speedUpEnabled,
                'bg-transparent border-[var(--color-border)] text-[var(--color-text)]': !speedUpEnabled
              }"
            >
              <span class="flex items-center gap-2">
                {{ speedUpEnabled ? '⚡ Speed Up ON' : '⚡ Speed Up OFF' }}
              </span>
            </BaseButton>
          </div>
        </div>
      </div>
    </main>
    
    
    <!-- Session Log -->
    <div v-if="gameId" :class="showLog ? 'fixed bottom-0 left-0 right-0 z-50' : 'fixed bottom-4 left-[5%] right-[5%] z-5'">
      <div :class="showLog ? 'bg-[var(--color-surface)] border-t border-[var(--color-border)] shadow-lg rounded-t-lg h-48' : 'bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg shadow-game'">
        
        <!-- Minimized View - Single line with latest message and + button -->
        <div v-if="!showLog" class="flex items-center justify-center px-3 py-0" style="padding-top: 0.165rem; padding-bottom: 0.165rem;">
          <div class="flex-1 text-xs text-[var(--color-text-secondary)] truncate pr-2">
            {{ logs.length > 0 ? logs[0].message : 'No messages yet...' }}
          </div>
          <BaseButton
            variant="ghost"
            size="xs"
            icon-only
            @click="showLog = true"
            class="flex-shrink-0"
          >
            <template #icon>+</template>
          </BaseButton>
        </div>
        
        <!-- Expanded View with Header -->
        <div v-if="showLog" class="h-full flex flex-col">
          <!-- Log Header -->
          <div class="flex justify-between items-center px-3 py-0 border-b border-[var(--color-border)] bg-[var(--color-background-secondary)] rounded-t-lg" style="padding-top: 0.33rem; padding-bottom: 0.33rem;">
            <span class="font-medium text-sm">Chat Logs</span>
            <BaseButton
              variant="ghost"
              size="xs"
              icon-only
              @click="showLog = false"
            >
              <template #icon>−</template>
            </BaseButton>
          </div>
          
          <!-- Messages -->
          <div class="flex-1 overflow-y-auto px-3 py-2 text-sm space-y-1 min-h-0">
            <div
              v-for="log in logs"
              :key="log.id"
              class="text-xs leading-tight text-[var(--color-text-secondary)]"
            >
              {{ log.message }}
            </div>
            <div v-if="!logs.length" class="text-xs text-[var(--color-text-tertiary)]">
              No messages yet...
            </div>
          </div>
          
          <!-- Chat Input -->
          <div class="flex border-t border-[var(--color-border)] bg-[var(--color-surface)]">
            <input
              v-model="chatMessage"
              @keydown.enter="sendChat"
              placeholder="Type message..."
              maxlength="128"
              class="flex-1 px-3 py-2 text-sm bg-transparent border-0 focus:outline-none focus:ring-1 focus:ring-[var(--color-primary)]"
            />
            <BaseButton
              variant="primary"
              size="sm"
              :disabled="!chatMessage.trim()"
              @click="sendChat"
              class="rounded-none"
            >
              Send
            </BaseButton>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Modals -->
    <!-- Registration Modal -->
    <BaseModal
      :is-open="showRegistration"
      title="Join the Game"
      size="md"
      :close-on-backdrop="false"
      :close-on-escape="false"
      :show-close="false"
    >
      <div class="space-y-4">
        <div>
          <label class="form-label">Your Name</label>
          <input
            v-model="playerName"
            @keydown.enter="register"
            placeholder="Enter your name"
            maxlength="20"
            class="form-input w-full"
            autofocus
          />
        </div>
      </div>
      
      <template #footer>
        <BaseButton
          variant="primary"
          :disabled="!playerName.trim()"
          @click="register"
          :loading="loading"
          block
        >
          Join Game
        </BaseButton>
      </template>
    </BaseModal>
    
    <!-- Share Game Modal -->
    <BaseModal
      :is-open="showShare"
      title="Share Game"
      @close="showShare = false"
    >
      <div class="space-y-4">
        <div>
          <label class="form-label">Game Link</label>
          <div class="flex">
            <input
              :value="gameShareUrl"
              readonly
              class="form-input flex-1 rounded-r-none"
            />
            <BaseButton
              variant="secondary"
              @click="copyGameLink"
              class="rounded-l-none"
            >
              Copy
            </BaseButton>
          </div>
        </div>
        <div class="text-sm text-[var(--color-text-secondary)]">
          Share this link with other players to join your game.
        </div>
      </div>
    </BaseModal>
    
    <!-- Settings Modal -->
    <BaseModal
      :is-open="showSettings"
      title="Settings"
      @close="showSettings = false"
    >
      <div class="space-y-4">
        <div>
          <label class="form-label">Card Back Color</label>
          <select v-model="cardBackColor" class="form-input w-full">
            <option v-for="color in cardBackColors" :key="color" :value="color">
              {{ color.charAt(0).toUpperCase() + color.slice(1) }}
            </option>
          </select>
        </div>
        
      </div>
    </BaseModal>
    
    <!-- Abandoned Game Modal -->
    <BaseModal
      :is-open="showAbandoned"
      title="Game Abandoned"
      :close-on-backdrop="false"
      :close-on-escape="false"
    >
      <p class="text-[var(--color-text-secondary)]">
        This game has been abandoned by the creator.
      </p>
      
      <template #footer>
        <BaseButton
          variant="primary"
          @click="closeAbandoned"
          block
        >
          Start New Game
        </BaseButton>
      </template>
    </BaseModal>
    
    <!-- About Modal -->
    <BaseModal
      :is-open="showAbout"
      title="About Donkey"
      @close="showAbout = false"
    >
      <div class="text-center space-y-4">
        <h2 class="text-xl font-bold text-[var(--color-text)]">Donkey v0.2</h2>
        <p class="text-[var(--color-text-secondary)]">Created by Deepak Amin</p>
        <p class="text-sm text-[var(--color-text-secondary)]">
          A multiplayer card game where players try to avoid collecting cards and getting DONKEY letters.
        </p>
      </div>
    </BaseModal>
    
    <!-- Help Modal -->
    <BaseModal
      :is-open="showHelp"
      title="How to Play Donkey"
      @close="showHelp = false"
    >
      <div class="space-y-4 text-[var(--color-text-secondary)]">
        <div>
          <h3 class="font-semibold text-[var(--color-text)] mb-2">Objective</h3>
          <p>Avoid collecting cards and getting DONKEY letters. The last player without all letters wins!</p>
        </div>
        
        <div>
          <h3 class="font-semibold text-[var(--color-text)] mb-2">Gameplay</h3>
          <ul class="list-disc list-inside space-y-1">
            <li>Cards are dealt evenly to all players</li>
            <li>Players take turns playing one card at a time</li>
            <li>Must follow suit if possible</li>
            <li>Highest card of the lead suit wins the trick</li>
            <li>Winner of the trick leads the next trick</li>
            <li>Player who collects the most cards in a round gets a DONKEY letter</li>
          </ul>
        </div>
        
        <div>
          <h3 class="font-semibold text-[var(--color-text)] mb-2">Winning</h3>
          <p>Keep playing rounds until all but one player has collected all 6 letters (D-O-N-K-E-Y). The remaining player wins!</p>
        </div>
        
        <div>
          <h3 class="font-semibold text-[var(--color-text)] mb-2">Tips</h3>
          <ul class="list-disc list-inside space-y-1">
            <li>Try to avoid taking tricks when possible</li>
            <li>Count cards to know what's still in play</li>
            <li>Watch what other players are void in</li>
          </ul>
        </div>
      </div>
    </BaseModal>
  </div>
</template>

<script>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import BaseButton from './components/BaseButton.vue'
import BaseModal from './components/BaseModal.vue'
import CardFan from './components/CardFan.vue'
import GameCard from './components/GameCard.vue'
import { useTheme } from './composables/useTheme.js'
import { useAnimationSystem } from './composables/useAnimationSystem.js'
import {
  getCookie, setCookie, cardBackColors, getPlayerSeats, sortCards, parseCards,
  registerUser, createGame, startGame, joinGame, finalizeGame as apiFinalizeGame, abandonGame, 
  sendChatMessage, getGameState, getGameLogs, getUser, getVersion, addBot, playCard as apiPlayCard
} from './utils/gameUtils.js'

export default {
  name: 'App',
  components: {
    BaseButton,
    BaseModal,
    CardFan,
    GameCard
  },
  
  setup() {
    // Theme management
    const { currentTheme, themeLabel, cycleTheme, setTheme } = useTheme()
    
    // Initialize animation system
    const animationSystem = useAnimationSystem()
    
    // Game state
    const user = ref({ id: null, name: null })
    const gameId = ref(null)
    const gameState = ref(null)
    const logs = ref([])
    const connected = ref(false)
    const loading = ref(false)
    
    // UI state
    const showMenu = ref(false)
    const showRegistration = ref(false)
    const showShare = ref(false)
    const showSettings = ref(false)
    const showAbandoned = ref(false)
    const showAbout = ref(false)
    const showHelp = ref(false)
    const dropdownOpen = ref(null)
    const showLog = ref(false)
    const speedUpEnabled = ref(false)
    
    // Turn result notification state
    const turnResultNotification = ref(null)
    const turnResultProgress = ref(100)
    let turnResultTimer = null
    let turnResultProgressTimer = null
    
    // Game end popup state
    const gameEndPopup = ref(null)
    const gameEndProgress = ref(100)
    let gameEndTimer = null
    let gameEndProgressTimer = null
    
    // Form data
    const playerName = ref('')
    const chatMessage = ref('')
    const cardBackColor = ref(getCookie('cardBack') || 'red')
    const selectedTheme = ref(getCookie('theme') || 'system')
    
    // Card selection state
    const selectedCard = ref(null)
    const disabledCardFeedback = ref(null)
    
    // Animation states
    const gameAnimationState = ref({
      isShowingCut: false,
      isMovingToDiscard: false,
      isMovingToWinner: false,
      cutWinnerPlayerId: null,
      animatingCards: []
    })
    
    // Backend now handles all timing - no frontend blocking needed
    
    // Use gameState directly - backend handles all timing
    const displayedGameState = computed(() => {
      return gameState.value
    })
    
    // Remove communication blocking - backend now handles all timing
    
    // Connection management
    const eventSource = ref(null)
    const connectionKey = ref(0)
    
    // Computed properties
    const isRequester = computed(() => {
      return gameState.value && gameState.value.game && gameState.value.game.requesterId === user.value.id
    })
    
    const playersWithSeats = computed(() => {
      if (!displayedGameState.value) return []
      return getPlayerSeats(displayedGameState.value.players, user.value.id)
    })
    
    // Separate opponents from current player for grid layout
    const opponentPlayers = computed(() => {
      if (!playersWithSeats.value) return []
      return playersWithSeats.value.filter(player => player.isOpponent)
    })
    
    const currentPlayerData = computed(() => {
      if (!playersWithSeats.value) return null
      return playersWithSeats.value.find(player => player.isCurrentUser)
    })
    
  const botPlayers = computed(() => {
    if (!displayedGameState.value) return []
    return displayedGameState.value.players.filter(player => player.isBot)
  })

  // Dynamic bot capacity: TOTAL - CREATOR - HUMAN_PLAYERS(excluding creator)
  const totalSlots = computed(() => displayedGameState.value?.game?.maxPlayers || 8)
  const creatorId = computed(() => displayedGameState.value?.game?.requesterId || null)
  const humanPlayers = computed(() => (displayedGameState.value?.players || []).filter(p => !p.isBot))
  const otherHumansCount = computed(() => humanPlayers.value.filter(p => p.id !== creatorId.value).length)
  const maxBotCap = computed(() => Math.max(0, totalSlots.value - 1 - otherHumansCount.value))
    
    const gameShareUrl = computed(() => {
      if (!gameId.value) return ''
      return `${window.location.origin}?gameId=${gameId.value}`
    })
    
    const isCurrentTurn = (playerId) => {
      if (!displayedGameState.value || !displayedGameState.value.currentTurn) return false
      return displayedGameState.value.currentTurn.expectedPlayerId === playerId
    }
    
    const isPlayerFinished = (playerId) => {
      if (!displayedGameState.value || !displayedGameState.value.roundPlayers) return false
      const roundPlayer = displayedGameState.value.roundPlayers.find(rp => rp.userId === playerId)
      return roundPlayer && roundPlayer.isFinished
    }
    
    // Helper function to get card rank value for comparison
    const getCardValue = (card) => {
      if (!card || !card.rank) return 0
      switch (card.rank) {
        case 'A': return 14
        case 'K': return 13
        case 'Q': return 12
        case 'J': return 11
        default: return parseInt(card.rank, 10) || 0
      }
    }
    
    // Determine which In-Game card has the highest value for this round
    const getHighestInGameCard = () => {
      if (!displayedGameState.value || !displayedGameState.value.inPlayCards || displayedGameState.value.inPlayCards.length === 0) {
        return null
      }
      
      const inPlayCards = displayedGameState.value.inPlayCards
      let highestCard = null
      let highestValue = 0
      
      for (const playedCard of inPlayCards) {
        const card = playedCard.card
        const cardValue = getCardValue(card)
        
        if (cardValue > highestValue) {
          highestValue = cardValue
          highestCard = playedCard
        }
      }
      
      return highestCard
    }
    
    // Check if a specific played card is the highest value
    const isHighestInGameCard = (playedCard) => {
      const highest = getHighestInGameCard()
      return highest && playedCard && 
             highest.playerId === playedCard.playerId &&
             highest.card.rank === playedCard.card.rank &&
             highest.card.suit === playedCard.card.suit
    }
    
    const hasAceOfSpades = (playerId) => {
      // Remove Ace of Spades indicator - we now have proper turn management
      // The turn indicator arrow clearly shows whose turn it is
      return false
    }
    
    const allRemainingPlayersAreBots = computed(() => {
      // Show only when the ONLY active players left in the round are bots
      if (!displayedGameState.value || displayedGameState.value.game?.status !== 'active') return false
      const activePlayers = (displayedGameState.value.players || []).filter(p => !isPlayerFinished(p.id))
      if (activePlayers.length === 0) return false
      return activePlayers.every(p => p.isBot === true)
    })
    
    const menuItems = computed(() => [
      !gameId.value && { label: 'New Game', action: 'newGame' },
      gameId.value && isRequester.value && { label: 'Abandon Game', action: 'abandon' },
      gameId.value && !isRequester.value && { label: 'New Game', disabled: true },
      { label: 'Settings', action: 'settings' },
      { label: 'Help', action: 'help' },
      { label: 'About', action: 'about' }
    ].filter(Boolean))
    
    // Watchers
    watch(cardBackColor, (newColor) => {
      setCookie('cardBack', newColor)
    })
    
    watch(selectedTheme, (newTheme) => {
      setTheme(newTheme)
      setCookie('theme', newTheme)
    })
    
    watch(() => gameState.value?.game?.status === 'abandoned', (isAbandoned) => {
      if (isAbandoned) {
        showAbandoned.value = true
      }
    })
    
    // Watch for bot turns and show indication (backend handles bot automation automatically)
    watch(() => gameState.value?.currentTurn?.expectedPlayerId, (expectedPlayerId) => {
      if (!expectedPlayerId || !gameState.value) return
      
      // Find the player who should play
      const expectedPlayer = gameState.value.players.find(p => p.id === expectedPlayerId)
      if (!expectedPlayer || !expectedPlayer.isBot) return
      
      // Note: Backend handles bot automation automatically when cards are played
    })
    
    // Clear disabled feedback when game state changes (new round, turn changes, etc.)
    watch(() => gameState.value?.currentRound?.roundNumber, () => {
      if (disabledCardFeedback.value) {
        disabledCardFeedback.value = null
      }
    })
    
    // Clear disabled feedback when current turn changes
    watch(() => gameState.value?.currentTurn?.turnNumber, () => {
      if (disabledCardFeedback.value) {
        disabledCardFeedback.value = null
      }
    })
    
    // Methods
    const toggleMenu = () => {
      showMenu.value = !showMenu.value
    }
    
    const handleMenuAction = (action) => {
      showMenu.value = false
      
      switch (action) {
        case 'newGame':
          startNewGame()
          break
        case 'abandon':
          abandon()
          break
        case 'settings':
          showSettings.value = true
          break
        case 'about':
          showAboutModal()
          break
        case 'help':
          showHelpModal()
          break
      }
    }
    
    const startNewGame = async () => {
      if (!user.value.id) {
        showRegistration.value = true
        return
      }
      
      try {
        loading.value = true
        const response = await createGame(user.value.id)
        gameId.value = response.gameId
      } catch (error) {
      } finally {
        loading.value = false
      }
    }
    
    const register = async () => {
      if (!playerName.value.trim()) return
      
      try {
        loading.value = true
        const userData = await registerUser(playerName.value.trim())
        user.value = userData
        setCookie('userId', userData.id)
        setCookie('userName', userData.name)
        showRegistration.value = false
        
        // Auto-join if there's a game ID in URL
        const urlParams = new URLSearchParams(window.location.search)
        const gameIdParam = urlParams.get('gameId')
        if (gameIdParam) {
          await joinGameById(gameIdParam)
          window.history.replaceState(null, '', window.location.pathname)
        }
      } catch (error) {
      } finally {
        loading.value = false
      }
    }
    
    const joinGameById = async (gid) => {
      try {
        await joinGame(gid, user.value.id)
        gameId.value = gid
      } catch (error) {
      }
    }
    
    const handleFinalizeGame = async () => {
      try {
        loading.value = true
        await startGame(gameId.value, user.value.id)
        // Refresh game state after starting
        await fetchGameState()
      } catch (error) {
      } finally {
        loading.value = false
      }
    }
    
    const abandon = async () => {
      try {
        await abandonGame(gameId.value, user.value.id)
        resetGame()
      } catch (error) {
      }
    }
    
    const sendChat = async () => {
      const message = chatMessage.value.trim()
      if (!message) return
      
      try {
        await sendChatMessage(gameId.value, user.value.id, message)
        chatMessage.value = ''
      } catch (error) {
      }
    }
    
    const fetchGameState = async () => {
      if (!gameId.value || !user.value.id) return
      
      try {
        const state = await getGameState(gameId.value, user.value.id)
        gameState.value = state
      } catch (error) {
      }
    }
    
    const fetchLogs = async () => {
      if (!gameId.value) return
      
      try {
        const logData = await getGameLogs(gameId.value)
        logs.value = logData
      } catch (error) {
      }
    }
    
    // Bot management functions
    const addBotPlayer = async (difficulty = 'easy') => {
      if (!gameId.value || !user.value.id) return
      
      try {
        loading.value = true
        await addBot(gameId.value, user.value.id, difficulty)
        await fetchGameState() // Refresh game state to show new bot
      } catch (error) {
      } finally {
        loading.value = false
      }
    }
    
    const removeBotPlayer = async (botId) => {
      if (!gameId.value || !user.value.id) return
      
      try {
        loading.value = true
        // Note: We'll need to implement a remove bot API endpoint
        
        // TODO: Call remove bot API when available
        await fetchGameState() // Refresh game state
      } catch (error) {
      } finally {
        loading.value = false
      }
    }
    
    const changeBotDifficulty = async (botId, newDifficulty) => {
      if (!gameId.value || !user.value.id) return
      
      dropdownOpen.value = null // Close dropdown
      
      try {
        loading.value = true
        // Note: We'll need to implement a change bot difficulty API endpoint
        
        // TODO: Call change bot difficulty API when available
        await fetchGameState() // Refresh game state
      } catch (error) {
      } finally {
        loading.value = false
      }
    }
    
    const toggleBotDifficultyDropdown = (botId) => {
      dropdownOpen.value = dropdownOpen.value === botId ? null : botId
    }
    
    // triggerBotPlay function removed - backend handles bot automation automatically
    
    const setupEventSource = () => {
      if (!gameId.value || !user.value.id) return
      
      if (eventSource.value) {
        eventSource.value.close()
      }
      
      eventSource.value = new EventSource(`/api/game/${gameId.value}/stream/${user.value.id}`)
      
      eventSource.value.onopen = () => {
        connected.value = true
      }
      
      eventSource.value.onerror = () => {
        connected.value = false
        eventSource.value.close()
        setTimeout(() => {
          connectionKey.value++
          setupEventSource()
        }, 2000)
      }
      
      eventSource.value.onmessage = (event) => {
        const data = JSON.parse(event.data)
        
        if (data.type === 'state') {
          // Backend handles timing - always fetch state updates
          fetchGameState()
        }
        
        if (data.type === 'log') {
          logs.value.unshift(data.log)

          // Trigger CUT indication animation and notifications if event contains structured data
          try {
            const log = data.log || {}
            if (log.type === 'turn_event' && log.eventData) {
              const ed = typeof log.eventData === 'string' ? JSON.parse(log.eventData) : log.eventData
              
              if (ed && ed.type === 'cut') {
                // Show CUT notification
                const cutterName = ed.cutPlayerName || ed.cutPlayerId || 'Unknown'
                const winnerName = ed.winnerName || ed.winnerPlayerId || 'Unknown'
                const message = `<strong>CUT by ${cutterName}.</strong><br>Previous highest card was played by ${winnerName}.`
                showTurnResultNotification(message, 'cut')
                
                // Add to chat logs
                const chatMessage = `CUT by ${cutterName}. Previous highest card was played by ${winnerName}.`
                const chatLog = {
                  id: `cut-${Date.now()}`,
                  gameId: gameId.value,
                  type: 'system',
                  message: chatMessage,
                  createdAt: new Date().toISOString()
                }
                logs.value.unshift(chatLog)
                
                // Trigger animation system if available
                if (animationSystem && animationSystem.showCutIndication) {
                  animationSystem.showCutIndication({
                    cutPlayerId: ed.cutPlayerId,
                    winnerPlayerId: ed.winnerPlayerId,
                    cardsCollected: ed.cardsCollected || (displayedGameState.value?.inPlayCards?.length || 0)
                  })
                }
              } else if (ed && ed.type === 'discard') {
                // Show discard notification
                const winnerName = ed.winnerName || ed.winnerPlayerId || 'Unknown'
                const message = `<strong>No CUT.</strong> Moving cards to discard pile.<br>Highest card played by ${winnerName}.`
                showTurnResultNotification(message, 'discard')
                
                // Add to chat logs
                const chatMessage = `No CUT. Moving cards to discard pile. Highest card played by ${winnerName}.`
                const chatLog = {
                  id: `discard-${Date.now()}`,
                  gameId: gameId.value,
                  type: 'system',
                  message: chatMessage,
                  createdAt: new Date().toISOString()
                }
                logs.value.unshift(chatLog)
              }
            } else if (log.type === 'turn_event' && log.message && log.message.includes('Discarded')) {
              // Fallback: detect discard from message text if eventData is missing
              const match = log.message.match(/Discarded \d+ cards\. (.+?) starts next turn\./)
              if (match) {
                const winnerName = match[1] || 'Unknown'
                const message = `<strong>No CUT.</strong> Moving cards to discard pile.<br>Highest card played by ${winnerName}.`
                showTurnResultNotification(message, 'discard')
                
                // Add to chat logs
                const chatMessage = `No CUT. Moving cards to discard pile. Highest card played by ${winnerName}.`
                const chatLog = {
                  id: `discard-fallback-${Date.now()}`,
                  gameId: gameId.value,
                  type: 'system',
                  message: chatMessage,
                  createdAt: new Date().toISOString()
                }
                logs.value.unshift(chatLog)
              }
            } else if (log.type === 'game_event' && log.eventData) {
              const ed = typeof log.eventData === 'string' ? JSON.parse(log.eventData) : log.eventData
              
              if (ed && ed.type === 'game_end') {
                // Show game end popup with scoreboard
                showGameEndPopup({
                  loserId: ed.loserId,
                  loserName: ed.loserName,
                  scoreboard: ed.scoreboard || []
                })
              }
            }
          } catch (e) {
            // Silently ignore log parsing errors
          }
        }
      }
    }
    
    const getPlayerCards = (player) => {
      if (player.id === user.value.id) {
        // Get user's cards directly from displayedGameState.myCards (they already include ID)
        if (displayedGameState.value && displayedGameState.value.myCards) {
          // Get played card IDs to filter them out from hand (backup safety check)
          const playedCardIds = new Set()
          ;(displayedGameState.value.inPlayCards || []).forEach(pc => {
            // Try both possible ID locations
            if (pc.id) playedCardIds.add(pc.id)
            if (pc.card && pc.card.id) playedCardIds.add(pc.card.id)
          })
          
          const availableCards = displayedGameState.value.myCards
            .filter(card => !playedCardIds.has(card.id)) // Remove any played cards from hand
          
          // Clear selection if selected card was played
          if (selectedCard.value && playedCardIds.has(selectedCard.value.id)) {
            selectedCard.value = null
          }
          
          // Check if it's the current player's turn
          const isMyTurn = isCurrentTurn(user.value.id)
          
          // If it's not my turn, all cards should be disabled
          if (!isMyTurn) {
            return availableCards.map((card, index) => ({
              id: card.id,
              rank: card.rank,
              suit: card.suit,
              code: card.code,
              selected: selectedCard.value && selectedCard.value.id === card.id,
              disabledFeedback: disabledCardFeedback.value && disabledCardFeedback.value.id === card.id,
              cardIndex: index,
              flipped: false,
              disabled: true // All cards disabled when not my turn
            }))
          }
          
          // Check for special Ace of Spades rule (first turn of first round)
          const isFirstTurnFirstRound = displayedGameState.value.currentRound && 
                                      displayedGameState.value.currentRound.roundNumber === 1 &&
                                      displayedGameState.value.currentTurn && 
                                      displayedGameState.value.currentTurn.turnNumber === 1 &&
                                      (!displayedGameState.value.inPlayCards || displayedGameState.value.inPlayCards.length === 0)
          
          if (isFirstTurnFirstRound) {
            // First turn of first round: Only Ace of Spades can be played
            
            return availableCards.map((card, index) => {
              const isAceOfSpades = card.rank === 'A' && card.suit === 'spades'
              const isDisabled = !isAceOfSpades
              
              
              
              return {
                id: card.id,
                rank: card.rank,
                suit: card.suit,
                code: card.code,
                selected: selectedCard.value && selectedCard.value.id === card.id,
                disabledFeedback: disabledCardFeedback.value && disabledCardFeedback.value.id === card.id,
                cardIndex: index,
                flipped: false,
                disabled: isDisabled // Only Ace of Spades is enabled
              }
            })
          }
          
          // Determine legal moves based on In-Game cards
          let requiredSuit = null
          let canPlayAnyCard = true
          
          if (displayedGameState.value.inPlayCards && displayedGameState.value.inPlayCards.length > 0) {
            // There are In-Game cards, so we must follow suit if possible
            requiredSuit = displayedGameState.value.inPlayCards[0].card.suit
            canPlayAnyCard = false
            
            // Check if player has any cards of the required suit
            const hasRequiredSuit = availableCards.some(card => card.suit === requiredSuit)
            
            // If player has cards of required suit, they must play one of them
            // If player has no cards of required suit, they can play any card
            if (!hasRequiredSuit) {
              canPlayAnyCard = true
            } else {
            }
          } else {
            
          }
          
          return availableCards.map((card, index) => {
            const isLegalMove = canPlayAnyCard || (requiredSuit && card.suit === requiredSuit)
            const isDisabled = !isLegalMove // Backend handles timing - no frontend pausing needed
            
            
            
            return {
              id: card.id, // Include the card ID for API calls
              rank: card.rank,
              suit: card.suit,
              code: card.code,
              selected: selectedCard.value && 
                       selectedCard.value.id === card.id,
              disabledFeedback: disabledCardFeedback.value && 
                               disabledCardFeedback.value.id === card.id,
              cardIndex: index,
              flipped: false,
              disabled: isDisabled // Disable cards during pause or illegal moves
            }
          })
        }
        return []
      } else {
        // Show card backs for other players
        return Array(player.cardsInHand || 0).fill(null).map((_, index) => ({
          rank: '2',
          suit: 'hearts',
          flipped: true,
          selected: false,
          disabled: true
        }))
      }
    }
    
    const handleCardClick = async (data) => {
      const { card, event } = data
      
      // Prevent default GameCard auto-selection behavior
      if (event) {
        event.preventDefault()
      }
      
      
      
      // Don't allow any interactions during animation pause
      if (animationSystem.isCommunicationBlocked()) {
        return
      }
      
      // Don't process clicks on disabled cards, show red border feedback
      if (card.disabled) {
        // Stop event propagation to prevent any side effects
        if (event) {
          event.stopPropagation()
          event.stopImmediatePropagation()
        }
        showDisabledCardFeedback(card)
        return
      }
      
      // Check if this card is already selected
      const isCurrentlySelected = selectedCard.value && 
                                 selectedCard.value.id === (card.id || `${card.rank}${card.suit}`)
      
      if (isCurrentlySelected) {
        // Second click on selected card - play the card with animation
        await playCardWithAnimation(card)
      } else {
        // First click - select the card
        selectCard(card)
      }
    }
    
    const handleCardSelect = (data) => {
      // Ignore GameCard's auto-selection - we handle selection manually in handleCardClick
      // This prevents the selection toggle conflict that was causing cards to deselect
      return
    }
    
    const selectCard = (card) => {
      // Don't allow selection of disabled cards
      if (card.disabled) {
        return
      }
      
      selectedCard.value = {
        id: card.id,
        rank: card.rank,
        suit: card.suit
      }
      
    }
    
    const showDisabledCardFeedback = (card) => {
      
      // Ensure no selection happens for disabled cards
      if (selectedCard.value && selectedCard.value.id === card.id) {
        selectedCard.value = null
      }
      
      // Show red border feedback for disabled card
      disabledCardFeedback.value = {
        id: card.id,
        rank: card.rank,
        suit: card.suit
      }
      
      // Clear feedback after 3 seconds
      setTimeout(() => {
        disabledCardFeedback.value = null
      }, 3000)
      
    }
    
    const playCardWithAnimation = async (card) => {
      if (!gameState.value || gameState.value.game.status !== 'active') return
      if (!user.value.id) return
      if (!card.id) return
      
      try {
        
        // 1. ANIMATE CARD TO IN-GAME PILE POSITION (keep selection during animation)
        const cardElement = document.querySelector(`[data-card-id="${card.id}"]`)
        
        if (cardElement) {
          // Get current card position in fan (not center position)
          const cardRect = cardElement.getBoundingClientRect()
          
          // Get In-Game target position (top-right of current player area)
          const targetX = window.innerWidth - (8 * 16) // 8rem from right edge (right-8 class)
          const targetY = (window.innerHeight * 0.6) + (4 * 16) // 60vh + 4rem (top-4 of current player area)
          
          // Calculate movement distance
          const deltaX = targetX - (cardRect.left + cardRect.width * 0.5)
          const deltaY = targetY - (cardRect.top + cardRect.height * 0.5)
          
          // Apply animation transform
          cardElement.classList.add('card-moving-to-game')
          cardElement.style.transform = `translate(${deltaX}px, ${deltaY}px) scale(0.8)`
          cardElement.style.zIndex = '1000'
          
          // Clear selection after animation starts to prevent visual confusion
          selectedCard.value = null
          
          // Start API call immediately (don't wait for animation)
          const apiPromise = apiPlayCard(gameId.value, user.value.id, card.id)
          
          // Wait for animation to complete (1.5 seconds)
          await Promise.all([
            new Promise(resolve => setTimeout(resolve, 1500)),
            apiPromise
          ])
        } else {
          // No animation element found, just call API
          selectedCard.value = null
          await apiPlayCard(gameId.value, user.value.id, card.id)
        }
        
        
      } catch (error) {
        
        // Clear any stuck states
        selectedCard.value = null
      }
    }
    
    // Legacy function for compatibility
    const playCard = playCardWithAnimation
    
    const deselectCard = () => {
      selectedCard.value = null
    }
    
    // Helper methods for new grid layout
    const getOpponentInGameCard = (playerId) => {
      if (!displayedGameState.value || !displayedGameState.value.inPlayCards) return null
      return displayedGameState.value.inPlayCards.find(card => card.playerId === playerId)
    }
    
    const getCurrentPlayerInGameCard = () => {
      if (!displayedGameState.value || !displayedGameState.value.inPlayCards || !user.value.id) return null
      return displayedGameState.value.inPlayCards.find(card => card.playerId === user.value.id)
    }
    
    const getInGameCardPosition = (index, totalCards, playedCard) => {
      // Find which player played this card
      const playerId = playedCard.playerId
      
      // Find the player's seat information to get their angle
      const players = playersWithSeats.value || []
      const player = players.find(p => p.id === playerId)
      
      if (!player || !player.seatAngle) {
        // Fallback to simple positioning if player not found
        return {
          left: `${index * 8}px`,
          top: `${index * 2}px`,
          zIndex: index + 10
        }
      }
      
      // Oval specifications (same as in getPlayerSeats)
      const radiusX = 25 // 25% of width (horizontal radius)
      const radiusY = 35 // 35% of height (vertical radius)
      
      // Position card 25% of the distance from center toward player
      const distanceRatio = 0.25 // 25% of distance from center
      
      // Convert to pixel offsets (assuming typical viewport size)
      const viewportWidth = typeof window !== 'undefined' ? window.innerWidth : 1200
      const viewportHeight = typeof window !== 'undefined' ? window.innerHeight : 800
      
      const cardOffsetX = ((radiusX / 100) * viewportWidth * Math.cos(player.seatAngle)) * distanceRatio
      const cardOffsetY = ((radiusY / 100) * viewportHeight * Math.sin(player.seatAngle)) * distanceRatio
      
      return {
        left: `${cardOffsetX}px`,
        top: `${cardOffsetY}px`,
        zIndex: index + 10,
        transform: 'translate(-50%, -50%)'
      }
    }
    
    const handleMainAreaClick = (event) => {
      // Check if click was not on a card or card-related element
      if (!event.target.closest('.game-card') && !event.target.closest('.card-fan')) {
        deselectCard()
      }
    }
    
    // Turn result notification functions
    const showTurnResultNotification = (message, type = 'info') => {
      // Clear existing timers
      if (turnResultTimer) clearTimeout(turnResultTimer)
      if (turnResultProgressTimer) clearInterval(turnResultProgressTimer)
      
      // Set notification  
      turnResultNotification.value = { message, type }
      turnResultProgress.value = 100
      
      // Start progress countdown
      let progress = 100
      turnResultProgressTimer = setInterval(() => {
        progress -= 3.33 // Decrease by 3.33% every 100ms for 3 second duration
        turnResultProgress.value = progress
        if (progress <= 0) {
          clearInterval(turnResultProgressTimer)
        }
      }, 100)
      
      // Auto dismiss after 3 seconds
      turnResultTimer = setTimeout(() => {
        dismissTurnResultNotification()
      }, 3000)
    }
    
    const dismissTurnResultNotification = () => {
      if (turnResultTimer) clearTimeout(turnResultTimer)
      if (turnResultProgressTimer) clearInterval(turnResultProgressTimer)
      turnResultNotification.value = null
      turnResultProgress.value = 100
    }
    
    // Game end popup functions
    const showGameEndPopup = (gameEndData) => {
      // Clear existing timers
      if (gameEndTimer) clearTimeout(gameEndTimer)
      if (gameEndProgressTimer) clearInterval(gameEndProgressTimer)
      
      // Set popup data
      gameEndPopup.value = gameEndData
      gameEndProgress.value = 100
      
      // Start progress countdown (10 seconds)
      let progress = 100
      gameEndProgressTimer = setInterval(() => {
        progress -= 1 // Decrease by 1% every 100ms for 10 second duration
        gameEndProgress.value = progress
        if (progress <= 0) {
          clearInterval(gameEndProgressTimer)
        }
      }, 100)
      
      // Auto dismiss after 10 seconds
      gameEndTimer = setTimeout(() => {
        dismissGameEndPopup()
      }, 10000)
    }
    
    const dismissGameEndPopup = () => {
      if (gameEndTimer) clearTimeout(gameEndTimer)
      if (gameEndProgressTimer) clearInterval(gameEndProgressTimer)
      gameEndPopup.value = null
      gameEndProgress.value = 100
    }
    
    const copyGameLink = async () => {
      try {
        await navigator.clipboard.writeText(gameShareUrl.value)
        // Could show a toast notification here
      } catch (error) {
        
      }
    }
    
    const toggleSpeedUp = () => {
      speedUpEnabled.value = !speedUpEnabled.value
      // TODO: Send speed up preference to backend via API call
      // For now, just toggle the local state
    }
    
    const showAboutModal = () => {
      showAbout.value = true
    }
    
    const showHelpModal = () => {
      showHelp.value = true
    }
    
    const closeAbandoned = () => {
      showAbandoned.value = false
      resetGame()
    }
    
    const resetGame = () => {
      gameId.value = null
      gameState.value = null
      logs.value = []
      window.history.replaceState(null, '', '/')
    }
    
    // Lifecycle
    onMounted(async () => {
      // Check for existing user
      const userId = getCookie('userId')
      if (userId) {
        try {
          const userData = await getUser(userId)
          user.value = userData
          setCookie('userId', userData.id)
          setCookie('userName', userData.name)
        } catch (error) {
          // Clear invalid user data and show registration
          document.cookie = 'userId=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT'
          document.cookie = 'userName=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT'
          showRegistration.value = true
        }
      } else {
        showRegistration.value = true
      }
      
      // Check for game ID in URL
      const urlParams = new URLSearchParams(window.location.search)
      const gameIdParam = urlParams.get('gameId')
      if (gameIdParam && user.value.id) {
        await joinGameById(gameIdParam)
        window.history.replaceState(null, '', window.location.pathname)
      }
    })
    
    // Watch for game changes to setup connections
    watch([gameId, () => user.value.id], () => {
      if (gameId.value && user.value.id) {
        fetchGameState()
        fetchLogs()
        setupEventSource()
      }
    })
    
    onUnmounted(() => {
      if (eventSource.value) {
        eventSource.value.close()
      }
    })
    
    // Close dropdown when clicking outside
    const closeDropdownOnClickOutside = (event) => {
      if (!event.target.closest('.relative')) {
        dropdownOpen.value = null
      }
    }
    
    onMounted(() => {
      document.addEventListener('click', closeDropdownOnClickOutside)
    })
    
    onUnmounted(() => {
      document.removeEventListener('click', closeDropdownOnClickOutside)
    })
    
    return {
      // Theme
      currentTheme,
      themeLabel,
      cycleTheme,
      selectedTheme,
      
      // State
      user,
      gameId,
      gameState,
      displayedGameState,
      logs,
      connected,
      loading,
      
      // UI
      showMenu,
      showRegistration,
      showShare,
      showSettings,
      showAbandoned,
      showLog,
      dropdownOpen,
      speedUpEnabled,
      gameEndPopup,
      gameEndProgress,
      
      // Form data
      playerName,
      chatMessage,
      cardBackColor,
      cardBackColors,
      
      // Computed
      isRequester,
      playersWithSeats,
      opponentPlayers,
      currentPlayerData,
      botPlayers,
      gameShareUrl,
      menuItems,
      isCurrentTurn,
      isPlayerFinished,
      hasAceOfSpades,
      allRemainingPlayersAreBots,
      maxBotCap,
      
      // Methods
      toggleMenu,
      handleMenuAction,
      startNewGame,
      register,
      handleFinalizeGame,
      abandon,
      sendChat,
      getPlayerCards,
      handleCardClick,
      handleCardSelect,
      selectCard,
      showDisabledCardFeedback,
      playCard,
      deselectCard,
      disabledCardFeedback,
      handleMainAreaClick,
      selectedCard,
      getInGameCardPosition,
      gameAnimationState,
      copyGameLink,
      addBotPlayer,
      removeBotPlayer,
      changeBotDifficulty,
      toggleBotDifficultyDropdown,
      toggleSpeedUp,
      showAbout,
      showHelp,
      
      // Turn result notifications
      turnResultNotification,
      turnResultProgress,
      showTurnResultNotification,
      dismissTurnResultNotification,
      
      // Game end popup
      showGameEndPopup,
      dismissGameEndPopup,
      
      showAboutModal,
      showHelpModal,
      closeAbandoned,
      
      // Grid layout helpers
      getOpponentInGameCard,
      getCurrentPlayerInGameCard,
      
      // Card highlight helpers  
      getHighestInGameCard,
      isHighestInGameCard
    }
  }
}
</script>

<style scoped>
/* Ensure all In-Game cards have consistent sizing at 70% of previous size */
.in-game-card {
  /* Force consistent size for all In-Game cards regardless of source player */
  width: 4.2rem !important;   /* 70% of previous 6rem */
  height: 5.6rem !important;  /* 70% of previous 8rem */
  min-width: 4.2rem !important;
  min-height: 5.6rem !important;
  max-width: 4.2rem !important;
  max-height: 5.6rem !important;
  transform-origin: center !important;
}

/* Ensure no scale inheritance from larger user cards */
.in-game-card * {
  transform: none !important;
  scale: none !important;
}

/* Card movement animation */
.card-moving-to-game {
  z-index: 1000 !important;
  transition: all 1.5s cubic-bezier(0.4, 0, 0.2, 1) !important;
  pointer-events: none !important;
  position: relative !important;
}

/* Visual feedback during animation */
.card-moving-to-game::after {
  content: '';
  position: absolute;
  top: -4px;
  left: -4px;
  right: -4px;
  bottom: -4px;
  border: 3px solid #10b981;
  border-radius: 12px;
  box-shadow: 0 0 25px rgba(16, 185, 129, 0.8);
  animation: card-play-glow 1.5s ease-in-out;
  pointer-events: none;
  z-index: 1;
}

@keyframes card-play-glow {
  0% {
    opacity: 1;
    transform: scale(1);
    box-shadow: 0 0 25px rgba(16, 185, 129, 0.8);
  }
  50% {
    opacity: 0.9;
    transform: scale(1.1);
    box-shadow: 0 0 35px rgba(16, 185, 129, 1);
  }
  100% {
    opacity: 0;
    transform: scale(1);
    box-shadow: 0 0 15px rgba(16, 185, 129, 0.4);
  }
}

/* New Grid Layout Styles */
.game-viewport {
  display: grid;
  grid-template-rows: 60vh 40vh;
  grid-template-columns: 1fr;
}

.opponent-area {
  grid-row: 1;
  grid-column: 1;
}


.current-player-area {
  grid-row: 2;
  grid-column: 1 / -1;
}

.opponent-cell {
  display: grid;
  grid-template-columns: 75% 25%;
  grid-template-rows: 80% 20%;
}

.opponent-cards-area {
  grid-column: 1;
  grid-row: 1;
}

.opponent-ingame-area {
  grid-column: 2;
  grid-row: 1;
}

.opponent-name-area {
  grid-column: 1 / -1;
  grid-row: 2;
}

/* Responsive Grid Layout */
@media (max-width: 768px) {
  .game-viewport {
    grid-template-columns: 1fr;
  }
  
  .opponent-cell {
    grid-template-columns: 70% 30%;
  }
}

@media (max-width: 480px) {
  .game-viewport {
    grid-template-rows: 55vh 45vh;
    grid-template-columns: 1fr;
  }
  
  .opponent-cell {
    grid-template-columns: 65% 35%;
    grid-template-rows: 75% 25%;
  }
}

/* Communication blocked overlay */
.game-paused-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  z-index: 999;
  display: flex;
  align-items: center;
  justify-content: center;
}

.game-paused-message {
  background: rgba(59, 130, 246, 0.95);
  color: white;
  padding: 1rem 2rem;
  border-radius: 0.5rem;
  font-weight: 600;
  font-size: 1.25rem;
  text-align: center;
  animation: pulse-message 2s infinite ease-in-out;
}

@keyframes pulse-message {
  0%, 100% { opacity: 0.95; transform: scale(1); }
  50% { opacity: 1; transform: scale(1.05); }
}

/* Horizontal bounce animation for turn indicator */
.animate-horizontal-bounce {
  animation: horizontal-bounce 1s infinite;
}

@keyframes horizontal-bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateX(0) translateY(-50%);
  }
  40% {
    transform: translateX(4px) translateY(-50%);
  }
  60% {
    transform: translateX(2px) translateY(-50%);
  }
}
</style>
