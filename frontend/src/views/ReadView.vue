<template>
  <div class="read-view">
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Loading book...</p>
    </div>

    <div v-if="error" class="error-state">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="15" y1="9" x2="9" y2="15"></line>
        <line x1="9" y1="9" x2="15" y2="15"></line>
      </svg>
      <h2>{{ error }}</h2>
      <button @click="$router.push('/')" class="back-btn">
        Go Back to Library
      </button>
    </div>

    <div class="reader-container" :class="{ hidden: loading || error }">
      <div class="reader-header">
        <button @click="$router.push('/')" class="back-link">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="19" y1="12" x2="5" y2="12"></line>
            <polyline points="12 19 5 12 12 5"></polyline>
          </svg>
          Back to Library
        </button>
        <h1 class="book-title">{{ bookTitle }}</h1>
      </div>

      <!-- EPUB Experimental Warning -->
      <div v-if="format === 'epub' && !epubWarningDismissed" class="epub-warning">
        <div class="warning-content">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
            <line x1="12" y1="9" x2="12" y2="13"></line>
            <line x1="12" y1="17" x2="12.01" y2="17"></line>
          </svg>
          <div class="warning-text">
            <strong>Experimental EPUB Reader</strong>
            <p>The EPUB reader is experimental and may have bugs. Initial loading can take a long time, especially for large books. For the best experience, consider downloading the book.</p>
          </div>
          <button @click="epubWarningDismissed = true" class="dismiss-btn" title="Dismiss warning">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
      </div>

      <div class="reader-content">
        <!-- Previous Page Button (Desktop) -->
        <button 
          @click="previousPage" 
          class="nav-btn nav-btn-side nav-btn-prev"
          title="Previous page (or use ← arrow key)"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="15 18 9 12 15 6"></polyline>
          </svg>
        </button>

        <!-- EPUB Reader -->
        <div v-show="format === 'epub'" class="epub-reader">
          <div ref="epubContainer" class="epub-container"></div>
        </div>

        <!-- PDF Reader -->
        <div v-show="format === 'pdf'" class="pdf-reader">
          <canvas ref="pdfCanvas" class="pdf-canvas"></canvas>
        </div>

        <!-- Next Page Button (Desktop) -->
        <button 
          @click="nextPage" 
          class="nav-btn nav-btn-side nav-btn-next"
          title="Next page (or use → arrow key)"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"></polyline>
          </svg>
        </button>
      </div>

      <!-- Progress Bar and Controls -->
      <div class="reader-footer">
        <div class="footer-controls">
          <button @click="previousPage" class="nav-btn nav-btn-footer">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="15 18 9 12 15 6"></polyline>
            </svg>
          </button>
          
          <div class="page-control">
            <input 
              v-if="format === 'pdf'"
              type="number" 
              v-model.number="pdfPageNum" 
              @change="onPageInputChange"
              :min="1" 
              :max="pdfNumPages"
              class="page-input"
            />
            <input 
              v-else-if="format === 'epub' && epubLocationsReady"
              type="number" 
              v-model.number="epubCurrentPage" 
              @change="onPageInputChange"
              :min="1" 
              :max="epubTotalPages"
              class="page-input"
            />
            <span v-else class="page-display">{{ format === 'epub' ? 'Page' : '—' }}</span>
            <span class="page-separator">/</span>
            <span class="page-total">{{ format === 'pdf' ? pdfNumPages : (format === 'epub' && epubLocationsReady ? epubTotalPages : '—') }}</span>
          </div>
          
          <button @click="nextPage" class="nav-btn nav-btn-footer">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMeta } from '@/composables/useMeta'
import ePub from 'epubjs'
import * as pdfjsLib from 'pdfjs-dist'
import type { Book as EpubBook, Rendition } from 'epubjs'
import type { PDFDocumentProxy, PDFPageProxy, RenderTask } from 'pdfjs-dist'

const { updateMeta } = useMeta({
  title: 'Read Book',
  description: 'Read your book in the browser'
})

// Set up PDF.js worker - proper way for Vite
pdfjsLib.GlobalWorkerOptions.workerSrc = new URL(
  'pdfjs-dist/build/pdf.worker.min.mjs',
  import.meta.url
).toString()

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const error = ref('')
const bookTitle = ref('')
const format = ref<'epub' | 'pdf' | ''>('')
const epubWarningDismissed = ref(false)

