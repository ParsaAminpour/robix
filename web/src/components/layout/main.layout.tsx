import ModalProvider from "@/providers/modal.provider";
import { BaseProps } from "@/types/global.types";
import { Flex, Layout } from "antd";
import Header from "./header/header";

const { Header: AntHeader, Footer, Sider, Content } = Layout;

const HEADER_HEIGHT = "64px";
const SIDER_WIDTH = "76px";
const FOOTER_HEIGHT = "78px";

const headerStyle: React.CSSProperties = {
	borderBottom: "1px solid #3E404C",
	height: HEADER_HEIGHT,
	padding: 0,
};

const contentStyle: React.CSSProperties = {
	flexGrow: 1,
};

const siderStyle: React.CSSProperties = {
	borderRight: "1px solid #3E404C",
};

const footerStyle: React.CSSProperties = {
	background: "#20222E66",
	width: "100%",
	height: FOOTER_HEIGHT,
	borderTop: "1px solid #3E404C",
};

const layoutStyle: React.CSSProperties = {
	borderRadius: 8,
	minHeight: "100svh",
	height: "100%",
	width: "100%",
};

const MainLayout: BaseProps = ({ children }) => {
	return (
		<Layout style={layoutStyle}>
			<AntHeader style={headerStyle}>
				<Header />
			</AntHeader>
			<Layout>
				<Sider
					width={SIDER_WIDTH}
					style={siderStyle}>
					Sider
				</Sider>
				<Flex
					vertical
					style={{
						width: "100%",
						// background: `url(/assets/images/layout-bg.png)`,
						background: `#20222E`,
						backgroundSize: "cover",
					}}>
					<Content style={contentStyle}>{children}</Content>
					<Footer style={footerStyle}>Footer</Footer>
				</Flex>
			</Layout>
			<ModalProvider />
		</Layout>
	);
};

export default MainLayout;
