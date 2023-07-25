export const getIcon = (type: number | undefined) => {
  switch (type) {
    case 1:
      return `<i class="fa-solid fa-circle-nodes">`;
    case 2:
      return `<i class="fa-solid fa-file-arrow-up">`;
    case 3:
      return `<i class="fa-solid fa-terminal">`;
    case 4:
      return `<i class="fa-solid fa-broom">`;
    case 5:
      return `<i class="fa-solid fa-check">`;
    default:
      return `<i class="fa-solid fa-question">`;
  }
};
