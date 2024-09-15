import { triggerModal } from "@/store/slices/modal/modal.slice";
import { useDispatch } from "@/store/store";
import { useWallet } from "@solana/wallet-adapter-react";
import { Button, Typography } from "antd";
import Image from "next/image";
const { Text } = Typography;

const ConnectButton = () => {
	const { connected, wallet } = useWallet();
	const dispatch = useDispatch();
	return (
		<>
			{connected && wallet ? (
				<Button
					type="primary"
					icon={
						<Image
							src={wallet.adapter.icon}
							alt="wallet"
							width={24}
							height={24}
						/>
					}
					iconPosition="start"
					color="primary">
					{wallet.adapter.publicKey?.toString()}
				</Button>
			) : (
				<Button
					onClick={() => dispatch(triggerModal({ modal: "wallet", trigger: true }))}
					type="primary">
					<Text>Connect & play</Text>
				</Button>
			)}
		</>
	);
};

export default ConnectButton;
