<script setup lang="ts">
import {
  ConvertFromJSON,
  ConvertToJSON,
} from "wailsjs/go/services/yamlService.js";
import { ref, computed } from "vue";
import { Codemirror } from "vue-codemirror";
import { json } from "@codemirror/lang-json";
import { yaml } from "@codemirror/lang-yaml";

const jsonContent = ref("");
const yamlContent = ref("");
const doConvert = (isProperties: boolean, content: string) => {
  if (isProperties) {
    ConvertFromJSON(content).then((resp) => {
      if (resp.success) {
        yamlContent.value = resp.data.yamlContent;
        jsonContent.value = resp.data.jsonContent;
      }
    });
  } else {
    ConvertToJSON(content).then((resp) => {
      if (resp.success) {
        yamlContent.value = resp.data.yamlContent;
        jsonContent.value = resp.data.jsonContent;
      }
    });
  }
};

const jsonExtensions = computed(() => {
  const result = [];
  result.push(json());
  return result;
});
const yamlExtensions = computed(() => {
  const result = [];
  result.push(yaml());
  return result;
});
</script>

<template>
  <n-space justify="space-around" size="large">
    <div>Properties</div>
    <div>Yaml</div>
  </n-space>
  <n-grid x-gap="12" :cols="2">
    <n-gi class="convertor-container">
      <codemirror
        v-model="jsonContent"
        placeholder="json"
        :extensions="jsonExtensions"
        :style="{ height: '100%' }"
        :autofocus="true"
        :indent-with-tab="true"
        :tab-size="2"
        @change="doConvert(true, $event)"
      />
    </n-gi>
    <n-gi class="convertor-container">
      <codemirror
        v-model="yamlContent"
        placeholder="yamlContent"
        :extensions="yamlExtensions"
        :style="{ height: '100%' }"
        :autofocus="true"
        :indent-with-tab="true"
        :tab-size="2"
        @change="doConvert(false, $event)"
      />
    </n-gi>
  </n-grid>
</template>

<style scoped>
.n-grid {
  height: 100%;
}
</style>
