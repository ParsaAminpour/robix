import { sendRequest } from "../base";
import { IErrorResponse, IResponse } from "../request-interface";
import { IGetRequestOption } from "./get-request-interface";

export default async function getRequest<T>(options: IGetRequestOption): Promise<IResponse<T> | IErrorResponse> {
	return sendRequest<T>({ method: "GET", ...options });
}
