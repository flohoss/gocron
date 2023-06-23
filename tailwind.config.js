/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/templates/**/*.html"],
  theme: {
    container: {
      center: true
    },
    extend: {
      maxWidth: {
        'xxs': '15rem',
      }
    }
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: [
      {
        light: {
          ...require("daisyui/src/theming/themes")["[data-theme=light]"],
          primary: "#14468c",
          secondary: "#EB912D",
          warning: "#EB912D",
          "primary-content": "#FFFFFF",
        },
      },
      {
        dark: {
          ...require("daisyui/src/theming/themes")["[data-theme=dark]"],
          primary: "#EB912D",
          secondary: "#14468c",
          warning: "#EB912D",
          "primary-content": "#000000",
        },
      },
    ],
  },
};
