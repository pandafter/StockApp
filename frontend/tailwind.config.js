/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        market: {
          bg: '#0f172a',
          surface: '#1e293b',
          border: '#334155',
        },
      },
    },
  },
  plugins: [],
}
