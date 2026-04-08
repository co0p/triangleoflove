<template>
  <div :class="['gallery-root', activeTheme]">
    <div class="container">

      <!-- Header -->
      <header class="gallery-header section">
        <div class="gallery-header-top">
          <img :src="logoSrc" alt="Triangle of Love mark" class="gallery-logo" />
          <div>
            <h1>Design Gallery</h1>
            <p class="text-muted">All tokens and component classes — 375 px reference.</p>
          </div>
        </div>
        <div class="gallery-theme-row">
          <span class="text-sm gallery-theme-label">Theme:</span>
          <button
            v-for="t in themes"
            :key="t.value"
            :class="['gallery-theme-btn', activeTheme === t.value ? 'gallery-theme-btn--active' : '']"
            @click="activeTheme = t.value"
          >{{ t.label }}</button>
        </div>
      </header>

      <!-- 1. Color Swatches -->
      <section class="section">
        <h2 class="gallery-section-title">Color Swatches</h2>

        <p class="gallery-section-sub">Brand palette</p>
        <div class="gallery-swatches">
          <div
            v-for="s in brandSwatches"
            :key="s.name"
            class="gallery-swatch"
            :style="{ background: s.hex }"
          >
            <span :style="{ color: s.light ? 'rgba(255,255,255,0.9)' : 'rgba(30,27,24,0.85)' }">
              {{ s.name }}<br /><small>{{ s.hex }}</small>
            </span>
          </div>
        </div>

        <p class="gallery-section-sub" style="margin-top: var(--space-4)">Neutrals</p>
        <div class="gallery-swatches gallery-swatches--tight">
          <div
            v-for="s in neutralSwatches"
            :key="s.name"
            class="gallery-swatch gallery-swatch--sm"
            :style="{ background: s.hex, border: '1px solid #e8e4dc' }"
          >
            <span :style="{ color: s.light ? 'rgba(255,255,255,0.9)' : 'rgba(30,27,24,0.85)' }">
              <small>{{ s.name }}</small>
            </span>
          </div>
        </div>

        <p class="gallery-section-sub" style="margin-top: var(--space-4)">
          Semantic tokens — switch theme to see them update
        </p>
        <div class="gallery-swatches gallery-swatches--tight">
          <div
            v-for="s in semanticSwatches"
            :key="s.name"
            class="gallery-swatch gallery-swatch--sm"
            :style="{ background: `var(${s.token})`, border: '1px solid var(--color-border)' }"
          >
            <span style="color: rgba(30,27,24,0.8)">
              <small>{{ s.name }}</small>
            </span>
          </div>
        </div>
      </section>

      <!-- 2. Typography -->
      <section class="section">
        <h2 class="gallery-section-title">Typography</h2>
        <div class="card gallery-type-card">
          <div class="gallery-type-row">
            <h1>Heading 1</h1>
            <code class="gallery-spec">3xl / bold / 1.25lh</code>
          </div>
          <div class="gallery-type-row">
            <h2>Heading 2</h2>
            <code class="gallery-spec">2xl / semibold</code>
          </div>
          <div class="gallery-type-row">
            <h3>Heading 3</h3>
            <code class="gallery-spec">xl / semibold</code>
          </div>
          <div class="gallery-type-row">
            <p>Body — long-form copy, 1.625 line-height.</p>
            <code class="gallery-spec">base / normal</code>
          </div>
          <div class="gallery-type-row">
            <p class="text-muted">Muted — helper text, hints, captions.</p>
            <code class="gallery-spec">sm / muted</code>
          </div>
          <div class="gallery-type-row">
            <p class="text-sm">Small — labels, metadata.</p>
            <code class="gallery-spec">sm</code>
          </div>
        </div>
      </section>

      <!-- 3. Buttons -->
      <section class="section">
        <h2 class="gallery-section-title">Buttons</h2>
        <div class="gallery-variants">
          <div v-for="v in buttonVariants" :key="v.label" class="gallery-variant">
            <span class="gallery-variant-label">{{ v.label }}</span>
            <button :class="['btn', v.cls]" :disabled="v.disabled">{{ v.text }}</button>
          </div>
        </div>
      </section>

      <!-- 4. Inputs -->
      <section class="section">
        <h2 class="gallery-section-title">Inputs</h2>
        <div class="gallery-variants">

          <div class="gallery-variant">
            <span class="gallery-variant-label">Default</span>
            <div class="input-group">
              <label class="input-label">Email</label>
              <input class="input" type="email" placeholder="you@example.com" />
            </div>
          </div>

          <div class="gallery-variant">
            <span class="gallery-variant-label">Focus (simulated)</span>
            <div class="input-group">
              <label class="input-label">Email</label>
              <input
                class="input gallery-state-focus"
                type="email"
                value="you@example.com"
                readonly
              />
            </div>
          </div>

          <div class="gallery-variant">
            <span class="gallery-variant-label">Error</span>
            <div class="input-group">
              <label class="input-label">Email</label>
              <input class="input input--error" type="email" value="not-an-email" readonly />
              <span class="input-hint input-hint--error">Enter a valid email address.</span>
            </div>
          </div>

          <div class="gallery-variant">
            <span class="gallery-variant-label">Disabled</span>
            <div class="input-group">
              <label class="input-label">Email</label>
              <input class="input" type="email" disabled value="you@example.com" />
            </div>
          </div>

        </div>
      </section>

      <!-- 5. Error Alert -->
      <section class="section">
        <h2 class="gallery-section-title">Error Alert</h2>
        <div class="alert-error" role="alert">
          Invalid email or password. Please try again.
        </div>
      </section>

      <!-- 6. Card -->
      <section class="section">
        <h2 class="gallery-section-title">Card</h2>
        <div class="card">
          <h3>Card title</h3>
          <p class="text-muted" style="margin-top: var(--space-2); margin-bottom: var(--space-4)">
            White surface with border, elevation shadow, and large border-radius.
            Contains any content.
          </p>
          <button class="btn btn-primary">Primary action</button>
        </div>
      </section>

      <!-- 7. Logo Mark -->
      <section class="section">
        <h2 class="gallery-section-title">Logo Mark</h2>
        <div class="gallery-logo-sizes">
          <div v-for="size in [80, 48, 32]" :key="size" class="gallery-logo-item">
            <img :src="logoSrc" :width="size" :height="size" alt="Triangle of Love" />
            <span class="text-muted">{{ size }}px</span>
          </div>
        </div>
      </section>

      <!-- 8. Avatar Badge -->
      <section class="section">
        <h2 class="gallery-section-title">Avatar Badge</h2>
        <div class="gallery-row">
          <div class="avatar">JG</div>
          <div class="avatar">AB</div>
          <div class="avatar">TL</div>
        </div>
      </section>

      <!-- 9. Focus Ring -->
      <section class="section">
        <h2 class="gallery-section-title">Focus Ring</h2>
        <p class="text-muted" style="margin-bottom: var(--space-4)">
          A warm gold ring appears on all interactive elements during keyboard navigation
          (<kbd>Tab</kbd> to experience it live on the right button).
        </p>
        <div class="gallery-row">
          <button class="btn btn-secondary gallery-demo-focused" style="width: auto">
            Focused (simulated)
          </button>
          <button class="btn btn-secondary" style="width: auto">Tab → me</button>
        </div>
      </section>

      <!-- 10. NavBar classes preview -->
      <section class="section">
        <h2 class="gallery-section-title">NavBar Classes</h2>
        <p class="text-muted" style="margin-bottom: var(--space-4)">
          Static preview of <code>.navbar</code>, <code>.navbar-brand</code>, and
          <code>.navbar-greeting</code>. The <code>NavBar.vue</code> component is Deliverable 3.
        </p>
        <div class="navbar gallery-navbar-preview">
          <div class="navbar-brand">
            <img :src="logoSrc" alt="" width="32" height="32" />
          </div>
          <span class="navbar-greeting">Hello, Alex</span>
          <div class="avatar">AK</div>
        </div>
      </section>

    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import logoSrc from '../assets/logo.svg';

