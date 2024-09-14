// import { YekanFont } from '@/styles/fonts/font.config'
import { TypographyOptions } from "@mui/material/styles/createTypography";

export const typography: TypographyOptions = {
	// fontFamily: [YekanFont.style.fontFamily].join(','),
	h1: {
		fontWeight: 700,
		fontStyle: "normal",
		fontSize: "4rem",
		lineHeight: `${4 * 1.6}rem`,
		color: "#334155",
	},
	h2: {
		fontWeight: 700,
		fontStyle: "normal",
		fontSize: "3rem",
		lineHeight: `${3 * 1.6}rem`,
		color: "#334155",
	},
	h3: {
		fontWeight: 400,
		fontStyle: "normal",
		fontSize: "1.75rem",
		lineHeight: `${1.75 * 1.6}rem`,
		color: "#334155",
	},
	body1: {
		fontWeight: 400,
		fontStyle: "normal",
		fontSize: "0.875rem",
		lineHeight: `${0.875 * 1.6}rem`,
		color: "#334155",
	},
	body2: {
		fontWeight: 400,
		fontStyle: "normal",
		fontSize: "1rem",
		lineHeight: `${1 * 1.6}rem`,
		color: "#334155",
	},
	button: {
		fontWeight: 400,
		fontStyle: "normal",
		fontSize: "0.75rem",
		lineHeight: `${0.75 * 1.6}rem`,
		textTransform: "none",
	},
	subtitle1: {
		fontWeight: 400,
		fontStyle: "normal",
		fontSize: "0.75rem",
		lineHeight: `${0.75 * 1.6}rem`,
	},
};
