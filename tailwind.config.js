/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./internal/components/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}

