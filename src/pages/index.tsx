import Head from "next/head";

export default function Home() {
	// const { wallets } = useWallet();
	// const availableWallets = wallets.filter((wallet) => ["Phantom", "Backpack"].includes(wallet.adapter.name));
	// console.log(availableWallets);
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
			<h1>Hello Robix</h1>
		</>
	);
}
