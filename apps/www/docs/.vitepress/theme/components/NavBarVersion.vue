<template>
  <div class="navbar-version">
    <div v-if="loading" class="version-loading">
      <span class="loading-spinner">⟳</span>
    </div>
    <div v-else-if="error" class="version-error" :title="error">
      <span class="version-label">Version</span>
      <span class="version-text">—</span>
    </div>
    <a 
      v-else 
      :href="releaseData?.html_url" 
      target="_blank" 
      class="version-link"
      :title="`Latest release: ${version} (${formatDate(releaseData?.published_at)})`"
    >
      <span class="version-label">v</span>
      <span class="version-text">{{ version }}</span>
    </a>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const version = ref('...')
const loading = ref(true)
const error = ref(null)
const releaseData = ref(null)

const GITHUB_API = 'https://api.github.com/repos/trymist/mist/releases/latest'

const fetchLatestVersion = async () => {
  try {
    loading.value = true
    error.value = null
    
    const response = await fetch(GITHUB_API, {
      headers: {
        'Accept': 'application/vnd.github.v3+json'
      }
    })
    
    if (!response.ok) {
      throw new Error(`Failed to fetch version: ${response.status}`)
    }
    
    const data = await response.json()
    releaseData.value = data
    
    // Extract version from tag_name (e.g., "v1.0.0" -> "1.0.0")
    version.value = data.tag_name.replace(/^v/, '')
    
  } catch (err) {
    console.error('Error fetching version:', err)
    error.value = 'Unable to fetch version'
    version.value = '—'
  } finally {
    loading.value = false
  }
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

onMounted(() => {
  fetchLatestVersion()
})
</script>

<style scoped>
.navbar-version {
  display: inline-flex;
  align-items: center;
  margin: 0;
  height: 100%;
}

.version-loading,
.version-error,
.version-link {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.375rem 0.75rem;
  border-radius: 6px;
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-border);
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.2s ease;
}

.version-link {
  text-decoration: none;
  color: var(--vp-c-text-1);
  cursor: pointer;
}

.version-link:hover {
  background: var(--vp-c-bg-elv);
  border-color: var(--vp-c-brand-1);
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.version-link:hover .version-label {
  color: var(--vp-c-brand-1);
}

.version-link:hover .version-text {
  color: var(--vp-c-brand-1);
}

.version-label {
  color: var(--vp-c-text-2);
  font-weight: 600;
  font-size: 0.8rem;
}

.version-text {
  color: var(--vp-c-text-1);
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.875rem;
  font-weight: 600;
}

.loading-spinner {
  display: inline-block;
  animation: spin 1s linear infinite;
  font-size: 1rem;
  color: var(--vp-c-text-3);
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.version-error {
  color: var(--vp-c-text-3);
  cursor: help;
}

@media (max-width: 960px) {
  .navbar-version {
    margin: 0 0.5rem;
  }
  
  .version-loading,
  .version-error,
  .version-link {
    padding: 0.25rem 0.5rem;
    font-size: 0.8rem;
  }
  
  .version-label {
    font-size: 0.75rem;
  }
  
  .version-text {
    font-size: 0.8rem;
  }
}

@media (max-width: 768px) {
  .navbar-version {
    display: none;
  }
}
</style>
