import { components } from "@/styles/mui/themes/components";
import { paletteTheme } from "@/styles/mui/themes/palette";
import { typography } from "@/styles/mui/themes/typography";
import { BaseProps } from "@/types/global.types";
import createCache from "@emotion/cache";
import { createTheme, ThemeProvider as MaterialThemeProvider } from "@mui/material";
import { AppCacheProvider } from "@mui/material-nextjs/v13-pagesRouter";

const ThemeProvider: BaseProps = ({ children }) => {
	const theme = createTheme({
		direction: "ltr",
		palette: paletteTheme,
		components: components,
		typography: {
			...typography,
			allVariants: {
				textAlign: "start",
			},
		},
	});
	const emotionCache = createCache({
		key: "mui",
	});

	return (
		<AppCacheProvider emotionCache={emotionCache}>
			<MaterialThemeProvider theme={theme}>{children}</MaterialThemeProvider>
		</AppCacheProvider>
	);
};

export default ThemeProvider;
