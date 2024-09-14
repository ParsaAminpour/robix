import { Button } from "@mui/material";
import { useWallet } from "@solana/wallet-adapter-react";

const DisconnectButton = () => {
	const { disconnect } = useWallet();
	return (
		<Button
			variant="outlined"
			color="error"
			onClick={disconnect}>
			disconnect
		</Button>
	);
};

export default DisconnectButton;
