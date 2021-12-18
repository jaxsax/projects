module.exports = {
  reactStrictMode: true,
  generatedBuildId: () => "0.0.1",
  webpack: (config) => {
    config.watchOptions = { poll: 300 };
    return config;
  },
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
