module.exports = {
  reactStrictMode: true,
  generatedBuildId: () => "0.0.1",
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://localhost:8081/api/:path*",
      },
    ];
  },
  webpack: (config) => {
    config.watchOptions = { poll: 300 };
    return config;
  },
  future: {
    webpack5: true,
  },
};
