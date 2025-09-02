<template>
  <div class="session-updates-container">
    <!-- Toggle Button -->
    <button
      @click="togglePanel"
      class="session-updates-toggle"
      :class="{ 'panel-open': isPanelOpen }"
      :title="isPanelOpen ? 'Hide Session Updates' : 'Show Session Updates'"
    >
      <div class="toggle-icon">
        <span v-if="!isPanelOpen">💬</span>
        <span v-else>✕</span>
      </div>
      <div class="unread-badge" v-if="unreadCount > 0 && !isPanelOpen">
        {{ unreadCount > 99 ? '99+' : unreadCount }}
      </div>
    </button>
    
    <!-- Session Updates Panel -->
    <Transition name="session-panel">
      <div 
        v-if="isPanelOpen"
        class="session-updates-panel"
        @click.stop
      >
        <!-- Panel Header -->
        <div class="panel-header">
          <div class="panel-title">
            <h3>Session Updates</h3>
            <span class="message-count">({{ allMessages.length }})</span>
          </div>
          
          <!-- Tab Switcher -->
          <div class="panel-tabs">
            <button
              @click="currentTab = 'all'"
              :class="{ active: currentTab === 'all' }"
              class="tab-button"
            >
              All
              <span class="tab-count">({{ allMessages.length }})</span>
            </button>
            <button
              @click="currentTab = 'events'"
              :class="{ active: currentTab === 'events' }"
              class="tab-button"
            >
              Events
              <span class="tab-count">({{ gameEvents.length }})</span>
            </button>
            <button
              @click="currentTab = 'chat'"
              :class="{ active: currentTab === 'chat' }"
              class="tab-button"
            >
              Chat
              <span class="tab-count">({{ chatMessages.length }})</span>
            </button>
          </div>
        </div>
        
        <!-- Messages List -->
        <div class="messages-container" ref="messagesContainer">
          <div 
            v-for="message in displayedMessages"
            :key="message.id"
            :class="getMessageClasses(message)"
            class="message-item"
          >
            <!-- Message Header -->
            <div class="message-header">
              <div class="message-meta">
                <span class="message-type-icon">
                  {{ getMessageIcon(message) }}
                </span>
                <span v-if="!message.isSystemMessage" class="message-author">
                  {{ message.authorName || 'Unknown' }}
                </span>
                <span class="message-time">
                  {{ message.displayTime }}
                </span>
              </div>
            </div>
            
            <!-- Message Content -->
            <div class="message-content">
              <span v-html="formatMessageContent(message.message)"></span>
            </div>
            
            <!-- Message Actions (for chat messages) -->
            <div 
              v-if="message.type === 'chat' && message.userId === user?.id"
              class="message-actions"
            >
              <button 
                @click="editMessage(message)"
                class="action-button edit"
                title="Edit message"
                disabled
              >
                ✏️
              </button>
              <button 
                @click="deleteMessage(message)"
                class="action-button delete"
                title="Delete message"
                disabled
              >
                🗑️
              </button>
            </div>
          </div>
          
          <!-- Empty State -->
          <div 
            v-if="displayedMessages.length === 0"
            class="empty-state"
          >
            <div class="empty-icon">
              {{ currentTab === 'chat' ? '💬' : '📋' }}
            </div>
            <div class="empty-text">
              {{ getEmptyStateText() }}
            </div>
          </div>
        </div>
        
        <!-- Chat Input (only for chat tab) -->
        <div 
          v-if="currentTab === 'chat'"
          class="chat-input-container"
        >
          <div class="chat-input-wrapper">
            <input
              v-model="chatMessage"
              @keydown="handleChatKeydown"
              @input="handleChatInput"
              placeholder="Type a message..."
              class="chat-input"
              :disabled="!canSendChat"
              :maxlength="chatMaxLength"
            />
            <button
              @click="sendChatMessage"
              :disabled="!canSendMessage"
              class="chat-send-button"
              title="Send message"
            >
              ➤
            </button>
          </div>
          <div class="chat-input-meta">
            <span class="character-count">
              {{ chatMessage.length }}/{{ chatMaxLength }}
            </span>
            <span v-if="rateLimitMessage" class="rate-limit-warning">
              {{ rateLimitMessage }}
            </span>
          </div>
        </div>
      </div>
    </Transition>
    
    <!-- Overlay to close panel when clicking outside -->
    <div
      v-if="isPanelOpen"
      class="panel-overlay"
      @click="closePanel"
    ></div>
  </div>
</template>

<script>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'

