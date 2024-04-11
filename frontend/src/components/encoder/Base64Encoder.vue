<script setup lang="ts">
import {
  EncodeFromBase64,
  DecodeFromBase64,
} from "wailsjs/go/services/encodeService.js";
import {types} from 'wailsjs/go/models';
import { ref, computed, onMounted } from "vue";

const decodedContent = ref("hello world");
const encodedContent = ref("");
const doEncode = () => {
  const content = decodedContent.value
  EncodeFromBase64(content).then((resp: types.JSResp) => {
      if (resp.success) {
        encodedContent.value = resp.data;
      }
    });
};
const doDecode = () => {
  const content = encodedContent.value
  DecodeFromBase64(content).then((resp: types.JSResp) => {
      if (resp.success) {
        decodedContent.value = resp.data;
      }
    });
};
onMounted(() => {
  doEncode()
})
</script>

<template>
  <n-form label-placement="top">
    <n-form-item :span="12" label="Decoded" path="decodedContent">
      <n-input
        v-model:value="decodedContent"
        placeholder="Textarea"
        type="textarea"
        :rows="8"
        @update:value="doEncode"
      />
    </n-form-item>
    <n-form-item :span="12" label="Encoded" path="decodedContent">
      <n-input
        v-model:value="encodedContent"
        placeholder="Textarea"
        type="textarea"
        :rows="8"
        @update:value="doDecode"
      />
    </n-form-item>
  </n-form>
</template>

<style scoped>
</style>
