import { triggerModal } from "@/store/slices/modal/modal.slice";
import { useDispatch } from "@/store/store";
import { Button } from "@mui/material";
import { useWallet } from "@solana/wallet-adapter-react";
import Image from "next/image";

const ConnectButton = () => {
	const { connected, wallet } = useWallet();
	const dispatch = useDispatch();
	return (
		<>
			{connected && wallet ? (
				<Button
					variant="contained"
					startIcon={
						<Image
							src={wallet.adapter.icon}
							alt="wallet"
							width={24}
							height={24}
						/>
					}
					color="primary">
					{wallet.adapter.publicKey?.toString()}
				</Button>
			) : (
				<Button
					onClick={() => dispatch(triggerModal({ modal: "wallet", trigger: true }))}
					variant="contained"
					color="primary">
					Connect Wallet
				</Button>
			)}
		</>
	);
};

export default ConnectButton;
