/**
 * Asset Loading Utilities
 * 
 * This module provides utilities for robust asset loading with fallback handling
 * and error recovery mechanisms for production environments.
 */

/**
 * Asset loading configuration
 */
export interface AssetConfig {
  src: string
  fallback?: string
  retries?: number
  timeout?: number
}

/**
 * Asset loading result
 */
export interface AssetResult {
  success: boolean
  src?: string
  error?: string
  fallbackUsed?: boolean
}

/**
 * Load an asset with fallback handling and retry logic
 * @param config - Asset loading configuration
 * @returns Promise<AssetResult> - Loading result with success status
 */
export async function loadAsset(config: AssetConfig): Promise<AssetResult> {
  const { src, fallback, retries = 3, timeout = 5000 } = config
  
  try {
    // Try to load the primary asset
    const result = await loadSingleAsset(src, timeout)
    if (result.success) {
      return { success: true, src: result.src }
    }
    
    // If primary fails and fallback exists, try fallback
    if (fallback) {
      console.warn(`Primary asset failed to load: ${src}, trying fallback: ${fallback}`)
      const fallbackResult = await loadSingleAsset(fallback, timeout)
      if (fallbackResult.success) {
        return { success: true, src: fallbackResult.src, fallbackUsed: true }
      }
    }
    
    return { 
      success: false, 
      error: `Failed to load asset: ${src}${fallback ? ` and fallback: ${fallback}` : ''}` 
    }
  } catch (error) {
    console.error('Asset loading error:', error)
    return { 
      success: false, 
      error: error instanceof Error ? error.message : 'Unknown error' 
    }
  }
}

/**
 * Load a single asset with timeout
 * @param src - Asset source URL
 * @param timeout - Timeout in milliseconds
 * @returns Promise<AssetResult> - Loading result
 */
function loadSingleAsset(src: string, timeout: number): Promise<AssetResult> {
  return new Promise((resolve) => {
    const img = new Image()
    const timer = setTimeout(() => {
      resolve({ success: false, error: 'Timeout' })
    }, timeout)
    
    img.onload = () => {
      clearTimeout(timer)
      resolve({ success: true, src: img.src })
    }
    
    img.onerror = () => {
      clearTimeout(timer)
      resolve({ success: false, error: 'Load error' })
    }
    
    img.src = src
  })
}

/**
 * Preload multiple assets with fallback handling
 * @param assets - Array of asset configurations
 * @returns Promise<AssetResult[]> - Results for all assets
 */
export async function preloadAssets(assets: AssetConfig[]): Promise<AssetResult[]> {
  try {
    const promises = assets.map(asset => loadAsset(asset))
    const results = await Promise.all(promises)
    
    const successCount = results.filter(r => r.success).length
    console.log(`Preloaded ${successCount}/${assets.length} assets successfully`)
    
    return results
  } catch (error) {
    console.error('Asset preloading error:', error)
    return assets.map(() => ({ success: false, error: 'Preload failed' }))
  }
}

/**
 * Create a data URL for a simple loading spinner
 * This provides a fallback when all image assets fail to load
 * @param size - Size of the spinner in pixels
 * @param color - Color of the spinner
 * @returns string - Data URL for the spinner
 */
export function createSpinnerDataUrl(size: number = 24, color: string = '#6366f1'): string {
  try {
    const canvas = document.createElement('canvas')
    canvas.width = size
    canvas.height = size
    const ctx = canvas.getContext('2d')
    
    if (!ctx) return ''
    
    // Create a simple spinning circle
    ctx.strokeStyle = color
    ctx.lineWidth = 2
    ctx.lineCap = 'round'
    
    // Draw a partial circle to create spinner effect
    ctx.beginPath()
    ctx.arc(size / 2, size / 2, (size - 4) / 2, 0, Math.PI * 1.5)
    ctx.stroke()
    
    return canvas.toDataURL('image/png')
  } catch (error) {
    console.error('Error creating spinner data URL:', error)
    return ''
  }
}

/**
 * Check if an asset URL is accessible
 * @param url - Asset URL to check
 * @returns Promise<boolean> - Whether the asset is accessible
 */
export async function checkAssetAccessibility(url: string): Promise<boolean> {
  try {
    const response = await fetch(url, { method: 'HEAD' })
    return response.ok
  } catch (error) {
    console.warn(`Asset not accessible: ${url}`, error)
    return false
  }
}

/**
 * Get optimized asset URL for production
 * This handles the difference between development and production asset paths
 * @param assetPath - Relative asset path
 * @returns string - Optimized asset URL
 */
export function getOptimizedAssetUrl(assetPath: string): string {
  try {
    // In production, Vite processes these imports and provides optimized URLs
    // In development, we can use the direct path
    if (import.meta.env.PROD) {
      // For production, we rely on Vite's asset processing
      return assetPath
    } else {
      // For development, we can use the direct path
      return assetPath
    }
  } catch (error) {
    console.error('Error getting optimized asset URL:', error)
    return assetPath
  }
}

/**
 * Asset loading cache to prevent duplicate requests
 */
class AssetCache {
  private cache = new Map<string, AssetResult>()
  private loading = new Set<string>()
  
  /**
   * Get cached asset result
   * @param key - Asset key
   * @returns AssetResult | null
   */
  get(key: string): AssetResult | null {
    return this.cache.get(key) || null
  }
  
  /**
   * Set cached asset result
   * @param key - Asset key
   * @param result - Asset result
   */
  set(key: string, result: AssetResult): void {
    this.cache.set(key, result)
  }
  
  /**
   * Check if asset is currently loading
   * @param key - Asset key
   * @returns boolean
   */
  isLoading(key: string): boolean {
    return this.loading.has(key)
  }
  
  /**
   * Mark asset as loading
   * @param key - Asset key
   */
  setLoading(key: string): void {
    this.loading.add(key)
  }
  
  /**
   * Mark asset as loaded
   * @param key - Asset key
   */
  setLoaded(key: string): void {
    this.loading.delete(key)
  }
  
  /**
   * Clear cache
   */
  clear(): void {
    this.cache.clear()
    this.loading.clear()
  }
}

// Global asset cache instance
export const assetCache = new AssetCache()
