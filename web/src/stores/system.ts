import { ref } from "vue";
import { defineStore } from "pinia";
import { SystemService, type system_Data } from "@/openapi";

const emptySystem: system_Data = {
  configuration: {
    hostname: "",
    rclone_config_file: "",
  },
  disk: {
    total: 0,
    unit: "",
    used: 0,
  },
  versions: {
    compose: "",
    docker: "",
    go: "",
    gobackup: "",
    rclone: "",
    restic: "",
  },
};

export const useSystemStore = defineStore("system", () => {
  const system = ref<system_Data>(emptySystem);

  async function getSystem() {
    const response = await SystemService.getSystem();
    system.value = response;
  }

  return { system, getSystem };
});
