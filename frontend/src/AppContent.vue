<script setup lang="ts">
import { computed, ref, reactive, onMounted } from "vue";
import { debounce } from "lodash";
import { useThemeVars } from "naive-ui";
import usePreferencesStore from "./stores/preferences";
import useTabStore from "./stores/tab";
import iconUrl from "@/assets/images/icon.png";
import { isMacOS } from "@/utils/platform";
import { extraTheme } from "@/utils/extra_theme";
import ContentValueTab from "@/components/content/ContentValueTab.vue";
import ToolbarControlWidget from "@/components/common/ToolbarControlWidget.vue";
import ResizeableWrapper from "@/components/common/ResizeableWrapper.vue";
import Ribbon from "@/components/sidebar/Ribbon.vue";
import DevContentPane from "./components/content/DevContentPane.vue";

import {
  WindowToggleMaximise,
  EventsOn,
  WindowIsFullscreen,
  WindowIsMaximised,
} from "wailsjs/runtime/runtime.js";

const themeVars = useThemeVars();
const prefStore = usePreferencesStore();
const tabStore = useTabStore();
const exThemeVars = computed(() => {
  return extraTheme(prefStore.isDark);
});
const saveSidebarWidth = debounce(prefStore.savePreferences, 1000, {
  trailing: true,
});
const handleResize = () => {
  saveSidebarWidth();
};
const logoPaddingLeft = ref(10);
const maximised = ref(false);

const props = defineProps({
  loading: Boolean,
});

const data = reactive({
  navMenuWidth: 50,
  toolbarHeight: 38,
});

const logoWrapperWidth = computed(() => {
  return `${prefStore.behavior.asideWidth - 4}px`;
});

const hideRadius = ref(false);
const wrapperStyle = computed(() => {
  return hideRadius.value
    ? {}
    : {
        border: `1px solid ${themeVars.value.borderColor}`,
        borderRadius: "10px",
      };
});
const spinStyle = computed(() => {
  return hideRadius.value
    ? {
        backgroundColor: themeVars.value.bodyColor,
      }
    : {
        backgroundColor: themeVars.value.bodyColor,
        borderRadius: "10px",
      };
});
const onToggleFullscreen = (fullscreen: boolean) => {
  hideRadius.value = fullscreen;
  if (fullscreen) {
    logoPaddingLeft.value = 10;
  } else {
    logoPaddingLeft.value = isMacOS() ? 70 : 10;
  }
};

const onToggleMaximize = (isMaximised: boolean) => {
  if (isMaximised) {
    maximised.value = true;
    if (!isMacOS()) {
      hideRadius.value = true;
    }
  } else {
    maximised.value = false;
    if (!isMacOS()) {
      hideRadius.value = false;
    }
  }
};

EventsOn("window_changed", (info: any) => {
  const { fullscreen, maximised } = info;
  onToggleFullscreen(fullscreen === true);
  onToggleMaximize(maximised);
});

onMounted(async () => {
  const fullscreen = await WindowIsFullscreen();
  onToggleFullscreen(fullscreen === true);
  const maximised = await WindowIsMaximised();
  onToggleMaximize(maximised);
});
</script>

<template>
  <n-spin
    :show="props.loading"
    :style="spinStyle"
    :theme-overrides="{ opacitySpinning: 0 }"
  >
    <div id="app-content-wrapper" :style="wrapperStyle" class="flex-box-v">
      <div
        id="app-toolbar"
        :style="{ height: data.toolbarHeight + 'px' }"
        class="flex-box-h"
        style="--wails-draggable: drag"
        @dblclick="WindowToggleMaximise"
      >
        <!-- title -->
        <div
          id="app-toolbar-title"
          :style="{
            width: logoWrapperWidth,
            minWidth: logoWrapperWidth,
            paddingLeft: `${logoPaddingLeft}px`,
          }"
        >
          <n-space :size="3" :wrap="false" :wrap-item="false" align="center">
            <n-avatar
              :size="32"
              :src="iconUrl"
              color="#0000"
              style="min-width: 32px"
            />
            <div style="min-width: 68px; font-weight: 800">DevTools</div>
          </n-space>
        </div>
        <!-- browser tabs -->
        <div class="app-toolbar-tab flex-item-expand">
          <content-value-tab />
        </div>
        <div class="flex-item-expand"></div>
        <!-- simulate window control buttons -->
        <toolbar-control-widget
          v-if="!isMacOS()"
          :maximised="maximised"
          :size="data.toolbarHeight"
          style="align-self: flex-start"
        />
      </div>

      <div
        id="app-content"
        :style="prefStore.generalFont"
        class="flex-box-h flex-item-expand"
        style="--wails-draggable: none"
      >
        <resizeable-wrapper
          v-model:size="prefStore.behavior.asideWidth"
          :min-size="200"
          :offset="data.navMenuWidth"
          class="flex-item dev-side-bar"
          @update:size="handleResize"
        >
          <ribbon />
        </resizeable-wrapper>
        <dev-content-pane
          v-for="t in tabStore.tabs"
          v-show="tabStore.currentTabName === t.name"
          :componentType="t.compontent"
          :key="t.name"
          :server="t.name"
          class="flex-item-expand"
        />
      </div>
    </div>
  </n-spin>
</template>

<style lang="scss" scoped>
#app-content-wrapper {
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  box-sizing: border-box;
  background-color: v-bind("themeVars.bodyColor");
  color: v-bind("themeVars.textColorBase");

  #app-toolbar {
    background-color: v-bind("exThemeVars.titleColor");
    border-bottom: 1px solid v-bind("exThemeVars.splitColor");
    &-title {
      padding-left: 10px;
      padding-right: 10px;
      box-sizing: border-box;
      align-self: center;
      align-items: baseline;
    }
  }

  .app-toolbar-tab {
    align-self: flex-end;
    margin-bottom: -1px;
    margin-left: 3px;
  }

  #app-content {
    height: calc(100% - 60px);

    .content-area {
      overflow: hidden;
    }
  }

  .app-side {
    //overflow: hidden;
    height: 100%;
    background-color: v-bind("exThemeVars.sidebarColor");
    border-right: 1px solid v-bind("exThemeVars.splitColor");
  }
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.dev-side-bar {
  background-color: v-bind("exThemeVars.titleColor");
}
</style>
