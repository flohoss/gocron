export const severityColor = (severity: number | undefined) => {
  switch (severity) {
    case 2:
      return 'text-warning';
    case 3:
      return 'text-error';
    default:
      return '';
  }
};

export const severityHTML = (severity: number | undefined) => {
  switch (severity) {
    case 1:
      return `<div class="text-success"><i class="fa-solid fa-check"></div>`;
    case 2:
      return `<div class="text-warning"><i class="fa-solid fa-triangle-exclamation"></i></div>`;
    case 2:
      return `<div class="text-error"><i class="fa-solid fa-exclamation"></i></div>`;
    default:
      return `<i class="fa-solid fa-question">`;
  }
};
