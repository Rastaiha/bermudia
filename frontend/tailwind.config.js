/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        'vazir': ['Vazirmatn', 'sans-serif'],
      },
      transitionTimingFunction: {
        'smooth-expand': 'cubic-bezier(0.25, 0.1, 0.25, 1.0)',
      },
      keyframes: {
        'boat-animation': {
          '0%, 100%': { transform: 'translate(0, 0) rotate(10deg) scale(0.3)' },
          '35%': { transform: 'translate(0.02px, 0.01px) rotate(-10deg) scale(0.3)' },
          '70%': { transform: 'translate(-0.02px, 0.01px) rotate(3deg) scale(0.3)' },
        }
      },
      animation: {
        'boat': 'boat-animation 10s linear infinite',
      }
    },
  },
  plugins: [],
}