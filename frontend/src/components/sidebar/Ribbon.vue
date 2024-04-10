<script setup lang="ts">
import { componentTypes } from "@/consts/support_component_type";
import { ref, h } from "vue";
import { NIcon } from "naive-ui";
import Properties from "@/components/icons/Properties.vue";
import useTabStore from "stores/tab";

const tabStore = useTabStore();
const data = ref([
  {
    key: "0",
    name: "Converters",
    type: 0,
    children: [
      { key: "00", type: 1, name: "Properties<>Yaml", compontent: componentTypes.PROPERTIES, icon: "filetype-yml" },
      { key: "01", type: 1, name: "JSON<>Yaml", compontent: componentTypes.JSON, icon: "filetype-json"},
      { key: "02", type: 1, name: "Date", compontent: componentTypes.Date, icon: "calendar2-week" },
    ],
  },
  {
    key: "1",
    name: "Network",
    type: 0,
    children: [
      { key: "10", type: 1, name: "IP", compontent: componentTypes.IP, icon: "modem"},
    ],
  },,
  {
    key: "2",
    name: "Security",
    type: 0,
    children: [
      { key: "20", type: 1, name: "BCryptPassword", compontent: componentTypes.BCryptPassword, icon: "shield-lock" },
    ],
  },
  {
    key: "3",
    name: "Generators",
    type: 0,
    children: [
      { key: "30", type: 1, name: "Hash", compontent: componentTypes.Hash, icon: "hash" },
    ],
  },
  {
    key: "4",
    name: "Encoders/Decoders",
    type: 0,
    children: [
      { key: "40", type: 1, name: "HTML", compontent: componentTypes.HTML, icon: "filetype-html"},
      { key: "41", type: 1, name: "URL", compontent: componentTypes.URL, icon: "file-earmark" },
      { key: "42", type: 1, name: "Base64", compontent: componentTypes.Base64, icon: "feather" },
      { key: "43", type: 1, name: "JWT", compontent: componentTypes.JWT, icon: "fingerprint" },
    ],
  }
]);

// render tree item label
const renderLabel = ({ option }) => {
  if (option.type == 0) {
    return h("span", { class: "classify-title" }, option.name);
  }
  return option.name;
};

const renderPrefix = ({ option }) => {
  if (option.type == 1) {
    return h("i", { class: "bi bi-" + option.icon });
  }
};
const renderSwitcherIcon = ({ option }) => {
  return null;
};
const defaultExpandedKeys = ref(["0", "1", "2", "3", "4"]);

const handleClick = ({ option }) => {
  if (option.type == 1) {
    tabStore.upsertTab(option.name, option.compontent, option.icon);
  }
};
</script>

<template>
  <div id="app-ribbon" class="flex-box-v">
    <n-tree
      block-line
      :data="data"
      :default-expanded-keys="defaultExpandedKeys"
      :render-label="renderLabel"
      :render-prefix="renderPrefix"
      :render-switcher-icon="renderSwitcherIcon"
      :override-default-node-click-behavior="handleClick"
      expand-on-click
    />
  </div>
</template>

<style lang="scss">
#app-ribbon {
  width: 100%;
  .n-tree-node-switcher--hide {
    display: none;
  }
  .classify-title {
    color: gray;
    margin-left: 15px;
  }
  .n-tree-node-switcher {
    display: none;
  }
}
</style>
