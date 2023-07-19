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
    themes: ['dark', 'light'],
  },
};
