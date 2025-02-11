/** @type {import('tailwindcss').Config} */

module.exports = {
  content: ["./internals/components/*.templ", "./**/*.html", "./**/*.templ", "./**/*.go","!./node_modules/**/*.html",],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}

