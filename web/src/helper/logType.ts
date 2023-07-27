import { database_LogType } from '@/openapi';

export const getIcon = (type: number | undefined) => {
  switch (type) {
    case database_LogType.LogGeneral:
      return `<i class="fa-solid fa-circle-nodes">`;
    case database_LogType.LogRestic:
      return `<i class="fa-solid fa-file-arrow-up">`;
    case database_LogType.LogCustom:
      return `<i class="fa-solid fa-terminal">`;
    case database_LogType.LogPrune:
      return `<i class="fa-solid fa-broom">`;
    case database_LogType.LogCheck:
      return `<i class="fa-solid fa-check">`;
    default:
      return `<i class="fa-solid fa-question">`;
  }
};
