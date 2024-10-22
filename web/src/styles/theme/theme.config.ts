import type { ThemeConfig } from "antd";
import { Poppins } from "next/font/google";

const PoppinsFont = Poppins({ weight: "400", subsets: ["latin"] });

const theme: ThemeConfig = {
	token: {
		fontSize: 16,
		colorPrimary: "#A5FFA3",
		fontFamily: PoppinsFont.style.fontFamily,
		colorText: "#000000",
	},
	components: {
		Layout: {
			headerBg: "#20222E",
			siderBg: "#20222E",
		},
	},
};

export default theme;
