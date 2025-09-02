<template>
  <component
    :is="tag"
    :class="buttonClasses"
    :disabled="disabled || loading"
    :type="tag === 'button' ? type : undefined"
    :href="tag === 'a' ? href : undefined"
    :to="tag === 'router-link' ? to : undefined"
    v-bind="$attrs"
    @click="handleClick"
  >
    <!-- Loading spinner -->
    <div
      v-if="loading"
      class="animate-spin mr-2 h-4 w-4 border-2 border-current border-t-transparent rounded-full"
      aria-hidden="true"
    />
    
    <!-- Icon (before text) -->
    <slot
      v-if="$slots.icon && !loading"
      name="icon"
      class="mr-2"
    />
    
    <!-- Button text/content -->
    <slot />
    
    <!-- Icon (after text) -->
    <slot
      v-if="$slots.iconAfter"
      name="iconAfter"
      class="ml-2"
    />
  </component>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'BaseButton',
  inheritAttrs: false,
  emits: ['click'],
  props: {
    variant: {
      type: String,
      default: 'primary',
      validator: (value) => [
        'primary',
        'secondary', 
        'danger',
        'ghost',
        'outline',
        'success',
        'warning'
      ].includes(value)
    },
    size: {
      type: String,
      default: 'md',
      validator: (value) => ['xs', 'sm', 'md', 'lg', 'xl'].includes(value)
    },
    tag: {
      type: String,
      default: 'button',
      validator: (value) => ['button', 'a', 'router-link'].includes(value)
    },
    type: {
      type: String,
      default: 'button',
      validator: (value) => ['button', 'submit', 'reset'].includes(value)
    },
    href: {
      type: String,
      default: null
    },
    to: {
      type: [String, Object],
      default: null
    },
    disabled: {
      type: Boolean,
      default: false
    },
    loading: {
      type: Boolean,
      default: false
    },
    block: {
      type: Boolean,
      default: false
    },
    iconOnly: {
      type: Boolean,
      default: false
    },
    rounded: {
      type: String,
      default: 'md',
      validator: (value) => ['none', 'sm', 'md', 'lg', 'full'].includes(value)
    }
  },
  
  setup(props, { emit }) {
    // Size classes mapping
    const sizeClasses = {
      xs: 'px-2 py-1 text-xs',
      sm: 'px-3 py-1.5 text-sm',
      md: 'px-4 py-2 text-sm',
      lg: 'px-6 py-3 text-base',
      xl: 'px-8 py-4 text-lg'
    }
    
    // Icon-only size classes
    const iconOnlySizeClasses = {
      xs: 'p-1',
      sm: 'p-1.5',
      md: 'p-2',
      lg: 'p-3',
      xl: 'p-4'
    }
    
    // Variant classes mapping
    const variantClasses = {
      primary: 'btn-primary',
      secondary: 'btn-secondary',
      danger: 'btn-danger',
      ghost: 'btn-ghost',
      success: 'btn-success',
      warning: 'btn-warning',
      outline: 'bg-transparent border-2 border-[var(--color-primary)] text-[var(--color-primary)] hover:bg-[var(--color-primary)] hover:text-white focus:ring-[var(--color-primary)]'
    }
    
    // Rounded classes mapping
    const roundedClasses = {
      none: 'rounded-none',
      sm: 'rounded-sm',
      md: 'rounded-lg',
      lg: 'rounded-xl',
      full: 'rounded-full'
    }
    
    // Computed button classes
    const buttonClasses = computed(() => {
      const classes = [
        'btn-base',
        variantClasses[props.variant],
        props.iconOnly ? iconOnlySizeClasses[props.size] : sizeClasses[props.size],
        roundedClasses[props.rounded]
      ]
      
      if (props.block) {
        classes.push('w-full')
      }
      
      if (props.iconOnly) {
        classes.push('flex items-center justify-center')
      }
      
      if (props.loading) {
        classes.push('cursor-wait')
      }
      
      return classes.join(' ')
    })
    
    // Click handler
    const handleClick = (event) => {
      if (!props.disabled && !props.loading) {
        emit('click', event)
      }
    }
    
    return {
      buttonClasses,
      handleClick
    }
  }
}
</script>

<style scoped>
/* Additional component-specific styles if needed */
.btn-base {
  /* Ensure smooth transitions for all interactive states */
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.btn-base:active {
  transform: translateY(1px);
}

.btn-base:disabled {
  transform: none;
}

/* Focus styles for better accessibility */
.btn-base:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}

/* Mobile touch optimization */
@media (max-width: 768px) {
  .btn-base {
    min-height: 44px; /* Minimum touch target size */
  }
}
</style>