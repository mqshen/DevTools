<script setup lang="ts">
// import { GenerateHash } from "wailsjs/go/services/dateService.js";
import { isYesterday } from "date-fns/esm";
import { onMounted, ref } from "vue";

const model = ref({
  date: Date.now(),
  uninxTime: 0,
  iso8601: "",
  calendar: 0,
});


const calcDate = (currentDate: Date) => {
  model.value.date = currentDate.getTime()
  model.value.uninxTime = Math.floor(model.value.date / 1000);
  model.value.iso8601 = currentDate.toISOString()
  model.value.calendar = model.value.date
}

const generateFromDate = () => {
  const currentDate = new Date(model.value.date)
  calcDate(currentDate);
};

const generateFromTimestamp = () => {
  const currentDate = new Date();
  currentDate.setTime(model.value.uninxTime * 1000)
  calcDate(currentDate);
};

const generateFromISO = () => {
  const currentDate = new Date(model.value.iso8601)
  calcDate(currentDate);
};

const handleUpdateValue = (
  _: number,
  { year, month, date }: { year: number; month: number; date: number }
) => {
  const currentDate = new Date(model.value.date)
  currentDate.setFullYear(year)
  currentDate.setMonth(month)
  currentDate.setDate(date)
  calcDate(currentDate);
};

const isDateDisabled = (timestamp: number) => {
  if (isYesterday(timestamp)) {
    return true;
  }
  return false;
};

onMounted(() => {
  generateFromDate();
});
</script>

<template>
  <n-form ref="formRef" :model="model">
    <n-form-item path="date" label="Date">
      <n-date-picker v-model:value="model.date" type="datetime" :actions="['now']" :update-value-on-close="true" @update:value="generateFromDate"/>
    </n-form-item>
    <n-form-item path="uninxTime" label="Unix Time">
      <n-input-number v-model:value="model.uninxTime" @update:value="generateFromTimestamp"/>
    </n-form-item>
    <n-form-item path="iso8601" label="ISO 8601">
      <n-input v-model:value="model.iso8601" @input="generateFromISO"/>
    </n-form-item>
    <n-form-item path="calendar" label="Calendar" >
      <n-calendar style="height: 300px;"
        v-model:value="model.calendar"
        #="{ year, month, date }"
        :is-date-disabled="isDateDisabled"
        @update:value="handleUpdateValue"
      >
      </n-calendar>
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
