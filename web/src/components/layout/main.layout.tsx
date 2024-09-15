import { BaseProps } from "@/types/global.types";
import { Flex, Layout } from "antd";

const { Header, Footer, Sider, Content } = Layout;

const headerStyle: React.CSSProperties = {
	textAlign: "center",
	color: "#fff",
	paddingInline: 48,
	lineHeight: "64px",
	backgroundColor: "blue",
};

const contentStyle: React.CSSProperties = {
	textAlign: "center",
	lineHeight: "120px",
	flexGrow: 1,
	width: "100%",
	height: "100%",
	color: "#fff",
};

const siderStyle: React.CSSProperties = {
	textAlign: "center",
	lineHeight: "120px",
	color: "#fff",
	height: "100%",
	backgroundColor: "green",
};

const footerStyle: React.CSSProperties = {
	textAlign: "center",
	color: "#fff",
	backgroundColor: "red",
	width: "100%",
	height: "78px",
};

const layoutStyle = {
	borderRadius: 8,
	height: "100svh",
};

const MainLayout: BaseProps = ({ children }) => {
	return (
		<Layout style={layoutStyle}>
			<Header style={headerStyle}>Header</Header>
			<Layout>
				<Sider
					width="100px"
					style={siderStyle}>
					Sider
				</Sider>
				<Flex
					vertical
					style={{ width: "100%" }}>
					<Content style={contentStyle}>{children}</Content>
					<Footer style={footerStyle}>Footer</Footer>
				</Flex>
			</Layout>
		</Layout>
	);
};

export default MainLayout;
