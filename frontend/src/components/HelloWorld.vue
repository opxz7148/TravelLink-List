<script setup lang="ts">
defineProps<{
  msg: string
}>()
</script>

<template>
  <div class="text-center lg:text-left">
    <h1 class="font-medium text-4xl -mt-2.5 text-green-600">{{ msg }}</h1>
    <h3 class="text-xl">
      You've successfully created a project with
      <a href="https://vite.dev/" target="_blank" rel="noopener" class="text-blue-600 hover:text-blue-700">Vite</a> +
      <a href="https://vuejs.org/" target="_blank" rel="noopener" class="text-blue-600 hover:text-blue-700">Vue 3</a>.
    </h3>
  </div>
</template>
