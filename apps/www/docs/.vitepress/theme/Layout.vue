<script setup>
import DefaultTheme from 'vitepress/theme'
import NavBarVersion from './components/NavBarVersion.vue'

const { Layout } = DefaultTheme
</script>

<template>
  <Layout>
    <template #nav-bar-title-after>
      <NavBarVersion />
    </template>
  </Layout>
</template>
