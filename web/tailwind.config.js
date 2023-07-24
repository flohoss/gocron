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
    themes: ["wireframe", "business"],
    darkTheme: "business"
  },
};
