<script setup lang="ts">
import { ConvertFromProperties, ConvertToProperties } from "wailsjs/go/services/yamlService.js";
import { ref } from "vue";
import { Codemirror } from "vue-codemirror";

const properties = ref("");
const yaml = ref("");
const doConvert = (isProperties: boolean, content: string) => {
  if (isProperties) {
    ConvertFromProperties(content).then((resp) => {
      console.log(resp);
      if (resp.success) {
        yaml.value = resp.data;
      }
    });
  } else {
    ConvertToProperties(content).then((resp) => {
      console.log(resp);
      if (resp.success) {
        properties.value = resp.data;
      }
    });
  }
  console.log(content);
};
</script>

<template>
  <n-space justify="space-around" size="large">
    <div>Properties</div>
    <div>Yaml</div>
  </n-space>
  <n-grid x-gap="12" :cols="2">
    <n-gi class="convertor-container">
      <codemirror
        v-model="properties"
        placeholder="properties"
        :style="{ height: '100%' }"
        :autofocus="true"
        :indent-with-tab="true"
        :tab-size="2"
        @change="doConvert(true, $event)"
      />
    </n-gi>
    <n-gi class="convertor-container">
      <codemirror
        v-model="yaml"
        placeholder="yaml"
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
.convertor-container {
  overflow: scroll;
}
</style>
