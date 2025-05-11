import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  images: {
    remotePatterns: [
      { hostname: "placebear.com" },
      { hostname: "picsum.photos" },
      { hostname: "fastly.picsum.photos" },
      { hostname: "placekitten.com" },
      { hostname: "loremflickr.com" },
    ],
    formats: [],
  },
};

export default nextConfig;
