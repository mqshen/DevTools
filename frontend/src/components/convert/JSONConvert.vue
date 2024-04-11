<script setup lang="ts">
import {
  ConvertFromJSON,
  ConvertToJSON,
} from "wailsjs/go/services/yamlService.js";
import { types } from "wailsjs/go/models";
import { ref, computed, onMounted } from "vue";
import { Codemirror } from "vue-codemirror";
import { json } from "@codemirror/lang-json";
import { yaml } from "@codemirror/lang-yaml";

const jsonContent = ref(`{
    "glossary": {
        "title": "example glossary",
		"GlossDiv": {
            "title": "S",
			"GlossList": {
                "GlossEntry": {
                    "ID": "SGML",
					"SortAs": "SGML",
					"GlossTerm": "Standard Generalized Markup Language",
					"Acronym": "SGML",
					"Abbrev": "ISO 8879:1986",
					"GlossDef": {
                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
						"GlossSeeAlso": ["GML", "XML"]
                    },
					"GlossSee": "markup"
                }
            }
        }
    }
}`);
const yamlContent = ref("");
const doConvert = (isJson: boolean) => {
  if (isJson) {
    const content = jsonContent.value
    ConvertFromJSON(content).then((resp: types.JSResp) => {
      if (resp.success) {
        yamlContent.value = resp.data.yamlContent;
        jsonContent.value = resp.data.jsonContent;
      }
    });
  } else {
    const content = yamlContent.value
    ConvertToJSON(content).then((resp: types.JSResp) => {
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
