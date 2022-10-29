/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  experimental: {
    appDir: true,
    // runtime: 'experimental-edge',
  },
  cleanDistDir: true,
  eslint: {
    dirs: ['**/**'],
    ignoreDuringBuilds: false,
  },
  typescript: {
    ignoreBuildErrors: false,
  },
  trailingSlash: false,
};

module.exports = nextConfig;
