import Container from "@/components/common/container/container";
import PageContainer from "@/components/common/pageContainer/pageContainer";
import ConnectButton from "@/components/common/wallet/connectButton/connectButton";
import DisconnectButton from "@/components/common/wallet/disconnectButton/disconnectButton";
import ShowBalance from "@/components/common/wallet/showBalance/showBalance";
import { Col, Row, Typography } from "antd";
import Head from "next/head";

const { Title } = Typography;

export default function Home() {
	return (
		<>
			<Head>
				<title>Robix</title>
				<meta
					name="description"
					content=""
				/>
				<meta
					name="viewport"
					content="width=device-width, initial-scale=1"
				/>
				<link
					rel="icon"
					href="/favicon.ico"
				/>
			</Head>

			<PageContainer>
				<Container maxWidth={"xl"}>
					<Title>Hello Robix</Title>
					<DisconnectButton />
					<ShowBalance />

					<Row gutter={[16, 8]}>
						<Col span={2}>
							<div>item </div>
						</Col>
						<Col span={2}>
							<div>item </div>
						</Col>
						<Col span={2}>
							<div>item </div>
						</Col>
						<Col span={2}>
							<div>item </div>
						</Col>
						<Col span={2}>
							<div>item </div>
						</Col>
						<Col span={2}>
							<div>item </div>
						</Col>
						<Col span={2}>
							<div>item </div>
						</Col>
						<Col span={2}>
							<div>item </div>
						</Col>
						<Col span={2}>
							<div>item </div>
						</Col>
						<Col span={2}>
							<div>item </div>
						</Col>
						<Col span={2}>
							<div>item </div>
						</Col>
						<Col span={2}>
							<div>item </div>
						</Col>
					</Row>
				</Container>
			</PageContainer>
		</>
	);
}
