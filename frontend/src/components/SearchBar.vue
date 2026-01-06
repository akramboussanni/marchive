<template>
  <div class="search-bar" :class="{ compact }">
    <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <circle cx="11" cy="11" r="8"></circle>
      <path d="m21 21-4.35-4.35"></path>
    </svg>
    <input 
      v-model="localQuery" 
      @keyup.enter="handleSearch"
      @input="handleInput"
      type="text" 
      :placeholder="placeholder"
      class="search-input"
    />
    <button 
      v-if="localQuery && showClear"
      @click="handleClear"
      class="clear-btn"
      title="Clear search"
    >
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <line x1="18" y1="6" x2="6" y2="18"></line>
        <line x1="6" y1="6" x2="18" y2="18"></line>
      </svg>
    </button>
    <button 
      v-if="showButton"
      @click="handleSearch" 
      :disabled="!localQuery || searching"
      class="search-btn"
    >
      <span v-if="!searching">Search</span>
      <div v-else class="spinner-small"></div>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

interface Props {
  modelValue?: string
  placeholder?: string
  searching?: boolean
  compact?: boolean
  showButton?: boolean
  showClear?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: 'Search for books...',
  searching: false,
  compact: false,
  showButton: true,
  showClear: true
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  search: [query: string]
  clear: []
}>()

const localQuery = ref(props.modelValue)

watch(() => props.modelValue, (newVal) => {
  localQuery.value = newVal
})

const handleInput = () => {
  emit('update:modelValue', localQuery.value)
}

const handleSearch = () => {
  if (!localQuery.value.trim() || props.searching) return
  emit('search', localQuery.value.trim())
}

const handleClear = () => {
  localQuery.value = ''
  emit('update:modelValue', '')
  emit('clear')
}
</script>

<style scoped>
.search-bar {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 12px;
  padding: 0.75rem 1rem;
  transition: all 0.3s ease;
}

.search-bar:focus-within {
  border-color: rgba(59, 130, 246, 0.6);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.search-bar.compact {
  padding: 0.625rem 0.875rem;
}

.search-icon {
  width: 20px;
  height: 20px;
  color: #64748b;
  flex-shrink: 0;
}

.search-bar.compact .search-icon {
  width: 18px;
  height: 18px;
}

.search-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: #e2e8f0;
  font-size: 1rem;
  padding: 0.25rem;
}

.search-bar.compact .search-input {
  font-size: 0.95rem;
}

.search-input::placeholder {
  color: #64748b;
}

.clear-btn {
  background: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  padding: 0.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.2s ease;
  flex-shrink: 0;
}

.clear-btn:hover {
  color: #94a3b8;
}

.clear-btn svg {
  width: 18px;
  height: 18px;
}

.search-btn {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  border: none;
  color: white;
  padding: 0.5rem 1.5rem;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 80px;
  flex-shrink: 0;
}

.search-bar.compact .search-btn {
  padding: 0.4rem 1.25rem;
  min-width: 70px;
  font-size: 0.9rem;
}

.search-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.search-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.spinner-small {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 768px) {
  .search-bar {
    padding: 0.625rem 0.875rem;
  }

  .search-icon {
    width: 18px;
    height: 18px;
  }

  .search-input {
    font-size: 0.9rem;
  }

  .clear-btn svg {
    width: 16px;
    height: 16px;
  }

  .search-btn {
    padding: 0.4rem 1.25rem;
    font-size: 0.9rem;
    min-width: 70px;
  }
}

@media (max-width: 480px) {
  .search-bar {
    padding: 0.5rem 0.75rem;
  }

  .search-icon {
    width: 16px;
    height: 16px;
  }

  .search-input {
    font-size: 0.875rem;
  }

  .search-btn {
    padding: 0.35rem 1rem;
    font-size: 0.85rem;
    min-width: 60px;
  }
}
</style>