export default {
  name: 'SessionUpdates',
  props: {
    messages: {
      type: Array,
      required: true,
      default: () => []
    },
    user: {
      type: Object,
      default: null
    },
    gameId: {
      type: String,
      required: true
    },
    isConnected: {
      type: Boolean,
      default: false
    },
    chatMaxLength: {
      type: Number,
      default: 280
    },
    rateLimitSeconds: {
      type: Number,
      default: 6 // 10 messages per minute = 6 seconds between messages
    }
  },
  emits: ['send-chat', 'mark-read'],
  setup(props, { emit }) {
    // Panel state
    const isPanelOpen = ref(false)
    const currentTab = ref('all')
    const messagesContainer = ref(null)
    
    // Chat input state
    const chatMessage = ref('')
    const lastChatTime = ref(0)
    const rateLimitMessage = ref('')
    
    // Unread tracking
    const unreadCount = ref(0)
    const lastReadTime = ref(Date.now())
    
    // Message filtering
    const allMessages = computed(() => {
      return props.messages.sort((a, b) => 
        new Date(a.createdAt) - new Date(b.createdAt)
      )
    })
    
    const gameEvents = computed(() => {
      return allMessages.value.filter(msg => msg.isSystemMessage)
    })
    
    const chatMessages = computed(() => {
      return allMessages.value.filter(msg => msg.type === 'chat')
    })
    
    const displayedMessages = computed(() => {
      switch (currentTab.value) {
        case 'events':
          return gameEvents.value
        case 'chat':
          return chatMessages.value
        default:
          return allMessages.value
      }
    })
    
    // Chat capabilities
    const canSendChat = computed(() => {
      return props.isConnected && props.user
    })
    
    const canSendMessage = computed(() => {
      return canSendChat.value && 
             chatMessage.value.trim().length > 0 && 
             chatMessage.value.length <= props.chatMaxLength &&
             !rateLimitMessage.value
    })
    
    // Update unread count
    const updateUnreadCount = () => {
      if (isPanelOpen.value) {
        unreadCount.value = 0
        lastReadTime.value = Date.now()
        emit('mark-read')
      } else {
        const newMessages = allMessages.value.filter(msg => 
          new Date(msg.createdAt) > new Date(lastReadTime.value)
        )
        unreadCount.value = newMessages.length
      }
    }
    
    // Panel controls
    const togglePanel = () => {
      isPanelOpen.value = !isPanelOpen.value
      if (isPanelOpen.value) {
        nextTick(() => {
          scrollToBottom()
          updateUnreadCount()
        })
      }
    }
    
    const closePanel = () => {
      isPanelOpen.value = false
    }
    
    // Chat functionality
    const sendChatMessage = () => {
      if (!canSendMessage.value) return
      
      const now = Date.now()
      const timeSinceLastMessage = now - lastChatTime.value
      
      if (timeSinceLastMessage < props.rateLimitSeconds * 1000) {
        const waitTime = Math.ceil((props.rateLimitSeconds * 1000 - timeSinceLastMessage) / 1000)
        rateLimitMessage.value = `Please wait ${waitTime}s before sending another message`
        setTimeout(() => {
          rateLimitMessage.value = ''
        }, waitTime * 1000)
        return
      }
      
      emit('send-chat', {
        message: chatMessage.value.trim(),
        gameId: props.gameId,
        userId: props.user.id
      })
      
      chatMessage.value = ''
      lastChatTime.value = now
    }
    
    const handleChatKeydown = (event) => {
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault()
        sendChatMessage()
      }
    }
    
    const handleChatInput = () => {
      // Auto-resize or other input handling could go here
    }
    
    // Message display utilities
    const getMessageClasses = (message) => {
      const classes = []
      
      if (message.isSystemMessage) {
        classes.push('message-system')
      } else {
        classes.push('message-chat')
      }
      
      if (message.userId === props.user?.id) {
        classes.push('message-own')
      }
      
      // Add type-specific classes
      classes.push(`message-${message.type}`)
      
      return classes.join(' ')
    }
    
    const getMessageIcon = (message) => {
      const iconMap = {
        'game_event': '🎮',
        'round_event': '🃏',
        'turn_event': '🔄',
        'chat': '💬',
        'system': 'ℹ️',
        'pause_started': '⏸️',
        'pause_completed': '▶️',
        'cut_performed': '⚔️',
        'player_joined': '👋',
        'player_disconnected': '📴',
        'player_reconnected': '📶'
      }
      
      // Check for specific event types in message content
      const content = message.message.toLowerCase()
      if (content.includes('cut')) return '⚔️'
      if (content.includes('pause')) return '⏸️'
      if (content.includes('joined')) return '👋'
      if (content.includes('disconnected')) return '📴'
      if (content.includes('bot')) return '🤖'
      if (content.includes('round')) return '🃏'
      if (content.includes('turn')) return '🔄'
      
      return iconMap[message.type] || '📋'
    }

    // Replace card codes like AS, 10D with emoji suits
    const formatMessageContent = (text) => {
      if (!text) return ''
      const suitMap = { 'S': '♠️', 'H': '♥️', 'D': '♦️', 'C': '♣️' }
      // Replace codes like AS, 10D, QH, 5C (word boundaries)
      return text.replace(/\b(10|[2-9]|[JQKA])([SHDC])\b/g, (_, rank, suit) => `${rank}${suitMap[suit]}`)
    }
    
    const formatMessageContent = (content) => {
      // Basic message formatting
      return content
        .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>') // **bold**
        .replace(/\*(.*?)\*/g, '<em>$1</em>') // *italic*
        .replace(/`(.*?)`/g, '<code>$1</code>') // `code`
    }
    
    const getEmptyStateText = () => {
      switch (currentTab.value) {
        case 'events':
          return 'No game events yet'
        case 'chat':
          return 'No chat messages yet. Start a conversation!'
        default:
          return 'No messages yet'
      }
    }
    
    // Scroll management
    const scrollToBottom = () => {
      if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
      }
    }
    
    // Watch for new messages and auto-scroll
    watch(() => displayedMessages.value.length, () => {
      if (isPanelOpen.value) {
        nextTick(scrollToBottom)
      }
      updateUnreadCount()
    })
    
    // Watch panel open/close for unread updates
    watch(isPanelOpen, (open) => {
      if (open) {
        updateUnreadCount()
      }
    })
    
    // Keyboard shortcut
    const handleKeydown = (event) => {
      // Toggle with Ctrl/Cmd + M
      if ((event.ctrlKey || event.metaKey) && event.key === 'm') {
        event.preventDefault()
        togglePanel()
      }
      
      // Close with Escape
      if (event.key === 'Escape' && isPanelOpen.value) {
        closePanel()
      }
    }
    
    // Placeholder functions for future features
    const editMessage = (message) => {
      // Future implementation
    }
    
    const deleteMessage = (message) => {
      // Future implementation
    }
    
    onMounted(() => {
      document.addEventListener('keydown', handleKeydown)
      updateUnreadCount()
    })
    
    onUnmounted(() => {
      document.removeEventListener('keydown', handleKeydown)
    })
    
    return {
      // State
      isPanelOpen,
      currentTab,
      messagesContainer,
      chatMessage,
      unreadCount,
      rateLimitMessage,
      
      // Computed
      allMessages,
      gameEvents,
      chatMessages,
      displayedMessages,
      canSendChat,
      canSendMessage,
      
      // Methods
      togglePanel,
      closePanel,
      sendChatMessage,
      handleChatKeydown,
      handleChatInput,
      getMessageClasses,
      getMessageIcon,
      formatMessageContent,
      getEmptyStateText,
      editMessage,
      deleteMessage
    }
  }
}
</script>

