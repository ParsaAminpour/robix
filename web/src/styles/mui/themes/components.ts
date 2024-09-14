import { Components, Theme } from "@mui/material";

export const components: Components<Omit<Theme, "components">> = {
	MuiButton: {
		styleOverrides: {
			sizeSmall: {
				borderRadius: "2px",
			},
			root: {
				borderRadius: "6px",
			},
			colorPrimary: {
				boxShadow: " 0px 4px 9px 0px #0F8F9226",
			},
			colorSecondary: {
				boxShadow: "0px 4px 9px 0px #FA896B26",
			},
			contained: {
				color: "#ffffff",
			},
		},
	},
	MuiAppBar: {
		styleOverrides: {
			root: {
				backgroundColor: "white",
				boxShadow: "none",
				borderBottom: "1.5px solid #D4D4D4",
				padding: "25px 15px",
			},
		},
	},
	MuiTextField: {
		styleOverrides: {
			root: {
				"> div": {
					backgroundColor: "#F9F9F9",
				},
				padding: "8px 5px",
				borderRadius: "6px",
				":before": {
					border: "0",
				},
				":hover:not(.Mui-disabled, .Mui-error):before": {
					border: "0",
				},
			},
		},
	},
	MuiChip: {
		styleOverrides: {
			root: {
				width: "fit-content",
				padding: "10px 25px 10px 25px",
			},
			filledPrimary: {
				backgroundColor: "#C6F1F2",
			},
			filledSecondary: {
				backgroundColor: "#FFEFEB",
			},
		},
	},
	MuiIconButton: {
		styleOverrides: {
			root: {
				backgroundColor: "white",
				boxShadow: " 0px 4px 25px 0px #292C7C14",
			},
		},
	},
	MuiAccordion: {
		styleOverrides: {
			root: {
				borderRadius: "5px",
				boxShadow: " 0px 6px 25px 0px #5555550F",
				"&.Mui-expanded": {
					boxShadow: "0px 4px 4px 0px #00000040",
				},
				"& .MuiAccordionSummary-root": {
					padding: "1.2rem 2rem",
					borderRadius: "5px 5px 0 0 ",
					border: "1px solid #ffffff",
				},
				"&.MuiAccordion-root:before": {
					backgroundColor: "unset",
				},
				"& .MuiAccordionSummary-root.Mui-expanded": {
					background: "#0F8F92",
				},
				"& .MuiAccordionSummary-content": {
					margin: 0,
					// flexGrow: 0,
					width: "fit-content",
				},
				"& .MuiAccordionSummary-root.Mui-expanded .MuiTypography-root": {
					color: "white",
				},
				"& .MuiAccordionSummary-expandIconWrapper": {
					width: "40px",
					height: "40px",
					display: "flex",
					justifyContent: "center",
					alignItems: "center",
					backgroundColor: "#0F8F92",
					borderRadius: "5px",
				},
				"& .MuiAccordionSummary-expandIconWrapper.Mui-expanded ": {
					width: "40px",
					height: "40px",
					display: "flex",
					justifyContent: "center",
					alignItems: "center",
					backgroundColor: "#fff",
				},
				"& .MuiAccordionSummary-expandIconWrapper.Mui-expanded svg": {
					fill: "#0F8F92",
				},
				"& .MuiAccordionSummary-expandIconWrapper svg": {
					fill: "white",
				},
				// '& .MuiAccordionSummary-expandIconWrapper': {},
			},
		},
	},
	MuiTab: {
		styleOverrides: {
			root: {
				transition: "all linear 300ms",
				backgroundColor: "#fff",
				".MuiTabs-flexContainer": {},
				"&.Mui-selected": {
					backgroundColor: "#0F8F92",
					color: "#fff",
				},
			},
		},
	},
	MuiInput: {
		styleOverrides: {
			root: {
				backgroundColor: "#F9F9F9",
				padding: "8px 5px",
				borderRadius: "6px",
				":before": {
					border: "0",
				},
				":hover:not(.Mui-disabled, .Mui-error):before": {
					border: "0",
				},
			},
			input: {
				padding: 0,
			},
		},
	},
};
