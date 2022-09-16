const colors = require('tailwindcss/colors')

// tailwind.config.js
module.exports = {
  content: [
    './src/**/*.tsx',
  ],
  purge: ["./src/**/*.{js,jsx,ts,tsx}", "./public/index.html"],
  darkMode: false, // or 'media' or 'class'
  theme: {
    colors: {
      primary: colors.blue,
      gray: colors.gray,
      red: colors.red,
      green: colors.green,
      white: colors.white,
    },
    extend: {},
  },
  variants: {
    extend: {},
  },
  plugins: [],
};
