import DefaultTheme from 'vitepress/theme'
import './custom.css'
import HeroInstall from './components/HeroInstall.vue'
import NavBarVersion from './components/NavBarVersion.vue'
import Layout from './Layout.vue'

export default {
  extends: DefaultTheme,
  Layout,
  enhanceApp({ app }) {
    app.component('HeroInstall', HeroInstall)
    app.component('NavBarVersion', NavBarVersion)
  }
}
