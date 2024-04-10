<script setup lang="ts">
import { UploadFileInfo } from "naive-ui";
const fileChnage = (options: { file: UploadFileInfo, fileList: Array<UploadFileInfo>, event?: Event }) => {
    console.log(options.file)
}
</script>
<template>
  <n-grid x-gap="12" :cols="3">
    <n-grid-item span="1">
      <n-upload
        multiple
        directory-dnd
        @change="fileChnage"
        :max="5"
      >
        <n-upload-dragger>
          <div style="margin-bottom: 12px">
            <n-icon size="48" :depth="3">
              <archive-icon />
            </n-icon>
          </div>
          <n-text style="font-size: 16px">
            Click or drag a file to this area to upload
          </n-text>
          <n-p depth="3" style="margin: 8px 0 0 0">
            Strictly prohibit from uploading sensitive information. For example,
            your bank card PIN or your credit card expiry date.
          </n-p>
        </n-upload-dragger>
      </n-upload>
    </n-grid-item>
    <n-grid-item span="2"> </n-grid-item>
  </n-grid>
</template>
