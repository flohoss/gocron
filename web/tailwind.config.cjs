/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    container: {
      center: true
    },
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: [
      {
        light: {
          ...require("daisyui/src/theming/themes")["lofi"],
          "success": "#28a745",
          "success-content": "white",
          "error": "#dc3545",
          "error-content": "white",
          "warning": "#ffc107",
          "warning-content": "black",
          "info": "#17a2b8",
          "info-content": "white",
        },
      },
      {
        dark: {
          ...require("daisyui/src/theming/themes")["black"],
          "success": "#28a745",
          "success-content": "white",
          "error": "#dc3545",
          "error-content": "white",
          "warning": "#ffc107",
          "warning-content": "black",
          "info": "#17a2b8",
          "info-content": "white",
        },
      },
    ],
  },
};
