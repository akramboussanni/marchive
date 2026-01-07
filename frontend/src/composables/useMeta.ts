import { onMounted, onUnmounted, watch } from 'vue'

interface MetaOptions {
  title?: string
  description?: string
  image?: string
  url?: string
}

export function useMeta(options: MetaOptions) {
  const setMeta = () => {
    const baseUrl = window.location.origin
    const currentUrl = options.url || window.location.href

    // Set title
    if (options.title) {
      document.title = `${options.title} - mArchive`
    }

    // Helper to update or create meta tag
    const updateMeta = (property: string, content: string, isName = false) => {
      const attr = isName ? 'name' : 'property'
      let meta = document.querySelector(`meta[${attr}="${property}"]`)
      if (!meta) {
        meta = document.createElement('meta')
        meta.setAttribute(attr, property)
        document.head.appendChild(meta)
      }
      meta.setAttribute('content', content)
    }

    // OpenGraph tags for Discord, Twitter, etc.
    updateMeta('og:title', options.title || 'mArchive')
    updateMeta('og:description', options.description || 'Your personal digital library')
    updateMeta('og:url', currentUrl)
    updateMeta('og:type', 'website')
    
    if (options.image) {
      const imageUrl = options.image.startsWith('http') ? options.image : `${baseUrl}${options.image}`
      updateMeta('og:image', imageUrl)
    } else {
      updateMeta('og:image', `${baseUrl}/favicon.ico`)
    }

    // Twitter cards
    updateMeta('twitter:card', options.image ? 'summary_large_image' : 'summary')
    updateMeta('twitter:title', options.title || 'mArchive')
    updateMeta('twitter:description', options.description || 'Your personal digital library')
    if (options.image) {
      const imageUrl = options.image.startsWith('http') ? options.image : `${baseUrl}${options.image}`
      updateMeta('twitter:image', imageUrl)
    }

    // Standard meta tags
    updateMeta('description', options.description || 'Your personal digital library', true)
  }

  onMounted(() => {
    setMeta()
  })

  // Return function to update meta dynamically
  return {
    updateMeta: (newOptions: MetaOptions) => {
      Object.assign(options, newOptions)
      setMeta()
    }
  }
}
