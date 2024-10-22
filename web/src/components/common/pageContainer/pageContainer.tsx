import { BaseProps } from "@/types/global.types";
import React from "react";

const PageContainer: BaseProps = ({ children }) => {
	return (
		<div
			style={{
				paddingTop: "53px",
				paddingBottom: "59px",
			}}>
			{children}
		</div>
	);
};

export default PageContainer;
