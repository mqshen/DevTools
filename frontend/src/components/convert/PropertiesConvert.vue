<script setup lang="ts">
import {
  ConvertFromProperties,
  ConvertToProperties,
} from "wailsjs/go/services/yamlService.js";
import { types } from "wailsjs/go/models";
import { ref, computed, onMounted } from "vue";
import { Codemirror } from "vue-codemirror";
import { yaml } from "@codemirror/lang-yaml";

const properties = ref(`spring.datasource.url=abc
spring.datasource.username=alice
spring.datasource.password=123456
spring.datasource.driver-class-name=com.mysql.cj.jdbc.Driver
spring.datasource.type=com.zaxxer.hikari.HikariDataSource`);
const yamlContent = ref("");
const doConvert = (isProperties: boolean) => {
  if (isProperties) {
    const content = properties.value
    ConvertFromProperties(content).then((resp: types.JSResp) => {
      if (resp.success) {
        yamlContent.value = resp.data;
      }
    });
  } else {
    const content = yamlContent.value
    ConvertToProperties(content).then((resp: types.JSResp) => {
      if (resp.success) {
        properties.value = resp.data;
      }
    });
  }
};

const yamlExtensions = computed(() => {
  const result = [];
  result.push(yaml());
  return result;
});

onMounted(() => {
  doConvert(true)
})

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
        v-model="yamlContent"
        placeholder="yaml"
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
