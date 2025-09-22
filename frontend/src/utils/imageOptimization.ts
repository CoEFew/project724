/**
 * Image Optimization Utilities
 * 
 * This module provides utilities for optimizing image loading and performance
 * across the application. It includes lazy loading, responsive images,
 * format conversion, and caching strategies.
 */

/**
 * Image loading configuration
 */
export interface ImageConfig {
  src: string
  alt: string
  loading?: 'lazy' | 'eager'
  sizes?: string
  srcset?: string
  placeholder?: string
  fallback?: string
  quality?: number
  format?: 'webp' | 'png' | 'jpg' | 'auto'
}

/**
 * Responsive image breakpoints
 */
export const BREAKPOINTS = {
  sm: '640px',
  md: '768px',
  lg: '1024px',
  xl: '1280px',
  '2xl': '1536px'
} as const

/**
 * Image quality settings
 */
export const QUALITY_SETTINGS = {
  high: 90,
  medium: 75,
  low: 60,
  thumbnail: 40
} as const

/**
 * Generate responsive srcset for different screen sizes
 * @param baseSrc - Base image source path
 * @param sizes - Array of sizes to generate
 * @param format - Image format (webp, png, jpg)
 * @returns Responsive srcset string
 */
export function generateSrcSet(
  baseSrc: string,
  sizes: number[] = [320, 640, 768, 1024, 1280],
  format: 'webp' | 'png' | 'jpg' = 'webp'
): string {
  try {
    const baseName = baseSrc.replace(/\.[^/.]+$/, '')
    const extension = format === 'webp' ? '.webp' : format === 'png' ? '.png' : '.jpg'
    
    return sizes
      .map(size => `${baseName}_${size}w${extension} ${size}w`)
      .join(', ')
  } catch (error) {
    console.warn('Error generating srcset:', error)
    return baseSrc
  }
}

/**
 * Generate responsive sizes attribute
 * @param breakpoints - Breakpoint configuration
 * @returns Sizes attribute string
 */
export function generateSizes(breakpoints: Record<string, string> = BREAKPOINTS): string {
  try {
    return Object.entries(breakpoints)
      .map(([bp, size]) => `(max-width: ${size}) ${size}`)
      .join(', ') + ', 100vw'
  } catch (error) {
    console.warn('Error generating sizes:', error)
    return '100vw'
  }
}

/**
 * Check if WebP is supported by the browser
 * @returns Promise<boolean>
 */
export async function isWebPSupported(): Promise<boolean> {
  try {
    if (!window.createImageBitmap) return false
    
    const webpData = 'data:image/webp;base64,UklGRhoBAABXRUJQVlA4TA4AAAAvAAAAABsYAAABAAEAAAAAAABAAEAABAAAAAA'
    const blob = await fetch(webpData).then(r => r.blob())
    return createImageBitmap(blob).then(() => true, () => false)
  } catch (error) {
    console.warn('WebP support check failed:', error)
    return false
  }
}

/**
 * Get optimized image source based on browser support
 * @param baseSrc - Base image source
 * @param preferredFormat - Preferred format (webp, png, jpg)
 * @returns Promise<string> - Optimized image source
 */
export async function getOptimizedSrc(
  baseSrc: string,
  preferredFormat: 'webp' | 'png' | 'jpg' = 'webp'
): Promise<string> {
  try {
    // If already optimized, return as is
    if (baseSrc.includes('.webp') || baseSrc.includes('_optimized')) {
      return baseSrc
    }

    // Check WebP support
    if (preferredFormat === 'webp' && await isWebPSupported()) {
      const webpSrc = baseSrc.replace(/\.(png|jpg|jpeg)$/i, '.webp')
      return webpSrc
    }

    // Fallback to original format
    return baseSrc
  } catch (error) {
    console.warn('Error getting optimized src:', error)
    return baseSrc
  }
}

/**
 * Create a placeholder image data URL
 * @param width - Image width
 * @param height - Image height
 * @param color - Background color
 * @returns Data URL string
 */
export function createPlaceholder(
  width: number = 100,
  height: number = 100,
  color: string = '#1e293b'
): string {
  try {
    const canvas = document.createElement('canvas')
    canvas.width = width
    canvas.height = height
    const ctx = canvas.getContext('2d')
    
    if (!ctx) return ''
    
    ctx.fillStyle = color
    ctx.fillRect(0, 0, width, height)
    
    return canvas.toDataURL('image/png', 0.1)
  } catch (error) {
    console.warn('Error creating placeholder:', error)
    return ''
  }
}

/**
 * Preload critical images
 * @param imageSrcs - Array of image sources to preload
 * @param priority - Whether to use high priority loading
 */
