const { fontFamily } = require('tailwindcss/defaultTheme');

/** @type {import('tailwindcss').Config} */

module.exports = {
  content: ["../views/**/*.{templ,go}"],
  theme: {
    container: {
      center: true,
      padding: {
        DEFAULT: '1rem',
        sm: '2rem',
        lg: '4rem',
        xl: '5rem',
        '2xl': '6rem',
      },
    },
  },
  plugins: [
    require("daisyui")
  ],
  daisyui: {
    themes: [
      {
        light: {
          ...require("daisyui/src/theming/themes")["light"],
          primary: "#75CEF9",
          secondary: "#FAFAFA",
        },
      },
      {
        dark: {
          ...require("daisyui/src/theming/themes")["dark"],
          primary: "#75CEF9",
          secondary: "#FAFAFA",
          "base-100": "#1A1A28",
        },
      },
    ],
  },
}

