<template>
  <div class="upload-container">
    <div class="upload-card">
      <h1>Upload Book</h1>
      <p class="subtitle">Add your own books to the library</p>

      <form @submit.prevent="handleSubmit" class="upload-form">
        <!-- Book File Upload -->
        <div class="form-group">
          <label for="book-file" class="required">Book File</label>
          <div class="file-input-wrapper">
            <input
              type="file"
              id="book-file"
              ref="bookFileInput"
              @change="handleBookFileChange"
              accept=".pdf,.epub,.mobi,.azw3,.djvu,.fb2,.txt"
              required
            />
            <div v-if="bookFile" class="file-info">
              <span class="file-name">{{ bookFile.name }}</span>
              <span class="file-size">{{ formatFileSize(bookFile.size) }}</span>
            </div>
          </div>
          <small>Supported formats: PDF, EPUB, MOBI, AZW3, DJVU, FB2, TXT (Max: 500MB)</small>
        </div>

        <!-- Cover Image Upload -->
        <div class="form-group">
          <label for="cover-file">Cover Image (Optional)</label>
          <div class="file-input-wrapper">
            <input
              type="file"
              id="cover-file"
              ref="coverFileInput"
              @change="handleCoverFileChange"
              accept=".jpg,.jpeg,.png,.webp"
            />
          </div>
          <small>Supported formats: JPG, PNG, WEBP (Max: 10MB)</small>
          
          <!-- Cover Preview -->
          <div v-if="coverPreview" class="cover-preview">
            <img :src="coverPreview" alt="Cover preview" />
            <button type="button" @click="editCover" class="btn-secondary">Edit Cover</button>
          </div>
        </div>

        <!-- Metadata Fields -->
        <div class="form-group">
          <label for="title" class="required">Title</label>
          <input
            type="text"
            id="title"
            v-model="formData.title"
            placeholder="Enter book title"
            required
          />
        </div>

        <div class="form-group">
          <label for="authors">Authors</label>
          <input
            type="text"
            id="authors"
            v-model="formData.authors"
            placeholder="Enter author names"
          />
        </div>

        <div class="form-group">
          <label for="publisher">Publisher</label>
          <input
            type="text"
            id="publisher"
            v-model="formData.publisher"
            placeholder="Enter publisher name"
          />
        </div>

        <div class="form-group">
          <label for="language">Language</label>
          <input
            type="text"
            id="language"
            v-model="formData.language"
            placeholder="e.g., English"
          />
        </div>

        <!-- Error Message -->
        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>

        <!-- Success Message -->
        <div v-if="successMessage" class="success-message">
          {{ successMessage }}
          <router-link :to="`/book/${uploadedBookHash}`" class="view-book-link">
            View Book
          </router-link>
        </div>

        <!-- Upload Progress -->
        <div v-if="uploading" class="upload-progress">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: uploadProgress + '%' }"></div>
          </div>
          <p>Uploading... {{ uploadProgress }}%</p>
        </div>

        <!-- Submit Button -->
        <div class="form-actions">
          <button type="submit" :disabled="uploading || !isFormValid" class="btn-primary">
            {{ uploading ? 'Uploading...' : 'Upload Book' }}
          </button>
          <button type="button" @click="resetForm" class="btn-secondary" :disabled="uploading">
            Reset
          </button>
        </div>
      </form>
    </div>

    <!-- Cover Editor Modal -->
    <CoverImageEditor
      v-if="showCoverEditor"
      :imageUrl="coverPreview"
      @save="handleCoverEdited"
      @close="showCoverEditor = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { booksApi } from '@/api/books'
import CoverImageEditor from '@/components/CoverImageEditor.vue'

const router = useRouter()

// Form data
const bookFile = ref<File | null>(null)
const coverFile = ref<File | null>(null)
const coverPreview = ref<string>('')
const formData = ref({
  title: '',
  authors: '',
  publisher: '',
  language: ''
})

// UI state
const uploading = ref(false)
const uploadProgress = ref(0)
const errorMessage = ref('')
const successMessage = ref('')
const uploadedBookHash = ref('')
const showCoverEditor = ref(false)

// Refs
const bookFileInput = ref<HTMLInputElement>()
const coverFileInput = ref<HTMLInputElement>()

// Computed
const isFormValid = computed(() => {
  return bookFile.value !== null && formData.value.title.trim() !== ''
})

// File handlers
const handleBookFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const file = target.files[0]
    
    // Validate file size
    if (file.size > 500 * 1024 * 1024) {
      errorMessage.value = 'Book file is too large. Maximum size is 500MB.'
      target.value = ''
      return
    }
    
    bookFile.value = file
    errorMessage.value = ''
  }
}

const handleCoverFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const file = target.files[0]
    
    // Validate file size
    if (file.size > 10 * 1024 * 1024) {
      errorMessage.value = 'Cover image is too large. Maximum size is 10MB.'
      target.value = ''
      return
    }
    
    coverFile.value = file
    
    // Create preview
    const reader = new FileReader()
    reader.onload = (e) => {
      coverPreview.value = e.target?.result as string
    }
    reader.readAsDataURL(file)
    
    errorMessage.value = ''
  }
}

const editCover = () => {
  if (coverPreview.value) {
    showCoverEditor.value = true
  }
}

const handleCoverEdited = (editedImageBlob: Blob) => {
  // Convert blob to file
  const file = new File([editedImageBlob], coverFile.value?.name || 'cover.jpg', {
    type: 'image/jpeg'
  })
  coverFile.value = file
  
  // Update preview
  const reader = new FileReader()
  reader.onload = (e) => {
    coverPreview.value = e.target?.result as string
  }
  reader.readAsDataURL(file)
  
  showCoverEditor.value = false
}

