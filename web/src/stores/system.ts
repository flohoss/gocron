import { computed, ref } from 'vue';
import { defineStore } from 'pinia';
import { SystemService, type system_SystemConfig } from '@/openapi';
import { emptySystemConfig } from '@/types';
import parser from 'cron-parser';
import moment from 'moment';

export const useSystemStore = defineStore('system', () => {
  const system = ref<system_SystemConfig>(emptySystemConfig);

  async function fetchSystem() {
    system.value = await SystemService.getSystem();
  }

  function getNextCronDate(date: string | undefined) {
    return moment(
      parser
        .parseExpression(date || '', { tz: system.value.config.time_zone })
        .next()
        .toDate()
    ).format('llll');
  }

  const nextBackup = computed(() => getNextCronDate(system.value.config.backup_cron));
  const nextCleanup = computed(() => getNextCronDate(system.value.config.cleanup_cron));
  const nextCheck = computed(() => getNextCronDate(system.value.config.check_cron));

  return { system, fetchSystem, nextBackup, nextCleanup, nextCheck };
});
