const pkg = require("./package.json");

module.exports = {
  reactStrictMode: true,
  generatedBuildId: () => pkg.version,
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://localhost:8080/api/:path*",
      },
      {
        source: "/images/:imageid*",
        destination: "http://localhost:8080/images/:imageid*",
      },
    ];
  },
};