const activeTheme = ref('');

const themes = [
  { value: '',           label: 'Brand' },
  { value: 'theme-warm', label: 'Warm'  },
];

const brandSwatches = [
  { name: 'Sage 500',  hex: '#7aab7a', light: false },
  { name: 'Sage 700',  hex: '#508a50', light: true  },
  { name: 'Rose 400',  hex: '#d97d8a', light: false },
  { name: 'Rose 700',  hex: '#b05568', light: true  },
  { name: 'Gold 500',  hex: '#d4a843', light: false },
  { name: 'Gold 700',  hex: '#a8832e', light: true  },
];

const neutralSwatches = [
  { name: '50',  hex: '#faf9f7', light: false },
  { name: '100', hex: '#f2f0ec', light: false },
  { name: '200', hex: '#e8e4dc', light: false },
  { name: '400', hex: '#b8b0a5', light: false },
  { name: '600', hex: '#6e665c', light: true  },
  { name: '800', hex: '#3a3530', light: true  },
  { name: '900', hex: '#1e1b18', light: true  },
];

const semanticSwatches = [
  { name: 'bg',        token: '--color-bg'        },
  { name: 'surface',   token: '--color-surface'   },
  { name: 'primary',   token: '--color-primary'   },
  { name: 'accent',    token: '--color-accent'    },
  { name: 'highlight', token: '--color-highlight' },
  { name: 'error',     token: '--color-error'     },
];

