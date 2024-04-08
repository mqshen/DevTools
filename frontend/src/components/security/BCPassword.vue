<script setup lang="ts">
import {
    GenerateBCPassword
} from "wailsjs/go/services/securityService.js";
import { ref, onMounted } from "vue";

const model = ref({
  content: "123456",
  password: null,
});

const generatePassword = () => {
    GenerateBCPassword(model.value.content).then((resp) => {
      console.log(resp);
      if (resp.success) {
        model.value.password = resp.data;
      }
    })
}
onMounted(() => {
  generatePassword();
});
</script>

<template>
  <n-form ref="formRef" :model="model" >
    <n-form-item path="content" label="content">
      <n-input v-model:value="model.content" @input="generatePassword"/>
    </n-form-item>
    <n-form-item path="password" label="password">
      <n-input
        v-model:value="model.password"
        placeholder="Textarea"
        type="textarea"
        :autosize="{
          minRows: 3,
          maxRows: 5,
        }"
      />
    </n-form-item>
  </n-form>
</template>

<style scoped>
.n-grid {
  height: 100%;
}
.convertor-container {
  overflow: scroll;
}
</style>
