<script setup>
const props = defineProps({
  inverse: {
    type: Boolean,
    default: false,
  },
  strokeWidth: {
    type: [Number, String],
    default: 3,
  },
});
</script>

<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    xmlns:xlink="http://www.w3.org/1999/xlink"
    viewBox="0 0 24 24"
  >
    <g
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <path
        d="M6 4h11a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H6a1 1 0 0 1-1-1V5a1 1 0 0 1 1-1m3 0v18"
      ></path>
      <path d="M13 8h2"></path>
      <path d="M13 12h2"></path>
    </g>
  </svg>
</template>

<style lang="scss" scoped></style>