const buttonVariants = [
  { label: 'Primary',             cls: 'btn-primary',   text: 'Continue', disabled: false },
  { label: 'Primary — disabled',  cls: 'btn-primary',   text: 'Continue', disabled: true  },
  { label: 'Secondary',           cls: 'btn-secondary', text: 'Sign up',  disabled: false },
  { label: 'Secondary — disabled',cls: 'btn-secondary', text: 'Sign up',  disabled: true  },
  { label: 'Ghost',               cls: 'btn-ghost',     text: 'Cancel',   disabled: false },
  { label: 'Ghost — disabled',    cls: 'btn-ghost',     text: 'Cancel',   disabled: true  },
];
</script>

<style scoped>
.gallery-root {
  min-height: 100vh;
  background-color: var(--color-bg);
}

/* ---- Header ---- */
.gallery-header {
  border-bottom: 1px solid var(--color-border);
}

.gallery-header-top {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  margin-bottom: var(--space-4);
}

.gallery-logo {
  width: 56px;
  height: 56px;
  flex-shrink: 0;
}

.gallery-theme-row {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  flex-wrap: wrap;
}

.gallery-theme-label {
  color: var(--color-text-muted);
}

.gallery-theme-btn {
  padding: var(--space-1) var(--space-3);
  border-radius: var(--radius-full);
  background: var(--color-neutral-100);
  border: 1.5px solid var(--color-border);
  font-family: var(--font-sans);
  font-size: var(--font-size-sm);
  color: var(--color-text);
  cursor: pointer;
  transition: background 0.15s ease, border-color 0.15s ease, color 0.15s ease;
}

.gallery-theme-btn--active {
  background: var(--color-primary);
  border-color: var(--color-primary);
  color: var(--color-text-inverse);
}

/* ---- Section titles ---- */
.gallery-section-title {
  margin-bottom: var(--space-4);
  padding-bottom: var(--space-2);
  border-bottom: 2px solid var(--color-border);
}

.gallery-section-sub {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  margin-bottom: var(--space-2);
}

/* ---- Swatches ---- */
.gallery-swatches {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-2);
}

.gallery-swatches--tight {
  grid-template-columns: repeat(4, 1fr);
}

.gallery-swatch {
  height: 72px;
  border-radius: var(--radius-md);
  padding: var(--space-2);
  display: flex;
  align-items: flex-end;
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
}

.gallery-swatch--sm {
  height: 48px;
  align-items: center;
}

/* ---- Typography card ---- */
.gallery-type-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.gallery-type-row {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: var(--space-2);
}

.gallery-spec {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  white-space: nowrap;
  flex-shrink: 0;
}

/* ---- Component variant rows ---- */
.gallery-variants {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.gallery-variant {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.gallery-variant-label {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-weight: var(--font-weight-semibold);
}

/* ---- Input demo states ---- */
/* Simulates the :focus appearance for static gallery display */
.gallery-state-focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(122, 171, 122, 0.22);
}

/* ---- Logo sizes ---- */
.gallery-logo-sizes {
  display: flex;
  align-items: flex-end;
  gap: var(--space-6);
  flex-wrap: wrap;
}

.gallery-logo-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-1);
}

/* ---- Row layout ---- */
.gallery-row {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  flex-wrap: wrap;
}

/* ---- Focus ring demo ---- */
/* Simulates the :focus-visible ring for gallery display */
.gallery-demo-focused {
  box-shadow: var(--focus-ring);
  outline: none;
}

/* ---- NavBar preview ---- */
.gallery-navbar-preview {
  border-radius: var(--radius-md);
  overflow: hidden;
  border: 1px solid var(--color-border);
}
</style>
