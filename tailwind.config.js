/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./**/*.templ', './**/*.html', './**/.go'],
  theme: {
    extend: {
      colors: {
        background: '#fcfafd',
        primary: '#8628d2',
        secondary: '#c18cec',
        accent: '#a74df0',
        text: '#070409',
      },
    },
  },
  plugins: [require('@tailwindcss/typography'), require('@tailwindcss/forms')],
};
