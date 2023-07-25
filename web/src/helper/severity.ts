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

export const severityIcons = (severity: number | undefined) => {
  switch (severity) {
    case 1:
      return '<i class="fa-solid fa-check"></i>';
    case 2:
      return `<i class="fa-solid fa-triangle-exclamation"></i>`;
    case 3:
      return `<i class="fa-solid fa-exclamation"></i>`;
    default:
      return `<span class="loading loading-spinner"></span>`;
  }
};
