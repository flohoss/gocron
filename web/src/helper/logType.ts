export const getIcon = (type: number | undefined) => {
  switch (type) {
    case 1:
      return `<i class="fa-solid fa-copy">`;
    default:
      return `<i class="fa-solid fa-question">`;
  }
};
