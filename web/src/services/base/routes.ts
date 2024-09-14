export const getApiRoutes = () => {
	const routes = {
		example: {
			slug: "/example",
			get: function (slug: string) {
				return `/${this.slug}/${slug}`;
			},
		},
	};
	return routes;
};
