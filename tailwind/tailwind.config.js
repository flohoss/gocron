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
          secondary: "#232332",
          "base-100": "#F3F4F6",
          neutral: "#A3A3B3",
          "neutral-content": "#1A1A28",
        },
      },
      {
        dark: {
          ...require("daisyui/src/theming/themes")["dark"],
          primary: "#75CEF9",
          secondary: "#F3F4F6",
          "base-100": "#1A1A28",
          neutral: "#A3A3B3",
          "neutral-content": "#1A1A28",
        },
      },
    ],
  },
}