// EPUB specific
const epubContainer = ref<HTMLElement | null>(null)
let epubBook: EpubBook | null = null
let rendition: Rendition | null = null
const epubCurrentPage = ref(1)
const epubTotalPages = ref(0)
const epubLocationsReady = ref(false)

// PDF specific
const pdfCanvas = ref<HTMLCanvasElement | null>(null)
let pdfDoc: PDFDocumentProxy | null = null
let currentRenderTask: RenderTask | null = null
const pdfPageNum = ref(1)
const pdfNumPages = ref(0)
const pdfScale = ref(1.5)
const pdfRendering = ref(false)

const canGoPrev = computed(() => {
  if (format.value === 'pdf') {
    return pdfPageNum.value > 1
  } else if (format.value === 'epub') {
    return epubCurrentPage.value > 1
  }
  return true
})

const canGoNext = computed(() => {
  if (format.value === 'pdf') {
    return pdfPageNum.value < pdfNumPages.value
  } else if (format.value === 'epub') {
    return epubCurrentPage.value < epubTotalPages.value
  }
  return true
})

const currentPageInfo = computed(() => {
  if (format.value === 'pdf') {
    return `${pdfPageNum.value} / ${pdfNumPages.value}`
  } else if (format.value === 'epub') {
    return epubLocationsReady.value ? `${epubCurrentPage.value} / ${epubTotalPages.value}` : 'Loading...'
  }
  return ''
})

async function loadBook() {
  const hash = route.params.hash as string
  if (!hash) {
    error.value = 'No book hash provided'
    loading.value = false
    return
  }

  try {
    const apiUrl = import.meta.env.VITE_API_URL || window.location.origin
    const bookUrl = `${apiUrl}/api/books/${hash}/download?fromreader=true`
    
    // Fetch book metadata first to get the title and format
    const metadataResponse = await fetch(`${apiUrl}/api/books/${hash}`)
    if (!metadataResponse.ok) {
      throw new Error('Failed to load book metadata')
    }
    const metadata = await metadataResponse.json()
    bookTitle.value = metadata.book?.title || 'Unknown Book'
    const bookFormat = (metadata.book?.format || '').toLowerCase()

    // Update meta tags with book info
    if (metadata.book) {
      const coverImage = metadata.book.cover_data || metadata.book.cover_url
      const authorInfo = metadata.book.authors ? ` by ${metadata.book.authors}` : ''
      updateMeta({
        title: `Read ${metadata.book.title} on mArchive`,
        description: `${metadata.book.title}${authorInfo} - Read online on mArchive.`,
        image: coverImage || undefined
      })
    }

    // Set format first so the DOM elements are created
    if (bookFormat === 'epub') {
      format.value = 'epub'
    } else if (bookFormat === 'pdf') {
      format.value = 'pdf'
    } else {
      throw new Error('Unsupported book format')
    }

    // Wait for Vue to update the DOM
    await nextTick()
    
    // Load the appropriate reader
    if (bookFormat === 'epub') {
      await loadEpub(bookUrl)
    } else if (bookFormat === 'pdf') {
      await loadPdf(bookUrl)
    }

    loading.value = false
  } catch (err: any) {
    console.error('Error loading book:', err)
    error.value = err.message || 'Failed to load book'
    loading.value = false
  }
}

