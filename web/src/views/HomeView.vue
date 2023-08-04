<script setup lang="ts">
import PageContent from '@/components/ui/PageContent.vue';
import PageHeader from '@/components/ui/PageHeader.vue';
import SystemOverview from '@/components/system/SystemOverview.vue';
import { useSystemStore } from '@/stores/system';
import NextCountdown from '@/components/system/NextCountdown.vue';

const systemStore = useSystemStore();
</script>

<template>
  <div>
    <PageHeader>
      <SystemOverview :versions="systemStore.system.versions" />
    </PageHeader>
    <PageContent>
      <div class="grid place-items-center grid-cols-1 xl:grid-cols-3 gap-3 md:gap-5">
        <NextCountdown title="Next backup" :cron="systemStore.nextBackup" />
        <NextCountdown title="Next cleanup" :cron="systemStore.nextCleanup" />
        <NextCountdown title="Next check" :cron="systemStore.nextCheck" />
      </div>
    </PageContent>
  </div>
</template>
