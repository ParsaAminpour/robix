import { WalletName } from "@solana/wallet-adapter-base";
import { useWallet } from "@solana/wallet-adapter-react";
import { Button, Divider, Modal, Row, Typography } from "antd";
import Image from "next/image";
import Link from "next/link";
import { useEffect } from "react";
import { IWalletModalProps } from "./walletModal.types";
const Title = Typography.Title;

const WalletModal = (props: IWalletModalProps) => {
	const { onClose, open } = props;
	const { wallets, select, connect, connecting, connected } = useWallet();
	const availableWallets = wallets.filter((wallet) => ["Phantom", "Backpack"].includes(wallet.adapter.name));
	const notDetectedWallets = wallets.filter((wallet) => wallet.readyState === "NotDetected");
	const handleSelectWallet = async (walletName: WalletName) => {
		select(walletName);
		await connect()
			.then(() => {
				// return notify({
				//   label: "Wallet Connected",
				//   message: "Connected to wallet successfully",
				//   type: "success",
				// });
			})
			.catch(() => {
				// return notify({
				//   label: "ERROR",
				//   message: error.message,
				//   type: "error",
				// });
			});
	};

	useEffect(() => {
		if (open && !connecting && connected) {
			onClose();
		}

		return () => {};
	}, [connecting, connected]);

	return (
		<Modal
			open={open}
			onClose={onClose}
			aria-labelledby="modal-modal-title"
			aria-describedby="modal-modal-description">
			<div
				style={{
					position: "absolute",
					top: "50%",
					left: "50%",
					transform: "translate(-50%, -50%)",
					width: 400,
					border: "2px solid #000",
					padding: 4,
				}}>
				<Row
					gutter={2}
					align={"top"}>
					{availableWallets.map((wallet) => {
						if (wallet.readyState !== "NotDetected") {
							return (
								<Button
									loading={wallet.adapter.connecting}
									disabled={connecting}
									icon={<></>}
									iconPosition="start"
									type="primary"
									color="primary"
									key={wallet.adapter.name}
									onClick={() => handleSelectWallet(wallet.adapter.name)}>
									<Row
										content={"space-between"}
										align={"middle"}>
										<Title>{wallet.adapter.name}</Title>
										<Image
											src={wallet.adapter.icon}
											alt={wallet.adapter.name}
											width={25}
											height={25}
										/>
									</Row>
								</Button>
							);
						} else {
							return <></>;
						}
					})}
					{notDetectedWallets.length ? (
						<div>
							<Title>Not Installed</Title>
							<Divider orientation="center" />
							{notDetectedWallets.map((wallet) => {
								return (
									<Row key={wallet.adapter.name}>
										<Link
											key={wallet.adapter.name}
											href={wallet.adapter.url}>
											<Row
												align={"middle"}
												content={"space-between"}>
												<Title>{wallet.adapter.name}</Title>
												<Image
													src={wallet.adapter.icon}
													alt={wallet.adapter.name}
													width={24}
													height={24}
												/>
											</Row>
										</Link>
									</Row>
								);
							})}
						</div>
					) : (
						<></>
					)}
				</Row>
			</div>
		</Modal>
	);
};

export default WalletModal;