async function loadEpub(url: string) {
  // Poll for the container to be ready (max 2 seconds)
  const maxAttempts = 40 // 40 * 50ms = 2 seconds
  let attempts = 0
  
  while (!epubContainer.value && attempts < maxAttempts) {
    await new Promise(resolve => setTimeout(resolve, 50))
    attempts++
  }
  
  if (!epubContainer.value) {
    console.error('EPUB container ref is null after waiting, DOM not ready')
    throw new Error('EPUB container not found in DOM')
  }

  try {
    console.log('EPUB container found, fetching EPUB from:', url)
    
    // Fetch the EPUB file and convert to ArrayBuffer
    const response = await fetch(url, {
      credentials: 'include' // Include cookies for authentication
    })
    
    if (!response.ok) {
      throw new Error(`Failed to fetch EPUB: ${response.status} ${response.statusText}`)
    }
    
    // Convert to ArrayBuffer instead of Blob
    const arrayBuffer = await response.arrayBuffer()
    console.log('EPUB fetched, size:', arrayBuffer.byteLength, 'bytes')
    
    // Create epub book from ArrayBuffer
    epubBook = ePub(arrayBuffer)
    
    console.log('EPUB book created, rendering to container')
    
    // Get container dimensions
    const containerRect = epubContainer.value.getBoundingClientRect()
    console.log('Container dimensions:', containerRect.width, 'x', containerRect.height)
    
    rendition = epubBook.renderTo(epubContainer.value, {
      width: '100%',
      height: '100%',
      spread: 'none',
      flow: 'paginated',
      allowScriptedContent: false // Security: disable scripts in EPUB
    })

    // Handle rendering errors
    rendition.on('rendered', () => {
      console.log('EPUB page rendered successfully')
    })

    rendition.on('displayError', (err: Error) => {
      console.error('EPUB display error:', err)
      error.value = `Failed to display EPUB: ${err.message}`
    })

    rendition.on('relocated', (location: any) => {
      console.log('EPUB relocated:', location)
      if (epubBook && epubBook.locations && location.start) {
        const currentLocation = epubBook.locations.locationFromCfi(location.start.cfi)
        // locationFromCfi returns a number representing the location index
        epubCurrentPage.value = typeof currentLocation === 'number' ? currentLocation : 1
        console.log('Current page:', epubCurrentPage.value, '/', epubTotalPages.value)
      }
    })

    console.log('Starting EPUB display...')
    await rendition.display()
    console.log('EPUB display initiated successfully')
    
    // Generate locations for pagination
    console.log('Generating EPUB locations...')
    await epubBook.locations.generate(1024) // Generate locations with ~1024 chars per page
    epubTotalPages.value = epubBook.locations.length()
    epubLocationsReady.value = true
    console.log('EPUB locations generated:', epubTotalPages.value, 'pages')
    
    // Update current page after locations are ready
    if (rendition.location) {
      const currentLocation = epubBook.locations.locationFromCfi(rendition.location.start.cfi)
      epubCurrentPage.value = typeof currentLocation === 'number' ? currentLocation : 1
    }
    
    // Wait a bit for the content to actually render
    await new Promise(resolve => setTimeout(resolve, 100))
  } catch (err: any) {
    console.error('Error initializing EPUB:', err)
    throw new Error(`Failed to load EPUB: ${err.message || 'Unknown error'}`)
  }
}

async function loadPdf(url: string) {
  try {
    const loadingTask = pdfjsLib.getDocument({
      url,
      // Add proper CORS and credentials handling if needed
      withCredentials: true,
      // Optimize memory usage
      maxImageSize: 1024 * 1024 * 10, // 10MB max per image
      disableAutoFetch: false,
      disableStream: false
    })

    // Handle loading progress if needed
    loadingTask.onProgress = (progress: { loaded: number; total: number }) => {
      console.log(`Loading PDF: ${progress.loaded}/${progress.total}`)
    }

    pdfDoc = await loadingTask.promise
    pdfNumPages.value = pdfDoc.numPages

    await renderPdfPage(pdfPageNum.value)
  } catch (err) {
    console.error('Error loading PDF:', err)
    throw err
  }
}

async function renderPdfPage(pageNumber: number) {
  if (!pdfDoc || !pdfCanvas.value || pdfRendering.value) return

  // Cancel any ongoing render task
  if (currentRenderTask) {
    currentRenderTask.cancel()
    currentRenderTask = null
  }

  pdfRendering.value = true

  try {
    const page: PDFPageProxy = await pdfDoc.getPage(pageNumber)
    
    // Calculate scale to fit the container
    const canvas = pdfCanvas.value
    const container = canvas.parentElement
    if (!container) return

    // Get viewport at scale 1 first to calculate proper scaling
    const viewport = page.getViewport({ scale: 1.0 })
    const containerWidth = container.clientWidth - 32 // Account for padding
    const containerHeight = container.clientHeight - 32
    
    // Calculate scale to fit container while maintaining aspect ratio
    const scaleX = containerWidth / viewport.width
    const scaleY = containerHeight / viewport.height
    const scale = Math.min(scaleX, scaleY, pdfScale.value) // Don't exceed desired scale

    // Get properly scaled viewport
    const scaledViewport = page.getViewport({ scale })

    // Prepare canvas using device pixel ratio for sharp rendering
    const context = canvas.getContext('2d')
    if (!context) return

    const devicePixelRatio = window.devicePixelRatio || 1
    canvas.height = scaledViewport.height * devicePixelRatio
    canvas.width = scaledViewport.width * devicePixelRatio
    canvas.style.height = `${scaledViewport.height}px`
    canvas.style.width = `${scaledViewport.width}px`

    // Scale context to account for device pixel ratio
    context.setTransform(devicePixelRatio, 0, 0, devicePixelRatio, 0, 0)

    const renderContext = {
      canvasContext: context,
      viewport: scaledViewport
    }

    // Start rendering
    currentRenderTask = page.render(renderContext)
    await currentRenderTask.promise
    currentRenderTask = null

    // Clean up page resources
    page.cleanup()
  } catch (err: any) {
    if (err?.name !== 'RenderingCancelledException') {
      console.error('Error rendering PDF page:', err)
      throw err
    }
  } finally {
    pdfRendering.value = false
  }
}

