<script setup lang="ts">
import PageContent from '@/components/ui/PageContent.vue';
import PageHeader from '@/components/ui/PageHeader.vue';
import SystemOverview from '@/components/system/SystemOverview.vue';
import { useSystemStore } from '@/stores/system';
import { computed } from 'vue';

const systemStore = useSystemStore();
const badges = computed(() => {
  return {
    'next backup': systemStore.nextBackup,
    'next cleanup': systemStore.nextCleanup,
    'next check': systemStore.nextCheck,
    'healthcheck url': systemStore.system.config.healthcheck_url,
  };
});
</script>

<template>
  <div>
    <PageHeader>
      <SystemOverview :versions="systemStore.system.versions" :badges="badges" />
    </PageHeader>
    <PageContent> </PageContent>
  </div>
</template>
