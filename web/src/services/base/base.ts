import axios, { AxiosInstance } from "axios";
import { deleteCookie, getCookie } from "cookies-next";
import { IErrorResponse, IRequestOption, IResponse } from "./request-interface";

export function successHandler<T>(response: IResponse<T>): T {
	return response.data;
}

export function errorHandler(error: IErrorResponse): void {
	throw error;
}

export async function sendRequest<T, D = undefined>({
	...restOptions
}: IRequestOption<D>): Promise<IResponse<T> | IErrorResponse> {
	const axiosInstance: AxiosInstance = axios.create({ baseURL: process.env.BASE_URL });

	axiosInstance.interceptors.request.use((config) => {
		const access_token = getCookie("token");
		if (access_token) {
			config.headers.set({
				authorization: `Bearer ${access_token}`,
				"Accept-Language": "fa",
				"Content-Type": "application/json",
				Accept: "application/json",
				...config.headers,
			});
		} else {
			config.headers.set({
				"Accept-Language": "fa",
				"Content-Type": "application/json",
				Accept: "application/json",
				...config.headers,
			});
		}

		return config;
	});

	axiosInstance.interceptors.response.use(
		(res) => {
			return res;
		},
		(error) => {
			if (error.response.status === 401) {
				// const domain = getDomain();
				deleteCookie("token");
			}
			return Promise.reject(error);
		},
	);

	try {
		const response: IResponse<T> = await axiosInstance({ ...restOptions });
		successHandler<T>(response);
		return response;
	} catch (error) {
		errorHandler(error as IErrorResponse);
		return error as IErrorResponse;
	}
}
