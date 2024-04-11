<script setup lang="ts">
import { EncodeFromJWT } from "wailsjs/go/services/encodeService.js";
import { types } from "wailsjs/go/models";
import { ref, onMounted } from "vue";

const decodedContent = ref("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c");
const header = ref("");
const payload = ref("");
const doEncode = () => {
  const content = decodedContent.value
  EncodeFromJWT(content).then((resp: types.JSResp) => {
    console.log(resp)
    if (resp.success) {
      header.value = resp.data.header;
      payload.value = resp.data.payload;
    }
    else {
      header.value = resp.msg;
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
    <n-form-item :span="12" label="Header" path="decodedContent">
      <n-input
        v-model:value="header"
        placeholder="Textarea"
        type="textarea"
        :rows="8"
        :disabled="true"
      />
    </n-form-item>
    <n-form-item :span="12" label="Payload" path="payload">
      <n-input
        v-model:value="payload"
        placeholder="Textarea"
        type="textarea"
        :rows="8"
        :disabled="true"
      />
    </n-form-item>
  </n-form>
</template>

<style scoped></style>
