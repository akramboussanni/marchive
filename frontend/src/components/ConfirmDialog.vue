<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isOpen" class="modal-overlay" @click="handleCancel">
        <div class="confirm-dialog" @click.stop>
          <div class="dialog-header">
            <svg 
              v-if="variant === 'danger'" 
              class="dialog-icon danger"
              viewBox="0 0 24 24" 
              fill="none" 
              stroke="currentColor" 
              stroke-width="2"
            >
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="15" y1="9" x2="9" y2="15"></line>
              <line x1="9" y1="9" x2="15" y2="15"></line>
            </svg>
            <svg 
              v-else-if="variant === 'warning'" 
              class="dialog-icon warning"
              viewBox="0 0 24 24" 
              fill="none" 
              stroke="currentColor" 
              stroke-width="2"
            >
              <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
              <line x1="12" y1="9" x2="12" y2="13"></line>
              <line x1="12" y1="17" x2="12.01" y2="17"></line>
            </svg>
            <svg 
              v-else 
              class="dialog-icon info"
              viewBox="0 0 24 24" 
              fill="none" 
              stroke="currentColor" 
              stroke-width="2"
            >
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="16" x2="12" y2="12"></line>
              <line x1="12" y1="8" x2="12.01" y2="8"></line>
            </svg>
          </div>
          
          <div class="dialog-content">
            <h3 class="dialog-title">{{ title }}</h3>
            <p class="dialog-message">{{ message }}</p>
          </div>

          <div class="dialog-actions">
            <button 
              class="dialog-btn cancel-btn" 
              @click="handleCancel"
            >
              {{ cancelText }}
            </button>
            <button 
              class="dialog-btn confirm-btn"
              :class="{ danger: variant === 'danger' }"
              @click="handleConfirm"
            >
              {{ confirmText }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
interface Props {
  isOpen: boolean
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  variant?: 'info' | 'warning' | 'danger'
}

const props = withDefaults(defineProps<Props>(), {
  confirmText: 'Confirm',
  cancelText: 'Cancel',
  variant: 'info'
})

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

const handleConfirm = () => {
  emit('confirm')
}

const handleCancel = () => {
  emit('cancel')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.confirm-dialog {
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.98) 0%, rgba(30, 41, 59, 0.98) 100%);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 16px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  max-width: 480px;
  width: 100%;
  padding: 1.5rem;
  animation: slideUp 0.3s ease;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1rem;
}

.dialog-icon {
  width: 48px;
  height: 48px;
  padding: 0.75rem;
  border-radius: 50%;
}

.dialog-icon.danger {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.dialog-icon.warning {
  color: #f59e0b;
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.dialog-icon.info {
  color: #3b82f6;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.dialog-content {
  text-align: center;
  margin-bottom: 1.5rem;
}

.dialog-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #e2e8f0;
  margin-bottom: 0.5rem;
}

.dialog-message {
  font-size: 0.938rem;
  color: #94a3b8;
  line-height: 1.5;
}

.dialog-actions {
  display: flex;
  gap: 0.75rem;
}

.dialog-btn {
  flex: 1;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 500;
  font-size: 0.938rem;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid;
}

.cancel-btn {
  background: rgba(71, 85, 105, 0.2);
  border-color: rgba(71, 85, 105, 0.4);
  color: #cbd5e1;
}

.cancel-btn:hover {
  background: rgba(71, 85, 105, 0.3);
  border-color: rgba(71, 85, 105, 0.6);
  transform: translateY(-1px);
}

.confirm-btn {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.4);
  color: #3b82f6;
}

.confirm-btn:hover {
  background: rgba(59, 130, 246, 0.3);
  border-color: rgba(59, 130, 246, 0.6);
  transform: translateY(-1px);
}

.confirm-btn.danger {
  background: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.4);
  color: #ef4444;
}

.confirm-btn.danger:hover {
  background: rgba(239, 68, 68, 0.3);
  border-color: rgba(239, 68, 68, 0.6);
}

/* Transitions */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .confirm-dialog,
.modal-leave-active .confirm-dialog {
  transition: transform 0.3s ease;
}

.modal-enter-from .confirm-dialog,
.modal-leave-to .confirm-dialog {
  transform: scale(0.9) translateY(20px);
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: scale(0.9) translateY(20px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

@media (max-width: 480px) {
  .confirm-dialog {
    padding: 1.25rem;
  }

  .dialog-icon {
    width: 40px;
    height: 40px;
    padding: 0.625rem;
  }

  .dialog-title {
    font-size: 1.125rem;
  }

  .dialog-message {
    font-size: 0.875rem;
  }

  .dialog-actions {
    flex-direction: column;
  }

  .dialog-btn {
    padding: 0.625rem 1.25rem;
    font-size: 0.875rem;
  }
}
</style>
