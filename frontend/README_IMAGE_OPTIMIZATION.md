# Image Optimization Guide

This document explains the image optimization system implemented in the PETTEXT project to significantly improve page load times while maintaining visual quality.

## 🚀 Performance Improvements

### Before Optimization
- **Total image size**: ~25MB (17 PNG files)
- **Largest files**: 2.4MB each
- **Format**: PNG only
- **Loading**: All images loaded immediately

### After Optimization
- **Total image size**: ~3-5MB (estimated 80-85% reduction)
- **Format**: WebP with PNG fallback
- **Loading**: Lazy loading with responsive sizes
- **Caching**: Intelligent browser caching

## 📁 File Structure

```
frontend/src/
├── components/
│   ├── LoadingOverlay.vue          # Consistent loading component
│   └── OptimizedImage.vue          # High-performance image component
├── utils/
│   └── imageOptimization.ts        # Image optimization utilities
├── assets/images/
│   ├── optimized/                  # Generated optimized images
│   │   ├── catwalk.webp
│   │   ├── catwalk_320w.webp
│   │   ├── catwalk_640w.webp
│   │   └── ...
│   └── index.js                    # Auto-generated image index
└── scripts/
    └── optimizeImages.js           # Image conversion script
```

## 🛠️ Setup Instructions

### 1. Install Dependencies

```bash
cd frontend
npm run optimize-images:install
```

### 2. Optimize Images

```bash
npm run optimize-images
```

This will:
- Convert all PNG images to WebP format
- Create responsive sizes (320w, 640w, 768w, 1024w, 1280w, 1920w)
- Generate thumbnails for quick loading
- Create an index file for easy imports

### 3. Use Optimized Images

#### Basic Usage
```vue
<template>
  <OptimizedImage 
    src="/src/assets/images/catwalk.webp" 
    alt="Loading cat" 
    class="w-24 h-24"
  />
</template>
```

#### Responsive Images
```vue
<template>
  <OptimizedImage 
    src="/src/assets/images/catwalk.webp" 
    alt="Loading cat"
    :sizes="[320, 640, 1024]"
    loading="lazy"
    class="w-full h-64"
  />
</template>
```

#### With Loading States
```vue
<template>
  <OptimizedImage 
    src="/src/assets/images/catwalk.webp" 
    alt="Loading cat"
    :show-placeholder="true"
    :show-spinner="true"
    fallback="/src/assets/images/fallback.png"
  />
</template>
```

## 🎯 Key Features

### 1. Automatic Format Detection
- Detects WebP support in browser
- Falls back to PNG/JPG if WebP not supported
- Maintains quality while reducing file size

### 2. Responsive Images
- Generates multiple sizes for different screen resolutions
- Uses `srcset` and `sizes` attributes for optimal loading
- Reduces bandwidth usage on mobile devices

### 3. Lazy Loading
- Images load only when they enter the viewport
- Uses Intersection Observer API
- Fallback for older browsers

### 4. Loading States
- Placeholder while image loads
- Loading spinner for better UX
- Error handling with fallback images

### 5. Caching
- Intelligent browser caching
- Memory-efficient image management
- Automatic cleanup of unused images

## 🔧 Configuration

### Image Quality Settings
```typescript
const QUALITY_SETTINGS = {
  high: 90,      // For hero images
  medium: 75,    // For regular content
  low: 60,       // For thumbnails
  thumbnail: 40  // For small previews
}
```

### Responsive Breakpoints
```typescript
const BREAKPOINTS = {
  sm: '640px',   // Mobile
  md: '768px',   // Tablet
  lg: '1024px',  // Desktop
  xl: '1280px',  // Large desktop
  '2xl': '1536px' // Extra large
}
```

## 📊 Performance Monitoring

### Before/After Comparison
```bash
# Check original file sizes
ls -la src/assets/images/*.png

# Check optimized file sizes
ls -la src/assets/images/optimized/*.webp
```

### Browser DevTools
1. Open Network tab
2. Filter by "Img"
3. Check file sizes and load times
4. Compare with/without optimization

## 🚨 Troubleshooting

### Images Not Loading
1. Check if WebP files exist in `optimized/` folder
2. Run `npm run optimize-images` to generate missing files
3. Check browser console for errors

### Large File Sizes
1. Adjust quality settings in `optimizeImages.js`
2. Reduce number of responsive sizes
3. Check if images are being compressed properly

### Lazy Loading Issues
1. Ensure `IntersectionObserver` is supported
2. Check if images have `data-src` attribute
3. Verify viewport settings

## 🔄 Maintenance

### Adding New Images
1. Place PNG files in `src/assets/images/`
2. Run `npm run optimize-images`
3. Use `OptimizedImage` component in templates

### Updating Existing Images
1. Replace PNG files
2. Run `npm run optimize-images`
3. Clear browser cache to see changes

### Regular Cleanup
- Remove unused optimized images
- Update quality settings based on performance needs
- Monitor file sizes and loading times

## 📈 Expected Results

- **80-85% reduction** in total image file size
- **3-5x faster** initial page load
- **Better mobile experience** with responsive images
- **Improved SEO** with faster loading times
- **Consistent loading experience** across all pages

## 🎨 Loading Visual

All pages now use the consistent `LoadingOverlay` component with:
- Animated cat loading visual
- Network status indicators
- Custom loading messages
- Smooth transitions
- Error handling

This ensures a cohesive user experience across the entire application.
