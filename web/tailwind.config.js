/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.{html,js}"],
  theme: {
    extend: {
      colors: {
        'primary': '#540333',
        'secondary': '#e73636',
        'darkest': '#331426',
        'dark': '#735064',
        'medium': '#a8879a',
        'light': '#f2dae8',
        'lightest': '#faf5f8'
      },
      fontFamily: {
        'roboto': ['Roboto', 'sans-serif']
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}

