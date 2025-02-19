/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{html,js}"
  ],
  theme: {
    extend: {
      colors: {
        'brand-primary': '#3B82F6',
        'brand-secondary': '#10B981',
        'brand-background': '#F3F4F6'
      },
      fontFamily: {
        'sans': ['Inter', 'ui-sans-serif', 'system-ui']
      }
    },
  },
  plugins: [],
}
