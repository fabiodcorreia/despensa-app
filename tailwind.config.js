/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["internal/views/**/*.templ"],
  theme: {
    extend: {},
  },
  daisyui: {
    darkTheme: "rosepine-moon",
    themes: [
      {
        'rosepine-dawn': {
          "primary": "#907aa9",
          "secondary": "#d7827e",
          "accent": "#907aa9",
          "neutral": "#faf4ed",
          "base-100": "#fffaf3",
          "info": "#286983",
          "success": "#56949f",
          "warning": "#ea9d34",
          "error": "#b4637a",
        },
        'rosepine-moon': {
          "primary": "#c4a7e7",
          "secondary": "#ea9a97",
          "accent": "#c4a7e7",
          "neutral": "#2a273f",
          "base-100": "#232136",
          "info": "#3e8fb0",
          "success": "#9ccfd8",
          "warning": "#f6c177",
          "error": "#eb6f92",
        },
      },
    ],
  },
  plugins: [
    require("daisyui"), 
    require("@tailwindcss/typography")
  ],
}

