<template>
  <div class="flex items-center justify-between rounded-lg p-4 w-full">
    <!-- Left: Text Skeleton -->
    <div class="flex flex-col gap-2">
      <div class="skeleton h-4 w-30"></div>
      <div class="skeleton h-3 w-20"></div>
    </div>

    <div class="flex">
      <div v-for="i in 3" :key="i" class="flex flex-col gap-1">
        <div class="flex items-center">
          <div class="skeleton rounded-full w-9 h-9"></div>
          <div v-if="i !== 3" class="skeleton h-1 w-6"></div>
        </div>
        <div class="flex pl-[0.1rem]">
          <div class="skeleton w-7 h-3"></div>
        </div>
      </div>
    </div>
  </div>
</template>
