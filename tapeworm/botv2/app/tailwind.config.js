module.exports = {
  purge: ["./pages/**/*.{js,ts,jsx,tsx}", "./components/**/*.{js,ts,jsx,tsx}"],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {},
    flex: {
      4: "0 1 24%",
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
};
