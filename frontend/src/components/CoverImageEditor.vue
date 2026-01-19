<template>
  <div class="editor-overlay" @click.self="$emit('close')">
    <div class="editor-modal">
      <div class="editor-header">
        <h2>Edit Cover Image</h2>
        <button @click="$emit('close')" class="close-btn">&times;</button>
      </div>

      <div class="editor-body">
        <!-- Canvas for image editing -->
        <div class="canvas-container">
          <canvas ref="canvas"></canvas>
        </div>

        <!-- Controls -->
        <div class="controls">
          <div class="control-group">
            <label>Crop Aspect Ratio</label>
            <select v-model="aspectRatio" @change="updateAspectRatio">
              <option value="free">Free</option>
              <option value="2:3">2:3 (Book Cover)</option>
              <option value="1:1">1:1 (Square)</option>
              <option value="3:4">3:4</option>
            </select>
          </div>

          <div class="control-group">
            <label>Rotation</label>
            <div class="rotation-buttons">
              <button @click="rotate(-90)" class="control-btn">↶ 90°</button>
              <button @click="rotate(90)" class="control-btn">↷ 90°</button>
            </div>
          </div>

          <div class="control-group">
            <label>Brightness: {{ brightness }}</label>
            <input
              type="range"
              v-model.number="brightness"
              min="-100"
              max="100"
              @input="applyFilters"
            />
          </div>

          <div class="control-group">
            <label>Contrast: {{ contrast }}</label>
            <input
              type="range"
              v-model.number="contrast"
              min="-100"
              max="100"
              @input="applyFilters"
            />
          </div>
        </div>
      </div>

      <div class="editor-footer">
        <button @click="reset" class="btn-secondary">Reset</button>
        <button @click="save" class="btn-primary">Save Changes</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'

const props = defineProps<{
  imageUrl: string
}>()

const emit = defineEmits<{
  save: [blob: Blob]
  close: []
}>()

// Refs
const canvas = ref<HTMLCanvasElement>()
const ctx = ref<CanvasRenderingContext2D | null>(null)
const image = ref<HTMLImageElement>()
const originalImage = ref<HTMLImageElement>()

// State
const aspectRatio = ref('2:3')
const rotation = ref(0)
const brightness = ref(0)
const contrast = ref(0)
const cropX = ref(0)
const cropY = ref(0)
const cropWidth = ref(0)
const cropHeight = ref(0)

// Dragging state for crop
const isDragging = ref(false)
const dragStartX = ref(0)
const dragStartY = ref(0)

onMounted(() => {
  if (canvas.value) {
    ctx.value = canvas.value.getContext('2d')
    loadImage()
  }
})

const loadImage = () => {
  const img = new Image()
  img.crossOrigin = 'anonymous'
  img.onload = () => {
    originalImage.value = img
    image.value = img
    
    // Set canvas size
    if (canvas.value) {
      canvas.value.width = Math.min(img.width, 800)
      canvas.value.height = (canvas.value.width / img.width) * img.height
      
      // Initialize crop area
      cropWidth.value = canvas.value.width
      cropHeight.value = canvas.value.height
      
      drawImage()
    }
  }
  img.src = props.imageUrl
}

const drawImage = () => {
  if (!ctx.value || !canvas.value || !image.value) return
  
  // Clear canvas
  ctx.value.clearRect(0, 0, canvas.value.width, canvas.value.height)
  
  // Save context state
  ctx.value.save()
  
  // Apply rotation
  if (rotation.value !== 0) {
    ctx.value.translate(canvas.value.width / 2, canvas.value.height / 2)
    ctx.value.rotate((rotation.value * Math.PI) / 180)
    ctx.value.translate(-canvas.value.width / 2, -canvas.value.height / 2)
  }
  
  // Draw image
  ctx.value.drawImage(image.value, 0, 0, canvas.value.width, canvas.value.height)
  
  // Restore context
  ctx.value.restore()
  
  // Apply filters
  applyFilters()
}

