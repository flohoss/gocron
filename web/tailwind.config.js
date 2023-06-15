/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.html"],
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
          ...require("daisyui/src/colors/themes")["[data-theme=light]"],
          primary: "#14468c",
          secondary: "#EB912D",
          warning: "#EB912D",
        },
      },
      {
        dark: {
          ...require("daisyui/src/colors/themes")["[data-theme=halloween]"],
          primary: "#EB912D",
          secondary: "#14468c",
          warning: "#EB912D",
        },
      },
    ],
    darkTheme: "dark",
    utils: true,
  },
};
