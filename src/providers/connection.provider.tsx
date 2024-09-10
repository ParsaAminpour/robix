import { BaseProps } from "@/types/global.types";
import { ConnectionProvider as SolanaConnectionProvider } from "@solana/wallet-adapter-react";
import { clusterApiUrl } from "@solana/web3.js";

const ConnectionProvider: BaseProps = ({ children }) => {
	const Mode = process.env.MODE || "development";

	const endpoint = Mode === "development" ? clusterApiUrl("devnet") : clusterApiUrl("mainnet-beta");

	return <SolanaConnectionProvider endpoint={endpoint}>{children}</SolanaConnectionProvider>;
};
export default ConnectionProvider;