function previousPage() {
  if (format.value === 'epub' && rendition) {
    rendition.prev()
  } else if (format.value === 'pdf' && pdfPageNum.value > 1 && !pdfRendering.value) {
    pdfPageNum.value--
    renderPdfPage(pdfPageNum.value)
  }
}

function nextPage() {
  if (format.value === 'epub' && rendition) {
    rendition.next()
  } else if (format.value === 'pdf' && pdfPageNum.value < pdfNumPages.value && !pdfRendering.value) {
    pdfPageNum.value++
    renderPdfPage(pdfPageNum.value)
  }
}

function onPageInputChange() {
  if (format.value === 'pdf' && pdfDoc && !pdfRendering.value) {
    // Clamp the value
    if (pdfPageNum.value < 1) pdfPageNum.value = 1
    if (pdfPageNum.value > pdfNumPages.value) pdfPageNum.value = pdfNumPages.value
    renderPdfPage(pdfPageNum.value)
  } else if (format.value === 'epub' && epubBook && rendition && epubLocationsReady.value) {
    // Clamp the value
    if (epubCurrentPage.value < 1) epubCurrentPage.value = 1
    if (epubCurrentPage.value > epubTotalPages.value) epubCurrentPage.value = epubTotalPages.value
    
    // Get CFI for the target page and navigate to it
    const targetCfi = epubBook.locations.cfiFromLocation(epubCurrentPage.value)
    if (targetCfi) {
      rendition.display(targetCfi)
    }
  }
}

function handleKeydown(e: KeyboardEvent) {
  // Ignore if user is typing in an input
  if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) {
    return
  }

  if (format.value === 'epub') {
    if (e.key === 'ArrowLeft') {
      e.preventDefault()
      previousPage()
    } else if (e.key === 'ArrowRight') {
      e.preventDefault()
      nextPage()
    }
  } else if (format.value === 'pdf') {
    if (e.key === 'ArrowLeft' || e.key === 'ArrowUp') {
      e.preventDefault()
      previousPage()
    } else if (e.key === 'ArrowRight' || e.key === 'ArrowDown') {
      e.preventDefault()
      nextPage()
    }
  }
}

// Handle window resize for PDF
let resizeTimeout: number | null = null
function handleResize() {
  if (format.value === 'pdf' && pdfDoc) {
    if (resizeTimeout) clearTimeout(resizeTimeout)
    resizeTimeout = window.setTimeout(() => {
      renderPdfPage(pdfPageNum.value)
    }, 250)
  }
}

onMounted(() => {
  loadBook()
  document.addEventListener('keydown', handleKeydown)
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  // Cleanup event listeners
  document.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('resize', handleResize)
  
  if (resizeTimeout) {
    clearTimeout(resizeTimeout)
  }

  // Cancel any ongoing render task
  if (currentRenderTask) {
    currentRenderTask.cancel()
    currentRenderTask = null
  }

  // Cleanup EPUB
  if (rendition) {
    rendition.destroy()
    rendition = null
  }
  if (epubBook) {
    epubBook.destroy()
    epubBook = null
  }

  // Cleanup PDF
  if (pdfDoc) {
    pdfDoc.cleanup()
    pdfDoc.destroy()
    pdfDoc = null
  }
})
</script>

<style scoped>
.read-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  color: #e2e8f0;
}

.hidden {
  display: none !important;
}

.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  gap: 1.5rem;
  padding: 2rem;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 4px solid rgba(59, 130, 246, 0.2);
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.error-state svg {
  width: 64px;
  height: 64px;
  color: #ef4444;
}

.back-btn {
  padding: 0.75rem 1.5rem;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 8px;
  color: #3b82f6;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.back-btn:hover {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.5);
}

