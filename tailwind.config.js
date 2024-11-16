/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "templates/**/*.{html,js}"
  ],
  theme: {
    extend: {
      colors: {
        primary: '#FF6363',
        secondary: {
          100: '#E2E2D5',
          200: '#888888',
        },
      },
      fontFamily: {
        'rubik': ['Rubik'],
      },
    },
  },
  plugins: [],
}