<style scoped>
.session-updates-container {
  position: fixed;
  bottom: 1rem;
  right: 1rem;
  z-index: 30;
}

/* Toggle Button */
.session-updates-toggle {
  position: relative;
  width: 3rem;
  height: 3rem;
  border-radius: 50%;
  background: var(--color-primary);
  color: white;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.session-updates-toggle:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 16px rgba(59, 130, 246, 0.5);
}

.session-updates-toggle.panel-open {
  background: var(--color-surface);
  color: var(--color-text);
  border: 2px solid var(--color-primary);
}

.toggle-icon {
  font-size: 1.25rem;
  transition: transform 0.3s ease;
}

.panel-open .toggle-icon {
  transform: rotate(90deg);
}

.unread-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  background: #ef4444;
  color: white;
  font-size: 0.625rem;
  font-weight: 700;
  padding: 0.125rem 0.25rem;
  border-radius: 0.5rem;
  min-width: 1rem;
  text-align: center;
  animation: unread-pulse 2s infinite;
}

@keyframes unread-pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); }
}

/* Panel */
.session-updates-panel {
  position: absolute;
  bottom: 4rem;
  right: 0;
  width: 24rem;
  max-height: 32rem;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 0.75rem;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.panel-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: -1;
  background: transparent;
}

/* Panel Header */
.panel-header {
  padding: 1rem 1rem 0.5rem;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-background-secondary);
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.panel-title h3 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.message-count {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  font-weight: 400;
}

/* Tabs */
.panel-tabs {
  display: flex;
  gap: 0.25rem;
}

.tab-button {
  flex: 1;
  padding: 0.375rem 0.75rem;
  border: 1px solid var(--color-border);
  background: var(--color-surface);
  color: var(--color-text);
  border-radius: 0.375rem;
  font-size: 0.75rem;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.25rem;
}

.tab-button:hover {
  background: var(--color-background);
}

