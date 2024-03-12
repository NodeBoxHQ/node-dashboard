export default {
  content: ["./public/**/*.{html,js}", "./services/**/*.go"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Jost", "sans-serif"],
      },
      colors: {
        secondaryColor: "#f8c4cf",
        tertiaryColor: "#cecece",
        backgroundColor: "#1a1b1b",
        textColor: "#F4F4F4",
        nodeTypesColorActive: "#E53325",
        progressBarFillColor: "#E53325",
        progressBarBackgroundColor: "#707D75",
        cardBackgroundColor: "#444444",
        cardTitleColor: "#E53325",
        cardSubBodyColor: "#BDC0BA",
        inputAccentColor: "#707D75",
        glowColor: "#ffffffb3",
        themeSwitcherColor: "#E53325",
      },
    },
  },
  plugins: [],
};