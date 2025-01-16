/** @type {import('tailwindcss').Config} */

module.exports = {
  content: ["views/**/*.{templ,go}"],
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
    utils: true,
    themes: [
      {
        light: {
          ...require("daisyui/src/theming/themes")["light"],
          primary: "#26a4e2",
          secondary: "#232332",
          "base-100": "#F3F4F6",
          neutral: "#E5E6E9",
          "neutral-content": "#232332",
        },
      },
      {
        dark: {
          ...require("daisyui/src/theming/themes")["dark"],
          primary: "#75CEF9",
          secondary: "#F3F4F6",
          "base-100": "#1A1A28",
          neutral: "#2E2E40",
          "neutral-content": "#F3F4F6",
        },
      },
    ],
  },
}

