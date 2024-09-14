import { PaletteOptions } from "@mui/material/styles";
declare module "@mui/material/styles" {
	interface Sizes {
		sideMenu: {
			width: number;
			itemHeight: number;
			logoHeight: number;
			paddingTop: number;
			marginBottom: number;
		};
		hex: number;
		hex_node: number;
		hex_event: number;
		hex_service: number;
		hex_structure: number;
		header: {
			iconTray: number;
			height: number;
		};
	}
	interface Palette {}
	interface PaletteOptions {}

	interface TypographyPalleteColorOptions {
		title: string;
		delete: string;
		secondary: string;
		success: string;
		info: string;
	}

	interface TypographyPalleteColorOptions {
		title: string;
		delete: string;
	}
}

export const paletteTheme: PaletteOptions = {
	primary: {
		main: "#0F8F92",
		"100": "#AFDADB",
		"200": "#87C7C8",
		"300": "#5FB4B6",
		"400": "#37A2A4",
		"500": "#0F8F92",
		"600": "#0D777A",
		"700": "#0A5F61",
		"800": "#084849",
		"900": "#053031",
	},
	secondary: {
		main: "#FA896B",
		"100": "#003e56",
		"200": "#004964",
		"300": "#005372",
		"400": "#005e81",
		"500": "#FA896B",
		"600": "#272A68",
		"700": "#1F2153",
		"800": "#18193F",
		"900": "#10112A",
	},
	warning: { main: "#FFA800" },
	success: { main: "#41C980" },
	info: { main: "#3395FF" },
	error: { main: "#EB5757" },
};
