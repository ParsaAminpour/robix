import { BaseProps } from "@/types/global.types";
import { ConnectionProvider as SolanaConnectionProvider } from "@solana/wallet-adapter-react";
import { clusterApiUrl } from "@solana/web3.js";

const ConnectionProvider: BaseProps = ({ children }) => {
	const Mode = process.env.MODE || "development";

	const endpoint =
		Mode === "development"
			? "https://devnet.helius-rpc.com/?api-key=48b598a7-1ea6-4667-8717-5dd3c5b31ed4"
			: clusterApiUrl("mainnet-beta");

	return <SolanaConnectionProvider endpoint={endpoint}>{children}</SolanaConnectionProvider>;
};
export default ConnectionProvider;
