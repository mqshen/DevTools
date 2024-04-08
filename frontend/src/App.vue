<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { darkTheme } from 'naive-ui'
import { darkThemeOverrides, themeOverrides } from '@/utils/theme'
import AppContent from './AppContent.vue'
import usePreferencesStore from "stores/preferences";
const prefStore = usePreferencesStore();

const initializing = ref(true)
onMounted(async () => {
    try {
        initializing.value = true
    } finally {
        initializing.value = false
    }
})
</script>

<template>
  <n-config-provider
    :inline-theme-disabled="true"
    :locale="prefStore.themeLocale"
    :theme="prefStore.isDark ? darkTheme : undefined"
    :theme-overrides="prefStore.isDark ? darkThemeOverrides : themeOverrides"
    class="fill-height"
  >
    <n-dialog-provider>
      <app-content :loading="initializing" />

      <!-- top modal dialogs -->
    </n-dialog-provider>
  </n-config-provider>
</template>

<style scoped>
</style>
