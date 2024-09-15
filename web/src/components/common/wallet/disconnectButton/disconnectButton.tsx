import { useWallet } from "@solana/wallet-adapter-react";
import { Button } from "antd";

const DisconnectButton = () => {
	const { disconnect } = useWallet();
	return (
		<Button
			type="default"
			color="error"
			onClick={disconnect}>
			disconnect
		</Button>
	);
};

export default DisconnectButton;