.tab-button.active {
  background: var(--color-primary);
  color: white;
  border-color: var(--color-primary);
}

.tab-count {
  font-size: 0.625rem;
  opacity: 0.8;
}

/* Messages */
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 0.5rem 0;
  max-height: 20rem;
}

.message-item {
  padding: 0.5rem 1rem;
  border-bottom: 1px solid rgba(var(--color-border-rgb), 0.1);
  transition: background 0.2s ease;
}

.message-item:hover {
  background: var(--color-background-secondary);
}

.message-item:last-child {
  border-bottom: none;
}

.message-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.25rem;
}

.message-meta {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.message-type-icon {
  font-size: 0.875rem;
}

.message-author {
  font-weight: 500;
  color: var(--color-text);
}

.message-time {
  opacity: 0.8;
}

.message-content {
  font-size: 0.875rem;
  line-height: 1.4;
  color: var(--color-text);
  word-wrap: break-word;
}

/* Message type styling */
.message-system .message-content {
  font-style: italic;
  color: var(--color-text-secondary);
}

.message-own .message-content {
  color: var(--color-primary);
}

.message-actions {
  margin-top: 0.5rem;
  display: flex;
  gap: 0.25rem;
}

.action-button {
  background: none;
  border: none;
  cursor: pointer;
  opacity: 0.6;
  font-size: 0.75rem;
  padding: 0.125rem;
  border-radius: 0.25rem;
  transition: opacity 0.2s ease;
}

.action-button:hover:not(:disabled) {
  opacity: 1;
}

.action-button:disabled {
  cursor: not-allowed;
  opacity: 0.3;
}

/* Empty State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem 1rem;
  color: var(--color-text-secondary);
  text-align: center;
}

.empty-icon {
  font-size: 2rem;
  margin-bottom: 0.5rem;
  opacity: 0.5;
}

.empty-text {
  font-size: 0.875rem;
}

/* Chat Input */
.chat-input-container {
  border-top: 1px solid var(--color-border);
  padding: 0.75rem;
  background: var(--color-background-secondary);
}

.chat-input-wrapper {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.chat-input {
  flex: 1;
  padding: 0.5rem;
  border: 1px solid var(--color-border);
  border-radius: 0.375rem;
  background: var(--color-surface);
  color: var(--color-text);
  font-size: 0.875rem;
  resize: none;
  outline: none;
  transition: border-color 0.2s ease;
}

.chat-input:focus {
  border-color: var(--color-primary);
}

.chat-input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.chat-send-button {
  padding: 0.5rem 0.75rem;
  background: var(--color-primary);
  color: white;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chat-send-button:hover:not(:disabled) {
  background: var(--color-primary-dark);
  transform: translateY(-1px);
}

.chat-send-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.chat-input-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 0.25rem;
  font-size: 0.625rem;
  color: var(--color-text-secondary);
}

.rate-limit-warning {
  color: #ef4444;
  font-weight: 500;
}

/* Panel Transitions */
.session-panel-enter-active,
.session-panel-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.session-panel-enter-from,
.session-panel-leave-to {
  opacity: 0;
  transform: translateY(1rem) scale(0.95);
}

/* Mobile Optimizations */
@media (max-width: 768px) {
  .session-updates-container {
    bottom: 0.5rem;
    right: 0.5rem;
  }
  
  .session-updates-panel {
    width: calc(100vw - 2rem);
    max-width: 20rem;
    bottom: 3.5rem;
    right: -0.5rem;
  }
  
  .messages-container {
    max-height: 16rem;
  }
  
  .panel-title h3 {
    font-size: 0.875rem;
  }
  
  .tab-button {
    padding: 0.25rem 0.5rem;
    font-size: 0.625rem;
  }
  
  .message-item {
    padding: 0.375rem 0.75rem;
  }
  
  .message-content {
    font-size: 0.8125rem;
  }
}

/* Scrollbar Styling */
.messages-container::-webkit-scrollbar {
  width: 6px;
}

.messages-container::-webkit-scrollbar-track {
  background: var(--color-background);
  border-radius: 3px;
}

.messages-container::-webkit-scrollbar-thumb {
  background: var(--color-border);
  border-radius: 3px;
}

.messages-container::-webkit-scrollbar-thumb:hover {
  background: var(--color-text-secondary);
}

/* Reduced Motion Support */
@media (prefers-reduced-motion: reduce) {
  .session-updates-toggle {
    transition: none;
  }
  
  .toggle-icon {
    transition: none;
  }
  
  .unread-badge {
    animation: none;
  }
  
  .session-panel-enter-active,
  .session-panel-leave-active {
    transition: opacity 0.2s ease;
  }
  
  .session-panel-enter-from,
  .session-panel-leave-to {
    transform: none;
  }
}
</style>
