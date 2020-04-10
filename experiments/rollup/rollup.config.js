import resolve from "@rollup/plugin-node-resolve";
import commonjs from "@rollup/plugin-commonjs";

export default {
  input: "src/main.js",
  output: {
    // file: "dist/bundle.js",
    dir: "dist",
    format: "cjs",
  },
  plugins: [
    resolve(),
    commonjs(),
  ],
};
