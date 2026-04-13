<template>
  <div class="mt-8 flex lg:mt-0 lg:py-4 lg:pl-4 lg:relative">
    <i class="flex items-center justify-center w-8 h-8 lg:absolute lg:top-1/2 lg:-translate-y-1/2 lg:-left-6 lg:w-12 lg:h-12 lg:border lg:border-gray-300 lg:bg-white lg:rounded-lg text-gray-800">
      <slot name="icon"></slot>
    </i>
    <div class="flex-1 lg:ml-4">
      <h3 class="text-xl font-medium mb-1 text-gray-900">
        <slot name="heading"></slot>
      </h3>
      <slot></slot>
    </div>

    <!-- Vertical connector lines for desktop -->
    <div class="hidden lg:block absolute left-0 border-l border-gray-300" style="bottom: 50%; height: calc(50% - 25px)"></div>
    <div class="hidden lg:block absolute left-0 border-l border-gray-300" style="top: 50%; height: calc(50% - 25px)"></div>
  </div>
</template>

<script setup lang="ts">
defineSlots<{
  icon(): unknown;
  heading(): unknown;
  default(): unknown;
}>();
</script>
