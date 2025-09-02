<template>
  <Teleport to="body">
    <Transition
      name="modal"
      appear
    >
      <div
        v-if="isOpen"
        class="modal-backdrop"
        @click="handleBackdropClick"
        @keydown.esc="handleEscapeKey"
        tabindex="-1"
        role="dialog"
        :aria-labelledby="titleId"
        :aria-describedby="descriptionId"
        aria-modal="true"
      >
        <div
          :class="modalClasses"
          @click.stop
          ref="modalRef"
        >
          <!-- Header -->
          <header
            v-if="$slots.header || title || showClose"
            class="flex items-center justify-between p-4 sm:p-6 border-b border-[var(--color-border)]"
          >
            <div class="flex-1">
              <slot name="header">
                <h2
                  v-if="title"
                  :id="titleId"
                  class="text-responsive-lg font-semibold text-[var(--color-text)]"
                >
                  {{ title }}
                </h2>
              </slot>
            </div>
            
            <BaseButton
              v-if="showClose"
              variant="ghost"
              size="sm"
              icon-only
              rounded="full"
              @click="close"
              :aria-label="closeLabel"
              class="ml-4 flex-shrink-0"
            >
              <template #icon>
                <svg
                  class="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </template>
            </BaseButton>
          </header>
          
          <!-- Body -->
          <main
            class="p-4 sm:p-6"
            :id="descriptionId"
          >
            <slot />
          </main>
          
          <!-- Footer -->
          <footer
            v-if="$slots.footer"
            class="flex flex-col-reverse sm:flex-row sm:justify-end gap-3 p-4 sm:p-6 border-t border-[var(--color-border)]"
          >
            <slot name="footer" />
          </footer>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script>
import { computed, ref, watch, nextTick, onUnmounted } from 'vue'
import BaseButton from './BaseButton.vue'

let modalCounter = 0

export default {
  name: 'BaseModal',
  components: {
    BaseButton
  },
  emits: ['close', 'open'],
  props: {
    isOpen: {
      type: Boolean,
      default: false
    },
    title: {
      type: String,
      default: ''
    },
    size: {
      type: String,
      default: 'md',
      validator: (value) => ['xs', 'sm', 'md', 'lg', 'xl', 'full'].includes(value)
    },
    showClose: {
      type: Boolean,
      default: true
    },
    closeOnBackdrop: {
      type: Boolean,
      default: true
    },
    closeOnEscape: {
      type: Boolean,
      default: true
    },
    closeLabel: {
      type: String,
      default: 'Close modal'
    },
    trapFocus: {
      type: Boolean,
      default: true
    },
    restoreFocus: {
      type: Boolean,
      default: true
    }
  },
  
  setup(props, { emit }) {
    const modalRef = ref(null)
    const previousActiveElement = ref(null)
    const modalId = ++modalCounter
    
    // Unique IDs for accessibility
    const titleId = computed(() => `modal-title-${modalId}`)
    const descriptionId = computed(() => `modal-description-${modalId}`)
    
    // Size classes mapping
    const sizeClasses = {
      xs: 'max-w-xs',
      sm: 'max-w-sm',
      md: 'max-w-md',
      lg: 'max-w-lg',
      xl: 'max-w-xl',
      full: 'max-w-full mx-4'
    }
    
    // Computed modal classes
    const modalClasses = computed(() => {
      return [
        'modal-content',
        sizeClasses[props.size],
        'animate-bounce-in'
      ].join(' ')
    })
    
    // Focus management
    const focusableSelectors = [
      'button:not([disabled])',
      '[href]',
      'input:not([disabled])',
      'select:not([disabled])',
      'textarea:not([disabled])',
      '[tabindex]:not([tabindex="-1"])'
    ].join(', ')
    
    const getFocusableElements = () => {
      if (!modalRef.value) return []
      return Array.from(modalRef.value.querySelectorAll(focusableSelectors))
    }
    
    const focusFirstElement = async () => {
      await nextTick()
      const focusableElements = getFocusableElements()
      if (focusableElements.length > 0) {
        focusableElements[0].focus()
      }
    }
    
    const handleTabKey = (event) => {
      if (!props.trapFocus) return
      
      const focusableElements = getFocusableElements()
      if (focusableElements.length === 0) return
      
      const firstElement = focusableElements[0]
      const lastElement = focusableElements[focusableElements.length - 1]
      
      if (event.shiftKey) {
        // Shift + Tab
        if (document.activeElement === firstElement) {
          event.preventDefault()
          lastElement.focus()
        }
      } else {
        // Tab
        if (document.activeElement === lastElement) {
          event.preventDefault()
          firstElement.focus()
        }
      }
    }
    
    // Event handlers
    const close = () => {
      emit('close')
    }
    
    const open = () => {
      emit('open')
    }
    
    const handleBackdropClick = () => {
      if (props.closeOnBackdrop) {
        close()
      }
    }
    
    const handleEscapeKey = (event) => {
      if (event.key === 'Escape' && props.closeOnEscape) {
        close()
      }
    }
    
    const handleKeydown = (event) => {
      if (event.key === 'Tab') {
        handleTabKey(event)
      }
    }
    
    // Body scroll management
    const toggleBodyScroll = (disable) => {
      if (typeof document === 'undefined') return
      
      if (disable) {
        document.body.style.overflow = 'hidden'
        document.body.style.paddingRight = `${window.innerWidth - document.documentElement.clientWidth}px`
      } else {
        document.body.style.overflow = ''
        document.body.style.paddingRight = ''
      }
    }
    
    // Watch for modal open/close
    watch(
      () => props.isOpen,
      async (isOpen) => {
        if (isOpen) {
          // Store previously focused element
          if (props.restoreFocus) {
            previousActiveElement.value = document.activeElement
          }
          
          // Disable body scroll
          toggleBodyScroll(true)
          
          // Focus first element
          if (props.trapFocus) {
            await focusFirstElement()
          }
          
          // Add keydown listener
          document.addEventListener('keydown', handleKeydown)
          
          emit('open')
        } else {
          // Enable body scroll
          toggleBodyScroll(false)
          
          // Restore focus
          if (props.restoreFocus && previousActiveElement.value) {
            previousActiveElement.value.focus()
            previousActiveElement.value = null
          }
          
          // Remove keydown listener
          document.removeEventListener('keydown', handleKeydown)
        }
      },
      { immediate: true }
    )
    
    // Cleanup on unmount
    onUnmounted(() => {
      toggleBodyScroll(false)
      document.removeEventListener('keydown', handleKeydown)
    })
    
    return {
      modalRef,
      titleId,
      descriptionId,
      modalClasses,
      close,
      open,
      handleBackdropClick,
      handleEscapeKey
    }
  }
}
</script>

<style scoped>
/* Modal transition animations */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal-content,
.modal-leave-active .modal-content {
  transition: transform 0.3s ease;
}

.modal-enter-from .modal-content,
.modal-leave-to .modal-content {
  transform: scale(0.9) translateY(-10px);
}

/* Ensure modal is above other content */
.modal-backdrop {
  z-index: 9999;
}

/* Mobile optimizations */
@media (max-width: 640px) {
  .modal-content {
    margin: 1rem;
    max-height: calc(100vh - 2rem);
  }
}

/* Prevent background scroll on iOS */
.modal-backdrop {
  position: fixed;
  overflow: auto;
  overscroll-behavior: contain;
}
</style>