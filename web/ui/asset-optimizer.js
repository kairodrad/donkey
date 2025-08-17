// Asset optimization utility
// This provides guidance for optimizing the themed assets

export const ASSET_OPTIMIZATION_SUGGESTIONS = {
  'donkey-background-light.png': {
    currentFormat: 'PNG',
    suggestedFormat: 'WebP',
    optimizations: [
      'Convert to WebP format for ~30% size reduction',
      'Use progressive encoding for faster loading',
      'Consider SVG format since this is vector-style art'
    ]
  },
  'donkey-background-dark.png': {
    currentFormat: 'PNG',
    suggestedFormat: 'WebP',
    optimizations: [
      'Convert to WebP format for ~30% size reduction',
      'Use progressive encoding for faster loading',
      'Consider SVG format since this is vector-style art'
    ]
  },
  'donkey-title-light.png': {
    currentFormat: 'PNG',
    suggestedFormat: 'WebP',
    optimizations: [
      'Convert to WebP format for ~30% size reduction',
      'Consider SVG format for text-based graphics',
      'Use smaller resolution if displaying at fixed size'
    ]
  },
  'donkey-title-dark.png': {
    currentFormat: 'PNG',
    suggestedFormat: 'WebP',
    optimizations: [
      'Convert to WebP format for ~30% size reduction',
      'Consider SVG format for text-based graphics',
      'Use smaller resolution if displaying at fixed size'
    ]
  }
};

// Enhanced image loading with format detection
export function createOptimizedImage(assetName, theme, fallbackTheme = 'light') {
  const img = new Image();
  
  // Try WebP first if supported
  const supportsWebP = (() => {
    const canvas = document.createElement('canvas');
    return canvas.toDataURL('image/webp').indexOf('data:image/webp') === 0;
  })();
  
  const extension = supportsWebP ? 'webp' : 'png';
  const primarySrc = `/assets/${assetName}-${theme}.${extension}`;
  const fallbackSrc = `/assets/${assetName}-${fallbackTheme}.png`;
  
  img.src = primarySrc;
  img.onerror = () => {
    if (img.src !== fallbackSrc) {
      img.src = fallbackSrc;
    }
  };
  
  return img;
}

// Lazy loading for background images
export function setupLazyBackgroundImage(element, assetName, theme) {
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const img = createOptimizedImage(assetName, theme);
        img.onload = () => {
          element.style.backgroundImage = `url(${img.src})`;
        };
        observer.unobserve(element);
      }
    });
  });
  
  observer.observe(element);
}

// Performance monitoring for asset loading
export function trackAssetPerformance(assetName, startTime) {
  const endTime = performance.now();
  const loadTime = endTime - startTime;
  
  console.log(`Asset ${assetName} loaded in ${loadTime.toFixed(2)}ms`);
  
  // Store performance data for optimization insights
  if (!window.assetPerformance) window.assetPerformance = {};
  window.assetPerformance[assetName] = loadTime;
}