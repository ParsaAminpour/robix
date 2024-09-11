/** @type {import('next').NextConfig} */
const nextConfig = {
	reactStrictMode: true,
	env: {
		MODE: process.env.NODE_ENV,
	},
};

export default nextConfig;
