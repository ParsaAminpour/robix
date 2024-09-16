import React from "react";
import { CSSProperties } from "react";
import { BreakpointsNames, ContainerProps, ContainerStyle } from "./container.types";
import { BaseProps } from "@/types/global.types";

// Define the maxWidth for each breakpoint
const maxWidthMap: Record<BreakpointsNames, ContainerStyle> = {
	sm: {
		maxWidth: "600px",
		paddingLeft: "",
		paddingRight: "",
	},
	md: {
		maxWidth: "960px",
		paddingLeft: "",
		paddingRight: "",
	},
	lg: {
		maxWidth: "1280px",
		paddingLeft: "",
		paddingRight: "",
	},
	xl: {
		maxWidth: "1920px",
		paddingLeft: "66px",
		paddingRight: "24px",
	},
};

const Container: BaseProps<ContainerProps> = ({ children, maxWidth = "xl", fluid = false }) => {
	const containerStyle: CSSProperties = {
		maxWidth: fluid ? "100%" : maxWidthMap[maxWidth].maxWidth, // Set maxWidth or make full width for fluid
		margin: "0 auto", // Center the container
		paddingRight: fluid ? "0" : maxWidthMap[maxWidth].paddingRight,
		paddingLeft: fluid ? "0" : maxWidthMap[maxWidth].paddingLeft,
		width: "100%", // Full width on smaller screens
		boxSizing: "border-box", // Ensure padding doesn’t affect the width
	};

	return <div style={containerStyle}>{children}</div>;
};

export default Container;
