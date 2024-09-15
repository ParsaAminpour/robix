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
};

export default theme;