const handleSubmit = async () => {
  if (!isFormValid.value) return
  
  uploading.value = true
  uploadProgress.value = 0
  errorMessage.value = ''
  successMessage.value = ''
  
  try {
    // Create FormData
    const formDataToSend = new FormData()
    formDataToSend.append('book', bookFile.value!)
    if (coverFile.value) {
      formDataToSend.append('cover', coverFile.value)
    }
    formDataToSend.append('title', formData.value.title)
    formDataToSend.append('authors', formData.value.authors)
    formDataToSend.append('publisher', formData.value.publisher)
    formDataToSend.append('language', formData.value.language)
    
    // Simulate progress (since we don't have real progress tracking)
    const progressInterval = setInterval(() => {
      if (uploadProgress.value < 90) {
        uploadProgress.value += 10
      }
    }, 200)
    
    // Upload
    const response = await booksApi.uploadBook(formDataToSend)
    
    clearInterval(progressInterval)
    uploadProgress.value = 100
    
    // Show success message
    successMessage.value = 'Book uploaded successfully!'
    uploadedBookHash.value = response.book.hash
    
    // Reset form after a delay
    setTimeout(() => {
      resetForm()
    }, 3000)
    
  } catch (error: any) {
    errorMessage.value = error.response?.data?.message || 'Failed to upload book. Please try again.'
  } finally {
    uploading.value = false
  }
}

const resetForm = () => {
  bookFile.value = null
  coverFile.value = null
  coverPreview.value = ''
  formData.value = {
    title: '',
    authors: '',
    publisher: '',
    language: ''
  }
  errorMessage.value = ''
  successMessage.value = ''
  uploadProgress.value = 0
  
  if (bookFileInput.value) bookFileInput.value.value = ''
  if (coverFileInput.value) coverFileInput.value.value = ''
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}
</script>

<style scoped>
.upload-container {
  min-height: calc(100vh - 73px);
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
}

.upload-card {
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 16px;
  padding: 2.5rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(12px);
  animation: fadeInUp 0.6s ease;
}

h1 {
  margin: 0 0 0.5rem 0;
  font-size: 2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #fff 0%, #94a3b8 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: -0.02em;
}

.subtitle {
  color: #64748b;
  margin: 0 0 2rem 0;
  font-size: 1.125rem;
}

.upload-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

label {
  font-weight: 600;
  color: #e2e8f0;
  font-size: 0.95rem;
}

label.required::after {
  content: ' *';
  color: #f87171;
}

input[type="text"],
input[type="file"] {
  padding: 0.875rem;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 8px;
  font-size: 1rem;
  color: #e2e8f0;
  transition: all 0.2s ease;
}

input[type="text"]:focus,
input[type="file"]:focus {
  outline: none;
  border-color: #3b82f6;
  background: rgba(15, 23, 42, 0.8);
}

input[type="text"]::placeholder {
  color: #64748b;
}

input[type="file"] {
  cursor: pointer;
}

input[type="file"]::file-selector-button {
  padding: 0.5rem 1rem;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 6px;
  color: #3b82f6;
  font-weight: 500;
  cursor: pointer;
  margin-right: 1rem;
  transition: all 0.2s ease;
}

input[type="file"]::file-selector-button:hover {
  background: rgba(59, 130, 246, 0.2);
  border-color: #3b82f6;
}

.file-input-wrapper {
  position: relative;
}

.file-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 0.5rem;
  padding: 0.75rem;
  background: rgba(59, 130, 246, 0.05);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 8px;
}

.file-name {
  font-weight: 500;
  color: #e2e8f0;
}

.file-size {
  color: #94a3b8;
  font-size: 0.9rem;
}

small {
  color: #64748b;
  font-size: 0.85rem;
}

.cover-preview {
  margin-top: 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  background: rgba(15, 23, 42, 0.4);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 12px;
}

.cover-preview img {
  max-width: 200px;
  max-height: 300px;
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}

.error-message {
  padding: 1rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  color: #fca5a5;
}

.success-message {
  padding: 1rem;
  background: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.3);
  border-radius: 8px;
  color: #86efac;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.view-book-link {
  color: #3b82f6;
  text-decoration: none;
  font-weight: 600;
  transition: color 0.2s ease;
}

.view-book-link:hover {
  color: #60a5fa;
  text-decoration: underline;
}

.upload-progress {
  margin: 1rem 0;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #3b82f6 0%, #60a5fa 100%);
  transition: width 0.3s ease;
  box-shadow: 0 0 10px rgba(59, 130, 246, 0.5);
}

.upload-progress p {
  margin-top: 0.5rem;
  text-align: center;
  color: #94a3b8;
  font-size: 0.9rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.btn-primary,
.btn-secondary {
  padding: 0.875rem 1.5rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-primary {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
  flex: 1;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(59, 130, 246, 0.4);
}

.btn-primary:disabled {
  background: rgba(59, 130, 246, 0.3);
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.btn-secondary {
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.3);
  color: #3b82f6;
}

.btn-secondary:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.2);
  border-color: #3b82f6;
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 768px) {
  .upload-container {
    padding: 1.5rem 1rem;
  }

  .upload-card {
    padding: 1.5rem;
  }

  h1 {
    font-size: 1.75rem;
  }

  .subtitle {
    font-size: 1rem;
  }
}

@media (max-width: 480px) {
  .upload-container {
    padding: 1rem 0.75rem;
  }

  .upload-card {
    padding: 1.25rem;
  }

  h1 {
    font-size: 1.5rem;
  }

  .form-actions {
    flex-direction: column;
  }

  .btn-primary,
  .btn-secondary {
    width: 100%;
  }
}
</style>
