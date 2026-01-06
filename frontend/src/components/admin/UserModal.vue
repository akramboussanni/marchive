<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal">
      <h2>{{ isEdit ? 'Edit User' : 'Create New User' }}</h2>
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label>Username</label>
          <input v-model="formData.username" type="text" required />
        </div>
        <div v-if="!isEdit" class="form-group">
          <label>Password</label>
          <input v-model="formData.password" type="password" required />
        </div>
        <div class="form-group">
          <label>Role</label>
          <select v-model="formData.role">
            <option value="user">User</option>
            <option value="admin">Admin</option>
          </select>
        </div>
        <div class="modal-actions">
          <button type="button" class="secondary-button" @click="$emit('close')">Cancel</button>
          <button type="submit" class="primary-button">{{ isEdit ? 'Update' : 'Create' }}</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  user?: any
  isEdit?: boolean
}>()

const emit = defineEmits<{
  close: []
  submit: [data: any]
}>()

const formData = ref({
  username: props.user?.username || '',
  password: '',
  role: props.user?.user_role || props.user?.role || 'user'
})

watch(() => props.user, (newUser) => {
  if (newUser) {
    formData.value = {
      username: newUser.username || '',
      password: '',
      role: newUser.user_role || newUser.role || 'user'
    }
  }
}, { immediate: true })

function handleSubmit() {
  emit('submit', formData.value)
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

.modal {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 16px;
  padding: 2rem;
  max-width: 500px;
  width: 90%;
  animation: slideUp 0.3s ease;
}

.modal h2 {
  color: #e2e8f0;
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  color: #94a3b8;
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.75rem;
  background: rgba(15, 23, 42, 0.8);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 8px;
  color: #e2e8f0;
  font-size: 1rem;
  transition: all 0.2s ease;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
}

.primary-button {
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.primary-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(59, 130, 246, 0.3);
}

.secondary-button {
  padding: 0.75rem 1.5rem;
  background: rgba(100, 116, 139, 0.1);
  color: #94a3b8;
  border: 1px solid rgba(100, 116, 139, 0.3);
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.secondary-button:hover {
  background: rgba(100, 116, 139, 0.2);
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