const applyFilters = () => {
  if (!ctx.value || !canvas.value) return
  
  const imageData = ctx.value.getImageData(0, 0, canvas.value.width, canvas.value.height)
  const data = imageData.data
  
  const brightnessValue = brightness.value
  const contrastValue = contrast.value * 2.55
  const factor = (259 * (contrastValue + 255)) / (255 * (259 - contrastValue))
  
  for (let i = 0; i < data.length; i += 4) {
    // Apply brightness
    data[i] += brightnessValue
    data[i + 1] += brightnessValue
    data[i + 2] += brightnessValue
    
    // Apply contrast
    data[i] = factor * (data[i] - 128) + 128
    data[i + 1] = factor * (data[i + 1] - 128) + 128
    data[i + 2] = factor * (data[i + 2] - 128) + 128
    
    // Clamp values
    data[i] = Math.max(0, Math.min(255, data[i]))
    data[i + 1] = Math.max(0, Math.min(255, data[i + 1]))
    data[i + 2] = Math.max(0, Math.min(255, data[i + 2]))
  }
  
  ctx.value.putImageData(imageData, 0, 0)
}

const rotate = (degrees: number) => {
  rotation.value = (rotation.value + degrees) % 360
  drawImage()
}

const updateAspectRatio = () => {
  // Update crop dimensions based on aspect ratio
  if (canvas.value) {
    const ratios: Record<string, number> = {
      'free': 0,
      '2:3': 2 / 3,
      '1:1': 1,
      '3:4': 3 / 4
    }
    
    const ratio = ratios[aspectRatio.value]
    if (ratio > 0) {
      cropWidth.value = canvas.value.width
      cropHeight.value = cropWidth.value / ratio
      
      if (cropHeight.value > canvas.value.height) {
        cropHeight.value = canvas.value.height
        cropWidth.value = cropHeight.value * ratio
      }
    }
  }
}

const reset = () => {
  rotation.value = 0
  brightness.value = 0
  contrast.value = 0
  aspectRatio.value = '2:3'
  loadImage()
}

const save = () => {
  if (!canvas.value) return
  
  canvas.value.toBlob((blob) => {
    if (blob) {
      emit('save', blob)
    }
  }, 'image/jpeg', 0.9)
}
</script>

<style scoped>
.editor-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.editor-modal {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 900px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e0e0e0;
}

.editor-header h2 {
  margin: 0;
  color: #2c3e50;
}

.close-btn {
  background: none;
  border: none;
  font-size: 2rem;
  color: #7f8c8d;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.close-btn:hover {
  background: #f0f0f0;
}

.editor-body {
  display: flex;
  gap: 2rem;
  padding: 1.5rem;
  overflow-y: auto;
}

.canvas-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  border-radius: 8px;
  padding: 1rem;
}

canvas {
  max-width: 100%;
  max-height: 500px;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.controls {
  width: 250px;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.control-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.control-group label {
  font-weight: 600;
  color: #2c3e50;
  font-size: 0.9rem;
}

select,
input[type="range"] {
  width: 100%;
}

select {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.9rem;
}

input[type="range"] {
  -webkit-appearance: none;
  appearance: none;
  height: 6px;
  background: #ddd;
  border-radius: 3px;
  outline: none;
}

input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 16px;
  height: 16px;
  background: #3498db;
  border-radius: 50%;
  cursor: pointer;
}

input[type="range"]::-moz-range-thumb {
  width: 16px;
  height: 16px;
  background: #3498db;
  border-radius: 50%;
  cursor: pointer;
  border: none;
}

.rotation-buttons {
  display: flex;
  gap: 0.5rem;
}

.control-btn {
  flex: 1;
  padding: 0.5rem;
  background: #ecf0f1;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 600;
  color: #2c3e50;
}

.control-btn:hover {
  background: #d5dbdb;
}

.editor-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding: 1.5rem;
  border-top: 1px solid #e0e0e0;
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
}

.btn-primary:hover {
  background: #2980b9;
}

.btn-secondary {
  background: #ecf0f1;
  color: #2c3e50;
}

.btn-secondary:hover {
  background: #d5dbdb;
}

@media (max-width: 768px) {
  .editor-body {
    flex-direction: column;
  }
  
  .controls {
    width: 100%;
  }
}
</style>
