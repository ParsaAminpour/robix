import ThemeProvider from "@/providers/theme.provider";
import "@/styles/globals.css";
import { NextPage } from "next";
import type { AppProps } from "next/app";
import { ReactElement, ReactNode } from "react";

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
	const defaultLayout = (page: ReactElement): ReactNode => <div>{page}</div>;
	const getLayout = Component.layout ?? defaultLayout;

	return (
		<>
			<ThemeProvider>{getLayout(<Component {...pageProps} />)}</ThemeProvider>
		</>
	);
}
