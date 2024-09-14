import { BaseProps } from "@/types/global.types";
import { BackpackWalletAdapter } from "@solana/wallet-adapter-backpack"; // Import Backpack
import { WalletProvider as SolanaWalletProvider } from "@solana/wallet-adapter-react";
import { PhantomWalletAdapter } from "@solana/wallet-adapter-wallets";
import { memo, useMemo } from "react";

const WalletProvider: BaseProps = ({ children, ...props }) => {
	const wallets = useMemo(() => [new PhantomWalletAdapter(), new BackpackWalletAdapter()], []);
	return (
		<SolanaWalletProvider
			wallets={wallets}
			autoConnect
			{...props}>
			{children}
		</SolanaWalletProvider>
	);
};
export default memo(WalletProvider);
