import svelte from "rollup-plugin-svelte";
import resolve from "@rollup/plugin-node-resolve";
import commonjs from "@rollup/plugin-commonjs";
import livereload from "rollup-plugin-livereload";
import { terser } from "rollup-plugin-terser";
// import sveltePreprocess from 'svelte-preprocess'
import html from "@rollup/plugin-html";

const makeHtmlAttributes = (attributes) => {
  if (!attributes) {
    return "";
  }

  const keys = Object.keys(attributes);
  // eslint-disable-next-line no-param-reassign
  return keys.reduce(
    (result, key) => (result += ` ${key}="${attributes[key]}"`),
    ""
  );
};

const defaultTemplate = async ({ attributes, files, publicPath, title }) => {
  const scripts = (files.js || [])
    .map(({ fileName }) => {
      const attrs = makeHtmlAttributes(attributes.script);
      return `<script src="${publicPath}${fileName}"${attrs}></script>`;
    })
    .join("\n");

  const links = (files.css || [])
    .map(({ fileName }) => {
      const attrs = makeHtmlAttributes(attributes.link);
      return `<link href="${publicPath}${fileName}" rel="stylesheet"${attrs}>`;
    })
    .join("\n");

  return `
  <!doctype html>
  <html${makeHtmlAttributes(attributes.html)}>
	<head>
	  <meta charset="utf-8">
	  <meta name="viewport" content="width=device-width,initial-scale=1" />

	  <title>${title}</title>
	  <link rel="stylesheet" href="/static/global.css" />
	  <link href="https://unpkg.com/tailwindcss@^1.0/dist/tailwind.min.css" rel="stylesheet">
	  ${links}
	</head>
	<body>
	  ${scripts}
	</body>
  </html>`;
};

const mode = process.env.NODE_ENV;
const dev = mode === 'development';

let output = {};
if (!dev) {
  output = {
    dir: "dist",
    name: "dist/bundle.js",
    sourcemap: true,
    entryFileNames: "[name]-[hash].js",
    format: "umd",
  };
} else {
  output = {
    dir: "dist",
    name: "dist/bundle.js",
    sourcemap: true,
    format: "umd",
  }
}

export default {
  input: "src/main.js",
  output: output,
  plugins: [
    svelte({
      // enable run-time checks when not in production
      dev: dev,
      // preprocess: sveltePreprocess({ postcss: true }),
      // we'll extract any component CSS out into
      // a separate file - better for performance
    }),
    resolve({browser: true}),
    commonjs(),
    html({
      template: defaultTemplate,
    }),
    // In dev mode, call `npm run start` once
    // the bundle has been generated
    dev && serve(),

    // Watch the `public` directory and refresh the
    // browser on changes when not in production
    dev && livereload("dist"),

    // If we're building for production (npm run build
    // instead of npm run dev), minify
    !dev && terser({
      module: true,
    }),
  ],
  watch: {
    clearScreen: false,
  },
};

function serve() {
  let started = false;

  return {
    writeBundle() {
      if (!started) {
        started = true;

        require("child_process").spawn("npm", ["run", "start", "--", "--dev"], {
          stdio: ["ignore", "inherit", "inherit"],
          shell: true,
        });
      }
    },
  };
}
