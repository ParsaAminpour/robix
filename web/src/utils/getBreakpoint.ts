import { BreakpointsNames } from "@/components/common/container/container.types";

export const getBreakpoint = (width: number): BreakpointsNames => {
	if (width < 600) return "sm";
	if (width < 960) return "md";
	if (width < 1280) return "lg";
	return "xl";
};
