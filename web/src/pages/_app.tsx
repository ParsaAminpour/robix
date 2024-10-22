import MainLayout from "@/components/layout/main.layout";
import ConnectionProvider from "@/providers/connection.provider";
import WalletProvider from "@/providers/wallet.provider";
import store from "@/store/store";
import "@/styles/globals.css";
import theme from "@/styles/theme/theme.config";
import { ConfigProvider } from "antd";
import { NextPage } from "next";
import type { AppProps } from "next/app";
import { ReactElement, ReactNode } from "react";
import { Provider } from "react-redux";

// Extend the NextPage type to include a layout property
type NextPageWithLayout = NextPage & {
	layout?: (page: ReactElement) => ReactNode;
};

// Update the AppProps type to use the extended NextPage type
type AppPropsWithLayout = AppProps & {
	Component: NextPageWithLayout;
};

export default function App(props: AppPropsWithLayout) {
	const { pageProps, Component } = props;
	const defaultLayout = (page: ReactElement): ReactNode => <MainLayout>{page}</MainLayout>;
	const getLayout = Component.layout ?? defaultLayout;

	return (
		<ConfigProvider theme={theme}>
			<ConnectionProvider>
				<WalletProvider>
					<Provider store={store}>{getLayout(<Component {...pageProps} />)}</Provider>
				</WalletProvider>
			</ConnectionProvider>
		</ConfigProvider>
	);
}
