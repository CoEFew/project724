#!/usr/bin/env node

/**
 * Image Optimization Script
 * 
 * This script converts PNG images to WebP format and creates multiple sizes
 * for responsive loading. It uses the sharp library for high-quality conversion.
 * 
 * Usage: node scripts/optimizeImages.js
 * 
 * Requirements:
 * - sharp: npm install sharp
 * - fs-extra: npm install fs-extra
 */

const fs = require('fs-extra')
const path = require('path')
const sharp = require('sharp')

// Configuration
const INPUT_DIR = path.join(__dirname, '../src/assets/images')
const OUTPUT_DIR = path.join(__dirname, '../src/assets/images/optimized')
const SIZES = [320, 640, 768, 1024, 1280, 1920]
const QUALITY = {
  webp: 80,
  png: 90,
  jpg: 85
}

/**
 * Get all PNG files in the input directory
 */
async function getImageFiles() {
  try {
    const files = await fs.readdir(INPUT_DIR)
    return files.filter(file => file.toLowerCase().endsWith('.png'))
  } catch (error) {
    console.error('Error reading input directory:', error)
    return []
  }
}

/**
 * Convert image to WebP format
 */
async function convertToWebP(inputPath, outputPath, quality = QUALITY.webp) {
  try {
    await sharp(inputPath)
      .webp({ quality })
      .toFile(outputPath)
    
    const stats = await fs.stat(outputPath)
    console.log(`✅ Created WebP: ${path.basename(outputPath)} (${formatBytes(stats.size)})`)
    return true
  } catch (error) {
    console.error(`❌ Error converting to WebP: ${inputPath}`, error.message)
    return false
  }
}

/**
 * Create responsive sizes for an image
 */
async function createResponsiveSizes(inputPath, baseName, format = 'webp') {
  const results = []
  
  for (const size of SIZES) {
    try {
      const outputPath = path.join(OUTPUT_DIR, `${baseName}_${size}w.${format}`)
      
      await sharp(inputPath)
        .resize(size, null, {
          withoutEnlargement: true,
          fit: 'inside'
        })
        .webp({ quality: QUALITY.webp })
        .toFile(outputPath)
      
      const stats = await fs.stat(outputPath)
      results.push({
        size,
        path: outputPath,
        sizeBytes: stats.size
      })
      
      console.log(`✅ Created ${size}w: ${path.basename(outputPath)} (${formatBytes(stats.size)})`)
    } catch (error) {
      console.error(`❌ Error creating ${size}w version:`, error.message)
    }
  }
  
  return results
}

/**
 * Create a thumbnail version
 */
async function createThumbnail(inputPath, baseName) {
  try {
    const outputPath = path.join(OUTPUT_DIR, `${baseName}_thumb.webp`)
    
    await sharp(inputPath)
      .resize(150, 150, {
        fit: 'cover',
        position: 'center'
      })
      .webp({ quality: 60 })
      .toFile(outputPath)
    
    const stats = await fs.stat(outputPath)
    console.log(`✅ Created thumbnail: ${path.basename(outputPath)} (${formatBytes(stats.size)})`)
    return outputPath
  } catch (error) {
    console.error(`❌ Error creating thumbnail:`, error.message)
    return null
  }
}

/**
 * Format bytes to human readable format
 */
function formatBytes(bytes, decimals = 2) {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}

/**
 * Get original file size
 */
async function getOriginalSize(filePath) {
  try {
    const stats = await fs.stat(filePath)
    return stats.size
  } catch (error) {
    return 0
  }
}

/**
 * Main optimization function
 */
async function optimizeImages() {
  console.log('🚀 Starting image optimization...\n')
  
  // Create output directory
  await fs.ensureDir(OUTPUT_DIR)
  
  // Get all PNG files
  const imageFiles = await getImageFiles()
  
  if (imageFiles.length === 0) {
    console.log('❌ No PNG files found in input directory')
    return
  }
  
  console.log(`📁 Found ${imageFiles.length} PNG files to optimize\n`)
  
  let totalOriginalSize = 0
  let totalOptimizedSize = 0
  let successCount = 0
  
  // Process each image
  for (const file of imageFiles) {
    const inputPath = path.join(INPUT_DIR, file)
    const baseName = path.parse(file).name
    const outputPath = path.join(OUTPUT_DIR, `${baseName}.webp`)
    
    console.log(`🔄 Processing: ${file}`)
    
    // Get original size
    const originalSize = await getOriginalSize(inputPath)
    totalOriginalSize += originalSize
    console.log(`   Original size: ${formatBytes(originalSize)}`)
    
    // Convert to WebP
    const webpSuccess = await convertToWebP(inputPath, outputPath)
    
    if (webpSuccess) {
      // Create responsive sizes
      await createResponsiveSizes(inputPath, baseName, 'webp')
      
      // Create thumbnail
      await createThumbnail(inputPath, baseName)
      
      // Get optimized size
      const optimizedSize = await getOriginalSize(outputPath)
      totalOptimizedSize += optimizedSize
      
      const savings = ((originalSize - optimizedSize) / originalSize * 100).toFixed(1)
      console.log(`   Optimized size: ${formatBytes(optimizedSize)} (${savings}% smaller)\n`)
      
      successCount++
    } else {
      console.log(`   ❌ Failed to optimize ${file}\n`)
    }
  }
  
  // Summary
  console.log('📊 Optimization Summary:')
  console.log(`   Files processed: ${successCount}/${imageFiles.length}`)
  console.log(`   Original total size: ${formatBytes(totalOriginalSize)}`)
  console.log(`   Optimized total size: ${formatBytes(totalOptimizedSize)}`)
  
  if (totalOriginalSize > 0) {
    const totalSavings = ((totalOriginalSize - totalOptimizedSize) / totalOriginalSize * 100).toFixed(1)
    console.log(`   Total savings: ${formatBytes(totalOriginalSize - totalOptimizedSize)} (${totalSavings}% smaller)`)
  }
  
  console.log('\n✅ Image optimization complete!')
  console.log(`📁 Optimized images saved to: ${OUTPUT_DIR}`)
}

/**
 * Create an index file for easy imports
 */
async function createIndexFile() {
  try {
    const files = await fs.readdir(OUTPUT_DIR)
    const webpFiles = files.filter(file => file.endsWith('.webp') && !file.includes('_'))
    
    const indexContent = `// Auto-generated image index
// Generated by optimizeImages.js

${webpFiles.map(file => {
  const name = path.parse(file).name
  return `export { default as ${name} } from './optimized/${file}'`
}).join('\n')}

// Responsive image utilities
export const getResponsiveSrc = (baseName, size) => \`/src/assets/images/optimized/\${baseName}_\${size}w.webp\`
export const getThumbnailSrc = (baseName) => \`/src/assets/images/optimized/\${baseName}_thumb.webp\`
`

    await fs.writeFile(path.join(INPUT_DIR, 'index.js'), indexContent)
    console.log('📝 Created index.js for easy imports')
  } catch (error) {
    console.error('❌ Error creating index file:', error.message)
  }
}

// Run optimization
if (require.main === module) {
  optimizeImages()
    .then(() => createIndexFile())
    .catch(error => {
      console.error('❌ Optimization failed:', error)
      process.exit(1)
    })
}

module.exports = {
  optimizeImages,
  createIndexFile,
  convertToWebP,
  createResponsiveSizes,
  createThumbnail
}
