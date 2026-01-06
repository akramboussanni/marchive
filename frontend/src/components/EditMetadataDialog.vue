<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isOpen" class="modal-overlay" @click="handleCancel">
        <div class="edit-dialog" @click.stop>
          <div class="dialog-header">
            <h3>Edit Book Metadata</h3>
            <button @click="handleCancel" class="close-btn" title="Close">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>

          <div class="dialog-content">
            <div class="form-group">
              <label for="title" class="form-label">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path>
                  <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path>
                </svg>
                Title
              </label>
              <input 
                id="title"
                v-model="editedTitle"
                type="text"
                class="form-input"
                placeholder="Enter book title"
                @keydown.enter="handleSave"
              />
            </div>

            <div class="form-group">
              <label for="authors" class="form-label">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                  <circle cx="12" cy="7" r="4"></circle>
                </svg>
                Authors
              </label>
              <input 
                id="authors"
                v-model="editedAuthors"
                type="text"
                class="form-input"
                placeholder="Enter authors"
                @keydown.enter="handleSave"
              />
            </div>

            <div class="form-group">
              <label for="publisher" class="form-label">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
                  <polyline points="9 22 9 12 15 12 15 22"></polyline>
                </svg>
                Publisher
              </label>
              <input 
                id="publisher"
                v-model="editedPublisher"
                type="text"
                class="form-input"
                placeholder="Enter publisher"
                @keydown.enter="handleSave"
              />
            </div>

            <div v-if="error" class="error-message">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="15" y1="9" x2="9" y2="15"></line>
                <line x1="9" y1="9" x2="15" y2="15"></line>
              </svg>
              <span>{{ error }}</span>
            </div>
          </div>

          <div class="dialog-actions">
            <button 
              class="dialog-btn cancel-btn" 
              @click="handleCancel"
            >
              Cancel
            </button>
            <button 
              class="dialog-btn save-btn"
              @click="handleSave"
              :disabled="!isValid"
            >
              Save Changes
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

interface Props {
  isOpen: boolean
  title: string
  authors: string
  publisher: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  save: [data: { title: string; authors: string; publisher: string }]
  cancel: []
}>()

const editedTitle = ref('')
const editedAuthors = ref('')
const editedPublisher = ref('')
const error = ref('')

const isValid = computed(() => {
  return editedTitle.value.trim().length > 0
})

watch(() => props.isOpen, (open) => {
  if (open) {
    editedTitle.value = props.title
    editedAuthors.value = props.authors
    editedPublisher.value = props.publisher
    error.value = ''
  }
})

const handleSave = () => {
  if (!isValid.value) {
    error.value = 'Title is required'
    return
  }

  emit('save', {
    title: editedTitle.value.trim(),
    authors: editedAuthors.value.trim(),
    publisher: editedPublisher.value.trim()
  })
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

.edit-dialog {
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.98) 0%, rgba(30, 41, 59, 0.98) 100%);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 16px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  max-width: 540px;
  width: 100%;
  animation: slideUp 0.3s ease;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem 1.5rem 1rem;
  border-bottom: 1px solid rgba(59, 130, 246, 0.2);
}

.dialog-header h3 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #e2e8f0;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.close-btn {
  background: rgba(71, 85, 105, 0.2);
  border: 1px solid rgba(71, 85, 105, 0.3);
  border-radius: 8px;
  padding: 0.5rem;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
}

.close-btn:hover {
  background: rgba(71, 85, 105, 0.3);
  border-color: rgba(71, 85, 105, 0.5);
  color: #cbd5e1;
}

.close-btn svg {
  width: 20px;
  height: 20px;
}

.dialog-content {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #cbd5e1;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.form-label svg {
  width: 18px;
  height: 18px;
  color: #3b82f6;
}

.form-input {
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 8px;
  padding: 0.75rem 1rem;
  font-size: 0.938rem;
  color: #e2e8f0;
  transition: all 0.2s ease;
  width: 100%;
}

.form-input::placeholder {
  color: #64748b;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  background: rgba(15, 23, 42, 0.8);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.error-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  color: #fca5a5;
  font-size: 0.875rem;
}

.error-message svg {
  width: 20px;
  height: 20px;
  color: #ef4444;
  flex-shrink: 0;
}

.dialog-actions {
  display: flex;
  gap: 0.75rem;
  padding: 1rem 1.5rem 1.5rem;
  border-top: 1px solid rgba(59, 130, 246, 0.1);
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

.save-btn {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.4);
  color: #3b82f6;
}

.save-btn:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.3);
  border-color: rgba(59, 130, 246, 0.6);
  transform: translateY(-1px);
}

.save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
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

.modal-enter-active .edit-dialog,
.modal-leave-active .edit-dialog {
  transition: transform 0.3s ease;
}

.modal-enter-from .edit-dialog,
.modal-leave-to .edit-dialog {
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

@media (max-width: 640px) {
  .edit-dialog {
    margin: 1rem;
  }

  .dialog-header {
    padding: 1.25rem 1.25rem 0.875rem;
  }

  .dialog-header h3 {
    font-size: 1.125rem;
  }

  .dialog-content {
    padding: 1.25rem;
    gap: 1rem;
  }

  .dialog-actions {
    flex-direction: column;
    padding: 0.875rem 1.25rem 1.25rem;
  }

  .dialog-btn {
    padding: 0.625rem 1.25rem;
    font-size: 0.875rem;
  }
}
</style>