export async function preloadImages(
  imageSrcs: string[],
  priority: boolean = false
): Promise<void> {
  try {
    const promises = imageSrcs.map(src => {
      return new Promise<void>((resolve, reject) => {
        const img = new Image()
        img.onload = () => resolve()
        img.onerror = () => reject(new Error(`Failed to load ${src}`))
        img.src = src
        img.loading = priority ? 'eager' : 'lazy'
      })
    })
    
    await Promise.all(promises)
  } catch (error) {
    console.warn('Error preloading images:', error)
  }
}

/**
 * Lazy load images with intersection observer
 * @param selector - CSS selector for images to lazy load
 * @param options - Intersection observer options
 */
export function setupLazyLoading(
  selector: string = 'img[data-src]',
  options: IntersectionObserverInit = {
    rootMargin: '50px 0px',
    threshold: 0.01
  }
): void {
  try {
    if (!('IntersectionObserver' in window)) {
      // Fallback for older browsers
      const images = document.querySelectorAll(selector)
      images.forEach(img => {
        const dataSrc = img.getAttribute('data-src')
        if (dataSrc) {
          img.setAttribute('src', dataSrc)
          img.removeAttribute('data-src')
        }
      })
      return
    }

    const imageObserver = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          const img = entry.target as HTMLImageElement
          const dataSrc = img.getAttribute('data-src')
          
          if (dataSrc) {
            img.src = dataSrc
            img.removeAttribute('data-src')
            img.classList.remove('lazy')
            imageObserver.unobserve(img)
          }
        }
      })
    }, options)

    const images = document.querySelectorAll(selector)
    images.forEach(img => imageObserver.observe(img))
  } catch (error) {
    console.warn('Error setting up lazy loading:', error)
  }
}

/**
 * Get image dimensions from file
 * @param file - Image file
 * @returns Promise<{width: number, height: number}>
 */
export function getImageDimensions(file: File): Promise<{width: number, height: number}> {
  return new Promise((resolve, reject) => {
    try {
      const img = new Image()
      img.onload = () => {
        resolve({ width: img.naturalWidth, height: img.naturalHeight })
      }
      img.onerror = () => reject(new Error('Failed to load image'))
      img.src = URL.createObjectURL(file)
    } catch (error) {
      reject(error)
    }
  })
}

/**
 * Compress image using canvas
 * @param file - Image file to compress
 * @param quality - Compression quality (0-1)
 * @param maxWidth - Maximum width
 * @param maxHeight - Maximum height
 * @returns Promise<Blob> - Compressed image blob
 */
export function compressImage(
  file: File,
  quality: number = 0.8,
  maxWidth: number = 1920,
  maxHeight: number = 1080
): Promise<Blob> {
  return new Promise((resolve, reject) => {
    try {
      const canvas = document.createElement('canvas')
      const ctx = canvas.getContext('2d')
      const img = new Image()
      
      if (!ctx) {
        reject(new Error('Canvas context not available'))
        return
      }
      
      img.onload = () => {
        // Calculate new dimensions
        let { width, height } = img
        if (width > maxWidth || height > maxHeight) {
          const ratio = Math.min(maxWidth / width, maxHeight / height)
          width *= ratio
          height *= ratio
        }
        
        canvas.width = width
        canvas.height = height
        
        // Draw and compress
        ctx.drawImage(img, 0, 0, width, height)
        canvas.toBlob(
          (blob) => {
            if (blob) {
              resolve(blob)
            } else {
              reject(new Error('Failed to compress image'))
            }
          },
          'image/webp',
          quality
        )
      }
      
      img.onerror = () => reject(new Error('Failed to load image'))
      img.src = URL.createObjectURL(file)
    } catch (error) {
      reject(error)
    }
  })
}

/**
 * Cache management for images
 */
export class ImageCache {
  private cache = new Map<string, string>()
  private maxSize = 50 // Maximum number of cached images
  
  /**
   * Get cached image URL
   * @param key - Cache key
   * @returns Cached URL or null
   */
  get(key: string): string | null {
    return this.cache.get(key) || null
  }
  
  /**
   * Set cached image URL
   * @param key - Cache key
   * @param url - Image URL
   */
  set(key: string, url: string): void {
    // Remove oldest entries if cache is full
    if (this.cache.size >= this.maxSize) {
      const firstKey = this.cache.keys().next().value
      this.cache.delete(firstKey || '')
    }
    
    this.cache.set(key, url)
  }
  
  /**
   * Clear cache
   */
  clear(): void {
    this.cache.clear()
  }
  
  /**
   * Get cache size
   */
  size(): number {
    return this.cache.size
  }
}

// Global image cache instance
export const imageCache = new ImageCache()
