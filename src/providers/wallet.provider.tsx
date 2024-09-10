import { BaseProps } from "@/types/global.types";
import { WalletProvider as SolanaWalletProvider } from "@solana/wallet-adapter-react";
import { PhantomWalletAdapter } from "@solana/wallet-adapter-wallets";

import { useMemo } from "react";

const WalletProvider: BaseProps = ({ children, ...props }) => {
	const wallets = useMemo(
		() => [
			new PhantomWalletAdapter(),
			//  new SolflareWalletAdapter({ network: WalletAdapterNetwork.Mainnet })
		],
		[],
	);
	return (
		<SolanaWalletProvider
			wallets={wallets}
			autoConnect
			{...props}>
			{children}
		</SolanaWalletProvider>
	);
};
export default WalletProvider;
