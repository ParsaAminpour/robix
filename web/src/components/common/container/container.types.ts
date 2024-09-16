export type Breakpoints = "600px" | "960px" | "1280px" | "1920px";
export type BreakpointsNames = "sm" | "md" | "lg" | "xl";

export type ContainerStyle = {
	maxWidth: Breakpoints;
	paddingLeft: string;
	paddingRight: string;
};

export type ContainerProps = {
	children: React.ReactNode;
	maxWidth?: BreakpointsNames;
	fluid?: boolean;
};
