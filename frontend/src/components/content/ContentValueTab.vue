<script setup lang="ts">
import { componentTypes } from "@/consts/support_component_type";
import Server from "@/components/icons/Server.vue";
import useTabStore from "stores/tab";
import { computed, shallowRef } from "vue";
import { useI18n } from "vue-i18n";
import { get, map } from "lodash";
import { useThemeVars } from "naive-ui";
import { extraTheme } from "@/utils/extra_theme";
import usePreferencesStore from "stores/preferences";
import Properties from "@/components/icons/Properties.vue";

const iconComponents = {
  [componentTypes.PROPERTIES]: Properties,
};
/**
 * Value content tab on head
 */
const themeVars = useThemeVars();
const i18n = useI18n();
const tabStore = useTabStore();
const prefStore = usePreferencesStore();

const onCloseTab = (tabIndex) => {
  const tab = get(tabStore.tabs, tabIndex);
  if (tab != null) {
    $dialog.warning(
      i18n.t("dialogue.close_confirm", { name: tab.name }),
      () => {
        tabStore.closeTab(tab.name)
      }
    );
  }
};



const tabClass = (idx: number) => {
  if (tabStore.activatedIndex === idx) {
    return [
      "value-tab",
      "value-tab-active"
    ];
  } else if (tabStore.activatedIndex - 1 === idx) {
    return ["value-tab", "value-tab-inactive"];
  } else {
    return ["value-tab", "value-tab-inactive", "value-tab-inactive2"];
  }
};

const tab = computed(() =>
  tabStore.tabs.map((item) => {
    return {
      name: item.name,
      icon: item.compontent,
    };
  })
);

const exThemeVars = computed(() => {
  return extraTheme(prefStore.isDark);
});
</script>

<template>
  <n-tabs
    v-model:value="tabStore.activatedIndex"
    :closable="true"
    :tabs-padding="3"
    :theme-overrides="{
      tabFontWeightActive: 800,
      tabGapSmallCard: 0,
      tabGapMediumCard: 0,
      tabGapLargeCard: 0,
      tabColor: '#0000',
      tabBorderColor: '#0000',
      tabTextColorCard: themeVars.closeIconColor,
    }"
    size="small"
    type="card"
    @close="onCloseTab"
  >
    <n-tab
      v-for="(t, i) in tab"
      :key="i"
      :class="tabClass(i)"
      :closable="true"
      :name="i"
      @dblclick.stop="() => {}"
    >
      <n-space
        :size="5"
        :wrap-item="false"
        align="center"
        inline
        justify="center"
      >
        <n-icon size="18">
          <component :is="iconComponents[t.icon]" />
        </n-icon>
        <n-ellipsis style="max-width: 150px">{{ t.name }}</n-ellipsis>
      </n-space>
    </n-tab>
  </n-tabs>
</template>

<style lang="scss">
.value-tab {
  --wails-draggable: none;
  position: relative;
  border: 1px solid v-bind("exThemeVars.splitColor") !important;
}

.value-tab-active {
  background-color: v-bind("themeVars.tabColor") !important;
  border-bottom-color: v-bind("themeVars.tabColor") !important;

}

.value-tab-inactive {
  background-color: #0000 !important;
  border-color: #0000 !important;

  &:hover {
    background-color: v-bind("exThemeVars.splitColor") !important;
  }
}

.value-tab-inactive2 {
  &:after {
    content: "";
    position: absolute;
    top: 25%;
    height: 50%;
    width: 1px;
    background-color: v-bind("themeVars.borderColor");
    right: -2px;
  }

  &:hover::after {
    background-color: #0000;
  }
}
</style>
