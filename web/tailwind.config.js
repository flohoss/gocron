/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{js,ts,jsx,tsx,mdx}'],
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
