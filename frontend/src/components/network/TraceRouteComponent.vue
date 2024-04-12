<script setup lang="ts">
import { GetTraceRoute } from "wailsjs/go/services/traceRouteService.js";
import { ref, h } from "vue";
import { types } from "wailsjs/go/models";
import { useI18n } from "vue-i18n";
const i18n = useI18n();
const host = ref("");

type RowData = {
  ttl: number;
  host: string;
  elapsedTime: number;
  replied: boolean;
  isPrivate: boolean;
};

const columns = [
  {
    title: "TTL",
    key: "ttl",
  },
  {
    title: "Host",
    key: "host",
    render(row: RowData) {
      if (row.isPrivate) {
        return [
          h("span", row.host),
          h(
            "span",
            { class: "private-address" },
            "(" + i18n.t("network.private") + ")"
          ),
        ];
      }
      return row.host;
    },
  },
  {
    title: "ElapsedTime",
    key: "elapsedTime",

    render(row: RowData) {
      return row.elapsedTime + "ms";
    },
  },
  {
    title: "Replied",
    key: "replied",
    render(row: RowData) {
      return "" + row.replied;
    },
  },
];

const data = ref([]);
const loading = ref(false);

const doTrace = async () => {
  loading.value = true;
  const resp = await GetTraceRoute(host.value); //.then((resp: types.JSResp) => {
  if (resp.msg == "noPermission") {
    $dialog.warning( i18n.t("network.no_permission" ));
  }
  if (resp.success) {
    const { complete, routes } = resp.data;
    data.value = routes || [];
    if (!complete) {
      setTimeout(doTrace, 2000, this);
    } else {
      loading.value = false;
    }
  }
};
</script>

<template>
  <n-form
    class="flex-item"
    label-align="left"
    label-placement="left"
    label-width="auto"
    size="small"
  >
    <n-form-item label="Domain or IP">
      <n-space :wrap="false" :wrap-item="false" style="width: 100%">
        <n-input v-model:value="host" clearable placeholder="" />
        <n-button type="primary" @click="doTrace"> Trace </n-button>
      </n-space>
    </n-form-item>
  </n-form>
  <div id="">
    <n-data-table :columns="columns" :data="data"> </n-data-table>
    <div class="loading-container">
      <n-spin :show="loading" :theme-overrides="{ opacitySpinning: 0 }"
        ><span></span
      ></n-spin>
    </div>
  </div>
</template>

<style lang="scss">
.loading-container {
  position: relative;
  margin-top: 40px;
}
.private-address {
  color: gray;
}
</style>
