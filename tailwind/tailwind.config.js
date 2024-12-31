const { fontFamily } = require('tailwindcss/defaultTheme');

/** @type {import('tailwindcss').Config} */

module.exports = {
  content: ["../views/**/*.{templ,go}"],
  plugins: [
    require("daisyui")
  ],
  daisyui: {
    themes: ["dark"]
  }
}

