import { LoadingButton } from "@mui/lab";
import { Box, Divider, Modal, Stack, Typography } from "@mui/material";
import { WalletName } from "@solana/wallet-adapter-base";
import { useWallet } from "@solana/wallet-adapter-react";
import Image from "next/image";
import Link from "next/link";
import { useEffect } from "react";
import { IWalletModalProps } from "./walletModal.types";

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
			<Box
				sx={{
					position: "absolute",
					top: "50%",
					left: "50%",
					transform: "translate(-50%, -50%)",
					width: 400,
					bgcolor: "background.paper",
					border: "2px solid #000",
					boxShadow: 24,
					p: 4,
				}}>
				<Stack
					gap={"10px"}
					width={"100%"}
					alignItems={"start"}>
					{availableWallets.map((wallet) => {
						if (wallet.readyState !== "NotDetected") {
							return (
								<LoadingButton
									loading={wallet.adapter.connecting}
									loadingPosition="start"
									disabled={connecting}
									startIcon={<></>}
									variant="contained"
									fullWidth
									color="primary"
									key={wallet.adapter.name}
									onClick={() => handleSelectWallet(wallet.adapter.name)}>
									<Stack
										direction={"row"}
										width={"100%"}
										justifyContent={"space-between"}
										alignItems={"center"}>
										<Typography>{wallet.adapter.name}</Typography>
										<Image
											src={wallet.adapter.icon}
											alt={wallet.adapter.name}
											width={25}
											height={25}
										/>
									</Stack>
								</LoadingButton>
							);
						} else {
							return <></>;
						}
					})}
					{notDetectedWallets.length ? (
						<Box>
							<Typography>Not Installed</Typography>
							<Divider />
							{notDetectedWallets.map((wallet) => {
								return (
									<Stack key={wallet.adapter.name}>
										<Link
											key={wallet.adapter.name}
											href={wallet.adapter.url}>
											<Stack
												direction={"row"}
												alignItems={"center"}
												justifyContent={"space-between"}>
												<Typography>{wallet.adapter.name}</Typography>
												<Image
													src={wallet.adapter.icon}
													alt={wallet.adapter.name}
													width={24}
													height={24}
												/>
											</Stack>
										</Link>
									</Stack>
								);
							})}
						</Box>
					) : (
						<></>
					)}
				</Stack>
			</Box>
		</Modal>
	);
};

export default WalletModal;
