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
  max-width: 800px;
  margin: 2rem auto;
  padding: 0 1rem;
}

.upload-card {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

h1 {
  margin: 0 0 0.5rem 0;
  color: #2c3e50;
}

.subtitle {
  color: #7f8c8d;
  margin: 0 0 2rem 0;
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
  color: #2c3e50;
}

label.required::after {
  content: ' *';
  color: #e74c3c;
}

input[type="text"],
input[type="file"] {
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
}

input[type="text"]:focus {
  outline: none;
  border-color: #3498db;
}

.file-input-wrapper {
  position: relative;
}

.file-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 0.5rem;
  padding: 0.5rem;
  background: #f8f9fa;
  border-radius: 4px;
}

.file-name {
  font-weight: 500;
  color: #2c3e50;
}

.file-size {
  color: #7f8c8d;
  font-size: 0.9rem;
}

small {
  color: #7f8c8d;
  font-size: 0.85rem;
}

.cover-preview {
  margin-top: 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.cover-preview img {
  max-width: 200px;
  max-height: 300px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.error-message {
  padding: 1rem;
  background: #fee;
  border: 1px solid #fcc;
  border-radius: 6px;
  color: #c33;
}

.success-message {
  padding: 1rem;
  background: #efe;
  border: 1px solid #cfc;
  border-radius: 6px;
  color: #3c3;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.view-book-link {
  color: #3498db;
  text-decoration: none;
  font-weight: 600;
}

.view-book-link:hover {
  text-decoration: underline;
}

.upload-progress {
  margin: 1rem 0;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: #ecf0f1;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #3498db;
  transition: width 0.3s ease;
}

.upload-progress p {
  margin-top: 0.5rem;
  text-align: center;
  color: #7f8c8d;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.btn-primary,
.btn-secondary {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-primary {
  background: #3498db;
  color: white;
  flex: 1;
}

.btn-primary:hover:not(:disabled) {
  background: #2980b9;
}

.btn-primary:disabled {
  background: #bdc3c7;
  cursor: not-allowed;
}

.btn-secondary {
  background: #ecf0f1;
  color: #2c3e50;
}

.btn-secondary:hover:not(:disabled) {
  background: #d5dbdb;
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
