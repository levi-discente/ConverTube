import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  experimental: {
    serverActions: {
      bodySizeLimit: '90mb', // Ajuste conforme necessário
    },
  },
};

export default nextConfig;
