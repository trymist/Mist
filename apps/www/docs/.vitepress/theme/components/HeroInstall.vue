<template>
  <div class="hero-install-section">
    <div class="install-box">
      <div class="install-code" @click="copyToClipboard">
        <span class="install-prompt">$</span>
        <span class="install-command">curl -fsSL https://trymist.cloud/install.sh | bash</span>
        <button class="copy-button" :class="{ copied }" @click.stop="copyToClipboard">
          <span v-if="!copied">📋</span>
          <span v-else>✓</span>
        </button>
      </div>
    </div>
    <div class="hero-actions">
      <a href="/guide/getting-started" class="action-button primary">Get Started</a>
      <a href="https://github.com/trymist/mist" class="action-button secondary" target="_blank">View on GitHub</a>
      <a href="https://discord.gg/hr6TCQDDkj" class="action-button discord" target="_blank">
        <svg width="20" height="20" viewBox="0 0 71 55" fill="none" xmlns="http://www.w3.org/2000/svg">
          <g clip-path="url(#clip0)">
            <path d="M60.1045 4.8978C55.5792 2.8214 50.7265 1.2916 45.6527 0.41542C45.5603 0.39851 45.468 0.440769 45.4204 0.525289C44.7963 1.6353 44.105 3.0834 43.6209 4.2216C38.1637 3.4046 32.7345 3.4046 27.3892 4.2216C26.905 3.0581 26.1886 1.6353 25.5617 0.525289C25.5141 0.443589 25.4218 0.40133 25.3294 0.41542C20.2584 1.2888 15.4057 2.8186 10.8776 4.8978C10.8384 4.9147 10.8048 4.9429 10.7825 4.9795C1.57795 18.7309 -0.943561 32.1443 0.293408 45.3914C0.299005 45.4562 0.335386 45.5182 0.385761 45.5576C6.45866 50.0174 12.3413 52.7249 18.1147 54.5195C18.2071 54.5477 18.305 54.5139 18.3638 54.4378C19.7295 52.5728 20.9469 50.6063 21.9907 48.5383C22.0523 48.4172 21.9935 48.2735 21.8676 48.2256C19.9366 47.4931 18.0979 46.6 16.3292 45.5858C16.1893 45.5041 16.1781 45.304 16.3068 45.2082C16.679 44.9293 17.0513 44.6391 17.4067 44.3461C17.471 44.2926 17.5606 44.2813 17.6362 44.3151C29.2558 49.6202 41.8354 49.6202 53.3179 44.3151C53.3935 44.2785 53.4831 44.2898 53.5502 44.3433C53.9057 44.6363 54.2779 44.9293 54.6529 45.2082C54.7816 45.304 54.7732 45.5041 54.6333 45.5858C52.8646 46.6197 51.0259 47.4931 49.0921 48.2228C48.9662 48.2707 48.9102 48.4172 48.9718 48.5383C50.038 50.6034 51.2554 52.5699 52.5959 54.435C52.6519 54.5139 52.7526 54.5477 52.845 54.5195C58.6464 52.7249 64.529 50.0174 70.6019 45.5576C70.6551 45.5182 70.6887 45.459 70.6943 45.3942C72.1747 30.0791 68.2147 16.7757 60.1968 4.9823C60.1772 4.9429 60.1437 4.9147 60.1045 4.8978ZM23.7259 37.3253C20.2276 37.3253 17.3451 34.1136 17.3451 30.1693C17.3451 26.225 20.1717 23.0133 23.7259 23.0133C27.308 23.0133 30.1626 26.2532 30.1066 30.1693C30.1066 34.1136 27.28 37.3253 23.7259 37.3253ZM47.3178 37.3253C43.8196 37.3253 40.9371 34.1136 40.9371 30.1693C40.9371 26.225 43.7636 23.0133 47.3178 23.0133C50.9 23.0133 53.7545 26.2532 53.6986 30.1693C53.6986 34.1136 50.9 37.3253 47.3178 37.3253Z" fill="currentColor"/>
          </g>
          <defs>
            <clipPath id="clip0">
              <rect width="71" height="55" fill="white"/>
            </clipPath>
          </defs>
        </svg>
        Join Discord
      </a>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const copied = ref(false)

const copyToClipboard = () => {
  const command = 'curl -sSL https://trymist.cloud/install.sh | bash'
  navigator.clipboard.writeText(command).then(() => {
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  })
}
</script>

<style scoped>
.hero-install-section {
  max-width: 800px;
  margin: 2rem auto 0;
  padding: 0 1.5rem;
}

.install-box {
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-border);
  border-radius: 12px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.install-box:hover {
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.install-code {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 1rem;
  cursor: pointer;
  user-select: none;
  position: relative;
}

.install-prompt {
  color: var(--vp-c-brand-1);
  font-weight: bold;
  opacity: 0.8;
}

.install-command {
  flex: 1;
  color: var(--vp-c-text-1);
  font-weight: 500;
}

.copy-button {
  background: var(--vp-c-brand-1);
  color: white;
  border: none;
  border-radius: 6px;
  padding: 0.5rem 0.75rem;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 48px;
}

.copy-button:hover {
  background: var(--vp-c-brand-2);
  transform: scale(1.05);
}

.copy-button.copied {
  background: #10b981;
}

.hero-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
}

.action-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 1rem;
  text-decoration: none;
  transition: all 0.2s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.action-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.action-button.primary {
  background: var(--vp-c-brand-1);
  color: white;
}

.action-button.primary:hover {
  background: var(--vp-c-brand-2);
}

.action-button.secondary {
  background: var(--vp-c-bg-soft);
  color: var(--vp-c-text-1);
  border: 1px solid var(--vp-c-border);
}

.action-button.secondary:hover {
  border-color: var(--vp-c-brand-1);
  color: var(--vp-c-brand-1);
}

.action-button.discord {
  background: #5865F2;
  color: white;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.action-button.discord:hover {
  background: #4752C4;
}

.action-button.discord svg {
  width: 20px;
  height: 20px;
}

@media (max-width: 768px) {
  .install-code {
    font-size: 0.875rem;
    flex-wrap: wrap;
  }
  
  .copy-button {
    position: absolute;
    right: 0;
    top: 0;
  }
  
  .install-command {
    padding-right: 60px;
  }
}
</style>
