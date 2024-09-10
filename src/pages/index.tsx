import { Button } from "@mui/material";
import Head from "next/head";

export default function Home() {
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

			<Button
				variant="text"
				color="primary">
				Text
			</Button>
			<Button
				variant="contained"
				color="primary">
				Contained
			</Button>
			<Button
				variant="outlined"
				color="primary">
				Outlined
			</Button>
		</>
	);
}
