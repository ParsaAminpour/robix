import { Stack, Typography } from "@mui/material";
import { useConnection, useWallet } from "@solana/wallet-adapter-react";
import { LAMPORTS_PER_SOL } from "@solana/web3.js";
import { memo, useEffect, useState } from "react";

const ShowBalance = () => {
	const { publicKey } = useWallet();
	const { connection } = useConnection();
	const [balance, setBalance] = useState<number>(0);

	useEffect(() => {
		if (publicKey) {
			connection
				.getBalance(publicKey)
				.then((res) => {
					setBalance(res / LAMPORTS_PER_SOL);
				})
				.catch((err) => console.log(err));
		}
	}, [publicKey, connection]);
	return (
		<Stack
			direction={"row"}
			alignItems={"center"}
			gap={"10px"}>
			<Typography variant="h6">Balance: </Typography>
			<Typography variant="body1">{balance}</Typography>
		</Stack>
	);
};

export default memo(ShowBalance);
