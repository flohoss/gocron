<script setup lang="ts">
import PageContent from '@/components/ui/PageContent.vue';
import PageHeader from '@/components/ui/PageHeader.vue';
import SystemOverview from '@/components/system/SystemOverview.vue';
</script>

<template>
  <div>
    <PageHeader>
      <SystemOverview />
    </PageHeader>
    <PageContent> </PageContent>
  </div>
</template>
