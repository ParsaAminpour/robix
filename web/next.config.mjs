/** @type {import('next').NextConfig} */
const nextConfig = {
	reactStrictMode: true,
	env: {
		MODE: process.env.NODE_ENV,
	},
	compress: true,
	transpilePackages: ["antd", "@ant-design/plots", "@ant-design/icons"],
};

export default nextConfig;
