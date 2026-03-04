import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  trailingSlash: true,
  skipTrailingSlashRedirect: true,

  images: {
    unoptimized: true,
  },

  compress: true,

  assetPrefix: process.env.NODE_ENV === 'production' ? '' : '',

  experimental: {
    optimizePackageImports: [
      '@radix-ui/react-icons',
      '@tabler/icons-react',
      'lucide-react'
    ],
  },

  async rewrites() {
    const apiHost = process.env.NEXT_PUBLIC_API_BASEURL?.endsWith('/') ? process.env.NEXT_PUBLIC_API_BASEURL.slice(0, -1) : process.env.NEXT_PUBLIC_API_BASEURL || ''
    return [
      {
        source: "/api/:path*",
        destination: `${apiHost}/api/:path*`
      }
    ]
  }
};

export default nextConfig;
