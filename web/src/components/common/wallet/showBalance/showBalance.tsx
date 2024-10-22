import { useConnection, useWallet } from "@solana/wallet-adapter-react";
import { LAMPORTS_PER_SOL } from "@solana/web3.js";
import { Flex, Typography } from "antd";
import { memo, useEffect, useState } from "react";
const { Title, Text } = Typography;

const ShowBalance = () => {
	const { publicKey } = useWallet();
	const { connection } = useConnection();
	const [balance, setBalance] = useState<number>(0);
	useEffect(() => {
		if (!publicKey) {
			setBalance(0);
			return;
		}
		// Function to fetch and update the balance
		const fetchBalance = async () => {
			try {
				const balanceInLamports = await connection.getBalance(publicKey);
				setBalance(balanceInLamports / LAMPORTS_PER_SOL);
			} catch (err) {
				console.log(err);
			}
		};
		fetchBalance();
		const intervalId = setInterval(fetchBalance, 10000); // 10 seconds
		return () => clearInterval(intervalId);
	}, [publicKey, connection]);

	return (
		<Flex>
			<Title level={2}>Balance: </Title>
			<Text>{balance}</Text>
		</Flex>
	);
};

export default memo(ShowBalance);
