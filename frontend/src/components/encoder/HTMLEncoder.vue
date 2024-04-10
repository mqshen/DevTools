<script setup lang="ts">
import {
    EncodeFromHTML,
    DecodeFromHTML,
} from "wailsjs/go/services/encodeService.js";
import {types} from 'wailsjs/go/models';
import { ref, computed } from "vue";

const decodedContent = ref("");
const encodedContent = ref("");
const doEncode = (content: string) => {
    EncodeFromHTML(content).then((resp: types.JSResp) => {
      if (resp.success) {
        encodedContent.value = resp.data;
      }
    });
};
const doDecode = (content: string) => {
    DecodeFromHTML(content).then((resp: types.JSResp) => {
      if (resp.success) {
        decodedContent.value = resp.data;
      }
    });
};
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
