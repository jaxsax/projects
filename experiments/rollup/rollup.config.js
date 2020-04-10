import resolve from "@rollup/plugin-node-resolve";
import commonjs from "@rollup/plugin-commonjs";
import html from '@rollup/plugin-html';

export default {
  output: {
    entryFileNames: '[name]-[hash].js',
    format: "es",
  },
  plugins: [
    resolve(),
    commonjs(),
    html(),
  ],
};