.reader-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.reader-header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  padding: 1rem 2rem;
  background: rgba(15, 23, 42, 0.8);
  border-bottom: 1px solid rgba(59, 130, 246, 0.2);
  backdrop-filter: blur(8px);
  flex-shrink: 0;
}

.back-link {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 6px;
  color: #3b82f6;
  text-decoration: none;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.back-link:hover {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.5);
}

.back-link svg {
  width: 20px;
  height: 20px;
}

.book-title {
  flex: 1;
  font-size: 1.25rem;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.epub-warning {
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.3);
  border-left: 4px solid #f59e0b;
  margin: 1rem 2rem;
  border-radius: 8px;
  flex-shrink: 0;
}

.warning-content {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1rem 1.25rem;
}

.warning-content > svg {
  width: 24px;
  height: 24px;
  color: #fbbf24;
  flex-shrink: 0;
  margin-top: 0.125rem;
}

.warning-text {
  flex: 1;
  color: #fef3c7;
}

.warning-text strong {
  display: block;
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 0.375rem;
  color: #fbbf24;
}

.warning-text p {
  margin: 0;
  font-size: 0.875rem;
  line-height: 1.5;
  color: #fde68a;
}

.dismiss-btn {
  background: transparent;
  border: none;
  color: #fbbf24;
  cursor: pointer;
  padding: 0.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.dismiss-btn:hover {
  background: rgba(245, 158, 11, 0.2);
}

.dismiss-btn svg {
  width: 20px;
  height: 20px;
}

.reader-content {
  flex: 1;
  display: flex;
  align-items: stretch;
  overflow: hidden;
  background: #1e293b;
  position: relative;
  min-height: 0;
}

.nav-btn {
  padding: 0.75rem;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 8px;
  color: #3b82f6;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 44px;
  min-height: 44px;
}

.nav-btn:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.5);
}

.nav-btn:active {
  transform: none;
  background: rgba(59, 130, 246, 0.15);
}

.nav-btn svg {
  width: 24px;
  height: 24px;
}

/* Side navigation buttons (Desktop) */
.nav-btn-side {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  z-index: 10;
  padding: 1rem;
  min-width: 48px;
  min-height: 48px;
}

.nav-btn-prev {
  left: 1rem;
}

.nav-btn-next {
  right: 1rem;
}

/* Footer navigation buttons - hidden on desktop, visible on mobile */
.nav-btn-footer {
  display: none;
}

.epub-reader,
.pdf-reader {
  flex: 1;
  width: 100%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem 2rem;
  background: #1e293b;
  min-height: 0;
}

.epub-container {
  width: 100%;
  height: 100%;
  max-width: 800px;
  max-height: 100%;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}

/* Ensure epub.js iframe fills the container */
.epub-container :deep(iframe) {
  border: none;
}

.pdf-canvas {
  max-width: 100%;
  max-height: 100%;
  border-radius: 8px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.3);
  background: white;
}

.reader-footer {
  padding: 1rem 2rem;
  background: rgba(15, 23, 42, 0.95);
  border-top: 1px solid rgba(59, 130, 246, 0.2);
  backdrop-filter: blur(8px);
  flex-shrink: 0;
}

.footer-controls {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1.5rem;
}

.page-control {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  color: #94a3b8;
}

.page-input {
  width: 60px;
  padding: 0.5rem;
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 6px;
  color: #e2e8f0;
  text-align: center;
  font-size: 0.875rem;
  transition: all 0.2s ease;
}

.page-input:focus {
  outline: none;
  border-color: rgba(59, 130, 246, 0.6);
  background: rgba(30, 41, 59, 1);
}

.page-input::-webkit-inner-spin-button,
.page-input::-webkit-outer-spin-button {
  opacity: 1;
}

.page-display {
  color: #64748b;
}

.page-separator {
  color: #64748b;
  margin: 0 0.25rem;
}

.page-total {
  color: #94a3b8;
  font-weight: 500;
}

@media (max-width: 768px) {
  .reader-header {
    padding: 1rem;
  }

  .book-title {
    font-size: 1rem;
  }

  /* Hide side navigation on mobile */
  .nav-btn-side {
    display: none;
  }

  /* Show footer navigation buttons on mobile */
  .nav-btn-footer {
    display: flex;
  }

  .epub-reader,
  .pdf-reader {
    padding: 1rem;
  }

  .reader-footer {
    padding: 0.75rem 1rem;
  }

  .footer-controls {
    gap: 1rem;
  }
}

</style>
