<script setup lang="ts">
import { GenerateHash } from "wailsjs/go/services/hashService.js";
import { onMounted, ref } from "vue";

const model = ref({
  input: "Hello World",
  md5: null,
  sha1: null,
  sha256: null,
  sha512: null,
});

const generateHash = () => {
  GenerateHash(model.value.input).then((resp) => {
    if (resp.success) {
      const { md5, sha1, sha256, sha512 } = resp.data;
      model.value.md5 = md5;
      model.value.sha1 = sha1;
      model.value.sha256 = sha256;
      model.value.sha512 = sha512;
    }
  });
};

onMounted(() => {
  generateHash();
});
</script>

<template>
  <n-form ref="formRef" :model="model">
    <n-form-item path="input" label="Input">
      <n-input
        v-model:value="model.input"
        placeholder="Textarea"
        @input="generateHash"
        type="textarea"
        :autosize="{
          minRows: 3,
          maxRows: 5,
        }"
      />
    </n-form-item>
    <n-form-item path="MD5" label="MD5">
      <n-input v-model:value="model.md5" />
    </n-form-item>
    <n-form-item path="SHA1" label="SHA1">
      <n-input v-model:value="model.sha1" />
    </n-form-item>
    <n-form-item path="SHA256" label="SHA256">
      <n-input v-model:value="model.sha256" />
    </n-form-item>
    <n-form-item path="SHA512" label="SHA512">
      <n-input v-model:value="model.sha512" />
    </n-form-item>
  </n-form>
</template>

<style scoped>
.n-grid {
  height: 100%;
}
</style>
