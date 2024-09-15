import ConnectButton from "@/components/common/wallet/connectButton/connectButton";
import DisconnectButton from "@/components/common/wallet/disconnectButton/disconnectButton";
import ShowBalance from "@/components/common/wallet/showBalance/showBalance";
import WalletModal from "@/components/common/wallet/walletModal/walletModal";
import { triggerModal } from "@/store/slices/modal/modal.slice";
import { useDispatch, useSelector } from "@/store/store";
import { Typography } from "antd";
import Head from "next/head";

const { Title } = Typography;

export default function Home() {
	const dispatch = useDispatch();
	const { modals } = useSelector((state) => state.modal);
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
			<Title>Hello Robix</Title>
			<ConnectButton />
			<DisconnectButton />
			<ShowBalance />
			<WalletModal
				onClose={() => dispatch(triggerModal({ modal: "wallet", trigger: false }))}
				open={modals.wallet}
			/>
		</>
	);
}
