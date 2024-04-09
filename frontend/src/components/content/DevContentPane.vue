<script setup lang="ts">
import { componentTypes } from "@/consts/support_component_type";
import { useThemeVars } from "naive-ui";
import PropertiesConvert from "@/components/convert/PropertiesConvert.vue";
import JSONConvert from "@/components/convert/JSONConvert.vue";
import DateConvert from "@/components/convert/DateConvert.vue";
import BCPassword from "@/components/security/BCPassword.vue";
import HashGenerator from "@/components/generators/HashGenerator.vue";
import IPComponent from "@/components/network/IPComponent.vue";
const themeVars = useThemeVars();

const props = defineProps({
    componentType: String,
})

const valueComponents = {
  [componentTypes.PROPERTIES]: PropertiesConvert,
  [componentTypes.JSON]: JSONConvert,
  [componentTypes.Date]: DateConvert,
  [componentTypes.BCryptPassword]: BCPassword,
  [componentTypes.Hash]: HashGenerator,
  [componentTypes.IP]: IPComponent,
};
</script>

<template>
  <div :class="'content-container flex-box-v '+ componentType + '_container'">
    <component :is="valueComponents[componentType]" />
  </div>
</template>

<style lang="scss">
.content-sub-tab {
  background-color: v-bind("themeVars.tabColor");
  height: 100%;
}

.content-sub-tab-pane {
  padding: 0 !important;
  height: 100%;
  background-color: v-bind("themeVars.bodyColor");
  overflow: hidden;
}

.n-tabs .n-tabs-bar {
  transition: none !important;
}
.content-container {
  padding: 20px;
}
.IP_container.content-container {
  padding: 20px;
  overflow-y: auto;
}
</style>
